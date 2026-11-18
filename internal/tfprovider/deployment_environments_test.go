package tfprovider_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// mockEnvironments simulates Bitbucket's deployment environment endpoints
// closely enough to exercise the behaviours the generic runtime depends on:
// create requires environment_type.name, updates are POSTed to /changes as a
// {"change": {...}} envelope, the environment type cannot be changed, and an
// accepted change (202) only becomes visible after a delayed read.
type mockEnvironments struct {
	mu          sync.Mutex
	nextID      int
	envs        map[string]map[string]any
	pending     map[string]map[string]any
	staleReads  map[string]int
	changeCalls []string
}

func newMockEnvironments() *mockEnvironments {
	return &mockEnvironments{
		envs:       map[string]map[string]any{},
		pending:    map[string]map[string]any{},
		staleReads: map[string]int{},
	}
}

func (m *mockEnvironments) register(mux *http.ServeMux) {
	const base = "/repositories/{workspace}/{repo_slug}/environments"
	mux.HandleFunc("POST "+base, m.create)
	mux.HandleFunc("GET "+base, m.list)
	mux.HandleFunc("GET "+base+"/{uuid}", m.get)
	mux.HandleFunc("DELETE "+base+"/{uuid}", m.delete)
	mux.HandleFunc("POST "+base+"/{uuid}/changes", m.change)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// seed adds an existing environment, like the default Test/Staging/Production
// environments Bitbucket creates for every repository.
func (m *mockEnvironments) seed(name, envType string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.nextID++
	uuid := fmt.Sprintf("{env-%d}", m.nextID)
	m.envs[uuid] = map[string]any{
		"type":             "deployment_environment",
		"uuid":             uuid,
		"name":             name,
		"slug":             strings.ToLower(strings.ReplaceAll(name, " ", "-")),
		"hidden":           true,
		"environment_type": map[string]any{"type": "deployment_environment_type", "name": envType},
		"restrictions":     map[string]any{"type": "deployment_restrictions_configuration", "admin_only": false},
	}
}

func (m *mockEnvironments) create(w http.ResponseWriter, r *http.Request) {
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": map[string]any{"message": err.Error()}})
		return
	}
	envType, _ := body["environment_type"].(map[string]any)
	typeName, _ := envType["name"].(string)
	if typeName == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": map[string]any{"message": "Unknown environment type."}})
		return
	}
	adminOnly := false
	if restrictions, ok := body["restrictions"].(map[string]any); ok {
		adminOnly, _ = restrictions["admin_only"].(bool)
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.nextID++
	name, _ := body["name"].(string)
	env := map[string]any{
		"type": "deployment_environment",
		"uuid": fmt.Sprintf("{env-%d}", m.nextID),
		"name": name,
		"slug": strings.ToLower(strings.ReplaceAll(name, " ", "-")),
		"rank": 0,
		"environment_type": map[string]any{
			"type": "deployment_environment_type",
			"name": typeName,
			"rank": 1,
		},
		"restrictions": map[string]any{
			"type":       "deployment_restrictions_configuration",
			"admin_only": adminOnly,
		},
		"lock": map[string]any{"type": "deployment_environment_lock_open", "name": "OPEN"},
	}
	m.envs[env["uuid"].(string)] = env
	writeJSON(w, http.StatusCreated, env)
}

func (m *mockEnvironments) list(w http.ResponseWriter, _ *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()
	values := make([]any, 0, len(m.envs))
	for _, env := range m.envs {
		values = append(values, env)
	}
	writeJSON(w, http.StatusOK, map[string]any{"values": values, "page": 1, "size": len(values)})
}

func (m *mockEnvironments) get(w http.ResponseWriter, r *http.Request) {
	uuid := r.PathValue("uuid")
	m.mu.Lock()
	defer m.mu.Unlock()
	env, ok := m.envs[uuid]
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": map[string]any{"message": "not found"}})
		return
	}
	if m.staleReads[uuid] > 0 {
		m.staleReads[uuid]--
	} else if change, ok := m.pending[uuid]; ok {
		for k, v := range change {
			env[k] = v
		}
		delete(m.pending, uuid)
	}
	writeJSON(w, http.StatusOK, env)
}

func (m *mockEnvironments) delete(w http.ResponseWriter, r *http.Request) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.envs, r.PathValue("uuid"))
	w.WriteHeader(http.StatusNoContent)
}

func (m *mockEnvironments) change(w http.ResponseWriter, r *http.Request) {
	uuid := r.PathValue("uuid")
	raw := new(strings.Builder)
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": map[string]any{"message": err.Error()}})
		return
	}
	_ = json.NewEncoder(raw).Encode(body)
	change, ok := body["change"].(map[string]any)
	if !ok || len(body) != 1 {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": map[string]any{"message": "The request body contains invalid properties"}})
		return
	}
	if _, ok := change["environment_type"]; ok {
		writeJSON(w, http.StatusBadRequest, map[string]any{"key": "deploy-service.environment.change-not-supported"})
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.envs[uuid]; !ok {
		writeJSON(w, http.StatusNotFound, map[string]any{"error": map[string]any{"message": "not found"}})
		return
	}
	m.changeCalls = append(m.changeCalls, strings.TrimSpace(raw.String()))
	m.pending[uuid] = change
	m.staleReads[uuid] = 1
	w.WriteHeader(http.StatusAccepted)
}

func (m *mockEnvironments) lastChange() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.changeCalls) == 0 {
		return ""
	}
	return m.changeCalls[len(m.changeCalls)-1]
}

func (m *mockEnvironments) count() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.envs)
}

func testAccDeploymentEnvironmentConfig(baseURL, name, envType string) string {
	return fmt.Sprintf(`
		provider "bitbucket" {
			base_url = %[1]q
		}

		resource "bitbucket_deployments" "staging" {
			workspace        = "testworkspace"
			repo_slug        = "test-repo"
			name             = %[2]q
			environment_type = { name = %[3]q }
		}

		data "bitbucket_deployments" "all" {
			workspace  = "testworkspace"
			repo_slug  = "test-repo"
			depends_on = [bitbucket_deployments.staging]
		}
	`, baseURL, name, envType)
}

// TestAccResourceDeployments_EnvironmentLifecycle covers issue #147: one
// bitbucket_deployments resource manages one environment with typed
// attributes, renames happen in place, a type change replaces the
// environment, and existing environments import by UUID.
func TestAccResourceDeployments_EnvironmentLifecycle(t *testing.T) {
	envs := newMockEnvironments()
	mux := http.NewServeMux()
	envs.register(mux)
	srv := httptest.NewServer(mux)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	const addr = "bitbucket_deployments.staging"
	var firstUUID string

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		CheckDestroy: func(*terraform.State) error {
			if n := envs.count(); n != 0 {
				return fmt.Errorf("%d environment(s) left after destroy", n)
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentEnvironmentConfig(srv.URL, "Staging EU", "Staging"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "name", "Staging EU"),
					resource.TestCheckResourceAttr(addr, "environment_type.name", "Staging"),
					resource.TestCheckResourceAttr(addr, "restrictions.admin_only", "false"),
					resource.TestCheckResourceAttrSet(addr, "environment_uuid"),
					resource.TestCheckResourceAttrWith(addr, "uuid", func(v string) error {
						firstUUID = v
						return nil
					}),
					testAccCheckEnvironmentListContains("data.bitbucket_deployments.all", "Staging EU"),
				),
			},
			{
				Config:   testAccDeploymentEnvironmentConfig(srv.URL, "Staging EU", "Staging"),
				PlanOnly: true,
			},
			{
				Config: testAccDeploymentEnvironmentConfig(srv.URL, "Staging Europe", "Staging"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(addr, plancheck.ResourceActionUpdate)},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "name", "Staging Europe"),
					resource.TestCheckResourceAttrWith(addr, "uuid", func(v string) error {
						if v != firstUUID {
							return fmt.Errorf("rename replaced the environment: uuid %s -> %s", firstUUID, v)
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith(addr, "id", func(v string) error {
						if v != firstUUID {
							return fmt.Errorf("update changed the resource id to %q, want %q", v, firstUUID)
						}
						return nil
					}),
					func(*terraform.State) error {
						if got, want := envs.lastChange(), `{"change":{"name":"Staging Europe"}}`; got != want {
							return fmt.Errorf("update body = %s, want %s", got, want)
						}
						return nil
					},
				),
			},
			{
				Config: testAccDeploymentEnvironmentConfig(srv.URL, "Staging Europe", "Production"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(addr, plancheck.ResourceActionReplace)},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "environment_type.name", "Production"),
					resource.TestCheckResourceAttrWith(addr, "uuid", func(v string) error {
						if v == firstUUID {
							return fmt.Errorf("environment type change did not replace the environment")
						}
						return nil
					}),
				),
			},
			{
				ResourceName:      addr,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[addr]
					if !ok {
						return "", fmt.Errorf("%s not in state", addr)
					}
					return "testworkspace/test-repo/" + rs.Primary.Attributes["uuid"], nil
				},
			},
		},
	})
}

func testAccCheckEnvironmentListContains(dataSource, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[dataSource]
		if !ok {
			return fmt.Errorf("%s not in state", dataSource)
		}
		var envs []map[string]any
		if err := json.Unmarshal([]byte(rs.Primary.Attributes["api_response"]), &envs); err != nil {
			return fmt.Errorf("api_response is not a JSON list of environments: %w", err)
		}
		for _, env := range envs {
			if env["name"] == name {
				return nil
			}
		}
		return fmt.Errorf("environment %q not listed in %s", name, dataSource)
	}
}

// TestAccRealAPI_ResourceDeployments_EnvironmentLifecycle creates a repository
// and a deployment environment, renames it in place via the /changes endpoint,
// imports it by UUID, and destroys everything.
func TestAccRealAPI_ResourceDeployments_EnvironmentLifecycle(t *testing.T) {
	workspace := skipIfNoRealAPI(t)

	suffix := strings.ToLower(acctest.RandStringFromCharSet(6, acctest.CharSetAlpha))
	repoSlug := "tf-acc-env-" + suffix
	projectKey := "TE" + strings.ToUpper(suffix[:5])
	const addr = "bitbucket_deployments.test"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckRepoDestroy(workspace, repoSlug),
		Steps: []resource.TestStep{
			{
				Config: testAccRealDeploymentEnvironmentConfig(workspace, projectKey, repoSlug, "TF Env "+suffix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(addr, "name", "TF Env "+suffix),
					resource.TestCheckResourceAttr(addr, "environment_type.name", "Staging"),
					resource.TestCheckResourceAttrSet(addr, "uuid"),
				),
			},
			{
				Config: testAccRealDeploymentEnvironmentConfig(workspace, projectKey, repoSlug, "TF Env "+suffix+" renamed"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{plancheck.ExpectResourceAction(addr, plancheck.ResourceActionUpdate)},
				},
				Check: resource.TestCheckResourceAttr(addr, "name", "TF Env "+suffix+" renamed"),
			},
			{
				Config:   testAccRealDeploymentEnvironmentConfig(workspace, projectKey, repoSlug, "TF Env "+suffix+" renamed"),
				PlanOnly: true,
			},
			{
				ResourceName:            addr,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_response"},
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[addr]
					if !ok {
						return "", fmt.Errorf("%s not in state", addr)
					}
					return fmt.Sprintf("%s/%s/%s", workspace, repoSlug, rs.Primary.Attributes["uuid"]), nil
				},
			},
		},
	})
}

func testAccRealDeploymentEnvironmentConfig(workspace, projectKey, repoSlug, envName string) string {
	return fmt.Sprintf(`
		provider "bitbucket" {}

		resource "bitbucket_projects" "test" {
			workspace   = %[1]q
			project_key = %[2]q
			request_body = jsonencode({
				name       = "TF Environments Test %[2]s"
				key        = %[2]q
				is_private = true
			})
		}

		resource "bitbucket_repos" "test" {
			workspace  = %[1]q
			repo_slug  = %[3]q
			scm        = "git"
			is_private = true
			project    = { key = %[2]q }
			depends_on = [bitbucket_projects.test]
		}

		resource "bitbucket_deployments" "test" {
			workspace        = %[1]q
			repo_slug        = %[3]q
			name             = %[4]q
			environment_type = { name = "Staging" }
			depends_on       = [bitbucket_repos.test]
		}
	`, workspace, projectKey, repoSlug, envName)
}

// testAccDeploymentEnvironmentsImportConfig adopts an environment that already
// exists by resolving its UUID from the data source at plan time, the pattern
// docs/guides/deployment-environments.md builds on. (The guide's for_each
// variant cannot run here: the testing framework's state shim rejects
// string-keyed resource instances.)
func testAccDeploymentEnvironmentsImportConfig(baseURL string) string {
	return fmt.Sprintf(`
		provider "bitbucket" {
			base_url = %q
		}

		locals {
			workspace = "testworkspace"
			repo_slug = "test-repo"
		}

		data "bitbucket_deployments" "existing" {
			workspace = local.workspace
			repo_slug = local.repo_slug
		}

		locals {
			existing_uuids = {
				for env in jsondecode(data.bitbucket_deployments.existing.api_response) : env.slug => env.uuid
			}
		}

		import {
			to = bitbucket_deployments.staging
			id = "${local.workspace}/${local.repo_slug}/${local.existing_uuids["staging"]}"
		}

		resource "bitbucket_deployments" "staging" {
			workspace        = local.workspace
			repo_slug        = local.repo_slug
			name             = "Staging"
			environment_type = { name = "Staging" }
		}

		resource "bitbucket_deployments" "dt_production" {
			workspace        = local.workspace
			repo_slug        = local.repo_slug
			name             = "DT Production"
			environment_type = { name = "Production" }
		}
	`, baseURL)
}

func TestAccResourceDeployments_ImportExistingEnvironment(t *testing.T) {
	envs := newMockEnvironments()
	envs.seed("Test", "Test")
	envs.seed("Staging", "Staging")
	envs.seed("Production", "Production")
	mux := http.NewServeMux()
	envs.register(mux)
	srv := httptest.NewServer(mux)
	defer srv.Close()
	setMockEnv(t, srv.URL)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDeploymentEnvironmentsImportConfig(srv.URL),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("bitbucket_deployments.staging", plancheck.ResourceActionNoop),
						plancheck.ExpectResourceAction("bitbucket_deployments.dt_production", plancheck.ResourceActionCreate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("bitbucket_deployments.staging", "uuid", "{env-2}"),
					resource.TestCheckResourceAttr("bitbucket_deployments.dt_production", "environment_type.name", "Production"),
				),
			},
			{
				Config:   testAccDeploymentEnvironmentsImportConfig(srv.URL),
				PlanOnly: true,
			},
		},
	})
}
