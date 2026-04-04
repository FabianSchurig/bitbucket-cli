---
page_title: "bitbucket_refs Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket refs via the Bitbucket Cloud API.
---

# bitbucket_refs (Resource)

Manages Bitbucket refs via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_refs" "example" {
  name = "main"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `name` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `default_merge_strategy` (String) The default merge strategy for pull requests targeting this branch.
- `name` (String) The name of the ref.
- `type` (String) type
