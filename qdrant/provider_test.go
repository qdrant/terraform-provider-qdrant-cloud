package qdrant

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var providerFactories = map[string]func() (*schema.Provider, error){
	"qdrant-cloud": func() (*schema.Provider, error) { //nolint:unparam // Interface is defined by TF, so we cannot remove error
		return Provider(), nil
	},
}

// Test the provider configuration with variables set.
func TestProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
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
