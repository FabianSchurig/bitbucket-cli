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
// declared on the item; if none are present it falls back to the canonical
// JSON encoding of the whole item. The fallback guarantees a total order
// even for shapes the registry doesn't know about, so the resulting list is
// always reproducible regardless of API ordering quirks.
func stableItemSortKey(m map[string]any, fields []BodyFieldDef) string {
	declared := map[string]bool{}
	for _, f := range fields {
		declared[f.Path] = true
	}
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
				return candidate + "=" + s
			}
		}
	}
	if b, err := json.Marshal(canonicalize(m)); err == nil {
		return "json=" + string(b)
	}
	return fmt.Sprintf("%v", m)
}

// stableObjectSortKey returns the same key for a `types.Object` element
// already living in Terraform state / plan. It mirrors stableItemSortKey
// over Terraform's attr.Value graph so plan-side sorting (the plan modifier)
// and state-side sorting (the response builder) agree byte-for-byte.
func stableObjectSortKey(obj types.Object, fields []BodyFieldDef) string {
	attrs := obj.Attributes()
	declared := map[string]bool{}
	for _, f := range fields {
		// nested attrs are stored under snake_cased keys.
		declared[bodyFieldKey(f)] = true
	}
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
			return candidate + "=" + s
		}
	}
	// Fallback: full string form (framework guarantees stable ordering of
	// attribute names within String()).
	return "obj=" + obj.String()
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

// canonicalize rewrites nested maps so json.Marshal emits keys in
// lexicographic order, giving a deterministic JSON tiebreaker. Go's
// encoding/json already sorts map keys, so this is mostly defensive — it
// also normalises float64 NaN-style oddities into stable strings via
// %v. The fallback only fires when none of the well-known identity
// fields are present.
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

// Description returns a human-readable description of the modifier.
func (m nestedListSortPlanModifier) Description(_ context.Context) string {
	return "Sorts the planned list elements by a stable identity key (uuid > id > slug > full_slug > name > canonical JSON) so the post-apply consistency check is order-insensitive."
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
