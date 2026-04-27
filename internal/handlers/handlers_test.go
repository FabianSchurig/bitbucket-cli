package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"

	"github.com/FabianSchurig/bitbucket-cli/internal/client"
	"github.com/FabianSchurig/bitbucket-cli/internal/handlers"
	"github.com/FabianSchurig/bitbucket-cli/internal/output"
)

const (
	headerContentType = "Content-Type"
	contentTypeJSON   = "application/json"
	dispatchErrFmt    = "Dispatch: %v"
	pullRequestsPath  = "/repositories/{workspace}/{repo_slug}/pullrequests"
)

// newTestClient creates a BBClient that points to the provided test server.
func newTestClient(t *testing.T, serverURL string) *client.BBClient {
	t.Helper()
	r := resty.New().SetBaseURL(serverURL).SetBasicAuth("u", "p")
	return &client.BBClient{Client: r}
}

func TestDispatch_GET_SingleResource(t *testing.T) {
	output.Format = "json"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/repositories/myorg/myrepo/pullrequests/42"
		if r.URL.Path != expected {
			t.Errorf("expected path %s, got %s", expected, r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set(headerContentType, contentTypeJSON)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": 42, "title": "My PR"})
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	err := handlers.Dispatch(context.Background(), c, handlers.Request{
		Method:      "GET",
		URLTemplate: "/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}",
		PathParams:  map[string]string{"workspace": "myorg", "repo_slug": "myrepo", "pull_request_id": "42"},
	})
	if err != nil {
		t.Fatalf(dispatchErrFmt, err)
	}
}

func TestDispatch_GET_Paginated_SinglePage(t *testing.T) {
	output.Format = "json"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repositories/myorg/myrepo/pullrequests" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("state") != "OPEN" {
			t.Errorf("expected state=OPEN query param, got %s", r.URL.Query().Get("state"))
		}
		w.Header().Set(headerContentType, contentTypeJSON)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"values": []any{
				map[string]any{"id": 1, "title": "First PR"},
			},
		})
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	err := handlers.Dispatch(context.Background(), c, handlers.Request{
		Method:      "GET",
		URLTemplate: pullRequestsPath,
		PathParams:  map[string]string{"workspace": "myorg", "repo_slug": "myrepo"},
		QueryParams: map[string]string{"state": "OPEN"},
	})
	if err != nil {
		t.Fatalf(dispatchErrFmt, err)
	}
}

func TestDispatch_GET_Paginated_AllPages(t *testing.T) {
	output.Format = "json"

	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set(headerContentType, contentTypeJSON)

		switch r.URL.Path {
		case "/repositories/myorg/myrepo/pullrequests":
			nextURL := "http://" + r.Host + "/page2"
			_ = json.NewEncoder(w).Encode(map[string]any{
				"values": []any{map[string]any{"id": 1, "title": "PR 1"}},
				"next":   nextURL,
			})
		case "/page2":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"values": []any{map[string]any{"id": 2, "title": "PR 2"}},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	err := handlers.Dispatch(context.Background(), c, handlers.Request{
		Method:      "GET",
		URLTemplate: pullRequestsPath,
		PathParams:  map[string]string{"workspace": "myorg", "repo_slug": "myrepo"},
		All:         true,
	})
	if err != nil {
		t.Fatalf(dispatchErrFmt, err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 pages fetched, got %d", callCount)
	}
}

func TestDispatch_POST_WithBody(t *testing.T) {
	output.Format = "json"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decoding body: %v", err)
		}
		if body["title"] != "My Feature" {
			t.Errorf("expected title 'My Feature', got %v", body["title"])
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set(headerContentType, contentTypeJSON)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": 1, "title": "My Feature"})
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	err := handlers.Dispatch(context.Background(), c, handlers.Request{
		Method:      "POST",
		URLTemplate: pullRequestsPath,
		PathParams:  map[string]string{"workspace": "myorg", "repo_slug": "myrepo"},
		Body:        `{"title":"My Feature","source":{"branch":{"name":"feature/x"}}}`,
	})
	if err != nil {
		t.Fatalf(dispatchErrFmt, err)
	}
}

func TestDispatch_APIError(t *testing.T) {
	output.Format = "json"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error": {"message": "Unauthorized"}}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	err := handlers.Dispatch(context.Background(), c, handlers.Request{
		Method:      "GET",
		URLTemplate: pullRequestsPath,
		PathParams:  map[string]string{"workspace": "myorg", "repo_slug": "myrepo"},
	})
	if err == nil {
		t.Error("expected error for 401 response, got nil")
	}
}

// TestDispatch_AbsoluteURLTemplate verifies that an absolute URL in
// URLTemplate is sent as-is, ignoring the client's BaseURL. This is required
// for Bitbucket's internal endpoints which live on a different host
// (https://bitbucket.org/!api/internal/...) than the public API base URL.
func TestDispatch_AbsoluteURLTemplate(t *testing.T) {
	output.Format = "json"

	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decoding body: %v", err)
		}
		if _, ok := body["values"]; !ok {
			t.Errorf("expected body with 'values' key, got %v", body)
		}
		w.Header().Set(headerContentType, contentTypeJSON)
		_ = json.NewEncoder(w).Encode(map[string]any{"ok": true})
	}))
	defer srv.Close()

	// Client points at a deliberately-wrong base URL; the absolute URLTemplate
	// must override it. Internal-API URLs require csrf/session cookies.
	c := newTestClient(t, "http://wrong-host.invalid")
	c.CSRFToken = "csrf-test"
	c.CloudSessionToken = "session-test"
	err := handlers.Dispatch(context.Background(), c, handlers.Request{
		Method:      "PUT",
		URLTemplate: srv.URL + "/!api/internal/workspaces/{workspace}/projects/{project_key}/branch-restrictions/by-pattern/{pattern}",
		PathParams: map[string]string{
			"workspace":   "myworkspace",
			"project_key": "MYPROJ",
			"pattern":     "main",
		},
		Body: `{"values":[{"pattern":"main","branch_match_kind":"glob","kind":"force"}]}`,
	})
	if err != nil {
		t.Fatalf(dispatchErrFmt, err)
	}
	wantPath := "/!api/internal/workspaces/myworkspace/projects/MYPROJ/branch-restrictions/by-pattern/main"
	if gotPath != wantPath {
		t.Errorf("expected path %s, got %s", wantPath, gotPath)
	}
}

func TestDispatch_DELETE_NoContent(t *testing.T) {
	output.Format = "json"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	err := handlers.Dispatch(context.Background(), c, handlers.Request{
		Method:      "DELETE",
		URLTemplate: "/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/approve",
		PathParams:  map[string]string{"workspace": "myorg", "repo_slug": "myrepo", "pull_request_id": "5"},
	})
	if err != nil {
		t.Fatalf(dispatchErrFmt, err)
	}
}
