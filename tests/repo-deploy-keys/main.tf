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

variable "key_id" {
  type    = string
  default = "123"
}

variable "repo_slug" {
  type    = string
  default = "my-repo"
}

provider "bitbucket" {}

data "bitbucket_repo_deploy_keys" "test" {
  key_id = var.key_id
  repo_slug = var.repo_slug
  workspace = var.workspace
}

resource "bitbucket_repo_deploy_keys" "test" {
  key_id = var.key_id
  repo_slug = var.repo_slug
  workspace = var.workspace
}
