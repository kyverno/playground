package models

type Registry struct {
	AllowInsecure     bool     `json:"allowInsecure"`
	PullSecrets       []string `json:"pullSecrets"`
	CredentialHelpers []string `json:"credentialHelpers"`
}
