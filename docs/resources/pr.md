---
page_title: "bitbucket_pr Resource - bitbucket"
subcategory: ""
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

## Example Usage

```hcl
resource "bitbucket_pr" "example" {
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
- `close_source_branch` (String) A boolean flag indicating if merging the pull request closes the source branch. (also computed from API response)
- `description` (String) Explains what the pull request does. (also computed from API response)
- `destination_branch_default_merge_strategy` (String) The default merge strategy, when this endpoint is the destination of the pull request. (also computed from API response)
- `destination_branch_name` (String) destination.branch.name (also computed from API response)
- `destination_commit_hash` (String) destination.commit.hash (also computed from API response)
- `draft` (String) A boolean flag indicating whether the pull request is a draft. (also computed from API response)
- `reason` (String) Explains why a pull request was declined. This field is only applicable to pull requests in rejected state. (also computed from API response)
- `source_branch_default_merge_strategy` (String) The default merge strategy, when this endpoint is the destination of the pull request. (also computed from API response)
- `source_branch_name` (String) source.branch.name (also computed from API response)
- `source_commit_hash` (String) source.commit.hash (also computed from API response)
- `state` (String) The pull request's current status. [OPEN, DRAFT, QUEUED, MERGED, DECLINED, SUPERSEDED] (also computed from API response)
- `title` (String) Title of the pull request. (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `comment_count` (String) The number of comments for a specific pull request.
- `created_on` (String) The ISO8601 timestamp the request was created.
- `merge_commit_hash` (String) merge_commit.hash
- `queued` (String) A boolean flag indicating whether the pull request is queued
- `summary_markup` (String) The type of markup language the raw content is to be interpreted in. [markdown, creole, plaintext]
- `summary_raw` (String) The text as it was typed by a user.
- `task_count` (String) The number of open tasks for a specific pull request.
- `updated_on` (String) The ISO8601 timestamp the request was last updated.
