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

const (
	tnByPattern    = "project-branch-restrictions-by-pattern"
	tnByBranchType = "project-branch-restrictions-by-branch-type"
)

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
		tnByPattern, newMockState(nil), in, &diag.Diagnostics{})
	if !reflect.DeepEqual(out, in) {
		t.Fatalf("unexpected mutation for non-target op: got %#v", out)
	}
}

func TestTransformProjectBranchRestrictionsRead_NilResult(t *testing.T) {
	out := transformProjectBranchRestrictionsRead(context.Background(),
		readOpGroupByBranch, tnByPattern, sourceWithPattern("*"), nil, &diag.Diagnostics{})
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
		readOpGroupByBranch, tnByPattern, sourceWithPattern("*"), resp, &diag.Diagnostics{})

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

func TestTransformProjectBranchRestrictionsRead_MultipleKindsSortedDeterministically(t *testing.T) {
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

	// Run the transform many times so that any nondeterministic Go map
	// iteration would surface as an out-of-order result on at least one run
	// (Go intentionally randomises map iteration to discourage relying on it).
	wantKinds := []string{"delete", "force", "push"}
	for i := 0; i < 50; i++ {
		out := transformProjectBranchRestrictionsRead(context.Background(),
			readOpGroupByBranch, tnByPattern, sourceWithPattern("main"), resp, &diag.Diagnostics{})

		values := out.(map[string]any)["values"].([]any)
		if len(values) != 3 {
			t.Fatalf("iteration %d: expected one row per kind, got %d", i, len(values))
		}
		gotKinds := make([]string, 0, len(values))
		for _, v := range values {
			row := v.(map[string]any)
			gotKinds = append(gotKinds, row["kind"].(string))
			if row["pattern"] != "main" {
				t.Errorf("iteration %d: pattern not propagated: %#v", i, row)
			}
		}
		if !reflect.DeepEqual(gotKinds, wantKinds) {
			t.Fatalf("iteration %d: expected deterministic sorted kind order %v, got %v", i, wantKinds, gotKinds)
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
		readOpGroupByBranch, tnByBranchType, sourceWithBranchType("production"), resp, &diag.Diagnostics{})

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
		readOpGroupByBranch, tnByPattern, sourceWithPattern("*"), resp, &diag.Diagnostics{})

	values := out.(map[string]any)["values"].([]any)
	if len(values) != 1 {
		t.Fatalf("expected 1 row, got %d", len(values))
	}
	row := values[0].(map[string]any)
	if v, ok := row["value"].(float64); !ok || v != 3 {
		t.Errorf("value: want 3, got %#v", row["value"])
	}
}

func TestTransformProjectBranchRestrictionsRead_ObjectShapedResponseSortedDeterministically(t *testing.T) {
	// The schema declares the response as an object whose values are arrays of
	// rules; ensure the transformer can handle that shape and that branch keys
	// are visited in a stable order.
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

	// Filter to "main"; should always yield the same single row regardless of
	// underlying map iteration order.
	for i := 0; i < 25; i++ {
		out := transformProjectBranchRestrictionsRead(context.Background(),
			readOpGroupByBranch, tnByPattern, sourceWithPattern("main"), resp, &diag.Diagnostics{})

		values := out.(map[string]any)["values"].([]any)
		if len(values) != 1 {
			t.Fatalf("iteration %d: expected 1 row filtered to main, got %d", i, len(values))
		}
		row := values[0].(map[string]any)
		if row["pattern"] != "main" || row["kind"] != "push" {
			t.Errorf("iteration %d: unexpected row: %#v", i, row)
		}
	}

	// And without a scope, both branch keys should be returned in sorted order.
	out := transformProjectBranchRestrictionsRead(context.Background(),
		readOpGroupByBranch, "unknown-typename", newMockState(nil), resp, &diag.Diagnostics{})
	values := out.(map[string]any)["values"].([]any)
	if len(values) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(values))
	}
	if values[0].(map[string]any)["pattern"] != "main" {
		t.Errorf("expected first row to be 'main' (sorted), got %v", values[0])
	}
	if values[1].(map[string]any)["pattern"] != "release/*" {
		t.Errorf("expected second row to be 'release/*' (sorted), got %v", values[1])
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
	// Use a typename that the gating logic does not recognise — both reads
	// should be skipped, leaving pattern/branch_type empty and matching all.
	out := transformProjectBranchRestrictionsRead(context.Background(),
		readOpGroupByBranch, "unknown-typename", newMockState(nil), resp, &diag.Diagnostics{})
	values := out.(map[string]any)["values"].([]any)
	if len(values) != 1 {
		t.Fatalf("expected entry to be kept when no scope is set, got %d", len(values))
	}
}

func TestTransformProjectBranchRestrictionsRead_TypeNameGatesAttributeReads(t *testing.T) {
	// When the active sub-resource is by-pattern, the transform must not call
	// GetAttribute on `branch_type` (which is not part of that sub-resource's
	// schema and would yield a diagnostic in the real plugin framework).
	// Mirror behaviour for by-branch-type and `pattern`.
	captured := newMockState(map[string]attr.Value{
		"pattern":     types.StringValue("main"),
		"branch_type": types.StringValue("production"),
	})

	resp := []any{
		map[string]any{
			"kind": map[string]any{
				"push": map[string]any{"users": []any{}, "groups": []any{}},
			},
			"branch_match_kind": "glob",
			"pattern":           "main",
			"branch_type":       "",
		},
	}

	// by-pattern: should match when only pattern is read. If branch_type were
	// also read it would be "production" and the entry (branch_type "") would
	// be filtered out.
	out := transformProjectBranchRestrictionsRead(context.Background(),
		readOpGroupByBranch, tnByPattern, captured, resp, &diag.Diagnostics{})
	if vs := out.(map[string]any)["values"].([]any); len(vs) != 1 {
		t.Fatalf("by-pattern: expected entry to match, got %d (branch_type was read by mistake?)", len(vs))
	}

	// by-branch-type: response entry has branch_type "" which doesn't match
	// "production"; verify pattern is NOT read (otherwise we'd still match).
	respBT := []any{
		map[string]any{
			"kind": map[string]any{
				"push": map[string]any{"users": []any{}, "groups": []any{}},
			},
			"branch_match_kind": "branching_model",
			"pattern":           "",
			"branch_type":       "production",
		},
	}
	out = transformProjectBranchRestrictionsRead(context.Background(),
		readOpGroupByBranch, tnByBranchType, captured, respBT, &diag.Diagnostics{})
	if vs := out.(map[string]any)["values"].([]any); len(vs) != 1 {
		t.Fatalf("by-branch-type: expected entry to match production, got %d (pattern was read by mistake?)", len(vs))
	}
}
