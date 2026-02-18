package policy

import (
	"github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	v1 "k8s.io/api/admissionregistration/v1"
)

type EvaluationMode string

const (
	Kubernetes EvaluationMode = "Kubernetes"
	JSON       EvaluationMode = "JSON"
	Envoy      EvaluationMode = "Envoy"
	HTTP       EvaluationMode = "HTTP"
)

type K8sPolicies struct {
	ValidatingAdmissionPolicies       []v1.ValidatingAdmissionPolicy
	ValidatingAdmissionPolicyBindings []v1.ValidatingAdmissionPolicyBinding

	Policies []kyvernov1.PolicyInterface

	ValidatingPolicies      []v1beta1.ValidatingPolicyLike
	ImageValidatingPolicies []v1beta1.ImageValidatingPolicyLike
	DeletingPolicies        []v1beta1.DeletingPolicyLike
	GeneratingPolicies      []v1beta1.GeneratingPolicyLike
	MutatingPolicies        []v1beta1.MutatingPolicyLike
}

func (p K8sPolicies) Length() int {
	return len(p.ValidatingAdmissionPolicies) + len(p.Policies) + len(p.ValidatingPolicies) + len(p.ImageValidatingPolicies) + len(p.DeletingPolicies) + len(p.GeneratingPolicies) + len(p.MutatingPolicies)
}

type JSONPolicies struct {
	ValidatingPolicies      []v1beta1.ValidatingPolicyLike
	ImageValidatingPolicies []v1beta1.ImageValidatingPolicyLike
	MutatingPolicies        []v1beta1.MutatingPolicyLike
}

func (p JSONPolicies) Length() int {
	return len(p.ValidatingPolicies) + len(p.ImageValidatingPolicies) + len(p.MutatingPolicies)
}

type AuthzPolicies struct {
	EnvoyPolicies []v1beta1.ValidatingPolicyLike
	HTTPPolicies  []v1beta1.ValidatingPolicyLike
}

func (p AuthzPolicies) Length() int {
	return len(p.EnvoyPolicies) + len(p.HTTPPolicies)
}
