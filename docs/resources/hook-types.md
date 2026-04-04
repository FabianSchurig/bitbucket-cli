---
page_title: "bitbucket_hook_types Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket hook-types via the Bitbucket Cloud API.
---

# bitbucket_hook_types (Resource)

Manages Bitbucket hook-types via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/hook_events` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-webhooks/#api-hook-events-get) |
| List | `GET` | `/hook_events/{subject_type}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-webhooks/#api-hook-events-subject-type-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | — |
| List | — |

## Example Usage

```hcl
resource "bitbucket_hook_types" "example" {
}
```

## Schema

### Required

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
