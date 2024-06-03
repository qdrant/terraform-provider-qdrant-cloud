package qdrant

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Test the provider configuration with variables set
func TestProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"qdrant-cloud": func() (*schema.Provider, error) {
				return Provider(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAPIKeyConfigBasic(),
			},
		},
	})
}
func testAccCheckAPIKeyConfigBasic() string {
	apiKey := os.Getenv("QDRANT_CLOUD_API_KEY")
	apiURL := os.Getenv("QDRANT_CLOUD_API_URL")

	return fmt.Sprintf(`
provider "qdrant-cloud" {
  alias = "qdrant_cloud"
  api_key = "%s"
  api_url = "%s"
}



`, apiKey, apiURL)
}

// Custom check function to validate that the provider's client is configured
func testAccCheckProviderConfigured(providerName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[providerName]
		if !ok {
			return fmt.Errorf("provider %q is not configured", providerName)
		}
		return nil
	}
}
