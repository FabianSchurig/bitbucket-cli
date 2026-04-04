---
page_title: "bitbucket_project_group_permissions Resource - bitbucket"
subcategory: "Projects"
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
- `group_full_slug` (String) The concatenation of the workspace's slug and the group's slug,
- `group_name` (String) group.name
- `group_slug` (String) The "sluggified" version of the group's name. This contains only ASCII
- `group_workspace_created_on` (String) group.workspace.created_on
- `group_workspace_forking_mode` (String) Controls the rules for forking repositories within this workspace.
- `group_workspace_is_privacy_enforced` (String) Indicates whether the workspace enforces private content, or whether it allows public content.
- `group_workspace_is_private` (String) Indicates whether the workspace is publicly accessible, or whether it is
- `group_workspace_name` (String) The name of the workspace.
- `group_workspace_slug` (String) The short label that identifies this workspace.
- `group_workspace_updated_on` (String) group.workspace.updated_on
- `group_workspace_uuid` (String) The workspace's immutable id.
- `permission` (String) [read, write, create-repo, admin, none]
- `project_created_on` (String) project.created_on
- `project_description` (String) project.description
- `project_has_publicly_visible_repos` (String) 
- `project_is_private` (String) 
- `project_key` (String) The project's key.
- `project_name` (String) The name of the project.
- `project_updated_on` (String) project.updated_on
- `project_uuid` (String) The project's immutable id.
- `type` (String) type
