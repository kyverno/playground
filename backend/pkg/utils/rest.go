package utils

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func RestConfig(kubeConfig string) (*rest.Config, error) {
	if kubeConfig == "" {
		return nil, nil
	}
	return clientcmd.BuildConfigFromFlags("", kubeConfig)
}
