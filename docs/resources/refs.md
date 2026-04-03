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

- `operation` (String) Override the default CRUD operation selection.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
