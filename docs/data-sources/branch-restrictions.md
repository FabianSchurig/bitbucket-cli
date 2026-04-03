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
  workspace = "my-workspace"
  repo_slug = "my-repo"
  param_id = "1"
}

output "branch_restrictions_response" {
  value = data.bitbucket_branch_restrictions.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `param_id` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
