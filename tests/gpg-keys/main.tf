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

variable "fingerprint" {
  type    = string
  default = "AA:BB:CC:DD"
}

variable "selected_user" {
  type    = string
  default = "jdoe"
}

provider "bitbucket" {}

data "bitbucket_gpg_keys" "test" {
  fingerprint = var.fingerprint
  selected_user = var.selected_user
}

resource "bitbucket_gpg_keys" "test" {
  fingerprint = var.fingerprint
  selected_user = var.selected_user
}
