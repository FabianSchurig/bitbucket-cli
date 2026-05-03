package tfprovider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestBuildListFromResponseSortsNestedObjectsByStableKey reproduces the
// `Provider produced inconsistent result after apply` failure described in
// the bitbucket_branch_restrictions `users` issue: the API returns the same
// elements in a non-deterministic order, so two equivalent responses must
// yield byte-identical Terraform state.
//
// The fix sorts nested-object items by a stable identity key (uuid > id >
// slug > full_slug > name > canonical JSON). Two responses containing the
// same set of users in different order must produce the same list.
func TestBuildListFromResponseSortsNestedObjectsByStableKey(t *testing.T) {
	userFields := []BodyFieldDef{
		{Path: "uuid", Type: "string"},
		{Path: "display_name", Type: "string"},
	}

	userA := map[string]any{"uuid": "{aaaa-aaaa}", "display_name": "Alice"}
	userB := map[string]any{"uuid": "{bbbb-bbbb}", "display_name": "Bob"}

	planOrder := buildListFromResponse([]any{userA, userB}, userFields)
	apiOrder := buildListFromResponse([]any{userB, userA}, userFields)

	if !planOrder.Equal(apiOrder) {
		t.Fatalf("nested-object list order must be deterministic regardless of input order:\n  plan-order = %s\n   api-order = %s", planOrder.String(), apiOrder.String())
	}

	// Sanity-check the canonical order: lower uuid first.
	first := planOrder.Elements()[0].(types.Object).Attributes()["uuid"].(types.String).ValueString()
	if first != "{aaaa-aaaa}" {
		t.Fatalf("expected sorted-by-uuid first element to be {aaaa-aaaa}, got %s", first)
	}
}

// TestBuildListFromResponseTiebreakerForDuplicateIdentity guards the total-
// ordering guarantee: when two items share the same identity-field value,
// the canonical JSON tiebreaker still produces a deterministic, reproducible
// order regardless of the input order.
func TestBuildListFromResponseTiebreakerForDuplicateIdentity(t *testing.T) {
	fields := []BodyFieldDef{
		{Path: "uuid", Type: "string"},
		{Path: "display_name", Type: "string"},
	}
	dupA := map[string]any{"uuid": "{same}", "display_name": "Alice"}
	dupB := map[string]any{"uuid": "{same}", "display_name": "Bob"}

	planOrder := buildListFromResponse([]any{dupA, dupB}, fields)
	apiOrder := buildListFromResponse([]any{dupB, dupA}, fields)
	if !planOrder.Equal(apiOrder) {
		t.Fatalf("duplicate-identity items must sort deterministically via tiebreaker:\n  plan-order = %s\n   api-order = %s", planOrder.String(), apiOrder.String())
	}
}

// TestBuildListFromResponseSortKeyFallbacks verifies the identity-key
// precedence: items without uuid fall back to id, then slug, then full_slug,
// then name, then a canonical JSON tiebreaker.
func TestBuildListFromResponseSortKeyFallbacks(t *testing.T) {
	cases := []struct {
		name   string
		fields []BodyFieldDef
		a, b   map[string]any
		want   string // expected first element value of the leading key
		key    string
	}{
		{
			name:   "id fallback",
			fields: []BodyFieldDef{{Path: "id", Type: "int"}},
			a:      map[string]any{"id": 2.0},
			b:      map[string]any{"id": 1.0},
			key:    "id",
			want:   "1",
		},
		{
			name:   "slug fallback",
			fields: []BodyFieldDef{{Path: "slug", Type: "string"}},
			a:      map[string]any{"slug": "zebra"},
			b:      map[string]any{"slug": "apple"},
			key:    "slug",
			want:   "apple",
		},
		{
			name:   "name fallback",
			fields: []BodyFieldDef{{Path: "name", Type: "string"}},
			a:      map[string]any{"name": "Bob"},
			b:      map[string]any{"name": "Alice"},
			key:    "name",
			want:   "Alice",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := buildListFromResponse([]any{tc.a, tc.b}, tc.fields)
			first := got.Elements()[0].(types.Object).Attributes()[tc.key]
			var actual string
			switch v := first.(type) {
			case types.String:
				actual = v.ValueString()
			case types.Int64:
				actual = "1"
				if v.ValueInt64() != 1 {
					actual = "(unexpected)"
				}
			default:
				t.Fatalf("unexpected attr type %T", first)
			}
			if actual != tc.want {
				t.Fatalf("first element %s = %q, want %q", tc.key, actual, tc.want)
			}
		})
	}
}

// TestNestedObjectArrayResourceAttrsAttachSortPlanModifier asserts that the
// schema builders attach the deterministic-order plan modifier to every
// nested-object array (`ListNestedAttribute`) they emit. Without the
// modifier, the user's plan keeps its config order while the post-apply
// state ends up canonicalised — so the two diverge and Terraform tags the
// result as inconsistent.
//
// Simple scalar lists (ListAttribute) intentionally do NOT get the modifier,
// because they have no per-item identity field to sort by.
func TestNestedObjectArrayResourceAttrsAttachSortPlanModifier(t *testing.T) {
	itemFields := []BodyFieldDef{{Path: "uuid", Type: "string"}}

	// 1. bodyFieldAttr — request body field for nested-object array.
	bodyAttr, ok := bodyFieldAttr(BodyFieldDef{Path: "users", IsArray: true, ItemFields: itemFields}).(resourceschema.ListNestedAttribute)
	if !ok {
		t.Fatalf("bodyFieldAttr returned %T, want ListNestedAttribute", bodyFieldAttr(BodyFieldDef{Path: "users", IsArray: true, ItemFields: itemFields}))
	}
	if !hasNestedListSortModifier(bodyAttr.PlanModifiers) {
		t.Fatalf("bodyFieldAttr nested-object array missing nestedListSortPlanModifier; got %#v", bodyAttr.PlanModifiers)
	}

	// 2. responseFieldAttr — computed-only response array.
	respAttr, ok := responseFieldAttr(BodyFieldDef{Path: "users", IsArray: true, ItemFields: itemFields}).(resourceschema.ListNestedAttribute)
	if !ok {
		t.Fatalf("responseFieldAttr returned %T, want ListNestedAttribute", responseFieldAttr(BodyFieldDef{Path: "users", IsArray: true, ItemFields: itemFields}))
	}
	if !hasNestedListSortModifier(respAttr.PlanModifiers) {
		t.Fatalf("responseFieldAttr nested-object array missing nestedListSortPlanModifier; got %#v", respAttr.PlanModifiers)
	}

	// 3. buildNestedItemAttrs — nested-object array inside a parent object.
	nested := buildNestedItemAttrs([]BodyFieldDef{
		{Path: "users", IsArray: true, ItemFields: itemFields},
	})
	listNested, ok := nested["users"].(resourceschema.ListNestedAttribute)
	if !ok {
		t.Fatalf("buildNestedItemAttrs[users] = %T, want ListNestedAttribute", nested["users"])
	}
	if !hasNestedListSortModifier(listNested.PlanModifiers) {
		t.Fatalf("buildNestedItemAttrs nested-object array missing nestedListSortPlanModifier; got %#v", listNested.PlanModifiers)
	}

	// 4. mergeListNestedResponseAttr — when a body field is later promoted
	//    to also satisfy a Read-side response field.
	merged, ok := mergeResponseAttr(
		resourceschema.ListNestedAttribute{Optional: true, NestedObject: resourceschema.NestedAttributeObject{Attributes: buildNestedItemAttrs(itemFields)}},
		BodyFieldDef{Path: "users", IsArray: true, ItemFields: itemFields},
	).(resourceschema.ListNestedAttribute)
	if !ok {
		t.Fatalf("mergeResponseAttr returned non-ListNestedAttribute")
	}
	if !hasNestedListSortModifier(merged.PlanModifiers) {
		t.Fatalf("mergeListNestedResponseAttr missing nestedListSortPlanModifier; got %#v", merged.PlanModifiers)
	}

	// Sanity: simple scalar lists (no ItemFields) do NOT get the sort modifier,
	// because there is no per-item identity field to sort by.
	tagsAttr, ok := bodyFieldAttr(BodyFieldDef{Path: "tags", IsArray: true}).(resourceschema.ListAttribute)
	if !ok {
		t.Fatalf("expected ListAttribute for scalar array, got %T", bodyFieldAttr(BodyFieldDef{Path: "tags", IsArray: true}))
	}
	if len(tagsAttr.PlanModifiers) != 0 {
		t.Fatalf("scalar ListAttribute must not carry the nested-object sort modifier; got %#v", tagsAttr.PlanModifiers)
	}
}

// TestNestedListSortPlanModifierSortsPlanValueByStableKey exercises the
// plan-side half of the fix: when a user writes `users = [B, A]` in their
// configuration and the provider sorts the API response by uuid, the planned
// value must be sorted the same way so Terraform's post-apply consistency
// check sees plan == state.
func TestNestedListSortPlanModifierSortsPlanValueByStableKey(t *testing.T) {
	itemFields := []BodyFieldDef{{Path: "uuid", Type: "string"}}
	objType := types.ObjectType{AttrTypes: itemAttrTypes(itemFields)}

	mkObj := func(uuid string) types.Object {
		return types.ObjectValueMust(itemAttrTypes(itemFields), map[string]attr.Value{
			"uuid": types.StringValue(uuid),
		})
	}
	configOrder := types.ListValueMust(objType, []attr.Value{
		mkObj("{bbbb-bbbb}"),
		mkObj("{aaaa-aaaa}"),
	})

	mod := newNestedListSortPlanModifier(itemFields)
	req := planmodifier.ListRequest{
		Path:      path.Root("users"),
		PlanValue: configOrder,
	}
	resp := &planmodifier.ListResponse{PlanValue: configOrder}
	mod.PlanModifyList(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected diagnostics: %#v", resp.Diagnostics)
	}
	first := resp.PlanValue.Elements()[0].(types.Object).Attributes()["uuid"].(types.String).ValueString()
	if first != "{aaaa-aaaa}" {
		t.Fatalf("plan modifier must sort plan elements by uuid; first uuid = %q, want {aaaa-aaaa}", first)
	}
}

// TestNestedListSortPlanModifierLeavesUnknownAndNullAlone guards the
// no-op edge cases — unknown plan values (e.g. "(known after apply)") and
// null values must be passed through untouched.
func TestNestedListSortPlanModifierLeavesUnknownAndNullAlone(t *testing.T) {
	itemFields := []BodyFieldDef{{Path: "uuid", Type: "string"}}
	objType := types.ObjectType{AttrTypes: itemAttrTypes(itemFields)}
	mod := newNestedListSortPlanModifier(itemFields)

	for _, v := range []types.List{
		types.ListUnknown(objType),
		types.ListNull(objType),
	} {
		req := planmodifier.ListRequest{Path: path.Root("users"), PlanValue: v}
		resp := &planmodifier.ListResponse{PlanValue: v}
		mod.PlanModifyList(context.Background(), req, resp)
		if !resp.PlanValue.Equal(v) {
			t.Fatalf("plan modifier must pass through %s untouched; got %s", v.String(), resp.PlanValue.String())
		}
	}
}
