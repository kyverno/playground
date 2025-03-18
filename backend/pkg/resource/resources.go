package resource

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

func AppendNamespaces(resources []unstructured.Unstructured, clusterResources []unstructured.Unstructured) []runtime.Object {
	namespaces := make(map[string]runtime.Object, 0)
	results := make([]runtime.Object, 0)

	for _, res := range clusterResources {
		if res.GetKind() == "Namespace" {
			namespaces[res.GetName()] = &res
			continue
		} else {
			results = append(results, &res)
		}

		ns := res.GetNamespace()
		if _, ok := namespaces[ns]; ok {
			continue
		}

		if ns != "" {
			namespaces[ns] = &unstructured.Unstructured{
				Object: map[string]any{
					"apiVersion": "v1",
					"kind":       "Namespace",
					"metadata": map[string]any{
						"name": res.GetNamespace(),
					},
				},
			}
		}
	}

	for _, res := range resources {
		ns := res.GetNamespace()
		if _, ok := namespaces[ns]; ok {
			continue
		}

		if ns != "" {
			namespaces[ns] = &unstructured.Unstructured{
				Object: map[string]any{
					"apiVersion": "v1",
					"kind":       "Namespace",
					"metadata": map[string]any{
						"name": res.GetNamespace(),
					},
				},
			}
		}
	}

	for _, ns := range namespaces {
		results = append(results, ns)
	}

	return results
}

func ConvertResources(objects []runtime.Object) []runtime.Object {
	resources := make([]runtime.Object, 0, len(objects))

	for _, r := range objects {
		res := r.(*unstructured.Unstructured)
		if res.GetKind() == "Namespace" {
			resources = append(resources, &v1.Namespace{
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
			})
		}

		if res.GetKind() == "ConfigMap" {
			data := make(map[string]string)
			if r, ok := res.Object["data"].(map[string]any); ok {
				for k, v := range r {
					if s, ok := v.(string); ok {
						data[k] = s
					}
				}
			}

			resources = append(resources, &v1.ConfigMap{
				TypeMeta: metav1.TypeMeta{
					Kind:       "ConfigMap",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:            res.GetName(),
					Namespace:       res.GetNamespace(),
					Labels:          res.GetLabels(),
					Annotations:     res.GetAnnotations(),
					OwnerReferences: res.GetOwnerReferences(),
				},
				Data: data,
			})
		}

		if res.GetKind() == "Secret" {
			data := make(map[string][]byte)
			if r, ok := res.Object["data"].(map[string]any); ok {
				for k, v := range r {
					if s, ok := v.(string); ok {
						data[k] = []byte(s)
					}
				}
			}

			resources = append(resources, &v1.Secret{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:            res.GetName(),
					Namespace:       res.GetNamespace(),
					Labels:          res.GetLabels(),
					Annotations:     res.GetAnnotations(),
					OwnerReferences: res.GetOwnerReferences(),
				},
				Data: data,
			})
		}
	}

	return resources
}
