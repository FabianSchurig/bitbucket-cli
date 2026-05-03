package tfprovider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// setLikeListType is a CustomType for ListNestedAttribute over object items
// where positional order is not semantically meaningful (e.g.
// branch_restrictions.users). It produces setLikeListValue values whose
// ListSemanticEquals treats two lists as equal iff they contain the same
// items, identified by the same stable identity-field precedence used by
// stableItemPrimaryKey / stableObjectPrimaryKey.
//
// Why a CustomType (not a plan modifier or response-side sort alone):
//
//   - Plan modifiers cannot mutate Required attributes (Terraform rejects with
//     "planned value does not match config value").
//   - Aligning response order to prior state alone makes Step "re-plan with
//     same config" pass, but Step "re-plan with reordered config" still fails:
//     refreshed state stays in prior order, then config (in new order) differs
//     from state and the framework reports an in-place update.
//   - With ListSemanticEquals the framework substitutes the prior state when
//     the planned value is set-equivalent, producing an empty diff for any
//     reordering. This works for plan, apply, and arbitrarily many subsequent
//     runs.
type setLikeListType struct {
	basetypes.ListType
	itemFields []BodyFieldDef
}

// setLikeListTypeFor constructs a setLikeListType whose underlying element
// type matches a ListNestedAttribute over the given item fields.
func setLikeListTypeFor(itemFields []BodyFieldDef) setLikeListType {
	return setLikeListType{
		ListType:   basetypes.ListType{ElemType: types.ObjectType{AttrTypes: itemAttrTypes(itemFields)}},
		itemFields: itemFields,
	}
}

// Equal reports type equality. Two setLikeListType values are equal iff
// their underlying element types match.
func (t setLikeListType) Equal(o attr.Type) bool {
	other, ok := o.(setLikeListType)
	if !ok {
		return false
	}
	return t.ListType.Equal(other.ListType)
}

// String returns a human-readable type name for diagnostics.
func (t setLikeListType) String() string {
	return "setLikeListType[" + t.ElemType.String() + "]"
}

// ValueType returns a zero-value of the matching value type so the framework
// knows what to instantiate.
func (t setLikeListType) ValueType(ctx context.Context) attr.Value {
	return setLikeListValue{ListValue: basetypes.NewListNull(t.ElemType), itemFields: t.itemFields}
}

// ValueFromList wraps a base ListValue in a setLikeListValue carrying the
// item-field schema needed for semantic equality.
func (t setLikeListType) ValueFromList(_ context.Context, list basetypes.ListValue) (basetypes.ListValuable, diag.Diagnostics) {
	return setLikeListValue{ListValue: list, itemFields: t.itemFields}, nil
}

// ValueFromTerraform decodes a tftypes.Value via the embedded ListType and
// then wraps it.
func (t setLikeListType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	v, err := t.ListType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}
	lv, ok := v.(basetypes.ListValue)
	if !ok {
		return nil, fmt.Errorf("setLikeListType: expected basetypes.ListValue from underlying decoder, got %T", v)
	}
	return setLikeListValue{ListValue: lv, itemFields: t.itemFields}, nil
}

// setLikeListValue is a basetypes.ListValuable whose ListSemanticEquals is
// order-insensitive for nested-object lists. Apart from semantic equality
// it behaves identically to basetypes.ListValue.
type setLikeListValue struct {
	basetypes.ListValue
	itemFields []BodyFieldDef
}

// Type returns the matching setLikeListType so framework operations pick
// the custom semantic-equality logic.
func (v setLikeListValue) Type(ctx context.Context) attr.Type {
	return setLikeListType{ListType: basetypes.ListType{ElemType: v.ElementType(ctx)}, itemFields: v.itemFields}
}

// Equal returns positional equality, matching basetypes.ListValue. Semantic
// (set-style) equality lives in ListSemanticEquals so Terraform can still
// detect genuine value changes via Equal in code paths that don't invoke
// semantic equality.
func (v setLikeListValue) Equal(o attr.Value) bool {
	other, ok := o.(setLikeListValue)
	if !ok {
		return false
	}
	return v.ListValue.Equal(other.ListValue)
}

// ListSemanticEquals reports whether two nested-object lists contain the
// same elements regardless of order. Items are paired by their stable
// identity key (uuid > id > slug > full_slug > name) and then compared
// with an *asymmetric* equality that treats Unknown attributes on the
// `other` (proposed-new) side as wildcards matching any concrete value
// on the `v` (prior) side — but NOT vice-versa.
//
// Why the asymmetry — and why both directions matter:
//
// terraform-plugin-framework calls SchemaSemanticEquality from two very
// different code paths, both passing PriorData as the receiver `v` and
// ProposedNewData as `other`, and in both, returning true causes the
// framework to REPLACE NewState with PriorData (see
// internal/fwschemadata/value_semantic_equality_list.go).
//
//  1. Plan / Read: PriorData = current state (concrete), ProposedNewData
//     = config-merged plan or API refresh response. The plan side often
//     carries Unknown for every Optional+Computed nested field (e.g.
//     users[*].created_on, display_name) — and the user may have
//     reordered the list in config. We WANT SemanticEqual=true here so
//     the framework substitutes the concrete prior state and produces
//     an empty diff. → tolerate Unknown on `other`.
//
//  2. Create / Update apply consistency: PriorData = req.PlannedState
//     (with Unknowns for computed fields), ProposedNewData = resp.NewState
//     (concrete API response). We MUST return false here, otherwise the
//     framework overwrites our concrete API response with the plan's
//     Unknowns and Terraform Core rejects the apply with "Provider
//     returned invalid result object after apply". → do NOT tolerate
//     Unknown on `v`.
//
// The asymmetric rule "Unknown on `other` is a wildcard, Unknown on `v`
// is a hard mismatch" satisfies both cases with one implementation.
//
// Order-insensitivity is also reinforced by refresh-time alignment in
// buildListFromResponse → alignResponseItemsToReference, which reorders
// the API response to match the prior list's primary-key order before
// SetAttribute. SemanticEquals here is the safety net for cases the
// alignment cannot cover (prior list null/empty, fresh import, etc.).
//
// Items whose identity key is empty fall back to a positional asymmetric
// search since we have no other way to pair them.
func (v setLikeListValue) ListSemanticEquals(ctx context.Context, other basetypes.ListValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	otherTyped, ok := other.(setLikeListValue)
	if !ok {
		// Fall back to positional equality when the peer isn't a
		// setLikeListValue (shouldn't happen at runtime — the schema pins
		// the CustomType — but defensive against framework changes).
		conv, convDiags := other.ToListValue(ctx)
		diags.Append(convDiags...)
		if diags.HasError() {
			return false, diags
		}
		return v.ListValue.Equal(conv), diags
	}

	// Trivially compare null/unknown — semantic equality is only
	// meaningful for known, non-null lists with the same length.
	if v.IsNull() != otherTyped.IsNull() || v.IsUnknown() != otherTyped.IsUnknown() {
		return false, diags
	}
	if v.IsNull() || v.IsUnknown() {
		return true, diags
	}
	left := v.Elements()
	right := otherTyped.Elements()
	if len(left) != len(right) {
		return false, diags
	}

	// Pair items by stable primary identity key. We use a per-key queue so
	// duplicate identities are still matched 1:1.
	type idxQueue struct{ items []int }
	rightByKey := map[string]*idxQueue{}
	rightFallback := []int{}
	for i, e := range right {
		obj, ok := e.(types.Object)
		if !ok || obj.IsNull() || obj.IsUnknown() {
			rightFallback = append(rightFallback, i)
			continue
		}
		k := stableObjectPrimaryKey(obj, v.itemFields)
		if k == "" {
			rightFallback = append(rightFallback, i)
			continue
		}
		q, ok := rightByKey[k]
		if !ok {
			q = &idxQueue{}
			rightByKey[k] = q
		}
		q.items = append(q.items, i)
	}

	usedRight := make([]bool, len(right))
	for _, l := range left {
		lObj, ok := l.(types.Object)
		if !ok {
			return false, diags
		}
		k := ""
		if !lObj.IsNull() && !lObj.IsUnknown() {
			k = stableObjectPrimaryKey(lObj, v.itemFields)
		}
		matched := false
		if k != "" {
			if q, ok := rightByKey[k]; ok {
				for len(q.items) > 0 {
					idx := q.items[0]
					q.items = q.items[1:]
					if usedRight[idx] {
						continue
					}
					rObj, rok := right[idx].(types.Object)
					if !rok {
						continue
					}
					// Pair-by-key + asymmetric Equal: Unknown on the
					// `right` (proposed-new) side is a wildcard for
					// any concrete left value; Unknown on the left
					// (prior) side is NOT — see the method comment
					// for the apply-consistency reasoning.
					if asymmetricObjectEqual(lObj, rObj) {
						usedRight[idx] = true
						matched = true
						break
					}
				}
			}
		}
		if matched {
			continue
		}
		// Fall back to a linear search over still-unused right items so
		// items without a primary key (or with an unknown nested attr that
		// blocks Equal) still get a chance to pair up. Same asymmetric
		// rule applies (Unknown on right side is wildcard).
		for _, idx := range rightFallback {
			if usedRight[idx] {
				continue
			}
			rObj, rok := right[idx].(types.Object)
			if !rok {
				if l.Equal(right[idx]) {
					usedRight[idx] = true
					matched = true
					break
				}
				continue
			}
			lObjFB, lok := l.(types.Object)
			if !lok {
				continue
			}
			if asymmetricObjectEqual(lObjFB, rObj) {
				usedRight[idx] = true
				matched = true
				break
			}
		}
		if !matched {
			return false, diags
		}
	}
	return true, diags
}

// asymmetricObjectEqual reports whether two nested-object values are
// equivalent under the asymmetric rule used by ListSemanticEquals:
// Unknown attributes on the `proposed` side match any concrete value on
// the `prior` side, but Unknown attributes on the `prior` side do NOT
// match concrete values on the `proposed` side.
//
// See ListSemanticEquals for the full rationale (Plan vs. apply
// consistency). This helper does the recursive walk: when an attribute
// is itself an Object it recurses, so deeply-nested computed fields
// (e.g. links.avatar.href) are handled correctly. Lists/Maps/Sets fall
// back to strict Equal — they are not the targeted shape here, and
// nested ListNestedAttributes carry their own semantic-equality.
func asymmetricObjectEqual(prior, proposed types.Object) bool {
	if prior.IsNull() != proposed.IsNull() {
		return false
	}
	if prior.IsNull() {
		return true
	}
	// Whole-object Unknown on proposed is a wildcard for anything on
	// prior (including a concrete or also-Unknown prior). The reverse
	// case — Unknown on prior, concrete on proposed — is rejected just
	// below.
	if proposed.IsUnknown() {
		return true
	}
	if prior.IsUnknown() {
		return false
	}
	pAttrs := prior.Attributes()
	nAttrs := proposed.Attributes()
	if len(pAttrs) != len(nAttrs) {
		return false
	}
	for name, pv := range pAttrs {
		nv, ok := nAttrs[name]
		if !ok {
			return false
		}
		// Unknown on proposed side → wildcard match.
		if nv.IsUnknown() {
			continue
		}
		// Unknown on prior side but proposed is concrete → mismatch.
		if pv.IsUnknown() {
			return false
		}
		// Recurse into nested Objects so deeper computed fields enjoy
		// the same asymmetric tolerance.
		if pObj, ok := pv.(types.Object); ok {
			if nObj, ok := nv.(types.Object); ok {
				if !asymmetricObjectEqual(pObj, nObj) {
					return false
				}
				continue
			}
			return false
		}
		if !pv.Equal(nv) {
			return false
		}
	}
	return true
}

// fields, used when a nested-object array attribute is absent from a
// response or input.
func setLikeListNull(itemFields []BodyFieldDef) setLikeListValue {
	objType := types.ObjectType{AttrTypes: itemAttrTypes(itemFields)}
	return setLikeListValue{ListValue: basetypes.NewListNull(objType), itemFields: itemFields}
}

// setLikeListValueMust wraps types.ListValueMust into a setLikeListValue so
// the value's concrete type matches the schema's CustomType.
func setLikeListValueMust(itemFields []BodyFieldDef, elements []attr.Value) setLikeListValue {
	objType := types.ObjectType{AttrTypes: itemAttrTypes(itemFields)}
	return setLikeListValue{ListValue: types.ListValueMust(objType, elements), itemFields: itemFields}
}

// Compile-time assertions: the value implements semantic-equality and the
// type implements ListTypable.
var (
	_ basetypes.ListTypable                    = setLikeListType{}
	_ basetypes.ListValuableWithSemanticEquals = setLikeListValue{}
)
