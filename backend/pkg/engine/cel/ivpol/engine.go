package ivpol

import (
	"github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	"github.com/kyverno/kyverno/pkg/cel/matching"
	"github.com/kyverno/kyverno/pkg/cel/policies/ivpol/engine"
	"github.com/kyverno/kyverno/pkg/clients/dclient"

	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/engine/utils"
)

func newIVPEngine(dClient dclient.Interface, policies []v1beta1.ImageValidatingPolicyLike, exceptions []*v1beta1.PolicyException) (engine.Engine, error) {
	provider, err := engine.NewProvider(policies, exceptions)
	if err != nil {
		return nil, err
	}

	var nsResolver engine.NamespaceResolver
	if dClient != nil {
		nsResolver = utils.NSResolver(dClient)
	}

	// Like the Kyverno CLI, run without an image verification cache: every
	// playground request is evaluated once, and imageverifycache.New() with no
	// options fails with "NumCounters can't be zero".
	return engine.NewEngine(
		provider,
		nsResolver,
		matching.NewMatcher(),
		cluster.NewSecretLister(dClient, ""),
		nil,
		nil,
	), nil
}
