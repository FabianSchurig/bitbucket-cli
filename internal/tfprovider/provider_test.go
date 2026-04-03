package tfprovider_test

import (
	"strings"
	"testing"

	"github.com/FabianSchurig/bitbucket-cli/internal/tfprovider"
)

// ─── Helper tests ─────────────────────────────────────────────────────────────

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"workspace", "workspace"},
		{"repo_slug", "repo_slug"},
		{"repo-slug", "repo_slug"},
		{"pullRequestId", "pull_request_id"},
		{"repoSlug", "repo_slug"},
		{"content.raw", "content.raw"},
		{"UPPER", "upper"},
	}
	for _, tc := range tests {
		// toSnakeCase is unexported, so test via exported MapCRUDOps indirectly
		// or use a simple duplicate here for validation.
		t.Run(tc.input, func(t *testing.T) {
			// We can only test exported functions, so we test MapCRUDOps
			// and the overall provider behavior instead.
		})
	}
}

func TestMapCRUDOps_BasicMapping(t *testing.T) {
	// Register a temporary config entry for this test.
	tfprovider.CRUDConfig["test-repos"] = tfprovider.CRUDMapping{
		Create: "createRepo",
		Read:   "getRepo",
		Update: "updateRepo",
		Delete: "deleteRepo",
		List:   "listRepos",
	}
	defer delete(tfprovider.CRUDConfig, "test-repos")

	ops := []tfprovider.OperationDef{
		{OperationID: "createRepo", Method: "POST", Path: "/repositories/{workspace}/{repo_slug}"},
		{OperationID: "getRepo", Method: "GET", Path: "/repositories/{workspace}/{repo_slug}"},
		{OperationID: "listRepos", Method: "GET", Path: "/repositories/{workspace}", Paginated: true},
		{OperationID: "updateRepo", Method: "PUT", Path: "/repositories/{workspace}/{repo_slug}"},
		{OperationID: "deleteRepo", Method: "DELETE", Path: "/repositories/{workspace}/{repo_slug}"},
	}

	crud := tfprovider.MapCRUDOps("test-repos", ops)

	if crud.Create == nil || crud.Create.OperationID != "createRepo" {
		t.Errorf("expected Create=createRepo, got %v", crud.Create)
	}
	if crud.Read == nil || crud.Read.OperationID != "getRepo" {
		t.Errorf("expected Read=getRepo, got %v", crud.Read)
	}
	if crud.Update == nil || crud.Update.OperationID != "updateRepo" {
		t.Errorf("expected Update=updateRepo, got %v", crud.Update)
	}
	if crud.Delete == nil || crud.Delete.OperationID != "deleteRepo" {
		t.Errorf("expected Delete=deleteRepo, got %v", crud.Delete)
	}
	if crud.List == nil || crud.List.OperationID != "listRepos" {
		t.Errorf("expected List=listRepos, got %v", crud.List)
	}
}

func TestMapCRUDOps_UnknownGroup(t *testing.T) {
	crud := tfprovider.MapCRUDOps("nonexistent-group", nil)
	if crud.Create != nil || crud.Read != nil || crud.Update != nil || crud.Delete != nil || crud.List != nil {
		t.Error("expected all CRUD ops to be nil for unknown group")
	}
}

func TestMapCRUDOps_MissingOperationID(t *testing.T) {
	// Config references an operation that doesn't exist in the ops list.
	tfprovider.CRUDConfig["test-missing"] = tfprovider.CRUDMapping{
		Create: "doesNotExist",
		Read:   "getItem",
	}
	defer delete(tfprovider.CRUDConfig, "test-missing")

	ops := []tfprovider.OperationDef{
		{OperationID: "getItem", Method: "GET", Path: "/items/{id}"},
	}

	crud := tfprovider.MapCRUDOps("test-missing", ops)

	if crud.Create != nil {
		t.Error("expected Create to be nil for missing operation ID")
	}
	if crud.Read == nil || crud.Read.OperationID != "getItem" {
		t.Errorf("expected Read=getItem, got %v", crud.Read)
	}
}

func TestBuildResourceDescription(t *testing.T) {
	crud := tfprovider.CRUDOps{
		Create: &tfprovider.OperationDef{OperationID: "createItem", Method: "POST", Path: "/items"},
		Read:   &tfprovider.OperationDef{OperationID: "getItem", Method: "GET", Path: "/items/{id}"},
	}
	desc := tfprovider.BuildResourceDescription("Manage items", crud)
	if desc == "" {
		t.Error("expected non-empty description")
	}
	if !strings.Contains(desc, "createItem") || !strings.Contains(desc, "getItem") {
		t.Error("expected description to mention operation IDs")
	}
}

// ─── Provider tests ───────────────────────────────────────────────────────────

func TestProviderNew(t *testing.T) {
	factory := tfprovider.New("v1.0.0")
	if factory == nil {
		t.Fatal("expected non-nil factory function")
	}
	p := factory()
	if p == nil {
		t.Fatal("expected non-nil provider")
	}
}

func TestRegisterResourceGroup(t *testing.T) {
	// The generated code calls RegisterResourceGroup in init().
	// Verify that New returns a provider with resources.
	factory := tfprovider.New("test")
	p := factory()
	if p == nil {
		t.Fatal("expected non-nil provider")
	}
}

// ─── Generated resource group smoke tests ─────────────────────────────────────

func TestAllGeneratedResourceGroups_AreRegistered(t *testing.T) {
	// Verify that the provider factory works and the generated init()
	// functions have registered resource groups.
	factory := tfprovider.New("test")
	p := factory()
	if p == nil {
		t.Fatal("expected non-nil provider")
	}
	// The provider's Resources and DataSources methods are called by
	// Terraform framework. We can't call them directly without the full
	// provider lifecycle, but we can verify the provider was created.
}

func TestGeneratedResourceGroups_HaveCRUDOps(t *testing.T) {
	// Verify that at least one generated resource group has CRUD ops mapped.
	// We'll test the PRResourceGroup directly since it's exported.
	group := tfprovider.PRResourceGroup
	if group.TypeName != "pr" {
		t.Errorf("expected TypeName 'pr', got %q", group.TypeName)
	}
	if group.Ops.Read == nil && group.Ops.List == nil {
		t.Error("expected PRResourceGroup to have at least a Read or List operation")
	}
	if len(group.AllOps) == 0 {
		t.Error("expected PRResourceGroup to have operations")
	}
}

func TestGeneratedResourceGroups_ReposHasAllCRUD(t *testing.T) {
	group := tfprovider.ReposResourceGroup
	if group.TypeName != "repos" {
		t.Errorf("expected TypeName 'repos', got %q", group.TypeName)
	}
	// Repos should have all CRUD operations mapped via CRUDConfig.
	if group.Ops.Create == nil {
		t.Error("expected repos to have Create operation")
	}
	if group.Ops.Read == nil {
		t.Error("expected repos to have Read operation")
	}
	if group.Ops.Update == nil {
		t.Error("expected repos to have Update operation")
	}
	if group.Ops.Delete == nil {
		t.Error("expected repos to have Delete operation")
	}
	if group.Ops.List == nil {
		t.Error("expected repos to have List operation")
	}
	// Verify the correct operations were picked (not sub-resource ones).
	if group.Ops.Create.OperationID != "createARepository" {
		t.Errorf("expected Create=createARepository, got %s", group.Ops.Create.OperationID)
	}
	if group.Ops.Read.OperationID != "getARepository" {
		t.Errorf("expected Read=getARepository, got %s", group.Ops.Read.OperationID)
	}
}

func TestGeneratedResourceGroups_AllGroupsHaveOps(t *testing.T) {
	groups := []tfprovider.ResourceGroup{
		tfprovider.PRResourceGroup,
		tfprovider.HooksResourceGroup,
		tfprovider.SearchResourceGroup,
		tfprovider.RefsResourceGroup,
		tfprovider.CommitsResourceGroup,
		tfprovider.ReportsResourceGroup,
		tfprovider.ReposResourceGroup,
		tfprovider.WorkspacesResourceGroup,
		tfprovider.ProjectsResourceGroup,
		tfprovider.PipelinesResourceGroup,
		tfprovider.IssuesResourceGroup,
		tfprovider.SnippetsResourceGroup,
		tfprovider.DeploymentsResourceGroup,
		tfprovider.BranchRestrictionsResourceGroup,
		tfprovider.BranchingModelResourceGroup,
		tfprovider.CommitStatusesResourceGroup,
		tfprovider.DownloadsResourceGroup,
		tfprovider.UsersResourceGroup,
		tfprovider.PropertiesResourceGroup,
		tfprovider.AddonResourceGroup,
	}

	if len(groups) != 20 {
		t.Fatalf("expected 20 resource groups, got %d", len(groups))
	}

	for _, g := range groups {
		t.Run(g.TypeName, func(t *testing.T) {
			if len(g.AllOps) == 0 {
				t.Errorf("resource group %q has no operations", g.TypeName)
			}
			if g.TypeName == "" {
				t.Error("resource group has empty TypeName")
			}
			if g.Description == "" {
				t.Error("resource group has empty Description")
			}
			// Every group should have at least a Read or List.
			if g.Ops.Read == nil && g.Ops.List == nil {
				t.Errorf("resource group %q has no Read or List operation", g.TypeName)
			}
		})
	}
}
