package policy

import (
	"os"
	"testing"

	"github.com/kyverno/kyverno/ext/resource/loader"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"

	"github.com/kyverno/playground/backend/data"
)

const (
	empty             string = "../../testdata/empty.yaml"
	singleResource    string = "../../testdata/namespace.yaml"
	multiplePolicy    string = "../../testdata/multiple-policies.yaml"
	policyWithComment string = "../../testdata/multiple-policies-with-comment.yaml"
	vapV1Alpha1       string = "../../testdata/vap-v1.yaml"
	vapV1Beta1        string = "../../testdata/vap-v1beta1.yaml"
	policyAndVap      string = "../../testdata/policy-and-vap.yaml"
)

func Test_Load(t *testing.T) {
	tests := []struct {
		name         string
		policies     string
		wantPolicies int
		wantVaps     int
		wantVapbs    int
		wantVpols    int
		wantIVpols   int
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
		name:         "vap v1alpha1",
		policies:     vapV1Alpha1,
		wantPolicies: 0,
		wantVaps:     1,
	}, {
		name:         "vap v1beta1",
		policies:     vapV1Beta1,
		wantPolicies: 0,
		wantVaps:     1,
	}, {
		name:         "policy and vap",
		policies:     policyAndVap,
		wantPolicies: 1,
		wantVaps:     1,
	}}

	schemas, _ := data.Schemas()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := os.ReadFile(tt.policies)
			require.NoError(t, err)
			loader, err := loader.New(
				openapiclient.NewComposite(
					openapiclient.NewGitHubBuiltins("1.31"),
					openapiclient.NewLocalSchemaFiles(schemas),
				),
			)
			require.NoError(t, err)
			if policies, vaps, vapbs, vpols, ivpols, err := Load(loader, bytes); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			} else if len(policies) != tt.wantPolicies {
				t.Errorf("Load() loaded amount = %v, wantLoaded %v", len(policies), tt.wantPolicies)
			} else if len(vaps) != tt.wantVaps {
				t.Errorf("Load() loaded amount = %v, wantLoaded %v", len(vaps), tt.wantVaps)
			} else if len(vapbs) != tt.wantVapbs {
				t.Errorf("Load() loaded amount = %v, wantLoaded %v", len(vaps), tt.wantVaps)
			} else if len(vpols) != tt.wantVpols {
				t.Errorf("Load() loaded amount = %v, wantLoaded %v", len(vpols), tt.wantVpols)
			} else if len(ivpols) != tt.wantIVpols {
				t.Errorf("Load() loaded amount = %v, wantLoaded %v", len(ivpols), tt.wantIVpols)
			}
		})
	}
}
