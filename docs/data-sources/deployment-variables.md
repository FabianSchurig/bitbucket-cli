---
page_title: "bitbucket_deployment_variables Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket deployment-variables via the Bitbucket Cloud API.
---

# bitbucket_deployment_variables (Data Source)

Reads Bitbucket deployment-variables via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_deployment_variables" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  environment_uuid = "env-uuid"
}

output "deployment_variables_response" {
  value = data.bitbucket_deployment_variables.example.api_response
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
