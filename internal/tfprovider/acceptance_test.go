package tfprovider_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

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

// ─── Real API acceptance tests (run when TF_ACC=1 and credentials are set) ──

func TestAccRealAPI_DataSourceRepos(t *testing.T) {
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
