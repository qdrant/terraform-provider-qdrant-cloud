package qdrant

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataAccountsAuthKeysV2(t *testing.T) {
	provider := fmt.Sprintf(`
provider "qdrant-cloud" {
  api_key = "%s"
}
	`, os.Getenv("QDRANT_CLOUD_API_KEY"))

	config := provider + fmt.Sprintf(`
data "qdrant-cloud_booking_packages" "test" {
	cloud_provider = "gcp"
	cloud_region = "us-east4"
}
locals {
  account_id = "%s"
  resource_data = data.qdrant-cloud_booking_packages.test.packages
  // Filter out the free tariffs
  free_tariffs = [
    for resource in local.resource_data : resource if resource.name == "free2"
  ]
  // Get the first free tariff
  first_free_tariff = local.free_tariffs[0]
}

resource "qdrant-cloud_accounts_cluster" "test" {
	name = "tf-acc-test-cluster-ds-auth-v2"
	account_id = local.account_id
	cloud_region = "us-east4"
	cloud_provider = "gcp"

	configuration {
		number_of_nodes = 1

		node_configuration {
			package_id = local.first_free_tariff.id
		}
	}
}

resource "qdrant-cloud_accounts_database_api_key_v2" "test" {
	name = "tf-acc-test-key-ds-v2"
	cluster_id = qdrant-cloud_accounts_cluster.test.id
	global_access_rule {
		access_type = "GLOBAL_ACCESS_RULE_ACCESS_TYPE_READ_ONLY"
	}
}

data "qdrant-cloud_accounts_database_api_keys_v2" "test" {
	cluster_id = qdrant-cloud_accounts_cluster.test.id
	depends_on = [qdrant-cloud_accounts_database_api_key_v2.test]
}
	`, os.Getenv("QDRANT_CLOUD_ACCOUNT_ID"))

	check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(
			"data.qdrant-cloud_accounts_database_api_keys_v2.test", "cluster_id",
			"qdrant-cloud_accounts_cluster.test", "id",
		),
		resource.TestCheckResourceAttrSet("data.qdrant-cloud_accounts_database_api_keys_v2.test", "keys.#"),
		resource.TestCheckResourceAttrSet("data.qdrant-cloud_accounts_database_api_keys_v2.test", "keys.0.id"),
	)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{Config: config, Check: check},
		},
	})
}
