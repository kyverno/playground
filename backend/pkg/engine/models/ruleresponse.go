package models

import (
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type RuleResponse struct {
	// name is the rule name specified in policy
	Name string `json:"name"`
	// ruleType is the rule type (Mutation,Generation,Validation) for Kyverno Policy
	RuleType engineapi.RuleType `json:"ruleType"`
	// message is the message response from the rule application
	Message string `json:"message"`
	// status rule status
	Status engineapi.RuleStatus `json:"status"`
	// generatedResource is the generated by the generate rules of a policy
	GeneratedResource string `json:"generatedResource"`
	// patchedTarget is the patched resource for mutate.targets
	PatchedTarget *unstructured.Unstructured `json:"patchedTarget"`
	// patchedTargetParentResourceGVR is the GVR of the parent resource of the PatchedTarget. This is only populated when PatchedTarget is a subresource.
	PatchedTargetParentResourceGVR metav1.GroupVersionResource `json:"patchedTargetParentResourceGVR"`
	// patchedTargetSubresourceName is the name of the subresource which is patched, empty if the resource patched is not a subresource.
	PatchedTargetSubresourceName string `json:"patchedTargetSubresourceName"`
	// podSecurityChecks contains pod security checks (only if this is a pod security rule)
	PodSecurityChecks *engineapi.PodSecurityChecks `json:"podSecurityChecks"`
	// exception is the exception applied (if any)
	Exceptions []engineapi.GenericException `json:"exceptions"`
	// properties are the additional properties from the rule that will be added to the policy report result
	Properties map[string]string `json:"properties"`
}
