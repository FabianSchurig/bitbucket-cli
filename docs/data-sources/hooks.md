---
page_title: "bitbucket_hooks Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket hooks via the Bitbucket Cloud API.
---

# bitbucket_hooks (Data Source)

Reads Bitbucket hooks via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_hooks" "example" {
  repo_slug = "my-repo"
  uid = "webhook-uuid"
  workspace = "my-workspace"
}

output "hooks_response" {
  value = data.bitbucket_hooks.example.api_response
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `uid` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
