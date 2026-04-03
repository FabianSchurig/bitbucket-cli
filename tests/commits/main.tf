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

variable "repo_slug" {
  type    = string
  default = "my-repo"
}

variable "commit" {
  type    = string
  default = "abc123def"
}

provider "bitbucket" {}

data "bitbucket_commits" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  commit = var.commit
}
