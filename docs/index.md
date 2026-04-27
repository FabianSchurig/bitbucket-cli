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
- `csrf_token` (String, Sensitive) CSRF token (`csrftoken` browser cookie) used to authenticate against Bitbucket's internal API (`https://bitbucket.org/!api/internal/...`). Required for resources marked **Internal API**, which reject HTTP Basic Auth. Can also be set via `BITBUCKET_CSRF_TOKEN`. Must be paired with `cloud_session_token`.
- `cloud_session_token` (String, Sensitive) Cloud session token (`cloud.session.token` browser cookie) used to authenticate against Bitbucket's internal API. Required for resources marked **Internal API**. Can also be set via `BITBUCKET_CLOUD_SESSION_TOKEN`. Must be paired with `csrf_token`.

## Authenticating against the internal API

A handful of resources (e.g. `bitbucket_project_branch_restrictions`) are
backed by Bitbucket's undocumented internal API at
`https://bitbucket.org/!api/internal/`. That endpoint **does not accept
HTTP Basic Auth** — it only accepts the same browser cookies the Bitbucket
web UI sends. Configure them like this:

```hcl
provider "bitbucket" {
  # Public REST API (optional if you only use internal-API resources)
  username = "your-email@example.com"
  token    = "your-api-token"

  # Internal API (required for resources marked "Internal API")
  csrf_token          = "value of the csrftoken cookie"
  cloud_session_token = "value of the cloud.session.token cookie"
}
```

Or via environment variables:

```bash
export BITBUCKET_CSRF_TOKEN="..."
export BITBUCKET_CLOUD_SESSION_TOKEN="..."
```

You can grab both values from your browser's developer tools while logged
in to bitbucket.org (Application → Cookies → bitbucket.org). The provider
inspects each request URL: requests to `/!api/internal/` automatically use
cookie auth (and `X-CSRFToken`, `X-Requested-With`, `Referer`,
`Sec-Fetch-*` headers); all other requests use Basic Auth.

~> **The `cloud.session.token` cookie is short-lived** — typically about a
month before Bitbucket invalidates it. Because of that, internal-API
resources (grouped under **Experimental** below) are best used **manually
and interactively**: copy a fresh cookie from your browser right before you
run `terraform apply`. They are generally not suitable for unattended CI
pipelines that may need to run weeks or months after the cookie was captured.

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


### Experimental

Resources in this group wrap **undocumented internal Bitbucket endpoints**
(`https://bitbucket.org/!api/internal/...`). They are not part of the public
REST API and have several important caveats:

- **Not auto-synced.** The rest of this provider is regenerated daily from
  Atlassian's published OpenAPI spec; internal-API resources are hand-curated
  and updated less frequently. Atlassian can change or remove these endpoints
  without notice.
- **Browser-cookie auth only.** They reject HTTP Basic Auth — you must
  configure `csrf_token` and `cloud_session_token` (or the matching
  `BITBUCKET_*` env vars). See
  [Authenticating against the internal API](#authenticating-against-the-internal-api).
- **Short-lived session token.** The `cloud.session.token` cookie typically
  expires after about a month, after which Terraform runs that touch these
  resources will start returning 401. Because of this, the practical
  recommendation is to use experimental resources **manually / interactively**:
  copy fresh values from your browser's developer tools (Application → Cookies
  → bitbucket.org), run `terraform apply`, then unset the variables. They
  are generally not suitable for unattended CI pipelines.

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
