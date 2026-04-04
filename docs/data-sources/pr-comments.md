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
- `updated_on` (String) updated_on
- `user` (Object) user
  Nested schema:
  - `created_on` (String) created_on
  - `display_name` (String) display_name
  - `uuid` (String) uuid

- `content` (Object) content
  Nested schema:
  - `raw` (String) The text as it was typed by a user.
  - `markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]

- `inline` (Object) inline
  Nested schema:
  - `start_to` (String) The starting line number in the new version of the file, if the comment is a multi-line comment. This is null otherwise.
  - `path` (String) The path of the file this comment is anchored to.
  - `from` (String) The comment's anchor line in the old version of the file. If the comment is a multi-line comment, this is the ending line number in the old version of the file.
  - `to` (String) The comment's anchor line in the new version of the file. If the comment is a multi-line comment, this is the ending line number in the new version of the file.
  - `start_from` (String) The starting line number in the old version of the file, if the comment is a multi-line comment. This is null otherwise.

- `pending` (String) pending
- `pullrequest` (Object) pullrequest
  Nested schema:
  - `reviewers` (List of Object) The list of users that were added as reviewers on this pull request when it was created. For performance reasons, the API only includes this list on a pull request's `self` URL.
    - `uuid` (String) uuid
    - `created_on` (String) created_on
    - `display_name` (String) display_name
  - `state` (String) The pull request's current status. [OPEN, DRAFT, QUEUED, MERGED, DECLINED, SUPERSEDED]
  - `description` (String) Explains what the pull request does.
  - `title` (String) Title of the pull request.
  - `close_source_branch` (String) A boolean flag indicating if merging the pull request closes the source branch.
  - `reason` (String) Explains why a pull request was declined. This field is only applicable to pull requests in rejected state.
  - `draft` (String) A boolean flag indicating whether the pull request is a draft.

- `resolution` (Object) The resolution object for a Comment.
  Nested schema:
  - `type` (String) type

