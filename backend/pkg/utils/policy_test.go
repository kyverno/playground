package utils_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"

	"github.com/kyverno/playground/backend/data"
	"github.com/kyverno/playground/backend/pkg/resource/loader"
	"github.com/kyverno/playground/backend/pkg/utils"
)

const (
	empty             string = "../../testdata/empty.yaml"
	singleResource    string = "../../testdata/namespace.yaml"
	multiplePolicy    string = "../../testdata/multiple-policies.yaml"
	policyWithComment string = "../../testdata/multiple-policies-with-comment.yaml"
)

func Test_LoadPolicies(t *testing.T) {
	tests := []struct {
		name       string
		policies   string
		wantLoaded int
		wantErr    bool
	}{{
		name:     "invalid policy",
		policies: "../../testdata/invalid-policy.yaml",
		wantErr:  true,
	}, {
		name:     "invalid cluster policy",
		policies: "../../testdata/invalid-cluster-policy.yaml",
		wantErr:  true,
	}, {
		name:     "load no policy with empy string",
		policies: empty,
	}, {
		name:     "load invalid string",
		policies: singleResource,
		wantErr:  true,
	}, {
		name:       "load single policy",
		policies:   "../../testdata/single-policy.yaml",
		wantLoaded: 1,
	}, {
		name:       "load multiple resources",
		policies:   multiplePolicy,
		wantLoaded: 2,
	}, {
		name:       "load policy with comment",
		policies:   policyWithComment,
		wantLoaded: 1,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := os.ReadFile(tt.policies)
			require.NoError(t, err)
			loader, err := loader.New(
				openapiclient.NewComposite(
					openapiclient.NewHardcodedBuiltins("1.27"),
					openapiclient.NewLocalFiles(data.Schemas(), "schemas"),
				),
			)
			require.NoError(t, err)
			if res, err := utils.LoadPolicies(loader, bytes); (err != nil) != tt.wantErr {
				t.Errorf("loader.Policies() error = %v, wantErr %v", err, tt.wantErr)
			} else if len(res) != tt.wantLoaded {
				t.Errorf("loader.Policies() loaded amount = %v, wantLoaded %v", len(res), tt.wantLoaded)
			}
		})
	}
}

func TestToPolicyInterface(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{{
		name:    "load single policy",
		file:    "../../testdata/single-policy.yaml",
		wantErr: true,
	}, {
		name:    "load single cluster policy",
		file:    "../../testdata/single-cluster-policy.yaml",
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := os.ReadFile(tt.file)
			require.NoError(t, err)
			loader, err := loader.New(
				openapiclient.NewComposite(
					openapiclient.NewHardcodedBuiltins("1.27"),
					openapiclient.NewLocalFiles(data.Schemas(), "schemas"),
				),
			)
			require.NoError(t, err)
			resource, err := loader.Load(bytes)
			require.NoError(t, err)
			err = unstructured.SetNestedField(resource.UnstructuredContent(), "foo", "spec", "bar")
			require.NoError(t, err)
			_, err = utils.ToPolicyInterface(resource)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToPolicyInterface() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
