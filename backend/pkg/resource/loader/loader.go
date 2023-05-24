package loader

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/openapi"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
	"sigs.k8s.io/kubectl-validate/pkg/validatorfactory"
	"sigs.k8s.io/yaml"

	"github.com/kyverno/playground/backend/data"
	"github.com/kyverno/playground/backend/pkg/cluster"
)

const mediaType = runtime.ContentTypeYAML

type Loader interface {
	Load([]byte) (unstructured.Unstructured, error)
}

type loader struct {
	factory *validatorfactory.ValidatorFactory
}

func New(cluster cluster.Cluster, kubeVersion string) (Loader, error) {
	version, err := semver.NewVersion(kubeVersion)
	if err != nil {
		kubeVersion = "1.27"
	} else {
		kubeVersion = fmt.Sprint(version.Major(), ".", version.Minor())
	}
	var clients []openapi.Client
	if cluster != nil {
		clients = append(clients, openapiclient.NewKubeConfig(clientcmd.ConfigOverrides{}))
	}
	clients = append(clients, openapiclient.NewLocalFiles(data.Schemas(), "schemas"))
	clients = append(clients, openapiclient.NewHardcodedBuiltins(kubeVersion))
	factory, err := validatorfactory.New(openapiclient.NewComposite(clients...))
	if err != nil {
		return nil, err
	}
	return &loader{
		factory: factory,
	}, nil
}

func (l *loader) Load(document []byte) (unstructured.Unstructured, error) {
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
