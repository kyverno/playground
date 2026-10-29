package mcp

import (
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/kyverno/playground/backend/pkg/mcp/tools"
)

func New() *server.StreamableHTTPServer {
	s := server.NewMCPServer(
		"Kyverno Playground",
		"0.1.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	s.AddTool(
		mcp.NewTool("execute",
			mcp.WithDescription("Execute a set of Kyverno policies in YAML format against a set of kubernetes resources in YAML format or JSON payloads"),
			mcp.WithInputSchema[tools.ExecuteInputSchema](),
			mcp.WithOutputSchema[tools.ExecuteOutputSchema](),
		),
		mcp.NewStructuredToolHandler(tools.HandleExecute),
	)

	s.AddTool(
		mcp.NewTool("validate",
			mcp.WithDescription("validate a set of kyverno policies in YAML format against the CRD schema and API definition. This tool does not validate the policies against any resources. It is used to validate the policies themselves."),
			mcp.WithString("policies", mcp.Title("Kyverno Policies"), mcp.Description("a set of kyverno policies and policyexceptions in YAML format"), mcp.Required()),
		),
		tools.HandleValidate,
	)

	return server.NewStreamableHTTPServer(s, server.WithEndpointPath(""))
}
