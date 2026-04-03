# Auto-generated Terraform test configuration for bitbucket_reports
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

variable "report_id" {
  type    = string
  default = "report-uuid"
}

provider "bitbucket" {}

data "bitbucket_reports" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  commit = var.commit
  report_id = var.report_id
}

resource "bitbucket_reports" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  commit = var.commit
  report_id = var.report_id
}
