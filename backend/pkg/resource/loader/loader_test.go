package loader

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/client-go/openapi"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
	"sigs.k8s.io/kubectl-validate/pkg/validatorfactory"
	"sigs.k8s.io/yaml"

	"github.com/kyverno/playground/backend/data"
)

type errClient struct{}

func (_ errClient) Paths() (map[string]openapi.GroupVersion, error) {
	return nil, errors.New("error")
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		client  openapi.Client
		want    Loader
		wantErr bool
	}{{
		name:    "err client",
		client:  errClient{},
		wantErr: true,
	}, {
		name:   "builtin",
		client: openapiclient.NewHardcodedBuiltins("1.27"),
		want: func() Loader {
			factory, err := validatorfactory.New(openapiclient.NewHardcodedBuiltins("1.27"))
			require.NoError(t, err)
			return &loader{
				factory: factory,
			}
		}(),
	}, {
		name:    "invalid local",
		client:  openapiclient.NewLocalSchemaFiles(data.Schemas(), "blam"),
		wantErr: true,
	}, {
		name:   "composite - no clients",
		client: openapiclient.NewComposite(),
		want: func() Loader {
			factory, err := validatorfactory.New(openapiclient.NewComposite())
			require.NoError(t, err)
			return &loader{
				factory: factory,
			}
		}(),
	}, {
		name:   "composite - err client",
		client: openapiclient.NewComposite(errClient{}),
		want: func() Loader {
			factory, err := validatorfactory.New(openapiclient.NewComposite(errClient{}))
			require.NoError(t, err)
			return &loader{
				factory: factory,
			}
		}(),
	}, {
		name:   "composite - with err client",
		client: openapiclient.NewComposite(openapiclient.NewHardcodedBuiltins("1.27"), errClient{}),
		want: func() Loader {
			factory, err := validatorfactory.New(openapiclient.NewComposite(openapiclient.NewHardcodedBuiltins("1.27"), errClient{}))
			require.NoError(t, err)
			return &loader{
				factory: factory,
			}
		}(),
	}, {
		name:   "composite - invalid local",
		client: openapiclient.NewComposite(openapiclient.NewLocalSchemaFiles(data.Schemas(), "blam")),
		want: func() Loader {
			factory, err := validatorfactory.New(openapiclient.NewComposite(openapiclient.NewLocalSchemaFiles(data.Schemas(), "blam")))
			require.NoError(t, err)
			return &loader{
				factory: factory,
			}
		}(),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.client)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_loader_Load(t *testing.T) {
	loadFile := func(path string) []byte {
		bytes, err := os.ReadFile(path)
		require.NoError(t, err)
		return bytes
	}
	newLoader := func(client openapi.Client) Loader {
		loader, err := New(client)
		require.NoError(t, err)
		return loader
	}
	toUnstructured := func(data []byte) unstructured.Unstructured {
		json, err := yaml.YAMLToJSON(data)
		require.NoError(t, err)
		var result unstructured.Unstructured
		require.NoError(t, result.UnmarshalJSON(json))
		if result.GetCreationTimestamp().Time.IsZero() {
			require.NoError(t, unstructured.SetNestedField(result.UnstructuredContent(), nil, "metadata", "creationTimestamp"))
		}
		return result
	}
	tests := []struct {
		name     string
		loader   Loader
		document []byte
		want     unstructured.Unstructured
		wantErr  bool
	}{{
		name:    "nil",
		loader:  newLoader(openapiclient.NewLocalSchemaFiles(data.Schemas(), "schemas")),
		wantErr: true,
	}, {
		name:     "empty GVK",
		loader:   newLoader(openapiclient.NewLocalSchemaFiles(data.Schemas(), "schemas")),
		document: []byte(`foo: bar`),
		wantErr:  true,
	}, {
		name:   "not yaml",
		loader: newLoader(openapiclient.NewLocalSchemaFiles(data.Schemas(), "schemas")),
		document: []byte(`
foo
  bar
  - baz`),
		wantErr: true,
	}, {
		name:     "unknown GVK",
		loader:   newLoader(openapiclient.NewLocalSchemaFiles(data.Schemas(), "schemas")),
		document: loadFile("../../../testdata/namespace.yaml"),
		wantErr:  true,
	}, {
		name:   "bad schema",
		loader: newLoader(openapiclient.NewHardcodedBuiltins("1.27")),
		document: []byte(`
apiVersion: v1
kind: Namespace
bad: field
metadata:
  name: prod-bus-app1
  labels:
    purpose: production`),
		wantErr: true,
	}, {
		name:     "ok",
		loader:   newLoader(openapiclient.NewHardcodedBuiltins("1.27")),
		document: loadFile("../../../testdata/namespace.yaml"),
		want:     toUnstructured(loadFile("../../../testdata/namespace.yaml")),
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.loader.Load(tt.document)
			if (err != nil) != tt.wantErr {
				t.Errorf("loader.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loader.Load() = %v, want %v", got, tt.want)
			}
		})
	}
}
