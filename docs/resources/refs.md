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
  workspace = "my-workspace"
  repo_slug = "my-repo"
  name = "main"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `name` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
