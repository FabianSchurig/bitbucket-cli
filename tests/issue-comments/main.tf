# Auto-generated Terraform test configuration for bitbucket_issue_comments
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

variable "issue_id" {
  type    = string
  default = "1"
}

variable "repo_slug" {
  type    = string
  default = "my-repo"
}

variable "comment_id" {
  type    = string
  default = "1"
}

provider "bitbucket" {}

data "bitbucket_issue_comments" "test" {
  issue_id = var.issue_id
  repo_slug = var.repo_slug
  workspace = var.workspace
  comment_id = var.comment_id
}

resource "bitbucket_issue_comments" "test" {
  issue_id = var.issue_id
  repo_slug = var.repo_slug
  workspace = var.workspace
}
