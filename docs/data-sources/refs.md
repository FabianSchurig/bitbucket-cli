---
page_title: "bitbucket_refs Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket refs via the Bitbucket Cloud API.
---

# bitbucket_refs (Data Source)

Reads Bitbucket refs via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_refs" "example" {
  name = "main"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "refs_response" {
  value = data.bitbucket_refs.example.api_response
}
```

## Schema

### Required
- `name` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `default_merge_strategy` (String) The default merge strategy for pull requests targeting this branch.
- `name` (String) The name of the ref.
- `type` (String) type
