package resource

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func GenerateNamespaces(objects []unstructured.Unstructured) []runtime.Object {
	namespaces := make(map[string]runtime.Object, 0)

	for _, res := range objects {
		if res.GetKind() != "Namespace" {
			continue
		}

		namespaces[res.GetName()] = &v1.Namespace{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Namespace",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:            res.GetName(),
				Labels:          res.GetLabels(),
				Annotations:     res.GetAnnotations(),
				OwnerReferences: res.GetOwnerReferences(),
			},
		}
	}

	for _, res := range objects {
		ns := res.GetNamespace()
		if _, ok := namespaces[ns]; ok {
			continue
		}

		if ns != "" {
			namespaces[ns] = &v1.Namespace{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Namespace",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: res.GetNamespace(),
				},
			}
		}
	}

	resources := make([]runtime.Object, 0, len(namespaces))
	for _, ns := range namespaces {
		resources = append(resources, ns)
	}

	return resources
}

func Combine(objects []unstructured.Unstructured, namespaces []runtime.Object) []runtime.Object {
	resources := make([]runtime.Object, 0, len(objects))

	for _, res := range objects {
		if res.GetKind() == "Namespace" {
			continue
		}

		resources = append(resources, &res)
	}

	return append(resources, namespaces...)
}
