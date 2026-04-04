---
page_title: "bitbucket_repo_user_permissions Resource - bitbucket"
subcategory: "Repositories"
description: |-
  Manages Bitbucket repo-user-permissions via the Bitbucket Cloud API.
---

# bitbucket_repo_user_permissions (Resource)

Manages Bitbucket repo-user-permissions via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/permissions-config/users/{selected_user_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-permissions-config-users-selected-user-id-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}/permissions-config/users/{selected_user_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-permissions-config-users-selected-user-id-put) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/permissions-config/users/{selected_user_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-permissions-config-users-selected-user-id-delete) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/permissions-config/users` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-permissions-config-users-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:repository:bitbucket` |
| Update | `admin:repository:bitbucket`, `write:permission:bitbucket` |
| Delete | `admin:repository:bitbucket`, `delete:permission:bitbucket` |
| List | `read:repository:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_repo_user_permissions" "example" {
  repo_slug = "my-repo"
  selected_user_id = "{user-uuid}"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `selected_user_id` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `permission` (String) [read, write, admin, none]
- `repository_created_on` (String) repository.created_on
- `repository_description` (String) repository.description
- `repository_fork_policy` (String) 
- `repository_full_name` (String) The concatenation of the repository owner's username and the slugified name, e.g. "evzijst/interruptingcow". This is the same string used in Bitbucket URLs.
- `repository_has_issues` (String) 
- `repository_has_wiki` (String) 
- `repository_is_private` (String) repository.is_private
- `repository_language` (String) repository.language
- `repository_name` (String) repository.name
- `repository_scm` (String) [git]
- `repository_size` (String) repository.size
- `repository_updated_on` (String) repository.updated_on
- `repository_uuid` (String) The repository's immutable id. This can be used as a substitute for the slug segment in URLs. Doing this guarantees your URLs will survive renaming of the repository by its owner, or even transfer of the repository to a different user.
- `type` (String) type
- `user_account_id` (String) The user's Atlassian account ID.
- `user_account_status` (String) The status of the account. Currently the only possible value is "active", but more values may be added in the future.
- `user_created_on` (String) user.created_on
- `user_display_name` (String) user.display_name
- `user_has_2fa_enabled` (String) user.has_2fa_enabled
- `user_is_staff` (String) user.is_staff
- `user_nickname` (String) Account name defined by the owner. Should be used instead of the "username" field. Note that "nickname" cannot be used in place of "username" in URLs and queries, as "nickname" is not guaranteed to be unique.
- `user_uuid` (String) user.uuid
