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

variable "param_id" {
  type    = string
  default = "1"
}

variable "repo_slug" {
  type    = string
  default = "my-repo"
}

provider "bitbucket" {}

data "bitbucket_branch_restrictions" "test" {
  param_id = var.param_id
  repo_slug = var.repo_slug
  workspace = var.workspace
}

resource "bitbucket_branch_restrictions" "test" {
  param_id = var.param_id
  repo_slug = var.repo_slug
  workspace = var.workspace
}
