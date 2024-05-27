package qdrant

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccDataAccountsAuthKeys(t *testing.T) {
	provider := fmt.Sprintf(`
provider "qdrant" {
  api_key = "%s"
}
	`, os.Getenv("QDRANT_CLOUD_API_KEY"))

	config := provider + `
data "qdrant_accounts_auth_keys" "test" {
	account_id = "c3e03ee1-b79b-443d-80f0-8eb8a2671978"
}
	`

	check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet("data.qdrant_accounts_auth_keys.test", "account_id"),
		resource.TestCheckResourceAttrSet("data.qdrant_accounts_auth_keys.test", "keys.#"),
	)

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"qdrant": func() (*schema.Provider, error) {
				return Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				Check:  check,
			},
		},
	})
}
