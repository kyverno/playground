package ivpol

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

func K8sProcess(ctx context.Context, dClient dclient.Interface, restMapper meta.RESTMapper, contextProvider libs.Context, params *models.Parameters, newResource, oldResource unstructured.Unstructured, ivpols []v1beta1.ImageValidatingPolicyLike) ([]models.Response, error) {
	request := utils.NewCELRequest(restMapper, contextProvider, params, newResource, oldResource)
	validations := make([]models.Response, 0)

	filtered := helper.Filter(ivpols, func(p v1beta1.ImageValidatingPolicyLike) bool {
		return p.GetNamespace() == "" || p.GetNamespace() == newResource.GetNamespace()
	})

	engine, err := newIVPEngine(dClient, filtered, nil)
	if err != nil {
		return nil, err
	}

	resp, _, err := engine.HandleMutating(ctx, request, nil)
	if err != nil {
		return nil, err
	}

	for _, result := range resp.Policies {
		resp := engineapi.NewEngineResponse(newResource, engineapi.NewImageValidatingPolicyFromLike(result.Policy), params.Context.NamespaceLabels).
			WithPolicyResponse(engineapi.PolicyResponse{Rules: []engineapi.RuleResponse{result.Result}})

		validations = append(validations, models.ConvertResponse(resp))
	}

	return validations, nil
}
