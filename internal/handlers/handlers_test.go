package handlers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"

	"github.com/FabianSchurig/bitbucket-cli/internal/client"
	"github.com/FabianSchurig/bitbucket-cli/internal/generated"
	"github.com/FabianSchurig/bitbucket-cli/internal/handlers"
	"github.com/FabianSchurig/bitbucket-cli/internal/output"
)

// newTestClient creates a BBClient that points to the provided test server.
func newTestClient(t *testing.T, serverURL string) *client.BBClient {
	t.Helper()
	r := resty.New().SetBaseURL(serverURL).SetBasicAuth("u", "p")
	return &client.BBClient{Client: r}
}

func TestListPullRequests_SinglePage(t *testing.T) {
	output.Format = "json"

	id1 := 1
	title1 := "First PR"
	state := generated.PullrequestStateOPEN

	page := generated.PaginatedPullrequests{
		Values: &[]generated.Pullrequest{
			{Id: &id1, Title: &title1, State: &state},
		},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repositories/myorg/myrepo/pullrequests" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(page); err != nil {
			t.Errorf("encoding response: %v", err)
		}
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	err := handlers.ListPullRequests(context.Background(), c, handlers.ListPRsInput{
		Workspace: "myorg",
		RepoSlug:  "myrepo",
	})
	if err != nil {
		t.Fatalf("ListPullRequests: %v", err)
	}
}

func TestListPullRequests_Pagination(t *testing.T) {
	output.Format = "json"

	id1, id2 := 1, 2
	t1, t2 := "PR 1", "PR 2"
	state := generated.PullrequestStateOPEN

	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/repositories/myorg/myrepo/pullrequests":
			// First page — return a "next" cursor pointing to the absolute URL
			next := r.Host
			if next == "" {
				next = "localhost"
			}
			nextURL := "http://" + r.Host + "/page2"
			page := generated.PaginatedPullrequests{
				Values: &[]generated.Pullrequest{{Id: &id1, Title: &t1, State: &state}},
				Next:   &nextURL,
			}
			if err := json.NewEncoder(w).Encode(page); err != nil {
				t.Errorf("encoding page 1: %v", err)
			}
		case "/page2":
			// Second page — no "next"
			page := generated.PaginatedPullrequests{
				Values: &[]generated.Pullrequest{{Id: &id2, Title: &t2, State: &state}},
			}
			if err := json.NewEncoder(w).Encode(page); err != nil {
				t.Errorf("encoding page 2: %v", err)
			}
		default:
			http.NotFound(w, r)
		}
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	err := handlers.ListPullRequests(context.Background(), c, handlers.ListPRsInput{
		Workspace: "myorg",
		RepoSlug:  "myrepo",
		All:       true,
	})
	if err != nil {
		t.Fatalf("ListPullRequests: %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 pages fetched, got %d", callCount)
	}
}

func TestGetPullRequest(t *testing.T) {
	output.Format = "json"

	id := 42
	title := "My PR"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/repositories/myorg/myrepo/pullrequests/42"
		if r.URL.Path != expected {
			t.Errorf("expected path %s, got %s", expected, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		pr := generated.Pullrequest{Id: &id, Title: &title}
		if err := json.NewEncoder(w).Encode(pr); err != nil {
			t.Errorf("encoding response: %v", err)
		}
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	err := handlers.GetPullRequest(context.Background(), c, handlers.GetPRInput{
		Workspace:     "myorg",
		RepoSlug:      "myrepo",
		PullRequestID: 42,
	})
	if err != nil {
		t.Fatalf("GetPullRequest: %v", err)
	}
}

func TestCreatePullRequest(t *testing.T) {
	output.Format = "json"

	id := 1
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body generated.Pullrequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decoding request body: %v", err)
		}
		if body.Title == nil || *body.Title != "My Feature" {
			t.Errorf("expected title 'My Feature', got %v", body.Title)
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		resp := generated.Pullrequest{Id: &id, Title: body.Title}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("encoding response: %v", err)
		}
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	err := handlers.CreatePullRequest(context.Background(), c, handlers.CreatePRInput{
		Workspace:         "myorg",
		RepoSlug:          "myrepo",
		Title:             "My Feature",
		SourceBranch:      "feature/x",
		DestinationBranch: "main",
	})
	if err != nil {
		t.Fatalf("CreatePullRequest: %v", err)
	}
}

func TestMergePullRequest(t *testing.T) {
	output.Format = "json"

	id := 5
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/repositories/myorg/myrepo/pullrequests/5/merge"
		if r.URL.Path != expected {
			t.Errorf("expected path %s, got %s", expected, r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		state := generated.PullrequestStateMERGED
		resp := generated.Pullrequest{Id: &id, State: &state}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("encoding response: %v", err)
		}
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	err := handlers.MergePullRequest(context.Background(), c, handlers.MergePRInput{
		Workspace:     "myorg",
		RepoSlug:      "myrepo",
		PullRequestID: 5,
		Strategy:      "squash",
	})
	if err != nil {
		t.Fatalf("MergePullRequest: %v", err)
	}
}

func TestListPullRequests_APIError(t *testing.T) {
	output.Format = "json"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		if _, err := w.Write([]byte(`{"error": {"message": "Unauthorized"}}`)); err != nil {
			t.Errorf("writing response: %v", err)
		}
	}))
	defer srv.Close()

	c := newTestClient(t, srv.URL)
	err := handlers.ListPullRequests(context.Background(), c, handlers.ListPRsInput{
		Workspace: "myorg",
		RepoSlug:  "myrepo",
	})
	if err == nil {
		t.Error("expected error for 401 response, got nil")
	}
}
