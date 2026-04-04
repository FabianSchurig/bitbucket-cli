# Auto-generated Terraform test configuration for bitbucket_repo_group_permissions
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

variable "repo_slug" {
  type    = string
  default = "my-repo"
}

provider "bitbucket" {}

data "bitbucket_repo_group_permissions" "test" {
  group_slug = var.group_slug
  repo_slug = var.repo_slug
  workspace = var.workspace
}
