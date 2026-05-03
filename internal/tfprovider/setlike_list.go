package tfprovider

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// setLikeListType is a custom ListType for nested-object arrays whose
// upstream Bitbucket API treats the collection as unordered. Wiring it via
// the schema's `CustomType` field makes the framework call
// `ListSemanticEquals` at every value-comparison site (config↔plan,
// plan↔refresh, prior↔new state). Two values containing the same multiset
// of elements are reported as equal regardless of order, so:
//
//   - The post-apply consistency check passes when the API returns elements
//     in a different order than they were submitted (the original
//     "Provider produced inconsistent result after apply" failure).
//   - The plan-time validity check passes when the operator wrote the
//     elements in a different order than the API returns them (the v0.15.6
//     "Provider produced invalid plan" regression).
//   - No silent reordering of the operator's HCL happens at plan time, so
//     adding or removing an element produces a surgical, comprehensible
//     diff in the operator's chosen order (the "perpetual diff on add"
//     UX failure that survived v0.15.6).
//
// This is the framework-native analog of the SDKv2 `schema.TypeSet` pattern
// the upstream issue suggests — it gives the same order-insensitive
// behaviour without changing the schema attribute kind, without forcing
// every nested attribute to be Computed-only, and without losing the
// ordered-array semantics for genuinely ordered scalar lists like `tags`.
type setLikeListType struct {
	basetypes.ListType

	// itemFields is the schema of the array's element objects. It is
	// captured here so `ListSemanticEquals` can find the canonical identity
	// field (uuid > id > slug > ...) and produce a stable per-element key
	// without re-deriving the schema from each `attr.Value`.
	itemFields []BodyFieldDef
}

// setLikeListTypeFor builds the custom list type for a nested-object array
// whose element schema is `itemFields`. The wrapped `basetypes.ListType`
// uses the same element type (`types.ObjectType{...}`) the schema would
// otherwise infer, so existing state files continue to load.
func setLikeListTypeFor(itemFields []BodyFieldDef) setLikeListType {
	return setLikeListType{
		ListType:   basetypes.ListType{ElemType: types.ObjectType{AttrTypes: itemAttrTypes(itemFields)}},
		itemFields: itemFields,
	}
}

// Equal honours the ListType element-type comparison and additionally
// distinguishes setLikeListType instances by the *full shape* of their
// nested-object schema. Two custom types whose paths happen to align but
// whose Type/IsArray/IsObject/Required/ItemFields differ must NOT compare
// equal — otherwise the framework's type-equality checks would let a value
// shaped for one schema slip into a slot expecting another, producing
// confusing downstream conversion errors.
func (t setLikeListType) Equal(o attr.Type) bool {
	other, ok := o.(setLikeListType)
	if !ok {
		return false
	}
	if !t.ListType.Equal(other.ListType) {
		return false
	}
	return bodyFieldsEqual(t.itemFields, other.itemFields)
}

// bodyFieldsEqual returns true when two slices of BodyFieldDef describe the
// same nested-object shape. Compared recursively over ItemFields so nested
// arrays/objects are not collapsed to a path-only comparison.
func bodyFieldsEqual(a, b []BodyFieldDef) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		x, y := a[i], b[i]
		if x.Path != y.Path || x.Type != y.Type ||
			x.IsArray != y.IsArray || x.IsObject != y.IsObject ||
			x.Required != y.Required {
			return false
		}
		if !bodyFieldsEqual(x.ItemFields, y.ItemFields) {
			return false
		}
	}
	return true
}

// String returns a human-readable identifier used in framework error
// messages. Including the field paths makes diagnostics actionable when the
// type comparison fails.
func (t setLikeListType) String() string {
	paths := make([]string, len(t.itemFields))
	for i, f := range t.itemFields {
		paths[i] = f.Path
	}
	return fmt.Sprintf("setLikeListType[%s]", joinPaths(paths))
}

// ValueFromList wraps a plain ListValue in our custom value type so the
// framework will call `ListSemanticEquals` on it. Required by
// `basetypes.ListTypable`.
func (t setLikeListType) ValueFromList(_ context.Context, in basetypes.ListValue) (basetypes.ListValuable, diag.Diagnostics) {
	return setLikeListValue{ListValue: in, itemFields: t.itemFields}, nil
}

// ValueFromTerraform is the deserialization entry point used by the plugin
// protocol. Without overriding it the framework would hand back a plain
// `basetypes.ListValue` during state refresh and the semantic-equality
// logic would never run.
func (t setLikeListType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	v, err := t.ListType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}
	listVal, ok := v.(basetypes.ListValue)
	if !ok {
		return nil, fmt.Errorf("setLikeListType: expected basetypes.ListValue, got %T", v)
	}
	return setLikeListValue{ListValue: listVal, itemFields: t.itemFields}, nil
}

// ValueType returns the zero value of the custom value type. The framework
// uses it to detect whether an `attr.Value` is of the expected type.
func (t setLikeListType) ValueType(_ context.Context) attr.Value {
	return setLikeListValue{itemFields: t.itemFields}
}

// setLikeListValue is the value half of the custom type. It embeds
// `basetypes.ListValue` so all standard list operations (Length, Elements,
// ToTerraformValue, ...) keep working unchanged, and adds
// `ListSemanticEquals` to make the framework treat the value as
// order-insensitive at every comparison point.
type setLikeListValue struct {
	basetypes.ListValue

	itemFields []BodyFieldDef
}

// Equal must distinguish setLikeListValue instances from plain ListValues so
// the framework's per-attribute equality probes route through this method
// (and from there to ListSemanticEquals).
func (v setLikeListValue) Equal(o attr.Value) bool {
	other, ok := o.(setLikeListValue)
	if !ok {
		return false
	}
	return v.ListValue.Equal(other.ListValue)
}

// Type returns the corresponding setLikeListType, completing the type↔value
// pairing required by the framework.
func (v setLikeListValue) Type(_ context.Context) attr.Type {
	return setLikeListType{
		ListType:   basetypes.ListType{ElemType: v.ElementType(context.Background())},
		itemFields: v.itemFields,
	}
}

// ToListValue converts back to a plain ListValue. Required by
// `basetypes.ListValuable`.
func (v setLikeListValue) ToListValue(_ context.Context) (basetypes.ListValue, diag.Diagnostics) {
	return v.ListValue, nil
}

// ListSemanticEquals returns true when v and other contain the same
// multiset of elements, irrespective of element order. Null/unknown values
// are never considered semantically equal to anything else — only known
// values participate in this comparison (per the framework contract on
// ListValuableWithSemanticEquals).
//
// For ordinary API objects with unique identity fields (uuid, id, slug, ...),
// equality is based on those identity keys only. This intentionally ignores
// API-computed fields that are Unknown in config/plan but populated in state.
// Ambiguous lists (missing identity keys or duplicate identities) fall back to
// the full canonical sort key so genuinely different duplicate payloads are
// still distinguished.
func (v setLikeListValue) ListSemanticEquals(_ context.Context, other basetypes.ListValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	o, ok := other.(setLikeListValue)
	if !ok {
		return false, diags
	}
	if v.IsNull() != o.IsNull() || v.IsUnknown() != o.IsUnknown() {
		return false, diags
	}
	if v.IsNull() || v.IsUnknown() {
		// Both sides are null or both are unknown — but the framework only
		// invokes semantic equality on known values, so this is a defensive
		// fallback.
		return v.IsNull() == o.IsNull() && v.IsUnknown() == o.IsUnknown(), diags
	}

	left := v.Elements()
	right := o.Elements()
	if len(left) != len(right) {
		return false, diags
	}

	leftKeys, leftPrimary := elementPrimaryKeys(left, v.itemFields)
	rightKeys, rightPrimary := elementPrimaryKeys(right, o.itemFields)
	if !leftPrimary || !rightPrimary {
		leftKeys = elementSortKeys(left, v.itemFields)
		rightKeys = elementSortKeys(right, o.itemFields)
	}
	sort.Strings(leftKeys)
	sort.Strings(rightKeys)
	for i := range leftKeys {
		if leftKeys[i] != rightKeys[i] {
			return false, diags
		}
	}
	return true, diags
}

// elementPrimaryKeys returns the identity-only key sequence for object
// elements when every element has a unique stable identity field (uuid, id,
// slug, ...). Computed fields such as display_name are intentionally excluded:
// config/plan often has them Unknown while Read refresh has them populated.
// If any element lacks a primary identity or a duplicate identity appears, the
// caller falls back to full canonical element keys so genuinely ambiguous
// lists are not collapsed.
func elementPrimaryKeys(elements []attr.Value, itemFields []BodyFieldDef) ([]string, bool) {
	keys := make([]string, len(elements))
	seen := make(map[string]bool, len(elements))
	for i, e := range elements {
		obj, ok := e.(types.Object)
		if !ok {
			return nil, false
		}
		key := stableObjectPrimaryKey(obj, itemFields)
		if isFallbackPrimaryKey(key) || seen[key] {
			return nil, false
		}
		keys[i] = key
		seen[key] = true
	}
	return keys, true
}

// isFallbackPrimaryKey reports whether stableObjectPrimaryKey had to fall
// back because no stable identity field (uuid, id, slug, ...) was available.
// Callers use this signal to avoid identity-only comparison and instead fall
// back to full canonical element comparison for ambiguous lists.
func isFallbackPrimaryKey(key string) bool {
	return strings.HasPrefix(key, fallbackObjectKeyPrefix) ||
		strings.HasPrefix(key, fallbackJSONKeyPrefix) ||
		strings.HasPrefix(key, fallbackRawKeyPrefix)
}

// elementSortKeys produces the canonical per-element key sequence used by
// ListSemanticEquals. Non-object elements are encoded via their default
// String() form so the comparison degrades gracefully if the element type
// is ever expanded.
func elementSortKeys(elements []attr.Value, itemFields []BodyFieldDef) []string {
	keys := make([]string, len(elements))
	for i, e := range elements {
		if obj, ok := e.(types.Object); ok {
			keys[i] = stableObjectSortKey(obj, itemFields)
			continue
		}
		keys[i] = "raw=" + e.String()
	}
	return keys
}

// joinPaths joins a slice of field paths with a comma. Tiny helper kept
// local so the type's String() doesn't pull in the full strings package
// surface.
func joinPaths(paths []string) string {
	out := ""
	for i, p := range paths {
		if i > 0 {
			out += ","
		}
		out += p
	}
	return out
}

// Compile-time interface assertions: any breakage of these contracts (e.g.
// a framework upgrade renaming a method) surfaces as a build error here
// rather than as a runtime panic during the first plan.
var (
	_ basetypes.ListTypable                    = setLikeListType{}
	_ basetypes.ListValuable                   = setLikeListValue{}
	_ basetypes.ListValuableWithSemanticEquals = setLikeListValue{}
)
