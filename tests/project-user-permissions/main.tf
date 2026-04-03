# Auto-generated Terraform test configuration for bitbucket_project_user_permissions
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

variable "selected_user_id" {
  type    = string
  default = "{user-uuid}"
}

provider "bitbucket" {}

data "bitbucket_project_user_permissions" "test" {
  workspace = var.workspace
  project_key = var.project_key
  selected_user_id = var.selected_user_id
}
