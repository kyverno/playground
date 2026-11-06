package playground

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/stretchr/testify/require"

	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/crd"
	"github.com/kyverno/playground/backend/pkg/serve"
)

// fakeServe records the requests the Playground sends and answers with the given responses.
func fakeServe(t *testing.T, respond func(serve.EvaluationRequest) (int, any)) *[]serve.EvaluationRequest {
	t.Helper()
	var requests []serve.EvaluationRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req serve.EvaluationRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		requests = append(requests, req)
		code, body := respond(req)
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(body)
	}))
	t.Cleanup(srv.Close)
	SetServeURL(srv.URL)
	t.Cleanup(func() { SetServeURL(serve.URLFromEnv()) })
	return &requests
}

const vpol = `apiVersion: policies.kyverno.io/v1
kind: ValidatingPolicy
metadata:
  name: check-labels
  labels:
    team: platform
spec:
  validations:
  - expression: "true"
`

const deployments = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: a
  namespace: default
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: b
  namespace: default
`

func TestRunKubernetes(t *testing.T) {
	requests := fakeServe(t, func(serve.EvaluationRequest) (int, any) {
		return 200, serve.EvaluationResponse{APIVersion: serve.APIVersion, Results: []serve.RuleResult{
			{Policy: serve.ObjectRef{APIVersion: "policies.kyverno.io/v1", Kind: "ValidatingPolicy", Name: "check-labels"}, Type: "Validation", Status: "fail", Message: "nope",
				Resource: &serve.ObjectRef{APIVersion: "apps/v1", Kind: "Deployment", Namespace: "default", Name: "b"}},
			{Policy: serve.ObjectRef{APIVersion: "policies.kyverno.io/v1", Kind: "ValidatingPolicy", Name: "check-labels"}, Type: "Validation", Status: "pass",
				Resource: &serve.ObjectRef{APIVersion: "apps/v1", Kind: "Deployment", Namespace: "default", Name: "a"}},
			{Policy: serve.ObjectRef{Kind: "MutatingPolicy", Name: "label-gw"}, Type: "Mutation", Status: "pass",
				Resource:      &serve.ObjectRef{APIVersion: "v1", Kind: "ConfigMap", Namespace: "default", Name: "cm"},
				PatchedTarget: map[string]any{"apiVersion": "v1", "kind": "ConfigMap", "metadata": map[string]any{"name": "cm", "namespace": "default", "labels": map[string]any{"x": "z"}}}},
			{Policy: serve.ObjectRef{Kind: "GeneratingPolicy", Name: "gen"}, Type: "Generation", Status: "pass",
				Resource:           &serve.ObjectRef{APIVersion: "apps/v1", Kind: "Deployment", Namespace: "default", Name: "a"},
				GeneratedResources: []map[string]any{{"apiVersion": "v1", "kind": "Secret"}}},
		}}
	})
	req := &EngineRequest{
		Policies:         vpol,
		VAPBindings:      "apiVersion: admissionregistration.k8s.io/v1\nkind: ValidatingAdmissionPolicyBinding\nmetadata:\n  name: b\n",
		Resources:        deployments,
		OldResources:     "apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: a\n  namespace: default\n",
		ClusterResources: "apiVersion: v1\nkind: ConfigMap\nmetadata:\n  name: cm\n  namespace: default\n",
		PolicyExceptions: "apiVersion: policies.kyverno.io/v1\nkind: PolicyException\nmetadata:\n  name: e\n",
		Context:          "context:\n  operation: UPDATE\n",
	}
	resp, err := Run(context.Background(), cluster.NewFake(), req, crd.APIConfiguration{})
	require.NoError(t, err)

	require.Len(t, *requests, 1)
	sent := (*requests)[0]
	require.Equal(t, "Kubernetes", sent.Mode)
	require.Equal(t, deployments, sent.Payload)
	require.Equal(t, req.OldResources, sent.OldPayload)
	require.Contains(t, sent.Policies, "ValidatingAdmissionPolicyBinding")
	require.Equal(t, req.PolicyExceptions, sent.Exceptions)
	require.Equal(t, "UPDATE", sent.Operation)
	require.Equal(t, req.ClusterResources, sent.ClusterResources)
	require.Equal(t, req.ClusterResources, sent.TargetResources)
	require.Contains(t, sent.Context, "kind: Context")
	require.Contains(t, sent.Context, "name: cm")

	require.Len(t, resp.Resources, 2)
	// payload order is kept (a before b) whatever order serve returned
	require.Len(t, resp.Validation, 2)
	require.Equal(t, "a", resp.Validation[0].Resource.GetName())
	require.Equal(t, "b", resp.Validation[1].Resource.GetName())
	require.Equal(t, engineapi.RuleStatusFail, resp.Validation[1].PolicyResponse.Rules[0].Status)
	require.Equal(t, "nope", resp.Validation[1].PolicyResponse.Rules[0].Message)
	require.Equal(t, "Kubernetes", resp.Validation[0].Policy.Mode)
	require.Equal(t, map[string]string{"team": "platform"}, resp.Validation[0].Policy.Labels)
	// the mutateExisting target is shown with its original and patched state
	require.Len(t, resp.Mutation, 1)
	require.Equal(t, "ConfigMap", resp.Mutation[0].Resource.GetKind())
	require.NotContains(t, resp.Mutation[0].OriginalResource, "x: z")
	require.Contains(t, resp.Mutation[0].PatchedResource, "x: z")
	require.NotNil(t, resp.Mutation[0].PolicyResponse.Rules[0].PatchedTarget)
	require.Len(t, resp.Generation, 1)
	require.Contains(t, resp.Generation[0].PolicyResponse.Rules[0].GeneratedResource, "kind: Secret")
}

func TestRunPayloadModes(t *testing.T) {
	tests := []struct {
		name     string
		policy   string
		payload  string
		mode     string
		requests int
	}{
		{"json array", "apiVersion: policies.kyverno.io/v1\nkind: ValidatingPolicy\nmetadata:\n  name: j\nspec:\n  evaluation:\n    mode: JSON\n", `[{"a":1},{"b":2}]`, "JSON", 2},
		{"envoy", "apiVersion: policies.kyverno.io/v1alpha1\nkind: ValidatingPolicy\nmetadata:\n  name: e\nspec:\n  evaluation:\n    mode: Envoy\n", `{"attributes":{"request":{"http":{"path":"/"}}}}`, "Envoy", 1},
		{"http", "apiVersion: policies.kyverno.io/v1alpha1\nkind: ValidatingPolicy\nmetadata:\n  name: h\nspec:\n  evaluation:\n    mode: HTTP\n", `[{"attributes":{"method":"GET"}},{"attributes":{"method":"POST"}}]`, "HTTP", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requests := fakeServe(t, func(serve.EvaluationRequest) (int, any) {
				return 200, serve.EvaluationResponse{APIVersion: serve.APIVersion, Results: []serve.RuleResult{{Policy: serve.ObjectRef{Name: "p"}, Type: "Validation", Status: "pass"}}}
			})
			resp, err := Run(context.Background(), cluster.NewFake(), &EngineRequest{Policies: tt.policy, Resources: tt.payload}, crd.APIConfiguration{})
			require.NoError(t, err)
			require.Len(t, *requests, tt.requests)
			for _, r := range *requests {
				require.Equal(t, tt.mode, r.Mode)
				require.Empty(t, r.OldPayload)
				require.Empty(t, r.TargetResources)
			}
			require.Len(t, resp.Validation, tt.requests)
		})
	}
}

func TestRunErrors(t *testing.T) {
	fakeServe(t, func(serve.EvaluationRequest) (int, any) {
		return 400, serve.Status{Status: "Failure", Code: 400, Message: "policies: legacy kind"}
	})
	_, err := Run(context.Background(), cluster.NewFake(), &EngineRequest{Policies: vpol, Resources: deployments}, crd.APIConfiguration{})
	require.EqualError(t, err, "policies: legacy kind")

	for _, ctx := range []string{
		"context:\n  username: alice\n",
		"context:\n  namespaceLabels:\n    a: b\n",
		"variables:\n  foo: bar\n",
	} {
		_, err := Run(context.Background(), cluster.NewFake(), &EngineRequest{Policies: vpol, Resources: deployments, Context: ctx}, crd.APIConfiguration{})
		require.ErrorContains(t, err, "not supported by the kyverno serve v1alpha1 evaluation backend")
	}
}

func TestRunNothingToEvaluate(t *testing.T) {
	requests := fakeServe(t, func(serve.EvaluationRequest) (int, any) { return 500, nil })
	for _, req := range []*EngineRequest{
		{Policies: "", Resources: deployments},
		{Policies: vpol, Resources: "  "},
		{Policies: "apiVersion: policies.kyverno.io/v1\nkind: PolicyException\nmetadata:\n  name: e\n", Resources: deployments},
	} {
		resp, err := Run(context.Background(), cluster.NewFake(), req, crd.APIConfiguration{})
		require.NoError(t, err)
		require.Nil(t, resp.Results)
	}
	require.Empty(t, *requests)
}

func TestPolicyModePrecedence(t *testing.T) {
	p, err := parsePolicies("apiVersion: policies.kyverno.io/v1\nkind: ValidatingPolicy\nmetadata:\n  name: h\nspec:\n  evaluation:\n    mode: HTTP\n---\napiVersion: policies.kyverno.io/v1\nkind: ValidatingPolicy\nmetadata:\n  name: j\nspec:\n  evaluation:\n    mode: JSON\n")
	require.NoError(t, err)
	require.Equal(t, modeJSON, p.mode())
	p, err = parsePolicies(vpol + "---\n" + "apiVersion: admissionregistration.k8s.io/v1\nkind: ValidatingAdmissionPolicy\nmetadata:\n  name: v\n")
	require.NoError(t, err)
	require.Equal(t, modeKubernetes, p.mode())
	require.Equal(t, "", p.describe(serve.ObjectRef{Kind: "ValidatingAdmissionPolicy", Name: "v"}).Mode)
	_, err = parsePolicies("a: [")
	require.Error(t, err)
}
