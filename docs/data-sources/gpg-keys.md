---
page_title: "bitbucket_gpg_keys Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket gpg-keys via the Bitbucket Cloud API.
---

# bitbucket_gpg_keys (Data Source)

Reads Bitbucket gpg-keys via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_gpg_keys" "example" {
  fingerprint = "AA:BB:CC:DD"
  selected_user = "jdoe"
}

output "gpg_keys_response" {
  value = data.bitbucket_gpg_keys.example.api_response
}
```

## Schema

### Required
- `fingerprint` (String) Path parameter.
- `selected_user` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `added_on` (String) added_on
- `comment` (String) The comment parsed from the GPG key (if present)
- `expires_on` (String) expires_on
- `fingerprint` (String) The GPG key fingerprint.
- `key` (String) The GPG key value in X format.
- `key_id` (String) The unique identifier for the GPG key
- `last_used` (String) last_used
- `name` (String) The user-defined label for the GPG key
- `parent_fingerprint` (String) The fingerprint of the parent key. This value is null unless the current key is a subkey.
