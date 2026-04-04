---
page_title: "bitbucket_repo_group_permissions Data Source - bitbucket"
subcategory: "Repositories"
description: |-
  Reads Bitbucket repo-group-permissions via the Bitbucket Cloud API.
---

# bitbucket_repo_group_permissions (Data Source)

Reads Bitbucket repo-group-permissions via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/permissions-config/groups/{group_slug}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-permissions-config-groups-group-slug-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/permissions-config/groups` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-permissions-config-groups-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:repository:bitbucket` |
| List | `read:repository:bitbucket` |

## Example Usage

```hcl
data "bitbucket_repo_group_permissions" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "repo_group_permissions_response" {
  value = data.bitbucket_repo_group_permissions.example.api_response
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `group_slug` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
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
  - `scm` (String) [git]
  - `description` (String) description
  - `size` (String) size
  - `fork_policy` (String) 
  - `is_private` (String) is_private
  - `name` (String) name
  - `updated_on` (String) updated_on
  - `has_wiki` (String) 
  - `full_name` (String) The concatenation of the repository owner's username and the slugified name, e.g. "evzijst/interruptingcow". This is the same string used in Bitbucket URLs.
  - `created_on` (String) created_on
  - `language` (String) language
  - `has_issues` (String) 

- `type` (String) type
