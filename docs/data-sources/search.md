---
page_title: "bitbucket_search Data Source - bitbucket"
subcategory: "Code Search"
description: |-
  Reads Bitbucket search via the Bitbucket Cloud API.
---

# bitbucket_search (Data Source)

Reads Bitbucket search via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| List | `GET` | `/workspaces/{workspace}/search/code` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-search/#api-workspaces-workspace-search-code-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| List | `read:repository:bitbucket` |

## Example Usage

```hcl
data "bitbucket_search" "example" {
  workspace = "my-workspace"
}

output "search_response" {
  value = data.bitbucket_search.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
