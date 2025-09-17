package qdrant

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceAccountsHybridCloudEnvironment_Basic(t *testing.T) {
	provider := fmt.Sprintf(`
provider "qdrant-cloud" {
  api_key = "%s"
}
`, os.Getenv("QDRANT_CLOUD_API_KEY"))

	config := provider + fmt.Sprintf(`
resource "qdrant-cloud_accounts_hybrid_cloud_environment" "test" {
  name       = "tf-acc-test-hc-env"
  account_id = "%s"

  configuration {
    namespace = "qdrant-hc-tf-acc"
  }
}
`, os.Getenv("QDRANT_CLOUD_ACCOUNT_ID"))

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			//nolint:unparam
			"qdrant-cloud": func() (*schema.Provider, error) {
				return Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("qdrant-cloud_accounts_hybrid_cloud_environment.test", "name", "tf-acc-test-hc-env"),
					resource.TestCheckResourceAttr("qdrant-cloud_accounts_hybrid_cloud_environment.test", "configuration.0.namespace", "qdrant-hc-tf-acc"),
					resource.TestCheckResourceAttrSet("qdrant-cloud_accounts_hybrid_cloud_environment.test", "id"),
					resource.TestCheckResourceAttrSet("qdrant-cloud_accounts_hybrid_cloud_environment.test", "account_id"),
					resource.TestCheckResourceAttrSet("qdrant-cloud_accounts_hybrid_cloud_environment.test", "created_at"),
					resource.TestCheckResourceAttrSet("qdrant-cloud_accounts_hybrid_cloud_environment.test", "last_modified_at"),
				),
			},
			{
				Config:  config,
				Destroy: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("qdrant-cloud_accounts_hybrid_cloud_environment.test", "name", "tf-acc-test-hc-env"),
					resource.TestCheckResourceAttr("qdrant-cloud_accounts_hybrid_cloud_environment.test", "configuration.0.namespace", "qdrant-hc-tf-acc"),
				),
			},
		},
	})
}

func TestAccResourceAccountsHybridCloudEnvironment_Import(t *testing.T) {
	provider := fmt.Sprintf(`
provider "qdrant-cloud" {
  api_key = "%s"
}
`, os.Getenv("QDRANT_CLOUD_API_KEY"))

	config := provider + fmt.Sprintf(`
resource "qdrant-cloud_accounts_hybrid_cloud_environment" "test" {
  name       = "tf-acc-test-hc-env-import"
  account_id = "%s"

  configuration {
    namespace = "qdrant-hc-tf-acc-import"
  }
}
`, os.Getenv("QDRANT_CLOUD_ACCOUNT_ID"))

	resourceName := "qdrant-cloud_accounts_hybrid_cloud_environment.test"

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			//nolint:unparam
			"qdrant-cloud": func() (*schema.Provider, error) {
				return Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{Config: config},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"bootstrap_commands", // ephemerally regenerated; ignore entire list
				},
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found in state: %s", resourceName)
					}
					return rs.Primary.ID, nil
				},
			},
			// Ensures we can still read it with the same config, then destroy
			{Config: config, Destroy: true},
		},
	})
}
