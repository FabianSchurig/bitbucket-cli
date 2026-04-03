---
page_title: "bitbucket_branch_restrictions Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket branch-restrictions via the Bitbucket Cloud API.
---

# bitbucket_branch_restrictions (Resource)

Manages Bitbucket branch-restrictions via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_branch_restrictions" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  id = "1"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `id` (String) Path parameter.

### Optional

- `operation` (String) Override the default CRUD operation selection.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
