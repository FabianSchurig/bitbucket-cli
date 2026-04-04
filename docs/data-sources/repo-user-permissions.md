---
page_title: "bitbucket_repo_user_permissions Data Source - bitbucket"
subcategory: "Repositories"
description: |-
  Reads Bitbucket repo-user-permissions via the Bitbucket Cloud API.
---

# bitbucket_repo_user_permissions (Data Source)

Reads Bitbucket repo-user-permissions via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/permissions-config/users/{selected_user_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-permissions-config-users-selected-user-id-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/permissions-config/users` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-permissions-config-users-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:repository:bitbucket` |
| List | `read:repository:bitbucket` |

## Example Usage

```hcl
data "bitbucket_repo_user_permissions" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "repo_user_permissions_response" {
  value = data.bitbucket_repo_user_permissions.example.api_response
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `selected_user_id` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `permission` (String) [read, write, admin, none]
- `repository` (Object) repository
  Nested schema:
  - `is_private` (String) is_private
  - `scm` (String) [git]
  - `has_wiki` (String) 
  - `full_name` (String) The concatenation of the repository owner's username and the slugified name, e.g. "evzijst/interruptingcow". This is the same string used in Bitbucket URLs.
  - `size` (String) size
  - `language` (String) language
  - `has_issues` (String) 
  - `created_on` (String) created_on
  - `uuid` (String) The repository's immutable id. This can be used as a substitute for the slug segment in URLs. Doing this guarantees your URLs will survive renaming of the repository by its owner, or even transfer of the repository to a different user.
  - `name` (String) name
  - `description` (String) description
  - `updated_on` (String) updated_on
  - `fork_policy` (String) 

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

