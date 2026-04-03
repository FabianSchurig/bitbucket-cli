---
page_title: "bitbucket_default_reviewers Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket default-reviewers via the Bitbucket Cloud API.
---

# bitbucket_default_reviewers (Data Source)

Reads Bitbucket default-reviewers via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_default_reviewers" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  target_username = "jdoe"
}

output "default_reviewers_response" {
  value = data.bitbucket_default_reviewers.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `target_username` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
