package resource

import "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

type Resources struct {
	Kubernetes    []unstructured.Unstructured `json:"kubernetes"`
	JSONPayload   []unstructured.Unstructured `json:"jsonPayload"`
	CheckRequests []any                       `json:"checkRequests"`
}
