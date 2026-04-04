---
page_title: "bitbucket_gpg_keys Data Source - bitbucket"
subcategory: "Users"
description: |-
  Reads Bitbucket gpg-keys via the Bitbucket Cloud API.
---

# bitbucket_gpg_keys (Data Source)

Reads Bitbucket gpg-keys via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/users/{selected_user}/gpg-keys/{fingerprint}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-gpg/#api-users-selected-user-gpg-keys-fingerprint-get) |
| List | `GET` | `/users/{selected_user}/gpg-keys` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-gpg/#api-users-selected-user-gpg-keys-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:gpg-key:bitbucket` |
| List | `read:gpg-key:bitbucket` |

## Example Usage

```hcl
data "bitbucket_gpg_keys" "example" {
  selected_user = "jdoe"
}

output "gpg_keys_response" {
  value = data.bitbucket_gpg_keys.example.api_response
}
```

## Schema

### Required
- `selected_user` (String) Path parameter.

### Optional
- `fingerprint` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `added_on` (String) added_on
- `comment` (String) The comment parsed from the GPG key (if present)
- `expires_on` (String) expires_on
- `key` (String) The GPG key value in X format.
- `key_id` (String) The unique identifier for the GPG key
- `last_used` (String) last_used
- `name` (String) The user-defined label for the GPG key
- `parent_fingerprint` (String) The fingerprint of the parent key. This value is null unless the current key is a subkey.
- `subkeys` (String) subkeys (JSON array)
