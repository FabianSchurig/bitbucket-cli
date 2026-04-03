---
page_title: "bitbucket_commits Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket commits via the Bitbucket Cloud API.
---

# bitbucket_commits (Data Source)

Reads Bitbucket commits via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_commits" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  commit = "abc123def"
}

output "commits_response" {
  value = data.bitbucket_commits.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `commit` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
