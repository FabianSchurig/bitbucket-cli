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

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
