package qdrant

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccDataClusterAccountsClusters(t *testing.T) {
	provider := fmt.Sprintf(`
provider "qdrant-cloud" {
  api_key = "%s"
}
`, os.Getenv("QDRANT_CLOUD_API_KEY"))

	config := provider + fmt.Sprintf(`
data "qdrant-cloud_accounts_clusters" "test" {
	account_id = "%s"
}
`, os.Getenv("QDRANT_CLOUD_ACCOUNT_ID"))

	check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet("data.qdrant-cloud_accounts_clusters.test", "clusters.#"),
	)

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"qdrant-cloud": func() (*schema.Provider, error) {
				return Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config:  config,
				Destroy: true,
				Check:   check,
			},
		},
	})
}
