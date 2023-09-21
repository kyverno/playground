package utils

import (
	kyvernov2alpha1 "github.com/kyverno/kyverno/api/kyverno/v2alpha1"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/resource/convert"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/resource/loader"

	"github.com/kyverno/playground/backend/pkg/resource"
)

func LoadPolicyExceptions(l loader.Loader, content []byte) ([]*kyvernov2alpha1.PolicyException, error) {
	untyped, err := resource.LoadResources(l, content)
	if err != nil {
		return nil, err
	}
	var exceptions []*kyvernov2alpha1.PolicyException
	for _, object := range untyped {
		exception, err := convert.To[kyvernov2alpha1.PolicyException](object)
		if err != nil {
			return nil, err
		}
		exceptions = append(exceptions, exception)
	}
	return exceptions, nil
}
