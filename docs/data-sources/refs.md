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
  workspace = "my-workspace"
  repo_slug = "my-repo"
  name = "main"
}

output "refs_response" {
  value = data.bitbucket_refs.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `name` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
