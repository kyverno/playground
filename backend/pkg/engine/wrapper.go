package engine

import (
	"context"
	"fmt"
	"io"

	"github.com/kyverno/kyverno/pkg/clients/dclient"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

// Overwrite write actions to dry run
// prevents performing actual operations
type Client struct {
	inner   dclient.Interface
	created map[string]*unstructured.Unstructured
}

func (c *Client) GetKubeClient() kubernetes.Interface {
	return c.inner.GetKubeClient()
}
func (c *Client) GetEventsInterface() corev1.EventInterface {
	return c.inner.GetEventsInterface()
}

func (c *Client) GetDynamicInterface() dynamic.Interface {
	return c.inner.GetDynamicInterface()
}

func (c *Client) Discovery() dclient.IDiscovery {
	return c.inner.Discovery()
}

func (c *Client) SetDiscovery(discoveryClient dclient.IDiscovery) {
	c.inner.SetDiscovery(discoveryClient)
}

func (c *Client) RawAbsPath(ctx context.Context, path string, method string, dataReader io.Reader) ([]byte, error) {
	if method == "" {
		method = "GET"
	}

	return c.inner.RawAbsPath(ctx, path, method, dataReader)
}

func (c *Client) GetResource(ctx context.Context, apiVersion string, kind string, namespace string, name string, subresources ...string) (*unstructured.Unstructured, error) {
	if obj, ok := c.created[keyFromValues(apiVersion, kind, namespace, name)]; ok {
		return obj, nil
	}

	return c.inner.GetResource(ctx, apiVersion, kind, namespace, name, subresources...)
}

func (c *Client) PatchResource(_ context.Context, _ string, _ string, _ string, _ string, _ []byte) (*unstructured.Unstructured, error) {
	return nil, nil
}

func (c *Client) ListResource(ctx context.Context, apiVersion string, kind string, namespace string, lselector *metav1.LabelSelector) (*unstructured.UnstructuredList, error) {
	return c.inner.ListResource(ctx, apiVersion, kind, namespace, lselector)
}

func (c *Client) DeleteResource(_ context.Context, _ string, _ string, _ string, _ string, _ bool) error {
	return nil
}

func (c *Client) CreateResource(_ context.Context, _ string, _ string, _ string, obj interface{}, _ bool) (*unstructured.Unstructured, error) {
	if o, ok := obj.(*unstructured.Unstructured); ok {
		c.created[keyFromObj(o)] = o
		return o, nil
	}

	return nil, nil
}

func (c *Client) UpdateResource(_ context.Context, _ string, _ string, _ string, obj interface{}, _ bool, _ ...string) (*unstructured.Unstructured, error) {
	if o, ok := obj.(*unstructured.Unstructured); ok {
		c.created[keyFromObj(o)] = o
		return o, nil
	}

	return nil, nil
}

func (c *Client) UpdateStatusResource(_ context.Context, _ string, _ string, _ string, obj interface{}, _ bool) (*unstructured.Unstructured, error) {
	if o, ok := obj.(*unstructured.Unstructured); ok {
		c.created[keyFromObj(o)] = o
		return o, nil
	}

	return nil, nil
}

func NewWrapper(client dclient.Interface) dclient.Interface {
	return &Client{
		inner:   client,
		created: make(map[string]*unstructured.Unstructured),
	}
}

func keyFromObj(obj *unstructured.Unstructured) string {
	return fmt.Sprintf("%s/%s/%s/%s", obj.GetAPIVersion(), obj.GetKind(), obj.GetNamespace(), obj.GetName())
}

func keyFromValues(apiVersion, kind, namespace, name string) string {
	return fmt.Sprintf("%s/%s/%s/%s", apiVersion, kind, namespace, name)
}
