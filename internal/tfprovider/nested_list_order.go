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

// stableItemPrimaryKey returns ONLY the identity-field prefix of an item's
// canonical sort key, omitting the canonical-JSON tiebreaker. It is the
// matching key used by `reorderResponseArrayBySourceKeys` to pair an API
// response item against its planned-state counterpart.
//
// The tiebreaker must be omitted because the planned-state side of that
// comparison routinely has Unknown values for `Optional+Computed` inner
// attributes (`display_name`, `created_on`, ...) that the API fills in,
// while the API-response side has all of them resolved. Including those
// fields in the matching key would make every pairing miss and the
// reorderer would silently fall back to API order — defeating the post-
// Create / post-Update consistency check.
//
// When no identity field is present (or none is declared on the schema),
// the canonical-JSON form is still used so two genuinely-different items
// don't collide; that's a degraded but defensible fallback for endpoints
// without a stable primary key.
func stableItemPrimaryKey(m map[string]any, fields []BodyFieldDef) string {
	declared := map[string]bool{}
	for _, f := range fields {
		declared[f.Path] = true
	}
	for _, candidate := range stableIdentityFieldOrder {
		if !declared[candidate] && len(fields) > 0 {
			continue
		}
		if v, ok := m[candidate]; ok && v != nil {
			s := stringifyIdentityValue(v)
			if s != "" {
				return candidate + "=" + s
			}
		}
	}
	// Fallback to the full canonical key when no identity field is present.
	return canonicalJSONKey(m)
}

// stableObjectPrimaryKey is the `types.Object` companion to
// `stableItemPrimaryKey`. It must produce byte-for-byte the same key for
// the same logical item so cross-domain matching works — see the godoc on
// `stableItemPrimaryKey` for why the tiebreaker is omitted.
func stableObjectPrimaryKey(obj types.Object, fields []BodyFieldDef) string {
	attrs := obj.Attributes()
	declared := map[string]bool{}
	for _, f := range fields {
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
	return "obj=" + obj.String()
}

// reorderResponseArrayBySourceKeys reshuffles `arr` so its elements appear
// in the same order as `sourceKeys` (which is the sequence of primary
// identity keys read from the planned/prior state). Items in the API
// response whose primary key isn't in `sourceKeys` are appended at the end
// in their original API order, so adding a new uuid to the HCL surfaces as
// a single-element diff in the operator's chosen position.
//
// This is the central piece that satisfies Terraform Core's post-apply
// consistency check: by reordering the API response to match what the
// operator wrote, the post-apply state is positionally identical to the
// planned state, eliminating the "Provider produced inconsistent result
// after apply" error without silently rewriting the operator's HCL at plan
// time.
func reorderResponseArrayBySourceKeys(arr []any, fields []BodyFieldDef, sourceKeys []string) []any {
	if len(sourceKeys) == 0 || len(arr) == 0 {
		return arr
	}
	keyToItem := make(map[string]any, len(arr))
	for _, it := range arr {
		m, ok := it.(map[string]any)
		if !ok {
			continue
		}
		k := stableItemPrimaryKey(m, fields)
		if _, exists := keyToItem[k]; !exists {
			keyToItem[k] = it
		}
	}
	out := make([]any, 0, len(arr))
	emitted := make(map[string]bool, len(arr))
	for _, k := range sourceKeys {
		if emitted[k] {
			continue
		}
		if it, ok := keyToItem[k]; ok {
			out = append(out, it)
			emitted[k] = true
		}
	}
	for _, it := range arr {
		m, ok := it.(map[string]any)
		if !ok {
			out = append(out, it)
			continue
		}
		k := stableItemPrimaryKey(m, fields)
		if emitted[k] {
			continue
		}
		out = append(out, it)
		emitted[k] = true
	}
	return out
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
