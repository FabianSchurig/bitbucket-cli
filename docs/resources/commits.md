---
page_title: "bitbucket_commits Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket commits via the Bitbucket Cloud API.
---

# bitbucket_commits (Resource)

Manages Bitbucket commits via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_commits" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  commit = "abc123def"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `commit` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
