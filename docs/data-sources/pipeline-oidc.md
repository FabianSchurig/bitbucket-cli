---
page_title: "bitbucket_pipeline_oidc Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket pipeline-oidc via the Bitbucket Cloud API.
---

# bitbucket_pipeline_oidc (Data Source)

Reads Bitbucket pipeline-oidc via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_pipeline_oidc" "example" {
  workspace = "my-workspace"
}

output "pipeline_oidc_response" {
  value = data.bitbucket_pipeline_oidc.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
