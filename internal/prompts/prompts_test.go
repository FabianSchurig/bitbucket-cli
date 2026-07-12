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

	expectedPrompts := map[string]bool{
		"bitbucket_pr_reviewer":      false,
		"bitbucket_comments_griller": false,
	}

	for prompt, err := range clientSession.Prompts(ctx, nil) {
		if err != nil {
			t.Fatalf("listing prompts: %v", err)
		}
		if _, ok := expectedPrompts[prompt.Name]; ok {
			expectedPrompts[prompt.Name] = true
			if prompt.Description == "" {
				t.Errorf("expected non-empty description for prompt %s", prompt.Name)
			}
			// Verify arguments are listed
			if len(prompt.Arguments) != 3 {
				t.Errorf("expected 3 arguments for prompt %s, got %d", prompt.Name, len(prompt.Arguments))
			}
		}
	}

	for name, found := range expectedPrompts {
		if !found {
			t.Errorf("expected prompt %s to be registered", name)
		}
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

	// Test bitbucket_pr_reviewer with templated args
	result, err := clientSession.GetPrompt(ctx, &mcp.GetPromptParams{
		Name: "bitbucket_pr_reviewer",
		Arguments: map[string]string{
			"pull_request_id": "456",
			"workspace":       "my-workspace",
			"repo_slug":       "my-repo",
		},
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
	if !strings.Contains(tc.Text, "456") {
		t.Error("expected templated Pull Request ID to be present in prompt output")
	}
	if !strings.Contains(tc.Text, "my-workspace") {
		t.Error("expected templated Workspace to be present in prompt output")
	}
	if !strings.Contains(tc.Text, "my-repo") {
		t.Error("expected templated Repository Slug to be present in prompt output")
	}

	// Test bitbucket_comments_griller
	resultGriller, err := clientSession.GetPrompt(ctx, &mcp.GetPromptParams{
		Name: "bitbucket_comments_griller",
		Arguments: map[string]string{
			"pull_request_id": "789",
		},
	})
	if err != nil {
		t.Fatalf("get griller prompt: %v", err)
	}
	tcGriller, ok := resultGriller.Messages[0].Content.(*mcp.TextContent)
	if !ok {
		t.Fatal("expected TextContent for griller message")
	}
	if !strings.Contains(tcGriller.Text, "Comments Griller Playbook") {
		t.Error("expected markdown content to contain griller playbook title")
	}
	if !strings.Contains(tcGriller.Text, "789") {
		t.Error("expected templated Pull Request ID to be present in griller output")
	}
}
