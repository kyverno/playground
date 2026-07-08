package mocks

import (
	"context"
	"encoding/json"
	"fmt"

	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	"github.com/kyverno/kyverno/pkg/engine/adapters"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/kyverno/kyverno/pkg/registryclient"

	"github.com/kyverno/playground/backend/pkg/engine/models"
)

type registryClientAdapter struct {
	engineapi.RegistryClient
	imageData map[string]models.ImageData
}

func ImageDataClient(next engineapi.RegistryClient, imageData map[string]models.ImageData) engineapi.RegistryClient {
	if next == nil {
		return nil
	}
	return registryClientAdapter{
		RegistryClient: next,
		imageData:      imageData,
	}
}

func (c registryClientAdapter) ForRef(ctx context.Context, ref string) (*engineapi.ImageData, error) {
	if data, err := c.RegistryClient.ForRef(ctx, ref); err == nil {
		return data, err
	}
	if c.imageData == nil {
		return nil, fmt.Errorf("failed to fetch image descriptor: %s", ref)
	}
	imageData, ok := c.imageData[ref]
	if !ok {
		return nil, fmt.Errorf("failed to fetch image descriptor: %s", ref)
	}
	configData, err := json.Marshal(imageData.ConfigData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config data for image descriptor: %s", ref)
	}
	manifest, err := json.Marshal(imageData.Manifest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config data for image descriptor: %s", ref)
	}
	return &engineapi.ImageData{
		Image:         imageData.Image,
		ResolvedImage: imageData.ResolvedImage,
		Registry:      imageData.Registry,
		Repository:    imageData.Repository,
		Identifier:    imageData.Identifier,
		Config:        configData,
		Manifest:      manifest,
	}, nil
}

type registryClientFactory struct {
	client engineapi.RegistryClient
}

func (f *registryClientFactory) GetClient(_ context.Context, _ *kyvernov1.ImageRegistryCredentials) (engineapi.RegistryClient, error) {
	return f.client, nil
}

func NewRegistryClientFactory(rclient registryclient.Client, imageData map[string]models.ImageData) engineapi.RegistryClientFactory {
	return &registryClientFactory{
		client: &registryClientAdapter{RegistryClient: adapters.RegistryClient(rclient), imageData: imageData},
	}
}
