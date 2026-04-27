# Auto-generated Terraform test configuration for bitbucket_project_branch_restrictions_by_branch_type
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

variable "project_key" {
  type    = string
  default = "PROJ"
}

provider "bitbucket" {}

data "bitbucket_project_branch_restrictions_by_branch_type" "test" {
  workspace = var.workspace
  project_key = var.project_key
}

resource "bitbucket_project_branch_restrictions_by_branch_type" "test" {
  workspace = var.workspace
  project_key = var.project_key
  branch_type = var.branch_type
}
