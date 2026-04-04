# Auto-generated Terraform test configuration for bitbucket_repo_settings
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

provider "bitbucket" {}

data "bitbucket_repo_settings" "test" {
  repo_slug = var.repo_slug
  workspace = var.workspace
}
