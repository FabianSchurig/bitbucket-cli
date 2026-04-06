package prompts

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestRegister_PromptsListedByClient(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "v0.0.1"}, nil)
	Register(server)

	ctx := context.Background()
	ct, st := mcp.NewInMemoryTransports()

	serverSession, err := server.Connect(ctx, st, nil)
	if err != nil {
		t.Fatalf("server connect: %v", err)
	}
	defer func() { _ = serverSession.Close() }()

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "v0.0.1"}, nil)
	clientSession, err := client.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer func() { _ = clientSession.Close() }()

	var found bool
	for prompt, err := range clientSession.Prompts(ctx, nil) {
		if err != nil {
			t.Fatalf("listing prompts: %v", err)
		}
		if prompt.Name == "bitbucket_pr_reviewer" {
			found = true
			if prompt.Description == "" {
				t.Error("expected non-empty description")
			}
		}
	}
	if !found {
		t.Error("expected bitbucket_pr_reviewer prompt to be registered")
	}
}

func TestRegister_GetPromptReturnsContent(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "v0.0.1"}, nil)
	Register(server)

	ctx := context.Background()
	ct, st := mcp.NewInMemoryTransports()

	serverSession, err := server.Connect(ctx, st, nil)
	if err != nil {
		t.Fatalf("server connect: %v", err)
	}
	defer func() { _ = serverSession.Close() }()

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "v0.0.1"}, nil)
	clientSession, err := client.Connect(ctx, ct, nil)
	if err != nil {
		t.Fatalf("client connect: %v", err)
	}
	defer func() { _ = clientSession.Close() }()

	result, err := clientSession.GetPrompt(ctx, &mcp.GetPromptParams{
		Name: "bitbucket_pr_reviewer",
	})
	if err != nil {
		t.Fatalf("get prompt: %v", err)
	}
	if len(result.Messages) == 0 {
		t.Fatal("expected at least one message")
	}
	msg := result.Messages[0]
	if msg.Role != "user" {
		t.Errorf("expected role 'user', got %q", msg.Role)
	}
	tc, ok := msg.Content.(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", msg.Content)
	}
	if !strings.Contains(tc.Text, "PR Reviewer Playbook") {
		t.Error("expected markdown content to contain playbook title")
	}
	if !strings.Contains(tc.Text, "listPullRequests") {
		t.Error("expected markdown content to mention listPullRequests operation")
	}
}
