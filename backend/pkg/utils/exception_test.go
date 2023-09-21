package utils_test

import (
	"os"
	"testing"

	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/resource/loader"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"

	"github.com/kyverno/playground/backend/data"
	"github.com/kyverno/playground/backend/pkg/utils"
)

func Test_LoadPolicyExceptions(t *testing.T) {
	tests := []struct {
		name       string
		policies   string
		wantLoaded int
		wantErr    bool
	}{{
		name:     "not a policy exception",
		policies: "../../testdata/namespace.yaml",
		wantErr:  true,
	}, {
		name:       "policy exception",
		policies:   "../../testdata/exception.yaml",
		wantLoaded: 1,
	}, {
		name:     "policy exception and policy",
		policies: "../../testdata/exception-and-policy.yaml",
		wantErr:  true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := os.ReadFile(tt.policies)
			require.NoError(t, err)
			loader, err := loader.New(
				openapiclient.NewComposite(
					openapiclient.NewLocalSchemaFiles(data.Schemas(), "schemas"),
				),
			)
			require.NoError(t, err)
			if res, err := utils.LoadPolicyExceptions(loader, bytes); (err != nil) != tt.wantErr {
				t.Errorf("loader.LoadPolicyExceptions() error = %v, wantErr %v", err, tt.wantErr)
			} else if len(res) != tt.wantLoaded {
				t.Errorf("loader.LoadPolicyExceptions() loaded amount = %v, wantLoaded %v", len(res), tt.wantLoaded)
			}
		})
	}
}
