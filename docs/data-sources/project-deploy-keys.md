---
page_title: "bitbucket_project_deploy_keys Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket project-deploy-keys via the Bitbucket Cloud API.
---

# bitbucket_project_deploy_keys (Data Source)

Reads Bitbucket project-deploy-keys via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_project_deploy_keys" "example" {
  key_id = "123"
  project_key = "PROJ"
  workspace = "my-workspace"
}

output "project_deploy_keys_response" {
  value = data.bitbucket_project_deploy_keys.example.api_response
}
```

## Schema

### Required
- `key_id` (String) Path parameter.
- `project_key` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
