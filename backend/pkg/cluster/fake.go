package cluster

import (
	"context"
	"errors"
	"time"

	v2 "github.com/kyverno/kyverno/api/kyverno/v2"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/fake"
	kubefake "k8s.io/client-go/kubernetes/fake"

	"github.com/kyverno/playground/backend/pkg/resource"
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

func (c fakeCluster) PolicyExceptionSelector(namespace string, exceptions ...*v2.PolicyException) engineapi.PolicyExceptionSelector {
	return NewPolicyExceptionSelector(namespace, nil, exceptions...)
}

func (c fakeCluster) DClient(resources []runtime.Object, objects ...runtime.Object) (dclient.Interface, error) {
	s := runtime.NewScheme()
	gvr := make(map[schema.GroupVersionResource]string)
	list := []schema.GroupVersionResource{}

	for _, o := range resources {
		plural, _ := meta.UnsafeGuessKindToResource(o.GetObjectKind().GroupVersionKind())
		if _, ok := gvr[plural]; ok {
			continue
		}

		s.AddKnownTypeWithName(o.GetObjectKind().GroupVersionKind(), o)

		gvr[plural] = o.GetObjectKind().GroupVersionKind().Kind + "List"

		list = append(list, plural)
	}

	for _, o := range objects {
		plural, _ := meta.UnsafeGuessKindToResource(o.GetObjectKind().GroupVersionKind())
		if _, ok := gvr[plural]; ok {
			continue
		}

		s.AddKnownTypeWithName(o.GetObjectKind().GroupVersionKind(), o)

		gvr[plural] = o.GetObjectKind().GroupVersionKind().Kind + "List"

		list = append(list, plural)
	}

	dyn := fake.NewSimpleDynamicClientWithCustomListKinds(s, gvr, objects...)
	kclient := kubefake.NewSimpleClientset(resource.FilterNamespaces(objects)...)

	dClient, _ := dclient.NewClient(context.Background(), dyn, kclient, time.Hour)
	dClient.SetDiscovery(dclient.NewFakeDiscoveryClient(list))

	return dClient, nil
}

func (c fakeCluster) IsFake() bool {
	return true
}
