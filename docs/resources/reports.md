---
page_title: "bitbucket_reports Resource - bitbucket"
subcategory: "Reports"
description: |-
  Manages Bitbucket reports via the Bitbucket Cloud API.
---

# bitbucket_reports (Resource)

Manages Bitbucket reports via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `PUT` | `/repositories/{workspace}/{repo_slug}/commit/{commit}/reports/{reportId}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-commits/#api-repositories-workspace-repo-slug-commit-commit-reports-reportId-put) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/commit/{commit}/reports/{reportId}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-commits/#api-repositories-workspace-repo-slug-commit-commit-reports-reportId-get) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/commit/{commit}/reports/{reportId}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-commits/#api-repositories-workspace-repo-slug-commit-commit-reports-reportId-delete) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/commit/{commit}/reports` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-commits/#api-repositories-workspace-repo-slug-commit-commit-reports-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `read:repository:bitbucket` |
| Read | `read:repository:bitbucket` |
| Delete | `read:repository:bitbucket` |
| List | `read:repository:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_reports" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  commit = "abc123def"
  report_id = "report-uuid"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `commit` (String) Path parameter.
- `report_id` (String) Path parameter.

### Optional
- `data` (List of Object) An array of data fields to display information on the report. Maximum 10. (also computed from API response)
  Nested schema:
  - `title` (String) A string describing what this data field represents.
  - `type` (String) The type of data contained in the value field. If not provided, then the value will be detected as a boolean, number or string. [BOOLEAN, DATE, DURATION, LINK, NUMBER, PERCENTAGE, TEXT]

- `details` (String) A string to describe the purpose of the report. (also computed from API response)
- `external_id` (String) ID of the report provided by the report creator. It can be used to identify the report as an alternative to it's generated uuid. It is not used by Bitbucket, but only by the report creator for updating or deleting this specific report. Needs to be unique. (also computed from API response)
- `link` (String) A URL linking to the results of the report in an external tool. (also computed from API response)
- `logo_url` (String) A URL to the report logo. If none is provided, the default insights logo will be used. (also computed from API response)
- `remote_link_enabled` (String) If enabled, a remote link is created in Jira for the work item associated with the commit the report belongs to. (also computed from API response)
- `report_type` (String) The type of the report. [SECURITY, COVERAGE, TEST, BUG] (also computed from API response)
- `reporter` (String) A string to describe the tool or company who created the report. (also computed from API response)
- `result` (String) The state of the report. May be set to PENDING and later updated. [PASSED, FAILED, PENDING] (also computed from API response)
- `title` (String) The title of the report. (also computed from API response)
- `uuid` (String) The UUID that can be used to identify the report. (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) The timestamp when the report was created.
- `updated_on` (String) The timestamp when the report was updated.
