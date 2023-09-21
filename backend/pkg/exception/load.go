package exception

import (
	"fmt"

	kyvernov2alpha1 "github.com/kyverno/kyverno/api/kyverno/v2alpha1"
	kyvernov2beta1 "github.com/kyverno/kyverno/api/kyverno/v2beta1"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/resource"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/resource/convert"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/resource/loader"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	exceptionV1 = schema.GroupVersion(kyvernov2alpha1.GroupVersion).WithKind("PolicyException")
	exceptionV2 = schema.GroupVersion(kyvernov2beta1.GroupVersion).WithKind("PolicyException")
)

func Load(l loader.Loader, content []byte) ([]*kyvernov2alpha1.PolicyException, error) {
	untyped, err := resource.LoadResources(l, content)
	if err != nil {
		return nil, err
	}
	var exceptions []*kyvernov2alpha1.PolicyException
	for _, object := range untyped {
		gvk := object.GroupVersionKind()
		switch gvk {
		case exceptionV1, exceptionV2:
			exception, err := convert.To[kyvernov2alpha1.PolicyException](object)
			if err != nil {
				return nil, err
			}
			exceptions = append(exceptions, exception)
		default:
			return nil, fmt.Errorf("policy exception type not supported %s", gvk)
		}
	}
	return exceptions, nil
}
