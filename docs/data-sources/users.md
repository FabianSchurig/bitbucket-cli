---
page_title: "bitbucket_users Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket users via the Bitbucket Cloud API.
---

# bitbucket_users (Data Source)

Reads Bitbucket users via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_users" "example" {
  selected_user = "jdoe"
}

output "users_response" {
  value = data.bitbucket_users.example.api_response
}
```

## Schema

### Required
- `selected_user` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `display_name` (String) display_name
- `uuid` (String) uuid
