package cluster

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestFakeCustomResourceMapping(t *testing.T) {
	definition := &apiextensionsv1.CustomResourceDefinition{
		TypeMeta:   metav1.TypeMeta{APIVersion: "apiextensions.k8s.io/v1", Kind: "CustomResourceDefinition"},
		ObjectMeta: metav1.ObjectMeta{Name: "gateways.gateway.networking.k8s.io"},
		Spec: apiextensionsv1.CustomResourceDefinitionSpec{
			Group:    "gateway.networking.k8s.io",
			Scope:    apiextensionsv1.NamespaceScoped,
			Names:    apiextensionsv1.CustomResourceDefinitionNames{Kind: "Gateway", Plural: "gateways"},
			Versions: []apiextensionsv1.CustomResourceDefinitionVersion{{Name: "v1", Served: true, Storage: true}},
		},
	}
	gateway := &unstructured.Unstructured{Object: map[string]interface{}{
		"apiVersion": "gateway.networking.k8s.io/v1",
		"kind":       "Gateway",
		"metadata":   map[string]interface{}{"name": "example", "namespace": "default"},
	}}
	previous := gateway.DeepCopy()
	previous.SetLabels(map[string]string{"revision": "previous"})
	for _, tc := range []struct {
		name    string
		objects []runtime.Object
		count   int
	}{
		{name: "no instances", objects: []runtime.Object{definition}, count: 0},
		{name: "instance before definition", objects: []runtime.Object{gateway, definition}, count: 1},
		{name: "duplicate identity", objects: []runtime.Object{gateway, previous, definition}, count: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			originalInputs := append([]runtime.Object(nil), tc.objects...)
			client, err := NewFake().DClient(tc.objects)
			require.Equal(t, originalInputs, tc.objects)
			require.NoError(t, err)
			listed, err := client.ListResource(context.Background(), "gateway.networking.k8s.io/v1", "Gateway", "default", nil)
			require.NoError(t, err)
			require.Len(t, listed.Items, tc.count)
			if tc.count > 0 {
				got, err := client.GetResource(context.Background(), "gateway.networking.k8s.io/v1", "Gateway", "default", "example")
				require.NoError(t, err)
				expected := gateway
				if tc.name == "duplicate identity" {
					expected = previous
				}
				require.Equal(t, expected.Object, got.Object)
				got.SetLabels(map[string]string{"changed": "locally"})
				again, err := client.GetResource(context.Background(), "gateway.networking.k8s.io/v1", "Gateway", "default", "example")
				require.NoError(t, err)
				require.Equal(t, expected.GetLabels(), again.GetLabels())
			}
		})
	}
}
