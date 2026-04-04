---
page_title: "bitbucket_forked_repository Data Source - bitbucket"
subcategory: "Repositories"
description: |-
  Reads Bitbucket forked-repository via the Bitbucket Cloud API.
---

# bitbucket_forked_repository (Data Source)

Reads Bitbucket forked-repository via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| List | `GET` | `/repositories/{workspace}/{repo_slug}/forks` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-forks-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| List | `read:repository:bitbucket` |

## Example Usage

```hcl
data "bitbucket_forked_repository" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "forked_repository_response" {
  value = data.bitbucket_forked_repository.example.api_response
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `project_created_on` (String) project.created_on
- `project_updated_on` (String) project.updated_on
- `updated_on` (String) updated_on
- `description` (String) description
- `fork_policy` (String) 
- `full_name` (String) The concatenation of the repository owner's username and the slugified name, e.g. "evzijst/interruptingcow". This is the same string used in Bitbucket URLs.
- `has_issues` (String) 
- `has_wiki` (String) 
- `is_private` (String) is_private
- `language` (String) language
- `mainbranch_default_merge_strategy` (String) The default merge strategy for pull requests targeting this branch.
- `mainbranch_merge_strategies` (List of String) Available merge strategies for pull requests targeting this branch. [merge_commit, squash, fast_forward, squash_fast_forward, rebase_fast_forward, rebase_merge]
- `mainbranch_name` (String) The name of the ref.
- `mainbranch_type` (String) mainbranch.type
- `name` (String) name
- `project_description` (String) project.description
- `project_has_publicly_visible_repos` (String) 
- `project_is_private` (String) 
- `project_key` (String) The project's key.
- `project_name` (String) The name of the project.
- `project_uuid` (String) The project's immutable id.
- `scm` (String) [git]
- `size` (String) size
- `uuid` (String) The repository's immutable id. This can be used as a substitute for the slug segment in URLs. Doing this guarantees your URLs will survive renaming of the repository by its owner, or even transfer of the repository to a different user.
