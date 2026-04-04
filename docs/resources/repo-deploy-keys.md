---
page_title: "bitbucket_repo_deploy_keys Resource - bitbucket"
subcategory: "Deployments"
description: |-
  Manages Bitbucket repo-deploy-keys via the Bitbucket Cloud API.
---

# bitbucket_repo_deploy_keys (Resource)

Manages Bitbucket repo-deploy-keys via the Bitbucket Cloud API.

## CRUD Operations
- **Create**: Supported
- **Read**: Supported
- **Update**: Supported
- **Delete**: Supported
- **List**: Supported (via data source)

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Create | `POST` | `/repositories/{workspace}/{repo_slug}/deploy-keys` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-deployments/#api-repositories-workspace-repo-slug-deploy-keys-post) |
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/deploy-keys/{key_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-deployments/#api-repositories-workspace-repo-slug-deploy-keys-key-id-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}/deploy-keys/{key_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-deployments/#api-repositories-workspace-repo-slug-deploy-keys-key-id-put) |
| Delete | `DELETE` | `/repositories/{workspace}/{repo_slug}/deploy-keys/{key_id}` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-deployments/#api-repositories-workspace-repo-slug-deploy-keys-key-id-delete) |
| List | `GET` | `/repositories/{workspace}/{repo_slug}/deploy-keys` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-deployments/#api-repositories-workspace-repo-slug-deploy-keys-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Create | `admin:repository:bitbucket`, `write:ssh-key:bitbucket` |
| Read | `admin:repository:bitbucket` |
| Update | `admin:repository:bitbucket`, `write:ssh-key:bitbucket` |
| Delete | `admin:repository:bitbucket`, `delete:ssh-key:bitbucket` |
| List | `admin:repository:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_repo_deploy_keys" "example" {
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

## Schema

### Required
- `repo_slug` (String) Path parameter.
- `workspace` (String) Path parameter.

### Optional
- `key_id` (String) Path parameter (auto-populated from API response).

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `added_on` (String) added_on
- `comment` (String) The comment parsed from the deploy key (if present)
- `key` (String) The deploy key value.
- `label` (String) The user-defined label for the deploy key
- `last_used` (String) last_used
- `owner_created_on` (String) owner.created_on
- `owner_display_name` (String) owner.display_name
- `owner_uuid` (String) owner.uuid
- `repository_created_on` (String) repository.created_on
- `repository_description` (String) repository.description
- `repository_fork_policy` (String) 
- `repository_full_name` (String) The concatenation of the repository owner's username and the slugified name, e.g. "evzijst/interruptingcow". This is the same string used in Bitbucket URLs.
- `repository_has_issues` (String) 
- `repository_has_wiki` (String) 
- `repository_is_private` (String) repository.is_private
- `repository_language` (String) repository.language
- `repository_name` (String) repository.name
- `repository_scm` (String) [git]
- `repository_size` (String) repository.size
- `repository_updated_on` (String) repository.updated_on
- `repository_uuid` (String) The repository's immutable id. This can be used as a substitute for the slug segment in URLs. Doing this guarantees your URLs will survive renaming of the repository by its owner, or even transfer of the repository to a different user.
