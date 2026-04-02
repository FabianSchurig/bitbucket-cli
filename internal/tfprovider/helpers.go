package tfprovider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// stateAccessor is a common interface for Terraform plan and state objects,
// used to read and write attributes generically.
type stateAccessor interface {
	GetAttribute(ctx context.Context, p path.Path, target interface{}) diag.Diagnostics
	SetAttribute(ctx context.Context, p path.Path, val interface{}) diag.Diagnostics
}

// attrPath creates a terraform-plugin-framework attribute path from a string name.
func attrPath(name string) path.Path {
	return path.Root(name)
}

// toSnakeCase converts parameter names like "repo_slug" or "repoSlug" to
// Terraform-compatible snake_case attribute names.
func toSnakeCase(s string) string {
	s = strings.ReplaceAll(s, "-", "_")

	// Handle camelCase by inserting underscores before uppercase letters.
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			prev := rune(s[i-1])
			if prev >= 'a' && prev <= 'z' {
				result.WriteRune('_')
			}
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// MapCRUDOps analyzes a list of operations and assigns them to CRUD slots based
// on HTTP method heuristics. Called at runtime by generated init() functions to
// map Bitbucket API operations to Terraform CRUD lifecycle methods.
func MapCRUDOps(ops []OperationDef) CRUDOps {
	var crud CRUDOps

	// Score operations by how well they match each CRUD slot.
	// Prefer specific endpoints (with more path params) for single-resource ops.
	var (
		bestCreate, bestRead, bestUpdate, bestDelete, bestList      *OperationDef
		createScore, readScore, updateScore, deleteScore, listScore int
	)

	for i := range ops {
		op := &ops[i]
		score := len(op.Params) // more specific paths score higher
		opLower := strings.ToLower(op.OperationID)

		switch op.Method {
		case "POST":
			if crud.Create == nil || score > createScore || strings.Contains(opLower, "create") {
				bestCreate = op
				createScore = score
			}
		case "GET":
			if op.Paginated || strings.Contains(opLower, "list") || strings.Contains(opLower, "getall") {
				if crud.List == nil || score < listScore { // prefer less-specific for list
					bestList = op
					listScore = score
				}
			} else {
				if crud.Read == nil || score > readScore || strings.Contains(opLower, "get") {
					bestRead = op
					readScore = score
				}
			}
		case "PUT", "PATCH":
			if crud.Update == nil || score > updateScore || strings.Contains(opLower, "update") {
				bestUpdate = op
				updateScore = score
			}
		case "DELETE":
			if crud.Delete == nil || score > deleteScore || strings.Contains(opLower, "delete") {
				bestDelete = op
				deleteScore = score
			}
		}
	}

	crud.Create = bestCreate
	crud.Read = bestRead
	crud.Update = bestUpdate
	crud.Delete = bestDelete
	crud.List = bestList

	return crud
}

// BuildResourceDescription builds a description for a Terraform resource
// from the command group description and its CRUD operations.
func BuildResourceDescription(groupDesc string, crud CRUDOps) string {
	var sb strings.Builder
	sb.WriteString(groupDesc)
	sb.WriteString("\n\nMapped CRUD operations:\n")
	if crud.Create != nil {
		fmt.Fprintf(&sb, "- Create: %s [%s %s]\n", crud.Create.OperationID, crud.Create.Method, crud.Create.Path)
	}
	if crud.Read != nil {
		fmt.Fprintf(&sb, "- Read: %s [%s %s]\n", crud.Read.OperationID, crud.Read.Method, crud.Read.Path)
	}
	if crud.Update != nil {
		fmt.Fprintf(&sb, "- Update: %s [%s %s]\n", crud.Update.OperationID, crud.Update.Method, crud.Update.Path)
	}
	if crud.Delete != nil {
		fmt.Fprintf(&sb, "- Delete: %s [%s %s]\n", crud.Delete.OperationID, crud.Delete.Method, crud.Delete.Path)
	}
	if crud.List != nil {
		fmt.Fprintf(&sb, "- List: %s [%s %s]\n", crud.List.OperationID, crud.List.Method, crud.List.Path)
	}
	return sb.String()
}
