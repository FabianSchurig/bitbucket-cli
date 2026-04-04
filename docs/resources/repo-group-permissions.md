---
page_title: "bitbucket_repo_group_permissions Resource - bitbucket"
subcategory: "Repositories"
description: |-
  Manages Bitbucket repo-group-permissions via the Bitbucket Cloud API.
---

# bitbucket_repo_group_permissions (Resource)

Manages Bitbucket repo-group-permissions via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/permissions-config/groups/{group_slug}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-permissions-config-groups-group-slug-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}/permissions-config/groups/{group_slug}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-permissions-config-groups-group-slug-put) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/permissions-config/groups/{group_slug}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-permissions-config-groups-group-slug-delete) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/permissions-config/groups` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-permissions-config-groups-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:repository:bitbucket` |
| Update | `admin:repository:bitbucket`, `write:permission:bitbucket` |
| Delete | `admin:repository:bitbucket`, `delete:permission:bitbucket` |
| List | `read:repository:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_repo_group_permissions" "example" {
  group_slug = "developers"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `group_slug` (String) Path parameter.
- `repo_slug` (String) Path parameter.
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
- `permission` (String) [read, write, admin, none]
- `type` (String) type
