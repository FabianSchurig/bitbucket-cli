---
page_title: "bitbucket_default_reviewers Data Source - bitbucket"
subcategory: ""
description: |-
  Reads Bitbucket default-reviewers via the Bitbucket Cloud API.
---

# bitbucket_default_reviewers (Data Source)

Reads Bitbucket default-reviewers via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/default-reviewers/{target_username}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-default-reviewers-target-username-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/default-reviewers` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pullrequests/#api-repositories-workspace-repo-slug-default-reviewers-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `read:pullrequest:bitbucket` |
| List | `read:pullrequest:bitbucket` |

## Example Usage

```hcl
data "bitbucket_default_reviewers" "example" {
  repo_slug = "my-repo"
  target_username = "jdoe"
  workspace = "my-workspace"
}

output "default_reviewers_response" {
  value = data.bitbucket_default_reviewers.example.api_response
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `target_username` (String) Path parameter.
- `workspace` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `created_on` (String) created_on
- `display_name` (String) display_name
- `uuid` (String) uuid
