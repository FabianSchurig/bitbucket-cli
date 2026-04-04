---
page_title: "bitbucket_gpg_keys Data Source - bitbucket"
subcategory: ""
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
