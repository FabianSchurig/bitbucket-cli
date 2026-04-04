# Auto-generated Terraform test configuration for bitbucket_annotations
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

variable "commit" {
  type    = string
  default = "abc123def"
}

variable "report_id" {
  type    = string
  default = "report-uuid"
}

variable "annotation_id" {
  type    = string
  default = "{annotation-id}"
}

provider "bitbucket" {}

data "bitbucket_annotations" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  commit = var.commit
  report_id = var.report_id
  annotation_id = var.annotation_id
}

resource "bitbucket_annotations" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  commit = var.commit
  report_id = var.report_id
  annotation_id = var.annotation_id
}
