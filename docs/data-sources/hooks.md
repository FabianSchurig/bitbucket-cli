---
page_title: "bitbucket_hooks Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket hooks via the Bitbucket Cloud API.
---

# bitbucket_hooks (Data Source)

Reads Bitbucket hooks via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/hooks/{uid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-webhooks/#api-repositories-workspace-repo-slug-hooks-uid-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/hooks` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-webhooks/#api-repositories-workspace-repo-slug-hooks-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:webhook:bitbucket` |
| List | `read:webhook:bitbucket` |

## Example Usage

```hcl
data "bitbucket_hooks" "example" {
  repo_slug = "my-repo"
  uid = "webhook-uuid"
  workspace = "my-workspace"
}

output "hooks_response" {
  value = data.bitbucket_hooks.example.api_response
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `uid` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
