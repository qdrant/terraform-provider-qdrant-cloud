package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"terraform-provider-qdrant-cloud/v1/qdrant"
)

// Examples folder formatting
//go:generate terraform fmt -recursive ./examples/

// Terraform plugin tool for documentation generation
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: qdrant.Provider})
}
