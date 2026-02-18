package vpol

import (
	"context"

	"github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	"github.com/kyverno/kyverno/pkg/cel/libs"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/kyverno/playground/backend/pkg/engine/models"
	"github.com/kyverno/playground/backend/pkg/engine/utils"
	helper "github.com/kyverno/playground/backend/pkg/utils"
)

func K8sProcess(ctx context.Context, dClient dclient.Interface, restMapper meta.RESTMapper, contextProvider libs.Context, params *models.Parameters, newResource, oldResource unstructured.Unstructured, vpols []v1beta1.ValidatingPolicyLike) ([]models.Response, error) {
	request := utils.NewCELRequest(restMapper, contextProvider, params, newResource, oldResource)

	filtered := helper.Filter(vpols, func(p v1beta1.ValidatingPolicyLike) bool {
		return p.GetNamespace() == "" || p.GetNamespace() == newResource.GetNamespace()
	})

	eng, err := newCELEngine(dClient, filtered, nil)
	if err != nil {
		return nil, err
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
