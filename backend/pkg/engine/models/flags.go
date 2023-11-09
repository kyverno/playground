package models

type Flags struct {
	Exceptions                        Exceptions                        `json:"exceptions"`
	Cosign                            Cosign                            `json:"cosign"`
	Registry                          Registry                          `json:"registry"`
	ProtectManagedResources           ProtectManagedResources           `json:"protectManagedResources"`
	ForceFailurePolicyIgnore          ForceFailurePolicyIgnore          `json:"forceFailurePolicyIgnore"`
	EnableDeferredLoading             EnableDeferredLoading             `json:"enableDeferredLoading"`
	GenerateValidatingAdmissionPolicy GenerateValidatingAdmissionPolicy `json:"generateValidatingAdmissionPolicy"`
}
