package tools

import (
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"

	"github.com/kyverno/playground/backend/pkg/playground"
)

type Result struct {
	PolicyType        string               `json:"policyType" jsonschema:"The type of the policy (Validation, Mutation, Generation, Deletion)"`
	Mode              string               `json:"mode" jsonschema:"The mode of the policy (e.g., Kubernetes, Envoy, HTTP, JSON)"`
	Policy            string               `json:"policy" jsonschema:"The name of the policy"`
	Resource          string               `json:"resource" jsonschema:"The name of the resource"`
	Namespace         string               `json:"namespace" jsonschema:"The namespace of the resource"`
	Message           string               `json:"message,omitempty" jsonschema:"The message returned by the policy execution"`
	Result            engineapi.RuleStatus `json:"result" jsonschema:"The result of the policy execution (Pass, Fail, Warn, Skip, Error)"`
	PatchedResource   string               `json:"patchedResource,omitempty" jsonschema:"The patched resource returned by mutating policies"`
	GeneratedResource string               `json:"generatedResource,omitempty" jsonschema:"The generated resource returned by generating policies"`
	Properties        map[string]string    `json:"properties,omitempty" jsonschema:"The properties returned by the policy execution"`
}

type Results[T any] struct {
	Results []T `json:"results"`
}

func MapResponse(response *playground.EngineResponse) ExecuteOutputSchema {
	result := make([]Result, 0)
	for _, v := range response.Validation {
		for _, ruleResponse := range v.PolicyResponse.Rules {
			result = append(result, Result{
				PolicyType: "Validation",
				Mode:       v.Policy.Mode,
				Policy:     v.Policy.Name,
				Resource:   v.Resource.GetName(),
				Namespace:  v.Resource.GetNamespace(),
				Result:     ruleResponse.Status,
				Properties: ruleResponse.Properties,
				Message:    ruleResponse.Message,
			})
		}
	}

	for _, v := range response.Deletion {
		for _, ruleResponse := range v.PolicyResponse.Rules {
			result = append(result, Result{
				PolicyType: "Deletion",
				Mode:       v.Policy.Mode,
				Policy:     v.Policy.Name,
				Resource:   v.Resource.GetName(),
				Namespace:  v.Resource.GetNamespace(),
				Result:     ruleResponse.Status,
				Properties: ruleResponse.Properties,
				Message:    ruleResponse.Message,
			})
		}
	}

	for _, v := range response.Mutation {
		for _, ruleResponse := range v.PolicyResponse.Rules {
			result = append(result, Result{
				PolicyType:      "Mutation",
				Mode:            v.Policy.Mode,
				Policy:          v.Policy.Name,
				Resource:        v.Resource.GetName(),
				Namespace:       v.Resource.GetNamespace(),
				PatchedResource: v.PatchedResource,
				Result:          ruleResponse.Status,
				Properties:      ruleResponse.Properties,
				Message:         ruleResponse.Message,
			})
		}
	}

	for _, v := range response.Generation {
		for _, ruleResponse := range v.PolicyResponse.Rules {
			result = append(result, Result{
				PolicyType:        "Generation",
				Mode:              v.Policy.Mode,
				Policy:            v.Policy.Name,
				Resource:          v.Resource.GetName(),
				Namespace:         v.Resource.GetNamespace(),
				GeneratedResource: ruleResponse.GeneratedResource,
				Result:            ruleResponse.Status,
				Properties:        ruleResponse.Properties,
				Message:           ruleResponse.Message,
			})
		}
	}

	return ExecuteOutputSchema{
		Results: result,
	}
}
