package policy

import (
	"os"
	"testing"

	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/resource/loader"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"

	"github.com/kyverno/playground/backend/data"
)

const (
	empty             string = "../../testdata/empty.yaml"
	singleResource    string = "../../testdata/namespace.yaml"
	multiplePolicy    string = "../../testdata/multiple-policies.yaml"
	policyWithComment string = "../../testdata/multiple-policies-with-comment.yaml"
	vap               string = "../../testdata/vap.yaml"
	policyAndVap      string = "../../testdata/policy-and-vap.yaml"
)

func Test_Load(t *testing.T) {
	tests := []struct {
		name         string
		policies     string
		wantPolicies int
		wantVaps     int
		wantErr      bool
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
		name:         "load single policy",
		policies:     "../../testdata/single-policy.yaml",
		wantPolicies: 1,
	}, {
		name:         "load multiple resources",
		policies:     multiplePolicy,
		wantPolicies: 2,
	}, {
		name:         "load policy with comment",
		policies:     policyWithComment,
		wantPolicies: 1,
	}, {
		name:         "vap",
		policies:     vap,
		wantPolicies: 0,
		wantVaps:     1,
	}, {
		name:         "policy and vap",
		policies:     policyAndVap,
		wantPolicies: 1,
		wantVaps:     1,
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
			if policies, vaps, err := Load(loader, bytes); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			} else if len(policies) != tt.wantPolicies {
				t.Errorf("Load() loaded amount = %v, wantLoaded %v", len(policies), tt.wantPolicies)
			} else if len(vaps) != tt.wantVaps {
				t.Errorf("Load() loaded amount = %v, wantLoaded %v", len(vaps), tt.wantVaps)
			}
		})
	}
}
