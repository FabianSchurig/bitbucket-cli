# Auto-generated Terraform test configuration for bitbucket_issues
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

variable "issue_id" {
  type    = string
  default = "1"
}

provider "bitbucket" {}

data "bitbucket_issues" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  issue_id = var.issue_id
}

resource "bitbucket_issues" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  issue_id = var.issue_id
}
