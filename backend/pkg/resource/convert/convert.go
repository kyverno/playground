package convert

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func Into[T any](untyped unstructured.Unstructured, result *T) error {
	return runtime.DefaultUnstructuredConverter.FromUnstructuredWithValidation(untyped.UnstructuredContent(), result, true)
}

func To[T any](untyped unstructured.Unstructured) (T, error) {
	var result T
	err := Into(untyped, &result)
	return result, err
}
