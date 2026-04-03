---
page_title: "bitbucket_snippets Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket snippets via the Bitbucket Cloud API.
---

# bitbucket_snippets (Resource)

Manages Bitbucket snippets via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_snippets" "example" {
  workspace = "my-workspace"
  encoded_id = "snippet-id"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `encoded_id` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
