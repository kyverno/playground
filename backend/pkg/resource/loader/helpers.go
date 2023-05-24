package loader

import (
	yamlutils "github.com/kyverno/kyverno/pkg/utils/yaml"
	"github.com/kyverno/playground/backend/pkg/resource/convert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Load[T any](l Loader, content []byte) (*T, error) {
	untyped, err := l.Load(content)
	if err != nil {
		return nil, err
	}
	result, err := convert.To[T](untyped)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func LoadResources(l Loader, content []byte) ([]unstructured.Unstructured, error) {
	documents, err := yamlutils.SplitDocuments(content)
	if err != nil {
		return nil, err
	}
	var resources []unstructured.Unstructured
	for _, document := range documents {
		untyped, err := l.Load(document)
		if err != nil {
			return nil, err
		}
		resources = append(resources, untyped)
	}
	return resources, nil
}
