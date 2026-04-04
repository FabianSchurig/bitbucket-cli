---
page_title: "bitbucket_annotations Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket annotations via the Bitbucket Cloud API.
---

# bitbucket_annotations (Data Source)

Reads Bitbucket annotations via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_annotations" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  commit = "abc123def"
  report_id = "report-uuid"
  annotation_id = "{annotation-id}"
}

output "annotations_response" {
  value = data.bitbucket_annotations.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `commit` (String) Path parameter.
- `report_id` (String) Path parameter.
- `annotation_id` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
