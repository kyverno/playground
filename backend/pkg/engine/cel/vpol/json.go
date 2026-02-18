package vpol

import (
	"context"

	"github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	celengine "github.com/kyverno/kyverno/pkg/cel/engine"
	"github.com/kyverno/kyverno/pkg/cel/libs"
	"github.com/kyverno/kyverno/pkg/cel/policies/vpol/engine"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/kyverno/playground/backend/pkg/engine/models"
)

func JSONProcess(ctx context.Context, contextProvider libs.Context, resource unstructured.Unstructured, vpols []v1beta1.ValidatingPolicyLike) ([]models.Response, error) {
	request := celengine.RequestFromJSON(contextProvider, &resource)

	provider, err := newVPOLProvider(vpols, nil)
	if err != nil {
		return nil, err
	}

	eng := engine.NewEngine(provider, nil, nil)
	resp, err := eng.Handle(ctx, request, nil)
	if err != nil {
		return nil, err
	}

	results := make([]models.Response, 0)
	for _, result := range resp.Policies {
		resp := engineapi.NewEngineResponse(resource, engineapi.NewValidatingPolicyFromLike(result.Policy), nil).
			WithPolicyResponse(engineapi.PolicyResponse{Rules: result.Rules})

		results = append(results, models.ConvertResponse(resp))
	}

	return results, nil
}
