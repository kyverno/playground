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

	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/crd"
	"github.com/kyverno/playground/backend/pkg/playground"
)

func Test_Serve(t *testing.T) {
	singleResource, err := os.ReadFile("../../../../testdata/namespace.yaml")
	require.NoError(t, err)
	singlePolicy, err := os.ReadFile("../../../../testdata/single-policy.yaml")
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
