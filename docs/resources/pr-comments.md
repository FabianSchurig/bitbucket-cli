---
page_title: "bitbucket_pr_comments Resource - bitbucket"
subcategory: "Pull Requests"
description: |-
  Manages Bitbucket pr-comments via the Bitbucket Cloud API.
---

# bitbucket_pr_comments (Resource)

Manages Bitbucket pr-comments via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/comments` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-pullrequests-pull-request-id-comments-post) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/comments/{comment_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-pullrequests-pull-request-id-comments-comment-id-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/comments/{comment_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-pullrequests-pull-request-id-comments-comment-id-put) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/comments/{comment_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-pullrequests-pull-request-id-comments-comment-id-delete) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}/comments` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-pullrequests-pull-request-id-comments-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `read:pullrequest:bitbucket` |
| Read | `read:pullrequest:bitbucket` |
| Update | `read:pullrequest:bitbucket` |
| Delete | `read:pullrequest:bitbucket` |
| List | `read:pullrequest:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_pr_comments" "example" {
  pull_request_id = "1"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `pull_request_id` (String) Path parameter.
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
  - `start_to` (String) The starting line number in the new version of the file, if the comment is a multi-line comment. This is null otherwise.
  - `path` (String) The path of the file this comment is anchored to.
  - `from` (String) The comment's anchor line in the old version of the file. If the comment is a multi-line comment, this is the ending line number in the old version of the file.
  - `to` (String) The comment's anchor line in the new version of the file. If the comment is a multi-line comment, this is the ending line number in the new version of the file.
  - `start_from` (String) The starting line number in the old version of the file, if the comment is a multi-line comment. This is null otherwise.

- `pending` (String) pending (also computed from API response)
- `pullrequest` (Object) pullrequest (also computed from API response)
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

- `resolution` (Object) The resolution object for a Comment. (also computed from API response)
  Nested schema:
  - `type` (String) type

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `deleted` (String) deleted
- `updated_on` (String) updated_on
- `user` (Object) user
  Nested schema:
  - `created_on` (String) created_on
  - `display_name` (String) display_name
  - `uuid` (String) uuid

