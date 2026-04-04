---
page_title: "bitbucket_pr_comments Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket pr-comments via the Bitbucket Cloud API.
---

# bitbucket_pr_comments (Data Source)

Reads Bitbucket pr-comments via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/comments/{comment_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-pullrequests-pull-request-id-comments-comment-id-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/comments` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-pullrequests-pull-request-id-comments-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:pullrequest:bitbucket` |
| List | `read:pullrequest:bitbucket` |

## Example Usage

```hcl
data "bitbucket_pr_comments" "example" {
  comment_id = "1"
  pull_request_id = "1"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "pr_comments_response" {
  value = data.bitbucket_pr_comments.example.api_response
}
```

## Schema

### Required
- `comment_id` (String) Path parameter.
- `pull_request_id` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
