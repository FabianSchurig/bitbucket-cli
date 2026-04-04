package main

import (
	"context"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestRegisterAllTools(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "v0.0.1"}, nil)
	registerAllTools(server)

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
