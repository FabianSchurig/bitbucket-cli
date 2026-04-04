# Auto-generated Terraform test configuration for bitbucket_commit_file
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

variable "path" {
  type    = string
  default = "README.md"
}

variable "repo_slug" {
  type    = string
  default = "my-repo"
}

provider "bitbucket" {}

data "bitbucket_commit_file" "test" {
  commit = var.commit
  path = var.path
  repo_slug = var.repo_slug
  workspace = var.workspace
}

resource "bitbucket_commit_file" "test" {
  repo_slug = var.repo_slug
  workspace = var.workspace
}
