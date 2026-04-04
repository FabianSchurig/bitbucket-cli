---
page_title: "bitbucket_branch_restrictions Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket branch-restrictions via the Bitbucket Cloud API.
---

# bitbucket_branch_restrictions (Data Source)

Reads Bitbucket branch-restrictions via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_branch_restrictions" "example" {
  param_id = "1"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "branch_restrictions_response" {
  value = data.bitbucket_branch_restrictions.example.api_response
}
```

## Schema

### Required
- `param_id` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `branch_match_kind` (String) Indicates how the restriction is matched against a branch. The default is `glob`. [branching_model, glob]
- `branch_type` (String) Apply the restriction to branches of this type. Active when `branch_match_kind` is `branching_model`. The branch type...
- `kind` (String) The type of restriction that is being applied. [push, delete, force, restrict_merges, require_tasks_to_be_completed, ...
- `pattern` (String) Apply the restriction to branches that match this pattern. Active when `branch_match_kind` is `glob`. Will be empty w...
- `value` (String) Value with kind-specific semantics:
