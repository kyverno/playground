package ivpol

import (
	"context"

	"github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"github.com/kyverno/kyverno/pkg/cel/libs"
	"github.com/kyverno/kyverno/pkg/cel/matching"
	"github.com/kyverno/kyverno/pkg/cel/policies/ivpol/engine"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	eval "github.com/kyverno/kyverno/pkg/imageverification/evaluator"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/kyverno/playground/backend/pkg/engine/models"
	"github.com/kyverno/playground/backend/pkg/engine/utils"
)

func Process(ctx context.Context, dClient dclient.Interface, restMapper meta.RESTMapper, contextProvider libs.Context, params *models.Parameters, newResource, oldResource unstructured.Unstructured, ivpols []v1alpha1.ImageValidatingPolicy) ([]models.Response, error) {
	request := utils.NewCELRequest(restMapper, contextProvider, params, newResource, oldResource)
	validations := make([]models.Response, 0)

	if request.JsonPayload == nil {
		engine, err := newIVPEngine(dClient, ivpols, nil)
		if err != nil {
			return nil, err
		}

		resp, _, err := engine.HandleMutating(ctx, request, nil)
		if err != nil {
			return nil, err
		}

		for _, result := range resp.Policies {
			resp := engineapi.NewEngineResponse(newResource, engineapi.NewImageValidatingPolicy(result.Policy), params.Context.NamespaceLabels).
				WithPolicyResponse(engineapi.PolicyResponse{Rules: []engineapi.RuleResponse{result.Result}})

			validations = append(validations, models.ConvertResponse(resp))
		}

		return validations, nil
	}

	compiled := make([]*eval.CompiledImageValidatingPolicy, 0)
	pMap := make(map[string]*v1alpha1.ImageValidatingPolicy)
	for i := range ivpols {
		p := ivpols[i]
		pMap[p.GetName()] = &p
		compiled = append(compiled, &eval.CompiledImageValidatingPolicy{Policy: &p})
	}

	results, err := eval.Evaluate(ctx, compiled, newResource.Object, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	resp := engineapi.EngineResponse{
		Resource:       newResource,
		PolicyResponse: engineapi.PolicyResponse{},
	}
	for p, rslt := range results {
		if rslt.Error != nil {
			resp.PolicyResponse.Rules = []engineapi.RuleResponse{
				*engineapi.RuleError("evaluation", engineapi.ImageVerify, "failed to evaluate policy for JSON", rslt.Error, nil),
			}
		} else if rslt.Result {
			resp.PolicyResponse.Rules = []engineapi.RuleResponse{
				*engineapi.RulePass(p, engineapi.ImageVerify, "success", nil),
			}
		} else {
			resp.PolicyResponse.Rules = []engineapi.RuleResponse{
				*engineapi.RuleFail(p, engineapi.ImageVerify, rslt.Message, nil),
			}
		}
		resp = resp.WithPolicy(engineapi.NewImageValidatingPolicy(pMap[p]))

		validations = append(validations, models.ConvertResponse(resp))
	}

	return validations, nil
}

func newIVPEngine(dClient dclient.Interface, policies []v1alpha1.ImageValidatingPolicy, exceptions []*v1alpha1.PolicyException) (engine.Engine, error) {
	provider, err := engine.NewProvider(policies, exceptions)
	if err != nil {
		return nil, err
	}
	return engine.NewEngine(
		provider,
		utils.NSResolver(dClient),
		matching.NewMatcher(),
		dClient.GetKubeClient().CoreV1().Secrets(""),
		nil,
	), nil
}
