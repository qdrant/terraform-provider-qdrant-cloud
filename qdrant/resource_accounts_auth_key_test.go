package qdrant

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccResourceAccountsAuthKeys(t *testing.T) {
	provider := fmt.Sprintf(`
provider "qdrant-cloud" {
  api_key = "%s"
}
	`, os.Getenv("QDRANT_CLOUD_API_KEY"))

	config := provider + fmt.Sprintf(`
resource "qdrant-cloud_accounts_auth_key" "test" {
	account_id = "%s"
	cluster_ids = ["9095da42-f4dc-4f03-93be-dd346ddc3302"]
}
	`, os.Getenv("QDRANT_CLOUD_ACCOUNT_ID"))

	check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet("qdrant_accounts_auth_keys.test", "account_id"),
		resource.TestCheckResourceAttrSet("qdrant_accounts_auth_keys.test", "keys.#"),
	)

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"qdrant-cloud": func() (*schema.Provider, error) {
				return Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config:             config,
				Check:              check,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
