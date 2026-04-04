---
page_title: "bitbucket_project_user_permissions Resource - bitbucket"
subcategory: "Projects"
description: |-
  Manages Bitbucket project-user-permissions via the Bitbucket Cloud API.
---

# bitbucket_project_user_permissions (Resource)

Manages Bitbucket project-user-permissions via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/workspaces/{workspace}/projects/{project_key}/permissions-config/users/{selected_user_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-projects/#api-workspaces-workspace-projects-project-key-permissions-config-users-selected-user-id-get) |
| Update | `PUT` | `/workspaces/{workspace}/projects/{project_key}/permissions-config/users/{selected_user_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-projects/#api-workspaces-workspace-projects-project-key-permissions-config-users-selected-user-id-put) |
| Delete | `DELETE` | `/workspaces/{workspace}/projects/{project_key}/permissions-config/users/{selected_user_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-projects/#api-workspaces-workspace-projects-project-key-permissions-config-users-selected-user-id-delete) |
| List | `GET` | `/workspaces/{workspace}/projects/{project_key}/permissions-config/users` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-projects/#api-workspaces-workspace-projects-project-key-permissions-config-users-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:project:bitbucket` |
| Update | `admin:project:bitbucket`, `write:permission:bitbucket` |
| Delete | `admin:project:bitbucket`, `delete:permission:bitbucket` |
| List | `read:project:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_project_user_permissions" "example" {
  project_key = "PROJ"
  selected_user_id = "{user-uuid}"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `project_key` (String) Path parameter.
- `selected_user_id` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `permission` (String) [read, write, create-repo, admin, none]
- `project` (Object) project
  Nested schema:
  - `key` (String) The project's key.
  - `name` (String) The name of the project.
  - `is_private` (String) 
  - `created_on` (String) created_on
  - `uuid` (String) The project's immutable id.
  - `description` (String) description
  - `updated_on` (String) updated_on
  - `has_publicly_visible_repos` (String) 

- `type` (String) type
- `user` (Object) user
  Nested schema:
  - `created_on` (String) created_on
  - `display_name` (String) display_name
  - `uuid` (String) uuid
  - `has_2fa_enabled` (String) has_2fa_enabled
  - `nickname` (String) Account name defined by the owner. Should be used instead of the "username" field. Note that "nickname" cannot be used in place of "username" in URLs and queries, as "nickname" is not guaranteed to be unique.
  - `is_staff` (String) is_staff
  - `account_id` (String) The user's Atlassian account ID.
  - `account_status` (String) The status of the account. Currently the only possible value is "active", but more values may be added in the future.

