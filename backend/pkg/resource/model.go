package resource

import (
	"github.com/kyverno/kyverno-authz/pkg/cel/libs/authz/http"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Resources struct {
	Kubernetes    []unstructured.Unstructured `json:"kubernetes"`
	JSONPayload   []unstructured.Unstructured `json:"jsonPayload"`
	CheckRequests []any                       `json:"checkRequests"`
}

type (
	header = map[string][]string
	query  = map[string][]string
)

type CheckRequest struct {
	Attributes CheckRequestAttributes `json:"attributes"`
}

type CheckRequestAttributes struct {
	Method        string `json:"method"`
	Header        header `json:"header"`
	Host          string `json:"host"`
	Protocol      string `json:"protocol"`
	ContentLength int64  `json:"contentLength"`
	Body          []byte `json:"body"`
	Scheme        string `json:"scheme"`
	Path          string `json:"path"`
	Query         query  `json:"query"`
	Fragment      string `json:"fragment"`
}

func (c CheckRequest) ToAuthz() *http.CheckRequest {
	return &http.CheckRequest{
		Attributes: http.CheckRequestAttributes{
			Method:        c.Attributes.Method,
			Header:        c.Attributes.Header,
			Host:          c.Attributes.Host,
			Protocol:      c.Attributes.Protocol,
			ContentLength: c.Attributes.ContentLength,
			Body:          c.Attributes.Body,
			Scheme:        c.Attributes.Scheme,
			Path:          c.Attributes.Path,
			Query:         c.Attributes.Query,
			Fragment:      c.Attributes.Fragment,
		},
	}
}
