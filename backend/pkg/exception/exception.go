package exception

import (
	"fmt"

	policiesv1 "github.com/kyverno/api/api/policies.kyverno.io/v1"
	policiesv1alpha1 "github.com/kyverno/api/api/policies.kyverno.io/v1alpha1"
	policiesv1beta1 "github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	kyvernov2 "github.com/kyverno/kyverno/api/kyverno/v2"
	kyvernov2beta1 "github.com/kyverno/kyverno/api/kyverno/v2beta1"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/data"
	"github.com/kyverno/kyverno/ext/resource/convert"
	resourceloader "github.com/kyverno/kyverno/ext/resource/loader"
	yamlutils "github.com/kyverno/kyverno/ext/yaml"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
)

var (
	exceptionV2      = schema.GroupVersion(kyvernov2.GroupVersion).WithKind("PolicyException")
	exceptionV2beta1 = schema.GroupVersion(kyvernov2beta1.GroupVersion).WithKind("PolicyException")

	polexV1alpha1 = policiesv1alpha1.SchemeGroupVersion.WithKind("PolicyException")
	polexV1beta1  = policiesv1beta1.SchemeGroupVersion.WithKind("PolicyException")
	polexV1       = policiesv1.SchemeGroupVersion.WithKind("PolicyException")
)

func Load(content []byte) ([]*kyvernov2.PolicyException, []*policiesv1beta1.PolicyException, error) {
	fs, err := data.Crds()
	if err != nil {
		return nil, nil, err
	}

	factory, err := resourceloader.New(openapiclient.NewComposite(openapiclient.NewLocalCRDFiles(fs)))
	if err != nil {
		return nil, nil, err
	}

	documents, err := yamlutils.SplitDocuments(content)
	if err != nil {
		return nil, nil, err
	}
	var exceptions []*kyvernov2.PolicyException
	var celExceptions []*policiesv1beta1.PolicyException
	for _, document := range documents {
		gvk, untyped, err := factory.Load(document)
		if err != nil {
			return nil, nil, err
		}
		switch gvk {
		case exceptionV2, exceptionV2beta1:
			exception, err := convert.To[kyvernov2.PolicyException](untyped)
			if err != nil {
				return nil, nil, err
			}
			exceptions = append(exceptions, exception)
		case polexV1alpha1, polexV1beta1, polexV1:
			exception, err := convert.To[policiesv1beta1.PolicyException](untyped)
			if err != nil {
				return nil, nil, err
			}
			celExceptions = append(celExceptions, exception)
		default:
			return nil, nil, fmt.Errorf("policy exception type not supported %s", gvk)
		}
	}
	return exceptions, celExceptions, nil
}
