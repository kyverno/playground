package playground_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"

	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/crd"
	"github.com/kyverno/playground/backend/pkg/playground"
)

func targetRequest(t *testing.T) *playground.EngineRequest {
	t.Helper()
	read := func(name string) string {
		data, err := os.ReadFile("../../testdata/mpol-targets/" + name + ".yaml")
		require.NoError(t, err)
		return string(data)
	}
	return &playground.EngineRequest{
		Policies:                  read("policy"),
		Resources:                 read("route"),
		ClusterResources:          read("gateways"),
		CustomResourceDefinitions: read("crds"),
	}
}

func TestMutatingPolicyTargets(t *testing.T) {
	req := targetRequest(t)
	response, err := playground.Run(context.Background(), cluster.NewFake(), req, crd.APIConfiguration{})
	require.NoError(t, err)
	require.Len(t, response.Mutation, 2)
	for i, name := range []string{"gateway-a", "gateway-b"} {
		row := response.Mutation[i]
		require.Equal(t, "Gateway", row.Resource.GetKind())
		require.Equal(t, name, row.Resource.GetName())
		require.Equal(t, "default", row.Resource.GetNamespace())
		var original, patched unstructured.Unstructured
		require.NoError(t, yaml.Unmarshal([]byte(row.OriginalResource), &original.Object))
		require.NoError(t, yaml.Unmarshal([]byte(row.PatchedResource), &patched.Object))
		require.Equal(t, row.Resource.Object, original.Object)
		require.NotContains(t, original.GetAnnotations(), "route")
		require.Equal(t, "route-a", patched.GetAnnotations()["route"])
		delete(patched.Object["metadata"].(map[string]interface{})["annotations"].(map[string]interface{}), "route")
		require.Equal(t, original.Object, patched.Object)
	}
	require.Len(t, response.Resources, 1)
	require.Equal(t, "HTTPRoute", response.Resources[0].GetKind())
	require.Equal(t, map[string]string{"revision": "current"}, response.Resources[0].GetAnnotations())
}

func TestMutatingPolicyTargetSelection(t *testing.T) {
	tests := []struct {
		name   string
		change func(*playground.EngineRequest, map[string]interface{})
		want   []string
		value  string
	}{
		{
			name: "overlapping rules select each target once",
			change: func(_ *playground.EngineRequest, spec map[string]interface{}) {
				constraints := spec["targetMatchConstraints"].(map[string]interface{})
				rules := constraints["resourceRules"].([]interface{})
				constraints["resourceRules"] = append(rules, rules[0])
			},
			want: []string{"gateway-a", "gateway-b"}, value: "route-a",
		},
		{
			name: "target namespace selector",
			change: func(req *playground.EngineRequest, spec map[string]interface{}) {
				req.ClusterResources += "\n---\napiVersion: v1\nkind: Namespace\nmetadata:\n  name: default\n  labels:\n    environment: enabled\n"
				req.ClusterResources += "\n---\n" + strings.ReplaceAll(strings.Split(req.ClusterResources, "---")[0], "namespace: default", "namespace: excluded")
				constraints := spec["targetMatchConstraints"].(map[string]interface{})
				constraints["namespaceSelector"] = map[string]interface{}{"matchLabels": map[string]interface{}{"environment": "enabled"}}
			},
			want: []string{"gateway-a", "gateway-b"}, value: "route-a",
		},
		{
			name: "no cluster resources of target kind",
			change: func(req *playground.EngineRequest, _ map[string]interface{}) {
				req.ClusterResources = ""
			},
		},
		{
			name: "admission evaluation disabled",
			change: func(_ *playground.EngineRequest, spec map[string]interface{}) {
				spec["evaluation"].(map[string]interface{})["admission"] = map[string]interface{}{"enabled": false}
			},
		},
		{
			name: "mutate existing evaluation disabled",
			change: func(_ *playground.EngineRequest, spec map[string]interface{}) {
				spec["evaluation"].(map[string]interface{})["mutateExisting"] = map[string]interface{}{"enabled": false}
			},
		},
		{
			name: "no matching targets",
			change: func(_ *playground.EngineRequest, spec map[string]interface{}) {
				constraints := spec["targetMatchConstraints"].(map[string]interface{})
				constraints["objectSelector"] = map[string]interface{}{"matchLabels": map[string]interface{}{"selected": "missing"}}
			},
		},
		{
			name: "trigger resource does not match",
			change: func(_ *playground.EngineRequest, spec map[string]interface{}) {
				spec["matchConstraints"].(map[string]interface{})["resourceRules"].([]interface{})[0].(map[string]interface{})["resourceNames"] = []interface{}{"another-route"}
			},
		},
		{
			name: "trigger condition is false",
			change: func(_ *playground.EngineRequest, spec map[string]interface{}) {
				spec["matchConditions"] = []interface{}{map[string]interface{}{"name": "disabled", "expression": "false"}}
			},
		},
		{
			name: "update retains original trigger request",
			change: func(req *playground.EngineRequest, spec map[string]interface{}) {
				req.Operation = "UPDATE"
				req.OldResources = strings.ReplaceAll(req.Resources, "revision: current", "revision: previous")
				spec["mutations"].([]interface{})[0].(map[string]interface{})["jsonPatch"].(map[string]interface{})["expression"] = `[JSONPatch{op: "add", path: "/metadata/annotations/route", value: request.oldObject.metadata.annotations.revision + ":" + request.object.metadata.annotations.revision}]`
			},
			want: []string{"gateway-a", "gateway-b"}, value: "previous:current",
		},
		{
			name: "expression selected object",
			change: func(_ *playground.EngineRequest, spec map[string]interface{}) {
				spec["targetMatchConstraints"] = map[string]interface{}{"expression": `resource.Get("gateway.networking.k8s.io/v1", "gateways", object.metadata.namespace, "gateway-a")`}
			},
			want: []string{"gateway-a"}, value: "route-a",
		},
		{
			name: "expression selected list",
			change: func(_ *playground.EngineRequest, spec map[string]interface{}) {
				spec["targetMatchConstraints"] = map[string]interface{}{"expression": `resource.List("gateway.networking.k8s.io/v1", "gateways", object.metadata.namespace)`}
			},
			want: []string{"gateway-a", "gateway-b", "gateway-excluded"}, value: "route-a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := targetRequest(t)
			var policy map[string]interface{}
			require.NoError(t, yaml.Unmarshal([]byte(req.Policies), &policy))
			tt.change(req, policy["spec"].(map[string]interface{}))
			policyYAML, err := yaml.Marshal(policy)
			require.NoError(t, err)
			req.Policies = string(policyYAML)
			response, err := playground.Run(context.Background(), cluster.NewFake(), req, crd.APIConfiguration{})
			require.NoError(t, err)
			require.Len(t, response.Mutation, len(tt.want))
			for i, name := range tt.want {
				row := response.Mutation[i]
				require.Equal(t, "Gateway", row.Resource.GetKind())
				require.Equal(t, name, row.Resource.GetName())
				var patched unstructured.Unstructured
				require.NoError(t, yaml.Unmarshal([]byte(row.PatchedResource), &patched.Object))
				require.Equal(t, tt.value, patched.GetAnnotations()["route"], "rules: %+v", row.PolicyResponse.Rules)
			}
		})
	}
}

func TestMutatingPolicyTargetConditions(t *testing.T) {
	req := targetRequest(t)
	req.Policies = strings.Replace(req.Policies, "  mutations:", "  targetMatchConditions:\n  - name: selected-gateway\n    expression: object.metadata.name == 'gateway-a'\n  mutations:", 1)
	response, err := playground.Run(context.Background(), cluster.NewFake(), req, crd.APIConfiguration{})
	require.NoError(t, err)
	require.Len(t, response.Mutation, 2)
	for _, row := range response.Mutation {
		require.Equal(t, "Gateway", row.Resource.GetKind())
		require.Len(t, row.PolicyResponse.Rules, 1)
		if row.Resource.GetName() == "gateway-a" {
			require.Equal(t, engineapi.RuleStatusPass, row.PolicyResponse.Rules[0].Status)
			require.Contains(t, row.PatchedResource, "route: route-a")
		} else {
			require.Equal(t, "gateway-b", row.Resource.GetName())
			require.Equal(t, engineapi.RuleStatusSkip, row.PolicyResponse.Rules[0].Status)
			require.NotContains(t, row.PatchedResource, "route: route-a")
		}
	}
}

func TestOrdinaryMutatingPolicy(t *testing.T) {
	for _, emptyConstraints := range []bool{false, true} {
		t.Run(fmt.Sprintf("empty constraints %t", emptyConstraints), func(t *testing.T) {
			req := targetRequest(t)
			var policy map[string]interface{}
			require.NoError(t, yaml.Unmarshal([]byte(req.Policies), &policy))
			spec := policy["spec"].(map[string]interface{})
			delete(spec, "evaluation")
			delete(spec, "targetMatchConstraints")
			if emptyConstraints {
				spec["targetMatchConstraints"] = map[string]interface{}{}
			}
			policyYAML, err := yaml.Marshal(policy)
			require.NoError(t, err)
			req.Policies = string(policyYAML)
			response, err := playground.Run(context.Background(), cluster.NewFake(), req, crd.APIConfiguration{})
			require.NoError(t, err)
			require.Len(t, response.Mutation, 1)
			row := response.Mutation[0]
			require.Equal(t, "HTTPRoute", row.Resource.GetKind())
			require.Equal(t, "route-a", row.Resource.GetName())
			require.NotContains(t, row.OriginalResource, "route: route-a")
			require.Contains(t, row.PatchedResource, "route: route-a")
		})
	}
}
