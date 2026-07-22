package tools

import (
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"

	"github.com/kyverno/playground/backend/pkg/playground"
)

type Result struct {
	PolicyType        string               `json:"policyType"`
	Mode              string               `json:"mode"`
	Policy            string               `json:"policy"`
	Resource          string               `json:"resource"`
	Namespace         string               `json:"namespace"`
	Message           string               `json:"message,omitempty"`
	Result            engineapi.RuleStatus `json:"result"`
	PatchedResource   string               `json:"patchedResource,omitempty"`
	GeneratedResource string               `json:"generatedResource,omitempty"`
	Properties        map[string]string    `json:"properties,omitempty"`
}

type Results[T any] struct {
	Results []T `json:"results"`
}

func MapResponse(response *playground.EngineResponse) Results[Result] {
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

	return Results[Result]{
		Results: result,
	}
}
