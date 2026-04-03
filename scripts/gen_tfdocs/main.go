// gen_tfdocs reads the hand-written CRUD config and all generated *.gen.go
// files to produce:
//   - docs/index.md              (provider documentation)
//   - docs/resources/<name>.md   (one per resource group)
//   - docs/data-sources/<name>.md (one per data source group)
//   - examples/provider/provider.tf
//   - examples/resources/<name>/resource.tf
//   - examples/data-sources/<name>/data-source.tf
//   - tests/<name>.tftest.hcl    (one per resource group)
//
// Usage: go run scripts/gen_tfdocs/main.go
//
// This follows the Terraform Registry documentation structure:
// https://developer.hashicorp.com/terraform/registry/providers/docs
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

// ─── CRUD config (mirrored from internal/tfprovider/crud_config.go) ───────────

type CRUDMapping struct {
	Create string
	Read   string
	Update string
	Delete string
	List   string
}

// crudConfig mirrors the hand-written config in internal/tfprovider/crud_config.go.
// We duplicate it here so the generator can run without importing internal packages.
var crudConfig = map[string]CRUDMapping{
	"repos": {
		Create: "createARepository",
		Read:   "getARepository",
		Update: "updateARepository",
		Delete: "deleteARepository",
		List:   "listRepositoriesInAWorkspace",
	},
	"pr": {
		Create: "createAPullRequest",
		Read:   "getAPullRequest",
		Update: "updateAPullRequest",
		List:   "listPullRequests",
	},
	"projects": {
		Create: "createAProjectInAWorkspace",
		Read:   "getAProjectForAWorkspace",
		Update: "updateAProjectForAWorkspace",
		Delete: "deleteAProjectForAWorkspace",
		List:   "listProjectsInAWorkspace",
	},
	"workspaces": {
		Read: "getAWorkspace",
		List: "listWorkspacesForUser",
	},
	"issues": {
		Create: "createAnIssue",
		Read:   "getAnIssue",
		Update: "updateAnIssue",
		Delete: "deleteAnIssue",
		List:   "listIssues",
	},
	"hooks": {
		Create: "createAWebhookForARepository",
		Read:   "getAWebhookForARepository",
		Update: "updateAWebhookForARepository",
		Delete: "deleteAWebhookForARepository",
		List:   "listWebhooksForARepository",
	},
	"snippets": {
		Create: "createASnippet",
		Read:   "getASnippet",
		Update: "updateASnippet",
		Delete: "deleteASnippet",
		List:   "listSnippets",
	},
	"refs": {
		Create: "createABranch",
		Read:   "getABranch",
		Delete: "deleteABranch",
		List:   "listBranchesAndTags",
	},
	"commits": {
		Read: "getACommit",
		List: "listCommits",
	},
	"pipelines": {
		Create: "createPipelineForRepository",
		Read:   "getPipelineForRepository",
		List:   "getPipelinesForRepository",
	},
	"deployments": {
		Create: "createEnvironment",
		Read:   "getEnvironmentForRepository",
		Delete: "deleteEnvironmentForRepository",
		List:   "getEnvironmentsForRepository",
	},
	"branch-restrictions": {
		Create: "createABranchRestrictionRule",
		Read:   "getABranchRestrictionRule",
		Update: "updateABranchRestrictionRule",
		Delete: "deleteABranchRestrictionRule",
		List:   "listBranchRestrictions",
	},
	"branching-model": {
		Read:   "getTheBranchingModelForARepository",
		Update: "updateTheBranchingModelConfigForARepository",
	},
	"commit-statuses": {
		Create: "createABuildStatusForACommit",
		Read:   "getABuildStatusForACommit",
		Update: "updateABuildStatusForACommit",
		List:   "listCommitStatusesForACommit",
	},
	"downloads": {
		Create: "uploadADownloadArtifact",
		Read:   "getADownloadArtifactLink",
		Delete: "deleteADownloadArtifact",
		List:   "listDownloadArtifacts",
	},
	"users": {
		Read: "getAUser",
		List: "listSshKeys",
	},
	"reports": {
		Create: "createOrUpdateReport",
		Read:   "getReport",
		Delete: "deleteReport",
		List:   "getReportsForCommit",
	},
	"search": {
		List: "searchWorkspace",
	},
	"properties": {
		Read:   "getRepositoryHostedPropertyValue",
		Update: "updateRepositoryHostedPropertyValue",
		Delete: "deleteRepositoryHostedPropertyValue",
	},
	"addon": {
		Update: "updateAnInstalledApp",
		Delete: "deleteAnApp",
		List:   "listLinkersForAnApp",
	},
	// ─── Sub-resource CRUD mappings ───────────────────────────────────────────
	"workspace-hooks": {
		Create: "createAWebhookForAWorkspace",
		Read:   "getAWebhookForAWorkspace",
		Update: "updateAWebhookForAWorkspace",
		Delete: "deleteAWebhookForAWorkspace",
		List:   "listWebhooksForAWorkspace",
	},
	"default-reviewers": {
		Read:   "getADefaultReviewer",
		Create: "addAUserToTheDefaultReviewers",
		Delete: "removeAUserFromTheDefaultReviewers",
		List:   "listDefaultReviewers",
	},
	"project-default-reviewers": {
		Read:   "getWorkspacesProjectsDefault-Reviewers",
		Create: "addTheSpecificUserAsADefaultReviewerForTheProject",
		Delete: "removeTheSpecificUserFromTheProjectsDefaultReviewers",
		List:   "listTheDefaultReviewersInAProject",
	},
	"pipeline-variables": {
		Create: "createRepositoryPipelineVariable",
		Read:   "getRepositoryPipelineVariable",
		Update: "updateRepositoryPipelineVariable",
		Delete: "deleteRepositoryPipelineVariable",
		List:   "getRepositoryPipelineVariables",
	},
	"workspace-pipeline-variables": {
		Create: "createPipelineVariableForWorkspace",
		Read:   "getPipelineVariableForWorkspace",
		Update: "updatePipelineVariableForWorkspace",
		Delete: "deletePipelineVariableForWorkspace",
		List:   "getPipelineVariablesForWorkspace",
	},
	"deployment-variables": {
		Create: "createDeploymentVariable",
		Read:   "getDeploymentVariables",
		Update: "updateDeploymentVariable",
		Delete: "deleteDeploymentVariable",
	},
	"repo-group-permissions": {
		Read:   "getAnExplicitGroupPermissionForARepository",
		Update: "updateAnExplicitGroupPermissionForARepository",
		Delete: "deleteAnExplicitGroupPermissionForARepository",
		List:   "listExplicitGroupPermissionsForARepository",
	},
	"repo-user-permissions": {
		Read:   "getAnExplicitUserPermissionForARepository",
		Update: "updateAnExplicitUserPermissionForARepository",
		Delete: "deleteAnExplicitUserPermissionForARepository",
		List:   "listExplicitUserPermissionsForARepository",
	},
	"project-group-permissions": {
		Read:   "getAnExplicitGroupPermissionForAProject",
		Update: "updateAnExplicitGroupPermissionForAProject",
		Delete: "deleteAnExplicitGroupPermissionForAProject",
		List:   "listExplicitGroupPermissionsForAProject",
	},
	"project-user-permissions": {
		Read:   "getAnExplicitUserPermissionForAProject",
		Update: "updateAnExplicitUserPermissionForAProject",
		Delete: "deleteAnExplicitUserPermissionForAProject",
		List:   "listExplicitUserPermissionsForAProject",
	},
	"repo-deploy-keys": {
		Create: "addARepositoryDeployKey",
		Read:   "getARepositoryDeployKey",
		Update: "updateARepositoryDeployKey",
		Delete: "deleteARepositoryDeployKey",
		List:   "listRepositoryDeployKeys",
	},
	"project-deploy-keys": {
		Create: "createAProjectDeployKey",
		Read:   "getAProjectDeployKey",
		Delete: "deleteADeployKeyFromAProject",
		List:   "listProjectDeployKeys",
	},
	// ─── Wave 2: additional sub-resource CRUD mappings ────────────────────────
	"tags": {
		Create: "createATag",
		Read:   "getATag",
		Delete: "deleteATag",
		List:   "listTags",
	},
	"pipeline-ssh-keys": {
		Read:   "getRepositoryPipelineSshKeyPair",
		Update: "updateRepositoryPipelineKeyPair",
		Delete: "deleteRepositoryPipelineKeyPair",
	},
	"pipeline-known-hosts": {
		Create: "createRepositoryPipelineKnownHost",
		Read:   "getRepositoryPipelineKnownHost",
		Update: "updateRepositoryPipelineKnownHost",
		Delete: "deleteRepositoryPipelineKnownHost",
		List:   "getRepositoryPipelineKnownHosts",
	},
	"pipeline-schedules": {
		Create: "createRepositoryPipelineSchedule",
		Read:   "getRepositoryPipelineSchedule",
		Update: "updateRepositoryPipelineSchedule",
		Delete: "deleteRepositoryPipelineSchedule",
		List:   "getRepositoryPipelineSchedules",
	},
	"pipeline-config": {
		Read:   "getRepositoryPipelineConfig",
		Update: "updateRepositoryPipelineConfig",
	},
	"ssh-keys": {
		Create: "addANewSshKey",
		Read:   "getASshKey",
		Update: "updateASshKey",
		Delete: "deleteASshKey",
		List:   "listSshKeys",
	},
	"current-user": {
		Read: "getCurrentUser",
	},
	"forked-repository": {
		Create: "forkARepository",
		List:   "listRepositoryForks",
	},
	"project-branching-model": {
		Read:   "getTheBranchingModelForAProject",
		Update: "updateTheBranchingModelConfigForAProject",
	},
	"pipeline-oidc": {
		Read: "getOIDCConfiguration",
	},
	"pipeline-oidc-keys": {
		Read: "getOIDCKeys",
	},
	"workspace-members": {
		Read: "getUserMembershipForAWorkspace",
		List: "listUsersInAWorkspace",
	},
	"annotations": {
		Create: "createOrUpdateAnnotation",
		Read:   "getAnnotation",
		Delete: "deleteAnnotation",
		List:   "getAnnotationsForReport",
	},
	"commit-file": {
		Create: "createACommitByUploadingAFile",
		Read:   "getFileOrDirectoryContents",
	},
	"pr-comments": {
		Create: "createACommentOnAPullRequest",
		Read:   "getACommentOnAPullRequest",
		Update: "updateACommentOnAPullRequest",
		Delete: "deleteACommentOnAPullRequest",
		List:   "listCommentsOnAPullRequest",
	},
	"issue-comments": {
		Create: "createACommentOnAnIssue",
		Read:   "getACommentOnAnIssue",
		Update: "updateACommentOnAnIssue",
		Delete: "deleteACommentOnAnIssue",
		List:   "listCommentsOnAnIssue",
	},
}

// ─── Param info per resource group (required path params for primary Read op) ─

// paramConfig maps resource groups to the required path parameters of their primary Read operation.
// This determines what attributes appear in examples and tests.
var paramConfig = map[string][]string{
	"repos":               {"workspace", "repo_slug"},
	"pr":                  {"workspace", "repo_slug", "pull_request_id"},
	"projects":            {"workspace", "project_key"},
	"workspaces":          {"workspace"},
	"issues":              {"workspace", "repo_slug", "issue_id"},
	"hooks":               {"workspace", "repo_slug", "uid"},
	"snippets":            {"workspace", "encoded_id"},
	"refs":                {"workspace", "repo_slug", "name"},
	"commits":             {"workspace", "repo_slug", "commit"},
	"pipelines":           {"workspace", "repo_slug", "pipeline_uuid"},
	"deployments":         {"workspace", "repo_slug", "environment_uuid"},
	"branch-restrictions": {"workspace", "repo_slug", "param_id"},
	"branching-model":     {"workspace", "repo_slug"},
	"commit-statuses":     {"workspace", "repo_slug", "commit", "key"},
	"downloads":           {"workspace", "repo_slug", "filename"},
	"users":               {"selected_user"},
	"reports":             {"workspace", "repo_slug", "commit", "report_id"},
	"search":              {"workspace"},
	"properties":          {"workspace", "repo_slug", "app_key", "property_name"},
	"addon":               {},
	// ─── Sub-resource params ──────────────────────────────────────────────────
	"workspace-hooks":              {"workspace", "uid"},
	"default-reviewers":            {"workspace", "repo_slug", "target_username"},
	"project-default-reviewers":    {"workspace", "project_key", "selected_user"},
	"pipeline-variables":           {"workspace", "repo_slug", "variable_uuid"},
	"workspace-pipeline-variables": {"workspace", "variable_uuid"},
	"deployment-variables":         {"workspace", "repo_slug", "environment_uuid"},
	"repo-group-permissions":       {"workspace", "repo_slug", "group_slug"},
	"repo-user-permissions":        {"workspace", "repo_slug", "selected_user_id"},
	"project-group-permissions":    {"workspace", "project_key", "group_slug"},
	"project-user-permissions":     {"workspace", "project_key", "selected_user_id"},
	"repo-deploy-keys":             {"workspace", "repo_slug", "key_id"},
	"project-deploy-keys":          {"workspace", "project_key", "key_id"},
	// ─── Wave 2: additional sub-resource params ──────────────────────────────
	"tags":                          {"workspace", "repo_slug", "name"},
	"pipeline-ssh-keys":             {"workspace", "repo_slug"},
	"pipeline-known-hosts":          {"workspace", "repo_slug", "known_host_uuid"},
	"pipeline-schedules":            {"workspace", "repo_slug", "schedule_uuid"},
	"pipeline-config":               {"workspace", "repo_slug"},
	"ssh-keys":                      {"selected_user", "key_id"},
	"current-user":                  {},
	"forked-repository":             {"workspace", "repo_slug"},
	"project-branching-model":       {"workspace", "project_key"},
	"pipeline-oidc":                 {"workspace"},
	"pipeline-oidc-keys":            {"workspace"},
	"workspace-members":             {"workspace", "member"},
	"annotations":                   {"workspace", "repo_slug", "commit", "reportId", "annotationId"},
	"commit-file":                   {"workspace", "repo_slug", "commit", "path"},
	"pr-comments":                   {"workspace", "repo_slug", "pull_request_id", "comment_id"},
	"issue-comments":                {"workspace", "repo_slug", "issue_id", "comment_id"},
}

// ─── Template data ────────────────────────────────────────────────────────────

type GroupData struct {
	Name        string
	TFName      string // e.g., "bitbucket_repos"
	HasCreate   bool
	HasRead     bool
	HasUpdate   bool
	HasDelete   bool
	HasList     bool
	HasIDParam  bool // true if "id" is a path parameter (avoids conflict with computed id)
	Params      []string
	ParamValues map[string]string
}

func exampleValue(param string) string {
	switch param {
	case "workspace":
		return "my-workspace"
	case "repo_slug":
		return "my-repo"
	case "pull_request_id":
		return "1"
	case "project_key":
		return "PROJ"
	case "issue_id":
		return "1"
	case "uid":
		return "webhook-uuid"
	case "encoded_id":
		return "snippet-id"
	case "name":
		return "main"
	case "commit":
		return "abc123def"
	case "pipeline_uuid":
		return "pipeline-uuid"
	case "environment_uuid":
		return "env-uuid"
	case "param_id":
		return "1"
	case "key":
		return "build-key"
	case "filename":
		return "artifact.zip"
	case "selected_user":
		return "jdoe"
	case "report_id":
		return "report-uuid"
	case "app_key":
		return "my-app"
	case "property_name":
		return "my-property"
	case "target_username":
		return "jdoe"
	case "variable_uuid":
		return "{variable-uuid}"
	case "group_slug":
		return "developers"
	case "selected_user_id":
		return "{user-uuid}"
	case "key_id":
		return "123"
	default:
		return "example-value"
	}
}

func buildGroups() []GroupData {
	var groups []GroupData
	for name, crud := range crudConfig {
		params := paramConfig[name]
		pv := make(map[string]string)
		hasIDParam := false
		for _, p := range params {
			pv[p] = exampleValue(p)
			if p == "param_id" {
				hasIDParam = true
			}
		}
		groups = append(groups, GroupData{
			Name:        name,
			TFName:      "bitbucket_" + strings.ReplaceAll(name, "-", "_"),
			HasCreate:   crud.Create != "",
			HasRead:     crud.Read != "",
			HasUpdate:   crud.Update != "",
			HasDelete:   crud.Delete != "",
			HasList:     crud.List != "",
			HasIDParam:  hasIDParam,
			Params:      params,
			ParamValues: pv,
		})
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i].Name < groups[j].Name })
	return groups
}

// ─── Templates ────────────────────────────────────────────────────────────────

var funcMap = template.FuncMap{
	"replace": strings.ReplaceAll,
	"snakeCase": func(s string) string {
		return strings.ReplaceAll(s, "-", "_")
	},
}

const providerDocTemplate = `---
page_title: "bitbucket Provider"
subcategory: ""
description: |-
  Terraform provider for Bitbucket Cloud. Auto-generated from the Bitbucket OpenAPI spec.
---

# bitbucket Provider

Terraform provider for Bitbucket Cloud, exposing all Bitbucket API operations as
generic resources and data sources. Auto-generated from the Bitbucket OpenAPI spec.

## Authentication

The provider authenticates via HTTP Basic Auth using an Atlassian API token.
Create a token at [id.atlassian.com/manage-profile/security/api-tokens](https://id.atlassian.com/manage-profile/security/api-tokens).

### Atlassian API Token (recommended)

` + "```" + `hcl
provider "bitbucket" {
  username = "your-email@example.com"  # Atlassian account email
  token    = "your-api-token"
}
` + "```" + `

Or via environment variables:

` + "```" + `bash
export BITBUCKET_USERNAME="your-email@example.com"
export BITBUCKET_TOKEN="your-api-token"
` + "```" + `

### Workspace Access Token

For workspace/repository access tokens, only the token is needed:

` + "```" + `hcl
provider "bitbucket" {
  token = "your-workspace-access-token"
}
` + "```" + `

## Example Usage

` + "```" + `hcl
terraform {
  required_providers {
    bitbucket = {
      source = "FabianSchurig/bitbucket"
    }
  }
}

provider "bitbucket" {
  # Authentication via environment variables recommended
}

# Read a repository
data "bitbucket_repos" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}

# Output the API response
output "repo_info" {
  value = data.bitbucket_repos.example.api_response
}
` + "```" + `

## Schema

### Optional

- ` + "`username`" + ` (String) Bitbucket username (Atlassian account email for API tokens). Can also be set via ` + "`BITBUCKET_USERNAME`" + ` environment variable.
- ` + "`token`" + ` (String, Sensitive) Bitbucket API token (Atlassian API token or workspace access token). Can also be set via ` + "`BITBUCKET_TOKEN`" + ` environment variable.
- ` + "`base_url`" + ` (String) Base URL for the Bitbucket API. Defaults to ` + "`https://api.bitbucket.org/2.0`" + `.

## Resources and Data Sources

This provider auto-generates resources and data sources for all Bitbucket API
operation groups. Each resource group maps to a set of CRUD operations.

| Resource | Data Source | CRUD |
|----------|-------------|------|
{{- range .Groups}}
| ` + "`" + `{{.TFName}}` + "`" + ` | ` + "`" + `{{.TFName}}` + "`" + ` | {{if .HasCreate}}C{{end}}{{if .HasRead}}R{{end}}{{if .HasUpdate}}U{{end}}{{if .HasDelete}}D{{end}}{{if .HasList}}L{{end}} |
{{- end}}

All resources share the same generic schema pattern:

- **Path parameters** become required/optional string attributes
- **Body fields** become optional string attributes
- ` + "`api_response`" + ` (Computed) contains the raw JSON API response
- ` + "`id`" + ` (Computed) is extracted from the response (uuid, id, slug, or name)
`

const resourceDocTemplate = `---
page_title: "{{.TFName}} Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket {{.Name}} via the Bitbucket Cloud API.
---

# {{.TFName}} (Resource)

Manages Bitbucket {{.Name}} via the Bitbucket Cloud API.

## CRUD Operations

{{- if .HasCreate}}
- **Create**: Supported
{{- end}}
{{- if .HasRead}}
- **Read**: Supported
{{- end}}
{{- if .HasUpdate}}
- **Update**: Supported
{{- end}}
{{- if .HasDelete}}
- **Delete**: Supported
{{- end}}
{{- if .HasList}}
- **List**: Supported (via data source)
{{- end}}

## Example Usage

` + "```" + `hcl
resource "{{.TFName}}" "example" {
{{- range .Params}}
  {{.}} = "{{index $.ParamValues .}}"
{{- end}}
}
` + "```" + `

## Schema

### Required

{{- range .Params}}
- ` + "`" + `{{.}}` + "`" + ` (String) Path parameter.
{{- end}}

### Optional

- ` + "`" + `request_body` + "`" + ` (String) Raw JSON request body for create/update operations. Use ` + "`" + `jsonencode({...})` + "`" + ` to pass fields not exposed as individual attributes.

### Read-Only
{{- if not .HasIDParam}}

- ` + "`" + `id` + "`" + ` (String) Resource identifier (extracted from API response).
{{- end}}
- ` + "`" + `api_response` + "`" + ` (String) The raw JSON response from the Bitbucket API.
`

const dataSourceDocTemplate = `---
page_title: "{{.TFName}} Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket {{.Name}} via the Bitbucket Cloud API.
---

# {{.TFName}} (Data Source)

Reads Bitbucket {{.Name}} via the Bitbucket Cloud API.

## Example Usage

` + "```" + `hcl
data "{{.TFName}}" "example" {
{{- range .Params}}
  {{.}} = "{{index $.ParamValues .}}"
{{- end}}
}

output "{{snakeCase .Name}}_response" {
  value = data.{{.TFName}}.example.api_response
}
` + "```" + `

## Schema

### Required

{{- range .Params}}
- ` + "`" + `{{.}}` + "`" + ` (String) Path parameter.
{{- end}}

### Read-Only

- ` + "`" + `id` + "`" + ` (String) Resource identifier.
- ` + "`" + `api_response` + "`" + ` (String) The raw JSON response from the Bitbucket API.
`

const exampleProviderTemplate = `terraform {
  required_providers {
    bitbucket = {
      source = "FabianSchurig/bitbucket"
    }
  }
}

# Configure via environment variables:
#   BITBUCKET_USERNAME (email) + BITBUCKET_TOKEN (Atlassian API token)
#   or BITBUCKET_TOKEN alone (workspace/repository access token)
provider "bitbucket" {}
`

const exampleResourceTemplate = `resource "{{.TFName}}" "example" {
{{- range .Params}}
  {{.}} = "{{index $.ParamValues .}}"
{{- end}}
}
`

const exampleDataSourceTemplate = `data "{{.TFName}}" "example" {
{{- range .Params}}
  {{.}} = "{{index $.ParamValues .}}"
{{- end}}
}

output "{{snakeCase .Name}}_response" {
  value = data.{{.TFName}}.example.api_response
}
`

const tfTestTemplate = `# Auto-generated Terraform test for bitbucket_{{snakeCase .Name}}
# Run with: terraform test
#
# These tests use mocked provider responses. For real API tests,
# set TF_ACC=1 with BITBUCKET_USERNAME and BITBUCKET_TOKEN.

{{- if or .HasRead .HasList}}

mock_provider "bitbucket" {}

{{- if .HasRead}}

run "read_{{snakeCase .Name}}" {
  command = apply

  variables {
{{- range .Params}}
    {{.}} = "{{index $.ParamValues .}}"
{{- end}}
  }

  # Data source read should succeed with mock provider
  assert {
    condition     = data.{{.TFName}}.test.id != ""
    error_message = "Expected non-empty id for data source {{.TFName}}"
  }
}
{{- end}}

{{- if .HasCreate}}

run "create_{{snakeCase .Name}}" {
  command = apply

  variables {
{{- range .Params}}
    {{.}} = "{{index $.ParamValues .}}"
{{- end}}
  }

  # Resource create should succeed with mock provider
  assert {
    condition     = {{.TFName}}.test.id != ""
    error_message = "Expected non-empty id for resource {{.TFName}}"
  }
}
{{- end}}

{{- end}}
`

const tfTestMainTemplate = `# Auto-generated Terraform test configuration for {{.TFName}}
# This file defines the resources/data sources referenced by the test assertions.

terraform {
  required_providers {
    bitbucket = {
      source = "FabianSchurig/bitbucket"
    }
  }
}

{{- if or .HasRead .HasList}}

variable "workspace" {
  type    = string
  default = "test-workspace"
}

{{- range .Params}}
{{- if ne . "workspace"}}

variable "{{.}}" {
  type    = string
  default = "{{index $.ParamValues .}}"
}
{{- end}}
{{- end}}

provider "bitbucket" {}

data "{{.TFName}}" "test" {
{{- range .Params}}
  {{.}} = var.{{.}}
{{- end}}
}

{{- if .HasCreate}}

resource "{{.TFName}}" "test" {
{{- range .Params}}
  {{.}} = var.{{.}}
{{- end}}
}
{{- end}}

{{- end}}
`

// ─── Main ─────────────────────────────────────────────────────────────────────

func main() {
	groups := buildGroups()

	// Create output directories.
	dirs := []string{
		"docs",
		"docs/resources",
		"docs/data-sources",
		"examples/provider",
		"examples/resources",
		"examples/data-sources",
		"tests",
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			fmt.Fprintf(os.Stderr, "mkdir %s: %v\n", d, err)
			os.Exit(1)
		}
	}

	// Generate provider doc.
	writeTemplate("docs/index.md", providerDocTemplate, map[string]any{"Groups": groups})

	// Generate provider example.
	writeFile("examples/provider/provider.tf", exampleProviderTemplate)

	for _, g := range groups {
		// Resource docs.
		writeTemplate(filepath.Join("docs/resources", g.Name+".md"), resourceDocTemplate, g)

		// Data source docs.
		writeTemplate(filepath.Join("docs/data-sources", g.Name+".md"), dataSourceDocTemplate, g)

		// Resource examples.
		resDir := filepath.Join("examples/resources", g.Name)
		_ = os.MkdirAll(resDir, 0o755)
		writeTemplate(filepath.Join(resDir, "resource.tf"), exampleResourceTemplate, g)

		// Data source examples.
		dsDir := filepath.Join("examples/data-sources", g.Name)
		_ = os.MkdirAll(dsDir, 0o755)
		writeTemplate(filepath.Join(dsDir, "data-source.tf"), exampleDataSourceTemplate, g)

		// Terraform test files.
		testDir := filepath.Join("tests", g.Name)
		_ = os.MkdirAll(testDir, 0o755)
		writeTemplate(filepath.Join(testDir, "main.tf"), tfTestMainTemplate, g)
		writeTemplate(filepath.Join(testDir, g.Name+".tftest.hcl"), tfTestTemplate, g)
	}

	fmt.Printf("Generated documentation for %d resource groups\n", len(groups))
	fmt.Println("  docs/index.md")
	fmt.Printf("  docs/resources/*.md (%d files)\n", len(groups))
	fmt.Printf("  docs/data-sources/*.md (%d files)\n", len(groups))
	fmt.Printf("  examples/provider/provider.tf\n")
	fmt.Printf("  examples/resources/*/ (%d dirs)\n", len(groups))
	fmt.Printf("  examples/data-sources/*/ (%d dirs)\n", len(groups))
	fmt.Printf("  tests/*/ (%d test suites)\n", len(groups))
}

func writeTemplate(path, tmplStr string, data any) {
	tmpl, err := template.New(path).Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parsing template for %s: %v\n", path, err)
		os.Exit(1)
	}
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		fmt.Fprintf(os.Stderr, "executing template for %s: %v\n", path, err)
		os.Exit(1)
	}
	if err := os.WriteFile(path, []byte(buf.String()), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "writing %s: %v\n", path, err)
		os.Exit(1)
	}
}

func writeFile(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "writing %s: %v\n", path, err)
		os.Exit(1)
	}
}
