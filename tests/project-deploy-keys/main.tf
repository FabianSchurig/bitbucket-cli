# Auto-generated Terraform test configuration for bitbucket_project_deploy_keys
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

variable "key_id" {
  type    = string
  default = "123"
}

provider "bitbucket" {}

data "bitbucket_project_deploy_keys" "test" {
  project_key = var.project_key
  workspace = var.workspace
  key_id = var.key_id
}

resource "bitbucket_project_deploy_keys" "test" {
  project_key = var.project_key
  workspace = var.workspace
}
