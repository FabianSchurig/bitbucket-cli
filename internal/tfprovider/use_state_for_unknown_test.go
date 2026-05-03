package tfprovider

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const useStateForUnknownDesc = useStateForUnknownDescription

// TestSetLikeListUseStateModifier_UnknownPlanFallsBackToPriorState guards
// the perpetual-diff regression on bitbucket_branch_restrictions: when the
// raw proposed plan differs from prior state (e.g. because users were
// reordered in config), the framework's MarkComputedNilsAsUnknown runs
// BEFORE plan modifiers and marks every Optional+Computed attribute whose
// config is null (e.g. `groups`) as Unknown. The plan modifier must
// substitute prior state in that case so the planned value matches the
// concrete prior state and the framework reports no change.
func TestSetLikeListUseStateModifier_UnknownPlanFallsBackToPriorState(t *testing.T) {
	itemFields := []BodyFieldDef{{Path: "uuid", Type: "string"}}
	objType := types.ObjectType{AttrTypes: itemAttrTypes(itemFields)}

	state := types.ListValueMust(objType, []attr.Value{})
	plan := basetypes.NewListUnknown(objType)

	req := planmodifier.ListRequest{
		Path:       path.Root("groups"),
		StateValue: state,
		PlanValue:  plan,
	}
	resp := &planmodifier.ListResponse{PlanValue: plan}

	mod := setLikeListUseStateIfSetEqual(itemFields).(setLikeListUseStateModifier)
	mod.PlanModifyList(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %#v", resp.Diagnostics)
	}
	if !resp.PlanValue.Equal(state) {
		t.Fatalf("plan modifier must substitute prior state when plan is Unknown; got %#v want %#v", resp.PlanValue, state)
	}
}

// TestComputedAttrsAttachUseStateForUnknown guards that all Computed-only
// resource attributes (id, api_response, scalar response fields, primitive
// list response fields) attach a UseStateForUnknown plan modifier. Without
// it, MarkComputedNilsAsUnknown causes those attrs to flip from a known
// prior state value to "(known after apply)" whenever any other attribute's
// raw plan differs from prior state — producing perpetual diffs on
// reorder-only configuration changes for set-like nested lists.
func TestComputedAttrsAttachUseStateForUnknown(t *testing.T) {
	base := resourceBaseAttrs()

	idAttr, ok := base["id"].(resourceschema.StringAttribute)
	if !ok {
		t.Fatalf("resourceBaseAttrs id must be StringAttribute; got %T", base["id"])
	}
	if !hasStringPlanModifierWithDesc(idAttr.PlanModifiers, useStateForUnknownDesc) {
		t.Fatalf("id attribute must attach a UseStateForUnknown plan modifier; got %#v", idAttr.PlanModifiers)
	}

	apiAttr, ok := base["api_response"].(resourceschema.StringAttribute)
	if !ok {
		t.Fatalf("resourceBaseAttrs api_response must be StringAttribute; got %T", base["api_response"])
	}
	if !hasStringPlanModifierWithDesc(apiAttr.PlanModifiers, useStateForUnknownDesc) {
		t.Fatalf("api_response attribute must attach a UseStateForUnknown plan modifier; got %#v", apiAttr.PlanModifiers)
	}

	strRF := BodyFieldDef{Path: "branch_type", Type: "string"}
	strAttr, ok := responseFieldAttr(strRF).(resourceschema.StringAttribute)
	if !ok {
		t.Fatalf("string response field must produce StringAttribute; got %T", responseFieldAttr(strRF))
	}
	if !hasStringPlanModifierWithDesc(strAttr.PlanModifiers, useStateForUnknownDesc) {
		t.Fatalf("Computed-only string response field must attach UseStateForUnknown; got %#v", strAttr.PlanModifiers)
	}

	intRF := BodyFieldDef{Path: "value", Type: "int"}
	intAttr, ok := responseFieldAttr(intRF).(resourceschema.Int64Attribute)
	if !ok {
		t.Fatalf("int response field must produce Int64Attribute; got %T", responseFieldAttr(intRF))
	}
	if !hasInt64PlanModifierWithDesc(intAttr.PlanModifiers, useStateForUnknownDesc) {
		t.Fatalf("Computed-only int response field must attach UseStateForUnknown; got %#v", intAttr.PlanModifiers)
	}

	listRF := BodyFieldDef{Path: "tags", IsArray: true}
	listAttr, ok := responseFieldAttr(listRF).(resourceschema.ListAttribute)
	if !ok {
		t.Fatalf("primitive list response field must produce ListAttribute; got %T", responseFieldAttr(listRF))
	}
	if !hasListPlanModifierWithDesc(listAttr.PlanModifiers, useStateForUnknownDesc) {
		t.Fatalf("Computed-only primitive list response field must attach UseStateForUnknown; got %#v", listAttr.PlanModifiers)
	}
}

func hasStringPlanModifierWithDesc(mods []planmodifier.String, desc string) bool {
	for _, m := range mods {
		if m != nil && strings.Contains(m.Description(context.Background()), desc) {
			return true
		}
	}
	return false
}

func hasInt64PlanModifierWithDesc(mods []planmodifier.Int64, desc string) bool {
	for _, m := range mods {
		if m != nil && strings.Contains(m.Description(context.Background()), desc) {
			return true
		}
	}
	return false
}

func hasListPlanModifierWithDesc(mods []planmodifier.List, desc string) bool {
	for _, m := range mods {
		if m != nil && strings.Contains(m.Description(context.Background()), desc) {
			return true
		}
	}
	return false
}
