# Auto-generated Terraform test configuration for bitbucket_pipelines
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

variable "pipeline_uuid" {
  type    = string
  default = "pipeline-uuid"
}

provider "bitbucket" {}

data "bitbucket_pipelines" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  pipeline_uuid = var.pipeline_uuid
}

resource "bitbucket_pipelines" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
}
