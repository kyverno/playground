package cluster

import (
	"context"
	"errors"

	"github.com/kyverno/kyverno/api/kyverno/v2alpha1"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/kubernetes"
	kubefake "k8s.io/client-go/kubernetes/fake"
)

type fakeCluster struct {
	kubeClient kubernetes.Interface
	dClient    dclient.Interface
}

func NewFake() Cluster {
	kubeClient := kubefake.NewSimpleClientset()

	return fakeCluster{
		kubeClient: kubeClient,
	}
}

func (c fakeCluster) Kinds(_ context.Context, excludeGroups ...string) ([]Resource, error) {
	return nil, errors.New("listing kinds not supported in fake cluster")
}

func (c fakeCluster) Namespaces(ctx context.Context) ([]string, error) {
	return nil, errors.New("listing namespaces not supported in fake cluster")
}

func (c fakeCluster) Search(ctx context.Context, apiVersion string, kind string, namespace string, labels map[string]string) ([]SearchResult, error) {
	return nil, errors.New("searching resources not supported in fake cluster")
}

func (c fakeCluster) Get(ctx context.Context, apiVersion string, kind string, namespace string, name string) (*unstructured.Unstructured, error) {
	return c.dClient.GetResource(ctx, apiVersion, kind, namespace, name)
}

func (c fakeCluster) PolicyExceptionSelector(exceptions []*v2alpha1.PolicyException) engineapi.PolicyExceptionSelector {
	return NewPolicyExceptionSelector(nil, exceptions)
}

func (c fakeCluster) DClient(objects []unstructured.Unstructured) dclient.Interface {
	c.dClient = dclient.NewEmptyFakeClient()
	for _, res := range objects {
		_, _ = c.dClient.CreateResource(context.TODO(), res.GetAPIVersion(), res.GetKind(), res.GetNamespace(), &res, false)
	}
	return c.dClient
}

func (c fakeCluster) IsFake() bool {
	return true
}
