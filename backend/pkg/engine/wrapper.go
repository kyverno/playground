package engine

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/kyverno/kyverno/pkg/clients/dclient"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	eventsv1 "k8s.io/client-go/kubernetes/typed/events/v1"
)

// Overwrite write actions to dry run
// prevents performing actual operations
type Client struct {
	inner   dclient.Interface
	created map[string]*unstructured.Unstructured
	mx      *sync.RWMutex
}

func (c *Client) GetKubeClient() kubernetes.Interface {
	return c.inner.GetKubeClient()
}

func (c *Client) GetEventsInterface() eventsv1.EventsV1Interface {
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
	if obj, ok := c.getObject(apiVersion, kind, namespace, name); ok {
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
		c.addObject(o)
		return o, nil
	}

	return nil, nil
}

func (c *Client) UpdateResource(_ context.Context, _ string, _ string, _ string, obj interface{}, _ bool, _ ...string) (*unstructured.Unstructured, error) {
	if o, ok := obj.(*unstructured.Unstructured); ok {
		c.addObject(o)
		return o, nil
	}

	return nil, nil
}

func (c *Client) UpdateStatusResource(_ context.Context, _ string, _ string, _ string, obj interface{}, _ bool) (*unstructured.Unstructured, error) {
	if o, ok := obj.(*unstructured.Unstructured); ok {
		c.addObject(o)
		return o, nil
	}

	return nil, nil
}

func (c *Client) ApplyResource(_ context.Context, _, _, _, _ string, obj interface{}, _ bool, _ string, _ ...string) (*unstructured.Unstructured, error) {
	if o, ok := obj.(*unstructured.Unstructured); ok {
		c.addObject(o)
		return o, nil
	}

	return nil, nil
}

func (c *Client) ApplyStatusResource(_ context.Context, _, _, _, _ string, obj interface{}, _ bool, _ string) (*unstructured.Unstructured, error) {
	if o, ok := obj.(*unstructured.Unstructured); ok {
		c.addObject(o)
		return o, nil
	}

	return nil, nil
}

func (c *Client) addObject(obj *unstructured.Unstructured) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.created[keyFromObj(obj)] = obj
}

func (c *Client) getObject(apiVersion, kind, namespace, name string) (*unstructured.Unstructured, bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	if obj, ok := c.created[keyFromValues(apiVersion, kind, namespace, name)]; ok {
		return obj, ok
	}

	return nil, false
}

func NewWrapper(client dclient.Interface) dclient.Interface {
	return &Client{
		inner:   client,
		created: make(map[string]*unstructured.Unstructured),
		mx:      new(sync.RWMutex),
	}
}

func keyFromObj(obj *unstructured.Unstructured) string {
	return fmt.Sprintf("%s/%s/%s/%s", obj.GetAPIVersion(), obj.GetKind(), obj.GetNamespace(), obj.GetName())
}

func keyFromValues(apiVersion, kind, namespace, name string) string {
	return fmt.Sprintf("%s/%s/%s/%s", apiVersion, kind, namespace, name)
}
