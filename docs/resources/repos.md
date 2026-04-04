---
page_title: "bitbucket_repos Resource - bitbucket"
subcategory: ""
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
- `full_name` (String) The concatenation of the repository owner's username and the slugified name, e.g. "evzijst/interruptingcow". This is ... (also computed from API response)
- `has_issues` (String)  (also computed from API response)
- `has_wiki` (String)  (also computed from API response)
- `is_private` (String) is_private (also computed from API response)
- `language` (String) language (also computed from API response)
- `mainbranch_default_merge_strategy` (String) The default merge strategy for pull requests targeting this branch. (also computed from API response)
- `mainbranch_name` (String) The name of the ref. (also computed from API response)
- `mainbranch_type` (String) mainbranch.type (also computed from API response)
- `name` (String) name (also computed from API response)
- `project_description` (String) project.description (also computed from API response)
- `project_has_publicly_visible_repos` (String)  (also computed from API response)
- `project_is_private` (String)  (also computed from API response)
- `project_key` (String) The project's key. (also computed from API response)
- `project_name` (String) The name of the project. (also computed from API response)
- `project_uuid` (String) The project's immutable id. (also computed from API response)
- `scm` (String) [git] (also computed from API response)
- `size` (String) size (also computed from API response)
- `uuid` (String) The repository's immutable id. This can be used as a substitute for the slug segment in URLs. Doing this guarantees y... (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `project_created_on` (String) project.created_on
- `project_updated_on` (String) project.updated_on
- `updated_on` (String) updated_on
