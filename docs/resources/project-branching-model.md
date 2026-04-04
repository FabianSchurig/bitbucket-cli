---
page_title: "bitbucket_project_branching_model Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket project-branching-model via the Bitbucket Cloud API.
---

# bitbucket_project_branching_model (Resource)

Manages Bitbucket project-branching-model via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **Update**: Supported

## Example Usage

```hcl
resource "bitbucket_project_branching_model" "example" {
  project_key = "PROJ"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `project_key` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `development_name` (String) Name of the target branch. If inherited by a repository, it will default to the main branch if the specified branch does not exist.
- `development_use_mainbranch` (String) Indicates if the setting points at an explicit branch (`false`) or tracks the main branch (`true`).
- `production_name` (String) Name of the target branch. If inherited by a repository, it will default to the main branch if the specified branch does not exist.
- `production_use_mainbranch` (String) Indicates if the setting points at an explicit branch (`false`) or tracks the main branch (`true`).
