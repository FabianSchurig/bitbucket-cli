---
page_title: "bitbucket_refs Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket refs via the Bitbucket Cloud API.
---

# bitbucket_refs (Data Source)

Reads Bitbucket refs via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/refs/branches/{name}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-refs/#api-repositories-workspace-repo-slug-refs-branches-name-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/refs` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-refs/#api-repositories-workspace-repo-slug-refs-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:repository:bitbucket` |
| List | `read:repository:bitbucket` |

## Example Usage

```hcl
data "bitbucket_refs" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
  name = "main"
}

output "refs_response" {
  value = data.bitbucket_refs.example.api_response
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.
- `name` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `default_merge_strategy` (String) The default merge strategy for pull requests targeting this branch.
- `merge_strategies` (String) Available merge strategies for pull requests targeting this branch. (JSON array)
- `type` (String) type
