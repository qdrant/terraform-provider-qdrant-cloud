package main

import (
	"qdrant-terraform-automation/qdrantcloud"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

{{/* Every TF provider needs a main function. */}}
func main() {
	opts := &plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return qdrantcloud.Provider()
		},
	}

	plugin.Serve(opts)
}
