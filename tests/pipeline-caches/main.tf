# Auto-generated Terraform test configuration for bitbucket_pipeline_caches
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

variable "cache_uuid" {
  type    = string
  default = "{cache-uuid}"
}

provider "bitbucket" {}

data "bitbucket_pipeline_caches" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  cache_uuid = var.cache_uuid
}
