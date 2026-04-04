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

// toSnakeCase converts parameter names like "repo_slug", "repoSlug", or
// "target.commit.hash" to Terraform-compatible snake_case attribute names.
func toSnakeCase(s string) string {
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, ".", "_")

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

// ParamAttrName returns the Terraform attribute name for an API parameter.
// If the snake_case name collides with reserved Terraform attributes (like "id"),
// it is prefixed with "param_" to avoid conflicts.
func ParamAttrName(paramName string) string {
	name := toSnakeCase(paramName)
	if name == "id" {
		return "param_id"
	}
	return name
}

// MapCRUDOps resolves CRUD operations for a resource group by looking up
// operation IDs from the hand-written CRUDConfig map. The typeName parameter
// identifies the resource group (e.g., "repos", "pr"). Called at runtime by
// generated init() functions to map Bitbucket API operations to Terraform
// CRUD lifecycle methods.
func MapCRUDOps(typeName string, ops []OperationDef) CRUDOps {
	cfg, ok := CRUDConfig[typeName]
	if !ok {
		return CRUDOps{}
	}

	// Build an index of operation ID → *OperationDef for fast lookup.
	index := make(map[string]*OperationDef, len(ops))
	for i := range ops {
		index[ops[i].OperationID] = &ops[i]
	}

	return CRUDOps{
		Create: index[cfg.Create],
		Read:   index[cfg.Read],
		Update: index[cfg.Update],
		Delete: index[cfg.Delete],
		List:   index[cfg.List],
	}
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
