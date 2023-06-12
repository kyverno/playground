package models

type PolicyResponse struct {
	// Rules contains policy rules responses
	Rules []RuleResponse `json:"rules"`
}
