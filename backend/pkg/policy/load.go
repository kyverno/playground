package policy

import (
	"fmt"

	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	kyvernov2beta1 "github.com/kyverno/kyverno/api/kyverno/v2beta1"
	"github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	policiesv1beta1 "github.com/kyverno/kyverno/api/policies.kyverno.io/v1beta1"
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
	dpolV1alpha1    = v1alpha1.SchemeGroupVersion.WithKind("DeletingPolicy")
	gpolV1alpha1    = v1alpha1.SchemeGroupVersion.WithKind("GeneratingPolicy")
	mpolV1alpha1    = v1alpha1.SchemeGroupVersion.WithKind("MutatingPolicy")
	vpolV1beta1     = policiesv1beta1.SchemeGroupVersion.WithKind("ValidatingPolicy")
	nvpolV1beta1    = policiesv1beta1.SchemeGroupVersion.WithKind("NamespacedValidatingPolicy")
	ivpolV1beta1    = policiesv1beta1.SchemeGroupVersion.WithKind("ImageValidatingPolicy")
	nivpolV1beta1   = policiesv1beta1.SchemeGroupVersion.WithKind("NamespacedImageValidatingPolicy")
	dpolV1beta1     = policiesv1beta1.SchemeGroupVersion.WithKind("DeletingPolicy")
	ndpolV1beta1    = policiesv1beta1.SchemeGroupVersion.WithKind("NamespacedDeletingPolicy")
	gpolV1beta1     = policiesv1beta1.SchemeGroupVersion.WithKind("GeneratingPolicy")
	mpolV1beta1     = policiesv1beta1.SchemeGroupVersion.WithKind("MutatingPolicy")
)

func Load(l loader.Loader, content []byte) ([]kyvernov1.PolicyInterface, []v1.ValidatingAdmissionPolicy, []v1.ValidatingAdmissionPolicyBinding, []policiesv1beta1.ValidatingPolicyLike, []policiesv1beta1.ImageValidatingPolicyLike, []policiesv1beta1.DeletingPolicyLike, []v1alpha1.GeneratingPolicy, []v1alpha1.MutatingPolicy, error) {
	untyped, err := resource.LoadResources(l, content)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, nil, err
	}
	var policies []kyvernov1.PolicyInterface
	var vaps []v1.ValidatingAdmissionPolicy
	var vapbs []v1.ValidatingAdmissionPolicyBinding
	var vpols []policiesv1beta1.ValidatingPolicyLike
	var ivpols []policiesv1beta1.ImageValidatingPolicyLike
	var dpols []policiesv1beta1.DeletingPolicyLike
	var gpols []v1alpha1.GeneratingPolicy
	var mpols []v1alpha1.MutatingPolicy
	for _, object := range untyped {
		gvk := object.GroupVersionKind()
		switch gvk {
		case policyV1, policyV2:
			typed, err := convert.To[kyvernov1.Policy](object)
			if err != nil {
				return nil, nil, nil, nil, nil, nil, nil, nil, err
			}
			policies = append(policies, typed)
		case clusterPolicyV1, clusterPolicyV2:
			typed, err := convert.To[kyvernov1.ClusterPolicy](object)
			if err != nil {
				return nil, nil, nil, nil, nil, nil, nil, nil, err
			}
			policies = append(policies, typed)
		case vapV1, vapV1beta1:
			typed, err := convert.To[v1.ValidatingAdmissionPolicy](object)
			if err != nil {
				return nil, nil, nil, nil, nil, nil, nil, nil, err
			}
			vaps = append(vaps, *typed)
		case vapbV1, vapbV1beta1:
			typed, err := convert.To[v1.ValidatingAdmissionPolicyBinding](object)
			if err != nil {
				return nil, nil, nil, nil, nil, nil, nil, nil, err
			}
			vapbs = append(vapbs, *typed)
		case vpolV1alpha1, vpolV1beta1:
			typed, err := convert.To[policiesv1beta1.ValidatingPolicy](object)
			if err != nil {
				return nil, nil, nil, nil, nil, nil, nil, nil, err
			}
			vpols = append(vpols, typed)
		case nvpolV1beta1:
			typed, err := convert.To[policiesv1beta1.NamespacedValidatingPolicy](object)
			if err != nil {
				return nil, nil, nil, nil, nil, nil, nil, nil, err
			}
			vpols = append(vpols, typed)
		case ivpolV1alpha1, ivpolV1beta1:
			typed, err := convert.To[policiesv1beta1.ImageValidatingPolicy](object)
			if err != nil {
				return nil, nil, nil, nil, nil, nil, nil, nil, err
			}
			ivpols = append(ivpols, typed)
		case nivpolV1beta1:
			typed, err := convert.To[policiesv1beta1.NamespacedImageValidatingPolicy](object)
			if err != nil {
				return nil, nil, nil, nil, nil, nil, nil, nil, err
			}
			ivpols = append(ivpols, typed)
		case dpolV1alpha1, dpolV1beta1:
			typed, err := convert.To[policiesv1beta1.DeletingPolicy](object)
			if err != nil {
				return nil, nil, nil, nil, nil, nil, nil, nil, err
			}
			dpols = append(dpols, typed)
		case ndpolV1beta1:
			typed, err := convert.To[policiesv1beta1.NamespacedDeletingPolicy](object)
			if err != nil {
				return nil, nil, nil, nil, nil, nil, nil, nil, err
			}
			dpols = append(dpols, typed)
		case gpolV1alpha1, gpolV1beta1:
			typed, err := convert.To[v1alpha1.GeneratingPolicy](object)
			if err != nil {
				return nil, nil, nil, nil, nil, nil, nil, nil, err
			}
			gpols = append(gpols, *typed)
		case mpolV1alpha1, mpolV1beta1:
			typed, err := convert.To[v1alpha1.MutatingPolicy](object)
			if err != nil {
				return nil, nil, nil, nil, nil, nil, nil, nil, err
			}
			mpols = append(mpols, *typed)
		default:
			return nil, nil, nil, nil, nil, nil, nil, nil, fmt.Errorf("policy type not supported %s", gvk)
		}
	}
	return policies, vaps, vapbs, vpols, ivpols, dpols, gpols, mpols, nil
}
