---
page_title: "bitbucket_issue_comments Resource - bitbucket"
subcategory: "Issues"
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
  issue_id = "1"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `issue_id` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `comment_id` (String) Path parameter (auto-populated from API response).
- `content` (Object) content (also computed from API response)
  Nested schema:
  - `raw` (String) The text as it was typed by a user.
  - `markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]

- `inline` (Object) inline (also computed from API response)
  Nested schema:
  - `path` (String) The path of the file this comment is anchored to.
  - `from` (String) The comment's anchor line in the old version of the file. If the comment is a multi-line comment, this is the ending line number in the old version of the file.
  - `to` (String) The comment's anchor line in the new version of the file. If the comment is a multi-line comment, this is the ending line number in the new version of the file.
  - `start_from` (String) The starting line number in the old version of the file, if the comment is a multi-line comment. This is null otherwise.
  - `start_to` (String) The starting line number in the new version of the file, if the comment is a multi-line comment. This is null otherwise.

- `issue` (Object) issue (also computed from API response)
  Nested schema:
  - `votes` (String) votes
  - `content` (Object) content
    - `markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]
    - `raw` (String) The text as it was typed by a user.
  - `title` (String) title
  - `state` (String) [submitted, new, open, resolved, on hold, invalid, duplicate, wontfix, closed]
  - `kind` (String) [bug, enhancement, proposal, task]
  - `priority` (String) [trivial, minor, major, critical, blocker]
  - `edited_on` (String) edited_on

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `deleted` (String) deleted
- `updated_on` (String) updated_on
- `user` (Object) user
  Nested schema:
  - `display_name` (String) display_name
  - `uuid` (String) uuid
  - `created_on` (String) created_on

