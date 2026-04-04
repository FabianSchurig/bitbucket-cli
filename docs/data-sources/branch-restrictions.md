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
