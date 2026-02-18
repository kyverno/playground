package http

import (
	"context"

	"github.com/kyverno/kyverno-authz/pkg/cel/libs/authz/http"
	"github.com/kyverno/kyverno/pkg/clients/dclient"

	"github.com/kyverno/playground/backend/pkg/engine/cel/vpol"
	"github.com/kyverno/playground/backend/pkg/engine/models"
	"github.com/kyverno/playground/backend/pkg/policy"
)

type HTTPProcessor struct {
	dClient dclient.Interface
}

func (p *HTTPProcessor) Run(ctx context.Context, policies policy.AuthzPolicies, resources []*http.CheckRequest) (*models.Results, error) {
	response := &models.Results{}

	for _, resource := range resources {
		if len(policies.HTTPPolicies) > 0 {
			results, err := vpol.HTTPProcess(context.TODO(), resource, policies.HTTPPolicies)
			if err != nil {
				return nil, err
			}

			response.Validation = append(response.Validation, results...)
		}
	}

	return response, nil
}

func NewProcessor(dClient dclient.Interface) *HTTPProcessor {
	if dClient == nil {
		dClient = dclient.NewEmptyFakeClient()
	}

	return &HTTPProcessor{
		dClient: dClient,
	}
}
