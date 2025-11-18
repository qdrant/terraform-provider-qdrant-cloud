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

	resourceName := "qdrant-cloud_accounts_hybrid_cloud_environment.test"

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
					resource.TestCheckResourceAttr(resourceName, "name", "tf-acc-test-hc-env"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.namespace", "qdrant-hc-tf-acc"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "account_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
					resource.TestCheckResourceAttrSet(resourceName, "last_modified_at"),
					// The user creating the resource via API key should have an email.
					resource.TestCheckResourceAttrSet(resourceName, "created_by_email"),
					// Bootstrap commands are generated on create, so this should be true.
					resource.TestCheckResourceAttr(resourceName, "bootstrap_commands_generated", "true"),
					// knob should be bumped to 1 on create (unless user explicitly set -1, which we don't here)
					resource.TestCheckResourceAttr(resourceName, "bootstrap_commands_version", "1"),
					// commands should be generated when version > 0
					testAccCheckListNonEmpty(resourceName, "bootstrap_commands"),
				),
			},
			{
				Config:  config,
				Destroy: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tf-acc-test-hc-env"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.namespace", "qdrant-hc-tf-acc"),
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
					// Ephemeral / local-only:
					"bootstrap_commands",
					"bootstrap_commands_version",
				},
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[resourceName]
					if !ok {
						return "", fmt.Errorf("resource not found in state: %s", resourceName)
					}
					return rs.Primary.ID, nil
				},
			},
			// Ensure read works with same config, then destroy
			{Config: config, Destroy: true},
		},
	})
}

func TestAccResourceAccountsHybridCloudEnvironment_UpdateVersion(t *testing.T) {
	provider := fmt.Sprintf(`
provider "qdrant-cloud" {
  api_key = "%s"
}
`, os.Getenv("QDRANT_CLOUD_API_KEY"))

	// Initial config with version = 0 (default)
	config := provider + fmt.Sprintf(`
resource "qdrant-cloud_accounts_hybrid_cloud_environment" "test" {
  name       = "tf-acc-test-hc-env-update"
  account_id = "%s"

  configuration {
    namespace = "qdrant-hc-tf-acc-update"
  }
}
`, os.Getenv("QDRANT_CLOUD_ACCOUNT_ID"))

	// Updated config bumps version knob
	configBump := provider + fmt.Sprintf(`
resource "qdrant-cloud_accounts_hybrid_cloud_environment" "test" {
  name       = "tf-acc-test-hc-env-update"
  account_id = "%s"

  configuration {
    namespace = "qdrant-hc-tf-acc-update"
  }

  bootstrap_commands_version = 1
}
`, os.Getenv("QDRANT_CLOUD_ACCOUNT_ID"))

	const resourceName = "qdrant-cloud_accounts_hybrid_cloud_environment.test"

	resource.Test(t, resource.TestCase{
		//nolint:unparam
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"qdrant-cloud": func() (*schema.Provider, error) { return Provider(), nil },
		},
		Steps: []resource.TestStep{
			// Step 1: Create
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "tf-acc-test-hc-env-update"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.namespace", "qdrant-hc-tf-acc-update"),
				),
			},
			// Step 2: Update version knob
			{
				Config: configBump,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "bootstrap_commands_version", "1"),
					// bootstrap_commands are sensitive; just check non-empty
					testCheckHasListAttr(resourceName, "bootstrap_commands"),
				),
			},
			// Step 3: Destroy
			{
				Config:  configBump,
				Destroy: true,
			},
		},
	})
}
