---
page_title: "bitbucket_reports Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket reports via the Bitbucket Cloud API.
---

# bitbucket_reports (Data Source)

Reads Bitbucket reports via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_reports" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  commit = "abc123def"
  reportId = "report-uuid"
}

output "reports_response" {
  value = data.bitbucket_reports.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `commit` (String) Path parameter.
- `reportId` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
