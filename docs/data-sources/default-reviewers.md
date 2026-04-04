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
  repo_slug = "my-repo"
  target_username = "jdoe"
  workspace = "my-workspace"
}

output "default_reviewers_response" {
  value = data.bitbucket_default_reviewers.example.api_response
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `target_username` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `display_name` (String) display_name
- `uuid` (String) uuid
