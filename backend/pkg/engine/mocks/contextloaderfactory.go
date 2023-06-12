package mocks

import (
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
)

func ContextLoaderFactory(cmResolver engineapi.ConfigmapResolver) engineapi.ContextLoaderFactory {
	next := engineapi.DefaultContextLoaderFactory(cmResolver)
	return func(policy kyvernov1.PolicyInterface, rule kyvernov1.Rule) engineapi.ContextLoader {
		return WithoutApiCalls(next(policy, rule))
	}
}
