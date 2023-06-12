package models

import (
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type Response struct {
	// OriginalResource is the original resource as YAML string
	OriginalResource string `json:"originalResource"`
	// Resource is the original resource
	Resource unstructured.Unstructured `json:"resource"`
	// Policy is the original policy
	Policy kyvernov1.PolicyInterface `json:"policy"`
	// namespaceLabels given by policy context
	NamespaceLabels map[string]string `json:"namespaceLabels"`
	// PatchedResource is the resource patched with the engine action changes
	PatchedResource string `json:"patchedResource"`
	// PolicyResponse contains the engine policy response
	PolicyResponse PolicyResponse `json:"policyResponse"`
}
