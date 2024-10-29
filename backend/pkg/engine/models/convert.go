package models

import (
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/kyverno/playground/backend/pkg/utils"
	"k8s.io/api/admissionregistration/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

func convertRuleResponse(in engineapi.RuleResponse) RuleResponse {
	var generatedResource []byte

	if len(in.GeneratedResources()) > 1 {
		generatedResource, _ = yaml.Marshal(utils.Map(in.GeneratedResources(), func(ob *unstructured.Unstructured) map[string]any {
			return ob.Object
		}))
	} else if len(in.GeneratedResources()) == 1 {
		generatedResource, _ = yaml.Marshal(in.GeneratedResources()[0].Object)
	}

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
		Exceptions:        in.Exceptions(),
	}
	return out
}

func ConvertResponse(in engineapi.EngineResponse) Response {
	var patchedResource, resource []byte

	var targets []map[string]interface{}
	for _, r := range in.PolicyResponse.Rules {
		if t, _, _ := r.PatchedTarget(); t != nil {
			targets = append(targets, t.Object)
		}
	}

	if len(targets) == 0 {
		patchedResource, _ = yaml.Marshal(in.PatchedResource.Object)
		resource, _ = yaml.Marshal(in.Resource.Object)
	} else if len(targets) == 1 {
		patchedResource, _ = yaml.Marshal(targets[0])
	} else if len(targets) > 1 {
		patchedResource, _ = yaml.Marshal(targets)
	}

	out := Response{
		OriginalResource: string(resource),
		Resource:         in.Resource,
		NamespaceLabels:  in.NamespaceLabels(),
		PatchedResource:  string(patchedResource),
	}
	if in.Policy().GetType() == engineapi.KyvernoPolicyType {
		out.Policy = in.Policy().MetaObject().(kyvernov1.PolicyInterface)
	} else {
		out.ValidatingAdmissionPolicy = in.Policy().MetaObject().(*v1beta1.ValidatingAdmissionPolicy)
	}
	for _, ruleresponse := range in.PolicyResponse.Rules {
		out.PolicyResponse.Rules = append(out.PolicyResponse.Rules, convertRuleResponse(ruleresponse))
	}
	return out
}
