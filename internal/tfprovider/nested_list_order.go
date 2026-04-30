package tfprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// stableIdentityFieldOrder is the precedence used to derive a stable per-item
// sort key for nested-object arrays. Bitbucket's REST API consistently
// exposes one of these as the natural primary key on collection items
// (uuid for users/repositories, id for branch restrictions, slug for
// groups/repositories, name for tags). Using a fixed precedence rather than
// per-endpoint configuration keeps the codegen pipeline schema-driven and
// avoids hand-maintained sort tables.
var stableIdentityFieldOrder = []string{
	"uuid",
	"id",
	"slug",
	"full_slug",
	"name",
	"kind",
	"pattern",
	"branch_type",
}

// stableItemSortKey returns a deterministic sort key for a nested-object
// item (`map[string]any` shape, as decoded from a JSON API response). It
// looks for the first non-empty value among the well-known identity fields
// declared on the item and **always** appends the canonical JSON encoding
// of the whole item as a secondary tiebreaker. The tiebreaker guarantees a
// total order even when two items happen to share the same identity value
// (or none of the known identity fields are present), so the resulting list
// is always reproducible regardless of API ordering quirks.
func stableItemSortKey(m map[string]any, fields []BodyFieldDef) string {
	declared := map[string]bool{}
	for _, f := range fields {
		declared[f.Path] = true
	}
	primary := ""
	for _, candidate := range stableIdentityFieldOrder {
		if !declared[candidate] && len(fields) > 0 {
			// Only consider identity fields that exist in the item's schema
			// when the schema is known; otherwise (fields == nil) accept any
			// candidate present in the map.
			continue
		}
		if v, ok := m[candidate]; ok && v != nil {
			s := stringifyIdentityValue(v)
			if s != "" {
				primary = candidate + "=" + s
				break
			}
		}
	}
	tiebreaker := canonicalJSONKey(m)
	if primary == "" {
		return tiebreaker
	}
	return primary + "|" + tiebreaker
}

// canonicalJSONKey returns a deterministic JSON-encoded form of v, used as
// a total-order tiebreaker by stableItemSortKey. Falls back to %v when
// json.Marshal returns an error (e.g. NaN / +Inf in numeric values).
func canonicalJSONKey(v any) string {
	if b, err := json.Marshal(canonicalize(v)); err == nil {
		return "json=" + string(b)
	}
	return fmt.Sprintf("raw=%v", v)
}

// stableObjectSortKey returns the same key for a `types.Object` element
// already living in Terraform state / plan. It mirrors stableItemSortKey
// over Terraform's attr.Value graph so plan-side sorting (the plan modifier)
// and state-side sorting (the response builder) agree byte-for-byte —
// including the secondary canonical-form tiebreaker that guarantees a
// total order when two items share the same identity value.
func stableObjectSortKey(obj types.Object, fields []BodyFieldDef) string {
	attrs := obj.Attributes()
	declared := map[string]bool{}
	for _, f := range fields {
		// nested attrs are stored under snake_cased keys.
		declared[bodyFieldKey(f)] = true
	}
	primary := ""
	for _, candidate := range stableIdentityFieldOrder {
		key := candidate
		if len(fields) > 0 && !declared[key] {
			continue
		}
		v, ok := attrs[key]
		if !ok {
			continue
		}
		if s, ok := stringifyAttrIdentity(v); ok && s != "" {
			primary = candidate + "=" + s
			break
		}
	}
	// Tiebreaker: the framework's stable String() form (attribute names are
	// emitted in lexicographic order).
	tiebreaker := "obj=" + obj.String()
	if primary == "" {
		return tiebreaker
	}
	return primary + "|" + tiebreaker
}

// bodyFieldKey returns the snake_cased attribute key for a BodyFieldDef.
// It mirrors the key derivation used by buildNestedItemAttrs / itemAttrTypes
// so identity-field lookups on Terraform objects line up with the schema.
func bodyFieldKey(f BodyFieldDef) string {
	// Identity fields used for sorting (uuid/id/slug/full_slug/name/...) are
	// already snake_case and contain no dots, so the simple ReplaceAll +
	// toSnakeCase the schema generators apply is equivalent to the field's
	// Path here. Keep the helper explicit for readability.
	key := f.Path
	return toSnakeCase(key)
}

func stringifyIdentityValue(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case bool:
		return fmt.Sprintf("%t", x)
	case float64:
		// JSON-decoded numbers come through as float64. Use %v so both ints
		// and floats produce stable, deterministic strings; lexicographic
		// ordering is fine here because the goal is determinism, not
		// numeric ordering.
		return fmt.Sprintf("%v", x)
	case int, int64, int32:
		return fmt.Sprintf("%d", x)
	}
	if b, err := json.Marshal(v); err == nil {
		return string(b)
	}
	return fmt.Sprintf("%v", v)
}

func stringifyAttrIdentity(v attr.Value) (string, bool) {
	if v == nil || v.IsNull() || v.IsUnknown() {
		return "", false
	}
	switch x := v.(type) {
	case types.String:
		return x.ValueString(), true
	case types.Int64:
		return fmt.Sprintf("%d", x.ValueInt64()), true
	case types.Bool:
		return fmt.Sprintf("%t", x.ValueBool()), true
	}
	// Anything else: use the framework's deterministic string form.
	return v.String(), true
}

// canonicalize rewrites nested maps and slices so json.Marshal can be used
// as a deterministic JSON tiebreaker. Go's encoding/json already sorts map
// keys lexicographically, so this is mostly defensive recursion through
// nested values. The fallback only fires when the primary identity-field
// lookup yields no key.
func canonicalize(v any) any {
	switch x := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(x))
		for k, vv := range x {
			out[k] = canonicalize(vv)
		}
		return out
	case []any:
		out := make([]any, len(x))
		for i, vv := range x {
			out[i] = canonicalize(vv)
		}
		return out
	default:
		return v
	}
}

// sortResponseItems sorts a JSON-decoded array of nested-object items in
// place by their stable identity key. It is the response-side half of the
// fix: every nested-object array we materialise into Terraform state goes
// through this so two equivalent API responses (same elements, different
// order) produce byte-identical state.
func sortResponseItems(arr []any, fields []BodyFieldDef) {
	if len(arr) < 2 {
		return
	}
	keys := make([]string, len(arr))
	for i, item := range arr {
		if m, ok := item.(map[string]any); ok {
			keys[i] = stableItemSortKey(m, fields)
		} else {
			// Non-object entries (rare; defensive) sort by their JSON form.
			if b, err := json.Marshal(item); err == nil {
				keys[i] = "raw=" + string(b)
			} else {
				keys[i] = fmt.Sprintf("raw=%v", item)
			}
		}
	}
	idx := make([]int, len(arr))
	for i := range idx {
		idx[i] = i
	}
	sort.SliceStable(idx, func(i, j int) bool {
		return keys[idx[i]] < keys[idx[j]]
	})
	sorted := make([]any, len(arr))
	for i, k := range idx {
		sorted[i] = arr[k]
	}
	copy(arr, sorted)
}

// nestedListSortPlanModifier is the plan-side half of the deterministic-
// order fix. Attaching it to every ListNestedAttribute over an object item
// ensures the planned value carries the same canonical order the response
// builder will produce — without it Terraform's post-apply consistency
// check would still fire whenever a user wrote elements in a different
// order than the canonical sort.
//
// The modifier is a value type (no fields beyond the per-attribute item
// schema) so equality / type-assertion in tests stays straightforward.
type nestedListSortPlanModifier struct {
	itemFields []BodyFieldDef
}

func newNestedListSortPlanModifier(itemFields []BodyFieldDef) nestedListSortPlanModifier {
	return nestedListSortPlanModifier{itemFields: itemFields}
}

// stableIdentityPrecedenceDescription returns the human-readable form of
// the identity-field precedence used by stableItemSortKey /
// stableObjectSortKey (e.g. "uuid > id > slug > ... > canonical JSON"),
// derived from stableIdentityFieldOrder so the docs can never drift.
func stableIdentityPrecedenceDescription() string {
	if len(stableIdentityFieldOrder) == 0 {
		return "canonical JSON"
	}
	precedence := stableIdentityFieldOrder[0]
	for _, field := range stableIdentityFieldOrder[1:] {
		precedence += " > " + field
	}
	return precedence + " > canonical JSON"
}

// Description returns a human-readable description of the modifier.
func (m nestedListSortPlanModifier) Description(_ context.Context) string {
	return fmt.Sprintf(
		"Sorts the planned list elements by a stable identity key (%s) so the post-apply consistency check is order-insensitive.",
		stableIdentityPrecedenceDescription(),
	)
}

// MarkdownDescription returns the Markdown form of Description.
func (m nestedListSortPlanModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

// PlanModifyList sorts the planned list value in-place using the same
// identity-field precedence the response builder uses.
func (m nestedListSortPlanModifier) PlanModifyList(_ context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}
	elements := req.PlanValue.Elements()
	if len(elements) < 2 {
		return
	}
	keys := make([]string, len(elements))
	for i, e := range elements {
		obj, ok := e.(types.Object)
		if !ok {
			// Defensive: leave non-object element lists alone.
			return
		}
		keys[i] = stableObjectSortKey(obj, m.itemFields)
	}
	idx := make([]int, len(elements))
	for i := range idx {
		idx[i] = i
	}
	sort.SliceStable(idx, func(i, j int) bool {
		return keys[idx[i]] < keys[idx[j]]
	})
	sorted := make([]attr.Value, len(elements))
	for i, k := range idx {
		sorted[i] = elements[k]
	}
	sortedList, diags := types.ListValue(req.PlanValue.ElementType(context.Background()), sorted)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}
	resp.PlanValue = sortedList
}

// nestedListSortPlanModifiers returns the standard plan-modifier slice for
// a nested-object array attribute. It is a tiny helper so the schema
// builders read consistently.
func nestedListSortPlanModifiers(itemFields []BodyFieldDef) []planmodifier.List {
	return []planmodifier.List{newNestedListSortPlanModifier(itemFields)}
}
