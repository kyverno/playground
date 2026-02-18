package models

import (
	"github.com/kyverno/api/api/policies.kyverno.io/v1beta1"
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
	if p := in.Policy().AsKyvernoPolicy(); p != nil {
		out.Policy = Policy{
			APIVersion:  "kyverno.io/v1",
			Kind:        p.GetKind(),
			Name:        p.GetName(),
			Namespace:   p.GetNamespace(),
			Labels:      p.GetLabels(),
			Annotations: p.GetAnnotations(),
		}
	} else if p := in.Policy().AsValidatingPolicy(); p != nil {
		out.Policy = Policy{
			APIVersion:  p.APIVersion,
			Kind:        p.GetKind(),
			Name:        p.GetName(),
			Namespace:   p.GetNamespace(),
			Labels:      p.GetLabels(),
			Annotations: p.GetAnnotations(),
		}
	} else if p := in.Policy().AsNamespacedValidatingPolicy(); p != nil {
		out.Policy = Policy{
			APIVersion:  p.APIVersion,
			Kind:        p.GetKind(),
			Name:        p.GetName(),
			Namespace:   p.GetNamespace(),
			Labels:      p.GetLabels(),
			Annotations: p.GetAnnotations(),
		}
	} else if p := in.Policy().AsValidatingAdmissionPolicy(); p != nil {
		out.Policy = Policy{
			APIVersion:  p.GetDefinition().APIVersion,
			Kind:        p.GetDefinition().Kind,
			Name:        p.GetDefinition().Name,
			Namespace:   p.GetDefinition().Namespace,
			Labels:      p.GetDefinition().Labels,
			Annotations: p.GetDefinition().Annotations,
		}
	} else if p := in.Policy().AsImageValidatingPolicy(); p != nil {
		out.Policy = Policy{
			APIVersion:  p.APIVersion,
			Kind:        p.GetKind(),
			Name:        p.GetName(),
			Namespace:   p.GetNamespace(),
			Labels:      p.GetLabels(),
			Annotations: p.GetAnnotations(),
		}
	} else if p := in.Policy().AsNamespacedImageValidatingPolicy(); p != nil {
		out.Policy = Policy{
			APIVersion:  p.APIVersion,
			Kind:        p.GetKind(),
			Name:        p.GetName(),
			Namespace:   p.GetNamespace(),
			Labels:      p.GetLabels(),
			Annotations: p.GetAnnotations(),
		}
	} else if p := in.Policy().AsDeletingPolicy(); p != nil {
		out.Policy = Policy{
			APIVersion:  "policies.kyverno.io/v1",
			Kind:        p.GetKind(),
			Name:        p.GetName(),
			Namespace:   p.GetNamespace(),
			Labels:      p.GetLabels(),
			Annotations: p.GetAnnotations(),
		}
	} else if p, ok := in.Policy().AsObject().(*v1beta1.NamespacedDeletingPolicy); ok {
		out.Policy = Policy{
			APIVersion:  "policies.kyverno.io/v1",
			Kind:        p.GetKind(),
			Name:        p.GetName(),
			Namespace:   p.GetNamespace(),
			Labels:      p.GetLabels(),
			Annotations: p.GetAnnotations(),
		}
	} else if p := in.Policy().AsMutatingPolicy(); p != nil {
		out.Policy = Policy{
			APIVersion:  p.APIVersion,
			Kind:        p.GetKind(),
			Name:        p.GetName(),
			Namespace:   p.GetNamespace(),
			Labels:      p.GetLabels(),
			Annotations: p.GetAnnotations(),
		}
	} else if p := in.Policy().AsGeneratingPolicy(); p != nil {
		out.Policy = Policy{
			APIVersion:  p.APIVersion,
			Kind:        p.GetKind(),
			Name:        p.GetName(),
			Namespace:   p.GetNamespace(),
			Labels:      p.GetLabels(),
			Annotations: p.GetAnnotations(),
		}
	} else if p := in.Policy().AsNamespacedMutatingPolicy(); p != nil {
		out.Policy = Policy{
			APIVersion:  p.APIVersion,
			Kind:        p.GetKind(),
			Name:        p.GetName(),
			Namespace:   p.GetNamespace(),
			Labels:      p.GetLabels(),
			Annotations: p.GetAnnotations(),
		}
	} else if p := in.Policy().AsNamespacedGeneratingPolicy(); p != nil {
		out.Policy = Policy{
			APIVersion:  p.APIVersion,
			Kind:        p.GetKind(),
			Name:        p.GetName(),
			Namespace:   p.GetNamespace(),
			Labels:      p.GetLabels(),
			Annotations: p.GetAnnotations(),
		}
	}

	for _, ruleresponse := range in.PolicyResponse.Rules {
		out.PolicyResponse.Rules = append(out.PolicyResponse.Rules, convertRuleResponse(in.Policy().GetName(), ruleresponse))
	}
	return out
}
