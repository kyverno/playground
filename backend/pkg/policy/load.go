package policy

import (
	"fmt"

	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	kyvernov2beta1 "github.com/kyverno/kyverno/api/kyverno/v2beta1"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/resource"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/resource/convert"
	"github.com/kyverno/kyverno/cmd/cli/kubectl-kyverno/resource/loader"
	"k8s.io/api/admissionregistration/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	policyV1        = schema.GroupVersion(kyvernov1.GroupVersion).WithKind("Policy")
	policyV2        = schema.GroupVersion(kyvernov2beta1.GroupVersion).WithKind("Policy")
	clusterPolicyV1 = schema.GroupVersion(kyvernov1.GroupVersion).WithKind("ClusterPolicy")
	clusterPolicyV2 = schema.GroupVersion(kyvernov2beta1.GroupVersion).WithKind("ClusterPolicy")
	vapV1alpha1     = v1alpha1.SchemeGroupVersion.WithKind("ValidatingAdmissionPolicy")
)

func Load(l loader.Loader, content []byte) ([]kyvernov1.PolicyInterface, []v1alpha1.ValidatingAdmissionPolicy, error) {
	untyped, err := resource.LoadResources(l, content)
	if err != nil {
		return nil, nil, err
	}
	var policies []kyvernov1.PolicyInterface
	var vaps []v1alpha1.ValidatingAdmissionPolicy
	for _, object := range untyped {
		gvk := object.GroupVersionKind()
		switch gvk {
		case policyV1, policyV2:
			typed, err := convert.To[kyvernov1.Policy](object)
			if err != nil {
				return nil, nil, err
			}
			policies = append(policies, typed)
		case clusterPolicyV1, clusterPolicyV2:
			typed, err := convert.To[kyvernov1.ClusterPolicy](object)
			if err != nil {
				return nil, nil, err
			}
			policies = append(policies, typed)
		case vapV1alpha1:
			typed, err := convert.To[v1alpha1.ValidatingAdmissionPolicy](object)
			if err != nil {
				return nil, nil, err
			}
			vaps = append(vaps, *typed)
		default:
			return nil, nil, fmt.Errorf("policy type not supported %s", gvk)
		}
	}
	return policies, vaps, nil
}
