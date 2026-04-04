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

variable "filename" {
  type    = string
  default = "artifact.zip"
}

variable "repo_slug" {
  type    = string
  default = "my-repo"
}

provider "bitbucket" {}

data "bitbucket_downloads" "test" {
  filename = var.filename
  repo_slug = var.repo_slug
  workspace = var.workspace
}

resource "bitbucket_downloads" "test" {
  filename = var.filename
  repo_slug = var.repo_slug
  workspace = var.workspace
}
