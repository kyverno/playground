package dpol

import (
	"context"

	"github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	"github.com/kyverno/kyverno/pkg/cel/libs"
	"github.com/kyverno/kyverno/pkg/cel/matching"
	"github.com/kyverno/kyverno/pkg/cel/policies/dpol/compiler"
	"github.com/kyverno/kyverno/pkg/cel/policies/dpol/engine"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/kyverno/playground/backend/pkg/engine/models"
	"github.com/kyverno/playground/backend/pkg/engine/utils"
	helper "github.com/kyverno/playground/backend/pkg/utils"
)

func Process(ctx context.Context, dClient dclient.Interface, restMapper meta.RESTMapper, contextProvider libs.Context, resource unstructured.Unstructured, dpols []v1beta1.DeletingPolicyLike) ([]models.Response, error) {
	filtered := helper.Filter(dpols, func(p v1beta1.DeletingPolicyLike) bool {
		return p.GetNamespace() == "" || p.GetNamespace() == resource.GetNamespace()
	})

	provider, err := engine.NewProvider(compiler.NewCompiler(), filtered, nil)
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
			result := engineapi.NewEngineResponse(resource, engineapi.NewDeletingPolicyFromLike(policy.Policy), nil)
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

		result := engineapi.NewEngineResponse(resource, engineapi.NewDeletingPolicyFromLike(policy.Policy), nil)
		result = result.WithPolicyResponse(engineapi.PolicyResponse{Rules: []engineapi.RuleResponse{
			*engineapi.NewRuleResponse("", engineapi.Deletion, message, status, nil),
		}})

		results = append(results, models.ConvertResponse(result))
	}

	return results, nil
}
