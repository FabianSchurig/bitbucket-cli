---
page_title: "bitbucket_addon Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket addon via the Bitbucket Cloud API.
---

# bitbucket_addon (Resource)

Manages Bitbucket addon via the Bitbucket Cloud API.

## CRUD Operations
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Update | `PUT` | `/addon` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-addon/#api-addon-put) |
| Delete | `DELETE` | `/addon` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-addon/#api-addon-delete) |
| List | `GET` | `/addon/linkers` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-addon/#api-addon-linkers-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Update | — |
| Delete | — |
| List | — |

## Example Usage

```hcl
resource "bitbucket_addon" "example" {
}
```

## Schema

### Required

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
