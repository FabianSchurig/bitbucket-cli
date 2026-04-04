---
page_title: "bitbucket_branch_restrictions Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket branch-restrictions via the Bitbucket Cloud API.
---

# bitbucket_branch_restrictions (Resource)

Manages Bitbucket branch-restrictions via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_branch_restrictions" "example" {
  param_id = "1"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `param_id` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `branch_match_kind` (String) Indicates how the restriction is matched against a branch. The default is `glob`. [branching_model, glob] (also computed from API response)
- `branch_type` (String) Apply the restriction to branches of this type. Active when `branch_match_kind` is `branching_model`. The branch type... (also computed from API response)
- `kind` (String) The type of restriction that is being applied. [push, delete, force, restrict_merges, require_tasks_to_be_completed, ... (also computed from API response)
- `pattern` (String) Apply the restriction to branches that match this pattern. Active when `branch_match_kind` is `glob`. Will be empty w... (also computed from API response)
- `value` (String) Value with kind-specific semantics: (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only
- `api_response` (String) The raw JSON response from the Bitbucket API.
