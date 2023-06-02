package engine

import (
	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	kyvernov2alpha1 "github.com/kyverno/kyverno/api/kyverno/v2alpha1"
	engineapi "github.com/kyverno/kyverno/pkg/engine/api"
	"github.com/mattbaird/jsonpatch"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

type Exceptions struct {
	Enabled   bool   `json:"enabled"`
	Namespace string `json:"namespace"`
}

type Cosign struct {
	ImageSignatureRepository string `json:"imageSignatureRepository"`
}

type Registry struct {
	AllowInsecure     bool     `json:"allowInsecure"`
	PullSecrets       []string `json:"pullSecrets"`
	CredentialHelpers []string `json:"credentialHelpers"`
}

type ProtectManagedResources struct {
	Enabled bool `json:"enabled"`
}

type ForceFailurePolicyIgnore struct {
	Enabled bool `json:"enabled"`
}

type Flags struct {
	Exceptions               Exceptions               `json:"exceptions"`
	Cosign                   Cosign                   `json:"cosign"`
	Registry                 Registry                 `json:"registry"`
	ProtectManagedResources  ProtectManagedResources  `json:"protectManagedResources"`
	ForceFailurePolicyIgnore ForceFailurePolicyIgnore `json:"forceFailurePolicyIgnore"`
}

type Parameters struct {
	Kubernetes Kubernetes             `json:"kubernetes"`
	Context    Context                `json:"context"`
	Variables  map[string]interface{} `json:"variables"`
	Flags      Flags                  `json:"flags"`
}

type Context struct {
	Username        string                       `json:"username"`
	Groups          []string                     `json:"groups"`
	Roles           []string                     `json:"roles"`
	ClusterRoles    []string                     `json:"clusterRoles"`
	Operation       kyvernov1.AdmissionOperation `json:"operation"`
	NamespaceLabels map[string]string            `json:"namespaceLabels"`
	DryRun          bool                         `json:"dryRun"`
}

type Kubernetes struct {
	Version string `json:"version"`
}

type Results struct {
	Mutation          []Response `json:"mutation"`
	ImageVerification []Response `json:"imageVerification"`
	Validation        []Response `json:"validation"`
	Generation        []Response `json:"generation"`
}

type Response struct {
	// OriginalResource is the original resource as YAML string
	OriginalResource string `json:"originalResource"`
	// Resource is the original resource
	Resource unstructured.Unstructured `json:"resource"`
	// Policy is the original policy
	Policy kyvernov1.PolicyInterface `json:"policy"`
	// namespaceLabels given by policy context
	NamespaceLabels map[string]string `json:"namespaceLabels"`
	// PatchedResource is the resource patched with the engine action changes
	PatchedResource string `json:"patchedResource"`
	// PolicyResponse contains the engine policy response
	PolicyResponse PolicyResponse `json:"policyResponse"`
}

type PolicyResponse struct {
	// Rules contains policy rules responses
	Rules []RuleResponse `json:"rules"`
}

type RuleResponse struct {
	// name is the rule name specified in policy
	Name string `json:"name"`
	// ruleType is the rule type (Mutation,Generation,Validation) for Kyverno Policy
	RuleType engineapi.RuleType `json:"ruleType"`
	// message is the message response from the rule application
	Message string `json:"message"`
	// status rule status
	Status engineapi.RuleStatus `json:"status"`
	// patches are JSON patches, for mutation rules
	Patches []jsonpatch.JsonPatchOperation `json:"patches"`
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
	Exception *kyvernov2alpha1.PolicyException `json:"exception"`
}

func ConvertRuleResponse(in engineapi.RuleResponse) RuleResponse {
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
	out.Patches = append(out.Patches, in.Patches()...)
	return out
}

func convertResponse(in engineapi.EngineResponse) Response {
	patchedResource, _ := yaml.Marshal(in.PatchedResource.Object)
	resource, _ := yaml.Marshal(in.Resource.Object)
	out := Response{
		OriginalResource: string(resource),
		Resource:         in.Resource,
		Policy:           in.Policy(),
		NamespaceLabels:  in.NamespaceLabels(),
		PatchedResource:  string(patchedResource),
	}
	for _, ruleresponse := range in.PolicyResponse.Rules {
		out.PolicyResponse.Rules = append(out.PolicyResponse.Rules, ConvertRuleResponse(ruleresponse))
	}
	return out
}
