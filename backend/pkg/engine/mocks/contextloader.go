package mocks

import (
	"context"

	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	enginecontext "github.com/kyverno/kyverno/pkg/engine/context"
	"github.com/kyverno/kyverno/pkg/engine/jmespath"
)

type withoutApiCalls struct {
	next engineapi.ContextLoader
}

func WithoutApiCalls(next engineapi.ContextLoader) engineapi.ContextLoader {
	return withoutApiCalls{
		next: next,
	}
}

func (cl withoutApiCalls) Load(
	ctx context.Context,
	jp jmespath.Interface,
	client engineapi.RawClient,
	imgClient engineapi.ImageDataClient,
	contextEntries []kyvernov1.ContextEntry,
	jsonContext enginecontext.Interface,
) error {
	var contextEntriesWithoutApiCalls []kyvernov1.ContextEntry
	for _, entry := range contextEntries {
		if entry.APICall == nil {
			contextEntriesWithoutApiCalls = append(contextEntriesWithoutApiCalls, entry)
		}
	}
	return cl.next.Load(ctx, jp, client, imgClient, contextEntriesWithoutApiCalls, jsonContext)
}
