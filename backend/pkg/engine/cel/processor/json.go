package processor

import (
	"context"

	"github.com/kyverno/kyverno/pkg/cel/libs"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	gctxstore "github.com/kyverno/kyverno/pkg/globalcontext/store"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	mpatch "k8s.io/apiserver/pkg/admission/plugin/policy/mutating/patch"

	"github.com/kyverno/playground/backend/pkg/engine/cel/ivpol"
	"github.com/kyverno/playground/backend/pkg/engine/cel/mpol"
	"github.com/kyverno/playground/backend/pkg/engine/cel/vpol"
	"github.com/kyverno/playground/backend/pkg/engine/models"
	"github.com/kyverno/playground/backend/pkg/policy"
)

type JSONProcessor struct {
	dClient dclient.Interface
	tcm     mpatch.TypeConverterManager
}

func (p *JSONProcessor) Run(ctx context.Context, policies policy.JSONPolicies, resources []unstructured.Unstructured) (*models.Results, error) {
	response := &models.Results{}

	contextProvider, err := libs.NewContextProvider(p.dClient, nil, gctxstore.New(), nil, false)
	if err != nil {
		return nil, err
	}

	for _, resource := range resources {
		if len(policies.MutatingPolicies) > 0 {
			results, err := mpol.JSONProcess(context.TODO(), p.tcm, contextProvider, resource, policies.MutatingPolicies)
			if err != nil {
				return nil, err
			}

			response.Mutation = append(response.Mutation, results...)
		}

		if len(policies.ImageValidatingPolicies) > 0 {
			results, err := ivpol.JSONProcess(context.TODO(), contextProvider, resource, policies.ImageValidatingPolicies)
			if err != nil {
				return nil, err
			}

			response.Validation = append(response.Validation, results...)
		}

		if len(policies.ValidatingPolicies) > 0 {
			results, err := vpol.JSONProcess(context.TODO(), contextProvider, resource, policies.ValidatingPolicies)
			if err != nil {
				return nil, err
			}

			response.Validation = append(response.Validation, results...)
		}
	}

	return response, nil
}

func NewJSON(dClient dclient.Interface, tcm mpatch.TypeConverterManager) *JSONProcessor {
	if dClient == nil {
		dClient = dclient.NewEmptyFakeClient()
	}

	return &JSONProcessor{
		dClient: dClient,
		tcm:     tcm,
	}
}
