package cluster

import (
	"context"

	"github.com/go-logr/logr"
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/utils/store"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	enginecontext "github.com/kyverno/kyverno/pkg/engine/context"
	"github.com/kyverno/kyverno/pkg/engine/jmespath"
	"github.com/kyverno/kyverno/pkg/logging"
	"github.com/kyverno/kyverno/pkg/registryclient"
)

func ContextLoaderFactory(fake bool, cmResolver engineapi.ConfigmapResolver) engineapi.ContextLoaderFactory {
	return func(policy kyvernov1.PolicyInterface, rule kyvernov1.Rule) engineapi.ContextLoader {
		if fake {
			return &fakeContextLoader{
				cmResolver: cmResolver,
				logger:     logging.WithName("MockContextLoaderFactory"),
				policyName: policy.GetName(),
				ruleName:   rule.Name,
			}
		}

		inner := engineapi.DefaultContextLoaderFactory(cmResolver)

		return inner(policy, rule)
	}
}

type fakeContextLoader struct {
	cmResolver engineapi.ConfigmapResolver
	logger     logr.Logger
	policyName string
	ruleName   string
}

func (l *fakeContextLoader) Load(
	ctx context.Context,
	jp jmespath.Interface,
	_ dclient.Interface,
	_ registryclient.Client,
	contextEntries []kyvernov1.ContextEntry,
	jsonContext enginecontext.Interface,
) error {
	rule := store.GetPolicyRule(l.policyName, l.ruleName)
	if rule != nil && len(rule.Values) > 0 {
		variables := rule.Values
		for key, value := range variables {
			if err := jsonContext.AddVariable(key, value); err != nil {
				return err
			}
		}
	}
	for _, entry := range contextEntries {
		if entry.ConfigMap != nil {
			_ = engineapi.LoadConfigMap(ctx, l.logger, entry, jsonContext, l.cmResolver)
		} else if entry.ImageRegistry != nil {
			rclient := store.GetRegistryClient()
			if err := engineapi.LoadImageData(ctx, jp, rclient, l.logger, entry, jsonContext); err != nil {
				return err
			}
		} else if entry.Variable != nil {
			if err := engineapi.LoadVariable(l.logger, jp, entry, jsonContext); err != nil {
				return err
			}
		}
	}
	if rule != nil && len(rule.ForEachValues) > 0 {
		for key, value := range rule.ForEachValues {
			if err := jsonContext.AddVariable(key, value[store.GetForeachElement()]); err != nil {
				return err
			}
		}
	}
	return nil
}
