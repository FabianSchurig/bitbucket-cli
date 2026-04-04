---
page_title: "bitbucket_reports Resource - bitbucket"
subcategory: ""
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
- `details` (String) A string to describe the purpose of the report. (also computed from API response)
- `external_id` (String) ID of the report provided by the report creator. It can be used to identify the report as an alternative to it's gene... (also computed from API response)
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
