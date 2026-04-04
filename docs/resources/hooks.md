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

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `active` (String) active
- `created_at` (String) created_at
- `description` (String) A user-defined description of the webhook.
- `secret` (String) The secret to associate with the hook. The secret is never returned via the API. As such, this field is only used during updates. The secret can be set to `null` or "" to remove the secret (or create a hook with no secret). Leaving out the secret field during updates will leave the secret unchanged. Leaving out the secret during creation will create a hook with no secret.
- `secret_set` (String) Indicates whether or not the hook has an associated secret. It is not possible to see the hook's secret. This field is ignored during updates.
- `subject_type` (String) The type of entity. Set to either `repository` or `workspace` based on where the subscription is defined. [repository, workspace]
- `url` (String) The URL events get delivered to.
- `uuid` (String) The webhook's id
