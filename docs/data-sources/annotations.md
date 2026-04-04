---
page_title: "bitbucket_annotations Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket annotations via the Bitbucket Cloud API.
---

# bitbucket_annotations (Data Source)

Reads Bitbucket annotations via the Bitbucket Cloud API.

## Example Usage

```hcl
data "bitbucket_annotations" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
  commit = "abc123def"
  report_id = "report-uuid"
  annotation_id = "{annotation-id}"
}

output "annotations_response" {
  value = data.bitbucket_annotations.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `commit` (String) Path parameter.
- `report_id` (String) Path parameter.
- `annotation_id` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) The timestamp when the report was created.
- `summary` (String) The message to display to users.
- `updated_on` (String) The timestamp when the report was updated.
- `annotation_type` (String) The type of the report. [VULNERABILITY, CODE_SMELL, BUG]
- `details` (String) The details to show to users when clicking on the annotation.
- `external_id` (String) ID of the annotation provided by the annotation creator. It can be used to identify the annotation as an alternative to it's generated uuid. It is not used by Bitbucket, but only by the annotation creator for updating or deleting this specific annotation. Needs to be unique.
- `line` (String) The line number that the annotation should belong to. If no line number is provided, then it will default to 0 and in a pull request it will appear at the top of the file specified by the path field.
- `link` (String) A URL linking to the annotation in an external tool.
- `path` (String) The path of the file on which this annotation should be placed. This is the path of the file relative to the git repository. If no path is provided, then it will appear in the overview modal on all pull requests where the tip of the branch is the given commit, regardless of which files were modified.
- `result` (String) The state of the report. May be set to PENDING and later updated. [PASSED, FAILED, SKIPPED, IGNORED]
- `severity` (String) The severity of the annotation. [CRITICAL, HIGH, MEDIUM, LOW]
- `uuid` (String) The UUID that can be used to identify the annotation.
