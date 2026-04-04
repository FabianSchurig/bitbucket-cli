---
page_title: "bitbucket_snippets Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket snippets via the Bitbucket Cloud API.
---

# bitbucket_snippets (Resource)

Manages Bitbucket snippets via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/snippets` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-snippets/#api-snippets-post) |
| Read | `GET` | `/snippets/{workspace}/{encoded_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-snippets/#api-snippets-workspace-encoded-id-get) |
| Update | `PUT` | `/snippets/{workspace}/{encoded_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-snippets/#api-snippets-workspace-encoded-id-put) |
| Delete | `DELETE` | `/snippets/{workspace}/{encoded_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-snippets/#api-snippets-workspace-encoded-id-delete) |
| List | `GET` | `/snippets` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-snippets/#api-snippets-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `read:snippet:bitbucket`, `write:snippet:bitbucket` |
| Read | `read:snippet:bitbucket` |
| Update | `read:snippet:bitbucket`, `write:snippet:bitbucket` |
| Delete | `delete:snippet:bitbucket` |
| List | `read:snippet:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_snippets" "example" {
  encoded_id = "snippet-id"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `encoded_id` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
