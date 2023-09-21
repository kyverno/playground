package models

import (
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"k8s.io/api/admissionregistration/v1alpha1"
	"sigs.k8s.io/yaml"
)

func convertRuleResponse(in engineapi.RuleResponse) RuleResponse {
	generatedResource, _ := yaml.Marshal(in.GeneratedResource().Object)
	out := RuleResponse{
		Name:              in.Name(),
		RuleType:          in.RuleType(),
		Message:           in.Message(),
		Status:            in.Status(),
		GeneratedResource: string(generatedResource),
		// PatchedTarget *unstructured.Unstructured
		// // patchedTargetParentResourceGVR is the GVR of the parent resource of the PatchedTarget. This is only populated when PatchedTarget is a subresource.
		// PatchedTargetParentResourceGVR metav1.GroupVersionResource
		// // patchedTargetSubresourceName is the name of the subresource which is patched, empty if the resource patched is not a subresource.
		// PatchedTargetSubresourceName string
		PodSecurityChecks: in.PodSecurityChecks(),
		Exception:         in.Exception(),
	}
	return out
}

func ConvertResponse(in engineapi.EngineResponse) Response {
	patchedResource, _ := yaml.Marshal(in.PatchedResource.Object)
	resource, _ := yaml.Marshal(in.Resource.Object)
	out := Response{
		OriginalResource: string(resource),
		Resource:         in.Resource,
		NamespaceLabels:  in.NamespaceLabels(),
		PatchedResource:  string(patchedResource),
	}
	if in.Policy().GetType() == engineapi.KyvernoPolicyType {
		out.Policy = in.Policy().MetaObject().(kyvernov1.PolicyInterface)
	} else {
		out.ValidatingAdmissionPolicy = in.Policy().MetaObject().(*v1alpha1.ValidatingAdmissionPolicy)
	}
	for _, ruleresponse := range in.PolicyResponse.Rules {
		out.PolicyResponse.Rules = append(out.PolicyResponse.Rules, convertRuleResponse(ruleresponse))
	}
	return out
}
