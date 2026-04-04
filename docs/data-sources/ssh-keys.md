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
- `created_on` (String) created_on
- `comment` (String) The comment parsed from the SSH key (if present)
- `expires_on` (String) expires_on
- `fingerprint` (String) The SSH key fingerprint in SHA-256 format.
- `key` (String) The SSH public key value in OpenSSH format.
- `label` (String) The user-defined label for the SSH key
- `last_used` (String) last_used
- `uuid` (String) The SSH key's immutable ID.
