package tfprovider_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// reproRepoServer is a stateful mock that stores the repository's
// description and reflects updates, mimicking Bitbucket's GET/PUT contract
// closely enough to exercise plan stability.
func reproRepoServer(t *testing.T) *httptest.Server {
	t.Helper()
	var mu sync.Mutex
	desc := "A test repository"

	respBody := func() map[string]any {
		return map[string]any{
			"uuid":        "{repo-uuid-123}",
			"slug":        "test-repo",
			"name":        "test-repo",
			"full_name":   "testworkspace/test-repo",
			"description": desc,
			"is_private":  true,
			"scm":         "git",
			"project": map[string]any{
				"key":  "TINF",
				"type": "project",
				"uuid": "{proj-uuid}",
				"name": "Infra",
			},
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/repositories/{workspace}/{repo_slug}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		mu.Lock()
		defer mu.Unlock()
		switch r.Method {
		case http.MethodPost, http.MethodPut:
			b, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("read request body: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			var body map[string]any
			if err := json.Unmarshal(b, &body); err != nil {
				t.Errorf("unmarshal request body: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if d, ok := body["description"].(string); ok {
				desc = d
			}
			if r.Method == http.MethodPost {
				w.WriteHeader(http.StatusCreated)
			}
			_ = json.NewEncoder(w).Encode(respBody())
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(respBody())
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv
}

// TestAccReproIssue111_RequestBodyModify reproduces
// https://github.com/FabianSchurig/bitbucket-cli/issues/111 — modifying a
// repository via request_body after it already exists produced
// "Provider produced invalid plan ... request_body: planned value
// cty.NullVal(cty.String) does not match config value".
func TestAccReproIssue111_RequestBodyModify(t *testing.T) {
	srv := reproRepoServer(t)
	setMockEnv(t, srv.URL)

	base := func(extra string) string {
		return fmt.Sprintf(`
			provider "bitbucket" {
				base_url = %q
			}

			resource "bitbucket_repos" "test" {
				workspace = "testworkspace"
				repo_slug = "test-repo"
				%s
			}
		`, srv.URL, extra)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// Create with a plain description attribute.
			{
				Config: base(`description = "A test repository"`),
				Check: resource.TestCheckResourceAttr(
					"bitbucket_repos.test", "description", "A test repository"),
			},
			// Now drive the update through request_body — the issue scenario.
			{
				Config: base(`request_body = jsonencode({ description = "Changed via body." })`),
				Check: resource.TestCheckResourceAttr(
					"bitbucket_repos.test", "request_body", `{"description":"Changed via body."}`),
			},
		},
	})
}
