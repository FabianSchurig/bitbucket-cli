terraform {
  required_providers {
    bitbucket = {
      source = "FabianSchurig/bitbucket"
    }
  }
}

# Configure via environment variables:
#   BITBUCKET_USERNAME + BITBUCKET_TOKEN (Basic auth)
#   or BITBUCKET_TOKEN (OAuth2)
provider "bitbucket" {}
