# Auto-generated Terraform test configuration for bitbucket_project_group_permissions
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

variable "group_slug" {
  type    = string
  default = "developers"
}

variable "project_key" {
  type    = string
  default = "PROJ"
}

provider "bitbucket" {}

data "bitbucket_project_group_permissions" "test" {
  group_slug = var.group_slug
  project_key = var.project_key
  workspace = var.workspace
}
