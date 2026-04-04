---
page_title: "bitbucket_pr Resource - bitbucket"
subcategory: "Pull Requests"
description: |-
  Manages Bitbucket pr via the Bitbucket Cloud API.
---

# bitbucket_pr (Resource)

Manages Bitbucket pr via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/repositories/{workspace}/{repo_slug}/pullrequests` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-pullrequests-post) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-pullrequests-pull-request-id-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-pullrequests-pull-request-id-put) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/pullrequests` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-pullrequests-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `read:pullrequest:bitbucket`, `write:pullrequest:bitbucket` |
| Read | `read:pullrequest:bitbucket` |
| Update | `read:pullrequest:bitbucket`, `write:pullrequest:bitbucket` |
| List | `read:pullrequest:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_pr" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `pull_request_id` (String) Path parameter (auto-populated from API response).
- `close_source_branch` (String) A boolean flag indicating if merging the pull request closes the source branch. (also computed from API response)
- `description` (String) Explains what the pull request does. (also computed from API response)
- `destination_commit_hash` (String) destination.commit.hash (also computed from API response)
- `draft` (String) A boolean flag indicating whether the pull request is a draft. (also computed from API response)
- `reason` (String) Explains why a pull request was declined. This field is only applicable to pull requests in rejected state. (also computed from API response)
- `reviewers` (List of Object) The list of users that were added as reviewers on this pull request when it was created. For performance reasons, the API only includes this list on a pull request's `self` URL. (also computed from API response)
  Nested schema:
  - `created_on` (String) created_on
  - `display_name` (String) display_name
  - `uuid` (String) uuid

- `source_commit_hash` (String) source.commit.hash (also computed from API response)
- `state` (String) The pull request's current status. [OPEN, DRAFT, QUEUED, MERGED, DECLINED, SUPERSEDED] (also computed from API response)
- `title` (String) Title of the pull request. (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `author_created_on` (String) author.created_on
- `author_display_name` (String) author.display_name
- `author_uuid` (String) author.uuid
- `closed_by_created_on` (String) closed_by.created_on
- `closed_by_display_name` (String) closed_by.display_name
- `closed_by_uuid` (String) closed_by.uuid
- `comment_count` (String) The number of comments for a specific pull request.
- `created_on` (String) The ISO8601 timestamp the request was created.
- `merge_commit_hash` (String) merge_commit.hash
- `participants` (List of Object) The list of users that are collaborating on this pull request.
  Nested schema:
  - `role` (String) [PARTICIPANT, REVIEWER]
  - `approved` (String) approved
  - `state` (String) [approved, changes_requested, <nil>]
  - `participated_on` (String) The ISO8601 timestamp of the participant's action. For approvers, this is the time of their approval. For commenters and pull request reviewers who are not approvers, this is the time they last commented, or null if they have not commented.

- `queued` (String) A boolean flag indicating whether the pull request is queued
- `summary_markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]
- `summary_raw` (String) The text as it was typed by a user.
- `task_count` (String) The number of open tasks for a specific pull request.
- `updated_on` (String) The ISO8601 timestamp the request was last updated.
