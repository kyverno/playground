package models

type PolicyValidation struct {
	PolicyName      string `json:"policyName"`
	PolicyNamespace string `json:"policyNamespace"`
	Field           string `json:"field"`
	Type            string `json:"type"`
	Detail          string `json:"detail"`
}
