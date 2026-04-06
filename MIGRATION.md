# Migration from `DrFaust92/terraform-provider-bitbucket`

This guide compares the legacy hand-written provider with the generated `FabianSchurig/bitbucket` provider in this repository.

It is intentionally a best-effort migration baseline: the generated docs sometimes list optional fields that are also computed by the API, so subtle cases still need manual review.

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
- Legacy endpoints: `POST /repositories/{workspace}/{repo_slug}/branch-restrictions`<br>`GET /repositories/{workspace}/{repo_slug}/branch-restrictions/{id}`<br>`PUT /repositories/{workspace}/{repo_slug}/branch-restrictions/{id}`<br>`DELETE /repositories/{workspace}/{repo_slug}/branch-restrictions/{id}`

#### Legacy HCL

```hcl
resource "bitbucket_branch_restriction" "legacy" {
  kind = "push"
  owner = "my-workspace"
  repository = "my-repo"

  # branch_match_kind = "glob"  # optional
  # branch_type = "feature"  # optional
  # groups = "example-groups"  # optional
  # pattern = "main"  # optional
  # users = "example-users"  # optional
  # value = "example-value"  # optional
}
```

- New operations: `Create POST /repositories/{workspace}/{repo_slug}/branch-restrictions`<br>`Read GET /repositories/{workspace}/{repo_slug}/branch-restrictions/{id}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/branch-restrictions/{id}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/branch-restrictions/{id}`<br>`List GET /repositories/{workspace}/{repo_slug}/branch-restrictions`

#### New HCL

##### `bitbucket_branch_restrictions`

```hcl
resource "bitbucket_branch_restrictions" "migrated" {
  kind = "push"
  repo_slug = "my-repo"
  workspace = "my-workspace"

  # branch_match_kind = "glob"  # optional
  # branch_type = "feature"  # optional
  # groups = "example-groups"  # optional
  # param_id = "example-param-id"  # new-only
  # pattern = "main"  # optional
  # request_body = jsonencode({})  # new-only
  # users = "example-users"  # optional
  # value = "example-value"  # optional
}
```

- Diff summary: renamed: `owner` → `workspace`, `repository` → `repo_slug`; new-only inputs: `param_id`, `request_body`

### `bitbucket_branching_model`

- New equivalent(s): `bitbucket_branching_model`
- Legacy endpoints: `Read GET /repositories/{workspace}/{repo_slug}/branching-model`<br>`Update PUT /repositories/{workspace}/{repo_slug}/branching-model/settings`

#### Legacy HCL

```hcl
resource "bitbucket_branching_model" "legacy" {
  branch_type = "feature"
  owner = "my-workspace"
  repository = "my-repo"

  # development = "example-development"  # optional
  # production = "example-production"  # optional
}
```

- New operations: `Read GET /repositories/{workspace}/{repo_slug}/branching-model`<br>`Update PUT /repositories/{workspace}/{repo_slug}/branching-model/settings`

#### New HCL

##### `bitbucket_branching_model`

```hcl
resource "bitbucket_branching_model" "migrated" {
  repo_slug = "my-repo"
  workspace = "my-workspace"

  # branch_type = "feature"  # legacy-only
  # development = "example-development"  # legacy-only
  # production = "example-production"  # legacy-only
}
```

- Diff summary: renamed: `owner` → `workspace`, `repository` → `repo_slug`; legacy-only inputs: `branch_type`, `development`, `production`

### `bitbucket_commit_file`

- New equivalent(s): `bitbucket_commit_file`
- Legacy endpoints: `GET /repositories/{workspace}/{repo_slug}/src/{commit}/{path}`

#### Legacy HCL

```hcl
resource "bitbucket_commit_file" "legacy" {
  branch = "main"
  commit_author = "Jane Doe <jane@example.com>"
  commit_message = "Example commit"
  content = "example content"
  filename = "example-filename"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

- New operations: `Create POST /repositories/{workspace}/{repo_slug}/src`<br>`Read GET /repositories/{workspace}/{repo_slug}/src/{commit}/{path}`

#### New HCL

##### `bitbucket_commit_file`

```hcl
resource "bitbucket_commit_file" "migrated" {
  repo_slug = "my-repo"
  workspace = "my-workspace"

  # branch = "main"  # legacy-only
  # commit = "main"  # new-only
  # commit_author = "Jane Doe <jane@example.com>"  # legacy-only
  # commit_message = "Example commit"  # legacy-only
  # content = "example content"  # legacy-only
  # filename = "example-filename"  # legacy-only
  # path = "README.md"  # new-only
}
```

- Diff summary: legacy-only inputs: `branch`, `commit_author`, `commit_message`, `content`, `filename`; new-only inputs: `commit`, `path`

### `bitbucket_default_reviewers`

- New equivalent(s): `bitbucket_default_reviewers`
- Legacy endpoints: `PUT /repositories/{workspace}/{repo_slug}/default-reviewers/{target}/{username}`<br>`GET /repositories/{workspace}/{repo_slug}/default-reviewers`<br>`DELETE /repositories/{workspace}/{repo_slug}/default-reviewers/{target}/{username}`

#### Legacy HCL

```hcl
resource "bitbucket_default_reviewers" "legacy" {
  owner = "my-workspace"
  repository = "my-repo"
  reviewers = ["example-user"]
}
```

- New operations: `Create PUT /repositories/{workspace}/{repo_slug}/default-reviewers/{target_username}`<br>`Read GET /repositories/{workspace}/{repo_slug}/default-reviewers/{target_username}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/default-reviewers/{target_username}`<br>`List GET /repositories/{workspace}/{repo_slug}/default-reviewers`

#### New HCL

##### `bitbucket_default_reviewers`

```hcl
resource "bitbucket_default_reviewers" "migrated" {
  repo_slug = "my-repo"
  target_username = "example-user"
  workspace = "my-workspace"

  # reviewers = ["example-user"]  # legacy-only
}
```

- Diff summary: renamed: `owner` → `workspace`, `repository` → `repo_slug`; legacy-only inputs: `reviewers`; new-only inputs: `target_username`

### `bitbucket_deploy_key`

- New equivalent(s): `bitbucket_repo_deploy_keys`
- Legacy endpoints: `GET /repositories/{workspace}/{repo_slug}/deploy-keys/{key}/{id}`<br>`DELETE /repositories/{workspace}/{repo_slug}/deploy-keys/{key}/{id}`

#### Legacy HCL

```hcl
resource "bitbucket_deploy_key" "legacy" {
  key = "example-key"
  repository = "my-repo"
  workspace = "my-workspace"

  # label = "Example label"  # optional
}
```

- New operations: `Create POST /repositories/{workspace}/{repo_slug}/deploy-keys`<br>`Read GET /repositories/{workspace}/{repo_slug}/deploy-keys/{key_id}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/deploy-keys/{key_id}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/deploy-keys/{key_id}`<br>`List GET /repositories/{workspace}/{repo_slug}/deploy-keys`

#### New HCL

##### `bitbucket_repo_deploy_keys`

```hcl
resource "bitbucket_repo_deploy_keys" "migrated" {
  repo_slug = "my-repo"
  workspace = "my-workspace"

  # key = "example-key"  # legacy-only
  # key_id = "example-key-id"  # new-only
  # label = "Example label"  # legacy-only
}
```

- Diff summary: renamed: `repository` → `repo_slug`; legacy-only inputs: `key`, `label`; new-only inputs: `key_id`
- Notes: The generated provider exposes deploy keys as `bitbucket_repo_deploy_keys` and also has separate project-level deploy key resources.

### `bitbucket_deployment`

- New equivalent(s): `bitbucket_deployments`
- Legacy endpoints: `Create POST /repositories/{workspace}/{repo_slug}/environments`<br>`Read GET /repositories/{workspace}/{repo_slug}/environments/{environment_uuid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/environments/{environment_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/environments`

#### Legacy HCL

```hcl
resource "bitbucket_deployment" "legacy" {
  name = "my-repo"
  repository = "my-repo"
  stage = "Test"

  # restrictions = "example-restrictions"  # optional
}
```

- New operations: `Create POST /repositories/{workspace}/{repo_slug}/environments`<br>`Read GET /repositories/{workspace}/{repo_slug}/environments/{environment_uuid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/environments/{environment_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/environments`

#### New HCL

##### `bitbucket_deployments`

```hcl
resource "bitbucket_deployments" "migrated" {
  name = "my-repo"
  repo_slug = "my-repo"
  workspace = "my-workspace"

  # environment_uuid = "{environment-uuid}"  # new-only
  # request_body = jsonencode({})  # new-only
  # restrictions = "example-restrictions"  # legacy-only
  # stage = "Test"  # legacy-only
}
```

- Diff summary: renamed: `repository` → `repo_slug`; legacy-only inputs: `restrictions`, `stage`; new-only inputs: `environment_uuid`, `request_body`, `uuid`, `workspace`

### `bitbucket_deployment_variable`

- New equivalent(s): `bitbucket_deployment_variables`
- Legacy endpoints: `Create POST /repositories/{workspace}/{repo_slug}/deployments_config/environments/{environment_uuid}/variables`<br>`Read GET /repositories/{workspace}/{repo_slug}/deployments_config/environments/{environment_uuid}/variables`<br>`Update PUT /repositories/{workspace}/{repo_slug}/deployments_config/environments/{environment_uuid}/variables/{variable_uuid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/deployments_config/environments/{environment_uuid}/variables/{variable_uuid}`

#### Legacy HCL

```hcl
resource "bitbucket_deployment_variable" "legacy" {
  deployment = "example-deployment"
  key = "example-key"
  value = "example-value"

  # secured = true  # optional
}
```

- New operations: `Create POST /repositories/{workspace}/{repo_slug}/deployments_config/environments/{environment_uuid}/variables`<br>`Read GET /repositories/{workspace}/{repo_slug}/deployments_config/environments/{environment_uuid}/variables`<br>`Update PUT /repositories/{workspace}/{repo_slug}/deployments_config/environments/{environment_uuid}/variables/{variable_uuid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/deployments_config/environments/{environment_uuid}/variables/{variable_uuid}`

#### New HCL

##### `bitbucket_deployment_variables`

```hcl
resource "bitbucket_deployment_variables" "migrated" {
  environment_uuid = "{environment-uuid}"
  key = "example-key"
  repo_slug = "my-repo"
  value = "example-value"
  workspace = "my-workspace"

  # deployment = "example-deployment"  # legacy-only
  # request_body = jsonencode({})  # new-only
  # secured = true  # optional
  # variable_uuid = "{variable-uuid}"  # new-only
}
```

- Diff summary: legacy-only inputs: `deployment`; new-only inputs: `environment_uuid`, `repo_slug`, `request_body`, `uuid`, `variable_uuid`, `workspace`

### `bitbucket_forked_repository`

- New equivalent(s): `bitbucket_forked_repository`
- Legacy endpoints: `POST /repositories/{workspace}/{repo_slug}/forks`<br>`GET /repositories/{workspace}/{repo_slug}`

#### Legacy HCL

```hcl
resource "bitbucket_forked_repository" "legacy" {
  name = "my-repo"
  owner = "my-workspace"

  # description = "Example description"  # optional
  # fork_policy = "example-fork-policy"  # optional
  # has_issues = true  # optional
  # has_wiki = true  # optional
  # is_private = true  # optional
  # language = "go"  # optional
  # link = "example-link"  # optional
  # pipelines_enabled = "example-pipelines-enabled"  # optional
  # project_key = "EXAMPLE"  # optional
  # slug = "my-repo"  # optional
  # website = "https://example.com"  # optional
}
```

- New operations: `Create POST /repositories/{workspace}/{repo_slug}/forks`<br>`List GET /repositories/{workspace}/{repo_slug}/forks`

#### New HCL

##### `bitbucket_forked_repository`

```hcl
resource "bitbucket_forked_repository" "migrated" {
  name = "my-repo"
  repo_slug = "my-repo"
  workspace = "my-workspace"

  # description = "Example description"  # optional
  # fork_policy = "example-fork-policy"  # optional
  # has_issues = true  # optional
  # has_wiki = true  # optional
  # is_private = true  # optional
  # language = "go"  # optional
  # link = "example-link"  # legacy-only
  # pipelines_enabled = "example-pipelines-enabled"  # legacy-only
  # project_key = "EXAMPLE"  # legacy-only
  # request_body = jsonencode({})  # new-only
  # slug = "my-repo"  # legacy-only
  # website = "https://example.com"  # legacy-only
}
```

- Diff summary: renamed: `owner` → `workspace`; legacy-only inputs: `link`, `pipelines_enabled`, `project_key`, `slug`, `website`; new-only inputs: `full_name`, `mainbranch`, `owner`, `project`, `repo_slug`, `request_body`, `scm`, `size`, `uuid`

### `bitbucket_hook`

- New equivalent(s): `bitbucket_hooks`
- Legacy endpoints: `Create POST /repositories/{workspace}/{repo_slug}/hooks`<br>`Read GET /repositories/{workspace}/{repo_slug}/hooks/{uid}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/hooks/{uid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/hooks/{uid}`<br>`List GET /repositories/{workspace}/{repo_slug}/hooks`

#### Legacy HCL

```hcl
resource "bitbucket_hook" "legacy" {
  description = "Example description"
  events = ["repo:push"]
  owner = "my-workspace"
  repository = "my-repo"
  url = "https://example.com/webhook"

  # active = true  # optional
  # secret = "example-secret"  # optional
  # skip_cert_verification = false  # optional
}
```

- New operations: `Create POST /repositories/{workspace}/{repo_slug}/hooks`<br>`Read GET /repositories/{workspace}/{repo_slug}/hooks/{uid}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/hooks/{uid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/hooks/{uid}`<br>`List GET /repositories/{workspace}/{repo_slug}/hooks`

#### New HCL

##### `bitbucket_hooks`

```hcl
resource "bitbucket_hooks" "migrated" {
  repo_slug = "my-repo"
  workspace = "my-workspace"

  # active = true  # legacy-only
  # description = "Example description"  # legacy-only
  # events = ["repo:push"]  # legacy-only
  # secret = "example-secret"  # legacy-only
  # skip_cert_verification = false  # legacy-only
  # uid = "example-uid"  # new-only
  # url = "https://example.com/webhook"  # legacy-only
}
```

- Diff summary: renamed: `owner` → `workspace`, `repository` → `repo_slug`; legacy-only inputs: `active`, `description`, `events`, `secret`, `skip_cert_verification`, `url`; new-only inputs: `uid`

### `bitbucket_pipeline_schedule`

- New equivalent(s): `bitbucket_pipeline_schedules`
- Legacy endpoints: `Create POST /repositories/{workspace}/{repo_slug}/pipelines_config/schedules`<br>`Read GET /repositories/{workspace}/{repo_slug}/pipelines_config/schedules/{schedule_uuid}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/pipelines_config/schedules/{schedule_uuid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/pipelines_config/schedules/{schedule_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/pipelines_config/schedules`

#### Legacy HCL

```hcl
resource "bitbucket_pipeline_schedule" "legacy" {
  cron_pattern = "0 0 * * *"
  enabled = true
  repository = "my-repo"
  target = jsonencode({ ref_name = "main" })
  workspace = "my-workspace"
}
```

- New operations: `Create POST /repositories/{workspace}/{repo_slug}/pipelines_config/schedules`<br>`Read GET /repositories/{workspace}/{repo_slug}/pipelines_config/schedules/{schedule_uuid}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/pipelines_config/schedules/{schedule_uuid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/pipelines_config/schedules/{schedule_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/pipelines_config/schedules`

#### New HCL

##### `bitbucket_pipeline_schedules`

```hcl
resource "bitbucket_pipeline_schedules" "migrated" {
  cron_pattern = "0 0 * * *"
  enabled = true
  repo_slug = "my-repo"
  target = jsonencode({ ref_name = "main" })
  workspace = "my-workspace"

  # request_body = jsonencode({})  # new-only
  # schedule_uuid = "{schedule-uuid}"  # new-only
}
```

- Diff summary: renamed: `repository` → `repo_slug`; new-only inputs: `request_body`, `schedule_uuid`

### `bitbucket_pipeline_ssh_key`

- New equivalent(s): `bitbucket_pipeline_ssh_keys`
- Legacy endpoints: `Read GET /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/key_pair`<br>`Update PUT /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/key_pair`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/key_pair`

#### Legacy HCL

```hcl
resource "bitbucket_pipeline_ssh_key" "legacy" {
  private_key = "---PRIVATE KEY---"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDemo"
  repository = "my-repo"
  workspace = "my-workspace"
}
```

- New operations: `Read GET /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/key_pair`<br>`Update PUT /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/key_pair`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/key_pair`

#### New HCL

##### `bitbucket_pipeline_ssh_keys`

```hcl
resource "bitbucket_pipeline_ssh_keys" "migrated" {
  private_key = "---PRIVATE KEY---"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDemo"
  repo_slug = "my-repo"
  workspace = "my-workspace"

  # request_body = jsonencode({})  # new-only
}
```

- Diff summary: renamed: `repository` → `repo_slug`; new-only inputs: `request_body`

### `bitbucket_pipeline_ssh_known_host`

- New equivalent(s): `bitbucket_pipeline_known_hosts`
- Legacy endpoints: `Create POST /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts`<br>`Read GET /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts/{known_host_uuid}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts/{known_host_uuid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts/{known_host_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts`

#### Legacy HCL

```hcl
resource "bitbucket_pipeline_ssh_known_host" "legacy" {
  hostname = "github.com"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDemo"
  repository = "my-repo"
  workspace = "my-workspace"
}
```

- New operations: `Create POST /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts`<br>`Read GET /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts/{known_host_uuid}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts/{known_host_uuid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts/{known_host_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/pipelines_config/ssh/known_hosts`

#### New HCL

##### `bitbucket_pipeline_known_hosts`

```hcl
resource "bitbucket_pipeline_known_hosts" "migrated" {
  hostname = "github.com"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDemo"
  repo_slug = "my-repo"
  workspace = "my-workspace"

  # known_host_uuid = "{known-host-uuid}"  # new-only
  # request_body = jsonencode({})  # new-only
}
```

- Diff summary: renamed: `repository` → `repo_slug`; new-only inputs: `known_host_uuid`, `request_body`, `uuid`

### `bitbucket_project`

- New equivalent(s): `bitbucket_projects`
- Legacy endpoints: `PUT /workspaces/{workspace}/projects/{project_key}`<br>`POST /workspaces/{workspace}/projects`<br>`GET /workspaces/{workspace}/projects/{project_key}`<br>`DELETE /workspaces/{workspace}/projects/{project_key}`

#### Legacy HCL

```hcl
resource "bitbucket_project" "legacy" {
  key = "example-key"
  name = "my-repo"
  owner = "my-workspace"

  # description = "Example description"  # optional
  # is_private = true  # optional
  # link = "example-link"  # optional
}
```

- New operations: `Create POST /workspaces/{workspace}/projects`<br>`Read GET /workspaces/{workspace}/projects/{project_key}`<br>`Update PUT /workspaces/{workspace}/projects/{project_key}`<br>`Delete DELETE /workspaces/{workspace}/projects/{project_key}`<br>`List GET /workspaces/{workspace}/projects`

#### New HCL

##### `bitbucket_projects`

```hcl
resource "bitbucket_projects" "migrated" {
  workspace = "my-workspace"

  # description = "Example description"  # legacy-only
  # is_private = true  # legacy-only
  # key = "example-key"  # legacy-only
  # link = "example-link"  # legacy-only
  # name = "my-repo"  # legacy-only
  # project_key = "EXAMPLE"  # new-only
  # request_body = jsonencode({})  # new-only
}
```

- Diff summary: renamed: `owner` → `workspace`; legacy-only inputs: `description`, `is_private`, `key`, `link`, `name`; new-only inputs: `project_key`, `request_body`

### `bitbucket_project_branching_model`

- New equivalent(s): `bitbucket_project_branching_model`
- Legacy endpoints: `Read GET /workspaces/{workspace}/projects/{project_key}/branching-model`<br>`Update PUT /workspaces/{workspace}/projects/{project_key}/branching-model/settings`

#### Legacy HCL

```hcl
resource "bitbucket_project_branching_model" "legacy" {
  branch_type = "feature"
  project = "EXAMPLE"
  workspace = "my-workspace"

  # development = "example-development"  # optional
  # production = "example-production"  # optional
}
```

- New operations: `Read GET /workspaces/{workspace}/projects/{project_key}/branching-model`<br>`Update PUT /workspaces/{workspace}/projects/{project_key}/branching-model/settings`

#### New HCL

##### `bitbucket_project_branching_model`

```hcl
resource "bitbucket_project_branching_model" "migrated" {
  project_key = "EXAMPLE"
  workspace = "my-workspace"

  # branch_type = "feature"  # legacy-only
  # development = "example-development"  # legacy-only
  # production = "example-production"  # legacy-only
  # project = "EXAMPLE"  # legacy-only
}
```

- Diff summary: legacy-only inputs: `branch_type`, `development`, `production`, `project`; new-only inputs: `project_key`

### `bitbucket_project_default_reviewers`

- New equivalent(s): `bitbucket_project_default_reviewers`
- Legacy endpoints: `PUT /workspaces/{workspace}/projects/{project_key}/default-reviewers/{selected_user}`<br>`GET /workspaces/{workspace}/projects/{project_key}/default-reviewers`<br>`DELETE /workspaces/{workspace}/projects/{project_key}/default-reviewers/{selected_user}`

#### Legacy HCL

```hcl
resource "bitbucket_project_default_reviewers" "legacy" {
  project = "EXAMPLE"
  reviewers = ["example-user"]
  workspace = "my-workspace"
}
```

- New operations: `Create PUT /workspaces/{workspace}/projects/{project_key}/default-reviewers/{selected_user}`<br>`Read GET /workspaces/{workspace}/projects/{project_key}/default-reviewers/{selected_user}`<br>`Delete DELETE /workspaces/{workspace}/projects/{project_key}/default-reviewers/{selected_user}`<br>`List GET /workspaces/{workspace}/projects/{project_key}/default-reviewers`

#### New HCL

##### `bitbucket_project_default_reviewers`

```hcl
resource "bitbucket_project_default_reviewers" "migrated" {
  project_key = "EXAMPLE"
  selected_user = "example-user"
  workspace = "my-workspace"

  # project = "EXAMPLE"  # legacy-only
  # reviewers = ["example-user"]  # legacy-only
}
```

- Diff summary: legacy-only inputs: `project`, `reviewers`; new-only inputs: `project_key`, `selected_user`

### `bitbucket_project_group_permission`

- New equivalent(s): `bitbucket_project_group_permissions`
- Legacy endpoints: `Read GET /workspaces/{workspace}/projects/{project_key}/permissions-config/groups/{group_slug}`<br>`Update PUT /workspaces/{workspace}/projects/{project_key}/permissions-config/groups/{group_slug}`<br>`Delete DELETE /workspaces/{workspace}/projects/{project_key}/permissions-config/groups/{group_slug}`<br>`List GET /workspaces/{workspace}/projects/{project_key}/permissions-config/groups`

#### Legacy HCL

```hcl
resource "bitbucket_project_group_permission" "legacy" {
  group_slug = "example-group"
  permission = "read"
  project_key = "EXAMPLE"
  workspace = "my-workspace"
}
```

- New operations: `Read GET /workspaces/{workspace}/projects/{project_key}/permissions-config/groups/{group_slug}`<br>`Update PUT /workspaces/{workspace}/projects/{project_key}/permissions-config/groups/{group_slug}`<br>`Delete DELETE /workspaces/{workspace}/projects/{project_key}/permissions-config/groups/{group_slug}`<br>`List GET /workspaces/{workspace}/projects/{project_key}/permissions-config/groups`

#### New HCL

##### `bitbucket_project_group_permissions`

```hcl
resource "bitbucket_project_group_permissions" "migrated" {
  group_slug = "example-group"
  project_key = "EXAMPLE"
  workspace = "my-workspace"

  # permission = "read"  # legacy-only
  # request_body = jsonencode({})  # new-only
}
```

- Diff summary: legacy-only inputs: `permission`; new-only inputs: `request_body`

### `bitbucket_project_user_permission`

- New equivalent(s): `bitbucket_project_user_permissions`
- Legacy endpoints: `Read GET /workspaces/{workspace}/projects/{project_key}/permissions-config/users/{selected_user_id}`<br>`Update PUT /workspaces/{workspace}/projects/{project_key}/permissions-config/users/{selected_user_id}`<br>`Delete DELETE /workspaces/{workspace}/projects/{project_key}/permissions-config/users/{selected_user_id}`<br>`List GET /workspaces/{workspace}/projects/{project_key}/permissions-config/users`

#### Legacy HCL

```hcl
resource "bitbucket_project_user_permission" "legacy" {
  permission = "read"
  project_key = "EXAMPLE"
  user_id = "{user-uuid}"
  workspace = "my-workspace"
}
```

- New operations: `Read GET /workspaces/{workspace}/projects/{project_key}/permissions-config/users/{selected_user_id}`<br>`Update PUT /workspaces/{workspace}/projects/{project_key}/permissions-config/users/{selected_user_id}`<br>`Delete DELETE /workspaces/{workspace}/projects/{project_key}/permissions-config/users/{selected_user_id}`<br>`List GET /workspaces/{workspace}/projects/{project_key}/permissions-config/users`

#### New HCL

##### `bitbucket_project_user_permissions`

```hcl
resource "bitbucket_project_user_permissions" "migrated" {
  project_key = "EXAMPLE"
  selected_user_id = "{user-uuid}"
  workspace = "my-workspace"

  # permission = "read"  # legacy-only
  # request_body = jsonencode({})  # new-only
  # user_id = "{user-uuid}"  # legacy-only
}
```

- Diff summary: legacy-only inputs: `permission`, `user_id`; new-only inputs: `request_body`, `selected_user_id`

### `bitbucket_repository`

- New equivalent(s): `bitbucket_repos`, `bitbucket_repo_settings`, `bitbucket_pipeline_config`
- Legacy endpoints: `PUT /repositories/{workspace}/{repo_slug}`<br>`POST /repositories/{workspace}/{repo_slug}`<br>`GET /repositories/{workspace}/{repo_slug}`<br>`DELETE /repositories/{workspace}/{repo_slug}`

#### Legacy HCL

```hcl
resource "bitbucket_repository" "legacy" {
  name = "my-repo"
  owner = "my-workspace"

  # description = "Example description"  # optional
  # fork_policy = "example-fork-policy"  # optional
  # has_issues = true  # optional
  # has_wiki = true  # optional
  # inherit_branching_model = true  # optional
  # inherit_default_merge_strategy = true  # optional
  # is_private = true  # optional
  # language = "go"  # optional
  # link = "example-link"  # optional
  # pipelines_enabled = "example-pipelines-enabled"  # optional
  # project_key = "EXAMPLE"  # optional
  # scm = "git"  # optional
  # slug = "my-repo"  # optional
  # website = "https://example.com"  # optional
}
```

- New operations: `Create POST /repositories/{workspace}/{repo_slug}`<br>`Read GET /repositories/{workspace}/{repo_slug}`<br>`Update PUT /repositories/{workspace}/{repo_slug}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}`<br>`List GET /repositories/{workspace}`<br>`Read GET /repositories/{workspace}/{repo_slug}/override-settings`<br>`Update PUT /repositories/{workspace}/{repo_slug}/override-settings`<br>`Read GET /repositories/{workspace}/{repo_slug}/pipelines_config`<br>`Update PUT /repositories/{workspace}/{repo_slug}/pipelines_config`

#### New HCL

##### `bitbucket_repos`

```hcl
resource "bitbucket_repos" "migrated" {
  name = "my-repo"
  repo_slug = "my-repo"
  workspace = "my-workspace"

  # description = "Example description"  # optional
  # fork_policy = "example-fork-policy"  # optional
  # has_issues = true  # optional
  # has_wiki = true  # optional
  # inherit_branching_model = true  # legacy-only
  # inherit_default_merge_strategy = true  # legacy-only
  # is_private = true  # optional
  # language = "go"  # optional
  # link = "example-link"  # legacy-only
  # pipelines_enabled = "example-pipelines-enabled"  # legacy-only
  # project_key = "EXAMPLE"  # legacy-only
  # request_body = jsonencode({})  # new-only
  # scm = "git"  # optional
  # slug = "my-repo"  # legacy-only
  # website = "https://example.com"  # legacy-only
}
```

##### `bitbucket_repo_settings`

```hcl
resource "bitbucket_repo_settings" "migrated" {
  repo_slug = "my-repo"
  workspace = "my-workspace"

  # description = "Example description"  # legacy-only
  # fork_policy = "example-fork-policy"  # legacy-only
  # has_issues = true  # legacy-only
  # has_wiki = true  # legacy-only
  # inherit_branching_model = true  # legacy-only
  # inherit_default_merge_strategy = true  # legacy-only
  # is_private = true  # legacy-only
  # language = "go"  # legacy-only
  # link = "example-link"  # legacy-only
  # name = "my-repo"  # legacy-only
  # pipelines_enabled = "example-pipelines-enabled"  # legacy-only
  # project_key = "EXAMPLE"  # legacy-only
  # scm = "git"  # legacy-only
  # slug = "my-repo"  # legacy-only
  # website = "https://example.com"  # legacy-only
}
```

##### `bitbucket_pipeline_config`

```hcl
resource "bitbucket_pipeline_config" "migrated" {
  repo_slug = "my-repo"
  workspace = "my-workspace"

  # description = "Example description"  # legacy-only
  # fork_policy = "example-fork-policy"  # legacy-only
  # has_issues = true  # legacy-only
  # has_wiki = true  # legacy-only
  # inherit_branching_model = true  # legacy-only
  # inherit_default_merge_strategy = true  # legacy-only
  # is_private = true  # legacy-only
  # language = "go"  # legacy-only
  # link = "example-link"  # legacy-only
  # name = "my-repo"  # legacy-only
  # pipelines_enabled = "example-pipelines-enabled"  # legacy-only
  # project_key = "EXAMPLE"  # legacy-only
  # request_body = jsonencode({})  # new-only
  # scm = "git"  # legacy-only
  # slug = "my-repo"  # legacy-only
  # website = "https://example.com"  # legacy-only
}
```

- Diff summary: renamed: `owner` → `workspace`; legacy-only inputs: `inherit_branching_model`, `inherit_default_merge_strategy`, `link`, `pipelines_enabled`, `project_key`, `slug`, `website`; new-only inputs: `enabled`, `full_name`, `mainbranch`, `owner`, `project`, `repo_slug`, `repository`, `request_body`, `size`, `uuid`
- Notes: The legacy repository resource bundled core repository CRUD, pipeline enablement, and override-settings flags. In the new provider, core CRUD stays on `bitbucket_repos`, pipeline enablement moves to `bitbucket_pipeline_config`, and repository settings have their own `bitbucket_repo_settings` resource.

### `bitbucket_repository_group_permission`

- New equivalent(s): `bitbucket_repo_group_permissions`
- Legacy endpoints: `Read GET /repositories/{workspace}/{repo_slug}/permissions-config/groups/{group_slug}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/permissions-config/groups/{group_slug}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/permissions-config/groups/{group_slug}`<br>`List GET /repositories/{workspace}/{repo_slug}/permissions-config/groups`

#### Legacy HCL

```hcl
resource "bitbucket_repository_group_permission" "legacy" {
  group_slug = "example-group"
  permission = "read"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

- New operations: `Read GET /repositories/{workspace}/{repo_slug}/permissions-config/groups/{group_slug}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/permissions-config/groups/{group_slug}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/permissions-config/groups/{group_slug}`<br>`List GET /repositories/{workspace}/{repo_slug}/permissions-config/groups`

#### New HCL

##### `bitbucket_repo_group_permissions`

```hcl
resource "bitbucket_repo_group_permissions" "migrated" {
  group_slug = "example-group"
  repo_slug = "my-repo"
  workspace = "my-workspace"

  # permission = "read"  # legacy-only
  # request_body = jsonencode({})  # new-only
}
```

- Diff summary: legacy-only inputs: `permission`; new-only inputs: `request_body`

### `bitbucket_repository_user_permission`

- New equivalent(s): `bitbucket_repo_user_permissions`
- Legacy endpoints: `Read GET /repositories/{workspace}/{repo_slug}/permissions-config/users/{selected_user_id}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/permissions-config/users/{selected_user_id}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/permissions-config/users/{selected_user_id}`<br>`List GET /repositories/{workspace}/{repo_slug}/permissions-config/users`

#### Legacy HCL

```hcl
resource "bitbucket_repository_user_permission" "legacy" {
  permission = "read"
  repo_slug = "my-repo"
  user_id = "{user-uuid}"
  workspace = "my-workspace"
}
```

- New operations: `Read GET /repositories/{workspace}/{repo_slug}/permissions-config/users/{selected_user_id}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/permissions-config/users/{selected_user_id}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/permissions-config/users/{selected_user_id}`<br>`List GET /repositories/{workspace}/{repo_slug}/permissions-config/users`

#### New HCL

##### `bitbucket_repo_user_permissions`

```hcl
resource "bitbucket_repo_user_permissions" "migrated" {
  repo_slug = "my-repo"
  selected_user_id = "{user-uuid}"
  workspace = "my-workspace"

  # permission = "read"  # legacy-only
  # request_body = jsonencode({})  # new-only
  # user_id = "{user-uuid}"  # legacy-only
}
```

- Diff summary: legacy-only inputs: `permission`, `user_id`; new-only inputs: `request_body`, `selected_user_id`

### `bitbucket_repository_variable`

- New equivalent(s): `bitbucket_pipeline_variables`
- Legacy endpoints: `Create POST /repositories/{workspace}/{repo_slug}/pipelines_config/variables`<br>`Read GET /repositories/{workspace}/{repo_slug}/pipelines_config/variables/{variable_uuid}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/pipelines_config/variables/{variable_uuid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/pipelines_config/variables/{variable_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/pipelines_config/variables`

#### Legacy HCL

```hcl
resource "bitbucket_repository_variable" "legacy" {
  key = "example-key"
  repository = "my-repo"
  value = "example-value"

  # secured = true  # optional
}
```

- New operations: `Create POST /repositories/{workspace}/{repo_slug}/pipelines_config/variables`<br>`Read GET /repositories/{workspace}/{repo_slug}/pipelines_config/variables/{variable_uuid}`<br>`Update PUT /repositories/{workspace}/{repo_slug}/pipelines_config/variables/{variable_uuid}`<br>`Delete DELETE /repositories/{workspace}/{repo_slug}/pipelines_config/variables/{variable_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/pipelines_config/variables`

#### New HCL

##### `bitbucket_pipeline_variables`

```hcl
resource "bitbucket_pipeline_variables" "migrated" {
  key = "example-key"
  repo_slug = "my-repo"
  value = "example-value"
  workspace = "my-workspace"

  # request_body = jsonencode({})  # new-only
  # secured = true  # optional
  # variable_uuid = "{variable-uuid}"  # new-only
}
```

- Diff summary: renamed: `repository` → `repo_slug`; new-only inputs: `request_body`, `uuid`, `variable_uuid`, `workspace`
- Notes: Legacy repository variables map to the pipelines variable API. Use `bitbucket_pipeline_variables` and rename `owner`/`repository` to `workspace`/`repo_slug`.

### `bitbucket_ssh_key`

- New equivalent(s): `bitbucket_ssh_keys`
- Legacy endpoints: `POST /users/{selected_user}/ssh-keys`<br>`GET /users/{selected_user}/ssh-keys/{key}/{id}`<br>`PUT /users/{selected_user}/ssh-keys/{key}/{id}`<br>`DELETE /users/{selected_user}/ssh-keys/{key}/{id}`

#### Legacy HCL

```hcl
resource "bitbucket_ssh_key" "legacy" {
  key = "example-key"
  user = "example-user"

  # label = "Example label"  # optional
}
```

- New operations: `Create POST /users/{selected_user}/ssh-keys`<br>`Read GET /users/{selected_user}/ssh-keys/{key_id}`<br>`Update PUT /users/{selected_user}/ssh-keys/{key_id}`<br>`Delete DELETE /users/{selected_user}/ssh-keys/{key_id}`<br>`List GET /users/{selected_user}/ssh-keys`

#### New HCL

##### `bitbucket_ssh_keys`

```hcl
resource "bitbucket_ssh_keys" "migrated" {
  key = "example-key"
  selected_user = "example-user"

  # key_id = "example-key-id"  # new-only
  # label = "Example label"  # optional
  # request_body = jsonencode({})  # new-only
  # user = "example-user"  # legacy-only
}
```

- Diff summary: legacy-only inputs: `user`; new-only inputs: `comment`, `expires_on`, `fingerprint`, `key_id`, `last_used`, `owner`, `request_body`, `selected_user`, `uuid`

### `bitbucket_workspace_hook`

- New equivalent(s): `bitbucket_workspace_hooks`
- Legacy endpoints: `Create POST /workspaces/{workspace}/hooks`<br>`Read GET /workspaces/{workspace}/hooks/{uid}`<br>`Update PUT /workspaces/{workspace}/hooks/{uid}`<br>`Delete DELETE /workspaces/{workspace}/hooks/{uid}`<br>`List GET /workspaces/{workspace}/hooks`

#### Legacy HCL

```hcl
resource "bitbucket_workspace_hook" "legacy" {
  description = "Example description"
  events = ["repo:push"]
  url = "https://example.com/webhook"
  workspace = "my-workspace"

  # active = true  # optional
  # secret = "example-secret"  # optional
  # skip_cert_verification = false  # optional
}
```

- New operations: `Create POST /workspaces/{workspace}/hooks`<br>`Read GET /workspaces/{workspace}/hooks/{uid}`<br>`Update PUT /workspaces/{workspace}/hooks/{uid}`<br>`Delete DELETE /workspaces/{workspace}/hooks/{uid}`<br>`List GET /workspaces/{workspace}/hooks`

#### New HCL

##### `bitbucket_workspace_hooks`

```hcl
resource "bitbucket_workspace_hooks" "migrated" {
  workspace = "my-workspace"

  # active = true  # legacy-only
  # description = "Example description"  # legacy-only
  # events = ["repo:push"]  # legacy-only
  # secret = "example-secret"  # legacy-only
  # skip_cert_verification = false  # legacy-only
  # uid = "example-uid"  # new-only
  # url = "https://example.com/webhook"  # legacy-only
}
```

- Diff summary: legacy-only inputs: `active`, `description`, `events`, `secret`, `skip_cert_verification`, `url`; new-only inputs: `uid`

### `bitbucket_workspace_variable`

- New equivalent(s): `bitbucket_workspace_pipeline_variables`
- Legacy endpoints: `Create POST /workspaces/{workspace}/pipelines-config/variables`<br>`Read GET /workspaces/{workspace}/pipelines-config/variables/{variable_uuid}`<br>`Update PUT /workspaces/{workspace}/pipelines-config/variables/{variable_uuid}`<br>`Delete DELETE /workspaces/{workspace}/pipelines-config/variables/{variable_uuid}`<br>`List GET /workspaces/{workspace}/pipelines-config/variables`

#### Legacy HCL

```hcl
resource "bitbucket_workspace_variable" "legacy" {
  key = "example-key"
  value = "example-value"
  workspace = "my-workspace"

  # secured = true  # optional
}
```

- New operations: `Create POST /workspaces/{workspace}/pipelines-config/variables`<br>`Read GET /workspaces/{workspace}/pipelines-config/variables/{variable_uuid}`<br>`Update PUT /workspaces/{workspace}/pipelines-config/variables/{variable_uuid}`<br>`Delete DELETE /workspaces/{workspace}/pipelines-config/variables/{variable_uuid}`<br>`List GET /workspaces/{workspace}/pipelines-config/variables`

#### New HCL

##### `bitbucket_workspace_pipeline_variables`

```hcl
resource "bitbucket_workspace_pipeline_variables" "migrated" {
  workspace = "my-workspace"

  # key = "example-key"  # legacy-only
  # request_body = jsonencode({})  # new-only
  # secured = true  # legacy-only
  # value = "example-value"  # legacy-only
  # variable_uuid = "{variable-uuid}"  # new-only
}
```

- Diff summary: legacy-only inputs: `key`, `secured`, `value`; new-only inputs: `request_body`, `variable_uuid`
- Notes: Workspace variables now live under the pipelines API as `bitbucket_workspace_pipeline_variables`.

## Legacy-only resources

### `bitbucket_group`

- New equivalent(s): none
- Legacy endpoints: none

#### Legacy HCL

```hcl
resource "bitbucket_group" "legacy" {
  name = "my-repo"
  workspace = "my-workspace"

  # auto_add = true  # optional
  # permission = "read"  # optional
}
```
- Notes: Workspace group management is not currently exposed by the generated provider because those endpoints are not represented in the generated Terraform docs.

### `bitbucket_group_membership`

- New equivalent(s): none
- Legacy endpoints: none

#### Legacy HCL

```hcl
resource "bitbucket_group_membership" "legacy" {
  group_slug = "example-group"
  uuid = "{resource-uuid}"
  workspace = "my-workspace"
}
```
- Notes: Group membership management is not currently exposed by the generated provider.

## Matched legacy data sources

### `bitbucket_current_user`

- New equivalent(s): `bitbucket_current_user`
- Legacy endpoints: `GET /user`<br>`GET /user/emails`

#### Legacy HCL

```hcl
data "bitbucket_current_user" "legacy" {
}
```

- New operations: `Read GET /user`

#### New HCL

##### `bitbucket_current_user`

```hcl
data "bitbucket_current_user" "migrated" {
}
```

- Diff summary: input names are effectively unchanged
- Notes: The legacy data source also fetched `/user/emails`. The generated provider splits that into `bitbucket_current_user` plus `bitbucket_user_emails` when you need email addresses.

### `bitbucket_deployment`

- New equivalent(s): `bitbucket_deployments`
- Legacy endpoints: `Read GET /repositories/{workspace}/{repo_slug}/environments/{environment_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/environments`

#### Legacy HCL

```hcl
data "bitbucket_deployment" "legacy" {
  repository = "my-repo"
  uuid = "{resource-uuid}"
  workspace = "my-workspace"
}
```

- New operations: `Read GET /repositories/{workspace}/{repo_slug}/environments/{environment_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/environments`

#### New HCL

##### `bitbucket_deployments`

```hcl
data "bitbucket_deployments" "migrated" {
  repo_slug = "my-repo"
  workspace = "my-workspace"

  # environment_uuid = "{environment-uuid}"  # new-only
  # uuid = "{resource-uuid}"  # legacy-only
}
```

- Diff summary: renamed: `repository` → `repo_slug`; legacy-only inputs: `uuid`; new-only inputs: `environment_uuid`
- Notes: Use `bitbucket_deployments` with the identifying path parameters for a single deployment; omit the single-resource expectation and treat the response as the generic deployment payload.

### `bitbucket_deployments`

- New equivalent(s): `bitbucket_deployments`
- Legacy endpoints: `Read GET /repositories/{workspace}/{repo_slug}/environments/{environment_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/environments`

#### Legacy HCL

```hcl
data "bitbucket_deployments" "legacy" {
  repository = "my-repo"
  workspace = "my-workspace"
}
```

- New operations: `Read GET /repositories/{workspace}/{repo_slug}/environments/{environment_uuid}`<br>`List GET /repositories/{workspace}/{repo_slug}/environments`

#### New HCL

##### `bitbucket_deployments`

```hcl
data "bitbucket_deployments" "migrated" {
  repo_slug = "my-repo"
  workspace = "my-workspace"

  # environment_uuid = "{environment-uuid}"  # new-only
}
```

- Diff summary: renamed: `repository` → `repo_slug`; new-only inputs: `environment_uuid`

### `bitbucket_file`

- New equivalent(s): `bitbucket_commit_file`
- Legacy endpoints: `Read GET /repositories/{workspace}/{repo_slug}/src/{commit}/{path}`

#### Legacy HCL

```hcl
data "bitbucket_file" "legacy" {
}
```

- New operations: `Read GET /repositories/{workspace}/{repo_slug}/src/{commit}/{path}`

#### New HCL

##### `bitbucket_commit_file`

```hcl
data "bitbucket_commit_file" "migrated" {
  commit = "main"
  path = "README.md"
  repo_slug = "my-repo"
  workspace = "my-workspace"
}
```

- Diff summary: new-only inputs: `commit`, `path`, `repo_slug`, `workspace`
- Notes: The legacy `bitbucket_file` data source maps most closely to `bitbucket_commit_file`, which reads file content via the commit-file endpoint.

### `bitbucket_hook_types`

- New equivalent(s): `bitbucket_hook_types`
- Legacy endpoints: `GET /hook-events-subject-type`

#### Legacy HCL

```hcl
data "bitbucket_hook_types" "legacy" {
}
```

- New operations: `Read GET /hook_events`<br>`List GET /hook_events/{subject_type}`

#### New HCL

##### `bitbucket_hook_types`

```hcl
data "bitbucket_hook_types" "migrated" {
}
```

- Diff summary: input names are effectively unchanged

### `bitbucket_pipeline_oidc_config`

- New equivalent(s): `bitbucket_pipeline_oidc`
- Legacy endpoints: `Read GET /workspaces/{workspace}/pipelines-config/identity/oidc/.well-known/openid-configuration`

#### Legacy HCL

```hcl
data "bitbucket_pipeline_oidc_config" "legacy" {
  workspace = "my-workspace"
}
```

- New operations: `Read GET /workspaces/{workspace}/pipelines-config/identity/oidc/.well-known/openid-configuration`

#### New HCL

##### `bitbucket_pipeline_oidc`

```hcl
data "bitbucket_pipeline_oidc" "migrated" {
  workspace = "my-workspace"
}
```

- Diff summary: input names are effectively unchanged

### `bitbucket_pipeline_oidc_config_keys`

- New equivalent(s): `bitbucket_pipeline_oidc_keys`
- Legacy endpoints: `Read GET /workspaces/{workspace}/pipelines-config/identity/oidc/keys.json`

#### Legacy HCL

```hcl
data "bitbucket_pipeline_oidc_config_keys" "legacy" {
  workspace = "my-workspace"
}
```

- New operations: `Read GET /workspaces/{workspace}/pipelines-config/identity/oidc/keys.json`

#### New HCL

##### `bitbucket_pipeline_oidc_keys`

```hcl
data "bitbucket_pipeline_oidc_keys" "migrated" {
  workspace = "my-workspace"
}
```

- Diff summary: input names are effectively unchanged

### `bitbucket_project`

- New equivalent(s): `bitbucket_projects`
- Legacy endpoints: `Read GET /workspaces/{workspace}/projects/{project_key}`<br>`List GET /workspaces/{workspace}/projects`

#### Legacy HCL

```hcl
data "bitbucket_project" "legacy" {
}
```

- New operations: `Read GET /workspaces/{workspace}/projects/{project_key}`<br>`List GET /workspaces/{workspace}/projects`

#### New HCL

##### `bitbucket_projects`

```hcl
data "bitbucket_projects" "migrated" {
  workspace = "my-workspace"

  # project_key = "EXAMPLE"  # new-only
}
```

- Diff summary: new-only inputs: `project_key`, `workspace`

### `bitbucket_repository`

- New equivalent(s): `bitbucket_repos`
- Legacy endpoints: `Read GET /repositories/{workspace}/{repo_slug}`<br>`List GET /repositories/{workspace}`

#### Legacy HCL

```hcl
data "bitbucket_repository" "legacy" {
}
```

- New operations: `Read GET /repositories/{workspace}/{repo_slug}`<br>`List GET /repositories/{workspace}`

#### New HCL

##### `bitbucket_repos`

```hcl
data "bitbucket_repos" "migrated" {
  workspace = "my-workspace"

  # repo_slug = "my-repo"  # new-only
}
```

- Diff summary: new-only inputs: `repo_slug`, `workspace`

### `bitbucket_user`

- New equivalent(s): `bitbucket_users`
- Legacy endpoints: `GET /users/{selected_user}`

#### Legacy HCL

```hcl
data "bitbucket_user" "legacy" {
  # uuid = "{resource-uuid}"  # optional
}
```

- New operations: `Read GET /users/{selected_user}`<br>`List GET /users/{selected_user}/ssh-keys`

#### New HCL

##### `bitbucket_users`

```hcl
data "bitbucket_users" "migrated" {
  selected_user = "example-user"

  # uuid = "{resource-uuid}"  # legacy-only
}
```

- Diff summary: legacy-only inputs: `uuid`; new-only inputs: `selected_user`

### `bitbucket_workspace`

- New equivalent(s): `bitbucket_workspaces`
- Legacy endpoints: `GET /workspaces/{workspace}`

#### Legacy HCL

```hcl
data "bitbucket_workspace" "legacy" {
  workspace = "my-workspace"
}
```

- New operations: `Read GET /workspaces/{workspace}`<br>`List GET /workspaces`

#### New HCL

##### `bitbucket_workspaces`

```hcl
data "bitbucket_workspaces" "migrated" {
  workspace = "my-workspace"
}
```

- Diff summary: input names are effectively unchanged

### `bitbucket_workspace_members`

- New equivalent(s): `bitbucket_workspace_members`
- Legacy endpoints: `GET /workspaces/{workspace}/members`

#### Legacy HCL

```hcl
data "bitbucket_workspace_members" "legacy" {
  workspace = "my-workspace"
}
```

- New operations: `Read GET /workspaces/{workspace}/members/{member}`<br>`List GET /workspaces/{workspace}/members`

#### New HCL

##### `bitbucket_workspace_members`

```hcl
data "bitbucket_workspace_members" "migrated" {
  workspace = "my-workspace"

  # member = "example-user"  # new-only
}
```

- Diff summary: new-only inputs: `member`

## Legacy-only data sources

### `bitbucket_group`

- New equivalent(s): none
- Legacy endpoints: none

#### Legacy HCL

```hcl
data "bitbucket_group" "legacy" {
  slug = "my-repo"
  workspace = "my-workspace"
}
```
- Notes: Group lookup is not currently exposed by the generated provider.

### `bitbucket_group_members`

- New equivalent(s): none
- Legacy endpoints: none

#### Legacy HCL

```hcl
data "bitbucket_group_members" "legacy" {
  slug = "my-repo"
  workspace = "my-workspace"
}
```
- Notes: Group member lookup is not currently exposed by the generated provider.

### `bitbucket_groups`

- New equivalent(s): none
- Legacy endpoints: none

#### Legacy HCL

```hcl
data "bitbucket_groups" "legacy" {
  workspace = "my-workspace"
}
```
- Notes: Group listing is not currently exposed by the generated provider.

### `bitbucket_ip_ranges`

- New equivalent(s): none
- Legacy endpoints: none

#### Legacy HCL

```hcl
data "bitbucket_ip_ranges" "legacy" {
}
```
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

