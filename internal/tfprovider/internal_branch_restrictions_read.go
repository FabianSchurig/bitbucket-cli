package tfprovider

import (
	"context"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// transformProjectBranchRestrictionsRead reshapes the response of the
// `getProjectBranchRestrictionsGroupedByBranch` GET endpoint into the
// `{"values": [...]}` form that the Create/Update PUT endpoints declare
// (and that the generic response-field extractor expects).
//
// The internal Bitbucket endpoint
// `/!api/internal/.../branch-restrictions/group-by-branch/` returns either
// an array of entries or an object whose values are arrays of entries.
// Each entry has shape:
//
//	{
//	  "kind": { "<kind_name>": { "users": [...], "groups": [...], "value": <int>? } },
//	  "branch_match_kind": "glob" | "branching_model",
//	  "pattern": "<glob>",
//	  "branch_type": "<branch-type>",
//	  ...
//	}
//
// Whereas the PUT endpoints accept and return:
//
//	{ "values": [
//	    { "kind": "<kind_name>", "branch_match_kind": "...", "pattern": "...",
//	      "branch_type": "...", "users": [...], "groups": [...], "value": <int> },
//	    ...
//	] }
//
// Without this transformation the generic Read function cannot map the GET
// response back to the `values` attribute, which Terraform then sees as
// `null` and flags as "Provider produced inconsistent result after apply",
// tainting the resource on every apply.
//
// The transform filters entries to those matching the resource's own
// `pattern` (or `branch_type`) — read from the source state — and expands
// each entry's nested `kind` object into one `values` row per kind name.
// `typeName` is the resource group's TypeName (e.g.
// `project-branch-restrictions-by-pattern`); it determines which scope
// attribute to read so we never call `GetAttribute` on a root attribute
// that the active sub-resource's schema does not declare.
func transformProjectBranchRestrictionsRead(
	ctx context.Context,
	op *OperationDef,
	typeName string,
	source stateAccessor,
	result any,
	diags *diag.Diagnostics,
) any {
	if op == nil || op.OperationID != "getProjectBranchRestrictionsGroupedByBranch" {
		return result
	}
	if result == nil {
		return map[string]any{"values": []any{}}
	}

	pattern, branchType := readBranchRestrictionScope(ctx, typeName, source, diags)

	entries := flattenGroupByBranchEntries(result)
	values := make([]any, 0, len(entries))
	for _, entry := range entries {
		if !matchesBranchRestrictionScope(entry, pattern, branchType) {
			continue
		}
		values = append(values, expandKindEntries(entry)...)
	}
	return map[string]any{"values": values}
}

// readBranchRestrictionScope reads the scope attribute (`pattern` or
// `branch_type`) declared by the active sub-resource. Reading an attribute
// that the schema does not declare would surface as a diagnostic error
// from the Terraform plugin framework, so we only call `GetAttribute` for
// the attribute that actually exists on the running sub-resource.
func readBranchRestrictionScope(ctx context.Context, typeName string, source stateAccessor, diags *diag.Diagnostics) (pattern, branchType string) {
	switch {
	case strings.HasSuffix(typeName, "by-pattern"):
		pattern, _ = readBodyStringValue(ctx, source, "pattern", diags)
	case strings.HasSuffix(typeName, "by-branch-type"):
		branchType, _ = readBodyStringValue(ctx, source, "branch_type", diags)
	}
	return pattern, branchType
}

// flattenGroupByBranchEntries normalises the GET response shape (either
// `[]any` or `map[string]any` with array values) into a flat `[]map[string]any`.
// When the response is an object (the schema-declared shape), the keys are
// visited in lexicographic order so the resulting `values` list is stable
// across runs — Go map iteration would otherwise randomise it and produce
// perpetual Terraform diffs on a ListNestedAttribute.
func flattenGroupByBranchEntries(result any) []map[string]any {
	var out []map[string]any
	switch v := result.(type) {
	case []any:
		for _, item := range v {
			if m, ok := item.(map[string]any); ok {
				out = append(out, m)
			}
		}
	case map[string]any:
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			arr, ok := v[k].([]any)
			if !ok {
				continue
			}
			for _, item := range arr {
				if m, ok := item.(map[string]any); ok {
					out = append(out, m)
				}
			}
		}
	}
	return out
}

// matchesBranchRestrictionScope returns true when entry belongs to the
// pattern (for by-pattern resources) or branch_type (for by-branch-type
// resources) declared on the Terraform resource. When the caller did not
// provide a scope (e.g. import-only state), all entries match so the
// resource can still be populated.
func matchesBranchRestrictionScope(entry map[string]any, pattern, branchType string) bool {
	if pattern == "" && branchType == "" {
		return true
	}
	if pattern != "" {
		ep, _ := entry["pattern"].(string)
		if ep != pattern {
			return false
		}
	}
	if branchType != "" {
		bt, _ := entry["branch_type"].(string)
		if bt != branchType {
			return false
		}
	}
	return true
}

// expandKindEntries turns a single GET entry (whose `kind` field is an
// object keyed by kind name) into one or more `values` rows whose `kind`
// field is the flat string the PUT schema expects. Kind names are visited
// in lexicographic order so the resulting `values` list is stable across
// runs — Go map iteration would otherwise randomise it and produce
// perpetual Terraform diffs on a ListNestedAttribute.
func expandKindEntries(entry map[string]any) []any {
	kindObj, ok := entry["kind"].(map[string]any)
	if !ok {
		return nil
	}
	kindNames := make([]string, 0, len(kindObj))
	for k := range kindObj {
		kindNames = append(kindNames, k)
	}
	sort.Strings(kindNames)
	rows := make([]any, 0, len(kindNames))
	for _, kindName := range kindNames {
		row := map[string]any{
			"kind":              kindName,
			"branch_match_kind": entry["branch_match_kind"],
			"pattern":           entry["pattern"],
			"branch_type":       entry["branch_type"],
		}
		applyKindData(row, kindObj[kindName])
		rows = append(rows, row)
	}
	return rows
}

// applyKindData copies users/groups/value from the per-kind payload into the
// flattened row. Bitbucket's GET nests these under each kind name; some
// kinds (e.g. require_approvals_to_merge) carry only a numeric threshold
// and may be encoded as a bare number rather than an object.
func applyKindData(row map[string]any, kindData any) {
	switch d := kindData.(type) {
	case map[string]any:
		if users, ok := d["users"]; ok {
			row["users"] = users
		}
		if groups, ok := d["groups"]; ok {
			row["groups"] = groups
		}
		if value, ok := d["value"]; ok {
			row["value"] = value
		}
	case float64, int, int64:
		row["value"] = d
	}
}
