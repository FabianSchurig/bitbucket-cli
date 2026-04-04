---
page_title: "bitbucket_pipeline_ssh_keys Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket pipeline-ssh-keys via the Bitbucket Cloud API.
---

# bitbucket_pipeline_ssh_keys (Data Source)

Reads Bitbucket pipeline-ssh-keys via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_pipeline_ssh_keys" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}

output "pipeline_ssh_keys_response" {
  value = data.bitbucket_pipeline_ssh_keys.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
