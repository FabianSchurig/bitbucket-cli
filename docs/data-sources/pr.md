---
page_title: "bitbucket_pr Data Source - bitbucket"
subcategory: ""
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
  pull_request_id = "1"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "pr_response" {
  value = data.bitbucket_pr.example.api_response
}
```

## Schema

### Required
- `pull_request_id` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `comment_count` (String) The number of comments for a specific pull request.
- `created_on` (String) The ISO8601 timestamp the request was created.
- `merge_commit_hash` (String) merge_commit.hash
- `queued` (String) A boolean flag indicating whether the pull request is queued
- `summary_markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]
- `summary_raw` (String) The text as it was typed by a user.
- `task_count` (String) The number of open tasks for a specific pull request.
- `updated_on` (String) The ISO8601 timestamp the request was last updated.
- `close_source_branch` (String) A boolean flag indicating if merging the pull request closes the source branch.
- `description` (String) Explains what the pull request does.
- `destination_branch_default_merge_strategy` (String) The default merge strategy, when this endpoint is the destination of the pull request.
- `destination_branch_name` (String) destination.branch.name
- `destination_commit_hash` (String) destination.commit.hash
- `draft` (String) A boolean flag indicating whether the pull request is a draft.
- `reason` (String) Explains why a pull request was declined. This field is only applicable to pull requests in rejected state.
- `source_branch_default_merge_strategy` (String) The default merge strategy, when this endpoint is the destination of the pull request.
- `source_branch_name` (String) source.branch.name
- `source_commit_hash` (String) source.commit.hash
- `state` (String) The pull request's current status. [OPEN, DRAFT, QUEUED, MERGED, DECLINED, SUPERSEDED]
- `title` (String) Title of the pull request.
