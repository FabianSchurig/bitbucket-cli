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
- `added_on` (String) added_on (also computed from API response)
- `comment` (String) The comment parsed from the GPG key (if present) (also computed from API response)
- `expires_on` (String) expires_on (also computed from API response)
- `fingerprint` (String) The GPG key fingerprint. (also computed from API response)
- `key` (String) The GPG key value in X format. (also computed from API response)
- `key_id` (String) The unique identifier for the GPG key (also computed from API response)
- `last_used` (String) last_used (also computed from API response)
- `name` (String) The user-defined label for the GPG key (also computed from API response)
- `parent_fingerprint` (String) The fingerprint of the parent key. This value is null unless the current key is a subkey. (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
