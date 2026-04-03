---
page_title: "bitbucket_default_reviewers Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket default-reviewers via the Bitbucket Cloud API.
---

# bitbucket_default_reviewers (Resource)

Manages Bitbucket default-reviewers via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_default_reviewers" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  target_username = "jdoe"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `target_username` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
