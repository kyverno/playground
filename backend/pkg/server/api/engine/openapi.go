package engine

import (
	"fmt"

	yamlutils "github.com/kyverno/kyverno/pkg/utils/yaml"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"k8s.io/apiextensions-apiserver/pkg/apiserver"
	structuralschema "k8s.io/apiextensions-apiserver/pkg/apiserver/schema"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/openapi"
	"k8s.io/kube-openapi/pkg/spec3"
	"k8s.io/kube-openapi/pkg/validation/spec"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient/groupversion"
)

type inmemoryClient struct {
	content []byte
}

func NewInMemory(content []byte) openapi.Client {
	return &inmemoryClient{content}
}

func (k *inmemoryClient) Paths() (map[string]openapi.GroupVersion, error) {
	codecs := serializer.NewCodecFactory(apiserver.Scheme).UniversalDecoder()
	crds := map[schema.GroupVersion]*spec3.OpenAPI{}

	documents, err := yamlutils.SplitDocuments(k.content)
	if err != nil {
		return nil, err
	}

	for _, document := range documents {
		crdObj, _, err := codecs.Decode(
			document,
			&schema.GroupVersionKind{
				Group:   "apiextensions.k8s.io",
				Version: runtime.APIVersionInternal,
				Kind:    "CustomResourceDefinition",
			}, nil)
		if err != nil {
			return nil, err
		}

		crd, ok := crdObj.(*apiextensions.CustomResourceDefinition)
		if !ok {
			return nil, fmt.Errorf("crd deserialized into incorrect type: %T", crdObj)
		}

		for _, v := range crd.Spec.Versions {
			// Convert schema to spec.Schema
			jsProps, err := apiextensions.GetSchemaForVersion(crd, v.Name)
			if err != nil {
				return nil, err
			}

			ss, err := structuralschema.NewStructural(jsProps.OpenAPIV3Schema)
			if err != nil {
				return nil, err
			}

			sch := ss.ToKubeOpenAPI()
			gvk := schema.GroupVersionKind{
				Group:   crd.Spec.Group,
				Version: v.Name,
				Kind:    crd.Spec.Names.Kind,
			}
			gvkObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(&gvk)
			if err != nil {
				return nil, err
			}

			gvr := gvk.GroupVersion().WithResource(crd.Spec.Names.Plural)
			sch.AddExtension("x-kubernetes-group-version-kind", []interface{}{gvkObj})

			key := fmt.Sprintf("%s/%s.%s", gvk.Group, gvk.Version, gvk.Kind)
			if existing, exists := crds[gvr.GroupVersion()]; exists {
				existing.Components.Schemas[key] = sch
			} else {
				crds[gvr.GroupVersion()] = &spec3.OpenAPI{
					Components: &spec3.Components{
						Schemas: map[string]*spec.Schema{
							key: sch,
						},
					},
				}
			}
		}
	}

	res := map[string]openapi.GroupVersion{}
	for k, v := range crds {
		res[fmt.Sprintf("apis/%s/%s", k.Group, k.Version)] = groupversion.NewForOpenAPI(v)
	}
	return res, nil
}
