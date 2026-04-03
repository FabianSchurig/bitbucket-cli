---
page_title: "bitbucket Provider"
subcategory: ""
description: |-
  Terraform provider for Bitbucket Cloud. Auto-generated from the Bitbucket OpenAPI spec.
---

# bitbucket Provider

Terraform provider for Bitbucket Cloud, exposing all Bitbucket API operations as
generic resources and data sources. Auto-generated from the Bitbucket OpenAPI spec.

## Authentication

The provider authenticates via HTTP Basic Auth using an Atlassian API token.
Create a token at [id.atlassian.com/manage-profile/security/api-tokens](https://id.atlassian.com/manage-profile/security/api-tokens).

### Atlassian API Token (recommended)

```hcl
provider "bitbucket" {
  username = "your-email@example.com"  # Atlassian account email
  token    = "your-api-token"
}
```

Or via environment variables:

```bash
export BITBUCKET_USERNAME="your-email@example.com"
export BITBUCKET_TOKEN="your-api-token"
```

### Workspace Access Token

For workspace/repository access tokens, only the token is needed:

```hcl
provider "bitbucket" {
  token = "your-workspace-access-token"
}
```

## Example Usage

```hcl
terraform {
  required_providers {
    bitbucket = {
      source = "FabianSchurig/bitbucket"
    }
  }
}

provider "bitbucket" {
  # Authentication via environment variables recommended
}

# Read a repository
data "bitbucket_repos" "example" {
  workspace = "my-workspace"
  repo_slug = "my-repo"
}

# Output the API response
output "repo_info" {
  value = data.bitbucket_repos.example.api_response
}
```

## Schema

### Optional

- `username` (String) Bitbucket username (Atlassian account email for API tokens). Can also be set via `BITBUCKET_USERNAME` environment variable.
- `token` (String, Sensitive) Bitbucket API token (Atlassian API token or workspace access token). Can also be set via `BITBUCKET_TOKEN` environment variable.
- `base_url` (String) Base URL for the Bitbucket API. Defaults to `https://api.bitbucket.org/2.0`.

## Resources and Data Sources

This provider auto-generates resources and data sources for all Bitbucket API
operation groups. Each resource group maps to a set of CRUD operations.

| Resource | Data Source | CRUD |
|----------|-------------|------|
| `bitbucket_addon` | `bitbucket_addon` | UDL |
| `bitbucket_branch_restrictions` | `bitbucket_branch_restrictions` | CRUDL |
| `bitbucket_branching_model` | `bitbucket_branching_model` | RU |
| `bitbucket_commit_statuses` | `bitbucket_commit_statuses` | CRUL |
| `bitbucket_commits` | `bitbucket_commits` | RL |
| `bitbucket_deployments` | `bitbucket_deployments` | CRDL |
| `bitbucket_downloads` | `bitbucket_downloads` | CRDL |
| `bitbucket_hooks` | `bitbucket_hooks` | CRUDL |
| `bitbucket_issues` | `bitbucket_issues` | CRUDL |
| `bitbucket_pipelines` | `bitbucket_pipelines` | CRL |
| `bitbucket_pr` | `bitbucket_pr` | CRUL |
| `bitbucket_projects` | `bitbucket_projects` | CRUDL |
| `bitbucket_properties` | `bitbucket_properties` | RUD |
| `bitbucket_refs` | `bitbucket_refs` | CRDL |
| `bitbucket_reports` | `bitbucket_reports` | CRDL |
| `bitbucket_repos` | `bitbucket_repos` | CRUDL |
| `bitbucket_search` | `bitbucket_search` | L |
| `bitbucket_snippets` | `bitbucket_snippets` | CRUDL |
| `bitbucket_users` | `bitbucket_users` | RL |
| `bitbucket_workspaces` | `bitbucket_workspaces` | RL |

All resources share the same generic schema pattern:

- **Path parameters** become required/optional string attributes
- **Body fields** become optional string attributes
- `api_response` (Computed) contains the raw JSON API response
- `id` (Computed) is extracted from the response (uuid, id, slug, or name)
