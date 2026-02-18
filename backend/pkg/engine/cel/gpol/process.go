package gpol

import (
	"context"
	"fmt"

	"github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	"github.com/kyverno/kyverno/pkg/cel/libs"
	"github.com/kyverno/kyverno/pkg/cel/matching"
	"github.com/kyverno/kyverno/pkg/cel/policies/gpol/compiler"
	"github.com/kyverno/kyverno/pkg/cel/policies/gpol/engine"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/kyverno/playground/backend/pkg/engine/models"
	"github.com/kyverno/playground/backend/pkg/engine/utils"
)

func Process(ctx context.Context, dClient dclient.Interface, restMapper meta.RESTMapper, contextProvider libs.Context, params *models.Parameters, resource unstructured.Unstructured, gpols []v1beta1.GeneratingPolicyLike) ([]models.Response, error) {
	comp := compiler.NewCompiler()
	policies := make([]engine.Policy, 0, len(gpols))

	for _, pol := range gpols {
		compiled, errs := comp.Compile(pol, nil)
		if len(errs) > 0 {
			return nil, fmt.Errorf("failed to compile policy %s (%w)", pol.GetName(), errs.ToAggregate())
		}

		policies = append(policies, engine.Policy{
			Policy:         pol,
			CompiledPolicy: compiled,
		})
	}

	eng := engine.NewEngine(utils.NSResolver(dClient), matching.NewMatcher())

	request := utils.NewCELRequest(restMapper, contextProvider, params, resource, unstructured.Unstructured{})
	results := make([]models.Response, 0)

	for _, policy := range policies {
		if policy.Policy.GetNamespace() != "" && policy.Policy.GetNamespace() != resource.GetNamespace() {
			continue
		}

		engineResponse, err := eng.Handle(request, policy, false)
		if err != nil {
			return nil, err
		}

		for _, res := range engineResponse.Policies {
			if res.Result == nil {
				continue
			}

			generateResponse := engineapi.EngineResponse{
				Resource: *engineResponse.Trigger,
				PolicyResponse: engineapi.PolicyResponse{
					Rules: []engineapi.RuleResponse{*res.Result},
				},
			}

			response := generateResponse.WithPolicy(engineapi.NewGeneratingPolicyFromLike(res.Policy))

			var newRuleResponse []engineapi.RuleResponse
			for _, rule := range response.PolicyResponse.Rules {
				if len(rule.GeneratedResources()) == 0 {
					continue
				}

				for _, g := range rule.GeneratedResources() {
					// cleanup metadata
					if meta, ok := g.Object["metadata"]; ok {
						delete(meta.(map[string]any), "managedFields")
					}
				}

				newRuleResponse = append(newRuleResponse, *rule.WithGeneratedResources(rule.GeneratedResources()))
			}

			response.PolicyResponse.Rules = newRuleResponse
			results = append(results, models.ConvertResponse(response))
		}
	}

	return results, nil
}
