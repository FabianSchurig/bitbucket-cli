package tfprovider_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/FabianSchurig/bitbucket-cli/internal/client"
	"github.com/FabianSchurig/bitbucket-cli/internal/handlers"
	"github.com/FabianSchurig/bitbucket-cli/internal/tfprovider"
)

// testAccProtoV6ProviderFactories creates provider factories for acceptance tests.
func testAccProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"bitbucket": providerserver.NewProtocol6WithError(tfprovider.New("test")()),
	}
}

// startMockServer starts a mock HTTP server simulating common Bitbucket API endpoints.
// It returns the server and its URL. The caller must defer srv.Close().
func startMockServer(t *testing.T) *httptest.Server {
	t.Helper()

	mux := http.NewServeMux()

	// Repository endpoints
	mux.HandleFunc("/repositories/{workspace}/{repo_slug}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"uuid":        "{repo-uuid-123}",
				"slug":        "test-repo",
				"name":        "test-repo",
				"full_name":   "testworkspace/test-repo",
				"description": "A test repository",
				"is_private":  true,
				"scm":         "git",
			})
		case http.MethodPost:
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"uuid":        "{repo-uuid-123}",
				"slug":        "test-repo",
				"name":        "test-repo",
				"full_name":   "testworkspace/test-repo",
				"description": "A test repository",
				"is_private":  true,
				"scm":         "git",
			})
		case http.MethodPut:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"uuid":        "{repo-uuid-123}",
				"slug":        "test-repo",
				"name":        "test-repo",
				"full_name":   "testworkspace/test-repo",
				"description": "Updated description",
				"is_private":  true,
				"scm":         "git",
			})
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})

	// Repository list endpoint (paginated)
	mux.HandleFunc("/repositories/{workspace}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"values": []any{
				map[string]any{
					"uuid":      "{repo-uuid-123}",
					"slug":      "test-repo",
					"name":      "test-repo",
					"full_name": "testworkspace/test-repo",
				},
			},
			"page": 1,
			"size": 1,
		})
	})

	// Project endpoints
	mux.HandleFunc("/workspaces/{workspace}/projects/{project_key}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"uuid":        "{project-uuid-123}",
				"key":         "TEST",
				"name":        "Test Project",
				"description": "A test project",
				"is_private":  true,
			})
		case http.MethodPut:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"uuid":        "{project-uuid-123}",
				"key":         "TEST",
				"name":        "Updated Project",
				"description": "Updated description",
				"is_private":  true,
			})
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})

	// Project create endpoint
	mux.HandleFunc("POST /workspaces/{workspace}/projects", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"uuid":        "{project-uuid-123}",
			"key":         "TEST",
			"name":        "Test Project",
			"description": "A test project",
			"is_private":  true,
		})
	})

	// Workspace endpoint
	mux.HandleFunc("/workspaces/{workspace}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"uuid":       "{workspace-uuid-123}",
			"slug":       "testworkspace",
			"name":       "Test Workspace",
			"is_private": false,
		})
	})

	// User endpoint
	mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"uuid":         "{user-uuid-123}",
			"username":     "testuser",
			"display_name": "Test User",
		})
	})

	// ─── Workspace webhook endpoints ──────────────────────────────────────────
	mux.HandleFunc("/workspaces/{workspace}/hooks/{uid}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"uuid":        "{hook-uuid-123}",
				"url":         "https://example.com/webhook",
				"description": "Test webhook",
				"active":      true,
			})
		case http.MethodPut:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"uuid":        "{hook-uuid-123}",
				"url":         "https://example.com/webhook-updated",
				"description": "Updated webhook",
				"active":      true,
			})
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})
	mux.HandleFunc("POST /workspaces/{workspace}/hooks", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"uuid":        "{hook-uuid-123}",
			"url":         "https://example.com/webhook",
			"description": "Test webhook",
			"active":      true,
		})
	})

	// ─── Default reviewer endpoints ───────────────────────────────────────────
	mux.HandleFunc("/repositories/{workspace}/{repo_slug}/default-reviewers/{target_username}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"uuid":         "{user-uuid-123}",
				"display_name": "Test User",
				"nickname":     "testuser",
			})
		case http.MethodPut:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"uuid":         "{user-uuid-123}",
				"display_name": "Test User",
				"nickname":     "testuser",
			})
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})

	// ─── Pipeline variable endpoints ──────────────────────────────────────────
	mux.HandleFunc("/repositories/{workspace}/{repo_slug}/pipelines_config/variables/{variable_uuid}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"uuid":    "{var-uuid-123}",
				"key":     "MY_VAR",
				"value":   "my-value",
				"secured": false,
			})
		case http.MethodPut:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"uuid":    "{var-uuid-123}",
				"key":     "MY_VAR",
				"value":   "updated-value",
				"secured": false,
			})
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})
	mux.HandleFunc("POST /repositories/{workspace}/{repo_slug}/pipelines_config/variables", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"uuid":    "{var-uuid-123}",
			"key":     "MY_VAR",
			"value":   "my-value",
			"secured": false,
		})
	})

	// ─── Repo deploy key endpoints ────────────────────────────────────────────
	mux.HandleFunc("/repositories/{workspace}/{repo_slug}/deploy-keys/{key_id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":      123,
				"key":     "ssh-rsa AAAA...",
				"label":   "test-key",
				"comment": "test@example.com",
			})
		case http.MethodPut:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":      123,
				"key":     "ssh-rsa AAAA...",
				"label":   "updated-key",
				"comment": "test@example.com",
			})
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})

	// ─── Repo explicit permissions endpoints ──────────────────────────────────
	mux.HandleFunc("/repositories/{workspace}/{repo_slug}/permissions-config/groups/{group_slug}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"permission": "read",
				"group": map[string]any{
					"slug": "developers",
					"name": "Developers",
				},
			})
		case http.MethodPut:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"permission": "write",
				"group": map[string]any{
					"slug": "developers",
					"name": "Developers",
				},
			})
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})
	mux.HandleFunc("/repositories/{workspace}/{repo_slug}/permissions-config/users/{selected_user_id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"permission": "read",
				"user": map[string]any{
					"uuid":         "{user-uuid-123}",
					"display_name": "Test User",
				},
			})
		case http.MethodPut:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"permission": "write",
				"user": map[string]any{
					"uuid":         "{user-uuid-123}",
					"display_name": "Test User",
				},
			})
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})

	// ─── Wave 2: mock endpoints for additional sub-resources ────────────────

	// Tags endpoints
	mux.HandleFunc("/repositories/{workspace}/{repo_slug}/refs/tags/{name}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{"name": "v1.0.0", "type": "tag"})
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})
	mux.HandleFunc("POST /repositories/{workspace}/{repo_slug}/refs/tags", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"name": "v1.0.0", "type": "tag"})
	})

	// Pipeline SSH keys endpoint
	mux.HandleFunc("/repositories/{workspace}/{repo_slug}/pipelines_config/ssh/key_pair", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{"uuid": "{ssh-key-uuid}", "public_key": "ssh-rsa AAAA..."})
		case http.MethodPut:
			_ = json.NewEncoder(w).Encode(map[string]any{"uuid": "{ssh-key-uuid}", "public_key": "ssh-rsa AAAA..."})
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})

	// Pipeline schedules endpoint
	mux.HandleFunc("/repositories/{workspace}/{repo_slug}/pipelines_config/schedules/{schedule_uuid}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{"uuid": "{schedule-uuid}", "enabled": true, "cron_pattern": "0 0 12 * * ? *", "target": map[string]any{}})
		case http.MethodPut:
			_ = json.NewEncoder(w).Encode(map[string]any{"uuid": "{schedule-uuid}", "enabled": true, "cron_pattern": "0 0 12 * * ? *", "target": map[string]any{}})
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})
	mux.HandleFunc("POST /repositories/{workspace}/{repo_slug}/pipelines_config/schedules", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"uuid": "{schedule-uuid}", "enabled": true, "cron_pattern": "0 0 12 * * ? *", "target": map[string]any{}})
	})

	// Pipeline known hosts endpoint
	mux.HandleFunc("/repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts/{known_host_uuid}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{"uuid": "{known-host-uuid}", "hostname": "bitbucket.org"})
		case http.MethodPut:
			_ = json.NewEncoder(w).Encode(map[string]any{"uuid": "{known-host-uuid}", "hostname": "bitbucket.org"})
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})
	mux.HandleFunc("POST /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"uuid": "{known-host-uuid}", "hostname": "bitbucket.org"})
	})

	// Pipeline config endpoint
	mux.HandleFunc("/repositories/{workspace}/{repo_slug}/pipelines_config", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"enabled": true, "type": "pipelines_config"})
	})

	// PR comments endpoint
	mux.HandleFunc("/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/comments/{comment_id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{"id": 1, "content": map[string]any{"raw": "test comment"}})
		case http.MethodPut:
			_ = json.NewEncoder(w).Encode(map[string]any{"id": 1, "content": map[string]any{"raw": "updated comment"}})
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})
	mux.HandleFunc("POST /repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/comments", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": 1, "content": map[string]any{"raw": "new comment"}})
	})

	// Issue comments endpoint
	mux.HandleFunc("/repositories/{workspace}/{repo_slug}/issues/{issue_id}/comments/{comment_id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{"id": 1, "content": map[string]any{"raw": "test issue comment"}})
		case http.MethodPut:
			_ = json.NewEncoder(w).Encode(map[string]any{"id": 1, "content": map[string]any{"raw": "updated issue comment"}})
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})
	mux.HandleFunc("POST /repositories/{workspace}/{repo_slug}/issues/{issue_id}/comments", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{"id": 1, "content": map[string]any{"raw": "new issue comment"}})
	})

	// ─── Branch restrictions endpoints ────────────────────────────────────────
	// The GET / PUT handlers always return `users` in lexicographic UUID
	// order regardless of the order the operator wrote them in HCL — this
	// reproduces the exact upstream Bitbucket behaviour the order-
	// insensitivity fix targets. Two distinct fixed UUIDs are used so the
	// test can write them in reverse order and assert that the plan still
	// matches and a follow-up plan is empty.
	branchRestrictionResponse := map[string]any{
		"id":                123,
		"kind":              "push",
		"branch_match_kind": "glob",
		"pattern":           "develop",
		"users": []any{
			map[string]any{"uuid": "{aaaaaaaa-0000-0000-0000-000000000001}", "display_name": "Alice", "created_on": "2024-01-01T00:00:00Z"},
			map[string]any{"uuid": "{bbbbbbbb-0000-0000-0000-000000000002}", "display_name": "Bob", "created_on": "2024-01-02T00:00:00Z"},
		},
		"groups": []any{},
	}
	mux.HandleFunc("/repositories/{workspace}/{repo_slug}/branch-restrictions/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.Method {
		case http.MethodGet, http.MethodPut:
			_ = json.NewEncoder(w).Encode(branchRestrictionResponse)
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})
	mux.HandleFunc("POST /repositories/{workspace}/{repo_slug}/branch-restrictions", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(branchRestrictionResponse)
	})

	// Catch-all for any other API calls during tests
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]any{"status": "ok"})
	})

	return httptest.NewServer(mux)
}

// setMockEnv configures environment variables to point at a mock server.
func setMockEnv(t *testing.T, serverURL string) {
	t.Helper()
	t.Setenv("BITBUCKET_USERNAME", "testuser")
	t.Setenv("BITBUCKET_TOKEN", "testtoken")
	t.Setenv("BITBUCKET_BASE_URL", serverURL)
}

// ─── Data source acceptance tests ─────────────────────────────────────────────

func TestAccDataSourceRepos_Read(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						base_url = %q
					}

					data "bitbucket_repos" "test" {
						workspace = "testworkspace"
						repo_slug = "test-repo"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_repos.test", "api_response"),
					resource.TestCheckResourceAttrSet("data.bitbucket_repos.test", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceWorkspaces_Read(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						base_url = %q
					}

					data "bitbucket_workspaces" "test" {
						workspace = "testworkspace"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_workspaces.test", "api_response"),
					resource.TestCheckResourceAttrSet("data.bitbucket_workspaces.test", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceUsers_Read(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						base_url = %q
					}

					data "bitbucket_users" "test" {
						selected_user = "testuser"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_users.test", "api_response"),
				),
			},
		},
	})
}

// ─── Resource acceptance tests ────────────────────────────────────────────────

func TestAccResourceRepos_CRUD(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// Create and Read
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						base_url = %q
					}

					resource "bitbucket_repos" "test" {
						workspace = "testworkspace"
						repo_slug = "test-repo"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("bitbucket_repos.test", "id"),
					resource.TestCheckResourceAttrSet("bitbucket_repos.test", "api_response"),
					resource.TestCheckResourceAttr("bitbucket_repos.test", "workspace", "testworkspace"),
					resource.TestCheckResourceAttr("bitbucket_repos.test", "repo_slug", "test-repo"),
				),
			},
		},
	})
}

func TestAccResourceProjects_CRUD(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						base_url = %q
					}

					resource "bitbucket_projects" "test" {
						workspace   = "testworkspace"
						project_key = "TEST"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("bitbucket_projects.test", "api_response"),
					resource.TestCheckResourceAttr("bitbucket_projects.test", "workspace", "testworkspace"),
					resource.TestCheckResourceAttr("bitbucket_projects.test", "project_key", "TEST"),
				),
			},
		},
	})
}

// ─── Sub-resource acceptance tests ────────────────────────────────────────────

func TestAccDataSourceWorkspaceHooks_Read(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						base_url = %q
					}

					data "bitbucket_workspace_hooks" "test" {
						workspace = "testworkspace"
						uid       = "hook-uuid"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_workspace_hooks.test", "api_response"),
				),
			},
		},
	})
}

func TestAccResourceWorkspaceHooks_CRUD(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						base_url = %q
					}

					resource "bitbucket_workspace_hooks" "test" {
						workspace = "testworkspace"
						uid       = "hook-uuid"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("bitbucket_workspace_hooks.test", "api_response"),
					resource.TestCheckResourceAttr("bitbucket_workspace_hooks.test", "workspace", "testworkspace"),
				),
			},
		},
	})
}

func TestAccDataSourceDefaultReviewers_Read(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						base_url = %q
					}

					data "bitbucket_default_reviewers" "test" {
						workspace       = "testworkspace"
						repo_slug       = "test-repo"
						target_username = "testuser"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_default_reviewers.test", "api_response"),
				),
			},
		},
	})
}

func TestAccResourceDefaultReviewers_CRUD(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						base_url = %q
					}

					resource "bitbucket_default_reviewers" "test" {
						workspace       = "testworkspace"
						repo_slug       = "test-repo"
						target_username = "testuser"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("bitbucket_default_reviewers.test", "api_response"),
					resource.TestCheckResourceAttr("bitbucket_default_reviewers.test", "workspace", "testworkspace"),
				),
			},
		},
	})
}

func TestAccDataSourcePipelineVariables_Read(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						base_url = %q
					}

					data "bitbucket_pipeline_variables" "test" {
						workspace     = "testworkspace"
						repo_slug     = "test-repo"
						variable_uuid = "{var-uuid}"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_pipeline_variables.test", "api_response"),
				),
			},
		},
	})
}

func TestAccResourcePipelineVariables_CRUD(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						base_url = %q
					}

					resource "bitbucket_pipeline_variables" "test" {
						workspace     = "testworkspace"
						repo_slug     = "test-repo"
						variable_uuid = "{var-uuid}"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("bitbucket_pipeline_variables.test", "api_response"),
					resource.TestCheckResourceAttr("bitbucket_pipeline_variables.test", "workspace", "testworkspace"),
				),
			},
		},
	})
}

func TestAccDataSourceRepoDeployKeys_Read(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						base_url = %q
					}

					data "bitbucket_repo_deploy_keys" "test" {
						workspace = "testworkspace"
						repo_slug = "test-repo"
						key_id    = "123"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_repo_deploy_keys.test", "api_response"),
				),
			},
		},
	})
}

func TestAccDataSourceRepoGroupPermissions_Read(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						base_url = %q
					}

					data "bitbucket_repo_group_permissions" "test" {
						workspace  = "testworkspace"
						repo_slug  = "test-repo"
						group_slug = "developers"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_repo_group_permissions.test", "api_response"),
				),
			},
		},
	})
}

func TestAccDataSourceRepoUserPermissions_Read(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						base_url = %q
					}

					data "bitbucket_repo_user_permissions" "test" {
						workspace        = "testworkspace"
						repo_slug        = "test-repo"
						selected_user_id = "{user-uuid}"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_repo_user_permissions.test", "api_response"),
				),
			},
		},
	})
}

// ─── Wave 2: additional sub-resource acceptance tests ─────────────────────────

func TestAccDataSourceTags_Read(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						base_url = %q
					}

					data "bitbucket_tags" "test" {
						workspace = "testworkspace"
						repo_slug = "test-repo"
						name      = "v1.0.0"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_tags.test", "api_response"),
				),
			},
		},
	})
}

func TestAccResourceTags_CRUD(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						base_url = %q
					}

					resource "bitbucket_tags" "test" {
						workspace = "testworkspace"
						repo_slug = "test-repo"
						name      = "v1.0.0"
						type      = "tag"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("bitbucket_tags.test", "api_response"),
				),
			},
		},
	})
}

func TestAccDataSourcePipelineSshKeys_Read(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						base_url = %q
					}

					data "bitbucket_pipeline_ssh_keys" "test" {
						workspace = "testworkspace"
						repo_slug = "test-repo"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_pipeline_ssh_keys.test", "api_response"),
				),
			},
		},
	})
}

func TestAccDataSourcePipelineSchedules_Read(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						base_url = %q
					}

					data "bitbucket_pipeline_schedules" "test" {
						workspace     = "testworkspace"
						repo_slug     = "test-repo"
						schedule_uuid = "{schedule-uuid}"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_pipeline_schedules.test", "api_response"),
				),
			},
		},
	})
}

func TestAccResourcePipelineSchedules_CRUD(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						base_url = %q
					}

					resource "bitbucket_pipeline_schedules" "test" {
						workspace     = "testworkspace"
						repo_slug     = "test-repo"
						schedule_uuid = "{schedule-uuid}"
						cron_pattern  = "0 0 12 * * ? *"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("bitbucket_pipeline_schedules.test", "api_response"),
				),
			},
		},
	})
}

func TestAccDataSourcePipelineKnownHosts_Read(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						base_url = %q
					}

					data "bitbucket_pipeline_known_hosts" "test" {
						workspace       = "testworkspace"
						repo_slug       = "test-repo"
						known_host_uuid = "{known-host-uuid}"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_pipeline_known_hosts.test", "api_response"),
				),
			},
		},
	})
}

func TestAccDataSourcePipelineConfig_Read(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						base_url = %q
					}

					data "bitbucket_pipeline_config" "test" {
						workspace = "testworkspace"
						repo_slug = "test-repo"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_pipeline_config.test", "api_response"),
				),
			},
		},
	})
}

func TestAccResourcePRComments_CRUD(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						base_url = %q
					}

					resource "bitbucket_pr_comments" "test" {
						workspace        = "testworkspace"
						repo_slug        = "test-repo"
						pull_request_id  = "1"
						comment_id       = "1"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("bitbucket_pr_comments.test", "api_response"),
				),
			},
		},
	})
}

func TestAccResourceIssueComments_CRUD(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						base_url = %q
					}

					resource "bitbucket_issue_comments" "test" {
						workspace  = "testworkspace"
						repo_slug  = "test-repo"
						issue_id   = "1"
						comment_id = "1"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("bitbucket_issue_comments.test", "api_response"),
				),
			},
		},
	})
}

// ─── Provider configuration tests ─────────────────────────────────────────────

func TestAccProvider_ConfigureWithUsername(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						username = "testuser"
						token    = "testtoken"
						base_url = %q
					}

					data "bitbucket_repos" "test" {
						workspace = "testworkspace"
						repo_slug = "test-repo"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_repos.test", "api_response"),
				),
			},
		},
	})
}

func TestAccProvider_ConfigureWithToken(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	// Only set token, not username
	t.Setenv("BITBUCKET_USERNAME", "")
	t.Setenv("BITBUCKET_TOKEN", "test-oauth-token")
	t.Setenv("BITBUCKET_BASE_URL", srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						token    = "test-oauth-token"
						base_url = %q
					}

					data "bitbucket_repos" "test" {
						workspace = "testworkspace"
						repo_slug = "test-repo"
					}
				`, srv.URL),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_repos.test", "api_response"),
				),
			},
		},
	})
}

// ─── Branch restrictions order-insensitivity acceptance test ─────────────────

// TestAccBitbucketBranchRestrictionsUsersOrderInsensitive exercises the
// fix for the bitbucket_branch_restrictions `users` ordering bug.
//
// The mock server returns `users` in lexicographic UUID order regardless
// of submission order — i.e. the exact upstream Bitbucket Cloud behaviour
// that triggers both the pre-0.15.6 "Provider produced inconsistent result
// after apply" and the v0.15.6 "Provider produced invalid plan" failures.
//
// The HCL config below deliberately writes the two users in REVERSE
// (bbbb before aaaa) so that:
//
//  1. plan + apply must succeed (no "invalid plan" diagnostic; the
//     setLikeListType.ListSemanticEquals lets the framework treat the
//     reordered API response as semantically equal to the config),
//  2. a follow-up plan must be empty (no perpetual diff — the operator's
//     order is preserved in state and the API order doesn't drift it),
//  3. the resource's persisted users length matches the API response,
//     proving state was actually saved (the pre-0.15.6 failure mode
//     silently dropped state and re-created the resource on every plan).
func TestAccBitbucketBranchRestrictionsUsersOrderInsensitive(t *testing.T) {
	srv := startMockServer(t)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	config := fmt.Sprintf(`
		provider "bitbucket" {
			base_url = %q
		}

		resource "bitbucket_branch_restrictions" "example" {
			workspace         = "testworkspace"
			repo_slug         = "test-repo"
			kind              = "push"
			branch_match_kind = "glob"
			pattern           = "develop"

			users = [
				{ uuid = "{bbbbbbbb-0000-0000-0000-000000000002}" },
				{ uuid = "{aaaaaaaa-0000-0000-0000-000000000001}" },
			]
		}
	`, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// Apply: must succeed even though the API will return users
			// in the opposite order. Without the order-insensitive custom
			// list type this step fails with either:
			//   - "Provider produced invalid plan" (v0.15.6), or
			//   - "Provider produced inconsistent result after apply" (pre-0.15.6).
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("bitbucket_branch_restrictions.example", "id"),
					resource.TestCheckResourceAttr("bitbucket_branch_restrictions.example", "users.#", "2"),
				),
			},
			// Re-plan with the same config: must be a no-op. This is the
			// "perpetual diff on add/reorder" guard — the framework only
			// reports the plan as empty when ListSemanticEquals returns
			// true for state vs. config.
			{
				Config:   config,
				PlanOnly: true,
			},
		},
	})
}

// ─── Real API acceptance tests (run when TF_ACC=1 and credentials are set) ──

// skipIfNoRealAPI skips the test if real API credentials are not configured.
// Returns the workspace name when credentials are available.
func skipIfNoRealAPI(t *testing.T) string {
	t.Helper()
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set, skipping real API acceptance test")
	}
	if os.Getenv("BITBUCKET_USERNAME") == "" && os.Getenv("BITBUCKET_TOKEN") == "" {
		t.Skip("No Bitbucket credentials set, skipping real API test")
	}
	workspace := os.Getenv("BITBUCKET_TEST_WORKSPACE")
	if workspace == "" {
		t.Skip("BITBUCKET_TEST_WORKSPACE not set, skipping real API test")
	}
	return workspace
}

// TestAccRealAPI_DataSourceWorkspaces reads a workspace and verifies the response.
func TestAccRealAPI_DataSourceWorkspaces(t *testing.T) {
	workspace := skipIfNoRealAPI(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {}

					data "bitbucket_workspaces" "test" {
						workspace = %q
					}
				`, workspace),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_workspaces.test", "api_response"),
					resource.TestCheckResourceAttrSet("data.bitbucket_workspaces.test", "id"),
				),
			},
		},
	})
}

// TestAccRealAPI_DataSourceCurrentUser reads the current authenticated user.
func TestAccRealAPI_DataSourceCurrentUser(t *testing.T) {
	skipIfNoRealAPI(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					provider "bitbucket" {}

					data "bitbucket_current_user" "me" {}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_current_user.me", "api_response"),
					resource.TestCheckResourceAttrSet("data.bitbucket_current_user.me", "id"),
				),
			},
		},
	})
}

// TestAccRealAPI_DataSourceUsers reads a user profile using the current user's UUID.
// The Bitbucket API (post-GDPR) requires a UUID in the {uuid} format for selected_user.
// We obtain it via the current-user data source and pass it through jsondecode.
func TestAccRealAPI_DataSourceUsers(t *testing.T) {
	skipIfNoRealAPI(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					provider "bitbucket" {}

					data "bitbucket_current_user" "me" {}

					data "bitbucket_users" "test" {
						selected_user = jsondecode(data.bitbucket_current_user.me.api_response).uuid
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_users.test", "api_response"),
					resource.TestCheckResourceAttrSet("data.bitbucket_users.test", "id"),
				),
			},
		},
	})
}

// TestAccRealAPI_ResourceProjects_CRUD creates, reads, updates, and deletes a project.
// Uses a random project key to ensure idempotency across test runs.
func TestAccRealAPI_ResourceProjects_CRUD(t *testing.T) {
	workspace := skipIfNoRealAPI(t)

	// Generate a unique project key (uppercase, max 10 chars) so tests are idempotent.
	suffix := strings.ToUpper(acctest.RandStringFromCharSet(5, acctest.CharSetAlpha))
	projectKey := "TF" + suffix
	projectName := "Terraform Test " + suffix

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckProjectDestroy(workspace, projectKey),
		Steps: []resource.TestStep{
			// Create
			{
				Config: testAccProjectConfig(workspace, projectKey, projectName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("bitbucket_projects.test", "id"),
					resource.TestCheckResourceAttrSet("bitbucket_projects.test", "api_response"),
					resource.TestCheckResourceAttr("bitbucket_projects.test", "workspace", workspace),
				),
			},
			// Update
			{
				Config: testAccProjectConfig(workspace, projectKey, projectName+" Updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("bitbucket_projects.test", "id"),
				),
			},
			// Destroy is handled automatically by the test framework
		},
	})
}

// testAccProjectConfig returns a Terraform config for a bitbucket_projects resource.
func testAccProjectConfig(workspace, key, name string) string {
	return fmt.Sprintf(`
		provider "bitbucket" {}

		resource "bitbucket_projects" "test" {
			workspace    = %q
			project_key  = %q
			request_body = jsonencode({
				name = %q
				key  = %q
			})
		}
	`, workspace, key, name, key)
}

// testAccCheckProjectDestroy verifies the project was deleted from the Bitbucket API.
func testAccCheckProjectDestroy(workspace, projectKey string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %v", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_, err = handlers.DispatchRaw(ctx, c, handlers.Request{
			Method:      "GET",
			URLTemplate: "/workspaces/{workspace}/projects/{project_key}",
			PathParams:  map[string]string{"workspace": workspace, "project_key": projectKey},
			All:         false,
		})
		if err == nil {
			return fmt.Errorf("project %s still exists in workspace %s after destroy", projectKey, workspace)
		}
		// Verify the error is a Bitbucket API 404 (not a network/auth error).
		if !strings.Contains(err.Error(), "bitbucket API error 404") {
			return fmt.Errorf("unexpected error checking project %s destroy: %v", projectKey, err)
		}
		return nil
	}
}

// TestAccRealAPI_DataSourceRepos reads a specific repository from the test workspace.
// Requires BITBUCKET_TEST_REPO to be set, otherwise lists the workspace.
func TestAccRealAPI_DataSourceRepos(t *testing.T) {
	workspace := skipIfNoRealAPI(t)
	repoSlug := os.Getenv("BITBUCKET_TEST_REPO")
	if repoSlug == "" {
		t.Skip("BITBUCKET_TEST_REPO not set, skipping repos read test")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {}

					data "bitbucket_repos" "test" {
						workspace = %q
						repo_slug = %q
					}
				`, workspace, repoSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_repos.test", "api_response"),
					resource.TestCheckResourceAttrSet("data.bitbucket_repos.test", "id"),
				),
			},
		},
	})
}

// TestAccRealAPI_ProviderAuth verifies the provider works with explicit auth config.
func TestAccRealAPI_ProviderAuth(t *testing.T) {
	workspace := skipIfNoRealAPI(t)
	username := os.Getenv("BITBUCKET_USERNAME")
	token := os.Getenv("BITBUCKET_TOKEN")
	if username == "" || token == "" {
		t.Skip("BITBUCKET_USERNAME or BITBUCKET_TOKEN not set")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {
						username = %q
						token    = %q
					}

					data "bitbucket_workspaces" "test" {
						workspace = %q
					}
				`, username, token, workspace),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_workspaces.test", "api_response"),
					resource.TestCheckResourceAttrSet("data.bitbucket_workspaces.test", "id"),
				),
			},
		},
	})
}

// TestAccRealAPI_DataSource_Commits reads a specific commit via its SHA.
// Chains through refs to discover the HEAD commit on "main".
// Requires BITBUCKET_TEST_REPO to be set.
func TestAccRealAPI_DataSource_Commits(t *testing.T) {
	workspace := skipIfNoRealAPI(t)
	repoSlug := os.Getenv("BITBUCKET_TEST_REPO")
	if repoSlug == "" {
		t.Skip("BITBUCKET_TEST_REPO not set, skipping commits test")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {}

					data "bitbucket_refs" "main" {
						workspace = %q
						repo_slug = %q
						name      = "main"
					}

					data "bitbucket_commits" "test" {
						workspace = %q
						repo_slug = %q
						commit    = jsondecode(data.bitbucket_refs.main.api_response).target.hash
					}
				`, workspace, repoSlug, workspace, repoSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_commits.test", "api_response"),
					resource.TestCheckResourceAttrSet("data.bitbucket_commits.test", "id"),
				),
			},
		},
	})
}

// TestAccRealAPI_DataSource_Refs reads the "main" branch for a repository.
// Requires BITBUCKET_TEST_REPO to be set.
func TestAccRealAPI_DataSource_Refs(t *testing.T) {
	workspace := skipIfNoRealAPI(t)
	repoSlug := os.Getenv("BITBUCKET_TEST_REPO")
	if repoSlug == "" {
		t.Skip("BITBUCKET_TEST_REPO not set, skipping refs test")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {}

					data "bitbucket_refs" "test" {
						workspace = %q
						repo_slug = %q
						name      = "main"
					}
				`, workspace, repoSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_refs.test", "api_response"),
					resource.TestCheckResourceAttrSet("data.bitbucket_refs.test", "id"),
				),
			},
		},
	})
}

// TestAccRealAPI_DataSource_BranchingModel reads the branching model for a repository.
// Requires BITBUCKET_TEST_REPO to be set.
func TestAccRealAPI_DataSource_BranchingModel(t *testing.T) {
	workspace := skipIfNoRealAPI(t)
	repoSlug := os.Getenv("BITBUCKET_TEST_REPO")
	if repoSlug == "" {
		t.Skip("BITBUCKET_TEST_REPO not set, skipping branching model test")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {}

					data "bitbucket_branching_model" "test" {
						workspace = %q
						repo_slug = %q
					}
				`, workspace, repoSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_branching_model.test", "api_response"),
					resource.TestCheckResourceAttrSet("data.bitbucket_branching_model.test", "id"),
				),
			},
		},
	})
}

// TestAccRealAPI_DataSource_HookTypes reads available webhook event types.
// No additional parameters required — GET /hook_events returns event categories.
func TestAccRealAPI_DataSource_HookTypes(t *testing.T) {
	skipIfNoRealAPI(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
					provider "bitbucket" {}

					data "bitbucket_hook_types" "test" {}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_hook_types.test", "api_response"),
					resource.TestCheckResourceAttrSet("data.bitbucket_hook_types.test", "id"),
				),
			},
		},
	})
}

// TestAccRealAPI_DataSource_WorkspacePermissions reads the current user's permission
// on the test workspace. Only requires workspace — GET /user/workspaces/{workspace}/permission.
func TestAccRealAPI_DataSource_WorkspacePermissions(t *testing.T) {
	workspace := skipIfNoRealAPI(t)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {}

					data "bitbucket_workspace_permissions" "test" {
						workspace = %q
					}
				`, workspace),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_workspace_permissions.test", "api_response"),
					resource.TestCheckResourceAttrSet("data.bitbucket_workspace_permissions.test", "id"),
				),
			},
		},
	})
}

// TestAccRealAPI_ResourceBranchRestrictions_OrderInsensitiveUsers exercises the
// real Bitbucket API end-to-end against the regression class fixed by the
// custom setLikeListType (multiset semantic equality + planned-order alignment
// of the response array).
//
// Three concrete bugs this test guards against on the real API:
//
//  1. "Provider produced inconsistent result after apply" (pre-0.15.6) — the API
//     echoes a fully-populated nested user object (display_name, created_on)
//     whose computed inner fields were Unknown in the plan; without the
//     planned-order alignment the post-apply state diverges from the plan.
//  2. "Provider produced invalid plan" (v0.15.6) — the lexicographic-sort plan
//     modifier reordered users at plan time, then refresh undid it.
//  3. Perpetual diff on add/reorder — config order ≠ API response order makes
//     every subsequent plan show a positional reshuffle.
//
// The mock-based TestAccBitbucketBranchRestrictionsUsersOrderInsensitive
// covers the same logic against a controlled response, but only the real API
// returns nested objects with all the computed inner fields present (and in
// whatever order Bitbucket chooses), which is exactly the surface where the
// regressions lived.
//
// Step 1 catches (1) on the single-user path (the user object's display_name /
// created_on are computed). Step 2 catches (3) for the same case. When
// BITBUCKET_TEST_USER_2 is also configured, steps 3-5 add full multi-element
// coverage (reorder + remove) that catches (1)+(2)+(3) on the multiset path.
//
// Required env: BITBUCKET_TEST_WORKSPACE, BITBUCKET_TEST_REPO, BITBUCKET_TEST_USER
// (UUID with surrounding braces, e.g. "{abcdef01-…}"). Optional:
// BITBUCKET_TEST_USER_2 for full multi-element coverage.
func TestAccRealAPI_ResourceBranchRestrictions_OrderInsensitiveUsers(t *testing.T) {
	workspace := skipIfNoRealAPI(t)
	repoSlug := os.Getenv("BITBUCKET_TEST_REPO")
	if repoSlug == "" {
		t.Skip("BITBUCKET_TEST_REPO not set, skipping branch restrictions real API test")
	}
	user1 := os.Getenv("BITBUCKET_TEST_USER")
	if user1 == "" {
		t.Skip("BITBUCKET_TEST_USER not set, skipping branch restrictions real API test")
	}
	user2 := os.Getenv("BITBUCKET_TEST_USER_2") // optional second user UUID

	pattern := "tf-acc-test-" + strings.ToLower(acctest.RandStringFromCharSet(8, acctest.CharSetAlphaNum))

	cfg := func(uuids ...string) string {
		var users strings.Builder
		for _, u := range uuids {
			fmt.Fprintf(&users, "    { uuid = %q },\n", u)
		}
		// `kind = "delete"` is used (not "push") because Bitbucket Cloud's
		// /branch-restrictions POST validates that every user listed in a
		// `kind=push` restriction has repository write access, and returns
		// `500 Something went wrong` (not 403) if any does not — a long-
		// standing Bitbucket quirk that's orthogonal to what this test
		// guards (multiset semantic equality + planned-order alignment of
		// the response array). `kind=delete` accepts the same `users` array
		// shape and exercises the exact same provider code paths, but does
		// not perform the per-user permission validation, so it succeeds
		// against any UUID the caller can read.
		//
		// `groups` is intentionally omitted: Bitbucket's branch-restrictions
		// POST returns 500 when an empty `groups` array is sent alongside a
		// non-empty `users` (the matching mock-server test in
		// TestAccBitbucketBranchRestrictionsUsersOrderInsensitive also omits
		// it).
		return fmt.Sprintf(`
			provider "bitbucket" {}

			resource "bitbucket_branch_restrictions" "test" {
				workspace         = %q
				repo_slug         = %q
				kind              = "delete"
				branch_match_kind = "glob"
				pattern           = %q

				users = [
%s				]
			}
		`, workspace, repoSlug, pattern, users.String())
	}

	steps := []resource.TestStep{
		// (1) Create with one user. The API echoes display_name / created_on
		// for that user; without setLikeListValue's planned-order alignment
		// this step previously failed with "Provider produced inconsistent
		// result after apply" because the computed inner fields are Unknown
		// in the plan.
		{
			Config: cfg(user1),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("bitbucket_branch_restrictions.test", "id"),
				resource.TestCheckResourceAttr("bitbucket_branch_restrictions.test", "users.#", "1"),
				resource.TestCheckResourceAttr("bitbucket_branch_restrictions.test", "users.0.uuid", user1),
				resource.TestCheckResourceAttrSet("bitbucket_branch_restrictions.test", "users.0.display_name"),
			),
		},
		// (2) Re-plan with the same config — must be empty. Catches the
		// perpetual-diff class.
		{
			Config:   cfg(user1),
			PlanOnly: true,
		},
	}

	if user2 != "" {
		steps = append(steps,
			// (3) Update to two users in {a, b} order. The API may echo
			// {b, a}; this exercises the multi-element multiset apply path
			// against the real response shape.
			resource.TestStep{
				Config: cfg(user1, user2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_branch_restrictions.test", "users.#", "2"),
				),
			},
			// (4) Reorder to {b, a} — must be a no-op plan. Catches the
			// v0.15.6 "Provider produced invalid plan" regression and the
			// silent-reorder perpetual-diff bug.
			resource.TestStep{
				Config:   cfg(user2, user1),
				PlanOnly: true,
			},
			// (5) Drop one user — must succeed and result in a clean
			// one-element list. Catches the "remove an element" path
			// through the multiset comparison.
			resource.TestStep{
				Config: cfg(user2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_branch_restrictions.test", "users.#", "1"),
					resource.TestCheckResourceAttr("bitbucket_branch_restrictions.test", "users.0.uuid", user2),
				),
			},
		)
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckBranchRestrictionDestroy(workspace, repoSlug, pattern),
		Steps:                    steps,
	})
}

// testAccCheckBranchRestrictionDestroy verifies no branch restriction matching
// the test pattern (kind=delete) remains in the workspace/repo after destroy.
// The resource ID is generated by Bitbucket and not stable across runs, so we
// query by the (kind, pattern) tuple that the test owns.
func testAccCheckBranchRestrictionDestroy(workspace, repoSlug, pattern string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		c, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %v", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		result, err := handlers.DispatchRaw(ctx, c, handlers.Request{
			Method:      "GET",
			URLTemplate: "/repositories/{workspace}/{repo_slug}/branch-restrictions",
			PathParams:  map[string]string{"workspace": workspace, "repo_slug": repoSlug},
			QueryParams: map[string]string{"kind": "delete", "pattern": pattern},
			All:         true,
		})
		if err != nil {
			return fmt.Errorf("failed to list branch restrictions for destroy check: %v", err)
		}
		if items, ok := result.([]any); ok && len(items) > 0 {
			return fmt.Errorf("branch restriction kind=delete, pattern=%q still exists in %s/%s after destroy (%d found)",
				pattern, workspace, repoSlug, len(items))
		}
		return nil
	}
}

// TestAccRealAPI_DataSource_UserEmails reads a specific email address for the current user.
// Uses BITBUCKET_USERNAME (the Atlassian account email) as the email parameter.
func TestAccRealAPI_DataSource_UserEmails(t *testing.T) {
	skipIfNoRealAPI(t)
	email := os.Getenv("BITBUCKET_USERNAME")
	if email == "" {
		t.Skip("BITBUCKET_USERNAME not set, skipping user emails test")
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					provider "bitbucket" {}

					data "bitbucket_user_emails" "test" {
						email = %q
					}
				`, email),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_user_emails.test", "api_response"),
					resource.TestCheckResourceAttrSet("data.bitbucket_user_emails.test", "id"),
				),
			},
		},
	})
}
