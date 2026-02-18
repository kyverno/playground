package resource

import (
	"encoding/json"
	"strings"

	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/resource"
	"github.com/kyverno/kyverno/ext/resource/loader"
	yamlutils "github.com/kyverno/kyverno/ext/yaml"
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

func LoadJSON(content string) ([]unstructured.Unstructured, error) {
	if strings.HasPrefix(strings.TrimSpace(content), "[") {
		var payload []any
		if err := json.Unmarshal([]byte(content), &payload); err != nil {
			return nil, err
		}

		resources := make([]unstructured.Unstructured, 0, len(payload))
		for _, item := range payload {
			if resource, ok := item.(map[string]any); ok {
				resources = append(resources, unstructured.Unstructured{Object: resource})
			}
		}

		return resources, nil
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(content), &payload); err != nil {
		return nil, err
	}

	return []unstructured.Unstructured{{Object: payload}}, nil
}
