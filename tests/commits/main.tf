# Auto-generated Terraform test configuration for bitbucket_commits
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

variable "commit" {
  type    = string
  default = "abc123def"
}

variable "repo_slug" {
  type    = string
  default = "my-repo"
}

provider "bitbucket" {}

data "bitbucket_commits" "test" {
  commit = var.commit
  repo_slug = var.repo_slug
  workspace = var.workspace
}
