package vpol

import (
	"context"

	"github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	"github.com/kyverno/kyverno/pkg/cel/libs"
	"github.com/kyverno/kyverno/pkg/cel/matching"
	"github.com/kyverno/kyverno/pkg/cel/policies/vpol/compiler"
	"github.com/kyverno/kyverno/pkg/cel/policies/vpol/engine"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/kyverno/playground/backend/pkg/engine/models"
	"github.com/kyverno/playground/backend/pkg/engine/utils"
	helper "github.com/kyverno/playground/backend/pkg/utils"
)

func Process(ctx context.Context, dClient dclient.Interface, restMapper meta.RESTMapper, contextProvider libs.Context, params *models.Parameters, newResource, oldResource unstructured.Unstructured, vpols []v1beta1.ValidatingPolicyLike) ([]models.Response, error) {
	var eng engine.Engine
	var err error

	request := utils.NewCELRequest(restMapper, contextProvider, params, newResource, oldResource)

	filtered := helper.Filter(vpols, func(p v1beta1.ValidatingPolicyLike) bool {
		return p.GetNamespace() == "" || p.GetNamespace() == newResource.GetNamespace()
	})

	if request.JsonPayload == nil {
		eng, err = newCELEngine(dClient, filtered, nil)
		if err != nil {
			return nil, err
		}
	} else {
		provider, err := newVPOLProvider(filtered, nil)
		if err != nil {
			return nil, err
		}

		eng = engine.NewEngine(provider, nil, nil)
	}

	// validate
	resp, err := eng.Handle(ctx, request, nil)
	if err != nil {
		return nil, err
	}

	results := make([]models.Response, 0)
	for _, result := range resp.Policies {
		resp := engineapi.NewEngineResponse(newResource, engineapi.NewValidatingPolicyFromLike(result.Policy), params.Context.NamespaceLabels).
			WithPolicyResponse(engineapi.PolicyResponse{Rules: result.Rules})

		results = append(results, models.ConvertResponse(resp))
	}

	return results, nil
}

func newCELEngine(dClient dclient.Interface, vpolicies []v1beta1.ValidatingPolicyLike, exceptions []*v1beta1.PolicyException) (engine.Engine, error) {
	provider, err := newVPOLProvider(vpolicies, exceptions)
	if err != nil {
		return nil, err
	}
	return engine.NewEngine(provider, utils.NSResolver(dClient), matching.NewMatcher()), nil
}

func newVPOLProvider(policies []v1beta1.ValidatingPolicyLike, exceptions []*v1beta1.PolicyException) (engine.Provider, error) {
	return engine.NewProvider(compiler.NewCompiler(), policies, exceptions)
}
