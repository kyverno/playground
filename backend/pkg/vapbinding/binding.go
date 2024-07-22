package vapbinding

import (
	"fmt"

	"github.com/kyverno/kyverno/ext/resource/convert"
	"github.com/kyverno/kyverno/ext/resource/loader"
	"k8s.io/api/admissionregistration/v1alpha1"
	"k8s.io/api/admissionregistration/v1beta1"

	"github.com/kyverno/playground/backend/pkg/resource"
)

var (
	vapV1beta1  = v1beta1.SchemeGroupVersion.WithKind("ValidatingAdmissionPolicyBinding")
	vapV1alpha1 = v1alpha1.SchemeGroupVersion.WithKind("ValidatingAdmissionPolicyBinding")
)

func Load(l loader.Loader, content []byte) ([]v1alpha1.ValidatingAdmissionPolicyBinding, error) {
	untyped, err := resource.LoadResources(l, content)
	if err != nil {
		return nil, err
	}
	var bindings []v1alpha1.ValidatingAdmissionPolicyBinding
	for _, object := range untyped {
		gvk := object.GroupVersionKind()
		switch gvk {
		case vapV1beta1, vapV1alpha1:
			typed, err := convert.To[v1alpha1.ValidatingAdmissionPolicyBinding](object)
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
