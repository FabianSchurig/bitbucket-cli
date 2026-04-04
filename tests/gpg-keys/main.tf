# Auto-generated Terraform test configuration for bitbucket_gpg_keys
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

variable "fingerprint" {
  type    = string
  default = "AA:BB:CC:DD"
}

provider "bitbucket" {}

data "bitbucket_gpg_keys" "test" {
  selected_user = var.selected_user
  fingerprint = var.fingerprint
}

resource "bitbucket_gpg_keys" "test" {
  selected_user = var.selected_user
}
