---
page_title: "bitbucket_pipeline_oidc_keys Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket pipeline-oidc-keys via the Bitbucket Cloud API.
---

# bitbucket_pipeline_oidc_keys (Data Source)

Reads Bitbucket pipeline-oidc-keys via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_pipeline_oidc_keys" "example" {
  workspace = "my-workspace"
}

output "pipeline_oidc_keys_response" {
  value = data.bitbucket_pipeline_oidc_keys.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
