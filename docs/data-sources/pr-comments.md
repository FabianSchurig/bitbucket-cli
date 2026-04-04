---
page_title: "bitbucket_pr_comments Data Source - bitbucket"
subcategory: "Pull Requests"
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
- `pull_request_id` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `comment_id` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `deleted` (String) deleted
- `pullrequest_comment_count` (String) The number of comments for a specific pull request.
- `pullrequest_created_on` (String) The ISO8601 timestamp the request was created.
- `pullrequest_id` (String) The pull request's unique ID. Note that pull request IDs are only unique within their associated repository.
- `pullrequest_merge_commit_hash` (String) pullrequest.merge_commit.hash
- `pullrequest_participants` (List of Object) The list of users that are collaborating on this pull request.
  Nested schema:
  - `role` (String) [PARTICIPANT, REVIEWER]
  - `approved` (String) approved
  - `state` (String) [approved, changes_requested, <nil>]
  - `participated_on` (String) The ISO8601 timestamp of the participant's action. For approvers, this is the time of their approval. For commenters and pull request reviewers who are not approvers, this is the time they last commented, or null if they have not commented.

- `pullrequest_queued` (String) A boolean flag indicating whether the pull request is queued
- `pullrequest_summary_markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]
- `pullrequest_summary_raw` (String) The text as it was typed by a user.
- `pullrequest_task_count` (String) The number of open tasks for a specific pull request.
- `pullrequest_updated_on` (String) The ISO8601 timestamp the request was last updated.
- `resolution_created_on` (String) The ISO8601 timestamp the resolution was created.
- `updated_on` (String) updated_on
- `user_created_on` (String) user.created_on
- `user_display_name` (String) user.display_name
- `user_uuid` (String) user.uuid
- `content_markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]
- `content_raw` (String) The text as it was typed by a user.
- `inline_from` (String) The comment's anchor line in the old version of the file. If the comment is a multi-line comment, this is the ending line number in the old version of the file.
- `inline_path` (String) The path of the file this comment is anchored to.
- `inline_start_from` (String) The starting line number in the old version of the file, if the comment is a multi-line comment. This is null otherwise.
- `inline_start_to` (String) The starting line number in the new version of the file, if the comment is a multi-line comment. This is null otherwise.
- `inline_to` (String) The comment's anchor line in the new version of the file. If the comment is a multi-line comment, this is the ending line number in the new version of the file.
- `pending` (String) pending
- `pullrequest_close_source_branch` (String) A boolean flag indicating if merging the pull request closes the source branch.
- `pullrequest_description` (String) Explains what the pull request does.
- `pullrequest_draft` (String) A boolean flag indicating whether the pull request is a draft.
- `pullrequest_reason` (String) Explains why a pull request was declined. This field is only applicable to pull requests in rejected state.
- `pullrequest_reviewers` (List of Object) The list of users that were added as reviewers on this pull request when it was created. For performance reasons, the API only includes this list on a pull request's `self` URL.
  Nested schema:
  - `uuid` (String) uuid
  - `created_on` (String) created_on
  - `display_name` (String) display_name

- `pullrequest_state` (String) The pull request's current status. [OPEN, DRAFT, QUEUED, MERGED, DECLINED, SUPERSEDED]
- `pullrequest_title` (String) Title of the pull request.
- `resolution_type` (String) resolution.type
