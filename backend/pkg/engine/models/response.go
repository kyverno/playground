package models

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Policy struct {
	APIVersion  string            `json:"apiVersion"`
	Kind        string            `json:"kind"`
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type Response struct {
	// OriginalResource is the original resource as YAML string
	OriginalResource string `json:"originalResource"`
	// Resource is the original resource
	Resource unstructured.Unstructured `json:"resource"`
	// Policy is the original policy
	Policy Policy `json:"policy"`
	// namespaceLabels given by policy context
	NamespaceLabels map[string]string `json:"namespaceLabels"`
	// PatchedResource is the resource patched with the engine action changes
	PatchedResource string `json:"patchedResource"`
	// PolicyResponse contains the engine policy response
	PolicyResponse PolicyResponse `json:"policyResponse"`
}
