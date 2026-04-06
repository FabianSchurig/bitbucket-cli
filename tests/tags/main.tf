# Auto-generated Terraform test configuration for bitbucket_tags
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

variable "name" {
  type    = string
  default = "main"
}

provider "bitbucket" {}

data "bitbucket_tags" "test" {
  repo_slug = var.repo_slug
  workspace = var.workspace
  name = var.name
}

resource "bitbucket_tags" "test" {
  repo_slug = var.repo_slug
  workspace = var.workspace
  type = "webhook"
}
