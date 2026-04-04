---
page_title: "bitbucket_workspace_hooks Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket workspace-hooks via the Bitbucket Cloud API.
---

# bitbucket_workspace_hooks (Resource)

Manages Bitbucket workspace-hooks via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/workspaces/{workspace}/hooks` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-webhooks/#api-workspaces-workspace-hooks-post) |
| Read | `GET` | `/workspaces/{workspace}/hooks/{uid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-webhooks/#api-workspaces-workspace-hooks-uid-get) |
| Update | `PUT` | `/workspaces/{workspace}/hooks/{uid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-webhooks/#api-workspaces-workspace-hooks-uid-put) |
| Delete | `DELETE` | `/workspaces/{workspace}/hooks/{uid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-webhooks/#api-workspaces-workspace-hooks-uid-delete) |
| List | `GET` | `/workspaces/{workspace}/hooks` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-webhooks/#api-workspaces-workspace-hooks-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `read:webhook:bitbucket`, `write:webhook:bitbucket` |
| Read | `read:webhook:bitbucket` |
| Update | `read:webhook:bitbucket`, `write:webhook:bitbucket` |
| Delete | `delete:webhook:bitbucket` |
| List | `read:webhook:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_workspace_hooks" "example" {
  uid = "webhook-uuid"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `uid` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
