---
page_title: "bitbucket_branching_model Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket branching-model via the Bitbucket Cloud API.
---

# bitbucket_branching_model (Resource)

Manages Bitbucket branching-model via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **Update**: Supported

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/branching-model` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-branching-model/#api-repositories-workspace-repo-slug-branching-model-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}/branching-model/settings` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-branching-model/#api-repositories-workspace-repo-slug-branching-model-settings-put) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:repository:bitbucket` |
| Update | `admin:repository:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_branching_model" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `development_branch_default_merge_strategy` (String) The default merge strategy for pull requests targeting this branch.
- `development_branch_name` (String) The name of the ref.
- `development_branch_type` (String) development.branch.type
- `development_name` (String) Name of the target branch. Will be listed here even when the target branch does not exist. Will be `null` if targeting the main branch and the repository is empty.
- `development_use_mainbranch` (String) Indicates if the setting points at an explicit branch (`false`) or tracks the main branch (`true`).
- `production_branch_default_merge_strategy` (String) The default merge strategy for pull requests targeting this branch.
- `production_branch_name` (String) The name of the ref.
- `production_branch_type` (String) production.branch.type
- `production_name` (String) Name of the target branch. Will be listed here even when the target branch does not exist. Will be `null` if targeting the main branch and the repository is empty.
- `production_use_mainbranch` (String) Indicates if the setting points at an explicit branch (`false`) or tracks the main branch (`true`).
