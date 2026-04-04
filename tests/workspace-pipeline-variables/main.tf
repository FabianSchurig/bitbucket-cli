# Auto-generated Terraform test configuration for bitbucket_workspace_pipeline_variables
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

variable "variable_uuid" {
  type    = string
  default = "{variable-uuid}"
}

provider "bitbucket" {}

data "bitbucket_workspace_pipeline_variables" "test" {
  workspace = var.workspace
  variable_uuid = var.variable_uuid
}

resource "bitbucket_workspace_pipeline_variables" "test" {
  workspace = var.workspace
  variable_uuid = var.variable_uuid
}
