# Auto-generated Terraform test configuration for bitbucket_repo_user_permissions
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

variable "selected_user_id" {
  type    = string
  default = "{user-uuid}"
}

provider "bitbucket" {}

data "bitbucket_repo_user_permissions" "test" {
  repo_slug = var.repo_slug
  selected_user_id = var.selected_user_id
  workspace = var.workspace
}
