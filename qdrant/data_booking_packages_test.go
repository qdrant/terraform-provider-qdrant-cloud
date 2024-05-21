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
provider "qdrant" {
  api_key = "%s"
}
`, os.Getenv("QDRANT_CLOUD_API_KEY"))

	config := provider + `
data "qdrant_booking_packages" "test" {
}
`

	check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet("data.qdrant_booking_packages.test", "packages.#"),
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
