---
page_title: "bitbucket_ssh_keys Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket ssh-keys via the Bitbucket Cloud API.
---

# bitbucket_ssh_keys (Data Source)

Reads Bitbucket ssh-keys via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_ssh_keys" "example" {
  key_id = "123"
  selected_user = "jdoe"
}

output "ssh_keys_response" {
  value = data.bitbucket_ssh_keys.example.api_response
}
```

## Schema

### Required
- `key_id` (String) Path parameter.
- `selected_user` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
