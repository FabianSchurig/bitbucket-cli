---
page_title: "bitbucket_branch_restrictions Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket branch-restrictions via the Bitbucket Cloud API.
---

# bitbucket_branch_restrictions (Data Source)

Reads Bitbucket branch-restrictions via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/branch-restrictions/{id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-branch-restrictions/#api-repositories-workspace-repo-slug-branch-restrictions-id-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/branch-restrictions` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-branch-restrictions/#api-repositories-workspace-repo-slug-branch-restrictions-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `admin:repository:bitbucket` |
| List | `admin:repository:bitbucket` |

## Example Usage

```hcl
data "bitbucket_branch_restrictions" "example" {
  param_id = "1"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "branch_restrictions_response" {
  value = data.bitbucket_branch_restrictions.example.api_response
}
```

## Schema

### Required
- `param_id` (String) Path parameter.
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
