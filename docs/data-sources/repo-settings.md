---
page_title: "bitbucket_repo_settings Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket repo-settings via the Bitbucket Cloud API.
---

# bitbucket_repo_settings (Data Source)

Reads Bitbucket repo-settings via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_repo_settings" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "repo_settings_response" {
  value = data.bitbucket_repo_settings.example.api_response
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `type` (String) type
