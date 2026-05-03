package tfprovider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestBuildListFromResponseWithoutPriorSortsByStableKey covers the
// data-source / no-prior-state path: when no reference order is available,
// nested-object items must be sorted by a stable identity key so two
// equivalent API responses (same elements in different order) yield
// byte-identical Terraform values.
func TestBuildListFromResponseWithoutPriorSortsByStableKey(t *testing.T) {
	userFields := []BodyFieldDef{
		{Path: "uuid", Type: "string"},
		{Path: "display_name", Type: "string"},
	}
	objType := types.ObjectType{AttrTypes: itemAttrTypes(userFields)}

	userA := map[string]any{"uuid": "{aaaa-aaaa}", "display_name": "Alice"}
	userB := map[string]any{"uuid": "{bbbb-bbbb}", "display_name": "Bob"}

	planOrder := buildListFromResponse([]any{userA, userB}, userFields, types.ListNull(objType))
	apiOrder := buildListFromResponse([]any{userB, userA}, userFields, types.ListNull(objType))

	if !planOrder.Equal(apiOrder) {
		t.Fatalf("nested-object list order must be deterministic regardless of input order:\n  plan-order = %s\n   api-order = %s", planOrder.String(), apiOrder.String())
	}
	first := planOrder.Elements()[0].(types.Object).Attributes()["uuid"].(types.String).ValueString()
	if first != "{aaaa-aaaa}" {
		t.Fatalf("expected sorted-by-uuid first element to be {aaaa-aaaa}, got %s", first)
	}
}

// TestBuildListFromResponseAlignsToPriorOrder covers the planning path:
// when a prior plan/state value is supplied, the API response must be
// reordered to match it. This is what keeps state == plan and avoids the
// "planned value does not match config value" framework error on Required
// nested attributes (e.g. branch_restrictions.users[].uuid).
func TestBuildListFromResponseAlignsToPriorOrder(t *testing.T) {
	userFields := []BodyFieldDef{{Path: "uuid", Type: "string"}}
	attrTypes := itemAttrTypes(userFields)
	objType := types.ObjectType{AttrTypes: attrTypes}

	mkObj := func(uuid string) types.Object {
		return types.ObjectValueMust(attrTypes, map[string]attr.Value{
			"uuid": types.StringValue(uuid),
		})
	}

	// Operator wrote users in B,A order; API returns them sorted A,B.
	prior := types.ListValueMust(objType, []attr.Value{
		mkObj("{bbbb-bbbb}"),
		mkObj("{aaaa-aaaa}"),
	})
	apiResp := []any{
		map[string]any{"uuid": "{aaaa-aaaa}"},
		map[string]any{"uuid": "{bbbb-bbbb}"},
	}

	got := buildListFromResponse(apiResp, userFields, prior)
	if got.Elements()[0].(types.Object).Attributes()["uuid"].(types.String).ValueString() != "{bbbb-bbbb}" {
		t.Fatalf("response must be aligned to prior order; got %s", got.String())
	}
	if got.Elements()[1].(types.Object).Attributes()["uuid"].(types.String).ValueString() != "{aaaa-aaaa}" {
		t.Fatalf("response must be aligned to prior order; got %s", got.String())
	}
}

// TestBuildListFromResponseAppendsLeftoverItemsFromAPI covers the case
// where the API returns an item that wasn't in the prior plan/state: the
// matched items keep their planned order and the new item is appended at
// the end so the operator still sees a deterministic diff.
func TestBuildListFromResponseAppendsLeftoverItemsFromAPI(t *testing.T) {
	userFields := []BodyFieldDef{{Path: "uuid", Type: "string"}}
	attrTypes := itemAttrTypes(userFields)
	objType := types.ObjectType{AttrTypes: attrTypes}

	mkObj := func(uuid string) types.Object {
		return types.ObjectValueMust(attrTypes, map[string]attr.Value{
			"uuid": types.StringValue(uuid),
		})
	}

	prior := types.ListValueMust(objType, []attr.Value{mkObj("{bbbb-bbbb}"), mkObj("{aaaa-aaaa}")})
	apiResp := []any{
		map[string]any{"uuid": "{aaaa-aaaa}"},
		map[string]any{"uuid": "{cccc-cccc}"},
		map[string]any{"uuid": "{bbbb-bbbb}"},
	}

	got := buildListFromResponse(apiResp, userFields, prior)
	uuids := make([]string, 0, len(got.Elements()))
	for _, e := range got.Elements() {
		uuids = append(uuids, e.(types.Object).Attributes()["uuid"].(types.String).ValueString())
	}
	if uuids[0] != "{bbbb-bbbb}" || uuids[1] != "{aaaa-aaaa}" || uuids[2] != "{cccc-cccc}" {
		t.Fatalf("matched items must keep prior order and new items append to end; got %v", uuids)
	}
}

// TestBuildListFromResponseTiebreakerForDuplicateIdentity guards the
// total-ordering guarantee of the no-prior fallback: when two items share
// the same identity-field value, the canonical JSON tiebreaker still
// produces a deterministic, reproducible order.
func TestBuildListFromResponseTiebreakerForDuplicateIdentity(t *testing.T) {
	fields := []BodyFieldDef{
		{Path: "uuid", Type: "string"},
		{Path: "display_name", Type: "string"},
	}
	objType := types.ObjectType{AttrTypes: itemAttrTypes(fields)}
	dupA := map[string]any{"uuid": "{same}", "display_name": "Alice"}
	dupB := map[string]any{"uuid": "{same}", "display_name": "Bob"}

	planOrder := buildListFromResponse([]any{dupA, dupB}, fields, types.ListNull(objType))
	apiOrder := buildListFromResponse([]any{dupB, dupA}, fields, types.ListNull(objType))
	if !planOrder.Equal(apiOrder) {
		t.Fatalf("duplicate-identity items must sort deterministically via tiebreaker:\n  plan-order = %s\n   api-order = %s", planOrder.String(), apiOrder.String())
	}
}

// TestNestedObjectArraySchemasAttachSetLikeListPlanModifier guards that the
// schema builders attach the setLikeListUseStateIfSetEqual plan modifier on
// nested-object array attributes. The framework only invokes
// ListSemanticEquals during Create/Update/Read — not during
// PlanResourceChange — so without this plan modifier a `terraform plan`
// against a reordered configuration would show a perpetual diff. The plan
// modifier closes that gap by substituting the prior state when the
// planned and prior lists contain the same items (compared by stable
// identity key).
func TestNestedObjectArraySchemasAttachSetLikeListPlanModifier(t *testing.T) {
	itemFields := []BodyFieldDef{{Path: "uuid", Type: "string"}}

	bodyAttr, ok := bodyFieldAttr(BodyFieldDef{Path: "users", IsArray: true, ItemFields: itemFields}).(resourceschema.ListNestedAttribute)
	if !ok {
		t.Fatalf("bodyFieldAttr returned non-ListNestedAttribute")
	}
	if len(bodyAttr.PlanModifiers) != 1 {
		t.Fatalf("bodyFieldAttr must attach exactly one plan modifier; got %#v", bodyAttr.PlanModifiers)
	}
	if _, ok := bodyAttr.PlanModifiers[0].(setLikeListUseStateModifier); !ok {
		t.Fatalf("bodyFieldAttr plan modifier must be setLikeListUseStateModifier; got %T", bodyAttr.PlanModifiers[0])
	}

	respAttr, ok := responseFieldAttr(BodyFieldDef{Path: "users", IsArray: true, ItemFields: itemFields}).(resourceschema.ListNestedAttribute)
	if !ok {
		t.Fatalf("responseFieldAttr returned non-ListNestedAttribute")
	}
	if len(respAttr.PlanModifiers) != 1 {
		t.Fatalf("responseFieldAttr must attach exactly one plan modifier; got %#v", respAttr.PlanModifiers)
	}
	if _, ok := respAttr.PlanModifiers[0].(setLikeListUseStateModifier); !ok {
		t.Fatalf("responseFieldAttr plan modifier must be setLikeListUseStateModifier; got %T", respAttr.PlanModifiers[0])
	}
}
