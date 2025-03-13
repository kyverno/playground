package vapbinding

import (
	"fmt"

	"github.com/kyverno/kyverno/ext/resource/convert"
	"github.com/kyverno/kyverno/ext/resource/loader"
	v1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/api/admissionregistration/v1beta1"

	"github.com/kyverno/playground/backend/pkg/resource"
)

var (
	vapV1beta1 = v1beta1.SchemeGroupVersion.WithKind("ValidatingAdmissionPolicyBinding")
	vapV1      = v1.SchemeGroupVersion.WithKind("ValidatingAdmissionPolicyBinding")
)

func Load(l loader.Loader, content []byte) ([]v1.ValidatingAdmissionPolicyBinding, error) {
	untyped, err := resource.LoadResources(l, content)
	if err != nil {
		return nil, err
	}
	var bindings []v1.ValidatingAdmissionPolicyBinding
	for _, object := range untyped {
		gvk := object.GroupVersionKind()
		switch gvk {
		case vapV1beta1, vapV1:
			typed, err := convert.To[v1.ValidatingAdmissionPolicyBinding](object)
			if err != nil {
				return nil, err
			}
			bindings = append(bindings, *typed)
		default:
			return nil, fmt.Errorf("ValidatingAdmissionPolicyBinding type not supported %s", gvk)
		}
	}
	return bindings, nil
}
