---
page_title: "bitbucket_branching_model Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket branching-model via the Bitbucket Cloud API.
---

# bitbucket_branching_model (Data Source)

Reads Bitbucket branching-model via the Bitbucket Cloud API.

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
- `development_branch_default_merge_strategy` (String) The default merge strategy for pull requests targeting this branch.
- `development_branch_name` (String) The name of the ref.
- `development_branch_type` (String) development.branch.type
- `development_name` (String) Name of the target branch. Will be listed here even when the target branch does not exist. Will be `null` if targetin...
- `development_use_mainbranch` (String) Indicates if the setting points at an explicit branch (`false`) or tracks the main branch (`true`).
- `production_branch_default_merge_strategy` (String) The default merge strategy for pull requests targeting this branch.
- `production_branch_name` (String) The name of the ref.
- `production_branch_type` (String) production.branch.type
- `production_name` (String) Name of the target branch. Will be listed here even when the target branch does not exist. Will be `null` if targetin...
- `production_use_mainbranch` (String) Indicates if the setting points at an explicit branch (`false`) or tracks the main branch (`true`).
