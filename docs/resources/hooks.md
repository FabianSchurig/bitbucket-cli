---
page_title: "bitbucket_hooks Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket hooks via the Bitbucket Cloud API.
---

# bitbucket_hooks (Resource)

Manages Bitbucket hooks via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_hooks" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  uid = "webhook-uuid"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `uid` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
