---
page_title: "bitbucket_commit_statuses Data Source - bitbucket"
subcategory: "Commit Statuses"
description: |-
  Reads Bitbucket commit-statuses via the Bitbucket Cloud API.
---

# bitbucket_commit_statuses (Data Source)

Reads Bitbucket commit-statuses via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/commit/{commit}/statuses/build/{key}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-commit-statuses/#api-repositories-workspace-repo-slug-commit-commit-statuses-build-key-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/commit/{commit}/statuses` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-commit-statuses/#api-repositories-workspace-repo-slug-commit-commit-statuses-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:repository:bitbucket` |
| List | `read:repository:bitbucket` |

## Example Usage

```hcl
data "bitbucket_commit_statuses" "example" {
  commit = "abc123def"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "commit_statuses_response" {
  value = data.bitbucket_commit_statuses.example.api_response
}
```

## Schema

### Required
- `commit` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `key` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `updated_on` (String) updated_on
- `description` (String) A description of the build (e.g. "Unit tests in Bamboo")
- `name` (String) An identifier for the build itself, e.g. BB-DEPLOY-1
- `refname` (String) 
- `state` (String) Provides some indication of the status of this commit [FAILED, INPROGRESS, STOPPED, SUCCESSFUL]
- `url` (String) A URL linking back to the vendor or build system, for providing more information about whatever process produced this status. Accepts context variables `repository` and `commit` that Bitbucket will evaluate at runtime whenever at runtime. For example, one could use https://foo.com/builds/{repository.full_name} which Bitbucket will turn into https://foo.com/builds/foo/bar at render time.
