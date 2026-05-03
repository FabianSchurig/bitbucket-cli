package tfprovider_test

import (
	"strings"
	"testing"

	"github.com/FabianSchurig/bitbucket-cli/internal/tfprovider"
)

const (
	testReposGroupName          = "test-repos"
	testMissingGroupName        = "test-missing"
	repoPathTemplate            = "/repositories/{workspace}/{repo_slug}"
	expectedNonNilProviderError = "expected non-nil provider"
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
		{"content.raw", "content_raw"},
		{"UPPER", "upper"},
	}
	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			// Test via ParamAttrName which uses toSnakeCase internally.
			// For non-"id" params, ParamAttrName == toSnakeCase.
			got := tfprovider.ParamAttrName(tc.input)
			if got != tc.expected {
				t.Errorf("ParamAttrName(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

func TestParamAttrName_IDCollision(t *testing.T) {
	// The API param "id" should be remapped to "param_id" to avoid collision
	// with Terraform's computed "id" attribute.
	got := tfprovider.ParamAttrName("id")
	if got != "param_id" {
		t.Errorf("ParamAttrName(\"id\") = %q, want \"param_id\"", got)
	}
}

func TestMapCRUDOps_BasicMapping(t *testing.T) {
	// Register a temporary config entry for this test.
	tfprovider.CRUDConfig[testReposGroupName] = tfprovider.CRUDMapping{
		Create: "createRepo",
		Read:   "getRepo",
		Update: "updateRepo",
		Delete: "deleteRepo",
		List:   "listRepos",
	}
	defer delete(tfprovider.CRUDConfig, testReposGroupName)

	ops := []tfprovider.OperationDef{
		{OperationID: "createRepo", Method: "POST", Path: repoPathTemplate},
		{OperationID: "getRepo", Method: "GET", Path: repoPathTemplate},
		{OperationID: "listRepos", Method: "GET", Path: "/repositories/{workspace}", Paginated: true},
		{OperationID: "updateRepo", Method: "PUT", Path: repoPathTemplate},
		{OperationID: "deleteRepo", Method: "DELETE", Path: repoPathTemplate},
	}

	crud := tfprovider.MapCRUDOps(testReposGroupName, ops)

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
	tfprovider.CRUDConfig[testMissingGroupName] = tfprovider.CRUDMapping{
		Create: "doesNotExist",
		Read:   "getItem",
	}
	defer delete(tfprovider.CRUDConfig, testMissingGroupName)

	ops := []tfprovider.OperationDef{
		{OperationID: "getItem", Method: "GET", Path: "/items/{id}"},
	}

	crud := tfprovider.MapCRUDOps(testMissingGroupName, ops)

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
	assertNonNilProvider(t, factory())
}

func TestRegisterResourceGroup(t *testing.T) {
	// The generated code calls RegisterResourceGroup in init().
	// Verify that New returns a provider with resources.
	assertNonNilProvider(t, tfprovider.New("test")())
}

// ─── Generated resource group smoke tests ─────────────────────────────────────

func TestAllGeneratedResourceGroups_AreRegistered(t *testing.T) {
	// Verify that the provider factory works and the generated init()
	// functions have registered resource groups.
	assertNonNilProvider(t, tfprovider.New("test")())
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

// ─── Sub-resource group tests ─────────────────────────────────────────────────

func TestSubResourceGroups_Registered(t *testing.T) {
	// Verify that sub-resource CRUDConfig entries resolve correctly
	// against the parent groups' operations.
	subResources := map[string]struct {
		wantRead   string
		wantCreate string
		wantList   string
	}{
		"workspace-hooks": {
			wantCreate: "createAWebhookForAWorkspace",
			wantRead:   "getAWebhookForAWorkspace",
			wantList:   "listWebhooksForAWorkspace",
		},
		"default-reviewers": {
			wantRead:   "getADefaultReviewer",
			wantCreate: "addAUserToTheDefaultReviewers",
			wantList:   "listDefaultReviewers",
		},
		"project-default-reviewers": {
			wantRead:   "getWorkspacesProjectsDefault-Reviewers",
			wantCreate: "addTheSpecificUserAsADefaultReviewerForTheProject",
			wantList:   "listTheDefaultReviewersInAProject",
		},
		"pipeline-variables": {
			wantCreate: "createRepositoryPipelineVariable",
			wantRead:   "getRepositoryPipelineVariable",
			wantList:   "getRepositoryPipelineVariables",
		},
		"workspace-pipeline-variables": {
			wantCreate: "createPipelineVariableForWorkspace",
			wantRead:   "getPipelineVariableForWorkspace",
			wantList:   "getPipelineVariablesForWorkspace",
		},
		"deployment-variables": {
			wantCreate: "createDeploymentVariable",
			wantRead:   "getDeploymentVariables",
		},
		"repo-group-permissions": {
			wantRead: "getAnExplicitGroupPermissionForARepository",
			wantList: "listExplicitGroupPermissionsForARepository",
		},
		"repo-user-permissions": {
			wantRead: "getAnExplicitUserPermissionForARepository",
			wantList: "listExplicitUserPermissionsForARepository",
		},
		"project-group-permissions": {
			wantRead: "getAnExplicitGroupPermissionForAProject",
			wantList: "listExplicitGroupPermissionsForAProject",
		},
		"project-user-permissions": {
			wantRead: "getAnExplicitUserPermissionForAProject",
			wantList: "listExplicitUserPermissionsForAProject",
		},
		"repo-deploy-keys": {
			wantCreate: "addARepositoryDeployKey",
			wantRead:   "getARepositoryDeployKey",
			wantList:   "listRepositoryDeployKeys",
		},
		"project-deploy-keys": {
			wantCreate: "createAProjectDeployKey",
			wantRead:   "getAProjectDeployKey",
			wantList:   "listProjectDeployKeys",
		},
		// ─── Wave 2: additional sub-resources ─────────────────────────────────
		"tags": {
			wantCreate: "createATag",
			wantRead:   "getATag",
			wantList:   "listTags",
		},
		"pipeline-ssh-keys": {
			wantRead: "getRepositoryPipelineSshKeyPair",
		},
		"pipeline-known-hosts": {
			wantCreate: "createRepositoryPipelineKnownHost",
			wantRead:   "getRepositoryPipelineKnownHost",
			wantList:   "getRepositoryPipelineKnownHosts",
		},
		"pipeline-schedules": {
			wantCreate: "createRepositoryPipelineSchedule",
			wantRead:   "getRepositoryPipelineSchedule",
			wantList:   "getRepositoryPipelineSchedules",
		},
		"pipeline-config": {
			wantRead: "getRepositoryPipelineConfig",
		},
		"ssh-keys": {
			wantCreate: "addANewSshKey",
			wantRead:   "getASshKey",
			wantList:   "listSshKeys",
		},
		"current-user": {
			wantRead: "getCurrentUser",
		},
		"forked-repository": {
			wantCreate: "forkARepository",
			wantList:   "listRepositoryForks",
		},
		"project-branching-model": {
			wantRead: "getTheBranchingModelForAProject",
		},
		"pipeline-oidc": {
			wantRead: "getOIDCConfiguration",
		},
		"pipeline-oidc-keys": {
			wantRead: "getOIDCKeys",
		},
		"workspace-members": {
			wantRead: "getUserMembershipForAWorkspace",
			wantList: "listUsersInAWorkspace",
		},
		"annotations": {
			wantCreate: "createOrUpdateAnnotation",
			wantRead:   "getAnnotation",
			wantList:   "getAnnotationsForReport",
		},
		"commit-file": {
			wantCreate: "createACommitByUploadingAFile",
			wantRead:   "getFileOrDirectoryContents",
		},
		"pr-comments": {
			wantCreate: "createACommentOnAPullRequest",
			wantRead:   "getACommentOnAPullRequest",
			wantList:   "listCommentsOnAPullRequest",
		},
		"issue-comments": {
			wantCreate: "createACommentOnAnIssue",
			wantRead:   "getACommentOnAnIssue",
			wantList:   "listCommentsOnAnIssue",
		},
	}

	for typeName, expected := range subResources {
		t.Run(typeName, func(t *testing.T) {
			cfg, ok := tfprovider.CRUDConfig[typeName]
			if !ok {
				t.Fatalf("CRUDConfig missing entry for %q", typeName)
			}

			ops := subResourceCRUDOps(typeName)
			assertExpectedOperation(t, "Read", expected.wantRead, ops.Read)
			assertExpectedOperation(t, "Create", expected.wantCreate, ops.Create)
			assertExpectedOperation(t, "List", expected.wantList, ops.List)
			_ = cfg // cfg verified via MapCRUDOps
		})
	}
}

func assertNonNilProvider(t *testing.T, provider any) {
	t.Helper()
	if provider == nil {
		t.Fatal(expectedNonNilProviderError)
	}
}

func subResourceCRUDOps(typeName string) tfprovider.CRUDOps {
	for _, parent := range subResourceParents() {
		candidate := tfprovider.MapCRUDOps(typeName, parent.AllOps)
		if candidate.Read != nil || candidate.List != nil || candidate.Create != nil {
			return candidate
		}
	}
	return tfprovider.CRUDOps{}
}

func subResourceParents() []tfprovider.ResourceGroup {
	return []tfprovider.ResourceGroup{
		tfprovider.WorkspacesResourceGroup,
		tfprovider.PRResourceGroup,
		tfprovider.ProjectsResourceGroup,
		tfprovider.PipelinesResourceGroup,
		tfprovider.ReposResourceGroup,
		tfprovider.DeploymentsResourceGroup,
		tfprovider.RefsResourceGroup,
		tfprovider.UsersResourceGroup,
		tfprovider.BranchingModelResourceGroup,
		tfprovider.ReportsResourceGroup,
		tfprovider.IssuesResourceGroup,
	}
}

func assertExpectedOperation(t *testing.T, label, expected string, op *tfprovider.OperationDef) {
	t.Helper()
	if expected == "" {
		return
	}
	if op == nil {
		t.Errorf("expected %s=%s, got nil", label, expected)
		return
	}
	if op.OperationID != expected {
		t.Errorf("expected %s=%s, got %s", label, expected, op.OperationID)
	}
}
