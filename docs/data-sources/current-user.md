---
page_title: "bitbucket_current_user Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket current-user via the Bitbucket Cloud API.
---

# bitbucket_current_user (Data Source)

Reads Bitbucket current-user via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_current_user" "example" {
}

output "current_user_response" {
  value = data.bitbucket_current_user.example.api_response
}
```

## Schema

### Required

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `display_name` (String) display_name
- `uuid` (String) uuid
