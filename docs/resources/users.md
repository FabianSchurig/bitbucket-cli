---
page_title: "bitbucket_users Resource - bitbucket"
subcategory: "Users"
description: |-
  Manages Bitbucket users via the Bitbucket Cloud API.
---

# bitbucket_users (Resource)

Manages Bitbucket users via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/users/{selected_user}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-users/#api-users-selected-user-get) |
| List | `GET` | `/users/{selected_user}/ssh-keys` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-ssh/#api-users-selected-user-ssh-keys-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:user:bitbucket` |
| List | `read:ssh-key:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_users" "example" {
  selected_user = "jdoe"
}
```

## Schema

### Required
- `selected_user` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `display_name` (String) display_name
- `uuid` (String) uuid
