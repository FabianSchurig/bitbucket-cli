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
- `project_created_on` (String) project.created_on
- `project_description` (String) project.description
- `project_has_publicly_visible_repos` (String) 
- `project_is_private` (String) 
- `project_key` (String) The project's key.
- `project_name` (String) The name of the project.
- `project_updated_on` (String) project.updated_on
- `project_uuid` (String) The project's immutable id.
- `type` (String) type
- `user_account_id` (String) The user's Atlassian account ID.
- `user_account_status` (String) The status of the account. Currently the only possible value is "active", but more values may be added in the future.
- `user_created_on` (String) user.created_on
- `user_display_name` (String) user.display_name
- `user_has_2fa_enabled` (String) user.has_2fa_enabled
- `user_is_staff` (String) user.is_staff
- `user_nickname` (String) Account name defined by the owner. Should be used instead of the "username" field. Note that "nickname" cannot be used in place of "username" in URLs and queries, as "nickname" is not guaranteed to be unique.
- `user_uuid` (String) user.uuid
