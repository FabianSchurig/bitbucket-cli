# Auto-generated Terraform test configuration for bitbucket_snippets
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

variable "encoded_id" {
  type    = string
  default = "snippet-id"
}

provider "bitbucket" {}

data "bitbucket_snippets" "test" {
  encoded_id = var.encoded_id
  workspace = var.workspace
}

resource "bitbucket_snippets" "test" {
  encoded_id = var.encoded_id
  workspace = var.workspace
}
