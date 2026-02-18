package policy

import (
	"fmt"

	policiesv1 "github.com/kyverno/api/api/policies.kyverno.io/v1"
	"github.com/kyverno/api/api/policies.kyverno.io/v1alpha1"
	policiesv1beta1 "github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	kyvernov2beta1 "github.com/kyverno/kyverno/api/kyverno/v2beta1"
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

	vpolV1beta1   = policiesv1beta1.SchemeGroupVersion.WithKind("ValidatingPolicy")
	nvpolV1beta1  = policiesv1beta1.SchemeGroupVersion.WithKind("NamespacedValidatingPolicy")
	ivpolV1beta1  = policiesv1beta1.SchemeGroupVersion.WithKind("ImageValidatingPolicy")
	nivpolV1beta1 = policiesv1beta1.SchemeGroupVersion.WithKind("NamespacedImageValidatingPolicy")
	dpolV1beta1   = policiesv1beta1.SchemeGroupVersion.WithKind("DeletingPolicy")
	ndpolV1beta1  = policiesv1beta1.SchemeGroupVersion.WithKind("NamespacedDeletingPolicy")
	gpolV1beta1   = policiesv1beta1.SchemeGroupVersion.WithKind("GeneratingPolicy")
	ngpolV1beta1  = policiesv1beta1.SchemeGroupVersion.WithKind("NamespacedGeneratingPolicy")
	mpolV1beta1   = policiesv1beta1.SchemeGroupVersion.WithKind("MutatingPolicy")
	nmpolV1beta1  = policiesv1beta1.SchemeGroupVersion.WithKind("NamespacedMutatingPolicy")

	vpolV1   = policiesv1.SchemeGroupVersion.WithKind("ValidatingPolicy")
	nvpolV1  = policiesv1.SchemeGroupVersion.WithKind("NamespacedValidatingPolicy")
	ivpolV1  = policiesv1.SchemeGroupVersion.WithKind("ImageValidatingPolicy")
	nivpolV1 = policiesv1.SchemeGroupVersion.WithKind("NamespacedImageValidatingPolicy")
	dpolV1   = policiesv1.SchemeGroupVersion.WithKind("DeletingPolicy")
	ndpolV1  = policiesv1.SchemeGroupVersion.WithKind("NamespacedDeletingPolicy")
	gpolV1   = policiesv1.SchemeGroupVersion.WithKind("GeneratingPolicy")
	ngpolV1  = policiesv1.SchemeGroupVersion.WithKind("NamespacedGeneratingPolicy")
	mpolV1   = policiesv1.SchemeGroupVersion.WithKind("MutatingPolicy")
	nmpolV1  = policiesv1.SchemeGroupVersion.WithKind("NamespacedMutatingPolicy")
)

func Load(l loader.Loader, content []byte) (K8sPolicies, JSONPolicies, AuthzPolicies, error) {
	k8sPolicies := K8sPolicies{}
	jsonPolicies := JSONPolicies{}
	authzPolicies := AuthzPolicies{}

	untyped, err := resource.LoadResources(l, content)
	if err != nil {
		return k8sPolicies, jsonPolicies, authzPolicies, err
	}

	for _, object := range untyped {
		gvk := object.GroupVersionKind()
		switch gvk {
		case policyV1, policyV2:
			typed, err := convert.To[kyvernov1.Policy](object)
			if err != nil {
				return k8sPolicies, jsonPolicies, authzPolicies, err
			}
			k8sPolicies.Policies = append(k8sPolicies.Policies, typed)
		case clusterPolicyV1, clusterPolicyV2:
			typed, err := convert.To[kyvernov1.ClusterPolicy](object)
			if err != nil {
				return k8sPolicies, jsonPolicies, authzPolicies, err
			}
			k8sPolicies.Policies = append(k8sPolicies.Policies, typed)
		case vapV1, vapV1beta1:
			typed, err := convert.To[v1.ValidatingAdmissionPolicy](object)
			if err != nil {
				return k8sPolicies, jsonPolicies, authzPolicies, err
			}
			k8sPolicies.ValidatingAdmissionPolicies = append(k8sPolicies.ValidatingAdmissionPolicies, *typed)
		case vapbV1, vapbV1beta1:
			typed, err := convert.To[v1.ValidatingAdmissionPolicyBinding](object)
			if err != nil {
				return k8sPolicies, jsonPolicies, authzPolicies, err
			}
			k8sPolicies.ValidatingAdmissionPolicyBindings = append(k8sPolicies.ValidatingAdmissionPolicyBindings, *typed)
		case vpolV1alpha1, vpolV1beta1, vpolV1:
			typed, err := convert.To[policiesv1beta1.ValidatingPolicy](object)
			if err != nil {
				return k8sPolicies, jsonPolicies, authzPolicies, err
			}

			mode := EvaluationMode(typed.Spec.EvaluationMode())

			switch mode {
			case Kubernetes:
				k8sPolicies.ValidatingPolicies = append(k8sPolicies.ValidatingPolicies, typed)
			case JSON:
				jsonPolicies.ValidatingPolicies = append(jsonPolicies.ValidatingPolicies, typed)
			case Envoy:
				authzPolicies.EnvoyPolicies = append(authzPolicies.EnvoyPolicies, typed)
			case HTTP:
				authzPolicies.HTTPPolicies = append(authzPolicies.HTTPPolicies, typed)
			}
		case nvpolV1beta1, nvpolV1:
			typed, err := convert.To[policiesv1beta1.NamespacedValidatingPolicy](object)
			if err != nil {
				return k8sPolicies, jsonPolicies, authzPolicies, err
			}
			mode := EvaluationMode(typed.Spec.EvaluationMode())

			switch mode {
			case Kubernetes:
				k8sPolicies.ValidatingPolicies = append(k8sPolicies.ValidatingPolicies, typed)
			case JSON:
				jsonPolicies.ValidatingPolicies = append(jsonPolicies.ValidatingPolicies, typed)
			}
		case ivpolV1alpha1, ivpolV1beta1, ivpolV1:
			typed, err := convert.To[policiesv1beta1.ImageValidatingPolicy](object)
			if err != nil {
				return k8sPolicies, jsonPolicies, authzPolicies, err
			}
			mode := EvaluationMode(typed.Spec.EvaluationMode())

			switch mode {
			case Kubernetes:
				k8sPolicies.ImageValidatingPolicies = append(k8sPolicies.ImageValidatingPolicies, typed)
			case JSON:
				jsonPolicies.ImageValidatingPolicies = append(jsonPolicies.ImageValidatingPolicies, typed)
			}
		case nivpolV1beta1, nivpolV1:
			typed, err := convert.To[policiesv1beta1.NamespacedImageValidatingPolicy](object)
			if err != nil {
				return k8sPolicies, jsonPolicies, authzPolicies, err
			}
			mode := EvaluationMode(typed.Spec.EvaluationMode())

			switch mode {
			case Kubernetes:
				k8sPolicies.ImageValidatingPolicies = append(k8sPolicies.ImageValidatingPolicies, typed)
			case JSON:
				jsonPolicies.ImageValidatingPolicies = append(jsonPolicies.ImageValidatingPolicies, typed)
			}
		case dpolV1alpha1, dpolV1beta1, dpolV1:
			typed, err := convert.To[policiesv1beta1.DeletingPolicy](object)
			if err != nil {
				return k8sPolicies, jsonPolicies, authzPolicies, err
			}
			k8sPolicies.DeletingPolicies = append(k8sPolicies.DeletingPolicies, typed)
		case ndpolV1beta1, ndpolV1:
			typed, err := convert.To[policiesv1beta1.NamespacedDeletingPolicy](object)
			if err != nil {
				return k8sPolicies, jsonPolicies, authzPolicies, err
			}
			k8sPolicies.DeletingPolicies = append(k8sPolicies.DeletingPolicies, typed)
		case gpolV1alpha1, gpolV1beta1, gpolV1:
			typed, err := convert.To[policiesv1beta1.GeneratingPolicy](object)
			if err != nil {
				return k8sPolicies, jsonPolicies, authzPolicies, err
			}
			k8sPolicies.GeneratingPolicies = append(k8sPolicies.GeneratingPolicies, typed)
		case mpolV1alpha1, mpolV1beta1, mpolV1:
			typed, err := convert.To[policiesv1beta1.MutatingPolicy](object)
			if err != nil {
				return k8sPolicies, jsonPolicies, authzPolicies, err
			}
			mode := EvaluationMode(typed.Spec.EvaluationMode())

			switch mode {
			case Kubernetes:
				k8sPolicies.MutatingPolicies = append(k8sPolicies.MutatingPolicies, typed)
			case JSON:
				jsonPolicies.MutatingPolicies = append(jsonPolicies.MutatingPolicies, typed)
			}
		case ngpolV1beta1, ngpolV1:
			typed, err := convert.To[policiesv1beta1.NamespacedGeneratingPolicy](object)
			if err != nil {
				return k8sPolicies, jsonPolicies, authzPolicies, err
			}
			k8sPolicies.GeneratingPolicies = append(k8sPolicies.GeneratingPolicies, typed)
		case nmpolV1beta1, nmpolV1:
			typed, err := convert.To[policiesv1beta1.NamespacedMutatingPolicy](object)
			if err != nil {
				return k8sPolicies, jsonPolicies, authzPolicies, err
			}
			mode := EvaluationMode(typed.Spec.EvaluationMode())

			switch mode {
			case Kubernetes:
				k8sPolicies.MutatingPolicies = append(k8sPolicies.MutatingPolicies, typed)
			case JSON:
				jsonPolicies.MutatingPolicies = append(jsonPolicies.MutatingPolicies, typed)
			}
		default:
			return k8sPolicies, jsonPolicies, authzPolicies, fmt.Errorf("policy type not supported %s", gvk)
		}
	}
	return k8sPolicies, jsonPolicies, authzPolicies, nil
}
