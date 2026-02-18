package vpol

import (
	"github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	"github.com/kyverno/kyverno/pkg/cel/policies/vpol/compiler"
	"github.com/kyverno/kyverno/pkg/cel/policies/vpol/engine"
)

func newVPOLProvider(policies []v1beta1.ValidatingPolicyLike, exceptions []*v1beta1.PolicyException) (engine.Provider, error) {
	return engine.NewProvider(compiler.NewCompiler(), policies, exceptions)
}
