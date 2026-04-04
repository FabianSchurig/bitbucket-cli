---
page_title: "bitbucket_current_user Resource - bitbucket"
subcategory: ""
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

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
