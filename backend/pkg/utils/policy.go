package utils

import (
	"fmt"

	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	kyvernov2beta1 "github.com/kyverno/kyverno/api/kyverno/v2beta1"
	"k8s.io/api/admissionregistration/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/kyverno/playground/backend/pkg/resource/loader"
)

var (
	policyV1        = schema.GroupVersion(kyvernov1.GroupVersion).WithKind("Policy")
	policyV2        = schema.GroupVersion(kyvernov2beta1.GroupVersion).WithKind("Policy")
	clusterPolicyV1 = schema.GroupVersion(kyvernov1.GroupVersion).WithKind("ClusterPolicy")
	clusterPolicyV2 = schema.GroupVersion(kyvernov2beta1.GroupVersion).WithKind("ClusterPolicy")
	vapV1alpha1     = v1alpha1.SchemeGroupVersion.WithKind("ValidatingAdmissionPolicy")
)

func LoadPolicies(l loader.Loader, content []byte) ([]kyvernov1.PolicyInterface, []v1alpha1.ValidatingAdmissionPolicy, error) {
	untyped, err := loader.LoadResources(l, content)
	if err != nil {
		return nil, nil, err
	}
	var policies []kyvernov1.PolicyInterface
	var vaps []v1alpha1.ValidatingAdmissionPolicy
	for _, policy := range untyped {
		gvk := policy.GroupVersionKind()
		switch gvk {
		case policyV1, policyV2:
			var typed kyvernov1.Policy
			if err := runtime.DefaultUnstructuredConverter.FromUnstructuredWithValidation(policy.UnstructuredContent(), &typed, true); err != nil {
				return nil, nil, err
			}
			policies = append(policies, &typed)
		case clusterPolicyV1, clusterPolicyV2:
			var typed kyvernov1.ClusterPolicy
			if err := runtime.DefaultUnstructuredConverter.FromUnstructuredWithValidation(policy.UnstructuredContent(), &typed, true); err != nil {
				return nil, nil, err
			}
			policies = append(policies, &typed)
		case vapV1alpha1:
			var typed v1alpha1.ValidatingAdmissionPolicy
			if err := runtime.DefaultUnstructuredConverter.FromUnstructuredWithValidation(policy.UnstructuredContent(), &typed, true); err != nil {
				return nil, nil, err
			}
			vaps = append(vaps, typed)
		default:
			return nil, nil, fmt.Errorf("policy type not supported %s", gvk)
		}
	}
	return policies, vaps, nil
}
