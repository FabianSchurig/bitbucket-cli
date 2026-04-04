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
- `active` (String) active
- `created_at` (String) created_at
- `description` (String) A user-defined description of the webhook.
- `secret` (String) The secret to associate with the hook. The secret is never returned via the API. As such, this field is only used dur...
- `secret_set` (String) Indicates whether or not the hook has an associated secret. It is not possible to see the hook's secret. This field i...
- `subject_type` (String) The type of entity. Set to either `repository` or `workspace` based on where the subscription is defined. [repository...
- `url` (String) The URL events get delivered to.
- `uuid` (String) The webhook's id
