---
page_title: "bitbucket_repos Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket repos via the Bitbucket Cloud API.
---

# bitbucket_repos (Data Source)

Reads Bitbucket repos via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_repos" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}

output "repos_response" {
  value = data.bitbucket_repos.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
