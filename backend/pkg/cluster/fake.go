package cluster

import (
	"context"
	"errors"
	"time"

	v2 "github.com/kyverno/kyverno/api/kyverno/v2"
	kdata "github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/data"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic/fake"
	kubefake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/openapi"
	"k8s.io/client-go/restmapper"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"

	"github.com/kyverno/playground/backend/data"
	"github.com/kyverno/playground/backend/pkg/resource"
	"github.com/kyverno/playground/backend/pkg/utils"
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

func (c fakeCluster) OpenAPIClient(version string) (openapi.Client, error) {
	kubeVersion, err := utils.ParseKubeVersion(version)
	if err != nil {
		return nil, err
	}
	schemas, err := data.Schemas()
	if err != nil {
		return nil, err
	}

	return openapiclient.NewComposite(
		openapiclient.NewHardcodedBuiltins(kubeVersion),
		openapiclient.NewLocalSchemaFiles(schemas),
	), nil
}

func (c fakeCluster) DClient(objects []runtime.Object) (dclient.Interface, error) {
	objects = append([]runtime.Object(nil), objects...)
	s := runtime.NewScheme()
	listKinds := make(map[schema.GroupVersionResource]string)
	resources := make(map[schema.GroupVersionKind]schema.GroupVersionResource)
	// A custom resource's plural cannot reliably be inferred from its kind
	// (Gateway, for example, is gateways rather than gatewaies). Register CRDs
	// before their instances, including kinds with no supplied instances.
	for i, object := range objects {
		if crd, ok := object.(*apiextensionsv1.CustomResourceDefinition); ok {
			for _, version := range crd.Spec.Versions {
				gvk := schema.GroupVersionKind{Group: crd.Spec.Group, Version: version.Name, Kind: crd.Spec.Names.Kind}
				gvr := gvk.GroupVersion().WithResource(crd.Spec.Names.Plural)
				resources[gvk] = gvr
				listKind := crd.Spec.Names.ListKind
				if listKind == "" {
					listKind = gvk.Kind + "List"
				}
				listKinds[gvr] = listKind
				s.AddKnownTypeWithName(gvk, &unstructured.Unstructured{})
			}
			data, err := runtime.DefaultUnstructuredConverter.ToUnstructured(crd)
			if err != nil {
				return nil, err
			}
			objects[i] = &unstructured.Unstructured{Object: data}
		}
	}
	for _, object := range objects {
		gvk := object.GetObjectKind().GroupVersionKind()
		if _, ok := resources[gvk]; ok {
			continue
		}
		gvr, _ := meta.UnsafeGuessKindToResource(gvk)
		resources[gvk] = gvr
		listKinds[gvr] = gvk.Kind + "List"
		s.AddKnownTypeWithName(gvk, &unstructured.Unstructured{})
	}

	dyn := fake.NewSimpleDynamicClientWithCustomListKinds(s, listKinds)
	for _, object := range objects {
		accessor, err := meta.Accessor(object)
		if err != nil {
			return nil, err
		}
		// Add would infer the plural again; seed the tracker with the CRD's GVR.
		gvr := resources[object.GetObjectKind().GroupVersionKind()]
		err = dyn.Tracker().Create(gvr, object, accessor.GetNamespace())
		if apierrors.IsAlreadyExists(err) {
			// Match the original tracker.Add behavior when UPDATE requests supply
			// both the current resource and its previous version.
			err = dyn.Tracker().Update(gvr, object, accessor.GetNamespace())
		}
		if err != nil {
			return nil, err
		}
	}
	registered := make([]schema.GroupVersionResource, 0, len(resources))
	for _, gvr := range resources {
		registered = append(registered, gvr)
	}
	discovery := dclient.NewFakeDiscoveryClient(registered)
	for gvk, gvr := range resources {
		discovery.AddGVRToGVKMapping(gvr, gvk)
	}
	kclient := kubefake.NewClientset(resource.ConvertResources(objects)...)
	dClient, _ := dclient.NewClient(context.Background(), dyn, kclient, time.Hour, false, nil)
	dClient.SetDiscovery(discovery)

	return dClient, nil
}

func (c fakeCluster) RESTMapper(crds []*apiextensionsv1.CustomResourceDefinition) meta.RESTMapper {
	apiGroupResources, _ := kdata.APIGroupResources()
	for _, crd := range crds {
		apiGroupResources = append(apiGroupResources, convertCRDToAPIGroupResources(crd))
	}

	return restmapper.NewDiscoveryRESTMapper(apiGroupResources)
}

func (c fakeCluster) IsFake() bool {
	return true
}

func convertCRDToAPIGroupResources(crd *apiextensionsv1.CustomResourceDefinition) *restmapper.APIGroupResources {
	groupResources := &restmapper.APIGroupResources{
		Group: metav1.APIGroup{
			Name:             crd.Spec.Group,
			Versions:         []metav1.GroupVersionForDiscovery{},
			PreferredVersion: metav1.GroupVersionForDiscovery{},
		},
		VersionedResources: make(map[string][]metav1.APIResource),
	}

	for _, v := range crd.Spec.Versions {
		groupResources.Group.Versions = append(groupResources.Group.Versions, metav1.GroupVersionForDiscovery{
			GroupVersion: crd.Spec.Group + "/" + v.Name,
			Version:      v.Name,
		})
		if v.Storage {
			groupResources.Group.PreferredVersion.GroupVersion = crd.Spec.Group + "/" + v.Name
			groupResources.Group.PreferredVersion.Version = v.Name
		}

		groupResources.VersionedResources[v.Name] = []metav1.APIResource{
			{
				Name:         crd.Spec.Names.Plural,
				SingularName: crd.Spec.Names.Singular,
				Namespaced:   crd.Spec.Scope == apiextensionsv1.NamespaceScoped,
				Kind:         crd.Spec.Names.Kind,
				Verbs:        metav1.Verbs{"get", "list", "watch", "create", "update", "patch", "delete"},
				ShortNames:   crd.Spec.Names.ShortNames,
			},
		}
	}
	return groupResources
}
