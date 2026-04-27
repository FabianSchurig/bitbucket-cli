# Auto-generated Terraform test configuration for bitbucket_project_branch_restrictions
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

data "bitbucket_project_branch_restrictions" "test" {
  workspace = var.workspace
  project_key = var.project_key
}
