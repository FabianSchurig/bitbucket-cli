---
page_title: "bitbucket_pipeline_schedules Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket pipeline-schedules via the Bitbucket Cloud API.
---

# bitbucket_pipeline_schedules (Resource)

Manages Bitbucket pipeline-schedules via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_pipeline_schedules" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  schedule_uuid = "{schedule-uuid}"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `schedule_uuid` (String) Path parameter.

### Optional
- `cron_pattern` (String) The cron expression with second precision (7 fields) that the schedule applies. For example, for expression: 0 0 12 *... (also computed from API response)
- `enabled` (String) Whether the schedule is enabled. (also computed from API response)
- `target_ref_name` (String) The name of the reference. (also computed from API response)
- `target_ref_type` (String) The type of reference (branch only). [branch] (also computed from API response)
- `target_selector_pattern` (String) The name of the matching pipeline definition. (also computed from API response)
- `target_selector_type` (String) The type of selector. [branches, tags, bookmarks, default, custom] (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) The timestamp when the schedule was created.
- `updated_on` (String) The timestamp when the schedule was updated.
- `uuid` (String) The UUID identifying the schedule.
