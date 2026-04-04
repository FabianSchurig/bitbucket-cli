# Auto-generated Terraform test configuration for bitbucket_pipeline_known_hosts
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

variable "known_host_uuid" {
  type    = string
  default = "{known-host-uuid}"
}

provider "bitbucket" {}

data "bitbucket_pipeline_known_hosts" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  known_host_uuid = var.known_host_uuid
}

resource "bitbucket_pipeline_known_hosts" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
}
