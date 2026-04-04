# Auto-generated Terraform test configuration for bitbucket_ssh_keys
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

variable "key_id" {
  type    = string
  default = "123"
}

variable "selected_user" {
  type    = string
  default = "jdoe"
}

provider "bitbucket" {}

data "bitbucket_ssh_keys" "test" {
  key_id = var.key_id
  selected_user = var.selected_user
}

resource "bitbucket_ssh_keys" "test" {
  key_id = var.key_id
  selected_user = var.selected_user
}
