// bb-tf is a Terraform provider for Bitbucket Cloud.
//
// It exposes all Bitbucket API operations as Terraform resources and data sources,
// auto-generated from the same OpenAPI schema that drives the CLI and MCP server.
//
// Install:
//
//	go install github.com/FabianSchurig/bitbucket-cli/cmd/bb-tf@latest
//
// Authentication:
//
//	App password (most common):
//	  export BITBUCKET_USERNAME=myuser
//	  export BITBUCKET_APP_PASSWORD=ATBBxxxxxxxx
//
//	OAuth2 access token:
//	  export BITBUCKET_TOKEN=<token>
//
// Usage in Terraform:
//
//	terraform {
//	  required_providers {
//	    bitbucket = {
//	      source  = "FabianSchurig/bitbucket"
//	    }
//	  }
//	}
//
//	provider "bitbucket" {
//	  username     = "myuser"
//	  app_password = "ATBBxxxxxxxx"
//	}
package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/FabianSchurig/bitbucket-cli/internal/tfprovider"
)

// Set via ldflags at build time (see goreleaser.yaml).
var version = "dev"

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/FabianSchurig/bitbucket",
		Debug:   debug,
	}

	if err := providerserver.Serve(context.Background(), tfprovider.New(version), opts); err != nil {
		log.Fatal(err.Error())
	}
}
