// Package mcptools provides the MCP server tool handler for Bitbucket API operations.
//
// Each command group (pr, hooks, repos, etc.) is registered as a single MCP tool
// with an "operation" parameter that selects which API operation to execute.
// This CRUD-combined design minimizes the number of tools while exposing all
// Bitbucket API operations, matching patterns used by Terraform providers.
package mcptools

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/FabianSchurig/bitbucket-cli/internal/client"
	"github.com/FabianSchurig/bitbucket-cli/internal/handlers"
)

// ─── Operation metadata (shared with generators) ─────────────────────────────

// ParamDef describes a single API parameter.
type ParamDef struct {
	Name     string // API parameter name (e.g., "workspace")
	In       string // "path" or "query"
	Type     string // "string", "integer", "boolean"
	Required bool
}

// BodyFieldDef describes a flattened request body field.
type BodyFieldDef struct {
	Path string // dot-separated path (e.g., "source.branch.name")
	Type string // "string", "integer", "boolean"
	Desc string // human-readable description
}

// OperationDef holds metadata for a single Bitbucket API operation.
type OperationDef struct {
	OperationID string
	Method      string
	Path        string
	Summary     string
	Description string
	Params      []ParamDef
	BodyFields  []BodyFieldDef
	HasBody     bool
	Paginated   bool
}

// ToolGroup holds a set of related operations registered as a single MCP tool.
type ToolGroup struct {
	Name        string
	Description string
	Operations  []OperationDef
}

// ─── Tool registration ───────────────────────────────────────────────────────

// RegisterToolGroup registers a ToolGroup as a single MCP tool on the server.
// The tool accepts an "operation" parameter to select which API operation to
// execute, plus all parameters from all operations (validated at runtime).
func RegisterToolGroup(server *mcp.Server, group ToolGroup) {
	properties := buildToolProperties(group)

	tool := mcp.Tool{
		Name:        group.Name,
		Description: group.Description,
		InputSchema: mustMarshal(map[string]any{
			"type":       "object",
			"properties": properties,
			"required":   []string{"operation"},
		}),
	}

	// Build operation lookup for the handler.
	opMap := make(map[string]OperationDef, len(group.Operations))
	for _, op := range group.Operations {
		opMap[op.OperationID] = op
	}

	server.AddTool(&tool, newToolHandler(opMap))
}

func buildToolProperties(group ToolGroup) map[string]any {
	properties := map[string]any{
		"operation": operationProperty(group.Operations),
		"body": map[string]any{
			"type":        "string",
			"description": "Raw JSON request body (for create/update operations). If provided, body field parameters are ignored.",
		},
	}
	addParamProperties(properties, group.Operations)
	addBodyFieldProperties(properties, group.Operations)
	return properties
}

func operationProperty(operations []OperationDef) map[string]any {
	opIDs := make([]any, 0, len(operations))
	opDescs := make([]string, 0, len(operations))
	for _, op := range operations {
		opIDs = append(opIDs, op.OperationID)
		opDescs = append(opDescs, fmt.Sprintf("- %s: %s [%s %s]", op.OperationID, opDescription(op), op.Method, op.Path))
	}
	return map[string]any{
		"type":        "string",
		"enum":        opIDs,
		"description": "The API operation to execute.\n\n" + strings.Join(opDescs, "\n"),
	}
}

func opDescription(op OperationDef) string {
	if op.Summary != "" {
		return op.Summary
	}
	return op.OperationID
}

func addParamProperties(properties map[string]any, operations []OperationDef) {
	paramSeen := map[string]bool{}
	for _, op := range operations {
		for _, p := range op.Params {
			if paramSeen[p.Name] {
				continue
			}
			paramSeen[p.Name] = true
			properties[p.Name] = map[string]any{
				"type":        jsonSchemaType(p.Type),
				"description": fmt.Sprintf("%s parameter", p.In),
			}
		}
	}
}

func addBodyFieldProperties(properties map[string]any, operations []OperationDef) {
	bodyFieldSeen := map[string]bool{}
	for _, op := range operations {
		for _, bf := range op.BodyFields {
			key := bodyFieldKey(bf.Path)
			if bodyFieldSeen[key] {
				continue
			}
			bodyFieldSeen[key] = true
			properties[key] = map[string]any{
				"type":        jsonSchemaType(bf.Type),
				"description": fmt.Sprintf("Body field: %s", bodyFieldDescription(bf)),
			}
		}
	}
}

// newToolHandler returns an MCP tool handler that dispatches to the Bitbucket API
// based on the "operation" parameter.
func newToolHandler(opMap map[string]OperationDef) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args, err := parseToolArgs(req)
		if err != nil {
			return errResult(fmt.Sprintf("invalid arguments: %v", err)), nil
		}

		opID, op, toolResult := resolveOperation(args, opMap)
		if toolResult != nil {
			return toolResult, nil
		}

		pathParams, queryParams, toolResult := buildRequestParams(args, op, opID)
		if toolResult != nil {
			return toolResult, nil
		}

		body, err := buildRequestBody(args, op)
		if err != nil {
			return errResult(fmt.Sprintf("invalid body: %v", err)), nil
		}

		c, err := client.NewClient()
		if err != nil {
			return errResult(fmt.Sprintf("authentication error: %v", err)), nil
		}

		result, err := handlers.DispatchRaw(ctx, c, handlers.Request{
			Method:      op.Method,
			URLTemplate: op.Path,
			PathParams:  pathParams,
			QueryParams: queryParams,
			Body:        body,
			All:         op.Paginated,
		})
		if err != nil {
			return errResult(fmt.Sprintf("API error: %v", err)), nil
		}

		return formatToolResult(result)
	}
}

func parseToolArgs(req *mcp.CallToolRequest) (map[string]any, error) {
	var args map[string]any
	if len(req.Params.Arguments) > 0 {
		if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
			return nil, err
		}
	}
	if args == nil {
		args = map[string]any{}
	}
	return args, nil
}

func resolveOperation(args map[string]any, opMap map[string]OperationDef) (string, OperationDef, *mcp.CallToolResult) {
	opID, _ := args["operation"].(string)
	if opID == "" {
		return "", OperationDef{}, errResult("missing required parameter: operation")
	}
	op, ok := opMap[opID]
	if !ok {
		return "", OperationDef{}, errResult(fmt.Sprintf("unknown operation: %s", opID))
	}
	return opID, op, nil
}

func buildRequestParams(args map[string]any, op OperationDef, opID string) (map[string]string, map[string]string, *mcp.CallToolResult) {
	pathParams := map[string]string{}
	queryParams := map[string]string{}
	for _, p := range op.Params {
		val := extractStringParam(args, p.Name, p.Type)
		if val == "" {
			if p.Required {
				return nil, nil, errResult(fmt.Sprintf("missing required parameter: %s (for operation %s)", p.Name, opID))
			}
			continue
		}
		assignRequestParam(pathParams, queryParams, p, val)
	}
	return pathParams, queryParams, nil
}

func assignRequestParam(pathParams, queryParams map[string]string, p ParamDef, val string) {
	switch p.In {
	case "path":
		pathParams[p.Name] = val
	case "query":
		queryParams[p.Name] = val
	}
}

func buildRequestBody(args map[string]any, op OperationDef) (string, error) {
	if rawBody, ok := args["body"].(string); ok && rawBody != "" {
		return rawBody, nil
	}
	if !op.HasBody {
		return "", nil
	}
	bodyObj := map[string]any{}
	for _, bf := range op.BodyFields {
		if val, ok := bodyFieldValue(args, bf.Path); ok {
			handlers.SetNested(bodyObj, bf.Path, val)
		}
	}
	if len(bodyObj) == 0 {
		return "", nil
	}
	b, err := json.Marshal(bodyObj)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func bodyFieldValue(args map[string]any, path string) (any, bool) {
	key := bodyFieldKey(path)
	val, ok := args[key]
	if !ok || val == nil || val == "" {
		return nil, false
	}
	return val, true
}

func bodyFieldKey(path string) string {
	return "body_" + strings.ReplaceAll(path, ".", "_")
}

func bodyFieldDescription(bf BodyFieldDef) string {
	if bf.Desc != "" {
		return bf.Desc
	}
	return bf.Path
}

func formatToolResult(result any) (*mcp.CallToolResult, error) {
	if result == nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "OK"}},
		}, nil
	}

	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return errResult(fmt.Sprintf("failed to format response: %v", err)), nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: string(jsonBytes)}},
	}, nil
}

// ─── Filtering ────────────────────────────────────────────────────────────────

// FilterOperations returns a copy of the ToolGroup with only operations whose
// HTTP method passes the filter function. The description is rebuilt to reflect
// the surviving operations.
func FilterOperations(group ToolGroup, allow func(method string) bool) ToolGroup {
	var kept []OperationDef
	for _, op := range group.Operations {
		if allow(op.Method) {
			kept = append(kept, op)
		}
	}
	filtered := ToolGroup{
		Name:       group.Name,
		Operations: kept,
	}
	// Rebuild description to reflect surviving operations.
	if len(kept) > 0 {
		filtered.Description = rebuildDescription(group.Description, kept)
	}
	return filtered
}

// rebuildDescription preserves the first line of the original description and
// rebuilds the operation list from the surviving operations.
func rebuildDescription(original string, ops []OperationDef) string {
	firstLine := original
	if idx := strings.Index(original, "\n"); idx >= 0 {
		firstLine = original[:idx]
	}
	var sb strings.Builder
	sb.WriteString(firstLine)
	sb.WriteString("\n\nAvailable operations:\n")
	for _, op := range ops {
		desc := op.Summary
		if desc == "" {
			desc = op.OperationID
		}
		fmt.Fprintf(&sb, "- %s: %s [%s]\n", op.OperationID, desc, op.Method)
	}
	return sb.String()
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func jsonSchemaType(oaType string) string {
	switch oaType {
	case "integer":
		return "integer"
	case "boolean":
		return "boolean"
	default:
		return "string"
	}
}

func extractStringParam(args map[string]any, name, oaType string) string {
	val, ok := args[name]
	if !ok || val == nil {
		return ""
	}
	switch v := val.(type) {
	case string:
		return v
	case float64:
		if oaType == "integer" {
			return strconv.Itoa(int(v))
		}
		return fmt.Sprintf("%g", v)
	case bool:
		return strconv.FormatBool(v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func errResult(msg string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{&mcp.TextContent{Text: msg}},
		IsError: true,
	}
}

func mustMarshal(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
