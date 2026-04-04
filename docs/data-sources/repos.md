---
page_title: "bitbucket_repos Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket repos via the Bitbucket Cloud API.
---

# bitbucket_repos (Data Source)

Reads Bitbucket repos via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-repo-slug-get) |
| List | `GET` | `/repositories/{workspace}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-repositories/#api-repositories-workspace-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:repository:bitbucket` |
| List | `read:repository:bitbucket` |

## Example Usage

```hcl
data "bitbucket_repos" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "repos_response" {
  value = data.bitbucket_repos.example.api_response
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
