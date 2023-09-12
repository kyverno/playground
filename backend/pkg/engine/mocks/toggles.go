package mocks

import "github.com/kyverno/kyverno/pkg/toggle"

type toggles struct {
	protectManagedResources  bool
	forceFailurePolicyIgnore bool
	enableDeferredLoading    bool
}

func Toggles(protectManagedResources, forceFailurePolicyIgnore, enableDeferredLoading bool) toggle.Toggles {
	return toggles{
		protectManagedResources:  protectManagedResources,
		forceFailurePolicyIgnore: forceFailurePolicyIgnore,
		enableDeferredLoading:    enableDeferredLoading,
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
