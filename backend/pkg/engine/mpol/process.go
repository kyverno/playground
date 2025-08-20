package mpol

import (
	"context"

	"github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"github.com/kyverno/kyverno/pkg/cel/libs"
	"github.com/kyverno/kyverno/pkg/cel/matching"
	"github.com/kyverno/kyverno/pkg/cel/policies/mpol/compiler"
	"github.com/kyverno/kyverno/pkg/cel/policies/mpol/engine"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apiserver/pkg/admission/plugin/policy/mutating/patch"

	"github.com/kyverno/playground/backend/pkg/engine/models"
	"github.com/kyverno/playground/backend/pkg/engine/utils"
)

func Process(ctx context.Context, dClient dclient.Interface, tcm patch.TypeConverterManager, restMapper meta.RESTMapper, contextProvider libs.Context, params *models.Parameters, resource, oldResource unstructured.Unstructured, mpols []v1alpha1.MutatingPolicy) ([]models.Response, error) {
	provider, err := NewProvider(compiler.NewCompiler(), mpols, nil)
	if err != nil {
		return nil, err
	}

	eng := engine.NewEngine(provider, utils.NSResolver(dClient), matching.NewMatcher(), tcm, contextProvider)
	request := utils.NewCELRequest(restMapper, contextProvider, params, resource, oldResource)
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

		results = append(results, models.ConvertResponse(response.WithPolicy(engineapi.NewMutatingPolicy(res.Policy))))
	}

	return results, nil
}
