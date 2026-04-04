---
page_title: "bitbucket_pr Data Source - bitbucket"
subcategory: "Pull Requests"
description: |-
  Reads Bitbucket pr via the Bitbucket Cloud API.
---

# bitbucket_pr (Data Source)

Reads Bitbucket pr via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/pullrequests/{pull_request_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-pullrequests-pull-request-id-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/pullrequests` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-pullrequests-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:pullrequest:bitbucket` |
| List | `read:pullrequest:bitbucket` |

## Example Usage

```hcl
data "bitbucket_pr" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "pr_response" {
  value = data.bitbucket_pr.example.api_response
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `pull_request_id` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
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
- `close_source_branch` (String) A boolean flag indicating if merging the pull request closes the source branch.
- `description` (String) Explains what the pull request does.
- `destination_branch_default_merge_strategy` (String) The default merge strategy, when this endpoint is the destination of the pull request.
- `destination_branch_merge_strategies` (List of String) Available merge strategies, when this endpoint is the destination of the pull request. [merge_commit, squash, fast_forward, squash_fast_forward, rebase_fast_forward, rebase_merge]
- `destination_branch_name` (String) destination.branch.name
- `destination_commit_hash` (String) destination.commit.hash
- `draft` (String) A boolean flag indicating whether the pull request is a draft.
- `reason` (String) Explains why a pull request was declined. This field is only applicable to pull requests in rejected state.
- `reviewers` (List of Object) The list of users that were added as reviewers on this pull request when it was created. For performance reasons, the API only includes this list on a pull request's `self` URL.
  Nested schema:
  - `created_on` (String) created_on
  - `display_name` (String) display_name
  - `uuid` (String) uuid

- `source_branch_default_merge_strategy` (String) The default merge strategy, when this endpoint is the destination of the pull request.
- `source_branch_merge_strategies` (List of String) Available merge strategies, when this endpoint is the destination of the pull request. [merge_commit, squash, fast_forward, squash_fast_forward, rebase_fast_forward, rebase_merge]
- `source_branch_name` (String) source.branch.name
- `source_commit_hash` (String) source.commit.hash
- `state` (String) The pull request's current status. [OPEN, DRAFT, QUEUED, MERGED, DECLINED, SUPERSEDED]
- `title` (String) Title of the pull request.
