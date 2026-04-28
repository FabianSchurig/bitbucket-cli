package tfprovider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestReadPriorIDAndRestore covers the helpers introduced to keep the
// resource id stable across Read.
func TestReadPriorIDAndRestore(t *testing.T) {
	ctx := context.Background()

	// Empty / null state → empty prior id, restore is a no-op.
	empty := newMockState(nil)
	var diags diag.Diagnostics
	if got := readPriorID(ctx, empty, &diags); got != "" {
		t.Fatalf("expected empty priorID for empty state, got %q", got)
	}
	if diags.HasError() {
		t.Fatalf("readPriorID on empty state should not error, got %#v", diags)
	}
	restorePriorID(ctx, empty, "", &diags)
	if _, ok := empty.set["id"]; ok {
		t.Fatalf("restorePriorID with empty id must not write state, set=%#v", empty.set)
	}

	// Populated state → priorID is read, restorePriorID writes it back.
	state := newMockState(map[string]attr.Value{
		"id": types.StringValue("replaceProjectBranchRestrictionsByPattern/ws/PROJ/*"),
	})
	got := readPriorID(ctx, state, &diags)
	if got != "replaceProjectBranchRestrictionsByPattern/ws/PROJ/*" {
		t.Fatalf("unexpected priorID: %q", got)
	}
	target := newMockState(nil)
	restorePriorID(ctx, target, got, &diags)
	if target.set["id"] != types.StringValue(got) {
		t.Fatalf("expected restorePriorID to write id to target, got %#v", target.set["id"])
	}
}

// TestReadPriorIDSurfacesDiagnostics confirms that GetAttribute diagnostics
// are appended to the caller's diag bag rather than being silently dropped.
func TestReadPriorIDSurfacesDiagnostics(t *testing.T) {
	ctx := context.Background()
	state := newMockState(nil)
	state.diags = map[string]diag.Diagnostics{
		"id": {diag.NewErrorDiagnostic("schema mismatch", "id attribute not declared")},
	}

	var diags diag.Diagnostics
	if got := readPriorID(ctx, state, &diags); got != "" {
		t.Fatalf("expected empty priorID on framework error, got %q", got)
	}
	if !diags.HasError() {
		t.Fatal("expected GetAttribute diagnostics to be surfaced")
	}
}

// TestRefreshStatePreservesID exercises Bug 2 from the issue: a Read-style
// dispatch (e.g. project branch restrictions' group-by-branch GET) used to
// overwrite the resource id with the GET operation's path, breaking the
// subsequent Delete which then targets the wrong endpoint. refreshState must
// preserve the prior id written by Create.
func TestRefreshStatePreservesID(t *testing.T) {
	ctx := context.Background()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/items/ws/5" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"title": "Hello"})
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	group := testResourceGroup()
	r := &GenericResource{group: group, client: testBBClient(srv.URL)}

	// Simulate post-Create state: id was set by Create's storeDispatchResult.
	priorID := "createSample/ws"
	state := newMockState(map[string]attr.Value{
		"workspace": types.StringValue("ws"),
		"param_id":  types.StringValue("5"),
		"id":        types.StringValue(priorID),
	})

	var diags diag.Diagnostics
	r.refreshState(ctx, group.Ops.Read, state, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected refreshState diagnostics: %#v", diags)
	}
	if got := state.set["id"]; got != types.StringValue(priorID) {
		t.Fatalf("refreshState must preserve prior id, got %#v", got)
	}
}

// TestRefreshStateNoPriorIDDoesNotRestore confirms that when no prior id is
// in state (e.g. an import where the id has not yet been computed), the
// dispatch-written id is left in place rather than being cleared.
func TestRefreshStateNoPriorIDDoesNotRestore(t *testing.T) {
	ctx := context.Background()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"title": "Hello"})
	}))
	defer srv.Close()

	group := testResourceGroup()
	r := &GenericResource{group: group, client: testBBClient(srv.URL)}

	state := newMockState(map[string]attr.Value{
		"workspace": types.StringValue("ws"),
		"param_id":  types.StringValue("5"),
	})

	var diags diag.Diagnostics
	r.refreshState(ctx, group.Ops.Read, state, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected refreshState diagnostics: %#v", diags)
	}
	// dispatch's storeDispatchResult writes a fallback id from op + path
	// params. With no prior id to restore, that fallback must remain.
	if got, ok := state.set["id"].(types.String); !ok || got.ValueString() == "" {
		t.Fatalf("expected dispatch-written id to remain, got %#v", state.set["id"])
	}
}
