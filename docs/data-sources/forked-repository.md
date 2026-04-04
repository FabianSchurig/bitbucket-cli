---
page_title: "bitbucket_forked_repository Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket forked-repository via the Bitbucket Cloud API.
---

# bitbucket_forked_repository (Data Source)

Reads Bitbucket forked-repository via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_forked_repository" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "forked_repository_response" {
  value = data.bitbucket_forked_repository.example.api_response
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
