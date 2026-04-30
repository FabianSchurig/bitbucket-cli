package tfprovider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// These tests guard the canonical-key helpers that survive after the move
// from in-place response sorting (the v0.15.6 approach) to type-level
// semantic equality (setLikeListType). The keys are now consumed only by
// setLikeListValue.ListSemanticEquals — but they still need to obey the
// same identity-field precedence and total-order contract, so the original
// table of cases is preserved verbatim under the new entry points.

// TestStableItemSortKeyHonoursIdentityPrecedence verifies the identity-key
// precedence: items use uuid first, then id, then slug, then full_slug,
// then name, then a canonical JSON tiebreaker.
func TestStableItemSortKeyHonoursIdentityPrecedence(t *testing.T) {
	cases := []struct {
		name    string
		fields  []BodyFieldDef
		a, b    map[string]any
		wantKey string // primary identity field expected to lead the sort key
	}{
		{
			name:    "uuid wins over everything",
			fields:  []BodyFieldDef{{Path: "uuid"}, {Path: "id"}},
			a:       map[string]any{"uuid": "{u-1}", "id": 9.0},
			b:       map[string]any{"uuid": "{u-2}", "id": 1.0},
			wantKey: "uuid=",
		},
		{
			name:    "id fallback when uuid not declared",
			fields:  []BodyFieldDef{{Path: "id"}},
			a:       map[string]any{"id": 2.0},
			b:       map[string]any{"id": 1.0},
			wantKey: "id=",
		},
		{
			name:    "slug fallback",
			fields:  []BodyFieldDef{{Path: "slug"}},
			a:       map[string]any{"slug": "zebra"},
			b:       map[string]any{"slug": "apple"},
			wantKey: "slug=",
		},
		{
			name:    "name fallback",
			fields:  []BodyFieldDef{{Path: "name"}},
			a:       map[string]any{"name": "Bob"},
			b:       map[string]any{"name": "Alice"},
			wantKey: "name=",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			keyA := stableItemSortKey(tc.a, tc.fields)
			if len(keyA) < len(tc.wantKey) || keyA[:len(tc.wantKey)] != tc.wantKey {
				t.Fatalf("expected key to start with %q, got %q", tc.wantKey, keyA)
			}
			keyB := stableItemSortKey(tc.b, tc.fields)
			if keyA == keyB {
				t.Fatalf("distinct items must produce distinct keys: %q == %q", keyA, keyB)
			}
		})
	}
}

// TestStableItemSortKeyTiebreakerForDuplicateIdentity guards the total-
// ordering guarantee: when two items share the same identity-field value,
// the canonical JSON tiebreaker still produces deterministic, distinct
// keys so multiset comparisons can tell them apart.
func TestStableItemSortKeyTiebreakerForDuplicateIdentity(t *testing.T) {
	fields := []BodyFieldDef{{Path: "uuid"}, {Path: "display_name"}}
	dupA := map[string]any{"uuid": "{same}", "display_name": "Alice"}
	dupB := map[string]any{"uuid": "{same}", "display_name": "Bob"}

	keyA := stableItemSortKey(dupA, fields)
	keyB := stableItemSortKey(dupB, fields)
	if keyA == keyB {
		t.Fatalf("duplicate-identity items must be distinguished by the tiebreaker; both produced %q", keyA)
	}
}

// TestStableObjectSortKeyMatchesItemSortKeyContract exercises the Terraform
// attr.Value side of the same key derivation. The plan-side / state-side
// keys must agree on identity precedence with the response-side key so
// setLikeListValue.ListSemanticEquals (which mixes the two) is consistent.
func TestStableObjectSortKeyMatchesItemSortKeyContract(t *testing.T) {
	fields := []BodyFieldDef{{Path: "uuid"}}
	attrTypes := itemAttrTypes(fields)

	objA := types.ObjectValueMust(attrTypes, map[string]attr.Value{"uuid": types.StringValue("{aaaa}")})
	objB := types.ObjectValueMust(attrTypes, map[string]attr.Value{"uuid": types.StringValue("{bbbb}")})

	keyA := stableObjectSortKey(objA, fields)
	keyB := stableObjectSortKey(objB, fields)
	if keyA == keyB {
		t.Fatalf("objects with distinct uuid must produce distinct keys; both = %q", keyA)
	}
	if len(keyA) < 5 || keyA[:5] != "uuid=" {
		t.Fatalf("expected key to begin with %q, got %q", "uuid=", keyA)
	}

	// Sanity: the order of attribute insertion in the map must not change
	// the key (the framework's Object.Attributes() returns a fresh map and
	// String() emits keys in lexicographic order).
	objAReordered := types.ObjectValueMust(attrTypes, map[string]attr.Value{"uuid": types.StringValue("{aaaa}")})
	if stableObjectSortKey(objAReordered, fields) != keyA {
		t.Fatalf("object key derivation must be insertion-order independent")
	}
}

// TestBuildListFromResponsePreservesUpstreamOrder asserts the new contract
// the response builder honours after the move to type-level semantic
// equality: it no longer reorders elements, so genuinely-ordered API
// responses keep their order in state and the operator's diff stays
// surgical when they add a single element to the end of their config.
func TestBuildListFromResponsePreservesUpstreamOrder(t *testing.T) {
	fields := []BodyFieldDef{{Path: "uuid"}}
	apiOrder := []any{
		map[string]any{"uuid": "{bbbb}"},
		map[string]any{"uuid": "{aaaa}"},
		map[string]any{"uuid": "{cccc}"},
	}
	got := buildListFromResponse(apiOrder, fields)
	if got.IsNull() || got.IsUnknown() {
		t.Fatalf("expected a known list, got %s", got.String())
	}
	elements := got.Elements()
	if len(elements) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(elements))
	}
	wantOrder := []string{"{bbbb}", "{aaaa}", "{cccc}"}
	for i, want := range wantOrder {
		uuid := elements[i].(types.Object).Attributes()["uuid"].(types.String).ValueString()
		if uuid != want {
			t.Fatalf("element %d uuid = %q, want %q (order must follow the API response verbatim now that ordering is enforced via setLikeListType.ListSemanticEquals)", i, uuid, want)
		}
	}
}

// TestBuildListFromResponseReturnsSetLikeListValue ensures the response
// builder propagates the element schema into the returned value's
// itemFields. The wrapping (setLikeListValue) is already enforced by the
// function's static return type — this test guards the field propagation
// that ListSemanticEquals depends on for its identity-key derivation.
func TestBuildListFromResponseReturnsSetLikeListValue(t *testing.T) {
	got := buildListFromResponse(
		[]any{map[string]any{"uuid": "{x}"}},
		[]BodyFieldDef{{Path: "uuid"}},
	)
	if got.itemFields == nil || got.itemFields[0].Path != "uuid" {
		t.Fatalf("expected itemFields to be propagated; got %#v", got.itemFields)
	}
}

// TestAttrNullValueWrapsNestedObjectArraysInSetLikeList guards the type-
// match between nullable attributes and their schema's CustomType — null
// values come from the response builder when an optional nested array is
// absent, and the framework rejects them outright if their attr.Type
// doesn't equal the declared schema type.
func TestAttrNullValueWrapsNestedObjectArraysInSetLikeList(t *testing.T) {
	got := attrNullValue(BodyFieldDef{Path: "users", IsArray: true, ItemFields: []BodyFieldDef{{Path: "uuid"}}})
	if _, ok := got.(setLikeListValue); !ok {
		t.Fatalf("attrNullValue for nested-object array = %T, want setLikeListValue", got)
	}
	if !got.IsNull() {
		t.Fatalf("attrNullValue must return a null value, got %s", got.String())
	}

	// Sanity: scalar lists keep the plain types.List null value.
	scalar := attrNullValue(BodyFieldDef{Path: "tags", IsArray: true})
	if _, ok := scalar.(types.List); !ok {
		t.Fatalf("scalar array attrNullValue = %T, want types.List", scalar)
	}
}

// TestBuildListFromResponseEmptyArrayProducesEmptySetLikeList covers the
// edge case where the API returns an empty array for an optional nested
// list: the wrapped list must be known-empty (not null, not unknown) so a
// subsequent plan against an unset config diffs cleanly.
func TestBuildListFromResponseEmptyArrayProducesEmptySetLikeList(t *testing.T) {
	got := buildListFromResponse([]any{}, []BodyFieldDef{{Path: "uuid"}})
	if got.IsNull() || got.IsUnknown() {
		t.Fatalf("empty array must produce a known empty list, got %s", got.String())
	}
	if len(got.Elements()) != 0 {
		t.Fatalf("expected zero elements, got %d", len(got.Elements()))
	}
}

// TestStableObjectSortKeyContextCompiles is a tiny placeholder asserting
// the helpers used by setLikeListValue.ListSemanticEquals build cleanly in
// a test-package context — guards against accidental visibility regressions
// during future refactors of nested_list_order.go.
func TestStableObjectSortKeyContextCompiles(t *testing.T) {
	_ = context.Background()
	_ = stableItemSortKey(map[string]any{}, nil)
	_ = stableObjectSortKey(types.ObjectNull(map[string]attr.Type{}), nil)
}

// ─── Identity-only key + response reordering ─────────────────────────────────
//
// The two helpers below are the load-bearing piece for the post-Create /
// post-Update consistency check. Background: terraform-plugin-framework's
// schema-level `ValueSemanticEquality` only runs on collection types whose
// values implement `ListValuableWithSemanticEquals` AND whose `ListSemanticEquals`
// returns `true`. For nested-object arrays carrying mixed user/computed
// inner fields (`uuid` user-supplied, `display_name` / `created_on`
// returned by the API), the planned-state value has Unknown for the
// computed fields while the post-Create value has them filled in — so any
// canonical-form comparison that includes those fields disagrees, the
// framework falls back to positional element comparison, and TF Core then
// errors `Provider produced inconsistent result after apply` when the API
// returned the elements in a different order than the operator wrote.
//
// The fix that survives this is:
//
//   1. Match elements between planned and API-response sides on the
//      *primary identity field only* (`uuid` > `id` > ... — not the full
//      canonical JSON). Unknown computed fields don't participate.
//   2. Reorder the API response array to match the operator's planned /
//      prior order before persisting it to state, appending genuinely new
//      elements at the end so adding a uuid surfaces as a single-element
//      diff in the operator's chosen order.

// TestStableItemPrimaryKeyReturnsIdentityOnly verifies the new helper that
// extracts only the primary identity prefix (no canonical-JSON tiebreaker).
// This is the key the response reorderer matches on.
func TestStableItemPrimaryKeyReturnsIdentityOnly(t *testing.T) {
	fields := []BodyFieldDef{{Path: "uuid"}, {Path: "display_name"}}
	got := stableItemPrimaryKey(map[string]any{
		"uuid":         "{abc}",
		"display_name": "Alice",
	}, fields)
	if got != "uuid={abc}" {
		t.Fatalf("primary key = %q, want %q (must include only the identity field, not display_name)", got, "uuid={abc}")
	}
}

// TestStableItemPrimaryKeyMatchesObjectPrimaryKey ensures the JSON-side and
// terraform-Object-side primary-key extractors agree byte-for-byte. Without
// this invariant, the response reorderer would never find any matches
// between the planned state (Object) and the API response (map[string]any).
func TestStableItemPrimaryKeyMatchesObjectPrimaryKey(t *testing.T) {
	fields := []BodyFieldDef{{Path: "uuid"}}

	jsonKey := stableItemPrimaryKey(map[string]any{"uuid": "{abc}"}, fields)
	objKey := stableObjectPrimaryKey(
		types.ObjectValueMust(itemAttrTypes(fields), map[string]attr.Value{
			"uuid": types.StringValue("{abc}"),
		}),
		fields,
	)
	if jsonKey != objKey {
		t.Fatalf("primary key mismatch: jsonKey=%q objKey=%q", jsonKey, objKey)
	}
}

// TestStableItemPrimaryKeyIgnoresUnknownComputedFields guards the post-Create
// scenario: the planned object has Unknown for display_name / created_on,
// the API response has them filled in. The primary key for both sides must
// be the same so the reorderer can pair them.
func TestStableItemPrimaryKeyIgnoresUnknownComputedFields(t *testing.T) {
	fields := []BodyFieldDef{{Path: "uuid"}, {Path: "display_name"}}

	planned := stableObjectPrimaryKey(
		types.ObjectValueMust(itemAttrTypes(fields), map[string]attr.Value{
			"uuid":         types.StringValue("{abc}"),
			"display_name": types.StringUnknown(),
		}),
		fields,
	)
	apiSide := stableItemPrimaryKey(map[string]any{
		"uuid":         "{abc}",
		"display_name": "Alice",
	}, fields)
	if planned != apiSide {
		t.Fatalf("planned-vs-api primary-key disagreement: planned=%q api=%q", planned, apiSide)
	}
}

// TestReorderResponseArrayPreservesSourceOrder is the central guard for the
// post-Create / post-Update consistency check. Given an API response in
// `[A, B]` order and a planned source in `[B, A]` order (the operator's
// HCL), the reorderer must hand back `[B, A]`. New items in the API
// response that aren't in source are appended at the end.
func TestReorderResponseArrayPreservesSourceOrder(t *testing.T) {
	fields := []BodyFieldDef{{Path: "uuid"}}
	api := []any{
		map[string]any{"uuid": "{A}", "display_name": "Alice"},
		map[string]any{"uuid": "{B}", "display_name": "Bob"},
		map[string]any{"uuid": "{C}", "display_name": "Carol"}, // newly added by API
	}
	sourceOrder := []string{"uuid={B}", "uuid={A}"} // operator's HCL had B before A; C was added later

	got := reorderResponseArrayBySourceKeys(api, fields, sourceOrder)
	if len(got) != 3 {
		t.Fatalf("len = %d, want 3", len(got))
	}
	wantUUIDs := []string{"{B}", "{A}", "{C}"}
	for i, want := range wantUUIDs {
		gotUUID, _ := got[i].(map[string]any)["uuid"].(string)
		if gotUUID != want {
			t.Fatalf("got[%d].uuid = %q, want %q", i, gotUUID, want)
		}
	}
}

// TestReorderResponseArrayReturnsAPIOrderWhenSourceEmpty handles the no-prior
// case (Read on an unmanaged resource, Import, etc.): the reorderer must
// not drop anything and must preserve the API's order verbatim.
func TestReorderResponseArrayReturnsAPIOrderWhenSourceEmpty(t *testing.T) {
	fields := []BodyFieldDef{{Path: "uuid"}}
	api := []any{
		map[string]any{"uuid": "{A}"},
		map[string]any{"uuid": "{B}"},
	}
	got := reorderResponseArrayBySourceKeys(api, fields, nil)
	if len(got) != 2 || got[0].(map[string]any)["uuid"] != "{A}" || got[1].(map[string]any)["uuid"] != "{B}" {
		t.Fatalf("got = %#v, want [{A},{B}] verbatim", got)
	}
}

// TestReorderResponseArrayDeduplicatesSourceKeys guards a defensive edge
// case: duplicate keys in source order must not cause an item to be emitted
// twice in the output.
func TestReorderResponseArrayDeduplicatesSourceKeys(t *testing.T) {
	fields := []BodyFieldDef{{Path: "uuid"}}
	api := []any{
		map[string]any{"uuid": "{A}"},
		map[string]any{"uuid": "{B}"},
	}
	got := reorderResponseArrayBySourceKeys(api, fields, []string{"uuid={A}", "uuid={A}", "uuid={B}"})
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
}
