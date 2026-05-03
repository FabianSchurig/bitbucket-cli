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
	return "setLikeListType[" + t.ListType.ElemType.String() + "]"
}

// ValueType returns a zero-value of the matching value type so the framework
// knows what to instantiate.
func (t setLikeListType) ValueType(ctx context.Context) attr.Value {
	return setLikeListValue{ListValue: basetypes.NewListNull(t.ListType.ElemType), itemFields: t.itemFields}
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
	return setLikeListType{ListType: basetypes.ListType{ElemType: v.ListValue.ElementType(ctx)}, itemFields: v.itemFields}
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
// same elements regardless of order. Items are matched by their stable
// identity key (uuid > id > slug > full_slug > name); two items match when
// their primary keys are equal AND every declared item field is positionally
// equal between them. Items whose identity key is empty fall back to deep
// per-field equality.
//
// This is what makes plan + apply idempotent under reordering: the
// framework substitutes prior state when SemanticEquals returns true,
// producing an empty diff for any permutation.
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
					if l.Equal(right[idx]) {
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
		// blocks Equal) still get a chance to pair up.
		for _, idx := range rightFallback {
			if usedRight[idx] {
				continue
			}
			if l.Equal(right[idx]) {
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

// Compile-time assertions: the value implements semantic-equality and the
// type implements ListTypable.
var (
	_ basetypes.ListTypable                  = setLikeListType{}
	_ basetypes.ListValuableWithSemanticEquals = setLikeListValue{}
)
