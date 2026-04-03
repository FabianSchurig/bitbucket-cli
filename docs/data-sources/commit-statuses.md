---
page_title: "bitbucket_commit_statuses Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket commit-statuses via the Bitbucket Cloud API.
---

# bitbucket_commit_statuses (Data Source)

Reads Bitbucket commit-statuses via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_commit_statuses" "example" {
  commit = "abc123def"
  key = "build-key"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "commit_statuses_response" {
  value = data.bitbucket_commit_statuses.example.api_response
}
```

## Schema

### Required
- `commit` (String) Path parameter.
- `key` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
