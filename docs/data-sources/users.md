---
page_title: "bitbucket_users Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket users via the Bitbucket Cloud API.
---

# bitbucket_users (Data Source)

Reads Bitbucket users via the Bitbucket Cloud API.

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
data "bitbucket_users" "example" {
  selected_user = "jdoe"
}

output "users_response" {
  value = data.bitbucket_users.example.api_response
}
```

## Schema

### Required
- `selected_user` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
