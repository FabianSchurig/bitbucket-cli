---
page_title: "bitbucket_workspace_runners Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket workspace-runners via the Bitbucket Cloud API.
---

# bitbucket_workspace_runners (Data Source)

Reads Bitbucket workspace-runners via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_workspace_runners" "example" {
  workspace = "my-workspace"
  runner_uuid = "{runner-uuid}"
}

output "workspace_runners_response" {
  value = data.bitbucket_workspace_runners.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `runner_uuid` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
