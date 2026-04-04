---
page_title: "bitbucket_projects Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket projects via the Bitbucket Cloud API.
---

# bitbucket_projects (Data Source)

Reads Bitbucket projects via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/workspaces/{workspace}/projects/{project_key}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-workspaces/#api-workspaces-workspace-projects-project-key-get) |
| List | `GET` | `/workspaces/{workspace}/projects` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-workspaces/#api-workspaces-workspace-projects-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:project:bitbucket` |
| List | `read:project:bitbucket` |

## Example Usage

```hcl
data "bitbucket_projects" "example" {
  project_key = "PROJ"
  workspace = "my-workspace"
}

output "projects_response" {
  value = data.bitbucket_projects.example.api_response
}
```

## Schema

### Required
- `project_key` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
