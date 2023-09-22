package mocks

import (
	"context"

	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	enginecontext "github.com/kyverno/kyverno/pkg/engine/context"
	"github.com/kyverno/kyverno/pkg/engine/jmespath"
	"github.com/kyverno/kyverno/pkg/imageverifycache"
)

type withoutAPICalls struct {
	next engineapi.ContextLoader
}

func WithoutAPICalls(next engineapi.ContextLoader) engineapi.ContextLoader {
	return withoutAPICalls{
		next: next,
	}
}

func (cl withoutAPICalls) Load(
	ctx context.Context,
	jp jmespath.Interface,
	client engineapi.RawClient,
	rclientFactory engineapi.RegistryClientFactory,
	ivCache imageverifycache.Client,
	contextEntries []kyvernov1.ContextEntry,
	jsonContext enginecontext.Interface,
) error {
	var contextEntriesWithoutAPICalls []kyvernov1.ContextEntry
	for _, entry := range contextEntries {
		if entry.APICall == nil {
			contextEntriesWithoutAPICalls = append(contextEntriesWithoutAPICalls, entry)
		}
	}
	return cl.next.Load(ctx, jp, client, rclientFactory, ivCache, contextEntriesWithoutAPICalls, jsonContext)
}

func WithCMCheck(resolver engineapi.ConfigmapResolver, next engineapi.ContextLoader) engineapi.ContextLoader {
	return withCMCheck{
		resolver: resolver,
		next:     next,
	}
}

type withCMCheck struct {
	resolver engineapi.ConfigmapResolver
	next     engineapi.ContextLoader
}

func (cl withCMCheck) Load(
	ctx context.Context,
	jp jmespath.Interface,
	client engineapi.RawClient,
	rclientFactory engineapi.RegistryClientFactory,
	ivCache imageverifycache.Client,
	contextEntries []kyvernov1.ContextEntry,
	jsonContext enginecontext.Interface,
) error {
	var contextEntriesWithExistingCMs []kyvernov1.ContextEntry
	for _, entry := range contextEntries {
		if entry.ConfigMap == nil {
			contextEntriesWithExistingCMs = append(contextEntriesWithExistingCMs, entry)
			continue
		}

		_, err := cl.resolver.Get(ctx, entry.ConfigMap.Namespace, entry.ConfigMap.Name)
		if err == nil {
			contextEntriesWithExistingCMs = append(contextEntriesWithExistingCMs, entry)
			continue
		}
	}
	return cl.next.Load(ctx, jp, client, rclientFactory, ivCache, contextEntriesWithExistingCMs, jsonContext)
}
