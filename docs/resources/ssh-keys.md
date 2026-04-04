---
page_title: "bitbucket_ssh_keys Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket ssh-keys via the Bitbucket Cloud API.
---

# bitbucket_ssh_keys (Resource)

Manages Bitbucket ssh-keys via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/users/{selected_user}/ssh-keys` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-ssh/#api-users-selected-user-ssh-keys-post) |
| Read | `GET` | `/users/{selected_user}/ssh-keys/{key_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-ssh/#api-users-selected-user-ssh-keys-key-id-get) |
| Update | `PUT` | `/users/{selected_user}/ssh-keys/{key_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-ssh/#api-users-selected-user-ssh-keys-key-id-put) |
| Delete | `DELETE` | `/users/{selected_user}/ssh-keys/{key_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-ssh/#api-users-selected-user-ssh-keys-key-id-delete) |
| List | `GET` | `/users/{selected_user}/ssh-keys` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-ssh/#api-users-selected-user-ssh-keys-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `read:ssh-key:bitbucket`, `write:ssh-key:bitbucket` |
| Read | `read:ssh-key:bitbucket` |
| Update | `read:ssh-key:bitbucket`, `write:ssh-key:bitbucket` |
| Delete | `delete:ssh-key:bitbucket` |
| List | `read:ssh-key:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_ssh_keys" "example" {
  key_id = "123"
  selected_user = "jdoe"
}
```

## Schema

### Required
- `key_id` (String) Path parameter.
- `selected_user` (String) Path parameter.

### Optional
- `comment` (String) The comment parsed from the SSH key (if present) (also computed from API response)
- `expires_on` (String) expires_on (also computed from API response)
- `fingerprint` (String) The SSH key fingerprint in SHA-256 format. (also computed from API response)
- `key` (String) The SSH public key value in OpenSSH format. (also computed from API response)
- `label` (String) The user-defined label for the SSH key (also computed from API response)
- `last_used` (String) last_used (also computed from API response)
- `uuid` (String) The SSH key's immutable ID. (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
