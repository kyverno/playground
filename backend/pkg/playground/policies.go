package playground

import (
	"fmt"
	"strings"

	yamlutils "github.com/kyverno/kyverno/ext/yaml"
	"sigs.k8s.io/yaml"

	"github.com/kyverno/playground/backend/pkg/engine/models"
	"github.com/kyverno/playground/backend/pkg/serve"
)

// policyDoc is the metadata the Playground shows for a policy. Policies are only inspected here
// to pick the evaluation mode and label results; evaluation happens in kyverno serve.
type policyDoc struct {
	apiVersion  string
	kind        string
	name        string
	namespace   string
	mode        string
	labels      map[string]string
	annotations map[string]string
}

type policyIndex []policyDoc

// evaluable lists the kinds kyverno serve evaluates (bindings and exceptions ride along).
var evaluable = map[string]bool{
	"ValidatingPolicy": true, "NamespacedValidatingPolicy": true,
	"MutatingPolicy": true, "NamespacedMutatingPolicy": true,
	"GeneratingPolicy": true, "NamespacedGeneratingPolicy": true,
	"ImageValidatingPolicy": true, "NamespacedImageValidatingPolicy": true,
	"DeletingPolicy": true, "NamespacedDeletingPolicy": true,
	"ValidatingAdmissionPolicy": true, "MutatingAdmissionPolicy": true,
	// legacy kinds are sent on so kyverno serve rejects them with its migration hint
	"ClusterPolicy": true, "Policy": true,
}

func parsePolicies(content string) (policyIndex, error) {
	docs, err := yamlutils.SplitDocuments([]byte(content))
	if err != nil {
		return nil, fmt.Errorf("unable to parse policies: %w", err)
	}
	var out policyIndex
	for _, doc := range docs {
		var obj struct {
			APIVersion string `json:"apiVersion"`
			Kind       string `json:"kind"`
			Metadata   struct {
				Name        string            `json:"name"`
				Namespace   string            `json:"namespace"`
				Labels      map[string]string `json:"labels"`
				Annotations map[string]string `json:"annotations"`
			} `json:"metadata"`
			Spec struct {
				Evaluation struct {
					Mode string `json:"mode"`
				} `json:"evaluation"`
			} `json:"spec"`
		}
		if err := yaml.Unmarshal(doc, &obj); err != nil {
			return nil, fmt.Errorf("unable to parse policies: %w", err)
		}
		if !evaluable[obj.Kind] {
			continue
		}
		mode := modeKubernetes
		switch strings.ToLower(obj.Spec.Evaluation.Mode) {
		case "json":
			mode = modeJSON
		case "http":
			mode = modeHTTP
		case "envoy":
			mode = modeEnvoy
		}
		out = append(out, policyDoc{
			apiVersion:  obj.APIVersion,
			kind:        obj.Kind,
			name:        obj.Metadata.Name,
			namespace:   obj.Metadata.Namespace,
			mode:        mode,
			labels:      obj.Metadata.Labels,
			annotations: obj.Metadata.Annotations,
		})
	}
	return out, nil
}

// mode picks the payload mode, with the same precedence the Playground has always used:
// Kubernetes, then JSON, then Envoy, then HTTP. Empty means there is nothing to evaluate.
func (p policyIndex) mode() string {
	for _, m := range []string{modeKubernetes, modeJSON, modeEnvoy, modeHTTP} {
		for _, d := range p {
			if d.mode == m {
				return m
			}
		}
	}
	return ""
}

func (p policyIndex) describe(ref serve.ObjectRef) models.Policy {
	out := models.Policy{APIVersion: ref.APIVersion, Kind: ref.Kind, Name: ref.Name, Namespace: ref.Namespace}
	for _, d := range p {
		if d.name == ref.Name && d.namespace == ref.Namespace && (d.kind == ref.Kind || ref.Kind == "") {
			if out.APIVersion == "" {
				out.APIVersion = d.apiVersion
			}
			if out.Kind == "" {
				out.Kind = d.kind
			}
			out.Labels = d.labels
			out.Annotations = d.annotations
			if d.kind != "ValidatingAdmissionPolicy" && d.kind != "MutatingAdmissionPolicy" {
				out.Mode = d.mode
			}
			break
		}
	}
	return out
}
