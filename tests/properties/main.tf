# Auto-generated Terraform test configuration for bitbucket_properties
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

variable "app_key" {
  type    = string
  default = "my-app"
}

variable "property_name" {
  type    = string
  default = "my-property"
}

provider "bitbucket" {}

data "bitbucket_properties" "test" {
  workspace = var.workspace
  repo_slug = var.repo_slug
  app_key = var.app_key
  property_name = var.property_name
}
