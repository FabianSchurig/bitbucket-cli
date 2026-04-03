# Auto-generated Terraform test configuration for bitbucket_branching_model
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

data "bitbucket_branching_model" "test" {
  repo_slug = var.repo_slug
  workspace = var.workspace
}
