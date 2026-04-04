---
page_title: "bitbucket_hooks Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket hooks via the Bitbucket Cloud API.
---

# bitbucket_hooks (Resource)

Manages Bitbucket hooks via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/repositories/{workspace}/{repo_slug}/hooks` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-webhooks/#api-repositories-workspace-repo-slug-hooks-post) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/hooks/{uid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-webhooks/#api-repositories-workspace-repo-slug-hooks-uid-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}/hooks/{uid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-webhooks/#api-repositories-workspace-repo-slug-hooks-uid-put) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/hooks/{uid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-webhooks/#api-repositories-workspace-repo-slug-hooks-uid-delete) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/hooks` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-webhooks/#api-repositories-workspace-repo-slug-hooks-get) |

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
resource "bitbucket_hooks" "example" {
  repo_slug = "my-repo"
  uid = "webhook-uuid"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `uid` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
