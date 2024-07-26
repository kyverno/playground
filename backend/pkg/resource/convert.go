package resource

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func ToObjects(objects []unstructured.Unstructured) []runtime.Object {
	list := make([]runtime.Object, 0, len(objects))
	for _, obj := range objects {
		list = append(list, &obj)
	}

	return list
}
