package tools

import (
	"context"

	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/crd"
	"github.com/kyverno/playground/backend/pkg/engine/models"
	"github.com/kyverno/playground/backend/pkg/playground"
	"github.com/kyverno/playground/backend/pkg/utils"
)

type Context struct {
	Username        string                       `json:"username,omitempty" jsonschema:"Used Username for the admission request"`
	Groups          []string                     `json:"groups,omitempty" jsonschema:"Used Groups for the admission request"`
	Roles           []string                     `json:"roles,omitempty" jsonschema:"Used Roles for the admission request, only used for ClusterPolicy and Policy types"`
	ClusterRoles    []string                     `json:"clusterRoles,omitempty" jsonschema:"Used ClusterRoles for the admission request, only used for ClusterPolicy and Policy types"`
	NamespaceLabels map[string]string            `json:"namespaceLabels,omitempty" jsonschema:"Used NamespaceLabels for the admission request"`
	DryRun          bool                         `json:"dryRun,omitempty" jsonschema:"Used DryRun for the admission request"`
	Operation       kyvernov1.AdmissionOperation `json:"operation,omitempty" jsonschema:"Used Operation for the admission request"`
}

type Parameters struct {
	Kubernetes string         `json:"kubernetes,omitempty" jsonschema:"the kubernetes version used for resource schema validation."`
	Variables  map[string]any `json:"variables,omitempty" jsonschema:"to fake variable substitution, you can provide a map of variable names to their values. The variables will be substituted in the policies during execution."`
	Context    *Context       `json:"context,omitempty" jsonschema:"The context for the admission request."`
}

type ExecuteInputSchema struct {
	Policies                  string      `json:"policies" jsonschema:"a set of kyverno policies and policyexceptions in YAML format to run against."`
	Resources                 string      `json:"resources" jsonschema:"a set of Kubernetes resources in YAML format or JSON payloads to run the policies against."`
	OldResources              string      `json:"oldResources,omitempty" jsonschema:"a set of the of Kubernetes resources in YAML format or JSON payloads which represents the previous state of 'resources' in case of an UPDATE request."`
	ClusterResources          string      `json:"clusterResources,omitempty" jsonschema:"a set of Kubernetes resources in YAML format which represents the existing resources in the cluster to use for resource lookups during policy execution."`
	CustomResourceDefinitions string      `json:"customResourceDefinitions,omitempty" jsonschema:"a set of custom resource definitions in YAML format to validate the custom resource schemas."`
	Parameters                *Parameters `json:"parameters,omitempty" jsonschema:"The parameters for the admission request."`
}

type ExecuteOutputSchema struct {
	Results []Result `json:"results" jsonschema:"The results of the policy execution."`
}

func HandleExecute(ctx context.Context, req mcp.CallToolRequest, args ExecuteInputSchema) (*mcp.CallToolResult, error) {
	if args.Parameters == nil {
		args.Parameters = &Parameters{}
	}
	if args.Parameters.Context == nil {
		args.Parameters.Context = &Context{}
	}
	op := utils.Fallback(args.Parameters.Context.Operation, kyvernov1.Create)

	request := &playground.EngineRequest{
		Policies:                  args.Policies,
		Resources:                 args.Resources,
		OldResources:              args.OldResources,
		ClusterResources:          args.ClusterResources,
		CustomResourceDefinitions: args.CustomResourceDefinitions,
		Operation:                 string(op),
		Parameters: &models.Parameters{
			Kubernetes: models.Kubernetes{
				Version: utils.Fallback(args.Parameters.Kubernetes, "v1.34.0"),
			},
			Variables: args.Parameters.Variables,
			Context: models.Context{
				Username:        args.Parameters.Context.Username,
				Groups:          args.Parameters.Context.Groups,
				Roles:           args.Parameters.Context.Roles,
				ClusterRoles:    args.Parameters.Context.ClusterRoles,
				NamespaceLabels: args.Parameters.Context.NamespaceLabels,
				DryRun:          args.Parameters.Context.DryRun,
				Operation:       op,
			},
		},
	}

	response, err := playground.Run(ctx, cluster.NewFake(), request, crd.APIConfiguration{
		BuiltInCrds: []string{"argocd", "cert-manager", "prometheus-operator", "tekton-pipeline", "wgpolicyk8s"},
		LocalCrds:   nil,
	})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(MapResponse(response))
}
