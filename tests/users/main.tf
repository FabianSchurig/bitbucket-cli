# Auto-generated Terraform test configuration for bitbucket_users
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

variable "selected_user" {
  type    = string
  default = "jdoe"
}

provider "bitbucket" {}

data "bitbucket_users" "test" {
  selected_user = var.selected_user
}
