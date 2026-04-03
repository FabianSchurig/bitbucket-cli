---
page_title: "bitbucket_issues Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket issues via the Bitbucket Cloud API.
---

# bitbucket_issues (Data Source)

Reads Bitbucket issues via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_issues" "example" {
  issue_id = "1"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "issues_response" {
  value = data.bitbucket_issues.example.api_response
}
```

## Schema

### Required
- `issue_id` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
