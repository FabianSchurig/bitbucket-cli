# Auto-generated Terraform test configuration for bitbucket_workspace_members
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

variable "member" {
  type    = string
  default = "{member-uuid}"
}

provider "bitbucket" {}

data "bitbucket_workspace_members" "test" {
  workspace = var.workspace
  member = var.member
}
