package mcptools_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/FabianSchurig/bitbucket-cli/internal/mcptools"
)

// ─── Test helpers ─────────────────────────────────────────────────────────────

// setupEnv configures Bitbucket auth env vars and returns a cleanup function.
func setupEnv(t *testing.T) {
	t.Helper()
	t.Setenv("BITBUCKET_USERNAME", "testuser")
	t.Setenv("BITBUCKET_TOKEN", "testtoken")
}

// testToolGroup returns a minimal ToolGroup for testing.
func testToolGroup() mcptools.ToolGroup {
	return mcptools.ToolGroup{
		Name:        "test_tool",
		Description: "Test tool for unit tests",
		Operations: []mcptools.OperationDef{
			{
				OperationID: "listItems",
				Method:      "GET",
				Path:        "/repositories/{workspace}/{repo_slug}/items",
				Summary:     "List items",
				Description: "Returns all items",
				Params: []mcptools.ParamDef{
					{Name: "workspace", In: "path", Type: "string", Required: true},
					{Name: "repo_slug", In: "path", Type: "string", Required: true},
					{Name: "state", In: "query", Type: "string"},
				},
				Paginated: true,
			},
			{
				OperationID: "getItem",
				Method:      "GET",
				Path:        "/repositories/{workspace}/{repo_slug}/items/{item_id}",
				Summary:     "Get an item",
				Params: []mcptools.ParamDef{
					{Name: "workspace", In: "path", Type: "string", Required: true},
					{Name: "repo_slug", In: "path", Type: "string", Required: true},
					{Name: "item_id", In: "path", Type: "integer", Required: true},
				},
			},
			{
				OperationID: "createItem",
				Method:      "POST",
				Path:        "/repositories/{workspace}/{repo_slug}/items",
				Summary:     "Create an item",
				Params: []mcptools.ParamDef{
					{Name: "workspace", In: "path", Type: "string", Required: true},
					{Name: "repo_slug", In: "path", Type: "string", Required: true},
				},
				BodyFields: []mcptools.BodyFieldDef{
					{Path: "title", Type: "string", Desc: "Item title"},
					{Path: "content.raw", Type: "string", Desc: "Item content"},
				},
				HasBody: true,
			},
			{
				OperationID: "deleteItem",
				Method:      "DELETE",
				Path:        "/repositories/{workspace}/{repo_slug}/items/{item_id}",
				Summary:     "Delete an item",
				Params: []mcptools.ParamDef{
					{Name: "workspace", In: "path", Type: "string", Required: true},
					{Name: "repo_slug", In: "path", Type: "string", Required: true},
					{Name: "item_id", In: "path", Type: "integer", Required: true},
				},
			},
		},
	}
}

// callTool creates an MCP server, registers a tool group, connects a client,
// and calls the tool with the given arguments.
func callTool(t *testing.T, group mcptools.ToolGroup, args map[string]any) *mcp.CallToolResult {
	t.Helper()

	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "v0.0.1"}, nil)
	mcptools.RegisterToolGroup(server, group)

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

	res, err := clientSession.CallTool(ctx, &mcp.CallToolParams{
		Name:      group.Name,
		Arguments: args,
	})
	if err != nil {
		t.Fatalf("CallTool: %v", err)
	}
	return res
}

// textContent extracts the text from the first TextContent in a result.
func textContent(t *testing.T, res *mcp.CallToolResult) string {
	t.Helper()
	if len(res.Content) == 0 {
		t.Fatal("empty result content")
	}
	tc, ok := res.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected TextContent, got %T", res.Content[0])
	}
	return tc.Text
}

// ─── Tests ────────────────────────────────────────────────────────────────────

func TestRegisterToolGroup_ToolListedByClient(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "v0.0.1"}, nil)
	mcptools.RegisterToolGroup(server, testToolGroup())

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

	// List tools and verify our tool is registered.
	var tools []*mcp.Tool
	for tool, err := range clientSession.Tools(ctx, nil) {
		if err != nil {
			t.Fatalf("listing tools: %v", err)
		}
		tools = append(tools, tool)
	}

	if len(tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(tools))
	}
	if tools[0].Name != "test_tool" {
		t.Errorf("expected tool name 'test_tool', got %q", tools[0].Name)
	}
}

func TestToolHandler_MissingOperation(t *testing.T) {
	res := callTool(t, testToolGroup(), map[string]any{})
	if !res.IsError {
		t.Error("expected error for missing operation")
	}
	text := textContent(t, res)
	if text != "missing required parameter: operation" {
		t.Errorf("unexpected error message: %s", text)
	}
}

func TestToolHandler_UnknownOperation(t *testing.T) {
	res := callTool(t, testToolGroup(), map[string]any{
		"operation": "doesNotExist",
	})
	if !res.IsError {
		t.Error("expected error for unknown operation")
	}
	text := textContent(t, res)
	if text != "unknown operation: doesNotExist" {
		t.Errorf("unexpected error message: %s", text)
	}
}

func TestToolHandler_MissingRequiredParam(t *testing.T) {
	res := callTool(t, testToolGroup(), map[string]any{
		"operation": "listItems",
		"workspace": "myws",
		// missing repo_slug
	})
	if !res.IsError {
		t.Error("expected error for missing required param")
	}
	text := textContent(t, res)
	if text != "missing required parameter: repo_slug (for operation listItems)" {
		t.Errorf("unexpected error message: %s", text)
	}
}

func TestToolHandler_NoAuth(t *testing.T) {
	// Clear all auth env vars.
	for _, k := range []string{"BITBUCKET_USERNAME", "BITBUCKET_TOKEN"} {
		if err := os.Unsetenv(k); err != nil {
			t.Fatalf("unsetenv %s: %v", k, err)
		}
	}

	res := callTool(t, testToolGroup(), map[string]any{
		"operation": "listItems",
		"workspace": "myws",
		"repo_slug": "myrepo",
	})
	if !res.IsError {
		t.Error("expected error for missing auth")
	}
	text := textContent(t, res)
	if text == "" {
		t.Error("expected non-empty error message")
	}
}

func TestToolHandler_GET_Success(t *testing.T) {
	setupEnv(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		expected := "/2.0/repositories/myws/myrepo/items/42"
		if r.URL.Path != expected {
			t.Errorf("expected path %s, got %s", expected, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"id": 42, "title": "Test Item"})
	}))
	defer srv.Close()
	t.Setenv("BITBUCKET_BASE_URL", srv.URL+"/2.0")

	res := callTool(t, testToolGroup(), map[string]any{
		"operation": "getItem",
		"workspace": "myws",
		"repo_slug": "myrepo",
		"item_id":   42,
	})
	if res.IsError {
		t.Fatalf("unexpected error: %s", textContent(t, res))
	}

	text := textContent(t, res)
	var result map[string]any
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if result["title"] != "Test Item" {
		t.Errorf("expected title 'Test Item', got %v", result["title"])
	}
}

func TestToolHandler_GET_PaginatedList(t *testing.T) {
	setupEnv(t)

	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/2.0/repositories/myws/myrepo/items" && callCount == 1:
			nextURL := "http://" + r.Host + "/page2"
			_ = json.NewEncoder(w).Encode(map[string]any{
				"values": []any{map[string]any{"id": 1, "title": "Item 1"}},
				"next":   nextURL,
			})
		case r.URL.Path == "/page2":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"values": []any{map[string]any{"id": 2, "title": "Item 2"}},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()
	t.Setenv("BITBUCKET_BASE_URL", srv.URL+"/2.0")

	res := callTool(t, testToolGroup(), map[string]any{
		"operation": "listItems",
		"workspace": "myws",
		"repo_slug": "myrepo",
	})
	if res.IsError {
		t.Fatalf("unexpected error: %s", textContent(t, res))
	}

	text := textContent(t, res)
	var items []any
	if err := json.Unmarshal([]byte(text), &items); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
	if callCount != 2 {
		t.Errorf("expected 2 HTTP calls (pagination), got %d", callCount)
	}
}

func TestToolHandler_POST_WithBodyFields(t *testing.T) {
	setupEnv(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decoding body: %v", err)
		}
		if body["title"] != "New Item" {
			t.Errorf("expected title 'New Item', got %v", body["title"])
		}
		content, _ := body["content"].(map[string]any)
		if content["raw"] != "Hello world" {
			t.Errorf("expected content.raw 'Hello world', got %v", content["raw"])
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": 99, "title": "New Item"})
	}))
	defer srv.Close()
	t.Setenv("BITBUCKET_BASE_URL", srv.URL+"/2.0")

	res := callTool(t, testToolGroup(), map[string]any{
		"operation":        "createItem",
		"workspace":        "myws",
		"repo_slug":        "myrepo",
		"body_title":       "New Item",
		"body_content_raw": "Hello world",
	})
	if res.IsError {
		t.Fatalf("unexpected error: %s", textContent(t, res))
	}

	text := textContent(t, res)
	var result map[string]any
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}
	if result["title"] != "New Item" {
		t.Errorf("expected title 'New Item', got %v", result["title"])
	}
}

func TestToolHandler_DELETE_Success(t *testing.T) {
	setupEnv(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()
	t.Setenv("BITBUCKET_BASE_URL", srv.URL+"/2.0")

	res := callTool(t, testToolGroup(), map[string]any{
		"operation": "deleteItem",
		"workspace": "myws",
		"repo_slug": "myrepo",
		"item_id":   42,
	})
	if res.IsError {
		t.Fatalf("unexpected error: %s", textContent(t, res))
	}

	text := textContent(t, res)
	if text != "OK" {
		t.Errorf("expected 'OK', got %q", text)
	}
}

func TestToolHandler_POST_WithRawBody(t *testing.T) {
	setupEnv(t)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decoding body: %v", err)
		}
		if body["custom"] != "field" {
			t.Errorf("expected custom field, got %v", body["custom"])
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": 1})
	}))
	defer srv.Close()
	t.Setenv("BITBUCKET_BASE_URL", srv.URL+"/2.0")

	res := callTool(t, testToolGroup(), map[string]any{
		"operation": "createItem",
		"workspace": "myws",
		"repo_slug": "myrepo",
		"body":      `{"custom":"field"}`,
	})
	if res.IsError {
		t.Fatalf("unexpected error: %s", textContent(t, res))
	}
}

// ─── Integration tests using InMemoryTransport ───────────────────────────────
// These tests validate the full MCP protocol round-trip without hitting
// Bitbucket API. They use modified test infrastructure to intercept HTTP calls.

func TestToolHandler_MultipleOperationsPerTool(t *testing.T) {
	// Verify that a single tool exposes multiple operations via the operation enum.
	group := testToolGroup()

	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "v0.0.1"}, nil)
	mcptools.RegisterToolGroup(server, group)

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

	// List tools - should have exactly one tool with all operations.
	var tools []*mcp.Tool
	for tool, err := range clientSession.Tools(ctx, nil) {
		if err != nil {
			t.Fatalf("listing tools: %v", err)
		}
		tools = append(tools, tool)
	}

	if len(tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(tools))
	}

	// Verify the tool has an InputSchema with an operation enum.
	var schema map[string]any
	raw, err := json.Marshal(tools[0].InputSchema)
	if err != nil {
		t.Fatalf("marshal input schema: %v", err)
	}
	if err := json.Unmarshal(raw, &schema); err != nil {
		t.Fatalf("unmarshal input schema: %v", err)
	}

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatal("expected properties in schema")
	}

	opProp, ok := props["operation"].(map[string]any)
	if !ok {
		t.Fatal("expected operation property in schema")
	}

	enumRaw, ok := opProp["enum"].([]any)
	if !ok {
		t.Fatal("expected enum in operation property")
	}

	if len(enumRaw) != 4 {
		t.Errorf("expected 4 operations in enum, got %d", len(enumRaw))
	}

	expectedOps := map[string]bool{
		"listItems":  true,
		"getItem":    true,
		"createItem": true,
		"deleteItem": true,
	}
	for _, op := range enumRaw {
		opStr, ok := op.(string)
		if !ok {
			t.Errorf("expected string enum value, got %T", op)
			continue
		}
		if !expectedOps[opStr] {
			t.Errorf("unexpected operation in enum: %s", opStr)
		}
	}
}

func TestToolHandler_GeneratedPRToolGroup(t *testing.T) {
	// Smoke test: verify the generated PRToolGroup can be registered.
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "v0.0.1"}, nil)
	mcptools.RegisterToolGroup(server, mcptools.PRToolGroup)

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

	var tools []*mcp.Tool
	for tool, err := range clientSession.Tools(ctx, nil) {
		if err != nil {
			t.Fatalf("listing tools: %v", err)
		}
		tools = append(tools, tool)
	}

	if len(tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(tools))
	}
	if tools[0].Name != "bitbucket_pr" {
		t.Errorf("expected tool name 'bitbucket_pr', got %q", tools[0].Name)
	}
}

func TestToolHandler_AllGeneratedToolGroups(t *testing.T) {
	// Smoke test: verify all generated tool groups can be registered.
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "v0.0.1"}, nil)

	groups := []mcptools.ToolGroup{
		mcptools.PRToolGroup,
		mcptools.HooksToolGroup,
		mcptools.SearchToolGroup,
		mcptools.RefsToolGroup,
		mcptools.CommitsToolGroup,
		mcptools.ReportsToolGroup,
		mcptools.ReposToolGroup,
		mcptools.WorkspacesToolGroup,
		mcptools.ProjectsToolGroup,
		mcptools.PipelinesToolGroup,
		mcptools.IssuesToolGroup,
		mcptools.SnippetsToolGroup,
		mcptools.DeploymentsToolGroup,
		mcptools.BranchRestrictionsToolGroup,
		mcptools.BranchingModelToolGroup,
		mcptools.CommitStatusesToolGroup,
		mcptools.DownloadsToolGroup,
		mcptools.UsersToolGroup,
		mcptools.PropertiesToolGroup,
		mcptools.AddonToolGroup,
	}

	for _, g := range groups {
		mcptools.RegisterToolGroup(server, g)
	}

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

	var tools []*mcp.Tool
	for tool, err := range clientSession.Tools(ctx, nil) {
		if err != nil {
			t.Fatalf("listing tools: %v", err)
		}
		tools = append(tools, tool)
	}

	if len(tools) != 20 {
		t.Errorf("expected 20 tools, got %d", len(tools))
	}
}
