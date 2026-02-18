package ivpol

import (
	"github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	"github.com/kyverno/kyverno/pkg/cel/matching"
	"github.com/kyverno/kyverno/pkg/cel/policies/ivpol/engine"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	"github.com/kyverno/playground/backend/pkg/engine/utils"
	k8scorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func newIVPEngine(dClient dclient.Interface, policies []v1beta1.ImageValidatingPolicyLike, exceptions []*v1beta1.PolicyException) (engine.Engine, error) {
	provider, err := engine.NewProvider(policies, exceptions)
	if err != nil {
		return nil, err
	}

	var secretsClient k8scorev1.SecretInterface
	if dClient != nil {
		secretsClient = dClient.GetKubeClient().CoreV1().Secrets("")
	}

	var nsResolver engine.NamespaceResolver
	if dClient != nil {
		nsResolver = utils.NSResolver(dClient)
	}

	return engine.NewEngine(
		provider,
		nsResolver,
		matching.NewMatcher(),
		secretsClient,
		nil,
	), nil
}
