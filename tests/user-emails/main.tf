# Auto-generated Terraform test configuration for bitbucket_user_emails
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

variable "email" {
  type    = string
  default = "user@example.com"
}

provider "bitbucket" {}

data "bitbucket_user_emails" "test" {
  email = var.email
}
