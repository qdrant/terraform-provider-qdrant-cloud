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
data "qdrant-cloud_booking_packages" "test" {
    	cloud_provider = "gcp"
		cloud_region = "us-east4"
}
locals {
  resource_data = data.qdrant-cloud_booking_packages.test.packages
  // Filter out the gpx1 
  gpx1_packages = [
    for resource in local.resource_data : resource if resource.name == "gpx1"
  ]
  // Get the first gpx1 package
  gpx1_package = local.gpx1_packages[0]
}

resource "qdrant-cloud_accounts_cluster" "test" {
	name = "test-cluster"
	account_id = "%s"
	cloud_region = "us-east4"
	cloud_provider = "gcp"

	configuration {
		number_of_nodes = 1

		node_configuration {
			package_id = local.gpx1_package.id
		}
	}
}

output "cluster_name" {
	value = qdrant-cloud_accounts_cluster.test.name
}

`, os.Getenv("QDRANT_CLOUD_ACCOUNT_ID"))

	config_update := provider + fmt.Sprintf(`
data "qdrant-cloud_booking_packages" "test" {
		cloud_provider = "gcp"
		cloud_region = "us-east4"
}
locals {
  resource_data = data.qdrant-cloud_booking_packages.test.packages
  // Filter out the gpx1 
  gpx1_packages = [
    for resource in local.resource_data : resource if resource.name == "gpx1"
  ]
  // Get the first gpx1 package
  gpx1_package = local.gpx1_packages[0]
}

resource "qdrant-cloud_accounts_cluster" "test" {
	name           = "test-cluster"
	account_id     = "%s"
	cloud_region   = "us-east4"
	cloud_provider = "gcp"

	configuration {
		number_of_nodes = 2   // Update the number of nodes to 2

		node_configuration {
			package_id = local.gpx1_package.id
		}
	}
}

output "cluster_name" {
	value = qdrant-cloud_accounts_cluster.test.name
}

output "cluster_cfg_num_nodes" {
	value = qdrant-cloud_accounts_cluster.test.configuration[0].number_of_nodes
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

func TestResourceClusterCreateWithExtraDisk(t *testing.T) {
	provider := fmt.Sprintf(`
provider "qdrant-cloud" {
  api_key = "%s"
}
	`, os.Getenv("QDRANT_CLOUD_API_KEY"))

	config := provider + fmt.Sprintf(`
data "qdrant-cloud_booking_packages" "test" {
  cloud_provider = "gcp"
  cloud_region   = "us-east4"
}
locals {
  resource_data = data.qdrant-cloud_booking_packages.test.packages
  // Filter out the gpx1
  gpx1_packages = [
    for resource in local.resource_data : resource if resource.name == "gpx1"
  ]
  // Get the first gpx1 package
  gpx1_package = local.gpx1_packages[0]
}

resource "qdrant-cloud_accounts_cluster" "test" {
  name           = "test-cluster-with-extra-disk"
  account_id     = "%s"
  cloud_region   = "us-east4"
  cloud_provider = "gcp"

  configuration {
    number_of_nodes = 1

    node_configuration {
      package_id = local.gpx1_package.id
      resource_configurations {
        amount        = 8
        resource_type = "disk"
        resource_unit = "Gi"
      }
    }
  }
}

output "cluster_name" {
  value = qdrant-cloud_accounts_cluster.test.name
}

output "extra_disk_amount" {
  value = qdrant-cloud_accounts_cluster.test.configuration[0].node_configuration[0].resource_configurations[0].amount
}

`, os.Getenv("QDRANT_CLOUD_ACCOUNT_ID"))

	t.Run("creates a cluster with extra disk", func(t *testing.T) {
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
							resource.TestCheckOutput("cluster_name", "test-cluster-with-extra-disk"),
							resource.TestCheckOutput("extra_disk_amount", "8"),
						),
					},
				},
			})
		}
		testCase(t, "apply")
	})
}
