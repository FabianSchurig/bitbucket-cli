---
page_title: "bitbucket_search Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket search via the Bitbucket Cloud API.
---

# bitbucket_search (Resource)

Manages Bitbucket search via the Bitbucket Cloud API.

## CRUD Operations
- **List**: Supported (via data source)

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
resource "bitbucket_search" "example" {
  workspace = "my-workspace"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
