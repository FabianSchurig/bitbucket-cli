---
page_title: "bitbucket_project_group_permissions Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket project-group-permissions via the Bitbucket Cloud API.
---

# bitbucket_project_group_permissions (Resource)

Manages Bitbucket project-group-permissions via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/workspaces/{workspace}/projects/{project_key}/permissions-config/groups/{group_slug}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-projects/#api-workspaces-workspace-projects-project-key-permissions-config-groups-group-slug-get) |
| Update | `PUT` | `/workspaces/{workspace}/projects/{project_key}/permissions-config/groups/{group_slug}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-projects/#api-workspaces-workspace-projects-project-key-permissions-config-groups-group-slug-put) |
| Delete | `DELETE` | `/workspaces/{workspace}/projects/{project_key}/permissions-config/groups/{group_slug}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-projects/#api-workspaces-workspace-projects-project-key-permissions-config-groups-group-slug-delete) |
| List | `GET` | `/workspaces/{workspace}/projects/{project_key}/permissions-config/groups` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-projects/#api-workspaces-workspace-projects-project-key-permissions-config-groups-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:project:bitbucket` |
| Update | `admin:project:bitbucket`, `write:permission:bitbucket` |
| Delete | `admin:project:bitbucket`, `delete:permission:bitbucket` |
| List | `read:project:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_project_group_permissions" "example" {
  group_slug = "developers"
  project_key = "PROJ"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `group_slug` (String) Path parameter.
- `project_key` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
