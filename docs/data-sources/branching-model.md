---
page_title: "bitbucket_branching_model Data Source - bitbucket"
subcategory: "Branching Model"
description: |-
  Reads Bitbucket branching-model via the Bitbucket Cloud API.
---

# bitbucket_branching_model (Data Source)

Reads Bitbucket branching-model via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/branching-model` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-branching-model/#api-repositories-workspace-repo-slug-branching-model-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:repository:bitbucket` |

## Example Usage

```hcl
data "bitbucket_branching_model" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "branching_model_response" {
  value = data.bitbucket_branching_model.example.api_response
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `branch_types` (List of Object) The active branch types.
  Nested schema:
  - `kind` (String) The kind of branch. [feature, bugfix, release, hotfix]
  - `prefix` (String) The prefix for this branch type. A branch with this prefix will be classified as per `kind`. The prefix must be a valid prefix for a branch and must always exist. It cannot be blank, empty or `null`.

- `development` (Object) development
  Nested schema:
  - `branch` (Object) branch
    - `type` (String) type
    - `name` (String) The name of the ref.
    - `merge_strategies` (String) Available merge strategies for pull requests targeting this branch. [merge_commit, squash, fast_forward, squash_fast_forward, rebase_fast_forward, rebase_merge]
    - `default_merge_strategy` (String) The default merge strategy for pull requests targeting this branch.
  - `name` (String) Name of the target branch. Will be listed here even when the target branch does not exist. Will be `null` if targeting the main branch and the repository is empty.
  - `use_mainbranch` (String) Indicates if the setting points at an explicit branch (`false`) or tracks the main branch (`true`).

- `production` (Object) production
  Nested schema:
  - `branch` (Object) branch
    - `type` (String) type
    - `name` (String) The name of the ref.
    - `merge_strategies` (String) Available merge strategies for pull requests targeting this branch. [merge_commit, squash, fast_forward, squash_fast_forward, rebase_fast_forward, rebase_merge]
    - `default_merge_strategy` (String) The default merge strategy for pull requests targeting this branch.
  - `name` (String) Name of the target branch. Will be listed here even when the target branch does not exist. Will be `null` if targeting the main branch and the repository is empty.
  - `use_mainbranch` (String) Indicates if the setting points at an explicit branch (`false`) or tracks the main branch (`true`).

