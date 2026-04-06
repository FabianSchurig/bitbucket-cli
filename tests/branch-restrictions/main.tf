# Auto-generated Terraform test configuration for bitbucket_branch_restrictions
# This file defines the resources/data sources referenced by the test assertions.

terraform {
  required_providers {
    bitbucket = {
      source = "FabianSchurig/bitbucket"
    }
  }
}

variable "workspace" {
  type    = string
  default = "test-workspace"
}

variable "repo_slug" {
  type    = string
  default = "my-repo"
}

variable "param_id" {
  type    = string
  default = "1"
}

provider "bitbucket" {}

data "bitbucket_branch_restrictions" "test" {
  repo_slug = var.repo_slug
  workspace = var.workspace
  param_id = var.param_id
}

resource "bitbucket_branch_restrictions" "test" {
  repo_slug = var.repo_slug
  workspace = var.workspace
  branch_match_kind = "glob"
  kind = "require_approvals_to_merge"
  pattern = "main"
}
