// bb-mcp is a Model Context Protocol (MCP) server for Bitbucket Cloud.
//
// It exposes all Bitbucket API operations as MCP tools, grouped by resource
// type (pull requests, repositories, pipelines, etc.). Each tool uses a
// CRUD-combined design with an "operation" parameter.
//
// Install:
//
//	go install github.com/FabianSchurig/bitbucket-cli/cmd/bb-mcp@latest
//
// Authentication:
//
//	App password (most common):
//	  export BITBUCKET_USERNAME=myuser
//	  export BITBUCKET_APP_PASSWORD=ATBBxxxxxxxx
//
//	OAuth2 access token:
//	  export BITBUCKET_TOKEN=<token>
//
// Run:
//
//	bb-mcp                  # stdio transport (default, for MCP clients)
//	bb-mcp --transport sse  # SSE transport (for HTTP-based MCP clients)
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/FabianSchurig/bitbucket-cli/internal/mcptools"
)

// Set via ldflags at build time (see goreleaser.yaml).
var (
	version = "dev"
)

func main() {
	transport := flag.String("transport", "stdio", "MCP transport: stdio or sse")
	addr := flag.String("addr", ":8080", "Address to listen on (SSE transport only)")
	flag.Parse()

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "bitbucket-mcp",
			Version: version,
		},
		nil,
	)

	// Register all Bitbucket tool groups.
	registerAllTools(server)

	ctx := context.Background()
	switch *transport {
	case "stdio":
		if err := server.Run(ctx, &mcp.StdioTransport{}); err != nil {
			log.Fatalf("MCP server error: %v", err)
		}
	case "sse":
		handler := mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server {
			return server
		}, nil)
		fmt.Fprintf(os.Stderr, "MCP SSE server listening on %s\n", *addr)
		if err := http.ListenAndServe(*addr, handler); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	default:
		log.Fatalf("unknown transport: %s (valid: stdio, sse)", *transport)
	}
}

func registerAllTools(server *mcp.Server) {
	mcptools.RegisterToolGroup(server, mcptools.PRToolGroup)
	mcptools.RegisterToolGroup(server, mcptools.HooksToolGroup)
	mcptools.RegisterToolGroup(server, mcptools.SearchToolGroup)
	mcptools.RegisterToolGroup(server, mcptools.RefsToolGroup)
	mcptools.RegisterToolGroup(server, mcptools.CommitsToolGroup)
	mcptools.RegisterToolGroup(server, mcptools.ReportsToolGroup)
	mcptools.RegisterToolGroup(server, mcptools.ReposToolGroup)
	mcptools.RegisterToolGroup(server, mcptools.WorkspacesToolGroup)
	mcptools.RegisterToolGroup(server, mcptools.ProjectsToolGroup)
	mcptools.RegisterToolGroup(server, mcptools.PipelinesToolGroup)
	mcptools.RegisterToolGroup(server, mcptools.IssuesToolGroup)
	mcptools.RegisterToolGroup(server, mcptools.SnippetsToolGroup)
	mcptools.RegisterToolGroup(server, mcptools.DeploymentsToolGroup)
	mcptools.RegisterToolGroup(server, mcptools.BranchRestrictionsToolGroup)
	mcptools.RegisterToolGroup(server, mcptools.BranchingModelToolGroup)
	mcptools.RegisterToolGroup(server, mcptools.CommitStatusesToolGroup)
	mcptools.RegisterToolGroup(server, mcptools.DownloadsToolGroup)
	mcptools.RegisterToolGroup(server, mcptools.UsersToolGroup)
	mcptools.RegisterToolGroup(server, mcptools.PropertiesToolGroup)
	mcptools.RegisterToolGroup(server, mcptools.AddonToolGroup)
}
