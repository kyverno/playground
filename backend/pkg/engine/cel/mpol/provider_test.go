package mpol

import (
	"context"
	"testing"

	"github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	"github.com/kyverno/kyverno/pkg/cel/libs"
	"github.com/kyverno/kyverno/pkg/cel/policies/mpol/compiler"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/stretchr/testify/require"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apiserver/pkg/admission/plugin/policy/mutating/patch"
	"sigs.k8s.io/yaml"
)

func mutationPolicy(t *testing.T) *v1beta1.MutatingPolicy {
	t.Helper()
	var policy v1beta1.MutatingPolicy
	require.NoError(t, yaml.Unmarshal([]byte(`
apiVersion: policies.kyverno.io/v1beta1
kind: MutatingPolicy
metadata:
  name: existing
spec:
  evaluation:
    mode: JSON
    mutateExisting:
      enabled: true
  mutations:
  - patchType: JSONPatch
    jsonPatch:
      expression: '[JSONPatch{op: "add", path: "", value: {"marker": "changed"}}]'
`), &policy))
	return &policy
}

func TestProviderFetch(t *testing.T) {
	existing := mutationPolicy(t)
	ordinary := existing.DeepCopy()
	ordinary.Name = "ordinary"
	ordinary.Spec.EvaluationConfiguration.MutateExistingConfiguration = nil
	provider, err := NewProvider(compiler.NewCompiler(), []v1beta1.MutatingPolicyLike{existing, ordinary}, nil, libs.NewFakeContextProvider())
	require.NoError(t, err)
	all := provider.Fetch(context.Background(), false)
	require.Len(t, all, 2)
	require.Equal(t, "existing", all[0].Policy.GetName())
	require.Equal(t, "ordinary", all[1].Policy.GetName())
	targets := provider.Fetch(context.Background(), true)
	require.Len(t, targets, 1)
	require.Equal(t, "existing", targets[0].Policy.GetName())
}

func TestJSONProcessDoesNotDuplicateSkippedMutateExisting(t *testing.T) {
	policy := mutationPolicy(t)
	policy.Spec.MatchConditions = []admissionregistrationv1.MatchCondition{{Name: "skip", Expression: "false"}}
	results, err := JSONProcess(context.Background(), patch.NewTypeConverterManager(nil, nil), libs.NewFakeContextProvider(), unstructured.Unstructured{
		Object: map[string]interface{}{"marker": "original"},
	}, []v1beta1.MutatingPolicyLike{policy})
	require.NoError(t, err)
	require.Len(t, results, 1)
	require.Equal(t, "existing", results[0].Policy.Name)
	require.Len(t, results[0].PolicyResponse.Rules, 1)
	require.Equal(t, engineapi.RuleStatusSkip, results[0].PolicyResponse.Rules[0].Status)
}
