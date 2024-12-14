package mocks

import "github.com/kyverno/kyverno/pkg/toggle"

type toggles struct {
	protectManagedResources           bool
	forceFailurePolicyIgnore          bool
	enableDeferredLoading             bool
	generateValidatingAdmissionPolicy bool
	dumpMutatePatches                 bool
	autogenV2                         bool
}

func Toggles(protectManagedResources, forceFailurePolicyIgnore, enableDeferredLoading, generateValidatingAdmissionPolicy, autogenV2 bool) toggle.Toggles {
	return toggles{
		protectManagedResources:           protectManagedResources,
		forceFailurePolicyIgnore:          forceFailurePolicyIgnore,
		enableDeferredLoading:             enableDeferredLoading,
		generateValidatingAdmissionPolicy: generateValidatingAdmissionPolicy,
		autogenV2:                         autogenV2,
	}
}

func (t toggles) ProtectManagedResources() bool {
	return t.protectManagedResources
}

func (t toggles) ForceFailurePolicyIgnore() bool {
	return t.forceFailurePolicyIgnore
}

func (t toggles) EnableDeferredLoading() bool {
	return t.enableDeferredLoading
}

func (t toggles) GenerateValidatingAdmissionPolicy() bool {
	return t.generateValidatingAdmissionPolicy
}

func (t toggles) DumpMutatePatches() bool {
	return t.dumpMutatePatches
}

func (t toggles) AutogenV2() bool {
	return t.dumpMutatePatches
}
