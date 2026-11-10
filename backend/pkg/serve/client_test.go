package serve

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEvaluate(t *testing.T) {
	var got EvaluationRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1alpha1/evaluate", r.URL.Path)
		require.NoError(t, json.NewDecoder(r.Body).Decode(&got))
		_ = json.NewEncoder(w).Encode(EvaluationResponse{APIVersion: APIVersion, Kind: "EvaluationResponse", Decision: "Deny"})
	}))
	defer srv.Close()
	resp, err := NewClient(srv.URL+"/").Evaluate(context.Background(), EvaluationRequest{Mode: "Kubernetes", Payload: "x"})
	require.NoError(t, err)
	require.Equal(t, "Deny", resp.Decision)
	require.Equal(t, APIVersion, got.APIVersion)
	require.Equal(t, "EvaluationRequest", got.Kind)
	require.Equal(t, "Kubernetes", got.Mode)
}

func TestEvaluateErrors(t *testing.T) {
	status := http.StatusBadRequest
	body := `{"kind":"Status","status":"Failure","code":400,"message":"payload: request body is empty"}`
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
	defer srv.Close()
	c := NewClient(srv.URL)

	_, err := c.Evaluate(context.Background(), EvaluationRequest{})
	var serveErr *Error
	require.True(t, errors.As(err, &serveErr))
	require.Equal(t, 400, serveErr.Code)
	require.Equal(t, "payload: request body is empty", serveErr.Message)

	status, body = http.StatusOK, `{"apiVersion":"serve.cli.kyverno.io/v2","kind":"EvaluationResponse"}`
	_, err = c.Evaluate(context.Background(), EvaluationRequest{})
	require.ErrorContains(t, err, "unexpected kyverno serve apiVersion")

	status, body = http.StatusBadGateway, "upstream down"
	_, err = c.Evaluate(context.Background(), EvaluationRequest{})
	require.ErrorContains(t, err, "upstream down")

	_, err = NewClient("http://127.0.0.1:1").Evaluate(context.Background(), EvaluationRequest{})
	require.ErrorContains(t, err, "kyverno serve unreachable")

	var nilClient *Client
	_, err = nilClient.Evaluate(context.Background(), EvaluationRequest{})
	require.Error(t, err)
}

func TestURLFromEnv(t *testing.T) {
	t.Setenv("KYVERNO_SERVE_URL", "")
	require.Equal(t, DefaultURL, URLFromEnv())
	t.Setenv("KYVERNO_SERVE_URL", " http://serve:9000 ")
	require.Equal(t, "http://serve:9000", URLFromEnv())
}
