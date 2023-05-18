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

func Test_Serve(t *testing.T) {
	jsonBody := api.Request{
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

	api.Serve(c)

	if w.Result().StatusCode != http.StatusOK {
		buf := new(bytes.Buffer)
		buf.ReadFrom(w.Result().Body) //nolint: errcheck

		t.Errorf("unexpected error: %s", buf.String())
	}
}
