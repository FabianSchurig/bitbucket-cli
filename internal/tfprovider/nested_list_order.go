package tfprovider

import (
	"encoding/json"
	"fmt"
	"sort"

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

// stableItemPrimaryKey returns just the primary identity key (without
// tiebreaker) for a JSON-decoded map item. Used by
// alignResponseItemsToReference to match response items against prior plan
// elements: matching cares only about the per-item identity field, not the
// tiebreaker (which intentionally differs in form between the map and
// types.Object representations).
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
	return ""
}

// stableObjectPrimaryKey is the types.Object counterpart of
// stableItemPrimaryKey: returns just the primary identity portion (no
// tiebreaker) so it can match items keyed via stableItemPrimaryKey.
func stableObjectPrimaryKey(obj types.Object, fields []BodyFieldDef) string {
	attrs := obj.Attributes()
	declared := map[string]bool{}
	for _, f := range fields {
		declared[bodyFieldKey(f)] = true
	}
	for _, candidate := range stableIdentityFieldOrder {
		if len(fields) > 0 && !declared[candidate] {
			continue
		}
		v, ok := attrs[candidate]
		if !ok {
			continue
		}
		if s, ok := stringifyAttrIdentity(v); ok && s != "" {
			return candidate + "=" + s
		}
	}
	return ""
}

// alignResponseItemsToReference reorders a JSON-decoded array of
// nested-object items in place so that items appear in the same order as
// the reference list (typically the prior plan/state value). Items that
// match a reference element by stable identity key are emitted in the
// reference order; items in the API response that have no matching
// reference element are appended at the end in their original API order
// (which is then itself stabilised by sortResponseItems for determinism).
//
// This is the linchpin that makes nested-object lists order-insensitive
// without ever mutating the planned/config value: the operator's order
// from configuration (carried through to the plan) is the order we save
// to state, regardless of the order Bitbucket happens to return items in.
//
// Reports true iff every reference element matched a response item — a
// signal that the response and reference differ only by ordering. When
// false, the caller may still want to apply the canonical sort to the
// trailing leftovers for stability.
func alignResponseItemsToReference(arr []any, ref []types.Object, fields []BodyFieldDef) bool {
	if len(arr) < 2 || len(ref) == 0 {
		return false
	}
	// Index response items by their primary identity key; tolerate
	// duplicates by maintaining a slice of indices per key (we consume
	// them in order). Items with no resolvable identity (empty key) are
	// skipped here and emitted at the end as leftovers.
	type idxQueue struct{ items []int }
	byKey := map[string]*idxQueue{}
	for i, item := range arr {
		m, ok := item.(map[string]any)
		if !ok {
			return false
		}
		k := stableItemPrimaryKey(m, fields)
		if k == "" {
			continue
		}
		q, ok := byKey[k]
		if !ok {
			q = &idxQueue{}
			byKey[k] = q
		}
		q.items = append(q.items, i)
	}
	out := make([]any, 0, len(arr))
	used := make([]bool, len(arr))
	allMatched := true
	for _, refObj := range ref {
		k := stableObjectPrimaryKey(refObj, fields)
		if k == "" {
			allMatched = false
			continue
		}
		q, ok := byKey[k]
		if !ok || len(q.items) == 0 {
			allMatched = false
			continue
		}
		idx := q.items[0]
		q.items = q.items[1:]
		out = append(out, arr[idx])
		used[idx] = true
	}
	// Append any leftover response items in their original API order.
	leftover := false
	for i, item := range arr {
		if !used[i] {
			out = append(out, item)
			leftover = true
		}
	}
	copy(arr, out)
	if leftover {
		return false
	}
	return allMatched
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
