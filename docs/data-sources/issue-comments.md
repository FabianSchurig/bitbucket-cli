---
page_title: "bitbucket_issue_comments Data Source - bitbucket"
subcategory: ""
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
  comment_id = "1"
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
- `comment_id` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `deleted` (String) deleted
- `issue_component_id` (String) issue.component.id
- `issue_created_on` (String) issue.created_on
- `issue_id` (String) issue.id
- `issue_milestone_id` (String) issue.milestone.id
- `issue_updated_on` (String) issue.updated_on
- `issue_version_id` (String) issue.version.id
- `updated_on` (String) updated_on
- `content_markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]
- `content_raw` (String) The text as it was typed by a user.
- `inline_from` (String) The comment's anchor line in the old version of the file. If the comment is a multi-line comment, this is the ending line number in the old version of the file.
- `inline_path` (String) The path of the file this comment is anchored to.
- `inline_start_from` (String) The starting line number in the old version of the file, if the comment is a multi-line comment. This is null otherwise.
- `inline_start_to` (String) The starting line number in the new version of the file, if the comment is a multi-line comment. This is null otherwise.
- `inline_to` (String) The comment's anchor line in the new version of the file. If the comment is a multi-line comment, this is the ending line number in the new version of the file.
- `issue_component_name` (String) issue.component.name
- `issue_content_markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]
- `issue_content_raw` (String) The text as it was typed by a user.
- `issue_edited_on` (String) issue.edited_on
- `issue_kind` (String) [bug, enhancement, proposal, task]
- `issue_milestone_name` (String) issue.milestone.name
- `issue_priority` (String) [trivial, minor, major, critical, blocker]
- `issue_state` (String) [submitted, new, open, resolved, on hold, invalid, duplicate, wontfix, closed]
- `issue_title` (String) issue.title
- `issue_version_name` (String) issue.version.name
- `issue_votes` (String) issue.votes
- `parent_id` (String) ID of referenced parent
