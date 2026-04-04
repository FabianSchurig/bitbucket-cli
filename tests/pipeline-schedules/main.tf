# Auto-generated Terraform test configuration for bitbucket_pipeline_schedules
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

variable "schedule_uuid" {
  type    = string
  default = "{schedule-uuid}"
}

provider "bitbucket" {}

data "bitbucket_pipeline_schedules" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  schedule_uuid = var.schedule_uuid
}

resource "bitbucket_pipeline_schedules" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  schedule_uuid = var.schedule_uuid
}
