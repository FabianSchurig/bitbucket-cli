---
page_title: "bitbucket_issue_comments Data Source - bitbucket"
subcategory: "Issues"
description: |-
  Reads Bitbucket issue-comments via the Bitbucket Cloud API.
---

# bitbucket_issue_comments (Data Source)

Reads Bitbucket issue-comments via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/issues/{issue_id}/comments/{comment_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-issue-tracker/#api-repositories-workspace-repo-slug-issues-issue-id-comments-comment-id-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/issues/{issue_id}/comments` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-issue-tracker/#api-repositories-workspace-repo-slug-issues-issue-id-comments-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:issue:bitbucket` |
| List | `read:issue:bitbucket` |

## Example Usage

```hcl
data "bitbucket_issue_comments" "example" {
  issue_id = "1"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "issue_comments_response" {
  value = data.bitbucket_issue_comments.example.api_response
}
```

## Schema

### Required
- `issue_id` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `comment_id` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `deleted` (String) deleted
- `updated_on` (String) updated_on
- `user` (Object) user
  Nested schema:
  - `created_on` (String) created_on
  - `display_name` (String) display_name
  - `uuid` (String) uuid

- `content` (Object) content
  Nested schema:
  - `markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]
  - `raw` (String) The text as it was typed by a user.

- `inline` (Object) inline
  Nested schema:
  - `from` (String) The comment's anchor line in the old version of the file. If the comment is a multi-line comment, this is the ending line number in the old version of the file.
  - `path` (String) The path of the file this comment is anchored to.
  - `start_from` (String) The starting line number in the old version of the file, if the comment is a multi-line comment. This is null otherwise.
  - `start_to` (String) The starting line number in the new version of the file, if the comment is a multi-line comment. This is null otherwise.
  - `to` (String) The comment's anchor line in the new version of the file. If the comment is a multi-line comment, this is the ending line number in the new version of the file.

- `issue` (Object) issue
  Nested schema:
  - `component` (Object) component
    - `id` (String) id
  - `content` (Object) content
    - `markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]
    - `raw` (String) The text as it was typed by a user.
  - `edited_on` (String) edited_on
  - `kind` (String) [bug, enhancement, proposal, task]
  - `milestone` (Object) milestone
    - `id` (String) id
  - `priority` (String) [trivial, minor, major, critical, blocker]
  - `state` (String) [submitted, new, open, resolved, on hold, invalid, duplicate, wontfix, closed]
  - `title` (String) title
  - `version` (Object) version
    - `id` (String) id
  - `votes` (String) votes
  - `id` (String) id

- `parent` (Object) parent
  Nested schema:
  - `id` (String) id

