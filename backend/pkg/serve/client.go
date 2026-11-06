// Package serve is a client for the experimental `kyverno serve` v1alpha1 HTTP API, which the
// Playground uses for all policy evaluation. The HTTP contract (documented in the kyverno repo at
// docs/user/cli/serve.md) is the integration point: there is no Go dependency on kyverno's CLI.
package serve

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// APIVersion is the kyverno serve API version this client speaks.
const APIVersion = "serve.cli.kyverno.io/v1alpha1"

// DefaultURL is used when neither a flag nor KYVERNO_SERVE_URL is set.
const DefaultURL = "http://127.0.0.1:8081"

// URLFromEnv returns KYVERNO_SERVE_URL, or DefaultURL.
func URLFromEnv() string {
	if u := strings.TrimSpace(os.Getenv("KYVERNO_SERVE_URL")); u != "" {
		return u
	}
	return DefaultURL
}

type EvaluationRequest struct {
	APIVersion         string `json:"apiVersion"`
	Kind               string `json:"kind"`
	Mode               string `json:"mode"`
	Payload            string `json:"payload"`
	OldPayload         string `json:"oldPayload,omitempty"`
	Policies           string `json:"policies,omitempty"`
	Exceptions         string `json:"exceptions,omitempty"`
	Context            string `json:"context,omitempty"`
	ClusterResources   string `json:"clusterResources,omitempty"`
	ParameterResources string `json:"parameterResources,omitempty"`
	TargetResources    string `json:"targetResources,omitempty"`
	// CustomResourceDefinitions apply to this request only.
	CustomResourceDefinitions string `json:"customResourceDefinitions,omitempty"`
	Operation                 string `json:"operation,omitempty"`
}

type ObjectRef struct {
	APIVersion string `json:"apiVersion,omitempty"`
	Kind       string `json:"kind,omitempty"`
	Namespace  string `json:"namespace,omitempty"`
	Name       string `json:"name,omitempty"`
}

type RuleResult struct {
	Policy             ObjectRef         `json:"policy"`
	Rule               string            `json:"rule,omitempty"`
	Type               string            `json:"type,omitempty"`
	Status             string            `json:"status"`
	Message            string            `json:"message,omitempty"`
	Resource           *ObjectRef        `json:"resource,omitempty"`
	Exceptions         []string          `json:"exceptions,omitempty"`
	Properties         map[string]string `json:"properties,omitempty"`
	PatchedResource    map[string]any    `json:"patchedResource,omitempty"`
	PatchedTarget      map[string]any    `json:"patchedTarget,omitempty"`
	GeneratedResources []map[string]any  `json:"generatedResources,omitempty"`
}

type Summary struct {
	Pass  int `json:"pass"`
	Fail  int `json:"fail"`
	Warn  int `json:"warn"`
	Error int `json:"error"`
	Skip  int `json:"skip"`
}

type EvaluationResponse struct {
	APIVersion      string       `json:"apiVersion"`
	Kind            string       `json:"kind"`
	Mode            string       `json:"mode"`
	KyvernoVersion  string       `json:"kyvernoVersion"`
	PolicySetDigest string       `json:"policySetDigest"`
	Decision        string       `json:"decision"`
	Summary         Summary      `json:"summary"`
	Results         []RuleResult `json:"results"`
	Warnings        []string     `json:"warnings,omitempty"`
	Errors          []string     `json:"errors,omitempty"`
}

// Status is the body of a non-2xx kyverno serve response.
type Status struct {
	Status  string `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error is returned when kyverno serve rejects a request (4xx) or fails to evaluate it (5xx).
type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("kyverno serve: %s (HTTP %d)", e.Message, e.Code)
}

// Client calls a kyverno serve instance.
type Client struct {
	BaseURL string
	HTTP    *http.Client
}

// NewClient returns a client for the kyverno serve instance at baseURL.
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		HTTP:    &http.Client{Timeout: 2 * time.Minute},
	}
}

// Evaluate sends one EvaluationRequest to POST /v1alpha1/evaluate.
func (c *Client) Evaluate(ctx context.Context, req EvaluationRequest) (*EvaluationResponse, error) {
	if c == nil || c.BaseURL == "" {
		return nil, errors.New("kyverno serve URL is not configured")
	}
	req.APIVersion, req.Kind = APIVersion, "EvaluationRequest"
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/v1alpha1/evaluate", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	res, err := c.HTTP.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("kyverno serve unreachable at %s: %w", c.BaseURL, err)
	}
	defer res.Body.Close()
	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		var st Status
		if json.Unmarshal(raw, &st) == nil && st.Message != "" {
			return nil, &Error{Code: res.StatusCode, Message: st.Message}
		}
		return nil, &Error{Code: res.StatusCode, Message: strings.TrimSpace(string(raw))}
	}
	var out EvaluationResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("failed to decode kyverno serve response: %w", err)
	}
	if out.APIVersion != APIVersion {
		return nil, fmt.Errorf("unexpected kyverno serve apiVersion %q, want %q", out.APIVersion, APIVersion)
	}
	return &out, nil
}
