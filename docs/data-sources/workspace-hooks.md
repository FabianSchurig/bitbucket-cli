---
page_title: "bitbucket_workspace_hooks Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket workspace-hooks via the Bitbucket Cloud API.
---

# bitbucket_workspace_hooks (Data Source)

Reads Bitbucket workspace-hooks via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/workspaces/{workspace}/hooks/{uid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-webhooks/#api-workspaces-workspace-hooks-uid-get) |
| List | `GET` | `/workspaces/{workspace}/hooks` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-webhooks/#api-workspaces-workspace-hooks-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:webhook:bitbucket` |
| List | `read:webhook:bitbucket` |

## Example Usage

```hcl
data "bitbucket_workspace_hooks" "example" {
  workspace = "my-workspace"
  uid = "webhook-uuid"
}

output "workspace_hooks_response" {
  value = data.bitbucket_workspace_hooks.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `uid` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
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
