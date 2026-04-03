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

variable "id" {
  type    = string
  default = "1"
}

provider "bitbucket" {}

data "bitbucket_branch_restrictions" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  id = var.id
}

resource "bitbucket_branch_restrictions" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  id = var.id
}
