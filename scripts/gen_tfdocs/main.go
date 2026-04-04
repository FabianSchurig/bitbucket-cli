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
	Name            string
	TFName          string // e.g., "bitbucket_repos"
	HasCreate       bool
	HasRead         bool
	HasUpdate       bool
	HasDelete       bool
	HasList         bool
	HasIDParam      bool // true if "id" is a path parameter (avoids conflict with computed id)
	HasBody         bool // true if any CRUD op accepts a body
	Params          []string
	ComputedParams  []string          // params from non-primary ops (Optional+Computed)
	ParamValues     map[string]string
	BodyFields      []FieldDoc   // writable body fields (Optional)
	ResponseFields  []FieldDoc   // computed response fields (Computed)
	OverlapFields   []FieldDoc   // fields that are both writable and computed (Optional+Computed)
	CRUDOps         []CRUDOpInfo // CRUD operation details (scopes, doc links)
}

// FieldDoc describes a Terraform attribute for documentation.
type FieldDoc struct {
	Name string // Terraform attribute name (snake_case)
	Desc string // Human-readable description
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
		// Derive path params: required (from primary op) and computed (from secondary ops).
		requiredParams, computedParams := deriveParams(name, groupIndex)

		pv := make(map[string]string)
		hasIDParam := false
		for _, p := range requiredParams {
			pv[p] = exampleValue(p)
			if p == "param_id" {
				hasIDParam = true
			}
		}
		for _, p := range computedParams {
			pv[p] = exampleValue(p)
			if p == "param_id" {
				hasIDParam = true
			}
		}

		// Derive body fields, response fields, and overlaps.
		bodyFields, responseFields, overlapFields, hasBody := deriveFields(name, groupIndex)

		// Collect CRUD operation details (scopes, doc links).
		crudOps := deriveCRUDOps(name, groupIndex)

		groups = append(groups, GroupData{
			Name:           name,
			TFName:         "bitbucket_" + strings.ReplaceAll(name, "-", "_"),
			HasCreate:      crud.Create != "",
			HasRead:        crud.Read != "",
			HasUpdate:      crud.Update != "",
			HasDelete:      crud.Delete != "",
			HasList:        crud.List != "",
			HasIDParam:     hasIDParam,
			HasBody:        hasBody,
			Params:         requiredParams,
			ComputedParams: computedParams,
			ParamValues:    pv,
			BodyFields:     bodyFields,
			ResponseFields: responseFields,
			OverlapFields:  overlapFields,
			CRUDOps:        crudOps,
		})
	}
	sort.Slice(groups, func(i, j int) bool { return groups[i].Name < groups[j].Name })
	return groups
}

// deriveParams extracts path parameters from a resource group, split into
// required params (from the primary Create/Read op) and computed params
// (from secondary ops like Read/Update/Delete that are not in the primary op).
func deriveParams(name string, index map[string]tfprovider.ResourceGroup) (required, computed []string) {
	rg, ok := index[name]
	if !ok {
		return nil, nil
	}

	// Determine the primary operation: Create if available, else Read.
	primaryOp := rg.Ops.Create
	if primaryOp == nil {
		primaryOp = rg.Ops.Read
		if primaryOp == nil {
			primaryOp = rg.Ops.List
		}
	}
	if primaryOp == nil {
		return nil, nil
	}

	// Collect required params from the primary op.
	primarySet := map[string]bool{}
	for _, p := range primaryOp.Params {
		if p.In == "path" {
			attrName := tfprovider.ParamAttrName(p.Name)
			primarySet[attrName] = true
			required = append(required, attrName)
		}
	}

	// Collect computed params from all other CRUD ops that are not in the primary op.
	computedSeen := map[string]bool{}
	crudOps := []*tfprovider.OperationDef{rg.Ops.Create, rg.Ops.Read, rg.Ops.Update, rg.Ops.Delete, rg.Ops.List}
	for _, op := range crudOps {
		if op == nil {
			continue
		}
		for _, p := range op.Params {
			if p.In != "path" {
				continue
			}
			attrName := tfprovider.ParamAttrName(p.Name)
			if !primarySet[attrName] && !computedSeen[attrName] {
				computedSeen[attrName] = true
				computed = append(computed, attrName)
			}
		}
	}
	return required, computed
}

// deriveFields extracts body fields, response fields, and their overlaps
// from a resource group's CRUD operations for documentation.
func deriveFields(name string, index map[string]tfprovider.ResourceGroup) (bodyFields, responseFields, overlapFields []FieldDoc, hasBody bool) {
	rg, ok := index[name]
	if !ok {
		return
	}

	// Collect body fields from all CRUD ops.
	bodyFieldMap := make(map[string]string) // key → description
	crudOps := []*tfprovider.OperationDef{rg.Ops.Create, rg.Ops.Read, rg.Ops.Update, rg.Ops.Delete, rg.Ops.List}
	for _, op := range crudOps {
		if op == nil {
			continue
		}
		if op.HasBody {
			hasBody = true
		}
		for _, bf := range op.BodyFields {
			key := snakeCaseField(bf.Path)
			if _, exists := bodyFieldMap[key]; !exists {
				desc := bf.Desc
				if desc == "" {
					desc = bf.Path
				}
				bodyFieldMap[key] = desc
			}
		}
	}

	// Collect response fields from Read (or Create) operation.
	responseFieldMap := make(map[string]string)
	responseOp := rg.Ops.Read
	if responseOp == nil {
		responseOp = rg.Ops.Create
	}
	if responseOp != nil {
		for _, rf := range responseOp.ResponseFields {
			key := snakeCaseField(rf.Path)
			if key == "id" || key == "api_response" || key == "request_body" {
				continue
			}
			desc := rf.Desc
			if desc == "" {
				desc = rf.Path
			}
			responseFieldMap[key] = desc
		}
	}

	// Categorize into body-only, response-only, and overlap.
	overlapSet := make(map[string]bool)
	for key, desc := range bodyFieldMap {
		if _, isResp := responseFieldMap[key]; isResp {
			overlapFields = append(overlapFields, FieldDoc{Name: key, Desc: truncateDesc(desc)})
			overlapSet[key] = true
		} else {
			bodyFields = append(bodyFields, FieldDoc{Name: key, Desc: truncateDesc(desc)})
		}
	}
	for key, desc := range responseFieldMap {
		if !overlapSet[key] {
			responseFields = append(responseFields, FieldDoc{Name: key, Desc: truncateDesc(desc)})
		}
	}

	// Sort all lists for deterministic output.
	sort.Slice(bodyFields, func(i, j int) bool { return bodyFields[i].Name < bodyFields[j].Name })
	sort.Slice(responseFields, func(i, j int) bool { return responseFields[i].Name < responseFields[j].Name })
	sort.Slice(overlapFields, func(i, j int) bool { return overlapFields[i].Name < overlapFields[j].Name })
	return
}

// snakeCaseField converts a dot-separated field path to snake_case attribute name.
func snakeCaseField(path string) string {
	s := strings.ReplaceAll(path, ".", "_")
	s = strings.ReplaceAll(s, "-", "_")
	return strings.ToLower(s)
}

// truncateDesc returns a single-line description for documentation.
func truncateDesc(desc string) string {
	// Take first line only.
	if idx := strings.IndexByte(desc, '\n'); idx >= 0 {
		desc = desc[:idx]
	}
	desc = strings.TrimSpace(desc)
	return desc
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
- **Body fields** become optional string attributes (writable)
- **Response fields** become computed string attributes (read-only, auto-populated from API response)
- Fields present in both request and response are **Optional+Computed** (can be set by user, also populated from API)
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
{{- if or .OverlapFields .BodyFields .HasBody .ComputedParams}}

### Optional
{{- end}}
{{- range .ComputedParams}}
- ` + "`" + `{{.}}` + "`" + ` (String) Path parameter (auto-populated from API response).
{{- end}}
{{- range .OverlapFields}}
- ` + "`" + `{{.Name}}` + "`" + ` (String) {{.Desc}} (also computed from API response)
{{- end}}
{{- range .BodyFields}}
- ` + "`" + `{{.Name}}` + "`" + ` (String) {{.Desc}}
{{- end}}
{{- if .HasBody}}
- ` + "`" + `request_body` + "`" + ` (String) Raw JSON request body for create/update operations. Use ` + "`" + `jsonencode({...})` + "`" + ` to pass fields not exposed as individual attributes.
{{- end}}

### Read-Only
{{- if not .HasIDParam}}

- ` + "`" + `id` + "`" + ` (String) Resource identifier (extracted from API response).
{{- end}}
- ` + "`" + `api_response` + "`" + ` (String) The raw JSON response from the Bitbucket API.
{{- range .ResponseFields}}
- ` + "`" + `{{.Name}}` + "`" + ` (String) {{.Desc}}
{{- end}}
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
{{- range .ComputedParams}}
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
{{- range .ComputedParams}}
- ` + "`" + `{{.}}` + "`" + ` (String) Path parameter.
{{- end}}

### Read-Only

- ` + "`" + `id` + "`" + ` (String) Resource identifier.
- ` + "`" + `api_response` + "`" + ` (String) The raw JSON response from the Bitbucket API.
{{- range .ResponseFields}}
- ` + "`" + `{{.Name}}` + "`" + ` (String) {{.Desc}}
{{- end}}
{{- range .OverlapFields}}
- ` + "`" + `{{.Name}}` + "`" + ` (String) {{.Desc}}
{{- end}}
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
{{- range .ComputedParams}}
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
{{- range .ComputedParams}}
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
{{- range .ComputedParams}}
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
{{- range .ComputedParams}}
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
