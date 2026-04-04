---
page_title: "bitbucket_users Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket users via the Bitbucket Cloud API.
---

# bitbucket_users (Resource)

Manages Bitbucket users via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_users" "example" {
  selected_user = "jdoe"
}
```

## Schema

### Required
- `selected_user` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `display_name` (String) display_name
- `uuid` (String) uuid
