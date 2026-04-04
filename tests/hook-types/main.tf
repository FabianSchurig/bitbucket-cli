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

variable "subject_type" {
  type    = string
  default = "repository"
}

provider "bitbucket" {}

data "bitbucket_hook_types" "test" {
  subject_type = var.subject_type
}
