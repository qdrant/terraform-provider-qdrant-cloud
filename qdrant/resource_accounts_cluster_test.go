package qdrant

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceClusterCreate(t *testing.T) {
	provider := fmt.Sprintf(`
provider "qdrant-cloud" {
  api_key = "%s"
}
	`, os.Getenv("QDRANT_CLOUD_API_KEY"))

	config := provider + fmt.Sprintf(`
data "qdrant-cloud_booking_packages" "test" {}
locals {
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
	account_id = "%s"
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

output "cluster_name" {
	value = qdrant-cloud_accounts_cluster.test.name
}

`, os.Getenv("QDRANT_CLOUD_ACCOUNT_ID"))

	config_update := provider + fmt.Sprintf(`
data "qdrant-cloud_booking_packages" "test" {}
locals {
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
	name           = "test-cluster"
	account_id     = "%s"
	cloud_region   = "us-east4"
	cloud_provider = "gcp"

	configuration {
		num_nodes_max = 2   // Update the number of nodes to 2
		num_nodes     = 2   // New number of nodes for the cluster

		node_configuration {
			package_id = local.first_free_tariff.id
		}
	}
}

output "cluster_name" {
	value = qdrant-cloud_accounts_cluster.test.name
}

output "cluster_cfg_num_nodes" {
	value = qdrant-cloud_accounts_cluster.test.configuration[0].num_nodes
}

`, os.Getenv("QDRANT_CLOUD_ACCOUNT_ID"))

	t.Run("creates a cluster", func(t *testing.T) {
		testCase := func(t *testing.T, mode string) {
			resource.Test(t, resource.TestCase{
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"qdrant-cloud": func() (*schema.Provider, error) {
						return Provider(), nil
					},
				},
				Steps: []resource.TestStep{
					{
						Config:  config,
						Destroy: false,
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckOutput("cluster_name", "test-cluster"),
						),
					},
					{
						Config:  config_update,
						Destroy: true,
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckOutput("cluster_cfg_num_nodes", "2"),
						),
					},
				},
			})
		}
		testCase(t, "apply")
	})
}
