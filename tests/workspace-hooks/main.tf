# Auto-generated Terraform test configuration for bitbucket_workspace_hooks
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

variable "uid" {
  type    = string
  default = "webhook-uuid"
}

provider "bitbucket" {}

data "bitbucket_workspace_hooks" "test" {
  uid = var.uid
  workspace = var.workspace
}

resource "bitbucket_workspace_hooks" "test" {
  uid = var.uid
  workspace = var.workspace
}
