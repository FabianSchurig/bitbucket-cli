---
page_title: "bitbucket_addon Resource - bitbucket"
subcategory: "Addon"
description: |-
  Manages Bitbucket addon via the Bitbucket Cloud API.
---

# bitbucket_addon (Resource)

Manages Bitbucket addon via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/addon/{addon_key}/client-key` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-addon/#api-addon-addon-key-client-key-get) |
| Update | `PUT` | `/addon` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-addon/#api-addon-put) |
| Delete | `DELETE` | `/addon` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-addon/#api-addon-delete) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `admin:workspace:bitbucket` |
| Update | — |
| Delete | — |

## Example Usage

```hcl
resource "bitbucket_addon" "example" {
  addon_key = "example-value"
}
```

## Schema

### Required
- `addon_key` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
