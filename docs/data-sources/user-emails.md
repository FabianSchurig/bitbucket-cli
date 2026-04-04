---
page_title: "bitbucket_user_emails Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket user-emails via the Bitbucket Cloud API.
---

# bitbucket_user_emails (Data Source)

Reads Bitbucket user-emails via the Bitbucket Cloud API.

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
data "bitbucket_user_emails" "example" {
  email = "user@example.com"
}

output "user_emails_response" {
  value = data.bitbucket_user_emails.example.api_response
}
```

## Schema

### Required
- `email` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
