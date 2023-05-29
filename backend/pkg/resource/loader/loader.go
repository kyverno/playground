package loader

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/openapi"
	"sigs.k8s.io/kubectl-validate/pkg/validatorfactory"
	"sigs.k8s.io/yaml"
)

type Loader interface {
	Load([]byte) (unstructured.Unstructured, error)
}

type loader struct {
	factory *validatorfactory.ValidatorFactory
}

func New(client openapi.Client) (Loader, error) {
	factory, err := validatorfactory.New(client)
	if err != nil {
		return nil, err
	}
	return &loader{
		factory: factory,
	}, nil
}

func (l *loader) Load(document []byte) (unstructured.Unstructured, error) {
	var metadata metav1.TypeMeta
	if err := yaml.Unmarshal(document, &metadata); err != nil {
		return unstructured.Unstructured{}, err
	}
	gvk := metadata.GetObjectKind().GroupVersionKind()
	if gvk.Empty() {
		return unstructured.Unstructured{}, fmt.Errorf("GVK cannot be empty")
	}
	validator, err := l.factory.ValidatorsForGVK(gvk)
	if err != nil {
		fmt.Println(err)
		return unstructured.Unstructured{}, err
	}
	decoder, err := validator.Decoder(gvk)
	if err != nil {
		return unstructured.Unstructured{}, err
	}
	info, ok := runtime.SerializerInfoForMediaType(decoder.SupportedMediaTypes(), runtime.ContentTypeYAML)
	if !ok {
		return unstructured.Unstructured{}, fmt.Errorf("unsupported media type %q", runtime.ContentTypeYAML)
	}
	var result unstructured.Unstructured
	_, _, err = decoder.DecoderToVersion(info.StrictSerializer, gvk.GroupVersion()).Decode(document, &gvk, &result)
	if err != nil {
		return unstructured.Unstructured{}, err
	}
	return result, err
}
