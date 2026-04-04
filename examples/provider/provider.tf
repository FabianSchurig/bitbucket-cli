terraform {
  required_providers {
    bitbucket = {
      source = "FabianSchurig/bitbucket"
    }
  }
}

# Configure via environment variables:
#   BITBUCKET_USERNAME (email) + BITBUCKET_TOKEN (Atlassian API token)
#   or BITBUCKET_TOKEN alone (workspace/repository access token)
provider "bitbucket" {}
