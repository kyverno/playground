package authz

import (
	"context"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/kyverno/kyverno/pkg/clients/dclient"

	"github.com/kyverno/playground/backend/pkg/engine/cel/vpol"
	"github.com/kyverno/playground/backend/pkg/engine/models"
	"github.com/kyverno/playground/backend/pkg/policy"
)

type EnvoyProcessor struct {
	dClient dclient.Interface
}

func (p *EnvoyProcessor) Run(ctx context.Context, policies policy.AuthzPolicies, resources []*authv3.CheckRequest) (*models.Results, error) {
	response := &models.Results{}

	for _, resource := range resources {
		if len(policies.EnvoyPolicies) > 0 {
			results, err := vpol.EnvoyProcess(context.TODO(), resource, policies.EnvoyPolicies)
			if err != nil {
				return nil, err
			}

			response.Validation = append(response.Validation, results...)
		}
	}

	return response, nil
}

func NewProcessor(dClient dclient.Interface) *EnvoyProcessor {
	if dClient == nil {
		dClient = dclient.NewEmptyFakeClient()
	}

	return &EnvoyProcessor{
		dClient: dClient,
	}
}
