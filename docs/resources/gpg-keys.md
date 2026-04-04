---
page_title: "bitbucket_gpg_keys Resource - bitbucket"
subcategory: "Users"
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
  selected_user = "jdoe"
}
```

## Schema

### Required
- `selected_user` (String) Path parameter.

### Optional
- `fingerprint` (String) Path parameter (auto-populated from API response).
- `added_on` (String) added_on (also computed from API response)
- `comment` (String) The comment parsed from the GPG key (if present) (also computed from API response)
- `expires_on` (String) expires_on (also computed from API response)
- `key` (String) The GPG key value in X format. (also computed from API response)
- `key_id` (String) The unique identifier for the GPG key (also computed from API response)
- `last_used` (String) last_used (also computed from API response)
- `name` (String) The user-defined label for the GPG key (also computed from API response)
- `owner_display_name` (String) owner.display_name (also computed from API response)
- `owner_uuid` (String) owner.uuid (also computed from API response)
- `parent_fingerprint` (String) The fingerprint of the parent key. This value is null unless the current key is a subkey. (also computed from API response)
- `subkeys` (String) subkeys (JSON array) (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `owner_created_on` (String) owner.created_on
