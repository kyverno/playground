package vpol

import (
	"context"

	"github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
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
)

func Process(ctx context.Context, dClient dclient.Interface, restMapper meta.RESTMapper, contextProvider libs.Context, params *models.Parameters, newResource, oldResource unstructured.Unstructured, vpols []v1alpha1.ValidatingPolicy) ([]models.Response, error) {
	var eng engine.Engine
	var err error

	request := utils.NewCELRequest(restMapper, contextProvider, params, newResource, oldResource)

	if request.JsonPayload == nil {
		eng, err = newCELEngine(dClient, vpols, nil)
		if err != nil {
			return nil, err
		}
	} else {
		provider, err := newVPOLProvider(vpols, nil)
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
		resp := engineapi.NewEngineResponse(newResource, engineapi.NewValidatingPolicy(&result.Policy), params.Context.NamespaceLabels).
			WithPolicyResponse(engineapi.PolicyResponse{Rules: result.Rules})

		results = append(results, models.ConvertResponse(resp))
	}

	return results, nil
}

func newCELEngine(dClient dclient.Interface, vpolicies []v1alpha1.ValidatingPolicy, exceptions []*v1alpha1.PolicyException) (engine.Engine, error) {
	provider, err := newVPOLProvider(vpolicies, exceptions)
	if err != nil {
		return nil, err
	}
	return engine.NewEngine(provider, utils.NSResolver(dClient), matching.NewMatcher()), nil
}

func newVPOLProvider(policies []v1alpha1.ValidatingPolicy, exceptions []*v1alpha1.PolicyException) (engine.Provider, error) {
	return engine.NewProvider(compiler.NewCompiler(), policies, exceptions)
}
