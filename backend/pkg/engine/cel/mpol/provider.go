package mpol

import (
	"context"

	"github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	"github.com/kyverno/kyverno/pkg/cel/policies/mpol/compiler"
	"github.com/kyverno/kyverno/pkg/cel/policies/mpol/engine"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiserver/pkg/admission"
)

type staticProvider struct {
	inner engine.Provider
}

func (p *staticProvider) Fetch(ctx context.Context, mutateExisting bool) []engine.Policy {
	p1 := p.inner.Fetch(ctx, true)
	p2 := p.inner.Fetch(ctx, false)

	return append(p1, p2...)
}

func (r *staticProvider) MatchesMutateExisting(ctx context.Context, attr admission.Attributes, namespace *corev1.Namespace) []string {
	return r.inner.MatchesMutateExisting(ctx, attr, namespace)
}

func NewProvider(compiler compiler.Compiler, policies []v1beta1.MutatingPolicyLike, exceptions []*v1beta1.PolicyException) (engine.Provider, error) {
	inner, err := engine.NewProvider(compiler, policies, exceptions)
	if err != nil {
		return nil, err
	}

	return &staticProvider{inner}, nil
}
