---
page_title: "bitbucket_repos Resource - bitbucket"
subcategory: "Repositories"
description: |-
  Manages Bitbucket repos via the Bitbucket Cloud API.
---

# bitbucket_repos (Resource)

Manages Bitbucket repos via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/repositories/{workspace}/{repo_slug}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-post) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-put) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-delete) |
| List | `GET` | `/repositories/{workspace}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `admin:repository:bitbucket` |
| Read | `read:repository:bitbucket` |
| Update | `admin:repository:bitbucket` |
| Delete | `delete:repository:bitbucket` |
| List | `read:repository:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_repos" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `description` (String) description (also computed from API response)
- `fork_policy` (String)  (also computed from API response)
- `full_name` (String) The concatenation of the repository owner's username and the slugified name, e.g. "evzijst/interruptingcow". This is the same string used in Bitbucket URLs. (also computed from API response)
- `has_issues` (String)  (also computed from API response)
- `has_wiki` (String)  (also computed from API response)
- `is_private` (String) is_private (also computed from API response)
- `language` (String) language (also computed from API response)
- `mainbranch` (Object) mainbranch (also computed from API response)
  Nested schema:
  - `type` (String) type
  - `name` (String) The name of the ref.
  - `merge_strategies` (List of String) Available merge strategies for pull requests targeting this branch. [merge_commit, squash, fast_forward, squash_fast_forward, rebase_fast_forward, rebase_merge]
  - `default_merge_strategy` (String) The default merge strategy for pull requests targeting this branch.

- `name` (String) name (also computed from API response)
- `owner` (Object) owner (also computed from API response)
  Nested schema:
  - `display_name` (String) display_name
  - `uuid` (String) uuid

- `project` (Object) project (also computed from API response)
  Nested schema:
  - `key` (String) The project's key.
  - `name` (String) The name of the project.
  - `is_private` (String) 
  - `uuid` (String) The project's immutable id.
  - `description` (String) description
  - `has_publicly_visible_repos` (String) 

- `scm` (String) [git] (also computed from API response)
- `size` (String) size (also computed from API response)
- `uuid` (String) The repository's immutable id. This can be used as a substitute for the slug segment in URLs. Doing this guarantees your URLs will survive renaming of the repository by its owner, or even transfer of the repository to a different user. (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `updated_on` (String) updated_on
