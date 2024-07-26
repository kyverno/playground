package policy

import (
	"fmt"

	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	kyvernov2beta1 "github.com/kyverno/kyverno/api/kyverno/v2beta1"
	"github.com/kyverno/kyverno/ext/resource/convert"
	"github.com/kyverno/kyverno/ext/resource/loader"
	"k8s.io/api/admissionregistration/v1alpha1"
	"k8s.io/api/admissionregistration/v1beta1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/kyverno/playground/backend/pkg/resource"
)

var (
	policyV1        = schema.GroupVersion(kyvernov1.GroupVersion).WithKind("Policy")
	policyV2        = schema.GroupVersion(kyvernov2beta1.GroupVersion).WithKind("Policy")
	clusterPolicyV1 = schema.GroupVersion(kyvernov1.GroupVersion).WithKind("ClusterPolicy")
	clusterPolicyV2 = schema.GroupVersion(kyvernov2beta1.GroupVersion).WithKind("ClusterPolicy")
	vapV1alpha1     = v1alpha1.SchemeGroupVersion.WithKind("ValidatingAdmissionPolicy")
	vapV1beta1      = v1beta1.SchemeGroupVersion.WithKind("ValidatingAdmissionPolicy")
	vapbV1alpha1    = v1alpha1.SchemeGroupVersion.WithKind("ValidatingAdmissionPolicyBinding")
	vapbV1beta1     = v1beta1.SchemeGroupVersion.WithKind("ValidatingAdmissionPolicyBinding")
)

func Load(l loader.Loader, content []byte) ([]kyvernov1.PolicyInterface, []v1alpha1.ValidatingAdmissionPolicy, []v1alpha1.ValidatingAdmissionPolicyBinding, error) {
	untyped, err := resource.LoadResources(l, content)
	if err != nil {
		return nil, nil, nil, err
	}
	var policies []kyvernov1.PolicyInterface
	var vaps []v1alpha1.ValidatingAdmissionPolicy
	var vapbs []v1alpha1.ValidatingAdmissionPolicyBinding
	for _, object := range untyped {
		gvk := object.GroupVersionKind()
		switch gvk {
		case policyV1, policyV2:
			typed, err := convert.To[kyvernov1.Policy](object)
			if err != nil {
				return nil, nil, nil, err
			}
			policies = append(policies, typed)
		case clusterPolicyV1, clusterPolicyV2:
			typed, err := convert.To[kyvernov1.ClusterPolicy](object)
			if err != nil {
				return nil, nil, nil, err
			}
			policies = append(policies, typed)
		case vapV1alpha1, vapV1beta1:
			typed, err := convert.To[v1alpha1.ValidatingAdmissionPolicy](object)
			if err != nil {
				return nil, nil, nil, err
			}
			vaps = append(vaps, *typed)
		case vapbV1alpha1, vapbV1beta1:
			typed, err := convert.To[v1alpha1.ValidatingAdmissionPolicyBinding](object)
			if err != nil {
				return nil, nil, nil, err
			}
			vapbs = append(vapbs, *typed)
		default:
			return nil, nil, nil, fmt.Errorf("policy type not supported %s", gvk)
		}
	}
	return policies, vaps, vapbs, nil
}
