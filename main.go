package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/qdrant/terraform-provider-qdrant-cloud/v1/qdrant"
)

// Examples folder formatting
//go:generate terraform fmt -recursive ./examples/

// Terraform plugin tool for documentation generation
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return qdrant.Provider()
		},
	})
}
