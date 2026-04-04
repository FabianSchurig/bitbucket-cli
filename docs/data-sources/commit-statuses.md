---
page_title: "bitbucket_commit_statuses Data Source - bitbucket"
subcategory: ""
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
  key = "build-key"
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
- `key` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
