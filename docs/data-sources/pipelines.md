---
page_title: "bitbucket_pipelines Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket pipelines via the Bitbucket Cloud API.
---

# bitbucket_pipelines (Data Source)

Reads Bitbucket pipelines via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_pipelines" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  pipeline_uuid = "pipeline-uuid"
}

output "pipelines_response" {
  value = data.bitbucket_pipelines.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `pipeline_uuid` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) The timestamp when the pipeline was created.
- `build_number` (String) The build number of the pipeline.
- `build_seconds_used` (String) The number of build seconds used by this pipeline.
- `completed_on` (String) The timestamp when the Pipeline was completed. This is not set if the pipeline is still in progress.
- `uuid` (String) The UUID identifying the pipeline.
