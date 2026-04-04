# Auto-generated Terraform test configuration for bitbucket_hook_types
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

provider "bitbucket" {}

data "bitbucket_hook_types" "test" {
}
