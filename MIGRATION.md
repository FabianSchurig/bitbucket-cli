# Migration from `DrFaust92/terraform-provider-bitbucket`

This guide compares the legacy hand-written provider with the generated `FabianSchurig/bitbucket` provider in this repository.

It intentionally avoids generated field-by-field or HCL diffs. Nested fields, computed attributes, and generated doc structure can otherwise produce misleading migration advice.

It was generated with `python3 scripts/gen_migration.py --output MIGRATION.md`, using:

- current docs from `./docs/`
- legacy docs and source from `https://github.com/DrFaust92/terraform-provider-bitbucket/tree/master`

## What changes first

1. Switch the provider source to `FabianSchurig/bitbucket`.
2. Update provider authentication fields.
3. Rename legacy resources/data sources to the generated equivalents below.
4. Rename common path inputs like `owner` → `workspace` and `repository` → `repo_slug`.
5. Review objects that split into multiple generated resources, especially repositories and variables.

## Provider block changes

### Example

```hcl
terraform {
  required_providers {
    bitbucket = {
      source = "FabianSchurig/bitbucket"
    }
  }
}

provider "bitbucket" {
  username = var.bitbucket_username # optional for workspace/repo access tokens
  token    = var.bitbucket_token
}
```

### Provider-level renames and removals

| Legacy | New | Notes |
|---|---|---|
| Provider `password` | Provider `token` | The new provider only accepts `token`; `BITBUCKET_PASSWORD` is replaced by `BITBUCKET_TOKEN`. |
| Provider `oauth_client_id`, `oauth_client_secret`, `oauth_token` | No direct equivalent | The generated provider currently supports API tokens and workspace/repository access tokens only. |
| `owner` | `workspace` | Most repository/project scoped resources renamed the workspace path parameter to match Bitbucket Cloud OpenAPI naming. |
| `repository` or legacy repository name/slug fields | `repo_slug` | The generated provider consistently uses the Bitbucket path parameter name `repo_slug`. |
| Singular resource names like `bitbucket_repository` | Plural/group-based names like `bitbucket_repos` | Generated resources follow API operation groups instead of the legacy hand-written naming scheme. |

## Coverage summary

- Matched legacy resources: **24 / 26**
- Legacy-only resources: **2**
- New-only resources: **30**
- Matched legacy data sources: **12 / 16**
- Legacy-only data sources: **4**
- New-only data sources: **45**

## Quick rename table for matched resources

| Legacy resource | New resource(s) |
|---|---|
| `bitbucket_branch_restriction` | `bitbucket_branch_restrictions` |
| `bitbucket_branching_model` | `bitbucket_branching_model` |
| `bitbucket_commit_file` | `bitbucket_commit_file` |
| `bitbucket_default_reviewers` | `bitbucket_default_reviewers` |
| `bitbucket_deploy_key` | `bitbucket_repo_deploy_keys` |
| `bitbucket_deployment` | `bitbucket_deployments` |
| `bitbucket_deployment_variable` | `bitbucket_deployment_variables` |
| `bitbucket_forked_repository` | `bitbucket_forked_repository` |
| `bitbucket_hook` | `bitbucket_hooks` |
| `bitbucket_pipeline_schedule` | `bitbucket_pipeline_schedules` |
| `bitbucket_pipeline_ssh_key` | `bitbucket_pipeline_ssh_keys` |
| `bitbucket_pipeline_ssh_known_host` | `bitbucket_pipeline_known_hosts` |
| `bitbucket_project` | `bitbucket_projects` |
| `bitbucket_project_branching_model` | `bitbucket_project_branching_model` |
| `bitbucket_project_default_reviewers` | `bitbucket_project_default_reviewers` |
| `bitbucket_project_group_permission` | `bitbucket_project_group_permissions` |
| `bitbucket_project_user_permission` | `bitbucket_project_user_permissions` |
| `bitbucket_repository` | `bitbucket_repos`, `bitbucket_repo_settings`, `bitbucket_pipeline_config` |
| `bitbucket_repository_group_permission` | `bitbucket_repo_group_permissions` |
| `bitbucket_repository_user_permission` | `bitbucket_repo_user_permissions` |
| `bitbucket_repository_variable` | `bitbucket_pipeline_variables` |
| `bitbucket_ssh_key` | `bitbucket_ssh_keys` |
| `bitbucket_workspace_hook` | `bitbucket_workspace_hooks` |
| `bitbucket_workspace_variable` | `bitbucket_workspace_pipeline_variables` |

## Quick rename table for matched data sources

| Legacy data source | New data source(s) |
|---|---|
| `bitbucket_current_user` | `bitbucket_current_user` |
| `bitbucket_deployment` | `bitbucket_deployments` |
| `bitbucket_deployments` | `bitbucket_deployments` |
| `bitbucket_file` | `bitbucket_commit_file` |
| `bitbucket_hook_types` | `bitbucket_hook_types` |
| `bitbucket_pipeline_oidc_config` | `bitbucket_pipeline_oidc` |
| `bitbucket_pipeline_oidc_config_keys` | `bitbucket_pipeline_oidc_keys` |
| `bitbucket_project` | `bitbucket_projects` |
| `bitbucket_repository` | `bitbucket_repos` |
| `bitbucket_user` | `bitbucket_users` |
| `bitbucket_workspace` | `bitbucket_workspaces` |
| `bitbucket_workspace_members` | `bitbucket_workspace_members` |

## Matched legacy resources

### `bitbucket_branch_restriction`

- New equivalent(s): `bitbucket_branch_restrictions`
- Legacy docs: [`bitbucket_branch_restriction`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/branch_restriction.md)
- New docs: [`bitbucket_branch_restrictions`](./docs/resources/bitbucket_branch_restrictions.md)

### `bitbucket_branching_model`

- New equivalent(s): `bitbucket_branching_model`
- Legacy docs: [`bitbucket_branching_model`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/branching_model.md)
- New docs: [`bitbucket_branching_model`](./docs/resources/bitbucket_branching_model.md)

### `bitbucket_commit_file`

- New equivalent(s): `bitbucket_commit_file`
- Legacy docs: [`bitbucket_commit_file`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/commit_file.md)
- New docs: [`bitbucket_commit_file`](./docs/resources/bitbucket_commit_file.md)

### `bitbucket_default_reviewers`

- New equivalent(s): `bitbucket_default_reviewers`
- Legacy docs: [`bitbucket_default_reviewers`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/default_reviewers.md)
- New docs: [`bitbucket_default_reviewers`](./docs/resources/bitbucket_default_reviewers.md)

### `bitbucket_deploy_key`

- New equivalent(s): `bitbucket_repo_deploy_keys`
- Legacy docs: [`bitbucket_deploy_key`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/deploy_key.md)
- New docs: [`bitbucket_repo_deploy_keys`](./docs/resources/bitbucket_repo_deploy_keys.md)
- Notes: The generated provider exposes deploy keys as `bitbucket_repo_deploy_keys` and also has separate project-level deploy key resources.

### `bitbucket_deployment`

- New equivalent(s): `bitbucket_deployments`
- Legacy docs: [`bitbucket_deployment`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/deployment.md)
- New docs: [`bitbucket_deployments`](./docs/resources/bitbucket_deployments.md)

### `bitbucket_deployment_variable`

- New equivalent(s): `bitbucket_deployment_variables`
- Legacy docs: [`bitbucket_deployment_variable`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/deployment_variable.md)
- New docs: [`bitbucket_deployment_variables`](./docs/resources/bitbucket_deployment_variables.md)

### `bitbucket_forked_repository`

- New equivalent(s): `bitbucket_forked_repository`
- Legacy docs: [`bitbucket_forked_repository`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/forked_repository.md)
- New docs: [`bitbucket_forked_repository`](./docs/resources/bitbucket_forked_repository.md)

### `bitbucket_hook`

- New equivalent(s): `bitbucket_hooks`
- Legacy docs: [`bitbucket_hook`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/hook.md)
- New docs: [`bitbucket_hooks`](./docs/resources/bitbucket_hooks.md)

### `bitbucket_pipeline_schedule`

- New equivalent(s): `bitbucket_pipeline_schedules`
- Legacy docs: [`bitbucket_pipeline_schedule`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/pipeline_schedule.md)
- New docs: [`bitbucket_pipeline_schedules`](./docs/resources/bitbucket_pipeline_schedules.md)

### `bitbucket_pipeline_ssh_key`

- New equivalent(s): `bitbucket_pipeline_ssh_keys`
- Legacy docs: [`bitbucket_pipeline_ssh_key`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/pipeline_ssh_key.md)
- New docs: [`bitbucket_pipeline_ssh_keys`](./docs/resources/bitbucket_pipeline_ssh_keys.md)

### `bitbucket_pipeline_ssh_known_host`

- New equivalent(s): `bitbucket_pipeline_known_hosts`
- Legacy docs: [`bitbucket_pipeline_ssh_known_host`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/pipeline_ssh_known_host.md)
- New docs: [`bitbucket_pipeline_known_hosts`](./docs/resources/bitbucket_pipeline_known_hosts.md)

### `bitbucket_project`

- New equivalent(s): `bitbucket_projects`
- Legacy docs: [`bitbucket_project`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/project.md)
- New docs: [`bitbucket_projects`](./docs/resources/bitbucket_projects.md)

### `bitbucket_project_branching_model`

- New equivalent(s): `bitbucket_project_branching_model`
- Legacy docs: [`bitbucket_project_branching_model`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/project_branching_model.md)
- New docs: [`bitbucket_project_branching_model`](./docs/resources/bitbucket_project_branching_model.md)

### `bitbucket_project_default_reviewers`

- New equivalent(s): `bitbucket_project_default_reviewers`
- Legacy docs: [`bitbucket_project_default_reviewers`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/project_default_reviewers.md)
- New docs: [`bitbucket_project_default_reviewers`](./docs/resources/bitbucket_project_default_reviewers.md)

### `bitbucket_project_group_permission`

- New equivalent(s): `bitbucket_project_group_permissions`
- Legacy docs: [`bitbucket_project_group_permission`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/project_group_permission.md)
- New docs: [`bitbucket_project_group_permissions`](./docs/resources/bitbucket_project_group_permissions.md)

### `bitbucket_project_user_permission`

- New equivalent(s): `bitbucket_project_user_permissions`
- Legacy docs: [`bitbucket_project_user_permission`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/project_user_permission.md)
- New docs: [`bitbucket_project_user_permissions`](./docs/resources/bitbucket_project_user_permissions.md)

### `bitbucket_repository`

- New equivalent(s): `bitbucket_repos`, `bitbucket_repo_settings`, `bitbucket_pipeline_config`
- Legacy docs: [`bitbucket_repository`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/repository.md)
- New docs: [`bitbucket_repos`](./docs/resources/bitbucket_repos.md), [`bitbucket_repo_settings`](./docs/resources/bitbucket_repo_settings.md), [`bitbucket_pipeline_config`](./docs/resources/bitbucket_pipeline_config.md)
- Notes: The legacy repository resource bundled core repository CRUD, pipeline enablement, and override-settings flags. In the new provider, core CRUD stays on `bitbucket_repos`, pipeline enablement moves to `bitbucket_pipeline_config`, and repository settings have their own `bitbucket_repo_settings` resource.

### `bitbucket_repository_group_permission`

- New equivalent(s): `bitbucket_repo_group_permissions`
- Legacy docs: [`bitbucket_repository_group_permission`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/repository_group_permission.md)
- New docs: [`bitbucket_repo_group_permissions`](./docs/resources/bitbucket_repo_group_permissions.md)

### `bitbucket_repository_user_permission`

- New equivalent(s): `bitbucket_repo_user_permissions`
- Legacy docs: [`bitbucket_repository_user_permission`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/repository_user_permission.md)
- New docs: [`bitbucket_repo_user_permissions`](./docs/resources/bitbucket_repo_user_permissions.md)

### `bitbucket_repository_variable`

- New equivalent(s): `bitbucket_pipeline_variables`
- Legacy docs: [`bitbucket_repository_variable`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/repository_variable.md)
- New docs: [`bitbucket_pipeline_variables`](./docs/resources/bitbucket_pipeline_variables.md)
- Notes: Legacy repository variables map to the pipelines variable API. Use `bitbucket_pipeline_variables` and rename `owner`/`repository` to `workspace`/`repo_slug`.

### `bitbucket_ssh_key`

- New equivalent(s): `bitbucket_ssh_keys`
- Legacy docs: [`bitbucket_ssh_key`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/ssh_key.md)
- New docs: [`bitbucket_ssh_keys`](./docs/resources/bitbucket_ssh_keys.md)

### `bitbucket_workspace_hook`

- New equivalent(s): `bitbucket_workspace_hooks`
- Legacy docs: [`bitbucket_workspace_hook`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/workspace_hook.md)
- New docs: [`bitbucket_workspace_hooks`](./docs/resources/bitbucket_workspace_hooks.md)

### `bitbucket_workspace_variable`

- New equivalent(s): `bitbucket_workspace_pipeline_variables`
- Legacy docs: [`bitbucket_workspace_variable`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/workspace_variable.md)
- New docs: [`bitbucket_workspace_pipeline_variables`](./docs/resources/bitbucket_workspace_pipeline_variables.md)
- Notes: Workspace variables now live under the pipelines API as `bitbucket_workspace_pipeline_variables`.

## Legacy-only resources

### `bitbucket_group`

- New equivalent(s): none
- Legacy docs: [`bitbucket_group`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/group.md)
- Notes: Workspace group management is not currently exposed by the generated provider because those endpoints are not represented in the generated Terraform docs.

### `bitbucket_group_membership`

- New equivalent(s): none
- Legacy docs: [`bitbucket_group_membership`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/resources/group_membership.md)
- Notes: Group membership management is not currently exposed by the generated provider.

## Matched legacy data sources

### `bitbucket_current_user`

- New equivalent(s): `bitbucket_current_user`
- Legacy docs: [`bitbucket_current_user`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/data-sources/current_user.md)
- New docs: [`bitbucket_current_user`](./docs/data-sources/bitbucket_current_user.md)
- Notes: The legacy data source also fetched `/user/emails`. The generated provider splits that into `bitbucket_current_user` plus `bitbucket_user_emails` when you need email addresses.

### `bitbucket_deployment`

- New equivalent(s): `bitbucket_deployments`
- Legacy docs: [`bitbucket_deployment`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/data-sources/deployment.md)
- New docs: [`bitbucket_deployments`](./docs/data-sources/bitbucket_deployments.md)
- Notes: Use `bitbucket_deployments` with the identifying path parameters for a single deployment; omit the single-resource expectation and treat the response as the generic deployment payload.

### `bitbucket_deployments`

- New equivalent(s): `bitbucket_deployments`
- Legacy docs: [`bitbucket_deployments`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/data-sources/deployments.md)
- New docs: [`bitbucket_deployments`](./docs/data-sources/bitbucket_deployments.md)

### `bitbucket_file`

- New equivalent(s): `bitbucket_commit_file`
- Legacy docs: [`bitbucket_file`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/data-sources/file.md)
- New docs: [`bitbucket_commit_file`](./docs/data-sources/bitbucket_commit_file.md)
- Notes: The legacy `bitbucket_file` data source maps most closely to `bitbucket_commit_file`, which reads file content via the commit-file endpoint.

### `bitbucket_hook_types`

- New equivalent(s): `bitbucket_hook_types`
- Legacy docs: [`bitbucket_hook_types`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/data-sources/hook_types.md)
- New docs: [`bitbucket_hook_types`](./docs/data-sources/bitbucket_hook_types.md)

### `bitbucket_pipeline_oidc_config`

- New equivalent(s): `bitbucket_pipeline_oidc`
- Legacy docs: [`bitbucket_pipeline_oidc_config`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/data-sources/pipeline_oidc_config.md)
- New docs: [`bitbucket_pipeline_oidc`](./docs/data-sources/bitbucket_pipeline_oidc.md)

### `bitbucket_pipeline_oidc_config_keys`

- New equivalent(s): `bitbucket_pipeline_oidc_keys`
- Legacy docs: [`bitbucket_pipeline_oidc_config_keys`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/data-sources/pipeline_oidc_config_keys.md)
- New docs: [`bitbucket_pipeline_oidc_keys`](./docs/data-sources/bitbucket_pipeline_oidc_keys.md)

### `bitbucket_project`

- New equivalent(s): `bitbucket_projects`
- Legacy docs: [`bitbucket_project`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/data-sources/project.md)
- New docs: [`bitbucket_projects`](./docs/data-sources/bitbucket_projects.md)

### `bitbucket_repository`

- New equivalent(s): `bitbucket_repos`
- Legacy docs: [`bitbucket_repository`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/data-sources/repository.md)
- New docs: [`bitbucket_repos`](./docs/data-sources/bitbucket_repos.md)

### `bitbucket_user`

- New equivalent(s): `bitbucket_users`
- Legacy docs: [`bitbucket_user`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/data-sources/user.md)
- New docs: [`bitbucket_users`](./docs/data-sources/bitbucket_users.md)

### `bitbucket_workspace`

- New equivalent(s): `bitbucket_workspaces`
- Legacy docs: [`bitbucket_workspace`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/data-sources/workspace.md)
- New docs: [`bitbucket_workspaces`](./docs/data-sources/bitbucket_workspaces.md)

### `bitbucket_workspace_members`

- New equivalent(s): `bitbucket_workspace_members`
- Legacy docs: [`bitbucket_workspace_members`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/data-sources/workspace_members.md)
- New docs: [`bitbucket_workspace_members`](./docs/data-sources/bitbucket_workspace_members.md)

## Legacy-only data sources

### `bitbucket_group`

- New equivalent(s): none
- Legacy docs: [`bitbucket_group`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/data-sources/group.md)
- Notes: Group lookup is not currently exposed by the generated provider.

### `bitbucket_group_members`

- New equivalent(s): none
- Legacy docs: [`bitbucket_group_members`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/data-sources/group_members.md)
- Notes: Group member lookup is not currently exposed by the generated provider.

### `bitbucket_groups`

- New equivalent(s): none
- Legacy docs: [`bitbucket_groups`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/data-sources/groups.md)
- Notes: Group listing is not currently exposed by the generated provider.

### `bitbucket_ip_ranges`

- New equivalent(s): none
- Legacy docs: [`bitbucket_ip_ranges`](https://github.com/DrFaust92/terraform-provider-bitbucket/blob/master/docs/data-sources/ip_ranges.md)
- Notes: The generated provider does not currently expose Bitbucket IP ranges as a Terraform data source.

## New provider-only resources

- [`bitbucket_addon`](./docs/resources/bitbucket_addon.md)
- [`bitbucket_annotations`](./docs/resources/bitbucket_annotations.md)
- [`bitbucket_commit_statuses`](./docs/resources/bitbucket_commit_statuses.md)
- [`bitbucket_commits`](./docs/resources/bitbucket_commits.md)
- [`bitbucket_current_user`](./docs/resources/bitbucket_current_user.md)
- [`bitbucket_downloads`](./docs/resources/bitbucket_downloads.md)
- [`bitbucket_gpg_keys`](./docs/resources/bitbucket_gpg_keys.md)
- [`bitbucket_hook_types`](./docs/resources/bitbucket_hook_types.md)
- [`bitbucket_issue_comments`](./docs/resources/bitbucket_issue_comments.md)
- [`bitbucket_issues`](./docs/resources/bitbucket_issues.md)
- [`bitbucket_pipeline_caches`](./docs/resources/bitbucket_pipeline_caches.md)
- [`bitbucket_pipeline_oidc`](./docs/resources/bitbucket_pipeline_oidc.md)
- [`bitbucket_pipeline_oidc_keys`](./docs/resources/bitbucket_pipeline_oidc_keys.md)
- [`bitbucket_pipelines`](./docs/resources/bitbucket_pipelines.md)
- [`bitbucket_pr`](./docs/resources/bitbucket_pr.md)
- [`bitbucket_pr_comments`](./docs/resources/bitbucket_pr_comments.md)
- [`bitbucket_project_deploy_keys`](./docs/resources/bitbucket_project_deploy_keys.md)
- [`bitbucket_properties`](./docs/resources/bitbucket_properties.md)
- [`bitbucket_refs`](./docs/resources/bitbucket_refs.md)
- [`bitbucket_repo_runners`](./docs/resources/bitbucket_repo_runners.md)
- [`bitbucket_reports`](./docs/resources/bitbucket_reports.md)
- [`bitbucket_search`](./docs/resources/bitbucket_search.md)
- [`bitbucket_snippets`](./docs/resources/bitbucket_snippets.md)
- [`bitbucket_tags`](./docs/resources/bitbucket_tags.md)
- [`bitbucket_user_emails`](./docs/resources/bitbucket_user_emails.md)
- [`bitbucket_users`](./docs/resources/bitbucket_users.md)
- [`bitbucket_workspace_members`](./docs/resources/bitbucket_workspace_members.md)
- [`bitbucket_workspace_permissions`](./docs/resources/bitbucket_workspace_permissions.md)
- [`bitbucket_workspace_runners`](./docs/resources/bitbucket_workspace_runners.md)
- [`bitbucket_workspaces`](./docs/resources/bitbucket_workspaces.md)

## New provider-only data sources

- [`bitbucket_addon`](./docs/data-sources/bitbucket_addon.md)
- [`bitbucket_annotations`](./docs/data-sources/bitbucket_annotations.md)
- [`bitbucket_branch_restrictions`](./docs/data-sources/bitbucket_branch_restrictions.md)
- [`bitbucket_branching_model`](./docs/data-sources/bitbucket_branching_model.md)
- [`bitbucket_commit_statuses`](./docs/data-sources/bitbucket_commit_statuses.md)
- [`bitbucket_commits`](./docs/data-sources/bitbucket_commits.md)
- [`bitbucket_default_reviewers`](./docs/data-sources/bitbucket_default_reviewers.md)
- [`bitbucket_deployment_variables`](./docs/data-sources/bitbucket_deployment_variables.md)
- [`bitbucket_downloads`](./docs/data-sources/bitbucket_downloads.md)
- [`bitbucket_forked_repository`](./docs/data-sources/bitbucket_forked_repository.md)
- [`bitbucket_gpg_keys`](./docs/data-sources/bitbucket_gpg_keys.md)
- [`bitbucket_hooks`](./docs/data-sources/bitbucket_hooks.md)
- [`bitbucket_issue_comments`](./docs/data-sources/bitbucket_issue_comments.md)
- [`bitbucket_issues`](./docs/data-sources/bitbucket_issues.md)
- [`bitbucket_pipeline_caches`](./docs/data-sources/bitbucket_pipeline_caches.md)
- [`bitbucket_pipeline_config`](./docs/data-sources/bitbucket_pipeline_config.md)
- [`bitbucket_pipeline_known_hosts`](./docs/data-sources/bitbucket_pipeline_known_hosts.md)
- [`bitbucket_pipeline_schedules`](./docs/data-sources/bitbucket_pipeline_schedules.md)
- [`bitbucket_pipeline_ssh_keys`](./docs/data-sources/bitbucket_pipeline_ssh_keys.md)
- [`bitbucket_pipeline_variables`](./docs/data-sources/bitbucket_pipeline_variables.md)
- [`bitbucket_pipelines`](./docs/data-sources/bitbucket_pipelines.md)
- [`bitbucket_pr`](./docs/data-sources/bitbucket_pr.md)
- [`bitbucket_pr_comments`](./docs/data-sources/bitbucket_pr_comments.md)
- [`bitbucket_project_branching_model`](./docs/data-sources/bitbucket_project_branching_model.md)
- [`bitbucket_project_default_reviewers`](./docs/data-sources/bitbucket_project_default_reviewers.md)
- [`bitbucket_project_deploy_keys`](./docs/data-sources/bitbucket_project_deploy_keys.md)
- [`bitbucket_project_group_permissions`](./docs/data-sources/bitbucket_project_group_permissions.md)
- [`bitbucket_project_user_permissions`](./docs/data-sources/bitbucket_project_user_permissions.md)
- [`bitbucket_properties`](./docs/data-sources/bitbucket_properties.md)
- [`bitbucket_refs`](./docs/data-sources/bitbucket_refs.md)
- [`bitbucket_repo_deploy_keys`](./docs/data-sources/bitbucket_repo_deploy_keys.md)
- [`bitbucket_repo_group_permissions`](./docs/data-sources/bitbucket_repo_group_permissions.md)
- [`bitbucket_repo_runners`](./docs/data-sources/bitbucket_repo_runners.md)
- [`bitbucket_repo_settings`](./docs/data-sources/bitbucket_repo_settings.md)
- [`bitbucket_repo_user_permissions`](./docs/data-sources/bitbucket_repo_user_permissions.md)
- [`bitbucket_reports`](./docs/data-sources/bitbucket_reports.md)
- [`bitbucket_search`](./docs/data-sources/bitbucket_search.md)
- [`bitbucket_snippets`](./docs/data-sources/bitbucket_snippets.md)
- [`bitbucket_ssh_keys`](./docs/data-sources/bitbucket_ssh_keys.md)
- [`bitbucket_tags`](./docs/data-sources/bitbucket_tags.md)
- [`bitbucket_user_emails`](./docs/data-sources/bitbucket_user_emails.md)
- [`bitbucket_workspace_hooks`](./docs/data-sources/bitbucket_workspace_hooks.md)
- [`bitbucket_workspace_permissions`](./docs/data-sources/bitbucket_workspace_permissions.md)
- [`bitbucket_workspace_pipeline_variables`](./docs/data-sources/bitbucket_workspace_pipeline_variables.md)
- [`bitbucket_workspace_runners`](./docs/data-sources/bitbucket_workspace_runners.md)

## Can this be automated?

Only partly. The rename tables are useful, but the actual migration still needs a human review against the authoritative docs.

Good candidates for an automated rewrite later:

- provider source replacement
- provider auth field rename (`password` → `token`)
- direct resource/data source renames where there is a 1:1 mapping

Cases that still need manual review:

- legacy objects that split into multiple generated resources
- objects missing from one provider or the other
- nested or computed fields
- fields whose semantics changed even when the name looks similar

