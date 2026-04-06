package mcptools

import (
	"encoding/json"
	"math"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestToolHelperFunctions(t *testing.T) {
	if got := opDescription(OperationDef{OperationID: "fallback"}); got != "fallback" {
		t.Fatalf("expected operation ID fallback, got %q", got)
	}
	if got := bodyFieldDescription(BodyFieldDef{Path: "title"}); got != "title" {
		t.Fatalf("expected body field path fallback, got %q", got)
	}
	if got := bodyFieldKey("content.raw"); got != "body_content_raw" {
		t.Fatalf("unexpected body field key %q", got)
	}

	props := buildToolProperties(ToolGroup{
		Operations: []OperationDef{
			{OperationID: "listItems", Method: "GET", Path: "/items", Params: []ParamDef{{Name: "workspace", In: "path", Type: "string"}}, BodyFields: []BodyFieldDef{{Path: "title", Type: "string"}}},
			{OperationID: "createItem", Method: "POST", Path: "/items", Params: []ParamDef{{Name: "workspace", In: "path", Type: "string"}}, BodyFields: []BodyFieldDef{{Path: "title", Type: "string"}}},
		},
	})
	if _, ok := props["workspace"]; !ok {
		t.Fatal("expected deduplicated parameter property")
	}
	if _, ok := props["body_title"]; !ok {
		t.Fatal("expected deduplicated body field property")
	}

	args, err := parseToolArgs(&mcp.CallToolRequest{Params: &mcp.CallToolParamsRaw{}})
	if err != nil || len(args) != 0 {
		t.Fatalf("expected empty args, got %#v err=%v", args, err)
	}
	if _, err := parseToolArgs(&mcp.CallToolRequest{Params: &mcp.CallToolParamsRaw{Arguments: json.RawMessage(`{`)}}); err == nil {
		t.Fatal("expected invalid JSON arguments error")
	}

	if _, _, result := resolveOperation(map[string]any{}, map[string]OperationDef{}); result == nil {
		t.Fatal("expected missing operation error result")
	}
	if _, _, result := buildRequestParams(map[string]any{}, OperationDef{Params: []ParamDef{{Name: "workspace", In: "path", Required: true}}}, "listItems"); result == nil {
		t.Fatal("expected missing required parameter result")
	}

	body, err := buildRequestBody(map[string]any{"body_title": "Demo"}, OperationDef{HasBody: true, BodyFields: []BodyFieldDef{{Path: "title"}}})
	if err != nil || body == "" {
		t.Fatalf("expected built request body, got %q err=%v", body, err)
	}
	if body, err := buildRequestBody(map[string]any{}, OperationDef{}); err != nil || body != "" {
		t.Fatalf("expected empty body for bodyless operation, got %q err=%v", body, err)
	}
	if _, err := buildRequestBody(map[string]any{"body_title": math.Inf(1)}, OperationDef{HasBody: true, BodyFields: []BodyFieldDef{{Path: "title"}}}); err == nil {
		t.Fatal("expected marshal error for unsupported JSON value")
	}

	if result, err := formatToolResult(nil); err != nil || result == nil || result.IsError {
		t.Fatalf("expected OK tool result, got %#v err=%v", result, err)
	}
	if result, err := formatToolResult(map[string]any{"value": math.Inf(1)}); err != nil || result == nil || !result.IsError {
		t.Fatalf("expected formatting error result, got %#v err=%v", result, err)
	}

	if got := extractStringParam(map[string]any{"value": 7}, "value", "integer"); got != "7" {
		t.Fatalf("expected integer conversion, got %q", got)
	}
	if got := extractStringParam(map[string]any{"value": true}, "value", "boolean"); got != "true" {
		t.Fatalf("expected boolean conversion, got %q", got)
	}
	if got := extractStringParam(map[string]any{"value": 7}, "value", "string"); got != "7" {
		t.Fatalf("expected default conversion, got %q", got)
	}

	if got := jsonSchemaType("unknown"); got != "string" {
		t.Fatalf("expected default JSON schema type, got %q", got)
	}
	if data := mustMarshal(map[string]any{"ok": true}); len(data) == 0 {
		t.Fatal("expected marshaled schema data")
	}
}

func TestFilterOperations(t *testing.T) {
	group := ToolGroup{
		Name:        "test_tool",
		Description: "Test tool\n\nAvailable operations:\n- listItems: List [GET]\n- deleteItem: Delete [DELETE]\n",
		Operations: []OperationDef{
			{OperationID: "listItems", Method: "GET", Path: "/items", Summary: "List"},
			{OperationID: "createItem", Method: "POST", Path: "/items", Summary: "Create"},
			{OperationID: "deleteItem", Method: "DELETE", Path: "/items/{id}", Summary: "Delete"},
		},
	}

	// Filter to only GET and POST.
	allowGETandPOST := func(method string) bool {
		return method == "GET" || method == "POST"
	}
	filtered := FilterOperations(group, allowGETandPOST)

	if len(filtered.Operations) != 2 {
		t.Fatalf("expected 2 operations, got %d", len(filtered.Operations))
	}
	if filtered.Operations[0].OperationID != "listItems" {
		t.Errorf("expected listItems, got %s", filtered.Operations[0].OperationID)
	}
	if filtered.Operations[1].OperationID != "createItem" {
		t.Errorf("expected createItem, got %s", filtered.Operations[1].OperationID)
	}
	if filtered.Name != "test_tool" {
		t.Errorf("expected name test_tool, got %s", filtered.Name)
	}

	// Filter out everything.
	allowNone := func(string) bool { return false }
	empty := FilterOperations(group, allowNone)
	if len(empty.Operations) != 0 {
		t.Errorf("expected 0 operations, got %d", len(empty.Operations))
	}

	// Allow all.
	allowAll := func(string) bool { return true }
	all := FilterOperations(group, allowAll)
	if len(all.Operations) != 3 {
		t.Errorf("expected 3 operations, got %d", len(all.Operations))
	}
}

func TestAllToolGroups_Registry(t *testing.T) {
	if len(AllToolGroups) == 0 {
		t.Fatal("expected AllToolGroups to be populated by generated init() functions")
	}
	found := false
	for _, g := range AllToolGroups {
		if g.Name == "bitbucket_pr" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected bitbucket_pr in AllToolGroups registry")
	}
}
