package engine

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"

	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/crd"
	"github.com/kyverno/playground/backend/pkg/playground"
)

// requireServe skips tests that need a real evaluation backend unless KYVERNO_SERVE_URL points at
// a `KYVERNO_EXPERIMENTAL=true kyverno serve --allow-request-policies` instance.
func requireServe(t *testing.T) {
	t.Helper()
	url := os.Getenv("KYVERNO_SERVE_URL")
	if url == "" {
		t.Skip("set KYVERNO_SERVE_URL to a `kyverno serve --allow-request-policies` instance")
	}
	playground.SetServeURL(url)
}

func Test_Serve(t *testing.T) {
	requireServe(t)
	singleResource, err := os.ReadFile("../../../../testdata/namespace.yaml")
	require.NoError(t, err)
	singlePolicy, err := os.ReadFile("../../../../testdata/require-ns-purpose-label.yaml")
	require.NoError(t, err)

	body := new(bytes.Buffer)
	require.NoError(t, json.NewEncoder(body).Encode(playground.EngineRequest{
		Resources: string(singleResource),
		Policies:  string(singlePolicy),
	}))

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/", body)
	c.Request.Header.Add("Content-Type", "application/json")

	handler, err := newEngineHandler(cluster.NewFake(), crd.APIConfiguration{})
	require.NoError(t, err)
	handler(c)

	if w.Result().StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(w.Result().Body) //nolint: errcheck
		t.Errorf("unexpected error: %s", buf.String())
	}
}

// TestServeLegacyPolicyRejected checks legacy kyverno.io policies get kyverno's 1.20 migration
// error as a 400, not an empty result.
func TestServeLegacyPolicyRejected(t *testing.T) {
	requireServe(t)
	resource, err := os.ReadFile("../../../../testdata/namespace.yaml")
	require.NoError(t, err)
	legacy, err := os.ReadFile("../../../../testdata/single-policy.yaml")
	require.NoError(t, err)
	body := new(bytes.Buffer)
	require.NoError(t, json.NewEncoder(body).Encode(playground.EngineRequest{Resources: string(resource), Policies: string(legacy)}))
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", body)
	ctx.Request.Header.Set("Content-Type", "application/json")
	handler, err := newEngineHandler(cluster.NewFake(), crd.APIConfiguration{})
	require.NoError(t, err)
	handler(ctx)
	require.Equal(t, http.StatusBadRequest, recorder.Code)
	require.Contains(t, recorder.Body.String(), "migration-to-cel")
}

func TestServeMutatingPolicyTargets(t *testing.T) {
	requireServe(t)
	read := func(name string) string {
		t.Helper()
		data, err := os.ReadFile("../../../../testdata/mpol-targets/" + name + ".yaml")
		require.NoError(t, err)
		return string(data)
	}
	body := new(bytes.Buffer)
	require.NoError(t, json.NewEncoder(body).Encode(playground.EngineRequest{
		Policies:                  read("policy"),
		Resources:                 read("route"),
		ClusterResources:          read("gateways"),
		CustomResourceDefinitions: read("crds"),
	}))
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	ctx.Request = httptest.NewRequest(http.MethodPost, "/", body)
	ctx.Request.Header.Set("Content-Type", "application/json")
	handler, err := newEngineHandler(cluster.NewFake(), crd.APIConfiguration{})
	require.NoError(t, err)
	handler(ctx)
	require.Equal(t, http.StatusOK, recorder.Code, recorder.Body.String())

	// Decode the actual wire fields to keep this assertion independent of the
	// backend response struct's JSON tags.
	var response struct {
		Mutation []struct {
			Resource         unstructured.Unstructured `json:"resource"`
			OriginalResource string                    `json:"originalResource"`
			PatchedResource  string                    `json:"patchedResource"`
		} `json:"mutation"`
	}
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &response))
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
		require.Equal(t, "Gateway", patched.GetKind())
		require.Equal(t, name, patched.GetName())
		require.Equal(t, "route-a", patched.GetAnnotations()["route"])
	}
}
