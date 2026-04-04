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
  workspace = "my-workspace"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.

### Optional
- `uid` (String) Path parameter (auto-populated from API response).

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `active` (String) active
- `created_at` (String) created_at
- `description` (String) A user-defined description of the webhook.
- `events` (String) The events this webhook is subscribed to. (JSON array)
- `secret` (String) The secret to associate with the hook. The secret is never returned via the API. As such, this field is only used during updates. The secret can be set to `null` or "" to remove the secret (or create a hook with no secret). Leaving out the secret field during updates will leave the secret unchanged. Leaving out the secret during creation will create a hook with no secret.
- `secret_set` (String) Indicates whether or not the hook has an associated secret. It is not possible to see the hook's secret. This field is ignored during updates.
- `subject_type` (String) The type of entity. Set to either `repository` or `workspace` based on where the subscription is defined. [repository, workspace]
- `url` (String) The URL events get delivered to.
- `uuid` (String) The webhook's id
