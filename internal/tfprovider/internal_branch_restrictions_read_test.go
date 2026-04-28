package tfprovider

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// readOpGroupByBranch is the operation definition under test. We only need the
// OperationID — the transform key — for this unit test.
var readOpGroupByBranch = &OperationDef{
	OperationID: "getProjectBranchRestrictionsGroupedByBranch",
}

func sourceWithPattern(pattern string) *mockState {
	return newMockState(map[string]attr.Value{
		"pattern": types.StringValue(pattern),
	})
}

func sourceWithBranchType(bt string) *mockState {
	return newMockState(map[string]attr.Value{
		"branch_type": types.StringValue(bt),
	})
}

func TestTransformProjectBranchRestrictionsRead_NotTargetOp(t *testing.T) {
	in := []any{map[string]any{"foo": "bar"}}
	out := transformProjectBranchRestrictionsRead(context.Background(),
		&OperationDef{OperationID: "somethingElse"},
		newMockState(nil), in, &diag.Diagnostics{})
	if !reflect.DeepEqual(out, in) {
		t.Fatalf("unexpected mutation for non-target op: got %#v", out)
	}
}

func TestTransformProjectBranchRestrictionsRead_NilResult(t *testing.T) {
	out := transformProjectBranchRestrictionsRead(context.Background(),
		readOpGroupByBranch, sourceWithPattern("*"), nil, &diag.Diagnostics{})
	m, ok := out.(map[string]any)
	if !ok {
		t.Fatalf("expected map result, got %T", out)
	}
	v, _ := m["values"].([]any)
	if len(v) != 0 {
		t.Fatalf("expected empty values list, got %#v", m["values"])
	}
}

func TestTransformProjectBranchRestrictionsRead_ArrayResponseByPattern(t *testing.T) {
	resp := []any{
		map[string]any{
			"kind": map[string]any{
				"push": map[string]any{
					"users": []any{
						map[string]any{"uuid": "{u-1}", "display_name": "Alice"},
					},
					"groups": []any{},
				},
			},
			"branch_match_kind":    "glob",
			"pattern":              "*",
			"branch_type":          "",
			"entity_type":          "project",
			"overlapping_patterns": []any{},
		},
		// A second entry that should be filtered out by pattern.
		map[string]any{
			"kind": map[string]any{
				"force": map[string]any{"users": []any{}, "groups": []any{}},
			},
			"branch_match_kind": "glob",
			"pattern":           "release/*",
			"branch_type":       "",
		},
	}

	out := transformProjectBranchRestrictionsRead(context.Background(),
		readOpGroupByBranch, sourceWithPattern("*"), resp, &diag.Diagnostics{})

	m := out.(map[string]any)
	values := m["values"].([]any)
	if len(values) != 1 {
		t.Fatalf("expected 1 value (filtered by pattern), got %d: %#v", len(values), values)
	}
	row := values[0].(map[string]any)
	if row["kind"] != "push" {
		t.Errorf("kind: want push, got %v", row["kind"])
	}
	if row["branch_match_kind"] != "glob" {
		t.Errorf("branch_match_kind: want glob, got %v", row["branch_match_kind"])
	}
	if row["pattern"] != "*" {
		t.Errorf("pattern: want *, got %v", row["pattern"])
	}
	users, ok := row["users"].([]any)
	if !ok || len(users) != 1 {
		t.Fatalf("users not flattened up to row: %#v", row["users"])
	}
	if u := users[0].(map[string]any); u["uuid"] != "{u-1}" {
		t.Errorf("user uuid: want {u-1}, got %v", u["uuid"])
	}
	groups, ok := row["groups"].([]any)
	if !ok || len(groups) != 0 {
		t.Errorf("groups should be empty list, got %#v", row["groups"])
	}
}

func TestTransformProjectBranchRestrictionsRead_MultipleKinds(t *testing.T) {
	resp := []any{
		map[string]any{
			"kind": map[string]any{
				"push":   map[string]any{"users": []any{}, "groups": []any{}},
				"delete": map[string]any{"users": []any{}, "groups": []any{}},
				"force":  map[string]any{"users": []any{}, "groups": []any{}},
			},
			"branch_match_kind": "glob",
			"pattern":           "main",
			"branch_type":       "",
		},
	}

	out := transformProjectBranchRestrictionsRead(context.Background(),
		readOpGroupByBranch, sourceWithPattern("main"), resp, &diag.Diagnostics{})

	values := out.(map[string]any)["values"].([]any)
	if len(values) != 3 {
		t.Fatalf("expected one row per kind, got %d", len(values))
	}
	gotKinds := map[string]bool{}
	for _, v := range values {
		row := v.(map[string]any)
		gotKinds[row["kind"].(string)] = true
		if row["pattern"] != "main" {
			t.Errorf("pattern not propagated: %#v", row)
		}
	}
	for _, want := range []string{"push", "delete", "force"} {
		if !gotKinds[want] {
			t.Errorf("missing expanded kind %q", want)
		}
	}
}

func TestTransformProjectBranchRestrictionsRead_ByBranchType(t *testing.T) {
	resp := []any{
		map[string]any{
			"kind": map[string]any{
				"require_approvals_to_merge": map[string]any{"value": float64(2)},
			},
			"branch_match_kind": "branching_model",
			"pattern":           "",
			"branch_type":       "production",
		},
		map[string]any{
			"kind": map[string]any{
				"push": map[string]any{"users": []any{}, "groups": []any{}},
			},
			"branch_match_kind": "branching_model",
			"pattern":           "",
			"branch_type":       "development",
		},
	}

	out := transformProjectBranchRestrictionsRead(context.Background(),
		readOpGroupByBranch, sourceWithBranchType("production"), resp, &diag.Diagnostics{})

	values := out.(map[string]any)["values"].([]any)
	if len(values) != 1 {
		t.Fatalf("expected 1 value filtered to production, got %d", len(values))
	}
	row := values[0].(map[string]any)
	if row["kind"] != "require_approvals_to_merge" {
		t.Errorf("kind mismatch: %v", row["kind"])
	}
	if row["branch_type"] != "production" {
		t.Errorf("branch_type mismatch: %v", row["branch_type"])
	}
	if v, ok := row["value"].(float64); !ok || v != 2 {
		t.Errorf("value: want 2, got %#v", row["value"])
	}
}

func TestTransformProjectBranchRestrictionsRead_BareNumericKindData(t *testing.T) {
	// Some kinds may be encoded with a bare numeric payload rather than an
	// object; the transformer should still surface it as `value`.
	resp := []any{
		map[string]any{
			"kind": map[string]any{
				"require_approvals_to_merge": float64(3),
			},
			"branch_match_kind": "glob",
			"pattern":           "*",
		},
	}

	out := transformProjectBranchRestrictionsRead(context.Background(),
		readOpGroupByBranch, sourceWithPattern("*"), resp, &diag.Diagnostics{})

	values := out.(map[string]any)["values"].([]any)
	if len(values) != 1 {
		t.Fatalf("expected 1 row, got %d", len(values))
	}
	row := values[0].(map[string]any)
	if v, ok := row["value"].(float64); !ok || v != 3 {
		t.Errorf("value: want 3, got %#v", row["value"])
	}
}

func TestTransformProjectBranchRestrictionsRead_ObjectShapedResponse(t *testing.T) {
	// The schema declares the response as an object whose values are arrays of
	// rules; ensure the transformer can handle that shape too.
	resp := map[string]any{
		"main": []any{
			map[string]any{
				"kind": map[string]any{
					"push": map[string]any{"users": []any{}, "groups": []any{}},
				},
				"branch_match_kind": "glob",
				"pattern":           "main",
			},
		},
		"release/*": []any{
			map[string]any{
				"kind": map[string]any{
					"push": map[string]any{"users": []any{}, "groups": []any{}},
				},
				"branch_match_kind": "glob",
				"pattern":           "release/*",
			},
		},
	}

	out := transformProjectBranchRestrictionsRead(context.Background(),
		readOpGroupByBranch, sourceWithPattern("main"), resp, &diag.Diagnostics{})

	values := out.(map[string]any)["values"].([]any)
	if len(values) != 1 {
		t.Fatalf("expected 1 row filtered to main, got %d", len(values))
	}
	row := values[0].(map[string]any)
	if row["pattern"] != "main" || row["kind"] != "push" {
		t.Errorf("unexpected row: %#v", row)
	}
}

func TestTransformProjectBranchRestrictionsRead_NoScopeMatchesAll(t *testing.T) {
	// When neither pattern nor branch_type is in source state (e.g. fresh
	// import before any attributes are set), the transformer should not drop
	// every entry.
	resp := []any{
		map[string]any{
			"kind": map[string]any{
				"push": map[string]any{"users": []any{}, "groups": []any{}},
			},
			"branch_match_kind": "glob",
			"pattern":           "*",
		},
	}
	out := transformProjectBranchRestrictionsRead(context.Background(),
		readOpGroupByBranch, newMockState(nil), resp, &diag.Diagnostics{})
	values := out.(map[string]any)["values"].([]any)
	if len(values) != 1 {
		t.Fatalf("expected entry to be kept when no scope is set, got %d", len(values))
	}
}
