---
page_title: "bitbucket_repo_runners Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket repo-runners via the Bitbucket Cloud API.
---

# bitbucket_repo_runners (Data Source)

Reads Bitbucket repo-runners via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_repo_runners" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  runner_uuid = "{runner-uuid}"
}

output "repo_runners_response" {
  value = data.bitbucket_repo_runners.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `runner_uuid` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
