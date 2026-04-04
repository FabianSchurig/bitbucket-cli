# Auto-generated Terraform test configuration for bitbucket_pipeline_variables
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

variable "variable_uuid" {
  type    = string
  default = "{variable-uuid}"
}

provider "bitbucket" {}

data "bitbucket_pipeline_variables" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  variable_uuid = var.variable_uuid
}

resource "bitbucket_pipeline_variables" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  variable_uuid = var.variable_uuid
}
