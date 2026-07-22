package tools

import (
	"context"

	"github.com/kyverno/kyverno/ext/resource/loader"
	"github.com/kyverno/kyverno/pkg/cel/policies/dpol"
	"github.com/kyverno/kyverno/pkg/cel/policies/gpol"
	"github.com/kyverno/kyverno/pkg/cel/policies/mpol"
	"github.com/kyverno/kyverno/pkg/cel/policies/vpol"
	"github.com/mark3labs/mcp-go/mcp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/playground"
	"github.com/kyverno/playground/backend/pkg/policy"
)

type Policy struct {
	metav1.TypeMeta `json:",inline"`
	Metadata        metav1.ObjectMeta `json:"metadata"`
}

type Validation struct {
	Policy   Policy   `json:"policy"`
	RuleName string   `json:"ruleName"`
	Message  string   `json:"message"`
	Warnings []string `json:"warnings,omitempty"`
	Valid    bool     `json:"valid"`
}

func HandleValidate(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	request := &playground.EngineRequest{}
	if err := req.BindArguments(request); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	openAPI, err := cluster.NewFake().OpenAPIClient("1.34.0")
	if err != nil {
		return nil, err
	}

	policyLoader, err := loader.New(openAPI)
	if err != nil {
		return nil, err
	}

	k8s, _, _, err := policy.Load(policyLoader, []byte(request.Policies))
	if err != nil {
		return nil, err
	}

	response := make([]Validation, 0, k8s.Length())
	for _, pol := range k8s.ValidatingPolicies {
		w, err := vpol.Validate(pol)
		if err != nil {
			response = append(response, Validation{
				Policy: Policy{
					TypeMeta: metav1.TypeMeta{
						Kind:       pol.GetObjectKind().GroupVersionKind().Kind,
						APIVersion: pol.GetObjectKind().GroupVersionKind().Version,
					},
					Metadata: metav1.ObjectMeta{
						Name:      pol.GetName(),
						Namespace: pol.GetNamespace(),
					},
				},
				Warnings: w,
				Valid:    false,
			})
		}
	}
	for _, pol := range k8s.MutatingPolicies {
		w, err := mpol.Validate(pol)
		if err != nil {
			response = append(response, Validation{
				Policy: Policy{
					TypeMeta: metav1.TypeMeta{
						Kind:       pol.GetObjectKind().GroupVersionKind().Kind,
						APIVersion: pol.GetObjectKind().GroupVersionKind().Version,
					},
					Metadata: metav1.ObjectMeta{
						Name:      pol.GetName(),
						Namespace: pol.GetNamespace(),
					},
				},
				Warnings: w,
				Valid:    false,
			})
		}
	}
	for _, pol := range k8s.GeneratingPolicies {
		w, err := gpol.Validate(pol)
		if err != nil {
			response = append(response, Validation{
				Policy: Policy{
					TypeMeta: metav1.TypeMeta{
						Kind:       pol.GetObjectKind().GroupVersionKind().Kind,
						APIVersion: pol.GetObjectKind().GroupVersionKind().Version,
					},
					Metadata: metav1.ObjectMeta{
						Name:      pol.GetName(),
						Namespace: pol.GetNamespace(),
					},
				},
				Warnings: w,
				Valid:    false,
			})
		}
	}
	for _, pol := range k8s.DeletingPolicies {
		w, err := dpol.Validate(pol)
		if err != nil {
			response = append(response, Validation{
				Policy: Policy{
					TypeMeta: metav1.TypeMeta{
						Kind:       pol.GetObjectKind().GroupVersionKind().Kind,
						APIVersion: pol.GetObjectKind().GroupVersionKind().Version,
					},
					Metadata: metav1.ObjectMeta{
						Name:      pol.GetName(),
						Namespace: pol.GetNamespace(),
					},
				},
				Warnings: w,
				Valid:    false,
			})
		}
	}

	return mcp.NewToolResultJSON(response)
}
