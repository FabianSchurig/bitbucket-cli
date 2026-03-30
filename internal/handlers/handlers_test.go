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
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"id": 42, "title": "My PR"})
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	err := handlers.Dispatch(context.Background(), c, handlers.Request{
		Method:      "GET",
		URLTemplate: "/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}",
		PathParams:  map[string]string{"workspace": "myorg", "repo_slug": "myrepo", "pull_request_id": "42"},
	})
	if err != nil {
		t.Fatalf("Dispatch: %v", err)
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
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"values": []any{
				map[string]any{"id": 1, "title": "First PR"},
			},
		})
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	err := handlers.Dispatch(context.Background(), c, handlers.Request{
		Method:      "GET",
		URLTemplate: "/repositories/{workspace}/{repo_slug}/pullrequests",
		PathParams:  map[string]string{"workspace": "myorg", "repo_slug": "myrepo"},
		QueryParams: map[string]string{"state": "OPEN"},
	})
	if err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
}

func TestDispatch_GET_Paginated_AllPages(t *testing.T) {
	output.Format = "json"

	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/repositories/myorg/myrepo/pullrequests":
			nextURL := "http://" + r.Host + "/page2"
			json.NewEncoder(w).Encode(map[string]any{
				"values": []any{map[string]any{"id": 1, "title": "PR 1"}},
				"next":   nextURL,
			})
		case "/page2":
			json.NewEncoder(w).Encode(map[string]any{
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
		URLTemplate: "/repositories/{workspace}/{repo_slug}/pullrequests",
		PathParams:  map[string]string{"workspace": "myorg", "repo_slug": "myrepo"},
		All:         true,
	})
	if err != nil {
		t.Fatalf("Dispatch: %v", err)
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
		json.NewDecoder(r.Body).Decode(&body)
		if body["title"] != "My Feature" {
			t.Errorf("expected title 'My Feature', got %v", body["title"])
		}
		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"id": 1, "title": "My Feature"})
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	err := handlers.Dispatch(context.Background(), c, handlers.Request{
		Method:      "POST",
		URLTemplate: "/repositories/{workspace}/{repo_slug}/pullrequests",
		PathParams:  map[string]string{"workspace": "myorg", "repo_slug": "myrepo"},
		Body:        `{"title":"My Feature","source":{"branch":{"name":"feature/x"}}}`,
	})
	if err != nil {
		t.Fatalf("Dispatch: %v", err)
	}
}

func TestDispatch_APIError(t *testing.T) {
	output.Format = "json"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": {"message": "Unauthorized"}}`))
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	err := handlers.Dispatch(context.Background(), c, handlers.Request{
		Method:      "GET",
		URLTemplate: "/repositories/{workspace}/{repo_slug}/pullrequests",
		PathParams:  map[string]string{"workspace": "myorg", "repo_slug": "myrepo"},
	})
	if err == nil {
		t.Error("expected error for 401 response, got nil")
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
		t.Fatalf("Dispatch: %v", err)
	}
}
