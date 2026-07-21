package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/kyverno/playground/backend/pkg/cluster"
	"github.com/kyverno/playground/backend/pkg/crd"
	"github.com/kyverno/playground/backend/pkg/playground"
)

func HandleExecute(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	request := &playground.EngineRequest{}
	if err := req.BindArguments(request); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	response, err := playground.Run(ctx, cluster.NewFake(), request, crd.APIConfiguration{
		BuiltInCrds: []string{"argocd", "cert-manager", "prometheus-operator", "tekton-pipeline", "wgpolicyk8s"},
		LocalCrds:   nil,
	})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultStructuredOnly(MapResponse(response)), nil
}
