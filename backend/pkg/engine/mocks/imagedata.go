package mocks

import (
	"context"
	"encoding/json"
	"fmt"

	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/kyverno/playground/backend/pkg/engine/models"
)

type imageDataClient struct {
	next      engineapi.ImageDataClient
	imageData map[string]models.ImageData
}

func ImageDataClient(next engineapi.ImageDataClient, imageData map[string]models.ImageData) engineapi.ImageDataClient {
	if next == nil {
		return nil
	}
	return imageDataClient{
		next:      next,
		imageData: imageData,
	}
}

func (c imageDataClient) ForRef(ctx context.Context, ref string) (*engineapi.ImageData, error) {
	if data, err := c.next.ForRef(ctx, ref); err == nil {
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
