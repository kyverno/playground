package mpol

import (
	"context"

	"github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"github.com/kyverno/kyverno/pkg/cel/policies/mpol/compiler"
	"github.com/kyverno/kyverno/pkg/cel/policies/mpol/engine"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiserver/pkg/admission"
)

type staticProvider struct {
	inner engine.Provider
}

func (p *staticProvider) Fetch(ctx context.Context, mutateExisting bool) ([]engine.Policy, error) {
	p1, err := p.inner.Fetch(ctx, true)
	if err != nil {
		return nil, err
	}

	p2, err := p.inner.Fetch(ctx, false)
	if err != nil {
		return nil, err
	}

	return append(p1, p2...), nil
}

func (r *staticProvider) MatchesMutateExisting(ctx context.Context, attr admission.Attributes, namespace *corev1.Namespace) []string {
	return r.inner.MatchesMutateExisting(ctx, attr, namespace)
}

func NewProvider(compiler compiler.Compiler, policies []v1alpha1.MutatingPolicy, exceptions []*v1alpha1.PolicyException) (engine.Provider, error) {
	inner, err := engine.NewProvider(compiler, policies, exceptions)
	if err != nil {
		return nil, err
	}

	return &staticProvider{inner}, nil
}
