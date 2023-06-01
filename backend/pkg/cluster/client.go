package cluster

import (
	"context"
	"io"

	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type fakeClient struct {
	objects []unstructured.Unstructured
}

func (c *fakeClient) RawAbsPath(ctx context.Context, path string, method string, dataReader io.Reader) ([]byte, error) {
	return nil, nil
}

func (c *fakeClient) GetResource(ctx context.Context, apiVersion string, kind string, namespace string, name string, subresources ...string) (*unstructured.Unstructured, error) {
	return nil, nil
}

func (c *fakeClient) GetResources(group, version, kind, subresource, namespace, name string) ([]engineapi.Resource, error) {
	var results []engineapi.Resource
	for _, object := range c.objects {
		// TODO: better matching logic
		if name == object.GetName() && namespace == object.GetNamespace() {
			results = append(results, engineapi.Resource{
				Unstructured: object,
			})
		}
	}
	return results, nil
}
