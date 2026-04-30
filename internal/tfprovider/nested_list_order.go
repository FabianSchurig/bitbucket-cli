package tfprovider

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// stableIdentityFieldOrder is the precedence used to derive a stable per-item
// sort key for nested-object arrays. Bitbucket's REST API consistently
// exposes one of these as the natural primary key on collection items
// (uuid for users/repositories, id for branch restrictions, slug for
// groups/repositories, name for tags). Using a fixed precedence rather than
// per-endpoint configuration keeps the codegen pipeline schema-driven and
// avoids hand-maintained sort tables.
//
// The keys are now consumed exclusively by setLikeListValue.ListSemanticEquals
// to decide whether two nested-object lists carry the same multiset of
// elements regardless of order. The runtime no longer reorders user-supplied
// or API-returned lists in place — order-insensitivity is expressed at the
// value-type level via setLikeListType.
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
// already living in Terraform state / plan / config. It mirrors
// stableItemSortKey over Terraform's attr.Value graph so the multiset
// comparison performed by setLikeListValue.ListSemanticEquals lines up
// byte-for-byte with the canonical key used for raw JSON responses —
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
