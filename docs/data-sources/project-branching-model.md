---
page_title: "bitbucket_project_branching_model Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket project-branching-model via the Bitbucket Cloud API.
---

# bitbucket_project_branching_model (Data Source)

Reads Bitbucket project-branching-model via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/workspaces/{workspace}/projects/{project_key}/branching-model` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-branching-model/#api-workspaces-workspace-projects-project-key-branching-model-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:repository:bitbucket` |

## Example Usage

```hcl
data "bitbucket_project_branching_model" "example" {
  project_key = "PROJ"
  workspace = "my-workspace"
}

output "project_branching_model_response" {
  value = data.bitbucket_project_branching_model.example.api_response
}
```

## Schema

### Required
- `project_key` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
