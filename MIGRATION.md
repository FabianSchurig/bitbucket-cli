# Migration from `DrFaust92/terraform-provider-bitbucket`

This guide compares the legacy hand-written provider with the generated `FabianSchurig/bitbucket` provider in this repository.

It is intentionally a best-effort migration baseline: the generated docs sometimes list optional fields that are also computed by the API, so subtle cases still need manual review.

It was generated with `python3 scripts/gen_migration.py --output MIGRATION.md`, using:

- current docs from `/home/runner/work/bitbucket-cli/bitbucket-cli/docs/`
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
- Legacy inputs: required: `owner`, `repository`, `kind`; optional: `branch_match_kind`, `branch_type`, `pattern`, `users`, `groups`, `value`
- Legacy endpoints: `POST /repositories/{workspace}/{repo_slug}/branch-restrictions`<br>`GET /repositories/{workspace}/{repo_slug}/branch-restrictions/{id}`<br>`PUT /repositories/{workspace}/{repo_slug}/branch-restrictions/{id}`<br>`DELETE /repositories/{workspace}/{repo_slug}/branch-restrictions/{id}`
- New inputs: required: `repo_slug`, `workspace`; optional: `param_id`, `branch_match_kind`, `branch_type`, `groups`, `kind`, `pattern`, `users`, `value`, `request_body`
- New operations: `Create POST /repositories/{workspace}/{repo_slug}/branch-restrictions`<br>`Read GET /repositories/{workspace}/{repo_slug}/branch-restrictions/{id}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/branch-restrictions/{id}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/branch-restrictions/{id}`<br>`List GET /repositories/{workspace}/{repo_slug}/branch-restrictions`
- Diff summary: renamed: `owner` → `workspace`, `repository` → `repo_slug`; new-only inputs: `param_id`, `request_body`

### `bitbucket_branching_model`

- New equivalent(s): `bitbucket_branching_model`
- Legacy inputs: required: `owner`, `repository`, `branch_type`; optional: `development`, `production`
- Legacy endpoints: `Read GET /repositories/{workspace}/{repo_slug}/branching-model`<br>`Update PUT /repositories/{workspace}/{repo_slug}/branching-model/settings`
- New inputs: required: `repo_slug`, `workspace`
- New operations: `Read GET /repositories/{workspace}/{repo_slug}/branching-model`<br>`Update PUT /repositories/{workspace}/{repo_slug}/branching-model/settings`
- Diff summary: renamed: `owner` → `workspace`, `repository` → `repo_slug`; legacy-only inputs: `branch_type`, `development`, `production`

### `bitbucket_commit_file`

- New equivalent(s): `bitbucket_commit_file`
- Legacy inputs: required: `workspace`, `repo_slug`, `filename`, `content`, `commit_author`, `branch`, `commit_message`
- Legacy endpoints: `GET /repositories/{workspace}/{repo_slug}/src/{commit}/{path}`
- New inputs: required: `repo_slug`, `workspace`; optional: `commit`, `path`
- New operations: `Create POST /repositories/{workspace}/{repo_slug}/src`<br>`Read GET /repositories/{workspace}/{repo_slug}/src/{commit}/{path}`
- Diff summary: legacy-only inputs: `branch`, `commit_author`, `commit_message`, `content`, `filename`; new-only inputs: `commit`, `path`

### `bitbucket_default_reviewers`

- New equivalent(s): `bitbucket_default_reviewers`
- Legacy inputs: required: `owner`, `repository`, `reviewers`
- Legacy endpoints: `PUT /repositories/{workspace}/{repo_slug}/default-reviewers/{target}/{username}`<br>`GET /repositories/{workspace}/{repo_slug}/default-reviewers`<br>`DELETE /repositories/{workspace}/{repo_slug}/default-reviewers/{target}/{username}`
- New inputs: required: `repo_slug`, `target_username`, `workspace`
- New operations: `Create PUT /repositories/{workspace}/{repo_slug}/default-reviewers/{target_username}`<br>`Read GET /repositories/{workspace}/{repo_slug}/default-reviewers/{target_username}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/default-reviewers/{target_username}`<br>`List GET /repositories/{workspace}/{repo_slug}/default-reviewers`
- Diff summary: renamed: `owner` → `workspace`, `repository` → `repo_slug`; legacy-only inputs: `reviewers`; new-only inputs: `target_username`

### `bitbucket_deploy_key`

- New equivalent(s): `bitbucket_repo_deploy_keys`
- Legacy inputs: required: `workspace`, `repository`, `key`; optional: `label`
- Legacy endpoints: `GET /repositories/{workspace}/{repo_slug}/deploy-keys/{key}/{id}`<br>`DELETE /repositories/{workspace}/{repo_slug}/deploy-keys/{key}/{id}`
- New inputs: required: `repo_slug`, `workspace`; optional: `key_id`
- New operations: `Create POST /repositories/{workspace}/{repo_slug}/deploy-keys`<br>`Read GET /repositories/{workspace}/{repo_slug}/deploy-keys/{key_id}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/deploy-keys/{key_id}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/deploy-keys/{key_id}`<br>`List GET /repositories/{workspace}/{repo_slug}/deploy-keys`
- Diff summary: renamed: `repository` → `repo_slug`; legacy-only inputs: `key`, `label`; new-only inputs: `key_id`
- Notes: The generated provider exposes deploy keys as `bitbucket_repo_deploy_keys` and also has separate project-level deploy key resources.

### `bitbucket_deployment`

- New equivalent(s): `bitbucket_deployments`
- Legacy inputs: required: `name`, `stage`, `repository`; optional: `restrictions`
- Legacy endpoints: `Create POST /repositories/{workspace}/{repo_slug}/environments`<br>`Read GET /repositories/{workspace}/{repo_slug}/environments/{environment_uuid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/environments/{environment_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/environments`
- New inputs: required: `workspace`, `repo_slug`; optional: `environment_uuid`, `name`, `uuid`, `request_body`
- New operations: `Create POST /repositories/{workspace}/{repo_slug}/environments`<br>`Read GET /repositories/{workspace}/{repo_slug}/environments/{environment_uuid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/environments/{environment_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/environments`
- Diff summary: renamed: `repository` → `repo_slug`; legacy-only inputs: `restrictions`, `stage`; new-only inputs: `environment_uuid`, `request_body`, `uuid`, `workspace`

### `bitbucket_deployment_variable`

- New equivalent(s): `bitbucket_deployment_variables`
- Legacy inputs: required: `deployment`, `key`, `value`; optional: `secured`
- Legacy endpoints: `Create POST /repositories/{workspace}/{repo_slug}/deployments_config/environments/{environment_uuid}/variables`<br>`Read GET /repositories/{workspace}/{repo_slug}/deployments_config/environments/{environment_uuid}/variables`<br>`Update PUT /repositories/{workspace}/{repo_slug}/deployments_config/environments/{environment_uuid}/variables/{variable_uuid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/deployments_config/environments/{environment_uuid}/variables/{variable_uuid}`
- New inputs: required: `workspace`, `repo_slug`, `environment_uuid`; optional: `variable_uuid`, `key`, `secured`, `uuid`, `value`, `request_body`
- New operations: `Create POST /repositories/{workspace}/{repo_slug}/deployments_config/environments/{environment_uuid}/variables`<br>`Read GET /repositories/{workspace}/{repo_slug}/deployments_config/environments/{environment_uuid}/variables`<br>`Update PUT /repositories/{workspace}/{repo_slug}/deployments_config/environments/{environment_uuid}/variables/{variable_uuid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/deployments_config/environments/{environment_uuid}/variables/{variable_uuid}`
- Diff summary: legacy-only inputs: `deployment`; new-only inputs: `environment_uuid`, `repo_slug`, `request_body`, `uuid`, `variable_uuid`, `workspace`

### `bitbucket_forked_repository`

- New equivalent(s): `bitbucket_forked_repository`
- Legacy inputs: required: `owner`, `name`; optional: `slug`, `is_private`, `website`, `language`, `has_issues`, `has_wiki`, `project_key`, `fork_policy`, `description`, `pipelines_enabled`, `link`
- Legacy endpoints: `POST /repositories/{workspace}/{repo_slug}/forks`<br>`GET /repositories/{workspace}/{repo_slug}`
- New inputs: required: `repo_slug`, `workspace`; optional: `description`, `fork_policy`, `full_name`, `has_issues`, `has_wiki`, `is_private`, `language`, `mainbranch`, `name`, `owner`, `project`, `scm`, `size`, `uuid`, `request_body`
- New operations: `Create POST /repositories/{workspace}/{repo_slug}/forks`<br>`List GET /repositories/{workspace}/{repo_slug}/forks`
- Diff summary: renamed: `owner` → `workspace`; legacy-only inputs: `link`, `pipelines_enabled`, `project_key`, `slug`, `website`; new-only inputs: `full_name`, `mainbranch`, `owner`, `project`, `repo_slug`, `request_body`, `scm`, `size`, `uuid`

### `bitbucket_hook`

- New equivalent(s): `bitbucket_hooks`
- Legacy inputs: required: `owner`, `repository`, `url`, `description`, `events`; optional: `active`, `skip_cert_verification`, `secret`
- Legacy endpoints: `Create POST /repositories/{workspace}/{repo_slug}/hooks`<br>`Read GET /repositories/{workspace}/{repo_slug}/hooks/{uid}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/hooks/{uid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/hooks/{uid}`<br>`List GET /repositories/{workspace}/{repo_slug}/hooks`
- New inputs: required: `repo_slug`, `workspace`; optional: `uid`
- New operations: `Create POST /repositories/{workspace}/{repo_slug}/hooks`<br>`Read GET /repositories/{workspace}/{repo_slug}/hooks/{uid}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/hooks/{uid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/hooks/{uid}`<br>`List GET /repositories/{workspace}/{repo_slug}/hooks`
- Diff summary: renamed: `owner` → `workspace`, `repository` → `repo_slug`; legacy-only inputs: `active`, `description`, `events`, `secret`, `skip_cert_verification`, `url`; new-only inputs: `uid`

### `bitbucket_pipeline_schedule`

- New equivalent(s): `bitbucket_pipeline_schedules`
- Legacy inputs: required: `workspace`, `repository`, `enabled`, `cron_pattern`, `target`
- Legacy endpoints: `Create POST /repositories/{workspace}/{repo_slug}/pipelines_config/schedules`<br>`Read GET /repositories/{workspace}/{repo_slug}/pipelines_config/schedules/{schedule_uuid}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/pipelines_config/schedules/{schedule_uuid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/pipelines_config/schedules/{schedule_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/pipelines_config/schedules`
- New inputs: required: `workspace`, `repo_slug`; optional: `schedule_uuid`, `cron_pattern`, `enabled`, `target`, `request_body`
- New operations: `Create POST /repositories/{workspace}/{repo_slug}/pipelines_config/schedules`<br>`Read GET /repositories/{workspace}/{repo_slug}/pipelines_config/schedules/{schedule_uuid}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/pipelines_config/schedules/{schedule_uuid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/pipelines_config/schedules/{schedule_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/pipelines_config/schedules`
- Diff summary: renamed: `repository` → `repo_slug`; new-only inputs: `request_body`, `schedule_uuid`

### `bitbucket_pipeline_ssh_key`

- New equivalent(s): `bitbucket_pipeline_ssh_keys`
- Legacy inputs: required: `workspace`, `repository`, `public_key`, `private_key`
- Legacy endpoints: `Read GET /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/key_pair`<br>`Update PUT /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/key_pair`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/key_pair`
- New inputs: required: `workspace`, `repo_slug`; optional: `private_key`, `public_key`, `request_body`
- New operations: `Read GET /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/key_pair`<br>`Update PUT /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/key_pair`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/key_pair`
- Diff summary: renamed: `repository` → `repo_slug`; new-only inputs: `request_body`

### `bitbucket_pipeline_ssh_known_host`

- New equivalent(s): `bitbucket_pipeline_known_hosts`
- Legacy inputs: required: `workspace`, `repository`, `hostname`, `public_key`
- Legacy endpoints: `Create POST /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts`<br>`Read GET /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts/{known_host_uuid}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts/{known_host_uuid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts/{known_host_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts`
- New inputs: required: `workspace`, `repo_slug`; optional: `known_host_uuid`, `hostname`, `public_key`, `uuid`, `request_body`
- New operations: `Create POST /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts`<br>`Read GET /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts/{known_host_uuid}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts/{known_host_uuid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts/{known_host_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts`
- Diff summary: renamed: `repository` → `repo_slug`; new-only inputs: `known_host_uuid`, `request_body`, `uuid`

### `bitbucket_project`

- New equivalent(s): `bitbucket_projects`
- Legacy inputs: required: `owner`, `name`, `key`; optional: `description`, `is_private`, `link`
- Legacy endpoints: `PUT /workspaces/{workspace}/projects/{project_key}`<br>`POST /workspaces/{workspace}/projects`<br>`GET /workspaces/{workspace}/projects/{project_key}`<br>`DELETE /workspaces/{workspace}/projects/{project_key}`
- New inputs: required: `workspace`; optional: `project_key`, `request_body`
- New operations: `Create POST /workspaces/{workspace}/projects`<br>`Read GET /workspaces/{workspace}/projects/{project_key}`<br>`Update PUT /workspaces/{workspace}/projects/{project_key}`<br>`Delete DELETE /workspaces/{workspace}/projects/{project_key}`<br>`List GET /workspaces/{workspace}/projects`
- Diff summary: renamed: `owner` → `workspace`; legacy-only inputs: `description`, `is_private`, `key`, `link`, `name`; new-only inputs: `project_key`, `request_body`

### `bitbucket_project_branching_model`

- New equivalent(s): `bitbucket_project_branching_model`
- Legacy inputs: required: `workspace`, `project`, `branch_type`; optional: `development`, `production`
- Legacy endpoints: `Read GET /workspaces/{workspace}/projects/{project_key}/branching-model`<br>`Update PUT /workspaces/{workspace}/projects/{project_key}/branching-model/settings`
- New inputs: required: `project_key`, `workspace`
- New operations: `Read GET /workspaces/{workspace}/projects/{project_key}/branching-model`<br>`Update PUT /workspaces/{workspace}/projects/{project_key}/branching-model/settings`
- Diff summary: legacy-only inputs: `branch_type`, `development`, `production`, `project`; new-only inputs: `project_key`

### `bitbucket_project_default_reviewers`

- New equivalent(s): `bitbucket_project_default_reviewers`
- Legacy inputs: required: `workspace`, `project`, `reviewers`
- Legacy endpoints: `PUT /workspaces/{workspace}/projects/{project_key}/default-reviewers/{selected_user}`<br>`GET /workspaces/{workspace}/projects/{project_key}/default-reviewers`<br>`DELETE /workspaces/{workspace}/projects/{project_key}/default-reviewers/{selected_user}`
- New inputs: required: `project_key`, `selected_user`, `workspace`
- New operations: `Create PUT /workspaces/{workspace}/projects/{project_key}/default-reviewers/{selected_user}`<br>`Read GET /workspaces/{workspace}/projects/{project_key}/default-reviewers/{selected_user}`<br>`Delete DELETE /workspaces/{workspace}/projects/{project_key}/default-reviewers/{selected_user}`<br>`List GET /workspaces/{workspace}/projects/{project_key}/default-reviewers`
- Diff summary: legacy-only inputs: `project`, `reviewers`; new-only inputs: `project_key`, `selected_user`

### `bitbucket_project_group_permission`

- New equivalent(s): `bitbucket_project_group_permissions`
- Legacy inputs: required: `workspace`, `project_key`, `group_slug`, `permission`
- Legacy endpoints: `Read GET /workspaces/{workspace}/projects/{project_key}/permissions-config/groups/{group_slug}`<br>`Update PUT /workspaces/{workspace}/projects/{project_key}/permissions-config/groups/{group_slug}`<br>`Delete DELETE /workspaces/{workspace}/projects/{project_key}/permissions-config/groups/{group_slug}`<br>`List GET /workspaces/{workspace}/projects/{project_key}/permissions-config/groups`
- New inputs: required: `group_slug`, `project_key`, `workspace`; optional: `request_body`
- New operations: `Read GET /workspaces/{workspace}/projects/{project_key}/permissions-config/groups/{group_slug}`<br>`Update PUT /workspaces/{workspace}/projects/{project_key}/permissions-config/groups/{group_slug}`<br>`Delete DELETE /workspaces/{workspace}/projects/{project_key}/permissions-config/groups/{group_slug}`<br>`List GET /workspaces/{workspace}/projects/{project_key}/permissions-config/groups`
- Diff summary: legacy-only inputs: `permission`; new-only inputs: `request_body`

### `bitbucket_project_user_permission`

- New equivalent(s): `bitbucket_project_user_permissions`
- Legacy inputs: required: `workspace`, `project_key`, `user_id`, `permission`
- Legacy endpoints: `Read GET /workspaces/{workspace}/projects/{project_key}/permissions-config/users/{selected_user_id}`<br>`Update PUT /workspaces/{workspace}/projects/{project_key}/permissions-config/users/{selected_user_id}`<br>`Delete DELETE /workspaces/{workspace}/projects/{project_key}/permissions-config/users/{selected_user_id}`<br>`List GET /workspaces/{workspace}/projects/{project_key}/permissions-config/users`
- New inputs: required: `project_key`, `selected_user_id`, `workspace`; optional: `request_body`
- New operations: `Read GET /workspaces/{workspace}/projects/{project_key}/permissions-config/users/{selected_user_id}`<br>`Update PUT /workspaces/{workspace}/projects/{project_key}/permissions-config/users/{selected_user_id}`<br>`Delete DELETE /workspaces/{workspace}/projects/{project_key}/permissions-config/users/{selected_user_id}`<br>`List GET /workspaces/{workspace}/projects/{project_key}/permissions-config/users`
- Diff summary: legacy-only inputs: `permission`, `user_id`; new-only inputs: `request_body`, `selected_user_id`

### `bitbucket_repository`

- New equivalent(s): `bitbucket_repos`, `bitbucket_repo_settings`, `bitbucket_pipeline_config`
- Legacy inputs: required: `owner`, `name`; optional: `slug`, `scm`, `is_private`, `website`, `language`, `has_issues`, `has_wiki`, `project_key`, `fork_policy`, `description`, `pipelines_enabled`, `link`, `inherit_default_merge_strategy`, `inherit_branching_model`
- Legacy endpoints: `PUT /repositories/{workspace}/{repo_slug}`<br>`POST /repositories/{workspace}/{repo_slug}`<br>`GET /repositories/{workspace}/{repo_slug}`<br>`DELETE /repositories/{workspace}/{repo_slug}`
- New inputs: required: `repo_slug`, `workspace`; optional: `description`, `fork_policy`, `full_name`, `has_issues`, `has_wiki`, `is_private`, `language`, `mainbranch`, `name`, `owner`, `project`, `scm`, `size`, `uuid`, `request_body`, `enabled`, `repository`
- New operations: `Create POST /repositories/{workspace}/{repo_slug}`<br>`Read GET /repositories/{workspace}/{repo_slug}`<br>`Update PUT /repositories/{workspace}/{repo_slug}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}`<br>`List GET /repositories/{workspace}`<br>`Read GET /repositories/{workspace}/{repo_slug}/override-settings`<br>`Update PUT /repositories/{workspace}/{repo_slug}/override-settings`<br>`Read GET /repositories/{workspace}/{repo_slug}/pipelines_config`<br>`Update PUT /repositories/{workspace}/{repo_slug}/pipelines_config`
- Diff summary: renamed: `owner` → `workspace`; legacy-only inputs: `inherit_branching_model`, `inherit_default_merge_strategy`, `link`, `pipelines_enabled`, `project_key`, `slug`, `website`; new-only inputs: `enabled`, `full_name`, `mainbranch`, `owner`, `project`, `repo_slug`, `repository`, `request_body`, `size`, `uuid`
- Notes: The legacy repository resource bundled core repository CRUD, pipeline enablement, and override-settings flags. In the new provider, core CRUD stays on `bitbucket_repos`, pipeline enablement moves to `bitbucket_pipeline_config`, and repository settings have their own `bitbucket_repo_settings` resource.

### `bitbucket_repository_group_permission`

- New equivalent(s): `bitbucket_repo_group_permissions`
- Legacy inputs: required: `workspace`, `repo_slug`, `group_slug`, `permission`
- Legacy endpoints: `Read GET /repositories/{workspace}/{repo_slug}/permissions-config/groups/{group_slug}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/permissions-config/groups/{group_slug}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/permissions-config/groups/{group_slug}`<br>`List GET /repositories/{workspace}/{repo_slug}/permissions-config/groups`
- New inputs: required: `group_slug`, `repo_slug`, `workspace`; optional: `request_body`
- New operations: `Read GET /repositories/{workspace}/{repo_slug}/permissions-config/groups/{group_slug}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/permissions-config/groups/{group_slug}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/permissions-config/groups/{group_slug}`<br>`List GET /repositories/{workspace}/{repo_slug}/permissions-config/groups`
- Diff summary: legacy-only inputs: `permission`; new-only inputs: `request_body`

### `bitbucket_repository_user_permission`

- New equivalent(s): `bitbucket_repo_user_permissions`
- Legacy inputs: required: `workspace`, `repo_slug`, `user_id`, `permission`
- Legacy endpoints: `Read GET /repositories/{workspace}/{repo_slug}/permissions-config/users/{selected_user_id}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/permissions-config/users/{selected_user_id}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/permissions-config/users/{selected_user_id}`<br>`List GET /repositories/{workspace}/{repo_slug}/permissions-config/users`
- New inputs: required: `repo_slug`, `selected_user_id`, `workspace`; optional: `request_body`
- New operations: `Read GET /repositories/{workspace}/{repo_slug}/permissions-config/users/{selected_user_id}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/permissions-config/users/{selected_user_id}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/permissions-config/users/{selected_user_id}`<br>`List GET /repositories/{workspace}/{repo_slug}/permissions-config/users`
- Diff summary: legacy-only inputs: `permission`, `user_id`; new-only inputs: `request_body`, `selected_user_id`

### `bitbucket_repository_variable`

- New equivalent(s): `bitbucket_pipeline_variables`
- Legacy inputs: required: `key`, `value`, `repository`; optional: `secured`
- Legacy endpoints: `Create POST /repositories/{workspace}/{repo_slug}/pipelines_config/variables`<br>`Read GET /repositories/{workspace}/{repo_slug}/pipelines_config/variables/{variable_uuid}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/pipelines_config/variables/{variable_uuid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/pipelines_config/variables/{variable_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/pipelines_config/variables`
- New inputs: required: `workspace`, `repo_slug`; optional: `variable_uuid`, `key`, `secured`, `uuid`, `value`, `request_body`
- New operations: `Create POST /repositories/{workspace}/{repo_slug}/pipelines_config/variables`<br>`Read GET /repositories/{workspace}/{repo_slug}/pipelines_config/variables/{variable_uuid}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/pipelines_config/variables/{variable_uuid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/pipelines_config/variables/{variable_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/pipelines_config/variables`
- Diff summary: renamed: `repository` → `repo_slug`; new-only inputs: `request_body`, `uuid`, `variable_uuid`, `workspace`
- Notes: Legacy repository variables map to the pipelines variable API. Use `bitbucket_pipeline_variables` and rename `owner`/`repository` to `workspace`/`repo_slug`.

### `bitbucket_ssh_key`

- New equivalent(s): `bitbucket_ssh_keys`
- Legacy inputs: required: `user`, `key`; optional: `label`
- Legacy endpoints: `POST /users/{selected_user}/ssh-keys`<br>`GET /users/{selected_user}/ssh-keys/{key}/{id}`<br>`PUT /users/{selected_user}/ssh-keys/{key}/{id}`<br>`DELETE /users/{selected_user}/ssh-keys/{key}/{id}`
- New inputs: required: `selected_user`; optional: `key_id`, `comment`, `expires_on`, `fingerprint`, `key`, `label`, `last_used`, `owner`, `uuid`, `request_body`
- New operations: `Create POST /users/{selected_user}/ssh-keys`<br>`Read GET /users/{selected_user}/ssh-keys/{key_id}`<br>`Update PUT /users/{selected_user}/ssh-keys/{key_id}`<br>`Delete DELETE /users/{selected_user}/ssh-keys/{key_id}`<br>`List GET /users/{selected_user}/ssh-keys`
- Diff summary: legacy-only inputs: `user`; new-only inputs: `comment`, `expires_on`, `fingerprint`, `key_id`, `last_used`, `owner`, `request_body`, `selected_user`, `uuid`

### `bitbucket_workspace_hook`

- New equivalent(s): `bitbucket_workspace_hooks`
- Legacy inputs: required: `workspace`, `url`, `description`, `events`; optional: `active`, `skip_cert_verification`, `secret`
- Legacy endpoints: `Create POST /workspaces/{workspace}/hooks`<br>`Read GET /workspaces/{workspace}/hooks/{uid}`<br>`Update PUT /workspaces/{workspace}/hooks/{uid}`<br>`Delete DELETE /workspaces/{workspace}/hooks/{uid}`<br>`List GET /workspaces/{workspace}/hooks`
- New inputs: required: `workspace`; optional: `uid`
- New operations: `Create POST /workspaces/{workspace}/hooks`<br>`Read GET /workspaces/{workspace}/hooks/{uid}`<br>`Update PUT /workspaces/{workspace}/hooks/{uid}`<br>`Delete DELETE /workspaces/{workspace}/hooks/{uid}`<br>`List GET /workspaces/{workspace}/hooks`
- Diff summary: legacy-only inputs: `active`, `description`, `events`, `secret`, `skip_cert_verification`, `url`; new-only inputs: `uid`

### `bitbucket_workspace_variable`

- New equivalent(s): `bitbucket_workspace_pipeline_variables`
- Legacy inputs: required: `workspace`, `key`, `value`; optional: `secured`
- Legacy endpoints: `Create POST /workspaces/{workspace}/pipelines-config/variables`<br>`Read GET /workspaces/{workspace}/pipelines-config/variables/{variable_uuid}`<br>`Update PUT /workspaces/{workspace}/pipelines-config/variables/{variable_uuid}`<br>`Delete DELETE /workspaces/{workspace}/pipelines-config/variables/{variable_uuid}`<br>`List GET /workspaces/{workspace}/pipelines-config/variables`
- New inputs: required: `workspace`; optional: `variable_uuid`, `request_body`
- New operations: `Create POST /workspaces/{workspace}/pipelines-config/variables`<br>`Read GET /workspaces/{workspace}/pipelines-config/variables/{variable_uuid}`<br>`Update PUT /workspaces/{workspace}/pipelines-config/variables/{variable_uuid}`<br>`Delete DELETE /workspaces/{workspace}/pipelines-config/variables/{variable_uuid}`<br>`List GET /workspaces/{workspace}/pipelines-config/variables`
- Diff summary: legacy-only inputs: `key`, `secured`, `value`; new-only inputs: `request_body`, `variable_uuid`
- Notes: Workspace variables now live under the pipelines API as `bitbucket_workspace_pipeline_variables`.

## Legacy-only resources

### `bitbucket_group`

- New equivalent(s): none
- Legacy inputs: required: `workspace`, `name`; optional: `auto_add`, `permission`
- Legacy endpoints: none
- Notes: Workspace group management is not currently exposed by the generated provider because those endpoints are not represented in the generated Terraform docs.

### `bitbucket_group_membership`

- New equivalent(s): none
- Legacy inputs: required: `workspace`, `group_slug`, `uuid`
- Legacy endpoints: none
- Notes: Group membership management is not currently exposed by the generated provider.

## Matched legacy data sources

### `bitbucket_current_user`

- New equivalent(s): `bitbucket_current_user`
- Legacy inputs: none
- Legacy endpoints: `GET /user`<br>`GET /user/emails`
- New inputs: none
- New operations: `Read GET /user`
- Diff summary: input names are effectively unchanged
- Notes: The legacy data source also fetched `/user/emails`. The generated provider splits that into `bitbucket_current_user` plus `bitbucket_user_emails` when you need email addresses.

### `bitbucket_deployment`

- New equivalent(s): `bitbucket_deployments`
- Legacy inputs: required: `uuid`, `repository`, `workspace`
- Legacy endpoints: `Read GET /repositories/{workspace}/{repo_slug}/environments/{environment_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/environments`
- New inputs: required: `workspace`, `repo_slug`; optional: `environment_uuid`
- New operations: `Read GET /repositories/{workspace}/{repo_slug}/environments/{environment_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/environments`
- Diff summary: renamed: `repository` → `repo_slug`; legacy-only inputs: `uuid`; new-only inputs: `environment_uuid`
- Notes: Use `bitbucket_deployments` with the identifying path parameters for a single deployment; omit the single-resource expectation and treat the response as the generic deployment payload.

### `bitbucket_deployments`

- New equivalent(s): `bitbucket_deployments`
- Legacy inputs: required: `repository`, `workspace`
- Legacy endpoints: `Read GET /repositories/{workspace}/{repo_slug}/environments/{environment_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/environments`
- New inputs: required: `workspace`, `repo_slug`; optional: `environment_uuid`
- New operations: `Read GET /repositories/{workspace}/{repo_slug}/environments/{environment_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/environments`
- Diff summary: renamed: `repository` → `repo_slug`; new-only inputs: `environment_uuid`

### `bitbucket_file`

- New equivalent(s): `bitbucket_commit_file`
- Legacy inputs: none
- Legacy endpoints: `Read GET /repositories/{workspace}/{repo_slug}/src/{commit}/{path}`
- New inputs: required: `commit`, `path`, `repo_slug`, `workspace`
- New operations: `Read GET /repositories/{workspace}/{repo_slug}/src/{commit}/{path}`
- Diff summary: new-only inputs: `commit`, `path`, `repo_slug`, `workspace`
- Notes: The legacy `bitbucket_file` data source maps most closely to `bitbucket_commit_file`, which reads file content via the commit-file endpoint.

### `bitbucket_hook_types`

- New equivalent(s): `bitbucket_hook_types`
- Legacy inputs: none
- Legacy endpoints: `GET /hook-events-subject-type`
- New inputs: none
- New operations: `Read GET /hook_events`<br>`List GET /hook_events/{subject_type}`
- Diff summary: input names are effectively unchanged

### `bitbucket_pipeline_oidc_config`

- New equivalent(s): `bitbucket_pipeline_oidc`
- Legacy inputs: required: `workspace`
- Legacy endpoints: `Read GET /workspaces/{workspace}/pipelines-config/identity/oidc/.well-known/openid-configuration`
- New inputs: required: `workspace`
- New operations: `Read GET /workspaces/{workspace}/pipelines-config/identity/oidc/.well-known/openid-configuration`
- Diff summary: input names are effectively unchanged

### `bitbucket_pipeline_oidc_config_keys`

- New equivalent(s): `bitbucket_pipeline_oidc_keys`
- Legacy inputs: required: `workspace`
- Legacy endpoints: `Read GET /workspaces/{workspace}/pipelines-config/identity/oidc/keys.json`
- New inputs: required: `workspace`
- New operations: `Read GET /workspaces/{workspace}/pipelines-config/identity/oidc/keys.json`
- Diff summary: input names are effectively unchanged

### `bitbucket_project`

- New equivalent(s): `bitbucket_projects`
- Legacy inputs: none
- Legacy endpoints: `Read GET /workspaces/{workspace}/projects/{project_key}`<br>`List GET /workspaces/{workspace}/projects`
- New inputs: required: `workspace`; optional: `project_key`
- New operations: `Read GET /workspaces/{workspace}/projects/{project_key}`<br>`List GET /workspaces/{workspace}/projects`
- Diff summary: new-only inputs: `project_key`, `workspace`

### `bitbucket_repository`

- New equivalent(s): `bitbucket_repos`
- Legacy inputs: none
- Legacy endpoints: `Read GET /repositories/{workspace}/{repo_slug}`<br>`List GET /repositories/{workspace}`
- New inputs: required: `workspace`; optional: `repo_slug`
- New operations: `Read GET /repositories/{workspace}/{repo_slug}`<br>`List GET /repositories/{workspace}`
- Diff summary: new-only inputs: `repo_slug`, `workspace`

### `bitbucket_user`

- New equivalent(s): `bitbucket_users`
- Legacy inputs: optional: `uuid`
- Legacy endpoints: `GET /users/{selected_user}`
- New inputs: required: `selected_user`
- New operations: `Read GET /users/{selected_user}`<br>`List GET /users/{selected_user}/ssh-keys`
- Diff summary: legacy-only inputs: `uuid`; new-only inputs: `selected_user`

### `bitbucket_workspace`

- New equivalent(s): `bitbucket_workspaces`
- Legacy inputs: required: `workspace`
- Legacy endpoints: `GET /workspaces/{workspace}`
- New inputs: optional: `workspace`
- New operations: `Read GET /workspaces/{workspace}`<br>`List GET /workspaces`
- Diff summary: input names are effectively unchanged

### `bitbucket_workspace_members`

- New equivalent(s): `bitbucket_workspace_members`
- Legacy inputs: required: `workspace`
- Legacy endpoints: `GET /workspaces/{workspace}/members`
- New inputs: required: `workspace`; optional: `member`
- New operations: `Read GET /workspaces/{workspace}/members/{member}`<br>`List GET /workspaces/{workspace}/members`
- Diff summary: new-only inputs: `member`

## Legacy-only data sources

### `bitbucket_group`

- New equivalent(s): none
- Legacy inputs: required: `workspace`, `slug`
- Legacy endpoints: none
- Notes: Group lookup is not currently exposed by the generated provider.

### `bitbucket_group_members`

- New equivalent(s): none
- Legacy inputs: required: `workspace`, `slug`
- Legacy endpoints: none
- Notes: Group member lookup is not currently exposed by the generated provider.

### `bitbucket_groups`

- New equivalent(s): none
- Legacy inputs: required: `workspace`
- Legacy endpoints: none
- Notes: Group listing is not currently exposed by the generated provider.

### `bitbucket_ip_ranges`

- New equivalent(s): none
- Legacy inputs: none
- Legacy endpoints: none
- Notes: The generated provider does not currently expose Bitbucket IP ranges as a Terraform data source.

## New provider-only resources

- `bitbucket_addon`
- `bitbucket_annotations`
- `bitbucket_commit_statuses`
- `bitbucket_commits`
- `bitbucket_current_user`
- `bitbucket_downloads`
- `bitbucket_gpg_keys`
- `bitbucket_hook_types`
- `bitbucket_issue_comments`
- `bitbucket_issues`
- `bitbucket_pipeline_caches`
- `bitbucket_pipeline_oidc`
- `bitbucket_pipeline_oidc_keys`
- `bitbucket_pipelines`
- `bitbucket_pr`
- `bitbucket_pr_comments`
- `bitbucket_project_deploy_keys`
- `bitbucket_properties`
- `bitbucket_refs`
- `bitbucket_repo_runners`
- `bitbucket_reports`
- `bitbucket_search`
- `bitbucket_snippets`
- `bitbucket_tags`
- `bitbucket_user_emails`
- `bitbucket_users`
- `bitbucket_workspace_members`
- `bitbucket_workspace_permissions`
- `bitbucket_workspace_runners`
- `bitbucket_workspaces`

## New provider-only data sources

- `bitbucket_addon`
- `bitbucket_annotations`
- `bitbucket_branch_restrictions`
- `bitbucket_branching_model`
- `bitbucket_commit_statuses`
- `bitbucket_commits`
- `bitbucket_default_reviewers`
- `bitbucket_deployment_variables`
- `bitbucket_downloads`
- `bitbucket_forked_repository`
- `bitbucket_gpg_keys`
- `bitbucket_hooks`
- `bitbucket_issue_comments`
- `bitbucket_issues`
- `bitbucket_pipeline_caches`
- `bitbucket_pipeline_config`
- `bitbucket_pipeline_known_hosts`
- `bitbucket_pipeline_schedules`
- `bitbucket_pipeline_ssh_keys`
- `bitbucket_pipeline_variables`
- `bitbucket_pipelines`
- `bitbucket_pr`
- `bitbucket_pr_comments`
- `bitbucket_project_branching_model`
- `bitbucket_project_default_reviewers`
- `bitbucket_project_deploy_keys`
- `bitbucket_project_group_permissions`
- `bitbucket_project_user_permissions`
- `bitbucket_properties`
- `bitbucket_refs`
- `bitbucket_repo_deploy_keys`
- `bitbucket_repo_group_permissions`
- `bitbucket_repo_runners`
- `bitbucket_repo_settings`
- `bitbucket_repo_user_permissions`
- `bitbucket_reports`
- `bitbucket_search`
- `bitbucket_snippets`
- `bitbucket_ssh_keys`
- `bitbucket_tags`
- `bitbucket_user_emails`
- `bitbucket_workspace_hooks`
- `bitbucket_workspace_permissions`
- `bitbucket_workspace_pipeline_variables`
- `bitbucket_workspace_runners`

## Can this be automated?

Partly. A comparison script is practical today, but a fully automatic HCL rewrite is only safe for the straightforward cases.

Good candidates for an automated rewrite later:

- provider source replacement
- provider auth field rename (`password` → `token`)
- direct resource/data source renames where there is a 1:1 mapping
- path argument renames like `owner` → `workspace` and `repository` → `repo_slug`

Cases that still need manual review:

- legacy objects that split into multiple generated resources
- objects missing from one provider or the other
- fields whose semantics changed even when the name looks similar
- places where the generated provider expects `request_body` for uncommon fields

