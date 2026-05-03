package tfprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

// setLikeListUseStateIfSetEqual returns a List plan modifier that replaces the
// planned value with the prior state value when the two lists contain the
// same items (compared by stable identity key via setLikeListValue's
// ListSemanticEquals). This prevents perpetual diffs when an operator
// reorders items in the configuration of a nested-object list whose order
// is not semantically meaningful (e.g. branch_restrictions.users).
//
// Why a plan modifier in addition to ListSemanticEquals on the type:
//
// terraform-plugin-framework only invokes ValueSemanticEquality during
// Create/Update/Read (server_createresource.go / server_updateresource.go /
// server_readresource.go). It does NOT invoke it during PlanResourceChange.
// As a result, ListSemanticEquals alone is enough to keep the post-apply
// state in the operator's configured order, but it cannot make a subsequent
// `terraform plan` against a reordered config show an empty diff —
// PlanResourceChange compares the proposed-new value (built from the new
// config) with the prior state directly. This plan modifier closes that
// gap by substituting the prior state into the plan when the set of items
// is unchanged.
func setLikeListUseStateIfSetEqual(itemFields []BodyFieldDef) planmodifier.List {
	return setLikeListUseStateModifier{itemFields: itemFields}
}

type setLikeListUseStateModifier struct {
	itemFields []BodyFieldDef
}

func (m setLikeListUseStateModifier) Description(_ context.Context) string {
	return "If the planned list contains the same set of items as the prior state " +
		"(compared by stable identity key, ignoring order), the prior state value " +
		"is preserved to avoid perpetual diffs from configuration reordering."
}

func (m setLikeListUseStateModifier) MarkdownDescription(ctx context.Context) string {
	return m.Description(ctx)
}

func (m setLikeListUseStateModifier) PlanModifyList(ctx context.Context, req planmodifier.ListRequest, resp *planmodifier.ListResponse) {
	// Nothing to do if either side is missing.
	if req.StateValue.IsNull() || req.StateValue.IsUnknown() {
		return
	}
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		return
	}

	// Wrap as setLikeListValue so we can reuse the type's ListSemanticEquals,
	// which performs identity-key based pairing with asymmetric Unknown
	// tolerance. Receiver = planned (proposed); argument = prior state.
	planned := setLikeListValue{ListValue: req.PlanValue, itemFields: m.itemFields}
	prior := setLikeListValue{ListValue: req.StateValue, itemFields: m.itemFields}

	equal, diags := planned.ListSemanticEquals(ctx, prior)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if equal {
		resp.PlanValue = req.StateValue
	}
}
