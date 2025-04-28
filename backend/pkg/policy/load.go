package policy

import (
	"fmt"

	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	kyvernov2beta1 "github.com/kyverno/kyverno/api/kyverno/v2beta1"
	"github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"github.com/kyverno/kyverno/ext/resource/convert"
	"github.com/kyverno/kyverno/ext/resource/loader"
	v1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/api/admissionregistration/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/kyverno/playground/backend/pkg/resource"
)

var (
	policyV1        = schema.GroupVersion(kyvernov1.GroupVersion).WithKind("Policy")
	policyV2        = schema.GroupVersion(kyvernov2beta1.GroupVersion).WithKind("Policy")
	clusterPolicyV1 = schema.GroupVersion(kyvernov1.GroupVersion).WithKind("ClusterPolicy")
	clusterPolicyV2 = schema.GroupVersion(kyvernov2beta1.GroupVersion).WithKind("ClusterPolicy")
	vapV1           = v1.SchemeGroupVersion.WithKind("ValidatingAdmissionPolicy")
	vapV1beta1      = v1beta1.SchemeGroupVersion.WithKind("ValidatingAdmissionPolicy")
	vapbV1          = v1.SchemeGroupVersion.WithKind("ValidatingAdmissionPolicyBinding")
	vapbV1beta1     = v1beta1.SchemeGroupVersion.WithKind("ValidatingAdmissionPolicyBinding")
	vpolV1alpha1    = v1alpha1.SchemeGroupVersion.WithKind("ValidatingPolicy")
	ivpolV1alpha1   = v1alpha1.SchemeGroupVersion.WithKind("ImageValidatingPolicy")
)

func Load(l loader.Loader, content []byte) ([]kyvernov1.PolicyInterface, []v1.ValidatingAdmissionPolicy, []v1.ValidatingAdmissionPolicyBinding, []v1alpha1.ValidatingPolicy, []v1alpha1.ImageValidatingPolicy, error) {
	untyped, err := resource.LoadResources(l, content)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	var policies []kyvernov1.PolicyInterface
	var vaps []v1.ValidatingAdmissionPolicy
	var vapbs []v1.ValidatingAdmissionPolicyBinding
	var vpols []v1alpha1.ValidatingPolicy
	var ivpols []v1alpha1.ImageValidatingPolicy
	for _, object := range untyped {
		gvk := object.GroupVersionKind()
		switch gvk {
		case policyV1, policyV2:
			typed, err := convert.To[kyvernov1.Policy](object)
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			policies = append(policies, typed)
		case clusterPolicyV1, clusterPolicyV2:
			typed, err := convert.To[kyvernov1.ClusterPolicy](object)
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			policies = append(policies, typed)
		case vapV1, vapV1beta1:
			typed, err := convert.To[v1.ValidatingAdmissionPolicy](object)
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			vaps = append(vaps, *typed)
		case vapbV1, vapbV1beta1:
			typed, err := convert.To[v1.ValidatingAdmissionPolicyBinding](object)
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			vapbs = append(vapbs, *typed)
		case vpolV1alpha1:
			typed, err := convert.To[v1alpha1.ValidatingPolicy](object)
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			vpols = append(vpols, *typed)
		case ivpolV1alpha1:
			typed, err := convert.To[v1alpha1.ImageValidatingPolicy](object)
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}
			ivpols = append(ivpols, *typed)
		default:
			return nil, nil, nil, nil, nil, fmt.Errorf("policy type not supported %s", gvk)
		}
	}
	return policies, vaps, vapbs, vpols, ivpols, nil
}
