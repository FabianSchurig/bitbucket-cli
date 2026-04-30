package tfprovider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// userFields is the canonical example used throughout these tests: the
// `users` nested-object array on `bitbucket_branch_restrictions`. Only
// `uuid` is a real identity field; `display_name` is a computed attribute
// the API returns alongside.
var userFields = []BodyFieldDef{
	{Path: "uuid", Type: "string"},
	{Path: "display_name", Type: "string"},
}

// mkUserList builds a setLikeListValue over the userFields schema with the
// given (uuid, display_name) pairs. Helper to keep tests readable.
func mkUserList(t *testing.T, items ...[2]string) setLikeListValue {
	t.Helper()
	attrTypes := itemAttrTypes(userFields)
	objType := types.ObjectType{AttrTypes: attrTypes}
	elements := make([]attr.Value, 0, len(items))
	for _, it := range items {
		obj, diags := types.ObjectValue(attrTypes, map[string]attr.Value{
			"uuid":         types.StringValue(it[0]),
			"display_name": types.StringValue(it[1]),
		})
		if diags.HasError() {
			t.Fatalf("failed to build object: %#v", diags)
		}
		elements = append(elements, obj)
	}
	listVal, diags := types.ListValue(objType, elements)
	if diags.HasError() {
		t.Fatalf("failed to build list: %#v", diags)
	}
	v, d := setLikeListTypeFor(userFields).ValueFromList(context.Background(), listVal)
	if d.HasError() {
		t.Fatalf("ValueFromList: %#v", d)
	}
	return v.(setLikeListValue)
}

// TestSetLikeListSemanticEqualsIgnoresElementOrder is the core promise of
// the custom list type: when the API returns the same set of users in a
// different order than the operator wrote in HCL, the framework must treat
// the two values as semantically equal — no "Provider produced invalid plan"
// or "inconsistent result after apply" diagnostic.
func TestSetLikeListSemanticEqualsIgnoresElementOrder(t *testing.T) {
	configOrder := mkUserList(t,
		[2]string{"{bbbb-bbbb}", "Bob"},
		[2]string{"{aaaa-aaaa}", "Alice"},
	)
	apiOrder := mkUserList(t,
		[2]string{"{aaaa-aaaa}", "Alice"},
		[2]string{"{bbbb-bbbb}", "Bob"},
	)
	eq, diags := configOrder.ListSemanticEquals(context.Background(), apiOrder)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %#v", diags)
	}
	if !eq {
		t.Fatalf("two lists with the same elements in different order must be semantically equal:\n  config = %s\n  api    = %s", configOrder.String(), apiOrder.String())
	}
}

// TestSetLikeListSemanticEqualsDistinguishesAddRemove guards the inverse:
// genuine differences (an extra/missing element) must still be visible to
// Terraform so the plan and apply machinery actually do something.
func TestSetLikeListSemanticEqualsDistinguishesAddRemove(t *testing.T) {
	twoUsers := mkUserList(t,
		[2]string{"{aaaa}", "Alice"},
		[2]string{"{bbbb}", "Bob"},
	)
	threeUsers := mkUserList(t,
		[2]string{"{aaaa}", "Alice"},
		[2]string{"{bbbb}", "Bob"},
		[2]string{"{cccc}", "Carol"},
	)
	eq, diags := twoUsers.ListSemanticEquals(context.Background(), threeUsers)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %#v", diags)
	}
	if eq {
		t.Fatalf("adding an element must change semantic equality:\n  before = %s\n  after  = %s", twoUsers.String(), threeUsers.String())
	}
}

// TestSetLikeListSemanticEqualsHandlesDuplicateIdentities verifies that two
// items sharing the same identity value (e.g. uuid) but differing in
// another attribute are still distinguished — i.e. semantic equality uses
// the canonical-JSON tiebreaker, not just the identity field. Otherwise an
// API change to one of the duplicate elements would be silently swallowed.
func TestSetLikeListSemanticEqualsHandlesDuplicateIdentities(t *testing.T) {
	a := mkUserList(t,
		[2]string{"{same}", "Alice"},
		[2]string{"{same}", "Bob"},
	)
	b := mkUserList(t,
		[2]string{"{same}", "Alice"},
		[2]string{"{same}", "Charlie"}, // different display_name
	)
	eq, diags := a.ListSemanticEquals(context.Background(), b)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %#v", diags)
	}
	if eq {
		t.Fatalf("differing element values must not be semantically equal even when identity field collides:\n  a = %s\n  b = %s", a.String(), b.String())
	}

	// And the symmetric same-elements case must be equal.
	c := mkUserList(t,
		[2]string{"{same}", "Bob"},
		[2]string{"{same}", "Alice"}, // reversed
	)
	eq, diags = a.ListSemanticEquals(context.Background(), c)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %#v", diags)
	}
	if !eq {
		t.Fatalf("identical multisets in different order must be semantically equal:\n  a = %s\n  c = %s", a.String(), c.String())
	}
}

// TestSetLikeListSemanticEqualsNullAndUnknown ensures the no-op edge cases
// don't panic and don't claim equality across known/null/unknown boundaries.
func TestSetLikeListSemanticEqualsNullAndUnknown(t *testing.T) {
	known := mkUserList(t, [2]string{"{aaaa}", "Alice"})
	objType := types.ObjectType{AttrTypes: itemAttrTypes(userFields)}

	null, d := setLikeListTypeFor(userFields).ValueFromList(context.Background(), types.ListNull(objType))
	if d.HasError() {
		t.Fatalf("ValueFromList(null): %#v", d)
	}
	unknown, d := setLikeListTypeFor(userFields).ValueFromList(context.Background(), types.ListUnknown(objType))
	if d.HasError() {
		t.Fatalf("ValueFromList(unknown): %#v", d)
	}

	for _, tc := range []struct {
		name string
		a, b basetypes.ListValuable
	}{
		{"known vs null", known, null},
		{"known vs unknown", known, unknown},
		{"null vs unknown", null, unknown},
	} {
		t.Run(tc.name, func(t *testing.T) {
			a := tc.a.(basetypes.ListValuableWithSemanticEquals)
			eq, diags := a.ListSemanticEquals(context.Background(), tc.b)
			if diags.HasError() {
				t.Fatalf("unexpected diagnostics: %#v", diags)
			}
			if eq {
				t.Fatalf("%s must not be semantically equal", tc.name)
			}
		})
	}
}

// TestNestedObjectArrayResourceAttrsUseSetLikeListType asserts that all four
// schema construction sites that previously attached the
// nestedListSortPlanModifier now emit a ListNestedAttribute whose CustomType
// is the order-insensitive setLikeListType — without the legacy plan
// modifier. This is the wiring contract the rest of the runtime depends on.
func TestNestedObjectArrayResourceAttrsUseSetLikeListType(t *testing.T) {
	itemFields := []BodyFieldDef{{Path: "uuid", Type: "string"}}

	checks := map[string]resourceschema.ListNestedAttribute{
		"bodyFieldAttr":     mustListNested(t, bodyFieldAttr(BodyFieldDef{Path: "users", IsArray: true, ItemFields: itemFields})),
		"responseFieldAttr": mustListNested(t, responseFieldAttr(BodyFieldDef{Path: "users", IsArray: true, ItemFields: itemFields})),
		"buildNestedItemAttrs": mustListNested(t, buildNestedItemAttrs([]BodyFieldDef{
			{Path: "users", IsArray: true, ItemFields: itemFields},
		})["users"]),
		"mergeListNestedResponseAttr": mustListNested(t, mergeResponseAttr(
			resourceschema.ListNestedAttribute{Optional: true, NestedObject: resourceschema.NestedAttributeObject{Attributes: buildNestedItemAttrs(itemFields)}},
			BodyFieldDef{Path: "users", IsArray: true, ItemFields: itemFields},
		)),
	}

	for name, attr := range checks {
		t.Run(name, func(t *testing.T) {
			if attr.CustomType == nil {
				t.Fatalf("%s: ListNestedAttribute.CustomType is nil; want setLikeListType", name)
			}
			if _, ok := attr.CustomType.(setLikeListType); !ok {
				t.Fatalf("%s: CustomType is %T; want setLikeListType", name, attr.CustomType)
			}
			if len(attr.PlanModifiers) != 0 {
				t.Fatalf("%s: PlanModifiers must be empty (semantic equality replaces the legacy sort modifier); got %#v", name, attr.PlanModifiers)
			}
		})
	}
}

// TestScalarListAttributeRemainsPlainListType guards the converse: simple
// scalar arrays (e.g. `tags`) are genuinely ordered and must not be wrapped
// in the order-insensitive custom type. They also keep no PlanModifiers.
func TestScalarListAttributeRemainsPlainListType(t *testing.T) {
	bodyAttr, ok := bodyFieldAttr(BodyFieldDef{Path: "tags", IsArray: true}).(resourceschema.ListAttribute)
	if !ok {
		t.Fatalf("bodyFieldAttr(scalar array) returned %T, want ListAttribute", bodyFieldAttr(BodyFieldDef{Path: "tags", IsArray: true}))
	}
	if bodyAttr.CustomType != nil {
		t.Fatalf("scalar ListAttribute must not carry CustomType; got %#v", bodyAttr.CustomType)
	}

	respAttr, ok := responseFieldAttr(BodyFieldDef{Path: "tags", IsArray: true}).(resourceschema.ListAttribute)
	if !ok {
		t.Fatalf("responseFieldAttr(scalar array) returned %T, want ListAttribute", responseFieldAttr(BodyFieldDef{Path: "tags", IsArray: true}))
	}
	if respAttr.CustomType != nil {
		t.Fatalf("scalar response ListAttribute must not carry CustomType; got %#v", respAttr.CustomType)
	}
}

// TestSetLikeListTypeRoundTripsThroughTerraformValue ensures the framework's
// state-marshaling round-trip preserves the custom type — without this the
// plugin protocol would silently downgrade values to the default ListValue
// and the semantic equality logic would never run during refresh.
func TestSetLikeListTypeRoundTripsThroughTerraformValue(t *testing.T) {
	ctx := context.Background()
	original := mkUserList(t,
		[2]string{"{aaaa}", "Alice"},
		[2]string{"{bbbb}", "Bob"},
	)
	tfVal, err := original.ToTerraformValue(ctx)
	if err != nil {
		t.Fatalf("ToTerraformValue: %v", err)
	}
	roundTripped, err := setLikeListTypeFor(userFields).ValueFromTerraform(ctx, tfVal)
	if err != nil {
		t.Fatalf("ValueFromTerraform: %v", err)
	}
	if _, ok := roundTripped.(setLikeListValue); !ok {
		t.Fatalf("round-trip yielded %T; want setLikeListValue", roundTripped)
	}
}

func mustListNested(t *testing.T, a resourceschema.Attribute) resourceschema.ListNestedAttribute {
	t.Helper()
	listAttr, ok := a.(resourceschema.ListNestedAttribute)
	if !ok {
		t.Fatalf("attribute %T is not a ListNestedAttribute", a)
	}
	return listAttr
}
