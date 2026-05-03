package tfprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure GenericResource implements the optional ModifyPlan hook.
var _ resource.ResourceWithModifyPlan = &GenericResource{}

// ModifyPlan provides a uniform "no semantic change ⇒ keep prior state"
// guarantee for every resource generated from the Bitbucket OpenAPI spec.
//
// Why a resource-level hook rather than per-attribute UseStateForUnknown:
//
// terraform-plugin-framework runs a step called MarkComputedNilsAsUnknown
// before plan modifiers. Whenever the raw plan differs from prior state at
// the schema level, every Optional+Computed attribute whose configuration
// is null is flipped to Unknown. For nested-object lists we declare with
// the order-insensitive setLikeListType, this fires whenever an operator
// merely reorders items in their configuration — even though the change
// is semantically a no-op. The schema-level setLikeListUseStateIfSetEqual
// modifier restores the list itself, but the framework has no way to
// "undo" the Unknown marks on sibling Computed attributes from there.
//
// Sprinkling stringplanmodifier.UseStateForUnknown across every Computed
// attribute is *not* a correct fix: it forces the planned post-apply
// value of fields like updated_on, name and api_response to equal prior
// state, which produces "Provider produced inconsistent result after
// apply" errors on legitimate Updates that actually change those fields.
//
// The right level for this decision is the whole resource: if the user's
// configuration is semantically equal to the prior state — using set-like
// semantic equality for set-like list attributes and tftypes equality
// elsewhere — there is no change to plan, and we can safely substitute
// the prior state into the plan. When any configured attribute really
// changes, we leave the plan untouched and the framework's normal
// Computed-Unknown promotion + per-attribute refresh proceeds.
//
// This is general for every list-like (and non-list) endpoint exposed by
// the generator: it inspects the current resource's schema and body
// definitions rather than special-casing any particular resource.
func (r *GenericResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Skip create (no prior state) and destroy (plan is null).
	if req.State.Raw.IsNull() || !req.State.Raw.IsKnown() {
		return
	}
	if req.Plan.Raw.IsNull() {
		return
	}

	configMap, err := topLevelObjectAsMap(req.Config.Raw)
	if err != nil {
		return
	}
	stateMap, err := topLevelObjectAsMap(req.State.Raw)
	if err != nil {
		return
	}

	for name, attr := range r.schemaAttrs {
		if !isConfigurableAttr(attr) {
			continue
		}
		cfg, cfgOk := configMap[name]
		st, stOk := stateMap[name]
		if !cfgOk || !stOk {
			// Attribute missing from one side — be safe and let the
			// framework's normal planning take over.
			return
		}
		equal, ok := attrValuesSemanticallyEqual(ctx, cfg, st, attr)
		if !ok {
			// Could not compare — fall back to default planning.
			return
		}
		if !equal {
			return
		}
	}

	// Every configurable attribute is semantically unchanged. Preserve the
	// prior state as the plan so Computed attributes stay concrete and we
	// produce an empty diff.
	resp.Plan.Raw = req.State.Raw
}

// topLevelObjectAsMap decodes a tftypes object value (e.g. the Raw value
// of a resource's Config or State) into a map keyed by attribute name.
func topLevelObjectAsMap(v tftypes.Value) (map[string]tftypes.Value, error) {
	m := map[string]tftypes.Value{}
	if err := v.As(&m); err != nil {
		return nil, err
	}
	return m, nil
}

// isConfigurableAttr reports whether an attribute is part of the user's
// configuration surface — i.e. Required or Optional. Computed-only
// attributes are never written by the user and are skipped: they have a
// null raw config value by definition and would always look "different"
// to a raw equality check.
func isConfigurableAttr(a schema.Attribute) bool {
	return a.IsRequired() || a.IsOptional()
}

// attrValuesSemanticallyEqual compares two raw tftypes values for the
// given schema attribute. For list-nested attributes whose CustomType is
// setLikeListType, it uses the type's order-insensitive ListSemanticEquals;
// for everything else it uses tftypes.Value.Equal. Returns (equal, ok)
// where ok=false signals a comparison error (caller should treat as a
// "real change" to be safe).
func attrValuesSemanticallyEqual(ctx context.Context, cfg, st tftypes.Value, a schema.Attribute) (bool, bool) {
	if itemFields, ok := setLikeItemFieldsFor(a); ok {
		return setLikeRawValuesEqual(ctx, cfg, st, itemFields)
	}
	return cfg.Equal(st), true
}

// setLikeItemFieldsFor returns the BodyFieldDef slice describing the
// items of a set-like list attribute, or (nil, false) if the attribute
// is not a set-like list.
func setLikeItemFieldsFor(a schema.Attribute) ([]BodyFieldDef, bool) {
	lna, ok := a.(schema.ListNestedAttribute)
	if !ok {
		return nil, false
	}
	slt, ok := lna.CustomType.(setLikeListType)
	if !ok {
		return nil, false
	}
	return slt.itemFields, true
}

// setLikeRawValuesEqual decodes both raw values via setLikeListType and
// reports whether the resulting setLikeListValues are semantically equal
// (same items by stable identity key, regardless of order).
func setLikeRawValuesEqual(ctx context.Context, cfg, st tftypes.Value, itemFields []BodyFieldDef) (bool, bool) {
	listType := setLikeListTypeFor(itemFields)
	cfgVal, err := listType.ValueFromTerraform(ctx, cfg)
	if err != nil {
		return false, false
	}
	stVal, err := listType.ValueFromTerraform(ctx, st)
	if err != nil {
		return false, false
	}
	cfgSet, ok := cfgVal.(setLikeListValue)
	if !ok {
		return false, false
	}
	stSet, ok := stVal.(setLikeListValue)
	if !ok {
		return false, false
	}
	equal, diags := cfgSet.ListSemanticEquals(ctx, stSet)
	if diags.HasError() {
		return false, false
	}
	return equal, true
}
