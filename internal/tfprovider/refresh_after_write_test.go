package tfprovider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TestRefreshAfterWrite_ReadDiffersTriggersGET asserts that any write
// operation whose Read counterpart differs (a generic property, not just
// project branch-restrictions) triggers a follow-up Read against the
// freshly-written state. This is the mechanism by which response-shape
// transformers — registered for the Read op — get a chance to populate
// state correctly when the write response shape diverges from the Read
// schema.
func TestRefreshAfterWrite_ReadDiffersTriggersGET(t *testing.T) {
	ctx := context.Background()

	var gets int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/items/ws/5" {
			atomic.AddInt32(&gets, 1)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"title": "Hello"})
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	group := testResourceGroup()
	r := &GenericResource{group: group, client: testBBClient(srv.URL)}

	// Both Create and Update are different ops from Read in the test group;
	// each must drive exactly one follow-up GET.
	for _, writeOp := range []*OperationDef{group.Ops.Create, group.Ops.Update} {
		atomic.StoreInt32(&gets, 0)
		state := newMockState(map[string]attr.Value{
			"workspace": types.StringValue("ws"),
			"param_id":  types.StringValue("5"),
			"id":        types.StringValue("priorID"),
		})
		var diags diag.Diagnostics
		r.refreshAfterWrite(ctx, writeOp, state, nil, &diags)
		if diags.HasError() {
			t.Fatalf("%s: unexpected diagnostics: %#v", writeOp.OperationID, diags)
		}
		if got := atomic.LoadInt32(&gets); got != 1 {
			t.Fatalf("%s: expected 1 follow-up GET, got %d", writeOp.OperationID, got)
		}
		// Prior id must be preserved across the refresh — the Read op's
		// fallback id (which would point at the wrong endpoint) must not
		// overwrite the canonical id written by the write op.
		if got := state.set["id"]; got != types.StringValue("priorID") {
			t.Fatalf("%s: refresh must preserve prior id, got %#v", writeOp.OperationID, got)
		}
	}
}

// TestRefreshAfterWrite_SameOpSkipsRefresh asserts that when the write op
// and the Read op share an OperationID (e.g. via CRUDConfig aliasing where
// a single endpoint is wired up for both roles), no redundant follow-up
// Read is performed. This avoids re-issuing the same call right after it
// just ran.
func TestRefreshAfterWrite_SameOpSkipsRefresh(t *testing.T) {
	ctx := context.Background()

	var hits int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&hits, 1)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{}`))
	}))
	defer srv.Close()

	group := testResourceGroup()
	// Force Read and Update to share an OperationID to exercise the skip
	// branch — a generic configuration the helper must handle.
	group.Ops.Read.OperationID = group.Ops.Update.OperationID
	r := &GenericResource{group: group, client: testBBClient(srv.URL)}

	state := newMockState(map[string]attr.Value{
		"workspace": types.StringValue("ws"),
		"param_id":  types.StringValue("5"),
	})
	var diags diag.Diagnostics
	r.refreshAfterWrite(ctx, group.Ops.Update, state, nil, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %#v", diags)
	}
	if got := atomic.LoadInt32(&hits); got != 0 {
		t.Fatalf("expected no follow-up call when read op == write op, got %d", got)
	}
}

// TestRefreshAfterWrite_NoReadOpIsNoOp guards the no-Read-op edge case
// (some resource groups expose only a write op). The helper must not
// panic and must not issue any HTTP calls.
func TestRefreshAfterWrite_NoReadOpIsNoOp(t *testing.T) {
	ctx := context.Background()
	srv := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		t.Fatal("no HTTP call expected when Read op is nil")
	}))
	defer srv.Close()

	group := testResourceGroup()
	group.Ops.Read = nil
	r := &GenericResource{group: group, client: testBBClient(srv.URL)}

	var diags diag.Diagnostics
	r.refreshAfterWrite(ctx, group.Ops.Create, newMockState(nil), nil, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %#v", diags)
	}
}

// TestRefreshAfterWrite_UsesParamFallback covers the Update case where a
// required Read path param is Computed-only and therefore unknown in the
// freshly-written state (e.g. a numeric "id" that was "(known after apply)"
// in the plan and only surfaces in the prior state). Without the fallback
// the post-write Read would fail with "Missing Required Parameter"; with
// it, the dispatcher consults the prior state and the refresh succeeds.
func TestRefreshAfterWrite_UsesParamFallback(t *testing.T) {
	ctx := context.Background()

	var gets int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/items/ws/5" {
			atomic.AddInt32(&gets, 1)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{"title": "Hello"})
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	group := testResourceGroup()
	r := &GenericResource{group: group, client: testBBClient(srv.URL)}

	// Freshly-written Update state: workspace is known, param_id is unknown
	// (Computed-only, was "(known after apply)" in the plan), id was set by
	// the write op.
	written := newMockState(map[string]attr.Value{
		"workspace": types.StringValue("ws"),
		"param_id":  types.StringUnknown(),
		"id":        types.StringValue("priorID"),
	})
	// Prior state carries the previously-assigned param_id that the post-
	// write Read needs in order to construct the GET URL.
	prior := newMockState(map[string]attr.Value{
		"workspace": types.StringValue("ws"),
		"param_id":  types.StringValue("5"),
		"id":        types.StringValue("priorID"),
	})

	var diags diag.Diagnostics
	r.refreshAfterWrite(ctx, group.Ops.Update, written, prior, &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %#v", diags)
	}
	if got := atomic.LoadInt32(&gets); got != 1 {
		t.Fatalf("expected 1 follow-up GET when param_id comes from prior state, got %d", got)
	}
	if got := written.set["id"]; got != types.StringValue("priorID") {
		t.Fatalf("refresh must preserve prior id when param_id is only available from prior state, got %#v", got)
	}
}

// TestRefreshAfterWrite_RetriesOn404 reproduces the eventual-consistency
// failure mode that surfaced on Bitbucket Cloud when a freshly created
// project/repository is read back immediately: the POST returns 2xx but
// the next GET for the same resource transiently returns 404 ("you may
// not have access ... or it no longer exists"). The retry loop in
// refreshState must convert that brief window into a successful refresh
// instead of failing the apply.
func TestRefreshAfterWrite_RetriesOn404(t *testing.T) {
	ctx := context.Background()

	var gets int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/items/ws/5" {
			http.NotFound(w, r)
			return
		}
		n := atomic.AddInt32(&gets, 1)
		// First two GETs return 404 (simulating eventual consistency);
		// the third succeeds. The retry budget must absorb this.
		if n < 3 {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"type":"error","error":{"message":"Not found"}}`))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"title": "Hello"})
	}))
	defer srv.Close()

	group := testResourceGroup()
	r := &GenericResource{group: group, client: testBBClient(srv.URL)}

	state := newMockState(map[string]attr.Value{
		"workspace": types.StringValue("ws"),
		"param_id":  types.StringValue("5"),
		"id":        types.StringValue("priorID"),
	})
	var diags diag.Diagnostics
	r.refreshAfterWrite(ctx, group.Ops.Create, state, nil, &diags)
	if diags.HasError() {
		t.Fatalf("expected refresh to succeed after retrying 404s, got: %#v", diags)
	}
	if got := atomic.LoadInt32(&gets); got != 3 {
		t.Fatalf("expected exactly 3 GETs (2 retried 404s + 1 success), got %d", got)
	}
	if got := state.set["id"]; got != types.StringValue("priorID") {
		t.Fatalf("refresh must preserve prior id across retries, got %#v", got)
	}
}

// TestRefreshAfterWrite_PersistentlyFails404 asserts that a resource that
// truly does not exist after the full retry budget surfaces the 404 to the
// user rather than silently succeeding. The retry exists only to smooth
// over eventual consistency; it must not hide real "missing resource"
// errors indefinitely.
func TestRefreshAfterWrite_PersistentlyFails404(t *testing.T) {
	ctx := context.Background()

	var gets int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&gets, 1)
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"type":"error","error":{"message":"Not found"}}`))
	}))
	defer srv.Close()

	group := testResourceGroup()
	r := &GenericResource{group: group, client: testBBClient(srv.URL)}

	state := newMockState(map[string]attr.Value{
		"workspace": types.StringValue("ws"),
		"param_id":  types.StringValue("5"),
		"id":        types.StringValue("priorID"),
	})
	var diags diag.Diagnostics
	r.refreshAfterWrite(ctx, group.Ops.Create, state, nil, &diags)
	if !diags.HasError() {
		t.Fatalf("expected diagnostics to carry the final 404, got none")
	}
	if got := atomic.LoadInt32(&gets); got != int32(readAfterWriteMaxAttempts) {
		t.Fatalf("expected exactly %d GETs after exhausting retries, got %d",
			readAfterWriteMaxAttempts, got)
	}
}

// TestRefreshAfterWrite_DoesNotRetryNon404 guards the scoping of the retry
// loop: only 404 is treated as eventual-consistency noise. Other API
// errors (400, 401, 403, 5xx) must propagate after a single attempt so
// genuine misconfiguration or auth problems are surfaced promptly instead
// of being hidden behind a multi-second backoff.
func TestRefreshAfterWrite_DoesNotRetryNon404(t *testing.T) {
	ctx := context.Background()

	var gets int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		atomic.AddInt32(&gets, 1)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"type":"error","error":{"message":"Bad request"}}`))
	}))
	defer srv.Close()

	group := testResourceGroup()
	r := &GenericResource{group: group, client: testBBClient(srv.URL)}

	state := newMockState(map[string]attr.Value{
		"workspace": types.StringValue("ws"),
		"param_id":  types.StringValue("5"),
		"id":        types.StringValue("priorID"),
	})
	var diags diag.Diagnostics
	r.refreshAfterWrite(ctx, group.Ops.Create, state, nil, &diags)
	if !diags.HasError() {
		t.Fatalf("expected diagnostics carrying the 400 error")
	}
	if got := atomic.LoadInt32(&gets); got != 1 {
		t.Fatalf("non-404 errors must not be retried; expected 1 GET, got %d", got)
	}
}
