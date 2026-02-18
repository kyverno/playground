package vpol

import (
	"github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	"github.com/kyverno/kyverno/pkg/cel/matching"
	"github.com/kyverno/kyverno/pkg/cel/policies/vpol/engine"
	"github.com/kyverno/kyverno/pkg/clients/dclient"

	"github.com/kyverno/playground/backend/pkg/engine/utils"
)

func newCELEngine(dClient dclient.Interface, vpolicies []v1beta1.ValidatingPolicyLike, exceptions []*v1beta1.PolicyException) (engine.Engine, error) {
	provider, err := newVPOLProvider(vpolicies, exceptions)
	if err != nil {
		return nil, err
	}
	return engine.NewEngine(provider, utils.NSResolver(dClient), matching.NewMatcher()), nil
}
