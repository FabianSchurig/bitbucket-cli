---
page_title: "bitbucket_ssh_keys Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket ssh-keys via the Bitbucket Cloud API.
---

# bitbucket_ssh_keys (Resource)

Manages Bitbucket ssh-keys via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_ssh_keys" "example" {
  key_id = "123"
  selected_user = "jdoe"
}
```

## Schema

### Required
- `key_id` (String) Path parameter.
- `selected_user` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
