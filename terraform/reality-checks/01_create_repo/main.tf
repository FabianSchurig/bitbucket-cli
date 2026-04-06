terraform {
  required_version = ">= 1.0.0"
  required_providers {
    bitbucket = { source = "FabianSchurig/bitbucket" }
    random = { source = "hashicorp/random" }
  }
}

provider "bitbucket" {}
provider "random" {}

resource "random_pet" "suffix" {
  length = 2
}

/* Resource: bitbucket_repos
   Docs: docs/resources/repos.md
   Notes: set `is_private = true` when the target workspace/project disallows public repositories.
*/
resource "bitbucket_repos" "repo" {
  repo_slug = "tf-reality-${random_pet.suffix.id}"
  workspace = var.workspace
  is_private = true
}

# Branch restriction: require 1 approval to merge into `main`
/* Resource: bitbucket_branch_restrictions
   Docs: docs/resources/branch-restrictions.md
*/
resource "bitbucket_branch_restrictions" "protect_main" {
  count             = var.create_branch_restriction ? 1 : 0
  repo_slug         = bitbucket_repos.repo.repo_slug
  workspace         = var.workspace
  kind              = "require_approvals_to_merge"
  branch_match_kind = "glob"
  pattern           = "main"
  value             = 1
}

/* Resource: bitbucket_hooks
   Docs: docs/resources/hooks.md
   Note: the provider examples expose only `repo_slug`/`workspace` for hooks;
         other fields (url/events/active/description) may be read-only. See docs.
*/
resource "bitbucket_hooks" "repo_hook" {
  count     = var.create_hook ? 1 : 0
  repo_slug = bitbucket_repos.repo.repo_slug
  workspace = var.workspace
  url       = var.hook_url
  events    = ["repo:push"]
}

/* Resource: bitbucket_pipeline_variables
   Docs: docs/resources/pipeline-variables.md
*/
resource "bitbucket_pipeline_variables" "pipeline_var" {
  workspace = var.workspace
  repo_slug = bitbucket_repos.repo.repo_slug
  key       = "TF_REALITY_VAR"
  value     = "hello-from-terraform"
  secured   = false
}

output "repo_slug" {
  value = bitbucket_repos.repo.repo_slug
}

output "branch_restriction_response" {
  value = length(bitbucket_branch_restrictions.protect_main) > 0 ? bitbucket_branch_restrictions.protect_main[0].api_response : null
}

output "hook_uuid" {
  value = length(bitbucket_hooks.repo_hook) > 0 ? bitbucket_hooks.repo_hook[0].id : null
}

output "pipeline_variable_uuid" {
  value = bitbucket_pipeline_variables.pipeline_var.uuid
}
