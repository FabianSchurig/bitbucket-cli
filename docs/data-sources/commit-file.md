---
page_title: "bitbucket_commit_file Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket commit-file via the Bitbucket Cloud API.
---

# bitbucket_commit_file (Data Source)

Reads Bitbucket commit-file via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_commit_file" "example" {
  commit = "abc123def"
  path = "README.md"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "commit_file_response" {
  value = data.bitbucket_commit_file.example.api_response
}
```

## Schema

### Required
- `commit` (String) Path parameter.
- `path` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
