---
page_title: "bitbucket_pr Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket pr via the Bitbucket Cloud API.
---

# bitbucket_pr (Data Source)

Reads Bitbucket pr via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_pr" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  pull_request_id = "1"
}

output "pr_response" {
  value = data.bitbucket_pr.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `pull_request_id` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
