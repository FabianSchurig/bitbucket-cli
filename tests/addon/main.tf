# Auto-generated Terraform test configuration for bitbucket_addon
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

variable "addon_key" {
  type    = string
  default = "example-value"
}

provider "bitbucket" {}

data "bitbucket_addon" "test" {
  addon_key = var.addon_key
}
