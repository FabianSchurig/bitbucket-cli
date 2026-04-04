---
page_title: "bitbucket_gpg_keys Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket gpg-keys via the Bitbucket Cloud API.
---

# bitbucket_gpg_keys (Resource)

Manages Bitbucket gpg-keys via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/users/{selected_user}/gpg-keys` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-gpg/#api-users-selected-user-gpg-keys-post) |
| Read | `GET` | `/users/{selected_user}/gpg-keys/{fingerprint}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-gpg/#api-users-selected-user-gpg-keys-fingerprint-get) |
| Delete | `DELETE` | `/users/{selected_user}/gpg-keys/{fingerprint}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-gpg/#api-users-selected-user-gpg-keys-fingerprint-delete) |
| List | `GET` | `/users/{selected_user}/gpg-keys` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-gpg/#api-users-selected-user-gpg-keys-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `read:gpg-key:bitbucket`, `write:gpg-key:bitbucket` |
| Read | `read:gpg-key:bitbucket` |
| Delete | `delete:gpg-key:bitbucket` |
| List | `read:gpg-key:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_gpg_keys" "example" {
  fingerprint = "AA:BB:CC:DD"
  selected_user = "jdoe"
}
```

## Schema

### Required
- `fingerprint` (String) Path parameter.
- `selected_user` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
