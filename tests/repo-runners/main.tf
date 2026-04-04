# Auto-generated Terraform test configuration for bitbucket_repo_runners
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

variable "runner_uuid" {
  type    = string
  default = "{runner-uuid}"
}

provider "bitbucket" {}

data "bitbucket_repo_runners" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  runner_uuid = var.runner_uuid
}

resource "bitbucket_repo_runners" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  runner_uuid = var.runner_uuid
}
