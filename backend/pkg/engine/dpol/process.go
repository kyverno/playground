package dpol

import (
	"context"

	"github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"github.com/kyverno/kyverno/pkg/cel/libs"
	"github.com/kyverno/kyverno/pkg/cel/matching"
	"github.com/kyverno/kyverno/pkg/cel/policies/dpol/compiler"
	"github.com/kyverno/kyverno/pkg/cel/policies/dpol/engine"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/kyverno/playground/backend/pkg/engine/models"
	"github.com/kyverno/playground/backend/pkg/engine/utils"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Process(ctx context.Context, dClient dclient.Interface, restMapper meta.RESTMapper, contextProvider libs.Context, resource unstructured.Unstructured, dpols []v1alpha1.DeletingPolicy) ([]models.Response, error) {
	provider, err := engine.NewProvider(compiler.NewCompiler(), dpols, nil)

	engine := engine.NewEngine(utils.NSResolver(dClient), restMapper, contextProvider, matching.NewMatcher())
	if err != nil {
		return nil, err
	}

	policies, err := provider.Fetch(context.Background())
	if err != nil {
		return nil, err
	}

	results := make([]models.Response, 0)

	for _, policy := range policies {
		// deleting
		resp, err := engine.Handle(ctx, policy, resource)
		if err != nil {
			result := engineapi.NewEngineResponse(resource, engineapi.NewDeletingPolicy(&policy.Policy), nil)
			result = result.WithPolicyResponse(engineapi.PolicyResponse{Rules: []engineapi.RuleResponse{
				*engineapi.NewRuleResponse("", engineapi.Deletion, err.Error(), engineapi.RuleStatusError, nil),
			}})

			results = append(results, models.ConvertResponse(result))
			continue
		}

		status := engineapi.RuleStatusPass
		message := "resource matched"
		if !resp.Match {
			status = engineapi.RuleStatusFail
			message = "resource did not match"
		}

		result := engineapi.NewEngineResponse(resource, engineapi.NewDeletingPolicy(&policy.Policy), nil)
		result = result.WithPolicyResponse(engineapi.PolicyResponse{Rules: []engineapi.RuleResponse{
			*engineapi.NewRuleResponse("", engineapi.Deletion, message, status, nil),
		}})

		results = append(results, models.ConvertResponse(result))
	}

	return results, nil
}
