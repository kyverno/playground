package mocks

import (
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"

	"github.com/kyverno/playground/backend/pkg/engine/models"
)

func ContextLoaderFactory(cmResolver engineapi.ConfigmapResolver, imageData map[string]models.ImageData) engineapi.ContextLoaderFactory {
	next := engineapi.DefaultContextLoaderFactory(cmResolver)
	return func(policy kyvernov1.PolicyInterface, rule kyvernov1.Rule) engineapi.ContextLoader {
		chain := next(policy, rule)
		chain = WithImageData(chain, imageData)
		chain = WithoutApiCalls(chain)
		return chain
	}
}
