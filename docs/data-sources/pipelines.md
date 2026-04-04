---
page_title: "bitbucket_pipelines Data Source - bitbucket"
subcategory: "Pipelines"
description: |-
  Reads Bitbucket pipelines via the Bitbucket Cloud API.
---

# bitbucket_pipelines (Data Source)

Reads Bitbucket pipelines via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines/{pipeline_uuid}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-pipeline-uuid-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:pipeline:bitbucket` |
| List | `read:pipeline:bitbucket` |

## Example Usage

```hcl
data "bitbucket_pipelines" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}

output "pipelines_response" {
  value = data.bitbucket_pipelines.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.

### Optional
- `pipeline_uuid` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) The timestamp when the pipeline was created.
- `build_number` (String) The build number of the pipeline.
- `build_seconds_used` (String) The number of build seconds used by this pipeline.
- `completed_on` (String) The timestamp when the Pipeline was completed. This is not set if the pipeline is still in progress.
- `configuration_sources` (List of Object) An ordered list of sources of the pipeline configuration
  Nested schema:
  - `source` (String) Identifier of the configuration source
  - `uri` (String) Link to the configuration source view or its immediate content

- `uuid` (String) The UUID identifying the pipeline.
- `variables` (List of Object) The variables for the pipeline.
  Nested schema:
  - `uuid` (String) The UUID identifying the variable.
  - `key` (String) The unique name of the variable.
  - `value` (String) The value of the variable. If the variable is secured, this will be empty.
  - `secured` (String) If true, this variable will be treated as secured. The value will never be exposed in the logs or the REST API.

