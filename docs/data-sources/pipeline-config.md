---
page_title: "bitbucket_pipeline_config Data Source - bitbucket"
subcategory: "Pipelines"
description: |-
  Reads Bitbucket pipeline-config via the Bitbucket Cloud API.
---

# bitbucket_pipeline_config (Data Source)

Reads Bitbucket pipeline-config via the Bitbucket Cloud API.

## API Endpoints

| Operation | Method | Path | API Docs |
|-----------|--------|------|----------|
| Read | `GET` | `/repositories/{workspace}/{repo_slug}/pipelines_config` | [View](https://developer.atlassian.com/cloud/bitbucket/rest/api-group-pipelines/#api-repositories-workspace-repo-slug-pipelines-config-get) |

## Required Permissions (OAuth2 Scopes)

| Operation | Required Scopes |
|-----------|----------------|
| Read | `admin:repository:bitbucket` |

## Example Usage

```hcl
data "bitbucket_pipeline_config" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}

output "pipeline_config_response" {
  value = data.bitbucket_pipeline_config.example.api_response
}
```

## Schema

### Required
- `workspace` (String) Path parameter.
- `repo_slug` (String) Path parameter.

### Read-Only

- `id` (String) Resource identifier.
- `api_response` (String) The raw JSON response from the Bitbucket API.
- `repository_created_on` (String) repository.created_on
- `repository_updated_on` (String) repository.updated_on
- `enabled` (String) Whether Pipelines is enabled for the repository.
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
- `repository_uuid` (String) The repository's immutable id. This can be used as a substitute for the slug segment in URLs. Doing this guarantees your URLs will survive renaming of the repository by its owner, or even transfer of the repository to a different user.
