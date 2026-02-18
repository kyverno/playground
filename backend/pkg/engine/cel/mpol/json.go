package mpol

import (
	"context"

	"github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	celengine "github.com/kyverno/kyverno/pkg/cel/engine"
	"github.com/kyverno/kyverno/pkg/cel/libs"
	"github.com/kyverno/kyverno/pkg/cel/matching"
	"github.com/kyverno/kyverno/pkg/cel/policies/mpol/compiler"
	"github.com/kyverno/kyverno/pkg/cel/policies/mpol/engine"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/kyverno/playground/backend/pkg/engine/models"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apiserver/pkg/admission/plugin/policy/mutating/patch"
)

func JSONProcess(ctx context.Context, tcm patch.TypeConverterManager, contextProvider libs.Context, resource unstructured.Unstructured, mpols []v1beta1.MutatingPolicyLike) ([]models.Response, error) {
	provider, err := NewProvider(compiler.NewCompiler(), mpols, nil)
	if err != nil {
		return nil, err
	}

	eng := engine.NewEngine(provider, nil, matching.NewMatcher(), tcm, contextProvider)
	request := celengine.RequestFromJSON(contextProvider, &resource)
	results := make([]models.Response, 0)

	engineResponse, err := eng.Handle(ctx, request, nil)
	if err != nil {
		return nil, err
	}

	for _, res := range engineResponse.Policies {
		var patched unstructured.Unstructured
		if engineResponse.PatchedResource != nil {
			patched = *engineResponse.PatchedResource
		}

		response := engineapi.EngineResponse{
			Resource:        *engineResponse.Resource,
			PatchedResource: patched,
			PolicyResponse: engineapi.PolicyResponse{
				Rules: res.Rules,
			},
		}

		results = append(results, models.ConvertResponse(response.WithPolicy(engineapi.NewMutatingPolicyFromLike(res.Policy))))
	}

	return results, nil
}
