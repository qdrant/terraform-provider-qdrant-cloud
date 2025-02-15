package qdrant

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccDataBookingPackages(t *testing.T) {
	provider := fmt.Sprintf(`
provider "qdrant-cloud" {
  api_key = "%s"
}
`, os.Getenv("QDRANT_CLOUD_API_KEY"))

	config := provider + `
data "qdrant-cloud_booking_packages" "test" {
	  cloud_provider = "gcp"
	  cloud_region = "europe-west3"
}

locals {
  resource_data = data.qdrant-cloud_booking_packages.test.packages

  // Filter out the free tariffs
  free_tariffs = [
    // TODO: Change the resource.name to resource.type when the API is updated
    for resource in local.resource_data : resource if resource.name == "free2"
  ]

  // Get the first free tariff
  first_free_tariff = local.free_tariffs[0]
}

output "first_free_tariff" {
	  value = local.first_free_tariff.name
}
`

	check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet("data.qdrant-cloud_booking_packages.test", "packages.#"),
		resource.TestCheckOutput("first_free_tariff", "free2"),
	)

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			//nolint:unparam // Ignoring unparam as we know error will always be nil
			"qdrant-cloud": func() (*schema.Provider, error) {
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
