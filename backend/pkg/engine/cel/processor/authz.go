package processor

import (
	"context"

	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/kyverno/kyverno-authz/pkg/cel/libs/authz/http"
	"github.com/kyverno/kyverno/pkg/clients/dclient"

	"github.com/kyverno/playground/backend/pkg/engine/cel/vpol"
	"github.com/kyverno/playground/backend/pkg/engine/models"
	"github.com/kyverno/playground/backend/pkg/policy"
)

type AuthzProcessor struct {
	dClient dclient.Interface
}

func (p *AuthzProcessor) RunEnvoy(ctx context.Context, policies policy.AuthzPolicies, resources []*authv3.CheckRequest) (*models.Results, error) {
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

func (p *AuthzProcessor) RunHTTP(ctx context.Context, policies policy.AuthzPolicies, resources []*http.CheckRequest) (*models.Results, error) {
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

func NewAuthz(dClient dclient.Interface) *AuthzProcessor {
	if dClient == nil {
		dClient = dclient.NewEmptyFakeClient()
	}

	return &AuthzProcessor{
		dClient: dClient,
	}
}
