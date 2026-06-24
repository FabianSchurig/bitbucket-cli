---
page_title: "bitbucket_addon Data Source - bitbucket"
subcategory: "Addon"
description: |-
  Reads Bitbucket addon via the Bitbucket Cloud API.
---

# bitbucket_addon (Data Source)

Reads Bitbucket addon via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| List | `GET` | `/addon/{addon_key}/client-key` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-addon/#api-addon-addon-key-client-key-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| List | `admin:workspace:bitbucket` |

## Example Usage

```hcl
data "bitbucket_addon" "example" {
  addon_key = "example-value"
}

output "addon_response" {
  value = data.bitbucket_addon.example.api_response
}
```

## Schema

### Required
- `addon_key` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
