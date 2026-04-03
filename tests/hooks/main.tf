# Auto-generated Terraform test configuration for bitbucket_hooks
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

variable "uid" {
  type    = string
  default = "webhook-uuid"
}

provider "bitbucket" {}

data "bitbucket_hooks" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  uid = var.uid
}

resource "bitbucket_hooks" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  uid = var.uid
}
