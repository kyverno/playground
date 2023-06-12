package mocks

import "github.com/kyverno/kyverno/pkg/toggle"

type toggles struct {
	protectManagedResources  bool
	forceFailurePolicyIgnore bool
}

func Toggles(protectManagedResources, forceFailurePolicyIgnore bool) toggle.Toggles {
	return toggles{
		protectManagedResources:  protectManagedResources,
		forceFailurePolicyIgnore: forceFailurePolicyIgnore,
	}
}

func (t toggles) ProtectManagedResources() bool {
	return t.protectManagedResources
}

func (t toggles) ForceFailurePolicyIgnore() bool {
	return t.forceFailurePolicyIgnore
}
