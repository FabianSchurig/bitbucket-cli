---
page_title: "bitbucket_current_user Resource - bitbucket"
subcategory: "Users"
description: |-
  Manages Bitbucket current-user via the Bitbucket Cloud API.
---

# bitbucket_current_user (Resource)

Manages Bitbucket current-user via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/user` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-users/#api-user-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:user:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_current_user" "example" {
}
```

## Schema

### Required

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `display_name` (String) display_name
- `uuid` (String) uuid
