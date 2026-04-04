---
page_title: "bitbucket_repos Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket repos via the Bitbucket Cloud API.
---

# bitbucket_repos (Data Source)

Reads Bitbucket repos via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_repos" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "repos_response" {
  value = data.bitbucket_repos.example.api_response
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
