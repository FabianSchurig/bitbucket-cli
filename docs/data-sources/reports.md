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
- `created_on` (String) The timestamp when the report was created.
- `updated_on` (String) The timestamp when the report was updated.
- `data` (List of Object) An array of data fields to display information on the report. Maximum 10.
  Nested schema:
  - `type` (String) The type of data contained in the value field. If not provided, then the value will be detected as a boolean, number or string. [BOOLEAN, DATE, DURATION, LINK, NUMBER, PERCENTAGE, TEXT]
  - `title` (String) A string describing what this data field represents.

- `details` (String) A string to describe the purpose of the report.
- `external_id` (String) ID of the report provided by the report creator. It can be used to identify the report as an alternative to it's generated uuid. It is not used by Bitbucket, but only by the report creator for updating or deleting this specific report. Needs to be unique.
- `link` (String) A URL linking to the results of the report in an external tool.
- `logo_url` (String) A URL to the report logo. If none is provided, the default insights logo will be used.
- `remote_link_enabled` (String) If enabled, a remote link is created in Jira for the work item associated with the commit the report belongs to.
- `report_type` (String) The type of the report. [SECURITY, COVERAGE, TEST, BUG]
- `reporter` (String) A string to describe the tool or company who created the report.
- `result` (String) The state of the report. May be set to PENDING and later updated. [PASSED, FAILED, PENDING]
- `title` (String) The title of the report.
- `uuid` (String) The UUID that can be used to identify the report.
