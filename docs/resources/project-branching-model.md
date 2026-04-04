---
page_title: "bitbucket_project_branching_model Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket project-branching-model via the Bitbucket Cloud API.
---

# bitbucket_project_branching_model (Resource)

Manages Bitbucket project-branching-model via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **Update**: Supported

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/workspaces/{workspace}/projects/{project_key}/branching-model` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-branching-model/#api-workspaces-workspace-projects-project-key-branching-model-get) |
| Update | `PUT` | `/workspaces/{workspace}/projects/{project_key}/branching-model/settings` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-branching-model/#api-workspaces-workspace-projects-project-key-branching-model-settings-put) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:repository:bitbucket` |
| Update | `admin:project:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_project_branching_model" "example" {
  project_key = "PROJ"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `project_key` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
