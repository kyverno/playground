package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/kyverno/playground/backend/pkg/api"
)

const (
	singleResource string = `apiVersion: v1
kind: Namespace
metadata:
  name: prod-bus-app1
  labels:
    purpose: production`

	singlePolicy string = `apiVersion: kyverno.io/v1
kind: ClusterPolicy
metadata:
  name: require-ns-purpose-label
spec:
  validationFailureAction: Enforce
  rules:
  - name: require-ns-purpose-label
    match:
      any:
      - resources:
          kinds:
          - Namespace
    validate:
      message: "You must have label 'purpose' with value 'production' set on all new namespaces."
      pattern:
        metadata:
          labels:
            purpose: production`
)

func Test_Serve(t *testing.T) {
	jsonBody := api.EngineRequest{
		Resources: singleResource,
		Policies:  singlePolicy,
	}

	body := new(bytes.Buffer)
	json.NewEncoder(body).Encode(jsonBody) //nolint: errcheck
	w := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Add("Content-Type", "application/json")

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler, _ := api.NewEngineHandler(nil)
	handler(c)

	if w.Result().StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(w.Result().Body) //nolint: errcheck

		t.Errorf("unexpected error: %s", buf.String())
	}
}
