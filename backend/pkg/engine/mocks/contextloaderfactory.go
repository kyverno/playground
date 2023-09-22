package mocks

import (
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/kyverno/kyverno/pkg/engine/factories"
)

func ContextLoaderFactory(cmResolver engineapi.ConfigmapResolver) engineapi.ContextLoaderFactory {
	next := factories.DefaultContextLoaderFactory(cmResolver)
	return func(policy kyvernov1.PolicyInterface, rule kyvernov1.Rule) engineapi.ContextLoader {
		chain := next(policy, rule)
		chain = WithoutAPICalls(chain)
		chain = WithCMCheck(cmResolver, chain)
		return chain
	}
}
