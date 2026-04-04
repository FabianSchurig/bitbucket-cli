---
page_title: "bitbucket_ssh_keys Data Source - bitbucket"
subcategory: "Users"
description: |-
  Reads Bitbucket ssh-keys via the Bitbucket Cloud API.
---

# bitbucket_ssh_keys (Data Source)

Reads Bitbucket ssh-keys via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/users/{selected_user}/ssh-keys/{key_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-ssh/#api-users-selected-user-ssh-keys-key-id-get) |
| List | `GET` | `/users/{selected_user}/ssh-keys` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-ssh/#api-users-selected-user-ssh-keys-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:ssh-key:bitbucket` |
| List | `read:ssh-key:bitbucket` |

## Example Usage

```hcl
data "bitbucket_ssh_keys" "example" {
  selected_user = "jdoe"
}

output "ssh_keys_response" {
  value = data.bitbucket_ssh_keys.example.api_response
}
```

## Schema

### Required
- `selected_user` (String) Path parameter.

### Optional
- `key_id` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

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
- `owner` (Object) owner
  Nested schema:
  - `display_name` (String) display_name
  - `uuid` (String) uuid

- `uuid` (String) The SSH key's immutable ID.
