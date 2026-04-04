---
page_title: "bitbucket_issue_comments Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket issue-comments via the Bitbucket Cloud API.
---

# bitbucket_issue_comments (Resource)

Manages Bitbucket issue-comments via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/repositories/{workspace}/{repo_slug}/issues/{issue_id}/comments` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-issue-tracker/#api-repositories-workspace-repo-slug-issues-issue-id-comments-post) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/issues/{issue_id}/comments/{comment_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-issue-tracker/#api-repositories-workspace-repo-slug-issues-issue-id-comments-comment-id-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}/issues/{issue_id}/comments/{comment_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-issue-tracker/#api-repositories-workspace-repo-slug-issues-issue-id-comments-comment-id-put) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/issues/{issue_id}/comments/{comment_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-issue-tracker/#api-repositories-workspace-repo-slug-issues-issue-id-comments-comment-id-delete) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/issues/{issue_id}/comments` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-issue-tracker/#api-repositories-workspace-repo-slug-issues-issue-id-comments-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `read:issue:bitbucket`, `write:issue:bitbucket` |
| Read | `read:issue:bitbucket` |
| Update | `read:issue:bitbucket`, `write:issue:bitbucket` |
| Delete | `write:issue:bitbucket` |
| List | `read:issue:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_issue_comments" "example" {
  comment_id = "1"
  issue_id = "1"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `comment_id` (String) Path parameter.
- `issue_id` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
