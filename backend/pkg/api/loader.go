package api

import (
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	yamlutils "github.com/kyverno/kyverno/pkg/utils/yaml"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
	"sigs.k8s.io/kubectl-validate/pkg/validatorfactory"
	"sigs.k8s.io/yaml"

	"github.com/kyverno/playground/backend/data"
)

const mediaType = runtime.ContentTypeYAML

type Loader struct {
	factory *validatorfactory.ValidatorFactory
}

func (l *Loader) Resources(resourcesContent string) ([]unstructured.Unstructured, error) {
	documents, err := splitDocuments([]byte(resourcesContent))
	if err != nil {
		return nil, err
	}

	var resources []unstructured.Unstructured
	for _, document := range documents {
		untyped, err := l.Unstructured(document)
		if err != nil {
			return nil, err
		}

		resources = append(resources, untyped)
	}

	return resources, nil
}

func (l *Loader) Policies(policyContent string) ([]kyvernov1.PolicyInterface, error) {
	documents, err := splitDocuments([]byte(policyContent))
	if err != nil {
		return nil, err
	}

	var policies []kyvernov1.PolicyInterface
	for _, document := range documents {
		untyped, err := l.Unstructured(document)
		if err != nil {
			return nil, err
		}

		policy, err := loadPolicy(untyped)
		if err != nil {
			return nil, err
		}

		policies = append(policies, policy)
	}

	return policies, nil
}

func (l *Loader) ConfigMap(content string) (*corev1.ConfigMap, error) {
	documents, err := splitDocuments([]byte(content))
	if err != nil {
		return nil, err
	}

	if len(documents) == 0 {
		return nil, nil
	}

	untyped, err := l.Unstructured(documents[0])
	if err != nil {
		return nil, err
	}

	configMap, err := fromUnstructured[corev1.ConfigMap](untyped)
	if err != nil {
		return nil, err
	}

	return &configMap, nil
}

func (l *Loader) Unstructured(document []byte) (unstructured.Unstructured, error) {
	var result unstructured.Unstructured
	var metadata metav1.TypeMeta
	if err := yaml.Unmarshal(document, &metadata); err != nil {
		return result, err
	}

	gvk := metadata.GetObjectKind().GroupVersionKind()
	if gvk.Empty() {
		return result, fmt.Errorf("GVK cannot be empty")
	}

	validator, err := l.factory.ValidatorsForGVK(gvk)
	if err != nil {
		return result, err
	}

	decoder, err := validator.Decoder(gvk)
	if err != nil {
		return result, err
	}

	info, ok := runtime.SerializerInfoForMediaType(decoder.SupportedMediaTypes(), mediaType)
	if !ok {
		return result, fmt.Errorf("unsupported media type %q", mediaType)
	}
	_, _, err = decoder.DecoderToVersion(info.StrictSerializer, gvk.GroupVersion()).Decode(document, &gvk, &result)

	return result, err
}

func loadPolicy(untyped unstructured.Unstructured) (kyvernov1.PolicyInterface, error) {
	kind := untyped.GetKind()
	if kind == "Policy" {
		policy, err := fromUnstructured[kyvernov1.Policy](untyped)
		if err != nil {
			return nil, err
		}

		return &policy, nil
	} else if kind == "ClusterPolicy" {
		policy, err := fromUnstructured[kyvernov1.ClusterPolicy](untyped)
		if err != nil {
			return nil, err
		}

		return &policy, nil
	}

	return nil, fmt.Errorf("invalid kind: %s", kind)
}

func fromUnstructured[T any](untyped unstructured.Unstructured) (T, error) {
	var result T
	err := runtime.DefaultUnstructuredConverter.FromUnstructured(untyped.UnstructuredContent(), &result)
	return result, err
}

func splitDocuments(yamlBytes []byte) ([][]byte, error) {
	documents, err := yamlutils.SplitDocuments(yamlBytes)
	if err != nil {
		return nil, err
	}

	var results [][]byte
	for _, document := range documents {
		onlyComments := true
		for _, line := range strings.Split(string(document), "\n") {
			if strings.TrimSpace(line) == "" {
				continue
			} else if !strings.HasPrefix(line, "#") {
				onlyComments = false
				break
			}
		}
		if !onlyComments {
			results = append(results, document)
		}
	}
	return results, nil
}

func NewLoader(kubeVersion string) (*Loader, error) {
	version, err := semver.NewVersion(kubeVersion)
	if err != nil {
		kubeVersion = "1.27"
	} else {
		kubeVersion = fmt.Sprint(version.Major(), ".", version.Minor())
	}

	factory, err := validatorfactory.New(
		openapiclient.NewComposite(
			openapiclient.NewLocalFiles(data.Schemas(), "schemas"),
			openapiclient.NewHardcodedBuiltins(kubeVersion),
		),
	)
	if err != nil {
		return nil, err
	}

	return &Loader{
		factory: factory,
	}, nil
}
