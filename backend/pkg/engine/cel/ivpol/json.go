package ivpol

import (
	"context"

	"github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	"github.com/kyverno/kyverno/pkg/cel/libs"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	eval "github.com/kyverno/kyverno/pkg/imageverification/evaluator"
	"github.com/kyverno/playground/backend/pkg/engine/models"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func JSONProcess(ctx context.Context, contextProvider libs.Context, resource unstructured.Unstructured, ivpols []v1beta1.ImageValidatingPolicyLike) ([]models.Response, error) {
	compiled := make([]*eval.CompiledImageValidatingPolicy, 0)
	pMap := make(map[string]v1beta1.ImageValidatingPolicyLike)
	for i := range ivpols {
		p := ivpols[i]
		pMap[p.GetName()] = p
		compiled = append(compiled, &eval.CompiledImageValidatingPolicy{Policy: p})
	}

	results, err := eval.Evaluate(ctx, compiled, resource.Object, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	resp := engineapi.EngineResponse{
		Resource:       resource,
		PolicyResponse: engineapi.PolicyResponse{},
	}

	validations := make([]models.Response, 0)
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
		resp = resp.WithPolicy(engineapi.NewImageValidatingPolicyFromLike(pMap[p]))

		validations = append(validations, models.ConvertResponse(resp))
	}

	return validations, nil
}
