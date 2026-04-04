# Auto-generated Terraform test configuration for bitbucket_workspace_runners
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

variable "runner_uuid" {
  type    = string
  default = "{runner-uuid}"
}

provider "bitbucket" {}

data "bitbucket_workspace_runners" "test" {
  workspace = var.workspace
  runner_uuid = var.runner_uuid
}

resource "bitbucket_workspace_runners" "test" {
  workspace = var.workspace
}
