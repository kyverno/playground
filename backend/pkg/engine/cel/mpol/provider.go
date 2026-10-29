package mpol

import (
	"context"

	"github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	"github.com/kyverno/kyverno/pkg/cel/libs"
	"github.com/kyverno/kyverno/pkg/cel/policies/mpol/compiler"
	"github.com/kyverno/kyverno/pkg/cel/policies/mpol/engine"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiserver/pkg/admission"
)

type staticProvider struct {
	inner engine.Provider
}

func (p *staticProvider) Fetch(ctx context.Context, mutateExisting bool) []engine.Policy {
	return p.inner.Fetch(ctx, mutateExisting)
}

func (r *staticProvider) MatchesMutateExisting(ctx context.Context, attr admission.Attributes, req *admissionv1.AdmissionRequest, namespace *corev1.Namespace, requestMapFn func() (map[string]any, error)) []string {
	return r.inner.MatchesMutateExisting(ctx, attr, req, namespace, requestMapFn)
}

func NewProvider(compiler compiler.Compiler, policies []v1beta1.MutatingPolicyLike, exceptions []*v1beta1.PolicyException, libCtx libs.Context) (engine.Provider, error) {
	inner, err := engine.NewProvider(compiler, policies, exceptions, libCtx)
	if err != nil {
		return nil, err
	}

	return &staticProvider{inner}, nil
}
