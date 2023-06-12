package mocks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-logr/logr"
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	enginecontext "github.com/kyverno/kyverno/pkg/engine/context"
	"github.com/kyverno/kyverno/pkg/engine/jmespath"
	"github.com/kyverno/kyverno/pkg/engine/variables"
	"github.com/kyverno/kyverno/pkg/registryclient"

	"github.com/kyverno/playground/backend/pkg/engine/models"
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
	client engineapi.Client,
	rclient registryclient.Client,
	contextEntries []kyvernov1.ContextEntry,
	jsonContext enginecontext.Interface,
) error {
	var contextEntriesWithoutApiCalls []kyvernov1.ContextEntry
	for _, entry := range contextEntries {
		if entry.APICall == nil {
			contextEntriesWithoutApiCalls = append(contextEntriesWithoutApiCalls, entry)
		}
	}
	return cl.next.Load(ctx, jp, client, rclient, contextEntriesWithoutApiCalls, jsonContext)
}

type withImageData struct {
	next      engineapi.ContextLoader
	imageData map[string]models.ImageData
}

func WithImageData(next engineapi.ContextLoader, imageData map[string]models.ImageData) engineapi.ContextLoader {
	return withImageData{
		next:      next,
		imageData: imageData,
	}
}

func (cl withImageData) Load(
	ctx context.Context,
	jp jmespath.Interface,
	client engineapi.Client,
	rclient registryclient.Client,
	contextEntries []kyvernov1.ContextEntry,
	jsonContext enginecontext.Interface,
) error {
	var contextEntriesWithoutImageRegistry []kyvernov1.ContextEntry
	for _, entry := range contextEntries {
		if entry.ImageRegistry != nil {
			jsonContext.AddDeferredLoader(entry.Name, func() error {
				if err := engineapi.LoadImageData(ctx, jp, rclient, logr.Discard(), entry, jsonContext); err != nil {
					if cl.imageData == nil {
						return err
					}
					imageData, err := cl.fetchImageData(ctx, jp, logr.Discard(), entry, jsonContext)
					if err != nil {
						return err
					}
					jsonBytes, err := json.Marshal(imageData)
					if err != nil {
						return err
					}
					if err := jsonContext.AddContextEntry(entry.Name, jsonBytes); err != nil {
						return fmt.Errorf("failed to add resource data to context: contextEntry: %v, error: %v", entry, err)
					}
					return nil
				}
				return nil
			})
		} else {
			contextEntriesWithoutImageRegistry = append(contextEntriesWithoutImageRegistry, entry)
		}
	}
	return cl.next.Load(ctx, jp, client, rclient, contextEntriesWithoutImageRegistry, jsonContext)
}

func (cl withImageData) fetchImageData(
	ctx context.Context,
	jp jmespath.Interface,
	logger logr.Logger,
	entry kyvernov1.ContextEntry,
	jsonContext enginecontext.Interface,
) (interface{}, error) {
	ref, err := variables.SubstituteAll(logger, jsonContext, entry.ImageRegistry.Reference)
	if err != nil {
		return nil, fmt.Errorf("ailed to substitute variables in context entry %s %s: %v", entry.Name, entry.ImageRegistry.Reference, err)
	}
	refString, ok := ref.(string)
	if !ok {
		return nil, fmt.Errorf("invalid image reference %s, image reference must be a string", ref)
	}
	path, err := variables.SubstituteAll(logger, jsonContext, entry.ImageRegistry.JMESPath)
	if err != nil {
		return nil, fmt.Errorf("failed to substitute variables in context entry %s %s: %v", entry.Name, entry.ImageRegistry.JMESPath, err)
	}
	imageData, err := cl.fetchImageDataMap(ctx, refString)
	if err != nil {
		return nil, err
	}
	if path != "" {
		imageData, err = cl.applyJMESPath(jp, path.(string), imageData)
		if err != nil {
			return nil, fmt.Errorf("failed to apply JMESPath (%s) results to context entry %s, error: %v", entry.ImageRegistry.JMESPath, entry.Name, err)
		}
	}
	return imageData, nil
}

func (cl withImageData) fetchImageDataMap(ctx context.Context, ref string) (interface{}, error) {
	if cl.imageData == nil {
		return nil, fmt.Errorf("failed to fetch image descriptor: %s", ref)
	}
	imageData, ok := cl.imageData[ref]
	if !ok {
		return nil, fmt.Errorf("failed to fetch image descriptor: %s", ref)
	}
	// we need to do the conversion from struct types to an interface type so that jmespath
	// evaluation works correctly. go-jmespath cannot handle function calls like max/sum
	// for types like integers for eg. the conversion to untyped allows the stdlib json
	// to convert all the types to types that are compatible with jmespath.
	jsonDoc, err := json.Marshal(imageData)
	if err != nil {
		return nil, err
	}
	var untyped interface{}
	err = json.Unmarshal(jsonDoc, &untyped)
	if err != nil {
		return nil, err
	}
	return untyped, nil
}

func (cl withImageData) applyJMESPath(jp jmespath.Interface, query string, data interface{}) (interface{}, error) {
	q, err := jp.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to compile JMESPath: %s, error: %v", query, err)
	}
	return q.Search(data)
}
