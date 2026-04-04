---
page_title: "bitbucket_projects Resource - bitbucket"
subcategory: "Projects"
description: |-
  Manages Bitbucket projects via the Bitbucket Cloud API.
---

# bitbucket_projects (Resource)

Manages Bitbucket projects via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/workspaces/{workspace}/projects` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-projects/#api-workspaces-workspace-projects-post) |
| Read | `GET` | `/workspaces/{workspace}/projects/{project_key}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-workspaces/#api-workspaces-workspace-projects-project-key-get) |
| Update | `PUT` | `/workspaces/{workspace}/projects/{project_key}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-projects/#api-workspaces-workspace-projects-project-key-put) |
| Delete | `DELETE` | `/workspaces/{workspace}/projects/{project_key}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-projects/#api-workspaces-workspace-projects-project-key-delete) |
| List | `GET` | `/workspaces/{workspace}/projects` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-workspaces/#api-workspaces-workspace-projects-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `admin:project:bitbucket` |
| Read | `read:project:bitbucket` |
| Update | `admin:project:bitbucket` |
| Delete | `admin:project:bitbucket` |
| List | `read:project:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_projects" "example" {
  workspace = "my-workspace"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.

### Optional
- `project_key` (String) Path parameter (auto-populated from API response).
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `description` (String) description
- `has_publicly_visible_repos` (String) 
- `is_private` (String) 
- `key` (String) The project's key.
- `name` (String) The name of the project.
- `updated_on` (String) updated_on
- `uuid` (String) The project's immutable id.
