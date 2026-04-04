---
page_title: "bitbucket_pipeline_config Resource - bitbucket"
subcategory: "Pipelines"
description: |-
  Manages Bitbucket pipeline-config via the Bitbucket Cloud API.
---

# bitbucket_pipeline_config (Resource)

Manages Bitbucket pipeline-config via the Bitbucket Cloud API.

## CRUD Operations
- **Read**: Supported
- **Update**: Supported

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines_config` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-get) |
| Update | `PUT` | `/repositories/{workspace}/{repo_slug}/pipelines_config` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-put) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `admin:repository:bitbucket` |
| Update | `admin:repository:bitbucket` |

## Example Usage

```hcl
resource "bitbucket_pipeline_config" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.

### Optional
- `enabled` (String) Whether Pipelines is enabled for the repository. (also computed from API response)
- `repository_description` (String) repository.description (also computed from API response)
- `repository_fork_policy` (String)  (also computed from API response)
- `repository_full_name` (String) The concatenation of the repository owner's username and the slugified name, e.g. "evzijst/interruptingcow". This is the same string used in Bitbucket URLs. (also computed from API response)
- `repository_has_issues` (String)  (also computed from API response)
- `repository_has_wiki` (String)  (also computed from API response)
- `repository_is_private` (String) repository.is_private (also computed from API response)
- `repository_language` (String) repository.language (also computed from API response)
- `repository_name` (String) repository.name (also computed from API response)
- `repository_scm` (String) [git] (also computed from API response)
- `repository_size` (String) repository.size (also computed from API response)
- `repository_uuid` (String) The repository's immutable id. This can be used as a substitute for the slug segment in URLs. Doing this guarantees your URLs will survive renaming of the repository by its owner, or even transfer of the repository to a different user. (also computed from API response)
- `request_body` (String) Raw JSON request body for create/update operations. Use `jsonencode({...})` to pass fields not exposed as individual attributes.

### Read-Only

- `id` (String) Resource identifier (extracted from API response).
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `repository_created_on` (String) repository.created_on
- `repository_updated_on` (String) repository.updated_on
