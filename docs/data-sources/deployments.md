---
page_title: "bitbucket_deployments Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket deployments via the Bitbucket Cloud API.
---

# bitbucket_deployments (Data Source)

Reads Bitbucket deployments via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_deployments" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  environment_uuid = "env-uuid"
}

output "deployments_response" {
  value = data.bitbucket_deployments.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `environment_uuid` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
