---
page_title: "bitbucket_workspace_hooks Data Source - bitbucket"
subcategory: "Workspaces"
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
}

output "workspace_hooks_response" {
  value = data.bitbucket_workspace_hooks.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.

### Optional
- `uid` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `active` (String) active
- `description` (String) A user-defined description of the webhook.
- `events` (List of String) The events this webhook is subscribed to. [issue:comment_created, issue:created, issue:updated, pipeline:span_created, project:updated, pullrequest:approved, pullrequest:changes_request_created, pullrequest:changes_request_removed, pullrequest:comment_created, pullrequest:comment_deleted, pullrequest:comment_reopened, pullrequest:comment_resolved, pullrequest:comment_updated, pullrequest:created, pullrequest:fulfilled, pullrequest:push, pullrequest:rejected, pullrequest:unapproved, pullrequest:updated, repo:commit_comment_created, repo:commit_status_created, repo:commit_status_updated, repo:created, repo:deleted, repo:fork, repo:imported, repo:push, repo:transfer, repo:updated]
- `secret` (String) The secret to associate with the hook. The secret is never returned via the API. As such, this field is only used during updates. The secret can be set to `null` or "" to remove the secret (or create a hook with no secret). Leaving out the secret field during updates will leave the secret unchanged. Leaving out the secret during creation will create a hook with no secret.
- `url` (String) The URL events get delivered to.
