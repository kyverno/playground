package config

import (
	"context"
	"time"

	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/kyverno/kyverno/pkg/engine/context/resolvers"
	"github.com/kyverno/playground/backend/pkg/utils"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
)

type Config interface {
	KubeClient() (kubernetes.Interface, error)
	DClient() (dclient.Interface, error)
	CMResolver() (engineapi.ConfigmapResolver, error)
}

type config struct {
	dClient    dclient.Interface
	kubeClient kubernetes.Interface
	cmResolver engineapi.ConfigmapResolver
}

func New(kubeConfig string) (Config, error) {
	restConfig, err := utils.RestConfig(kubeConfig)
	if err != nil {
		return nil, err
	}
	if restConfig == nil {
		return &config{}, nil
	}
	kubeClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	cmResolver, err := resolvers.NewClientBasedResolver(kubeClient)
	if err != nil {
		return nil, err
	}
	dClient, err := dclient.NewClient(context.Background(), dynamicClient, kubeClient, 15*time.Minute)
	if err != nil {
		return nil, err
	}
	return &config{
		dClient:    dClient,
		kubeClient: kubeClient,
		cmResolver: cmResolver,
	}, nil
}

func (r *config) DClient() (dclient.Interface, error) {
	return r.dClient, nil
}

func (r *config) CMResolver() (engineapi.ConfigmapResolver, error) {
	return r.cmResolver, nil
}

func (r *config) KubeClient() (kubernetes.Interface, error) {
	return r.kubeClient, nil
}
