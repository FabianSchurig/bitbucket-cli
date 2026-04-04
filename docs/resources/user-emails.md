---
page_title: "bitbucket_user_emails Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket user-emails via the Bitbucket Cloud API.
---

# bitbucket_user_emails (Resource)

Manages Bitbucket user-emails via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_user_emails" "example" {
  email = "user@example.com"
}
```

## Schema

### Required
- `email` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
