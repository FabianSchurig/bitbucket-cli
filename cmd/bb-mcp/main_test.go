package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/FabianSchurig/bitbucket-cli/internal/config"
)

func TestRegisterAllTools(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "v0.0.1"}, nil)
	cfg := config.DefaultConfig()
	registerAllTools(server, cfg)

	ctx := context.Background()
	ct, st := mcp.NewInMemoryTransports()

	serverSession, err := server.Connect(ctx, st, nil)
	if err != nil {
		t.Fatalf("server connect: %v", err)
	}
	defer func() { _ = serverSession.Close() }()

	client := mcp.NewClient(&mcp.Implementation{Name: "client", Version: "v0.0.1"}, nil)
	clientSession, err := client.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer func() { _ = clientSession.Close() }()

	count := 0
	for _, err := range clientSession.Tools(ctx, nil) {
		if err != nil {
			t.Fatalf("tool listing failed: %v", err)
		}
		count++
	}

	if count != 20 {
		t.Fatalf("expected 20 registered MCP tools, got %d", count)
	}
}

func TestRegisterAllTools_FiltersDELETE(t *testing.T) {
	cfg, err := config.Parse([]byte(`
server:
  allowed_methods: ["GET", "POST", "PUT", "PATCH"]
`))
	if err != nil {
		t.Fatal(err)
	}

	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "v0.0.1"}, nil)
	registerAllTools(server, cfg)

	ctx := context.Background()
	ct, st := mcp.NewInMemoryTransports()

	serverSession, err := server.Connect(ctx, st, nil)
	if err != nil {
		t.Fatalf("server connect: %v", err)
	}
	defer func() { _ = serverSession.Close() }()

	client := mcp.NewClient(&mcp.Implementation{Name: "client", Version: "v0.0.1"}, nil)
	clientSession, err := client.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer func() { _ = clientSession.Close() }()

	for tool, err := range clientSession.Tools(ctx, nil) {
		if err != nil {
			t.Fatalf("tool listing failed: %v", err)
		}
		// Check that no tool description mentions DELETE.
		desc := tool.Description
		if strings.Contains(desc, "[DELETE]") {
			t.Errorf("tool %q still contains DELETE operations after filtering", tool.Name)
		}
	}
}

func TestRegisterAllTools_IgnoredTool(t *testing.T) {
	cfg, err := config.Parse([]byte(`
server:
  ignored_tools:
    - bitbucket_pr
`))
	if err != nil {
		t.Fatal(err)
	}

	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "v0.0.1"}, nil)
	registerAllTools(server, cfg)

	ctx := context.Background()
	ct, st := mcp.NewInMemoryTransports()

	serverSession, err := server.Connect(ctx, st, nil)
	if err != nil {
		t.Fatalf("server connect: %v", err)
	}
	defer func() { _ = serverSession.Close() }()

	client := mcp.NewClient(&mcp.Implementation{Name: "client", Version: "v0.0.1"}, nil)
	clientSession, err := client.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer func() { _ = clientSession.Close() }()

	for tool, err := range clientSession.Tools(ctx, nil) {
		if err != nil {
			t.Fatalf("tool listing failed: %v", err)
		}
		if tool.Name == "bitbucket_pr" {
			t.Error("bitbucket_pr should be ignored but was registered")
		}
	}
}

func TestRegisterAllTools_DescriptionOverride(t *testing.T) {
	cfg, err := config.Parse([]byte(`
tool_overrides:
  bitbucket_pr:
    description: "Custom PR description"
`))
	if err != nil {
		t.Fatal(err)
	}

	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "v0.0.1"}, nil)
	registerAllTools(server, cfg)

	ctx := context.Background()
	ct, st := mcp.NewInMemoryTransports()

	serverSession, err := server.Connect(ctx, st, nil)
	if err != nil {
		t.Fatalf("server connect: %v", err)
	}
	defer func() { _ = serverSession.Close() }()

	client := mcp.NewClient(&mcp.Implementation{Name: "client", Version: "v0.0.1"}, nil)
	clientSession, err := client.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer func() { _ = clientSession.Close() }()

	for tool, err := range clientSession.Tools(ctx, nil) {
		if err != nil {
			t.Fatalf("tool listing failed: %v", err)
		}
		if tool.Name == "bitbucket_pr" {
			if tool.Description != "Custom PR description" {
				t.Errorf("expected overridden description, got %q", tool.Description)
			}
			return
		}
	}
	t.Error("bitbucket_pr tool not found")
}

func TestRegisterAllTools_ConfigFromFile(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "mcp_config.yaml")
	if err := os.WriteFile(cfgPath, []byte(`
server:
  ignored_tools:
    - bitbucket_search
`), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load(cfgPath)
	if err != nil {
		t.Fatal(err)
	}

	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "v0.0.1"}, nil)
	registerAllTools(server, cfg)

	ctx := context.Background()
	ct, st := mcp.NewInMemoryTransports()

	serverSession, err := server.Connect(ctx, st, nil)
	if err != nil {
		t.Fatalf("server connect: %v", err)
	}
	defer func() { _ = serverSession.Close() }()

	client := mcp.NewClient(&mcp.Implementation{Name: "client", Version: "v0.0.1"}, nil)
	clientSession, err := client.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer func() { _ = clientSession.Close() }()

	count := 0
	for tool, err := range clientSession.Tools(ctx, nil) {
		if err != nil {
			t.Fatalf("tool listing failed: %v", err)
		}
		if tool.Name == "bitbucket_search" {
			t.Error("bitbucket_search should have been filtered out")
		}
		count++
	}
	if count != 19 {
		t.Errorf("expected 19 tools (20 minus 1 ignored), got %d", count)
	}
}
