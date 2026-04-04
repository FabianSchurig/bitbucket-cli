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
- `group` (Object) group
  Nested schema:
  - `name` (String) name
  - `slug` (String) The "sluggified" version of the group's name. This contains only ASCII
  - `full_slug` (String) The concatenation of the workspace's slug and the group's slug,

- `permission` (String) [read, write, admin, none]
- `repository` (Object) repository
  Nested schema:
  - `uuid` (String) The repository's immutable id. This can be used as a substitute for the slug segment in URLs. Doing this guarantees your URLs will survive renaming of the repository by its owner, or even transfer of the repository to a different user.
  - `name` (String) name
  - `description` (String) description
  - `updated_on` (String) updated_on
  - `fork_policy` (String) 
  - `is_private` (String) is_private
  - `scm` (String) [git]
  - `has_wiki` (String) 
  - `full_name` (String) The concatenation of the repository owner's username and the slugified name, e.g. "evzijst/interruptingcow". This is the same string used in Bitbucket URLs.
  - `size` (String) size
  - `language` (String) language
  - `has_issues` (String) 
  - `created_on` (String) created_on

- `type` (String) type
