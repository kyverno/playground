package cluster

import (
	"context"
	"errors"

	"github.com/kyverno/kyverno/api/kyverno/v2beta1"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type fakeCluster struct{}

func NewFake() Cluster {
	return fakeCluster{}
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
	return nil, errors.New("getting resource not supported in fake cluster")
}

func (c fakeCluster) PolicyExceptionSelector(namespace string, exceptions ...*v2beta1.PolicyException) engineapi.PolicyExceptionSelector {
	return NewPolicyExceptionSelector(namespace, nil, exceptions...)
}

func (c fakeCluster) DClient(objects ...unstructured.Unstructured) (dclient.Interface, error) {
	dClient := dclient.NewEmptyFakeClient()
	for i := range objects {
		res := objects[i]
		_, err := dClient.CreateResource(context.TODO(), res.GetAPIVersion(), res.GetKind(), res.GetNamespace(), &res, false)
		if err != nil {
			return nil, err
		}
	}
	return dClient, nil
}

func (c fakeCluster) IsFake() bool {
	return true
}
