---
page_title: "bitbucket_reports Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket reports via the Bitbucket Cloud API.
---

# bitbucket_reports (Data Source)

Reads Bitbucket reports via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/commit/{commit}/reports/{reportId}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-commits/#api-repositories-workspace-repo-slug-commit-commit-reports-reportId-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/commit/{commit}/reports` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-commits/#api-repositories-workspace-repo-slug-commit-commit-reports-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:repository:bitbucket` |
| List | `read:repository:bitbucket` |

## Example Usage

```hcl
data "bitbucket_reports" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  commit = "abc123def"
  report_id = "report-uuid"
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
- `report_id` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
