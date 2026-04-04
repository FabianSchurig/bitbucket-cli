# Auto-generated Terraform test configuration for bitbucket_deployments
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

variable "environment_uuid" {
  type    = string
  default = "env-uuid"
}

provider "bitbucket" {}

data "bitbucket_deployments" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  environment_uuid = var.environment_uuid
}

resource "bitbucket_deployments" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  environment_uuid = var.environment_uuid
}
