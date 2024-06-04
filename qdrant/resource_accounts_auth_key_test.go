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
data "qdrant-cloud_booking_packages" "test" {}
locals {
  account_id = "%s"
  resource_data = data.qdrant-cloud_booking_packages.test.packages
  // Filter out the free tariffs
  // TODO: Change the resource.name to resource.type when the API is updated
  free_tariffs = [
    for resource in local.resource_data : resource if resource.name == "free2"
  ]
  // Get the first free tariff
  first_free_tariff = local.free_tariffs[0]
}

resource "qdrant-cloud_accounts_cluster" "test" {
	name = "test-cluster"
	account_id = local.account_id
	cloud_region = "us-east4"
	cloud_provider = "gcp"

	configuration {
		num_nodes_max = 1
		num_nodes = 1

		node_configuration {
			package_id = local.first_free_tariff.id
		}
	}
}

resource "qdrant-cloud_accounts_auth_key" "test" {
	account_id = local.account_id
	cluster_ids = [qdrant-cloud_accounts_cluster.test.id]
}
	`, os.Getenv("QDRANT_CLOUD_ACCOUNT_ID"))

	check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet("qdrant-cloud_accounts_auth_key.test", "account_id"),
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
				ExpectNonEmptyPlan: false,
				PlanOnly:           false,
			},
			{
				Destroy: true,
				Config:  config,
			},
		},
	})
}
