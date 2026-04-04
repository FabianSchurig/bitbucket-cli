// gen_tfdocs reads the CRUD config and registered resource groups from
// internal/tfprovider to produce:
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

	// Import the tfprovider package to access shared CRUDConfig
	// and registered resource groups (triggers init() registration).
	"github.com/FabianSchurig/bitbucket-cli/internal/tfprovider"
)

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
	CRUDOps     []CRUDOpInfo // CRUD operation details (scopes, doc links)
}

// CRUDOpInfo holds details about a single CRUD operation for documentation.
type CRUDOpInfo struct {
	Label  string   // "Create", "Read", "Update", "Delete", "List"
	Scopes []string // OAuth2 scopes required
	DocURL string   // Atlassian REST API documentation URL
	Method string   // HTTP method (GET, POST, etc.)
	Path   string   // API path template
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
	case "known_host_uuid":
		return "{known-host-uuid}"
	case "schedule_uuid":
		return "{schedule-uuid}"
	case "member":
		return "{member-uuid}"
	case "annotation_id", "annotationId":
		return "{annotation-id}"
	case "report_id_path", "reportId":
		return "report-uuid"
	case "path":
		return "README.md"
	case "comment_id":
		return "1"
	case "runner_uuid":
		return "{runner-uuid}"
	case "cache_uuid":
		return "{cache-uuid}"
	case "fingerprint":
		return "AA:BB:CC:DD"
	case "email":
		return "user@example.com"
	case "subject_type":
		return "repository"
	default:
		return "example-value"
	}
}

func buildGroups() []GroupData {
	// Build a lookup from group name → registered ResourceGroup so we can
	// derive path params from the Read (or Create/List) operation.
	groupIndex := make(map[string]tfprovider.ResourceGroup)
	for _, g := range tfprovider.RegisteredGroups() {
		groupIndex[g.TypeName] = g
	}

	var groups []GroupData
	for name, crud := range tfprovider.CRUDConfig {
		// Derive path params from the best available operation (Read > Create > List).
		params := deriveParams(name, groupIndex)

		pv := make(map[string]string)
		hasIDParam := false
		for _, p := range params {
			pv[p] = exampleValue(p)
			if p == "param_id" {
				hasIDParam = true
			}
		}

		// Collect CRUD operation details (scopes, doc links).
		crudOps := deriveCRUDOps(name, groupIndex)

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
			CRUDOps:     crudOps,
		})
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i].Name < groups[j].Name })
	return groups
}

// deriveParams extracts the required path parameters from a resource group's
// primary operation (Read preferred, then Create, then List). This avoids
// having to maintain a separate paramConfig map.
func deriveParams(name string, index map[string]tfprovider.ResourceGroup) []string {
	rg, ok := index[name]
	if !ok {
		return nil
	}

	// Pick the best operation to derive params from.
	var op *tfprovider.OperationDef
	switch {
	case rg.Ops.Read != nil:
		op = rg.Ops.Read
	case rg.Ops.Create != nil:
		op = rg.Ops.Create
	case rg.Ops.List != nil:
		op = rg.Ops.List
	case rg.Ops.Update != nil:
		op = rg.Ops.Update
	default:
		return nil
	}

	var params []string
	for _, p := range op.Params {
		if p.In == "path" {
			params = append(params, tfprovider.ParamAttrName(p.Name))
		}
	}
	return params
}

// deriveCRUDOps collects details (scopes, doc URL) for each CRUD operation.
func deriveCRUDOps(name string, index map[string]tfprovider.ResourceGroup) []CRUDOpInfo {
	rg, ok := index[name]
	if !ok {
		return nil
	}

	type entry struct {
		label string
		op    *tfprovider.OperationDef
	}
	entries := []entry{
		{"Create", rg.Ops.Create},
		{"Read", rg.Ops.Read},
		{"Update", rg.Ops.Update},
		{"Delete", rg.Ops.Delete},
		{"List", rg.Ops.List},
	}

	var ops []CRUDOpInfo
	for _, e := range entries {
		if e.op == nil {
			continue
		}
		ops = append(ops, CRUDOpInfo{
			Label:  e.label,
			Scopes: e.op.Scopes,
			DocURL: e.op.DocURL,
			Method: e.op.Method,
			Path:   e.op.Path,
		})
	}
	return ops
}

// ─── Templates ────────────────────────────────────────────────────────────────

var funcMap = template.FuncMap{
	"replace": strings.ReplaceAll,
	"snakeCase": func(s string) string {
		return strings.ReplaceAll(s, "-", "_")
	},
	"joinScopes": func(scopes []string) string {
		quoted := make([]string, len(scopes))
		for i, s := range scopes {
			quoted[i] = "`" + s + "`"
		}
		return strings.Join(quoted, ", ")
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
{{- if .CRUDOps}}

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
{{- range .CRUDOps}}
| {{.Label}} | ` + "`" + `{{.Method}}` + "`" + ` | ` + "`" + `{{.Path}}` + "`" + ` | {{if .DocURL}}[View]({{.DocURL}}){{end}} |
{{- end}}

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
{{- range .CRUDOps}}
| {{.Label}} | {{if .Scopes}}{{joinScopes .Scopes}}{{else}}—{{end}} |
{{- end}}
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
{{- if .CRUDOps}}

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
{{- range .CRUDOps}}
{{- if or (eq .Label "Read") (eq .Label "List")}}
| {{.Label}} | ` + "`" + `{{.Method}}` + "`" + ` | ` + "`" + `{{.Path}}` + "`" + ` | {{if .DocURL}}[View]({{.DocURL}}){{end}} |
{{- end}}
{{- end}}

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
{{- range .CRUDOps}}
{{- if or (eq .Label "Read") (eq .Label "List")}}
| {{.Label}} | {{if .Scopes}}{{joinScopes .Scopes}}{{else}}—{{end}} |
{{- end}}
{{- end}}
{{- end}}

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
