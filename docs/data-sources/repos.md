---
page_title: "bitbucket_repos Data Source - bitbucket"
subcategory: "Repositories"
description: |-
  Reads Bitbucket repos via the Bitbucket Cloud API.
---

# bitbucket_repos (Data Source)

Reads Bitbucket repos via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-get) |
| List | `GET` | `/repositories/{workspace}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:repository:bitbucket` |
| List | `read:repository:bitbucket` |

## Example Usage

```hcl
data "bitbucket_repos" "example" {
  workspace = "my-workspace"
}

output "repos_response" {
  value = data.bitbucket_repos.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.

### Optional
- `repo_slug` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `updated_on` (String) updated_on
- `description` (String) description
- `fork_policy` (String) 
- `full_name` (String) The concatenation of the repository owner's username and the slugified name, e.g. "evzijst/interruptingcow". This is the same string used in Bitbucket URLs.
- `has_issues` (String) 
- `has_wiki` (String) 
- `is_private` (String) is_private
- `language` (String) language
- `mainbranch` (Object) mainbranch
  Nested schema:
  - `type` (String) type
  - `name` (String) The name of the ref.
  - `merge_strategies` (List of String) Available merge strategies for pull requests targeting this branch. [merge_commit, squash, fast_forward, squash_fast_forward, rebase_fast_forward, rebase_merge]
  - `default_merge_strategy` (String) The default merge strategy for pull requests targeting this branch.

- `name` (String) name
- `owner` (Object) owner
  Nested schema:
  - `display_name` (String) display_name
  - `uuid` (String) uuid

- `project` (Object) project
  Nested schema:
  - `key` (String) The project's key.
  - `name` (String) The name of the project.
  - `is_private` (String) 
  - `uuid` (String) The project's immutable id.
  - `description` (String) description
  - `has_publicly_visible_repos` (String) 

- `scm` (String) [git]
- `size` (String) size
- `uuid` (String) The repository's immutable id. This can be used as a substitute for the slug segment in URLs. Doing this guarantees your URLs will survive renaming of the repository by its owner, or even transfer of the repository to a different user.
