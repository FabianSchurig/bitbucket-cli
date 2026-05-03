package tfprovider_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sort"
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

const testAccRealAPITimeout = 30 * time.Second

var testAccRepoUserPermissionOps = tfprovider.MapCRUDOps("repo-user-permissions", tfprovider.ReposResourceGroup.AllOps)

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
		ctx, cancel := context.WithTimeout(context.Background(), testAccRealAPITimeout)
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
// created_on are computed). Step 2 catches (3) for the same case. Steps 3-4
// model the real-world Terraform graph: a bitbucket_repo_user_permissions
// resource grants the second workspace member repository write access first,
// the branch restriction depends on it, and Terraform destroy tears the branch
// restriction down before removing the temporary permission.
//
// Required env: BITBUCKET_TEST_WORKSPACE, BITBUCKET_TEST_REPO.
func TestAccRealAPI_ResourceBranchRestrictions_OrderInsensitiveUsers(t *testing.T) {
	workspace := skipIfNoRealAPI(t)
	repoSlug := os.Getenv("BITBUCKET_TEST_REPO")
	if repoSlug == "" {
		t.Skip("BITBUCKET_TEST_REPO not set, skipping branch restrictions real API test")
	}

	c, err := client.NewClient()
	if err != nil {
		t.Fatalf("failed to create Bitbucket client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), testAccRealAPITimeout)
	defer cancel()

	user1, err := testAccCurrentUserUUID(ctx, c)
	if err != nil {
		t.Fatalf("failed to read current user UUID: %v", err)
	}

	user2, restorePermission, err := testAccPrepareSecondBranchRestrictionUser(ctx, c, workspace, repoSlug, user1)
	if err != nil {
		t.Fatalf("failed to prepare a second repository user for branch restriction testing: %v", err)
	}
	defer func() {
		restoreCtx, restoreCancel := context.WithTimeout(context.Background(), testAccRealAPITimeout)
		defer restoreCancel()
		if err := restorePermission(restoreCtx); err != nil {
			t.Logf("failed to restore repository permission for %s: %v", user2, err)
		}
	}()

	const restrictionKind = "push"
	patternUser := strings.Trim(user1, "{}")
	if len(patternUser) > 8 {
		patternUser = patternUser[:8]
	}
	patternBase := "tf-acc-order-insensitive-users-" + strings.ToLower(patternUser)
	patternSingle := patternBase + "-single"
	patternMulti := patternBase + "-multi"
	if err := testAccDeleteBranchRestrictionsByPattern(ctx, c, workspace, repoSlug, restrictionKind, patternSingle, patternMulti); err != nil {
		t.Fatalf("failed to clean up existing branch restrictions before test: %v", err)
	}
	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), testAccRealAPITimeout)
		defer cleanupCancel()
		if err := testAccDeleteBranchRestrictionsByPattern(cleanupCtx, c, workspace, repoSlug, restrictionKind, patternSingle, patternMulti); err != nil {
			t.Logf("failed to clean up branch restrictions after test: %v", err)
		}
	}()

	cfg := func(resourceName, pattern string, createUser2PermissionResource bool, uuids ...string) string {
		var users strings.Builder
		for _, u := range uuids {
			fmt.Fprintf(&users, "    { uuid = %q },\n", u)
		}
		permissionResource := ""
		dependsOn := ""
		if createUser2PermissionResource {
			permissionResource = fmt.Sprintf(`
			resource "bitbucket_repo_user_permissions" "branch_restriction_user2" {
				workspace        = %[1]q
				repo_slug        = %[2]q
				selected_user_id = %[3]q
				request_body     = jsonencode({ permission = "write" })
			}
`, workspace, repoSlug, user2)
			dependsOn = `
				depends_on = [bitbucket_repo_user_permissions.branch_restriction_user2]
`
		}
		// `kind = "push"` is the branch-restriction kind that supports a
		// user allow-list. Bitbucket validates that every listed user has
		// repository write access before accepting the restriction. The
		// multi-user steps grant the second user's repository permission with
		// the Terraform resource exposed by this provider, matching the
		// real-world dependency graph users should write.
		//
		// `groups` is intentionally omitted: Bitbucket's branch-restrictions
		// POST returns 500 when an empty `groups` array is sent alongside a
		// non-empty `users` (the matching mock-server test in
		// TestAccBitbucketBranchRestrictionsUsersOrderInsensitive also omits
		// it).
		return fmt.Sprintf(`
			provider "bitbucket" {}
%s

			resource "bitbucket_branch_restrictions" %q {
				workspace         = %q
				repo_slug         = %q
				kind              = %q
				branch_match_kind = "glob"
				pattern           = %q

				users = [
%s				]
%s
			}
		`, permissionResource, resourceName, workspace, repoSlug, restrictionKind, pattern, users.String(), dependsOn)
	}

	steps := []resource.TestStep{
		// (1) Create with one user. The API echoes display_name / created_on
		// for that user; without setLikeListValue's planned-order alignment
		// this step previously failed with "Provider produced inconsistent
		// result after apply" because the computed inner fields are Unknown
		// in the plan.
		{
			Config: cfg("test_single", patternSingle, false, user1),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttrSet("bitbucket_branch_restrictions.test_single", "id"),
				resource.TestCheckResourceAttr("bitbucket_branch_restrictions.test_single", "users.#", "1"),
				resource.TestCheckResourceAttr("bitbucket_branch_restrictions.test_single", "users.0.uuid", user1),
				resource.TestCheckResourceAttrSet("bitbucket_branch_restrictions.test_single", "users.0.display_name"),
			),
		},
		// (2) Re-plan with the same config — must be empty. Catches the
		// perpetual-diff class.
		{
			Config:   cfg("test_single", patternSingle, false, user1),
			PlanOnly: true,
		},
		// (3) Create a fresh two-user restriction in {a, b} order. The API
		// may echo {b, a}; the provider must still preserve the configured
		// order in state while keeping computed fields populated.
		{
			Config: cfg("test_multi", patternMulti, true, user1, user2),
			Check: resource.ComposeTestCheckFunc(
				resource.TestCheckResourceAttr("bitbucket_branch_restrictions.test_multi", "users.#", "2"),
				resource.TestCheckResourceAttr("bitbucket_branch_restrictions.test_multi", "users.0.uuid", user1),
				resource.TestCheckResourceAttr("bitbucket_branch_restrictions.test_multi", "users.1.uuid", user2),
				resource.TestCheckResourceAttrSet("bitbucket_branch_restrictions.test_multi", "users.0.display_name"),
				resource.TestCheckResourceAttrSet("bitbucket_branch_restrictions.test_multi", "users.1.display_name"),
			),
		},
		// (4) Reorder to {b, a} — must be a no-op plan. Catches the
		// v0.15.6 "Provider produced invalid plan" regression and the
		// silent-reorder perpetual-diff bug on the real multi-user path.
		{
			Config:   cfg("test_multi", patternMulti, true, user2, user1),
			PlanOnly: true,
		},
	}

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckBranchRestrictionDestroy(workspace, repoSlug, restrictionKind, patternSingle, patternMulti),
		Steps:                    steps,
	})
}

// testAccCheckBranchRestrictionDestroy verifies no branch restriction matching
// the test pattern remains in the workspace/repo after destroy.
// The resource ID is generated by Bitbucket and not stable across runs, so we
// query by the (kind, pattern) tuple that the test owns.
func testAccCheckBranchRestrictionDestroy(workspace, repoSlug, kind string, patterns ...string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		c, err := client.NewClient()
		if err != nil {
			return fmt.Errorf("failed to create client: %v", err)
		}
		ctx, cancel := context.WithTimeout(context.Background(), testAccRealAPITimeout)
		defer cancel()
		for _, pattern := range patterns {
			if pattern == "" {
				continue
			}
			result, err := handlers.DispatchRaw(ctx, c, handlers.Request{
				Method:      "GET",
				URLTemplate: "/repositories/{workspace}/{repo_slug}/branch-restrictions",
				PathParams:  map[string]string{"workspace": workspace, "repo_slug": repoSlug},
				QueryParams: map[string]string{"kind": kind, "pattern": pattern},
				All:         true,
			})
			if err != nil {
				return fmt.Errorf("failed to list branch restrictions for destroy check: %v", err)
			}
			if items, ok := result.([]any); ok && len(items) > 0 {
				return fmt.Errorf("branch restriction kind=%s, pattern=%q still exists in %s/%s after destroy (%d found)",
					kind, pattern, workspace, repoSlug, len(items))
			}
		}
		return nil
	}
}

// testAccCurrentUserUUID returns the UUID of the currently authenticated user
// from the Bitbucket `/user` endpoint.
func testAccCurrentUserUUID(ctx context.Context, c *client.BBClient) (string, error) {
	result, err := handlers.DispatchRaw(ctx, c, handlers.Request{
		Method:      "GET",
		URLTemplate: "/user",
	})
	if err != nil {
		return "", err
	}
	user, ok := result.(map[string]any)
	if !ok {
		return "", fmt.Errorf("unexpected /user response type %T", result)
	}
	uuid, ok := user["uuid"].(string)
	if !ok || uuid == "" {
		return "", fmt.Errorf("current user response does not contain uuid")
	}
	return uuid, nil
}

// testAccPrepareSecondBranchRestrictionUser returns a second non-owner
// workspace member UUID for the real branch-restrictions acceptance test and a
// restore callback for the repository permission resource managed by Terraform
// during the test. The test itself grants write access via
// bitbucket_repo_user_permissions so Terraform exercises the real dependency
// order: permission first, branch restriction second, reverse on destroy.
func testAccPrepareSecondBranchRestrictionUser(ctx context.Context, c *client.BBClient, workspace, repoSlug, excludeUUID string) (string, func(context.Context) error, error) {
	userUUID, err := testAccFindAnotherWorkspaceMemberUUID(ctx, c, workspace, excludeUUID)
	if err != nil {
		return "", nil, err
	}
	if userUUID == "" {
		return "", nil, fmt.Errorf("no second non-owner workspace member found for workspace %q", workspace)
	}

	restore, err := testAccRepoUserPermissionRestoreFunc(ctx, c, workspace, repoSlug, userUUID)
	if err != nil {
		return "", nil, err
	}
	return userUUID, restore, nil
}

// testAccFindAnotherRepoWriterUUID lists explicit repository user permissions
// and returns a different user UUID with write/admin access. The UUID in
// excludeUUID is skipped. When multiple candidates exist, the UUIDs are sorted
// so the test consistently selects the same secondary user across runs.
func testAccFindAnotherRepoWriterUUID(ctx context.Context, c *client.BBClient, workspace, repoSlug, excludeUUID string) (string, error) {
	listOp, err := testAccRequireRepoUserPermissionOp(testAccRepoUserPermissionOps.List, "list")
	if err != nil {
		return "", err
	}
	result, err := handlers.DispatchRaw(ctx, c, handlers.Request{
		Method:      listOp.Method,
		URLTemplate: listOp.Path,
		PathParams:  map[string]string{"workspace": workspace, "repo_slug": repoSlug},
		All:         true,
	})
	if err != nil {
		return "", err
	}
	items, ok := result.([]any)
	if !ok {
		return "", fmt.Errorf("unexpected repo user permissions response type %T", result)
	}

	var candidates []string
	for _, raw := range items {
		item, ok := raw.(map[string]any)
		if !ok {
			return "", fmt.Errorf("unexpected repo user permission item response type %T", raw)
		}
		permission, _ := item["permission"].(string)
		if permission != "write" && permission != "admin" {
			continue
		}
		user, ok := item["user"].(map[string]any)
		if !ok {
			continue
		}
		uuid, _ := user["uuid"].(string)
		if uuid == "" || uuid == excludeUUID {
			continue
		}
		candidates = append(candidates, uuid)
	}
	sort.Strings(candidates)
	if len(candidates) == 0 {
		return "", nil
	}
	return candidates[0], nil
}

// testAccFindAnotherWorkspaceMemberUUID returns a deterministic non-owner
// workspace member UUID that differs from excludeUUID. This gives the real-API
// test a stable candidate whose repository permission can be managed inside the
// test itself.
func testAccFindAnotherWorkspaceMemberUUID(ctx context.Context, c *client.BBClient, workspace, excludeUUID string) (string, error) {
	result, err := handlers.DispatchRaw(ctx, c, handlers.Request{
		Method:      "GET",
		URLTemplate: "/workspaces/{workspace}/permissions",
		PathParams:  map[string]string{"workspace": workspace},
		All:         true,
	})
	if err != nil {
		return "", err
	}
	items, ok := result.([]any)
	if !ok {
		return "", fmt.Errorf("unexpected workspace permissions response type %T", result)
	}

	var candidates []string
	for _, raw := range items {
		item, ok := raw.(map[string]any)
		if !ok {
			return "", fmt.Errorf("unexpected workspace permission item response type %T", raw)
		}
		permission, _ := item["permission"].(string)
		if permission == "owner" {
			continue
		}
		user, ok := item["user"].(map[string]any)
		if !ok {
			continue
		}
		uuid, _ := user["uuid"].(string)
		if uuid == "" || uuid == excludeUUID {
			continue
		}
		candidates = append(candidates, uuid)
	}
	sort.Strings(candidates)
	if len(candidates) == 0 {
		return "", nil
	}
	return candidates[0], nil
}

// testAccEnsureRepoUserWritePermission makes sure the selected workspace user
// has repository write/admin access for the duration of the test and returns a
// callback that restores the previous explicit permission state afterwards.
func testAccEnsureRepoUserWritePermission(ctx context.Context, c *client.BBClient, workspace, repoSlug, selectedUserID string) (func(context.Context) error, error) {
	restore, err := testAccRepoUserPermissionRestoreFunc(ctx, c, workspace, repoSlug, selectedUserID)
	if err != nil {
		return nil, err
	}
	oldPermission, _, err := testAccRepoUserPermission(ctx, c, workspace, repoSlug, selectedUserID)
	if err != nil {
		return nil, err
	}
	if oldPermission == "write" || oldPermission == "admin" {
		return func(context.Context) error { return nil }, nil
	}
	if err := testAccSetRepoUserPermission(ctx, c, workspace, repoSlug, selectedUserID, "write"); err != nil {
		return nil, err
	}
	return restore, nil
}

func testAccRepoUserPermissionRestoreFunc(ctx context.Context, c *client.BBClient, workspace, repoSlug, selectedUserID string) (func(context.Context) error, error) {
	oldPermission, hadExplicitPermission, err := testAccRepoUserPermission(ctx, c, workspace, repoSlug, selectedUserID)
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context) error {
		// "none" is equivalent to no explicit repository permission in this
		// restore path, so it must use DELETE just like an absent original
		// permission.
		if hadExplicitPermission && oldPermission != "none" {
			return testAccSetRepoUserPermission(ctx, c, workspace, repoSlug, selectedUserID, oldPermission)
		}
		// Bitbucket may report "none" for a user without an effective
		// explicit repository grant, but the PUT endpoint rejects
		// permission=none. Deleting the permission is the API's restore path
		// for both "none" and absent original state.
		return testAccDeleteRepoUserPermission(ctx, c, workspace, repoSlug, selectedUserID)
	}, nil
}

// testAccRepoUserPermission returns the explicit repository permission, if any,
// for the selected user.
func testAccRepoUserPermission(ctx context.Context, c *client.BBClient, workspace, repoSlug, selectedUserID string) (permission string, explicit bool, err error) {
	readOp, err := testAccRequireRepoUserPermissionOp(testAccRepoUserPermissionOps.Read, "read")
	if err != nil {
		return "", false, err
	}
	result, err := handlers.DispatchRaw(ctx, c, handlers.Request{
		Method:      readOp.Method,
		URLTemplate: readOp.Path,
		PathParams: map[string]string{
			"workspace":        workspace,
			"repo_slug":        repoSlug,
			"selected_user_id": selectedUserID,
		},
	})
	if err != nil {
		if testAccBitbucketAPIStatus(err, http.StatusNotFound) {
			return "", false, nil
		}
		return "", false, err
	}
	permissionResponse, ok := result.(map[string]any)
	if !ok {
		return "", false, fmt.Errorf("unexpected repo user permission response type %T", result)
	}
	permission, ok = permissionResponse["permission"].(string)
	if !ok || permission == "" {
		return "", true, fmt.Errorf("repo user permission response does not contain permission")
	}
	return permission, true, nil
}

func testAccDeleteRepoUserPermission(ctx context.Context, c *client.BBClient, workspace, repoSlug, selectedUserID string) error {
	deleteOp, err := testAccRequireRepoUserPermissionOp(testAccRepoUserPermissionOps.Delete, "delete")
	if err != nil {
		return err
	}
	_, err = handlers.DispatchRaw(ctx, c, handlers.Request{
		Method:      deleteOp.Method,
		URLTemplate: deleteOp.Path,
		PathParams: map[string]string{
			"workspace":        workspace,
			"repo_slug":        repoSlug,
			"selected_user_id": selectedUserID,
		},
	})
	if testAccBitbucketAPIStatus(err, http.StatusNotFound) {
		return nil
	}
	return err
}

// testAccSetRepoUserPermission sets an explicit repository permission for the
// selected user.
func testAccSetRepoUserPermission(ctx context.Context, c *client.BBClient, workspace, repoSlug, selectedUserID, permission string) error {
	updateOp, err := testAccRequireRepoUserPermissionOp(testAccRepoUserPermissionOps.Update, "update")
	if err != nil {
		return err
	}
	body, err := json.Marshal(map[string]string{"permission": permission})
	if err != nil {
		return err
	}
	_, err = handlers.DispatchRaw(ctx, c, handlers.Request{
		Method:      updateOp.Method,
		URLTemplate: updateOp.Path,
		PathParams: map[string]string{
			"workspace":        workspace,
			"repo_slug":        repoSlug,
			"selected_user_id": selectedUserID,
		},
		Body: string(body),
	})
	return err
}

func testAccRequireRepoUserPermissionOp(op *tfprovider.OperationDef, action string) (*tfprovider.OperationDef, error) {
	if op == nil {
		return nil, fmt.Errorf("repo-user-permissions %s operation is not configured", action)
	}
	return op, nil
}

// testAccDeleteBranchRestrictionsByPattern deletes every repository branch
// restriction matching the provided kind/pattern pairs. It is used to make the
// real-API acceptance test idempotent before and after execution.
func testAccDeleteBranchRestrictionsByPattern(ctx context.Context, c *client.BBClient, workspace, repoSlug, kind string, patterns ...string) error {
	for _, pattern := range patterns {
		if pattern == "" {
			continue
		}
		result, err := handlers.DispatchRaw(ctx, c, handlers.Request{
			Method:      "GET",
			URLTemplate: "/repositories/{workspace}/{repo_slug}/branch-restrictions",
			PathParams:  map[string]string{"workspace": workspace, "repo_slug": repoSlug},
			QueryParams: map[string]string{"kind": kind, "pattern": pattern},
			All:         true,
		})
		if err != nil {
			return err
		}
		items, ok := result.([]any)
		if !ok {
			return fmt.Errorf("unexpected branch restriction list response type %T", result)
		}
		for _, raw := range items {
			item, ok := raw.(map[string]any)
			if !ok {
				return fmt.Errorf("unexpected branch restriction item response type %T", raw)
			}
			if gotPattern, _ := item["pattern"].(string); gotPattern != pattern {
				continue
			}
			id, ok := item["id"]
			if !ok {
				return fmt.Errorf("branch restriction response does not contain id")
			}
			_, err := handlers.DispatchRaw(ctx, c, handlers.Request{
				Method:      "DELETE",
				URLTemplate: "/repositories/{workspace}/{repo_slug}/branch-restrictions/{id}",
				PathParams: map[string]string{
					"workspace": workspace,
					"repo_slug": repoSlug,
					"id":        fmt.Sprint(id),
				},
			})
			if err != nil && !testAccBitbucketAPIStatus(err, http.StatusNotFound) {
				return err
			}
		}
	}
	return nil
}

func testAccBitbucketAPIStatus(err error, status int) bool {
	return err != nil && strings.Contains(err.Error(), fmt.Sprintf("bitbucket API error %d", status))
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
