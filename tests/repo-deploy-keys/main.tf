# Auto-generated Terraform test configuration for bitbucket_repo_deploy_keys
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

variable "key_id" {
  type    = string
  default = "123"
}

provider "bitbucket" {}

data "bitbucket_repo_deploy_keys" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  key_id = var.key_id
}

resource "bitbucket_repo_deploy_keys" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  key_id = var.key_id
}
