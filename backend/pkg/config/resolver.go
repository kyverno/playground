package config

import (
	"github.com/golobby/container/v3"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Resolver struct {
	cluster   bool
	container container.Container
}

func (r *Resolver) DClient() (dclient.Interface, error) {
	return resolve[dclient.Interface](r.container, r.cluster)
}

func (r *Resolver) CMResolver() (engineapi.ConfigmapResolver, error) {
	return resolve[engineapi.ConfigmapResolver](r.container, r.cluster)
}

func (r *Resolver) KubeClient() (kubernetes.Interface, error) {
	return resolve[kubernetes.Interface](r.container, r.cluster)
}

func resolve[T any](c container.Container, cluster bool) (T, error) {
	var value T

	if !cluster {
		return value, nil
	}

	err := c.Resolve(&value)

	return value, err
}

func NewResolver(config *rest.Config) (*Resolver, error) {
	c, err := InitContainer(config)
	if err != nil {
		return nil, err
	}

	return &Resolver{
		cluster:   config != nil,
		container: c,
	}, nil
}
