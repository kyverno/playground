package cluster

import (
	"context"
	"strings"
	"time"

	"github.com/kyverno/kyverno/pkg/client/clientset/versioned"
	"github.com/kyverno/kyverno/pkg/clients/dclient"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type SearchResult struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type Resource struct {
	APIVersion    string `json:"apiVersion"`
	Kind          string `json:"kind"`
	ClusterScoped bool   `json:"clusterScoped"`
}

type Cluster interface {
	Kinds(context.Context, ...string) ([]Resource, error)
	Namespaces(context.Context) ([]string, error)
	Search(context.Context, string, string, string, map[string]string) ([]SearchResult, error)
	Get(context.Context, string, string, string, string) (*unstructured.Unstructured, error)
	KubeClient() kubernetes.Interface
	KyvernoClient() versioned.Interface
	DClient() dclient.Interface
}

type cluster struct {
	kubeClient    kubernetes.Interface
	kyvernoClient versioned.Interface
	dClient       dclient.Interface
}

func New(restConfig *rest.Config) (Cluster, error) {
	kubeClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	kyvernoClient, err := versioned.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}
	dClient, err := dclient.NewClient(context.Background(), dynamicClient, kubeClient, 15*time.Minute)
	if err != nil {
		return nil, err
	}
	return cluster{kubeClient, kyvernoClient, dClient}, nil
}

func (c cluster) Kinds(ctx context.Context, excludeGroups ...string) ([]Resource, error) {
	excluded := sets.New(excludeGroups...)
	disco := c.kubeClient.Discovery()
	_, resources, err := disco.ServerGroupsAndResources()
	var kinds []Resource
	for _, group := range resources {
		gv, err := schema.ParseGroupVersion(group.GroupVersion)
		if err != nil {
			continue
		}
		if excluded.Has(gv.Group) {
			continue
		}
		for _, resource := range group.APIResources {
			if strings.Contains(resource.Name, "/") {
				continue
			}
			verbs := sets.New(resource.Verbs...)
			if verbs.Has("get") && verbs.Has("list") {
				kinds = append(kinds, Resource{
					APIVersion:    group.GroupVersion,
					Kind:          resource.Kind,
					ClusterScoped: !resource.Namespaced,
				})
			}
		}
	}
	return kinds, err
}

func (c cluster) Namespaces(ctx context.Context) ([]string, error) {
	nsClient := c.kubeClient.CoreV1().Namespaces()
	list, err := nsClient.List(ctx, v1.ListOptions{})
	if err != nil {
		return nil, err
	}
	namespaces := make([]string, 0, len(list.Items))
	for _, item := range list.Items {
		namespaces = append(namespaces, item.GetName())
	}
	return namespaces, nil
}

func (c cluster) Search(ctx context.Context, apiVersion string, kind string, namespace string, labels map[string]string) ([]SearchResult, error) {
	var selector *v1.LabelSelector
	if labels != nil {
		selector = &v1.LabelSelector{MatchLabels: labels}
	}
	list, err := c.dClient.ListResource(ctx, apiVersion, kind, namespace, selector)
	if err != nil {
		return nil, err
	}
	resources := make([]SearchResult, 0, len(list.Items))
	for _, item := range list.Items {
		resources = append(resources, SearchResult{
			Namespace: item.GetNamespace(),
			Name:      item.GetName(),
		})
	}
	return resources, nil
}

func (c cluster) Get(ctx context.Context, apiVersion string, kind string, namespace string, name string) (*unstructured.Unstructured, error) {
	return c.dClient.GetResource(ctx, apiVersion, kind, namespace, name)
}

func (c cluster) KubeClient() kubernetes.Interface {
	return c.kubeClient
}

func (c cluster) KyvernoClient() versioned.Interface {
	return c.kyvernoClient
}

func (c cluster) DClient() dclient.Interface {
	return c.dClient
}