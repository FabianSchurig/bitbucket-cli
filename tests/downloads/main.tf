# Auto-generated Terraform test configuration for bitbucket_downloads
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

variable "filename" {
  type    = string
  default = "artifact.zip"
}

provider "bitbucket" {}

data "bitbucket_downloads" "test" {
  repo_slug = var.repo_slug
  workspace = var.workspace
  filename = var.filename
}

resource "bitbucket_downloads" "test" {
  repo_slug = var.repo_slug
  workspace = var.workspace
}
