---
page_title: "bitbucket_commit_statuses Resource - bitbucket"
subcategory: ""
description: |-
  Manages Bitbucket commit-statuses via the Bitbucket Cloud API.
---

# bitbucket_commit_statuses (Resource)

Manages Bitbucket commit-statuses via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **List**: Supported (via data source)

## Example Usage

```hcl
resource "bitbucket_commit_statuses" "example" {
  commit = "abc123def"
  key = "build-key"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `commit` (String) Path parameter.
- `key` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `description` (String) A description of the build (e.g. "Unit tests in Bamboo") (also computed from API response)
- `key` (String) An identifier for the status that's unique to (also computed from API response)
- `name` (String) An identifier for the build itself, e.g. BB-DEPLOY-1 (also computed from API response)
- `refname` (String)  (also computed from API response)
- `state` (String) Provides some indication of the status of this commit [FAILED, INPROGRESS, STOPPED, SUCCESSFUL] (also computed from API response)
- `url` (String) A URL linking back to the vendor or build system, for providing more information about whatever process produced this status. Accepts context variables `repository` and `commit` that Bitbucket will evaluate at runtime whenever at runtime. For example, one could use https://foo.com/builds/{repository.full_name} which Bitbucket will turn into https://foo.com/builds/foo/bar at render time. (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `updated_on` (String) updated_on
