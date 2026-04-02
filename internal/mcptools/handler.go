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
	// Build the JSON Schema for the tool input.
	properties := map[string]any{}
	required := []string{"operation"}

	// Operation enum from all operation IDs.
	opIDs := make([]any, 0, len(group.Operations))
	opDescs := make([]string, 0, len(group.Operations))
	for _, op := range group.Operations {
		opIDs = append(opIDs, op.OperationID)
		desc := op.Summary
		if desc == "" {
			desc = op.OperationID
		}
		opDescs = append(opDescs, fmt.Sprintf("- %s: %s [%s %s]", op.OperationID, desc, op.Method, op.Path))
	}

	properties["operation"] = map[string]any{
		"type":        "string",
		"enum":        opIDs,
		"description": "The API operation to execute.\n\n" + strings.Join(opDescs, "\n"),
	}

	// Collect all unique parameter names across operations.
	paramSeen := map[string]bool{}
	for _, op := range group.Operations {
		for _, p := range op.Params {
			if paramSeen[p.Name] {
				continue
			}
			paramSeen[p.Name] = true
			prop := map[string]any{
				"type":        jsonSchemaType(p.Type),
				"description": fmt.Sprintf("%s parameter", p.In),
			}
			properties[p.Name] = prop
		}
	}

	// Add body parameter for raw JSON body.
	properties["body"] = map[string]any{
		"type":        "string",
		"description": "Raw JSON request body (for create/update operations). If provided, body field parameters are ignored.",
	}

	// Collect unique body fields across operations.
	bodyFieldSeen := map[string]bool{}
	for _, op := range group.Operations {
		for _, bf := range op.BodyFields {
			key := "body_" + strings.ReplaceAll(bf.Path, ".", "_")
			if bodyFieldSeen[key] {
				continue
			}
			bodyFieldSeen[key] = true
			desc := bf.Desc
			if desc == "" {
				desc = bf.Path
			}
			properties[key] = map[string]any{
				"type":        jsonSchemaType(bf.Type),
				"description": fmt.Sprintf("Body field: %s", desc),
			}
		}
	}

	tool := mcp.Tool{
		Name:        group.Name,
		Description: group.Description,
		InputSchema: mustMarshal(map[string]any{
			"type":       "object",
			"properties": properties,
			"required":   required,
		}),
	}

	// Build operation lookup for the handler.
	opMap := make(map[string]OperationDef, len(group.Operations))
	for _, op := range group.Operations {
		opMap[op.OperationID] = op
	}

	server.AddTool(&tool, newToolHandler(opMap))
}

// newToolHandler returns an MCP tool handler that dispatches to the Bitbucket API
// based on the "operation" parameter.
func newToolHandler(opMap map[string]OperationDef) mcp.ToolHandler {
	return func(ctx context.Context, req *mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Unmarshal raw JSON arguments into a map.
		var args map[string]any
		if len(req.Params.Arguments) > 0 {
			if err := json.Unmarshal(req.Params.Arguments, &args); err != nil {
				return errResult(fmt.Sprintf("invalid arguments: %v", err)), nil
			}
		}
		if args == nil {
			args = map[string]any{}
		}

		opID, _ := args["operation"].(string)
		if opID == "" {
			return errResult("missing required parameter: operation"), nil
		}

		op, ok := opMap[opID]
		if !ok {
			return errResult(fmt.Sprintf("unknown operation: %s", opID)), nil
		}

		// Build path parameters.
		pathParams := map[string]string{}
		queryParams := map[string]string{}
		for _, p := range op.Params {
			val := extractStringParam(args, p.Name, p.Type)
			if val == "" {
				if p.Required {
					return errResult(fmt.Sprintf("missing required parameter: %s (for operation %s)", p.Name, opID)), nil
				}
				continue
			}
			switch p.In {
			case "path":
				pathParams[p.Name] = val
			case "query":
				queryParams[p.Name] = val
			}
		}

		// Build request body.
		body := ""
		if rawBody, ok := args["body"].(string); ok && rawBody != "" {
			body = rawBody
		} else if op.HasBody {
			bodyObj := map[string]any{}
			for _, bf := range op.BodyFields {
				key := "body_" + strings.ReplaceAll(bf.Path, ".", "_")
				if val, ok := args[key]; ok && val != nil && val != "" {
					handlers.SetNested(bodyObj, bf.Path, val)
				}
			}
			if len(bodyObj) > 0 {
				b, _ := json.Marshal(bodyObj)
				body = string(b)
			}
		}

		// Create Bitbucket client (uses env vars).
		c, err := client.NewClient()
		if err != nil {
			return errResult(fmt.Sprintf("authentication error: %v", err)), nil
		}

		// Dispatch the API call.
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

		// Format result as JSON text content.
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
