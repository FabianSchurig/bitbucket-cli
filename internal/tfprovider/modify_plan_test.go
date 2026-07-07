package tfprovider

import (
	"context"
	"testing"

	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// TestAttrValuesSemanticallyEqual_StringAttribute verifies that the helper
// reports tftypes equality for plain attributes — i.e. a normal Update
// where the user changes an attribute is detected as a real change so
// ModifyPlan does NOT substitute prior state and Computed siblings get
// refreshed from the API. This guards against the regression where an
// over-broad UseStateForUnknown caused "Provider produced inconsistent
// result after apply" errors on legitimate updates.
func TestAttrValuesSemanticallyEqual_StringAttribute(t *testing.T) {
	ctx := context.Background()
	a := resourceschema.StringAttribute{Optional: true}

	same := tftypes.NewValue(tftypes.String, "foo")
	other := tftypes.NewValue(tftypes.String, "foo")
	if equal, ok := attrValuesSemanticallyEqual(ctx, same, other, a); !ok || !equal {
		t.Fatalf("identical strings should be equal, got equal=%v ok=%v", equal, ok)
	}

	changed := tftypes.NewValue(tftypes.String, "bar")
	if equal, ok := attrValuesSemanticallyEqual(ctx, same, changed, a); !ok || equal {
		t.Fatalf("differing strings should be unequal, got equal=%v ok=%v", equal, ok)
	}
}

// TestAttrValuesSemanticallyEqual_SetLikeList verifies that set-like
// nested-object lists are treated as equal under reordering — the
// general fix for "user reordered list ⇒ empty plan", applied uniformly
// across every list-like API endpoint via the resource ModifyPlan.
func TestAttrValuesSemanticallyEqual_SetLikeList(t *testing.T) {
	ctx := context.Background()
	itemFields := []BodyFieldDef{{Path: "uuid", Type: "string"}}
	objType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{"uuid": tftypes.String}}
	listTfType := tftypes.List{ElementType: objType}

	mkItem := func(uuid string) tftypes.Value {
		return tftypes.NewValue(objType, map[string]tftypes.Value{
			"uuid": tftypes.NewValue(tftypes.String, uuid),
		})
	}
	abc := tftypes.NewValue(listTfType, []tftypes.Value{mkItem("a"), mkItem("b"), mkItem("c")})
	cba := tftypes.NewValue(listTfType, []tftypes.Value{mkItem("c"), mkItem("b"), mkItem("a")})
	abd := tftypes.NewValue(listTfType, []tftypes.Value{mkItem("a"), mkItem("b"), mkItem("d")})

	a := resourceschema.ListNestedAttribute{
		Optional:   true,
		Computed:   true,
		CustomType: setLikeListTypeFor(itemFields),
		NestedObject: resourceschema.NestedAttributeObject{
			Attributes: buildNestedItemAttrs(itemFields),
		},
	}

	if equal, ok := attrValuesSemanticallyEqual(ctx, abc, cba, a); !ok || !equal {
		t.Fatalf("reordered set-like lists should be equal, got equal=%v ok=%v", equal, ok)
	}
	if equal, ok := attrValuesSemanticallyEqual(ctx, abc, abd, a); !ok || equal {
		t.Fatalf("set-like lists with different items should be unequal, got equal=%v ok=%v", equal, ok)
	}
}

// TestAttrValuesSemanticallyEqual_SetLikeList_NullComputedInnerFields is a
// regression test for the branch-restrictions OrderInsensitiveUsers failure.
// ModifyPlan compares Config vs State: the config side of a nested-object list
// carries only user-supplied fields (e.g. {uuid}), with Computed inner fields
// (display_name, links, …) left NULL — whereas prior state has them populated
// from the API. When the list is reordered, the set must still compare equal so
// ModifyPlan substitutes prior state and does not promote sibling Computed
// attributes to Unknown (which produced a perpetual in-place diff).
func TestAttrValuesSemanticallyEqual_SetLikeList_NullComputedInnerFields(t *testing.T) {
	ctx := context.Background()
	itemFields := []BodyFieldDef{
		{Path: "uuid", Type: "string"},
		{Path: "display_name", Type: "string"},
	}
	objType := tftypes.Object{AttributeTypes: map[string]tftypes.Type{
		"uuid":         tftypes.String,
		"display_name": tftypes.String,
	}}
	listTfType := tftypes.List{ElementType: objType}

	// config item: only uuid set, display_name NULL (unspecified computed field)
	cfgItem := func(uuid string) tftypes.Value {
		return tftypes.NewValue(objType, map[string]tftypes.Value{
			"uuid":         tftypes.NewValue(tftypes.String, uuid),
			"display_name": tftypes.NewValue(tftypes.String, nil),
		})
	}
	// state item: uuid + populated display_name (from the API)
	stItem := func(uuid, name string) tftypes.Value {
		return tftypes.NewValue(objType, map[string]tftypes.Value{
			"uuid":         tftypes.NewValue(tftypes.String, uuid),
			"display_name": tftypes.NewValue(tftypes.String, name),
		})
	}
	cfg := tftypes.NewValue(listTfType, []tftypes.Value{cfgItem("b"), cfgItem("a")})
	st := tftypes.NewValue(listTfType, []tftypes.Value{stItem("a", "Alice"), stItem("b", "Bob")})

	a := resourceschema.ListNestedAttribute{
		Optional:   true,
		Computed:   true,
		CustomType: setLikeListTypeFor(itemFields),
		NestedObject: resourceschema.NestedAttributeObject{
			Attributes: buildNestedItemAttrs(itemFields),
		},
	}
	if equal, ok := attrValuesSemanticallyEqual(ctx, cfg, st, a); !ok || !equal {
		t.Fatalf("reordered set-like list with null computed inner fields should be equal, got equal=%v ok=%v", equal, ok)
	}
	// A genuine change to a user-supplied field must still be detected.
	stChanged := tftypes.NewValue(listTfType, []tftypes.Value{stItem("a", "Alice"), stItem("c", "Carol")})
	if equal, ok := attrValuesSemanticallyEqual(ctx, cfg, stChanged, a); !ok || equal {
		t.Fatalf("set-like list with a different uuid must be unequal, got equal=%v ok=%v", equal, ok)
	}
}

// TestIsConfigurableAttr ensures Computed-only attributes (id,
// api_response, response fields) are skipped by the ModifyPlan walk,
// otherwise the raw equality check would always flag them as "changed"
// (config null vs state populated by the API).
func TestIsConfigurableAttr(t *testing.T) {
	cases := map[string]struct {
		attr resourceschema.Attribute
		want bool
	}{
		"required":          {resourceschema.StringAttribute{Required: true}, true},
		"optional":          {resourceschema.StringAttribute{Optional: true}, true},
		"optional+computed": {resourceschema.StringAttribute{Optional: true, Computed: true}, true},
		"computed_only":     {resourceschema.StringAttribute{Computed: true}, false},
	}
	for name, c := range cases {
		if got := isConfigurableAttr(c.attr); got != c.want {
			t.Errorf("%s: isConfigurableAttr=%v want=%v", name, got, c.want)
		}
	}
}
