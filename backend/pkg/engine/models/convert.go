package models

import (
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"

	"github.com/kyverno/playground/backend/pkg/utils"
)

func convertRuleResponse(policy string, in engineapi.RuleResponse) RuleResponse {
	var generatedResource []byte

	if len(in.GeneratedResources()) > 1 {
		generatedResource, _ = yaml.Marshal(utils.Map(in.GeneratedResources(), func(ob *unstructured.Unstructured) map[string]any {
			return ob.Object
		}))
	} else if len(in.GeneratedResources()) == 1 {
		generatedResource, _ = yaml.Marshal(in.GeneratedResources()[0].Object)
	}

	name := in.Name()
	if name == "" {
		name = policy
	}

	properties := make(map[string]string)
	for k, v := range in.Properties() {
		if v == "" {
			continue
		}

		properties[k] = v
	}

	out := RuleResponse{
		Name:              name,
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
		Properties:        properties,
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
	if in.Policy().AsKyvernoPolicy() != nil {
		out.Policy = in.Policy().AsKyvernoPolicy()
	} else if in.Policy().AsValidatingPolicy() != nil {
		out.ValidatingPolicy = in.Policy().AsValidatingPolicy()
	} else if in.Policy().AsValidatingAdmissionPolicy() != nil {
		out.ValidatingAdmissionPolicy = in.Policy().AsValidatingAdmissionPolicy().GetDefinition()
	} else if in.Policy().AsImageValidatingPolicy() != nil {
		out.ImageValidatingPolicy = in.Policy().AsImageValidatingPolicy()
	} else if in.Policy().AsDeletingPolicy() != nil {
		out.DeletingPolicy = in.Policy().AsDeletingPolicy()
	}

	for _, ruleresponse := range in.PolicyResponse.Rules {
		out.PolicyResponse.Rules = append(out.PolicyResponse.Rules, convertRuleResponse(in.Policy().GetName(), ruleresponse))
	}
	return out
}
