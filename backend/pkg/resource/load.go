package resource

import (
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/resource"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/resource/loader"
	yamlutils "github.com/kyverno/kyverno/pkg/utils/yaml"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func Load[T any](l loader.Loader, content []byte) (*T, error) {
	return resource.Load[T](l, content)
}

func LoadResources(l loader.Loader, content []byte) ([]unstructured.Unstructured, error) {
	documents, err := yamlutils.SplitDocuments(content)
	if err != nil {
		return nil, err
	}
	var resources []unstructured.Unstructured
	for _, document := range documents {
		_, untyped, err := l.Load(document)
		if err != nil {
			return nil, err
		}
		resources = append(resources, untyped)
	}
	return resources, nil
}
