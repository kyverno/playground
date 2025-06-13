package utils

import (
	"context"

	"github.com/kyverno/kyverno/pkg/cel/engine"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NSResolver(dClient dclient.Interface) engine.NamespaceResolver {
	return func(name string) *corev1.Namespace {
		if name == "" {
			return nil
		}
		ns, err := dClient.GetKubeClient().CoreV1().Namespaces().Get(context.TODO(), name, metav1.GetOptions{})
		if err != nil {
			return nil
		}
		return ns
	}
}
