package mpol

import (
	"context"
	"fmt"

	"github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"github.com/kyverno/kyverno/pkg/cel/libs"
	"github.com/kyverno/kyverno/pkg/cel/matching"
	"github.com/kyverno/kyverno/pkg/cel/policies/mpol/compiler"
	"github.com/kyverno/kyverno/pkg/cel/policies/mpol/engine"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/kyverno/playground/backend/pkg/engine/models"
	"github.com/kyverno/playground/backend/pkg/engine/utils"
)

func Process(ctx context.Context, dClient dclient.Interface, restMapper meta.RESTMapper, contextProvider libs.Context, params *models.Parameters, resource, oldResource unstructured.Unstructured, mpols []v1alpha1.MutatingPolicy) ([]models.Response, error) {
	provider, err := NewProvider(compiler.NewCompiler(), mpols, nil)
	if err != nil {
		return nil, err
	}

	eng := engine.NewEngine(provider, utils.NSResolver(dClient), dClient.GetKubeClient(), matching.NewMatcher())

	request := utils.NewCELRequest(restMapper, contextProvider, params, resource, oldResource)
	results := make([]models.Response, 0)

	engineResponse, err := eng.Handle(ctx, request, nil)
	if err != nil {
		return nil, err
	}

	for _, res := range engineResponse.Policies {
		response := engineapi.EngineResponse{
			Resource: *engineResponse.Resource,
			PolicyResponse: engineapi.PolicyResponse{
				Rules: res.Rules,
			},
		}

		results = append(results, models.ConvertResponse(response.WithPolicy(engineapi.NewMutatingPolicy(res.Policy))))
	}

	return results, nil
}

func NewProvider(
	compiler compiler.Compiler,
	policies []v1alpha1.MutatingPolicy,
	exceptions []*v1alpha1.PolicyException,
) (engine.ProviderFunc, error) {
	out := make([]engine.Policy, 0, len(policies))
	for _, policy := range policies {
		var matchedExceptions []*v1alpha1.PolicyException
		for _, polex := range exceptions {
			for _, ref := range polex.Spec.PolicyRefs {
				if ref.Name == policy.GetName() && ref.Kind == policy.GetKind() {
					matchedExceptions = append(matchedExceptions, polex)
				}
			}
		}
		compiled, errs := compiler.Compile(&policy, matchedExceptions)
		if len(errs) > 0 {
			return nil, fmt.Errorf("failed to compile policy %s (%w)", policy.GetName(), errs.ToAggregate())
		}
		out = append(out, engine.Policy{
			Policy:         policy,
			CompiledPolicy: compiled,
		})
	}
	return func(context.Context) ([]engine.Policy, error) {
		return out, nil
	}, nil
}
