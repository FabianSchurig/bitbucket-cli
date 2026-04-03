---
page_title: "bitbucket_repo_deploy_keys Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket repo-deploy-keys via the Bitbucket Cloud API.
---

# bitbucket_repo_deploy_keys (Data Source)

Reads Bitbucket repo-deploy-keys via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_repo_deploy_keys" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  key_id = "123"
}

output "repo_deploy_keys_response" {
  value = data.bitbucket_repo_deploy_keys.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `key_id` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
