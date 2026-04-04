---
page_title: "bitbucket_repo_deploy_keys Data Source - bitbucket"
subcategory: "Deployments"
description: |-
  Reads Bitbucket repo-deploy-keys via the Bitbucket Cloud API.
---

# bitbucket_repo_deploy_keys (Data Source)

Reads Bitbucket repo-deploy-keys via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/deploy-keys/{key_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-deployments/#api-repositories-workspace-repo-slug-deploy-keys-key-id-get) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/deploy-keys` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-deployments/#api-repositories-workspace-repo-slug-deploy-keys-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `admin:repository:bitbucket` |
| List | `admin:repository:bitbucket` |

## Example Usage

```hcl
data "bitbucket_repo_deploy_keys" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}

output "repo_deploy_keys_response" {
  value = data.bitbucket_repo_deploy_keys.example.api_response
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `key_id` (String) Path parameter. Provide to fetch a specific resource; omit to list all.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `added_on` (String) added_on
- `comment` (String) The comment parsed from the deploy key (if present)
- `key` (String) The deploy key value.
- `label` (String) The user-defined label for the deploy key
- `last_used` (String) last_used
