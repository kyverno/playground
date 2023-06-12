package models

type Registry struct {
	AllowInsecure     bool     `json:"allowInsecure"`
	PullSecrets       []string `json:"pullSecrets"`
	CredentialHelpers []string `json:"credentialHelpers"`
}

type ProtectManagedResources struct {
	Enabled bool `json:"enabled"`
}

type ForceFailurePolicyIgnore struct {
	Enabled bool `json:"enabled"`
}

type Flags struct {
	Exceptions               Exceptions               `json:"exceptions"`
	Cosign                   Cosign                   `json:"cosign"`
	Registry                 Registry                 `json:"registry"`
	ProtectManagedResources  ProtectManagedResources  `json:"protectManagedResources"`
	ForceFailurePolicyIgnore ForceFailurePolicyIgnore `json:"forceFailurePolicyIgnore"`
}
