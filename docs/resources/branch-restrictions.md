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
  param_id = "1"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `param_id` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only
- `api_response` (String) The raw JSON response from the Bitbucket API.
