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
  workspace = "my-workspace"
  repo_slug = "my-repo"
  commit = "abc123def"
  key = "build-key"
}

output "commit_statuses_response" {
  value = data.bitbucket_commit_statuses.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `commit` (String) Path parameter.
- `key` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
