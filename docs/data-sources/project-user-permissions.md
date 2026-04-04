---
page_title: "bitbucket_project_user_permissions Data Source - bitbucket"
subcategory: "Projects"
description: |-
  Reads Bitbucket project-user-permissions via the Bitbucket Cloud API.
---

# bitbucket_project_user_permissions (Data Source)

Reads Bitbucket project-user-permissions via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/workspaces/{workspace}/projects/{project_key}/permissions-config/users/{selected_user_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-projects/#api-workspaces-workspace-projects-project-key-permissions-config-users-selected-user-id-get) |
| List | `GET` | `/workspaces/{workspace}/projects/{project_key}/permissions-config/users` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-projects/#api-workspaces-workspace-projects-project-key-permissions-config-users-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:project:bitbucket` |
| List | `read:project:bitbucket` |

## Example Usage

```hcl
data "bitbucket_project_user_permissions" "example" {
  project_key = "PROJ"
  workspace = "my-workspace"
}

output "project_user_permissions_response" {
  value = data.bitbucket_project_user_permissions.example.api_response
}
```

## Schema

### Required
- `project_key` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `selected_user_id` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `permission` (String) [read, write, create-repo, admin, none]
- `project` (Object) project
  Nested schema:
  - `uuid` (String) The project's immutable id.
  - `key` (String) The project's key.
  - `name` (String) The name of the project.
  - `description` (String) description
  - `is_private` (String) 
  - `updated_on` (String) updated_on
  - `created_on` (String) created_on
  - `has_publicly_visible_repos` (String) 

- `type` (String) type
- `user` (Object) user
  Nested schema:
  - `created_on` (String) created_on
  - `display_name` (String) display_name
  - `uuid` (String) uuid
  - `account_id` (String) The user's Atlassian account ID.
  - `account_status` (String) The status of the account. Currently the only possible value is "active", but more values may be added in the future.
  - `has_2fa_enabled` (String) has_2fa_enabled
  - `nickname` (String) Account name defined by the owner. Should be used instead of the "username" field. Note that "nickname" cannot be used in place of "username" in URLs and queries, as "nickname" is not guaranteed to be unique.
  - `is_staff` (String) is_staff

