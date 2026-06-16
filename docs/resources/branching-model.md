---
page_title: "bitbucket_branching_model Resource - bitbucket"
subcategory: "Branching Model"
description: |-
  Manages Bitbucket branching-model via the Bitbucket Cloud API.
---

# bitbucket_branching_model (Resource)

Manages Bitbucket branching-model via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `PUT` | `/repositories/{workspace}/{repo_slug}/branching-model/settings` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-branching-model/#api-repositories-workspace-repo-slug-branching-model-settings-put) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/branching-model/settings` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-branching-model/#api-repositories-workspace-repo-slug-branching-model-settings-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}/branching-model/settings` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-branching-model/#api-repositories-workspace-repo-slug-branching-model-settings-put) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `admin:repository:bitbucket` |
| Read | `admin:repository:bitbucket` |
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

### Optional
- `branch_types` (List of Object) branch_types (also computed from API response)
  Nested schema:
  - `enabled` (String) Whether the branch type is enabled or not. A disabled branch type may contain an invalid `prefix`.
  - `kind` (String) The kind of the branch type. [feature, bugfix, release, hotfix]
  - `prefix` (String) The prefix for this branch type. A branch with this prefix will be classified as per `kind`. The `prefix` of an enabled branch type must be a valid branch prefix.Additionally, it cannot be blank, empty or `null`. The `prefix` for a disabled branch type can be empty or invalid.

- `development` (Object) development (also computed from API response)
  Nested schema:
  - `name` (String) The configured branch. It must be `null` when `use_mainbranch` is `true`. Otherwise it must be a non-empty value. It is possible for the configured branch to not exist (e.g. it was deleted after the settings are set). In this case `is_valid` will be `false`. The branch must exist when updating/setting the `name` or an error will occur.
  - `use_mainbranch` (String) Indicates if the setting points at an explicit branch (`false`) or tracks the main branch (`true`). When `true` the `name` must be `null` or not provided. When `false` the `name` must contain a non-empty branch name.

- `production` (Object) production (also computed from API response)
  Nested schema:
  - `enabled` (String) Indicates if branch is enabled or not.
  - `name` (String) The configured branch. It must be `null` when `use_mainbranch` is `true`. Otherwise it must be a non-empty value. It is possible for the configured branch to not exist (e.g. it was deleted after the settings are set). In this case `is_valid` will be `false`. The branch must exist when updating/setting the `name` or an error will occur.
  - `use_mainbranch` (String) Indicates if the setting points at an explicit branch (`false`) or tracks the main branch (`true`). When `true` the `name` must be `null` or not provided. When `false` the `name` must contain a non-empty branch name.

- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
