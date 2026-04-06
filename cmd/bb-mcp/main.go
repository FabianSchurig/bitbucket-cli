// bb-mcp is a Model Context Protocol (MCP) server for Bitbucket Cloud.
//
// It exposes all Bitbucket API operations as MCP tools, grouped by resource
// type (pull requests, repositories, pipelines, etc.). Each tool uses a
// CRUD-combined design with an "operation" parameter.
//
// An optional mcp_config.yaml file in the working directory controls which
// tools and HTTP methods are exposed at runtime.
//
// Install:
//
//	go install github.com/FabianSchurig/bitbucket-cli/cmd/bb-mcp@latest
//
// Authentication:
//
//	API token (recommended):
//	  export BITBUCKET_USERNAME=myuser
//	  export BITBUCKET_TOKEN=<token>
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

	"github.com/FabianSchurig/bitbucket-cli/internal/config"
	"github.com/FabianSchurig/bitbucket-cli/internal/mcptools"
	"github.com/FabianSchurig/bitbucket-cli/internal/prompts"
)

// Set via ldflags at build time (see goreleaser.yaml).
var (
	version = "dev"
)

func main() {
	transport := flag.String("transport", "stdio", "MCP transport: stdio or sse")
	addr := flag.String("addr", ":8080", "Address to listen on (SSE transport only)")
	configFile := flag.String("config", config.DefaultConfigFile, "Path to MCP configuration file")
	outputFmt := flag.String("output", "markdown", "Response format for tool results: markdown, table, json, id")
	flag.Parse()

	mcptools.Format = *outputFmt

	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatalf("loading config: %v", err)
	}

	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "bitbucket-mcp",
			Version: version,
		},
		nil,
	)

	// Register filtered Bitbucket tool groups.
	registerAllTools(server, cfg)

	// Register MCP prompts (playbooks).
	prompts.Register(server)

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

func registerAllTools(server *mcp.Server, cfg *config.Config) {
	for _, group := range mcptools.AllToolGroups {
		// Filter 1: Skip entirely if the tool name is ignored.
		if cfg.IsToolIgnored(group.Name) {
			continue
		}

		// Filter 2: Remove operations whose HTTP method is not allowed.
		filtered := mcptools.FilterOperations(group, cfg.IsMethodAllowed)
		if len(filtered.Operations) == 0 {
			continue
		}

		// Override: Replace description if configured.
		if override, ok := cfg.ToolOverrides[filtered.Name]; ok && override.Description != "" {
			filtered.Description = override.Description
		}

		mcptools.RegisterToolGroup(server, filtered)
	}
}
