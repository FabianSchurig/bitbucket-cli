---
page_title: "bitbucket_pipeline_schedules Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket pipeline-schedules via the Bitbucket Cloud API.
---

# bitbucket_pipeline_schedules (Data Source)

Reads Bitbucket pipeline-schedules via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_pipeline_schedules" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  schedule_uuid = "{schedule-uuid}"
}

output "pipeline_schedules_response" {
  value = data.bitbucket_pipeline_schedules.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `schedule_uuid` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) The timestamp when the schedule was created.
- `updated_on` (String) The timestamp when the schedule was updated.
- `uuid` (String) The UUID identifying the schedule.
- `cron_pattern` (String) The cron expression with second precision (7 fields) that the schedule applies. For example, for expression: 0 0 12 *...
- `enabled` (String) Whether the schedule is enabled.
- `target_ref_name` (String) The name of the reference.
- `target_ref_type` (String) The type of reference (branch only). [branch]
- `target_selector_pattern` (String) The name of the matching pipeline definition.
- `target_selector_type` (String) The type of selector. [branches, tags, bookmarks, default, custom]
