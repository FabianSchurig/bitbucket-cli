---
page_title: "bitbucket_annotations Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket annotations via the Bitbucket Cloud API.
---

# bitbucket_annotations (Resource)

Manages Bitbucket annotations via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_annotations" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  commit = "abc123def"
  report_id = "report-uuid"
  annotation_id = "{annotation-id}"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `commit` (String) Path parameter.
- `report_id` (String) Path parameter.
- `annotation_id` (String) Path parameter.

### Optional
- `annotation_type` (String) The type of the report. [VULNERABILITY, CODE_SMELL, BUG] (also computed from API response)
- `details` (String) The details to show to users when clicking on the annotation. (also computed from API response)
- `external_id` (String) ID of the annotation provided by the annotation creator. It can be used to identify the annotation as an alternative to it's generated uuid. It is not used by Bitbucket, but only by the annotation creator for updating or deleting this specific annotation. Needs to be unique. (also computed from API response)
- `line` (String) The line number that the annotation should belong to. If no line number is provided, then it will default to 0 and in a pull request it will appear at the top of the file specified by the path field. (also computed from API response)
- `link` (String) A URL linking to the annotation in an external tool. (also computed from API response)
- `path` (String) The path of the file on which this annotation should be placed. This is the path of the file relative to the git repository. If no path is provided, then it will appear in the overview modal on all pull requests where the tip of the branch is the given commit, regardless of which files were modified. (also computed from API response)
- `result` (String) The state of the report. May be set to PENDING and later updated. [PASSED, FAILED, SKIPPED, IGNORED] (also computed from API response)
- `severity` (String) The severity of the annotation. [CRITICAL, HIGH, MEDIUM, LOW] (also computed from API response)
- `uuid` (String) The UUID that can be used to identify the annotation. (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) The timestamp when the report was created.
- `summary` (String) The message to display to users.
- `updated_on` (String) The timestamp when the report was updated.
