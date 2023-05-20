package api

import (
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/kyverno/playground/backend/pkg/engine"
)

type ConfigResponse struct {
	Cluster bool   `json:"cluster"`
	Sponsor string `json:"sponsor"`
}

type ListResourcesRequest struct {
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Namespace  string            `json:"namespace"`
	Selector   map[string]string `json:"selector"`
}

type ListResourcesResponse struct {
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

type GetResourceRequest struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
	Namespace  string `json:"namespace"`
	Name       string `json:"name"`
}

type EngineRequest struct {
	Policies  string `json:"policies"`
	Resources string `json:"resources"`
	Context   string `json:"context"`
	Config    string `json:"config"`
}

type EngineResponse struct {
	Policies          []kyvernov1.PolicyInterface `json:"policies"`
	Resources         []unstructured.Unstructured `json:"resources"`
	Mutation          []engine.Response           `json:"mutation"`
	ImageVerification []engine.Response           `json:"imageVerification"`
	Validation        []engine.Response           `json:"validation"`
	Generation        []engine.Response           `json:"generation"`
}
