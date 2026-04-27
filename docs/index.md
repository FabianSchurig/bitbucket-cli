---
page_title: "bitbucket Provider"
subcategory: ""
description: |-
  Terraform provider for Bitbucket Cloud. Auto-generated from the Bitbucket OpenAPI spec.
---

# bitbucket Provider

Terraform provider for Bitbucket Cloud, exposing all Bitbucket API operations as
generic resources and data sources. Auto-generated from the Bitbucket OpenAPI spec.

Migrating from the legacy `DrFaust92/terraform-provider-bitbucket` provider? See
[`MIGRATION.md`](https://github.com/FabianSchurig/bitbucket-cli/blob/main/MIGRATION.md).

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


### Addon

| Resource | Data Source | CRUD |
|----------|-------------|------|
| `bitbucket_addon` | `bitbucket_addon` | UDL |


### Branch Restrictions

| Resource | Data Source | CRUD |
|----------|-------------|------|
| `bitbucket_branch_restrictions` | `bitbucket_branch_restrictions` | CRUDL |


### Branching Model

| Resource | Data Source | CRUD |
|----------|-------------|------|
| `bitbucket_branching_model` | `bitbucket_branching_model` | RU |
| `bitbucket_project_branching_model` | `bitbucket_project_branching_model` | RU |


### Code Search

| Resource | Data Source | CRUD |
|----------|-------------|------|
| `bitbucket_search` | `bitbucket_search` | L |


### Commit Statuses

| Resource | Data Source | CRUD |
|----------|-------------|------|
| `bitbucket_commit_statuses` | `bitbucket_commit_statuses` | CRUL |


### Commits

| Resource | Data Source | CRUD |
|----------|-------------|------|
| `bitbucket_commits` | `bitbucket_commits` | RL |


### Deployments

| Resource | Data Source | CRUD |
|----------|-------------|------|
| `bitbucket_deployments` | `bitbucket_deployments` | CRDL |
| `bitbucket_project_deploy_keys` | `bitbucket_project_deploy_keys` | CRDL |
| `bitbucket_repo_deploy_keys` | `bitbucket_repo_deploy_keys` | CRUDL |


### Downloads

| Resource | Data Source | CRUD |
|----------|-------------|------|
| `bitbucket_downloads` | `bitbucket_downloads` | CRDL |


### Internal Project Branch Restrictions

| Resource | Data Source | CRUD |
|----------|-------------|------|
| `bitbucket_project_branch_restrictions` | `bitbucket_project_branch_restrictions` | RUL |
| `bitbucket_project_branch_restrictions_by_branch_type` | `bitbucket_project_branch_restrictions_by_branch_type` | CRUDL |
| `bitbucket_project_branch_restrictions_by_pattern` | `bitbucket_project_branch_restrictions_by_pattern` | CRUDL |


### Issues

| Resource | Data Source | CRUD |
|----------|-------------|------|
| `bitbucket_issue_comments` | `bitbucket_issue_comments` | CRUDL |
| `bitbucket_issues` | `bitbucket_issues` | CRUDL |


### Pipelines

| Resource | Data Source | CRUD |
|----------|-------------|------|
| `bitbucket_deployment_variables` | `bitbucket_deployment_variables` | CRUD |
| `bitbucket_pipeline_caches` | `bitbucket_pipeline_caches` | RDL |
| `bitbucket_pipeline_config` | `bitbucket_pipeline_config` | RU |
| `bitbucket_pipeline_known_hosts` | `bitbucket_pipeline_known_hosts` | CRUDL |
| `bitbucket_pipeline_oidc` | `bitbucket_pipeline_oidc` | R |
| `bitbucket_pipeline_oidc_keys` | `bitbucket_pipeline_oidc_keys` | R |
| `bitbucket_pipeline_schedules` | `bitbucket_pipeline_schedules` | CRUDL |
| `bitbucket_pipeline_ssh_keys` | `bitbucket_pipeline_ssh_keys` | RUD |
| `bitbucket_pipeline_variables` | `bitbucket_pipeline_variables` | CRUDL |
| `bitbucket_pipelines` | `bitbucket_pipelines` | CRL |
| `bitbucket_repo_runners` | `bitbucket_repo_runners` | CRUDL |
| `bitbucket_workspace_pipeline_variables` | `bitbucket_workspace_pipeline_variables` | CRUDL |
| `bitbucket_workspace_runners` | `bitbucket_workspace_runners` | CRUDL |


### Projects

| Resource | Data Source | CRUD |
|----------|-------------|------|
| `bitbucket_project_default_reviewers` | `bitbucket_project_default_reviewers` | CRDL |
| `bitbucket_project_group_permissions` | `bitbucket_project_group_permissions` | RUDL |
| `bitbucket_project_user_permissions` | `bitbucket_project_user_permissions` | RUDL |
| `bitbucket_projects` | `bitbucket_projects` | CRUDL |


### Properties

| Resource | Data Source | CRUD |
|----------|-------------|------|
| `bitbucket_properties` | `bitbucket_properties` | RUD |


### Pull Requests

| Resource | Data Source | CRUD |
|----------|-------------|------|
| `bitbucket_default_reviewers` | `bitbucket_default_reviewers` | CRDL |
| `bitbucket_pr` | `bitbucket_pr` | CRUL |
| `bitbucket_pr_comments` | `bitbucket_pr_comments` | CRUDL |


### Refs

| Resource | Data Source | CRUD |
|----------|-------------|------|
| `bitbucket_refs` | `bitbucket_refs` | CRDL |
| `bitbucket_tags` | `bitbucket_tags` | CRDL |


### Reports

| Resource | Data Source | CRUD |
|----------|-------------|------|
| `bitbucket_annotations` | `bitbucket_annotations` | CRDL |
| `bitbucket_reports` | `bitbucket_reports` | CRDL |


### Repositories

| Resource | Data Source | CRUD |
|----------|-------------|------|
| `bitbucket_commit_file` | `bitbucket_commit_file` | CR |
| `bitbucket_forked_repository` | `bitbucket_forked_repository` | CL |
| `bitbucket_repo_group_permissions` | `bitbucket_repo_group_permissions` | RUDL |
| `bitbucket_repo_settings` | `bitbucket_repo_settings` | RU |
| `bitbucket_repo_user_permissions` | `bitbucket_repo_user_permissions` | RUDL |
| `bitbucket_repos` | `bitbucket_repos` | CRUDL |


### Snippets

| Resource | Data Source | CRUD |
|----------|-------------|------|
| `bitbucket_snippets` | `bitbucket_snippets` | CRUDL |


### Users

| Resource | Data Source | CRUD |
|----------|-------------|------|
| `bitbucket_current_user` | `bitbucket_current_user` | R |
| `bitbucket_gpg_keys` | `bitbucket_gpg_keys` | CRDL |
| `bitbucket_ssh_keys` | `bitbucket_ssh_keys` | CRUDL |
| `bitbucket_user_emails` | `bitbucket_user_emails` | RL |
| `bitbucket_users` | `bitbucket_users` | RL |


### Webhooks

| Resource | Data Source | CRUD |
|----------|-------------|------|
| `bitbucket_hook_types` | `bitbucket_hook_types` | RL |
| `bitbucket_hooks` | `bitbucket_hooks` | CRUDL |


### Workspaces

| Resource | Data Source | CRUD |
|----------|-------------|------|
| `bitbucket_workspace_hooks` | `bitbucket_workspace_hooks` | CRUDL |
| `bitbucket_workspace_members` | `bitbucket_workspace_members` | RL |
| `bitbucket_workspace_permissions` | `bitbucket_workspace_permissions` | RL |
| `bitbucket_workspaces` | `bitbucket_workspaces` | RL |

All resources share the same generic schema pattern:

- **Path parameters** become required/optional string attributes
- **Body fields** become optional string attributes (writable)
- **Response fields** become computed string attributes (read-only, auto-populated from API response)
- Fields present in both request and response are **Optional+Computed** (can be set by user, also populated from API)
- `api_response` (Computed) contains the raw JSON API response
- `id` (Computed) is extracted from the response (uuid, id, slug, or name)
