# Auto-generated Terraform test configuration for bitbucket_commit_statuses
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

variable "key" {
  type    = string
  default = "build-key"
}

variable "repo_slug" {
  type    = string
  default = "my-repo"
}

provider "bitbucket" {}

data "bitbucket_commit_statuses" "test" {
  commit = var.commit
  key = var.key
  repo_slug = var.repo_slug
  workspace = var.workspace
}

resource "bitbucket_commit_statuses" "test" {
  commit = var.commit
  key = var.key
  repo_slug = var.repo_slug
  workspace = var.workspace
}
