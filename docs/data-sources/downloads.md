---
page_title: "bitbucket_downloads Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket downloads via the Bitbucket Cloud API.
---

# bitbucket_downloads (Data Source)

Reads Bitbucket downloads via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_downloads" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  filename = "artifact.zip"
}

output "downloads_response" {
  value = data.bitbucket_downloads.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `filename` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
