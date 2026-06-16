---
page_title: "bitbucket_project_branching_model Data Source - bitbucket"
subcategory: "Branching Model"
description: |-
  Reads Bitbucket project-branching-model via the Bitbucket Cloud API.
---

# bitbucket_project_branching_model (Data Source)

Reads Bitbucket project-branching-model via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/workspaces/{workspace}/projects/{project_key}/branching-model/settings` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-branching-model/#api-workspaces-workspace-projects-project-key-branching-model-settings-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `admin:project:bitbucket` |

## Example Usage

```hcl
data "bitbucket_project_branching_model" "example" {
  project_key = "PROJ"
  workspace = "my-workspace"
}

output "project_branching_model_response" {
  value = data.bitbucket_project_branching_model.example.api_response
}
```

## Schema

### Required
- `project_key` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `branch_types` (List of Object) branch_types
  Nested schema:
  - `enabled` (String) Whether the branch type is enabled or not. A disabled branch type may contain an invalid `prefix`.
  - `kind` (String) The kind of the branch type. [feature, bugfix, release, hotfix]
  - `prefix` (String) The prefix for this branch type. A branch with this prefix will be classified as per `kind`. The `prefix` of an enabled branch type must be a valid branch prefix.Additionally, it cannot be blank, empty or `null`. The `prefix` for a disabled branch type can be empty or invalid.

- `development` (Object) development
  Nested schema:
  - `name` (String) The configured branch. It must be `null` when `use_mainbranch` is `true`. Otherwise it must be a non-empty value. It is possible for the configured branch to not exist (e.g. it was deleted after the settings are set). In this case `is_valid` will be `false`. The branch must exist when updating/setting the `name` or an error will occur.
  - `use_mainbranch` (String) Indicates if the setting points at an explicit branch (`false`) or tracks the main branch (`true`). When `true` the `name` must be `null` or not provided. When `false` the `name` must contain a non-empty branch name.

- `production` (Object) production
  Nested schema:
  - `enabled` (String) Indicates if branch is enabled or not.
  - `name` (String) The configured branch. It must be `null` when `use_mainbranch` is `true`. Otherwise it must be a non-empty value. It is possible for the configured branch to not exist (e.g. it was deleted after the settings are set). In this case `is_valid` will be `false`. The branch must exist when updating/setting the `name` or an error will occur.
  - `use_mainbranch` (String) Indicates if the setting points at an explicit branch (`false`) or tracks the main branch (`true`). When `true` the `name` must be `null` or not provided. When `false` the `name` must contain a non-empty branch name.

