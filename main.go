package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"terraform-provider-qdrant-cloud/v1/qdrant"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: qdrant.Provider})
}
