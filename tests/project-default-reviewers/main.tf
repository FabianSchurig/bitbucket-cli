# Auto-generated Terraform test configuration for bitbucket_project_default_reviewers
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

variable "project_key" {
  type    = string
  default = "PROJ"
}

variable "selected_user" {
  type    = string
  default = "jdoe"
}

provider "bitbucket" {}

data "bitbucket_project_default_reviewers" "test" {
  workspace = var.workspace
  project_key = var.project_key
  selected_user = var.selected_user
}

resource "bitbucket_project_default_reviewers" "test" {
  workspace = var.workspace
  project_key = var.project_key
  selected_user = var.selected_user
}
