---
page_title: "bitbucket_user_emails Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket user-emails via the Bitbucket Cloud API.
---

# bitbucket_user_emails (Resource)

Manages Bitbucket user-emails via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/user/emails/{email}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-users/#api-user-emails-email-get) |
| List | `GET` | `/user/emails` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-users/#api-user-emails-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:user:bitbucket` |
| List | `read:user:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_user_emails" "example" {
  email = "user@example.com"
}
```

## Schema

### Required
- `email` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
