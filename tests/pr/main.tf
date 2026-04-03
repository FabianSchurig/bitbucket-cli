# Auto-generated Terraform test configuration for bitbucket_pr
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

variable "pull_request_id" {
  type    = string
  default = "1"
}

variable "repo_slug" {
  type    = string
  default = "my-repo"
}

provider "bitbucket" {}

data "bitbucket_pr" "test" {
  pull_request_id = var.pull_request_id
  repo_slug = var.repo_slug
  workspace = var.workspace
}

resource "bitbucket_pr" "test" {
  pull_request_id = var.pull_request_id
  repo_slug = var.repo_slug
  workspace = var.workspace
}
