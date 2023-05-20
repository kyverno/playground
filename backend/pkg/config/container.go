package config

import (
	"context"
	"time"

	"github.com/golobby/container/v3"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/kyverno/kyverno/pkg/engine/context/resolvers"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"github.com/kyverno/playground/backend/pkg/engine"
)

func InitContainer(config *rest.Config) (container.Container, error) {
	c := container.New()

	if config == nil {
		return c, nil
	}

	err := c.Singleton(func() (dynamic.Interface, error) {
		return dynamic.NewForConfig(config)
	})
	if err != nil {
		return c, err
	}

	err = c.Singleton(func() (kubernetes.Interface, error) {
		return kubernetes.NewForConfig(config)
	})
	if err != nil {
		return c, err
	}

	err = c.Singleton(func(kube kubernetes.Interface, dyn dynamic.Interface) (dclient.Interface, error) {
		c, err := dclient.NewClient(context.Background(), dyn, kube, 15*time.Minute)
		if err != nil {
			return nil, err
		}

		return engine.NewWrapper(c), nil
	})
	if err != nil {
		return c, err
	}

	err = c.Singleton(func(kube kubernetes.Interface) (engineapi.ConfigmapResolver, error) {
		return resolvers.NewClientBasedResolver(kube)
	})
	return c, err
}
