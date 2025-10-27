package qdrant

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccResourceClusterCreate(t *testing.T) {
	apiKey := os.Getenv("QDRANT_CLOUD_API_KEY")
	accountID := os.Getenv("QDRANT_CLOUD_ACCOUNT_ID")
	cloudProvider := getEnvDefault("QDRANT_CLOUD_CLOUD_PROVIDER", "gcp")
	cloudRegion := getEnvDefault("QDRANT_CLOUD_REGION", "europe-west3")

	provider := fmt.Sprintf(`
provider "qdrant-cloud" {
  api_key = "%s"
}
	`, apiKey)

	config := provider + fmt.Sprintf(`
data "qdrant-cloud_booking_packages" "test" {
	cloud_provider = "%s"
	cloud_region   = "%s"
}
locals {
  resource_data = data.qdrant-cloud_booking_packages.test.packages
  // Filter on gpx1 
  gpx1_packages = [
    for resource in local.resource_data : resource if resource.name == "gpx1"
  ]
  // Get the first gpx1 package
  gpx1_package = local.gpx1_packages[0]
}

resource "qdrant-cloud_accounts_cluster" "test" {
	name           = "test-cluster"
	account_id     = "%s"
	cloud_region   = "%s"
	cloud_provider = "%s"

	lifecycle {
	  ignore_changes = [
	    configuration[0].gpu_type,
	    configuration[0].rebalance_strategy,
	    configuration[0].restart_policy,
	    configuration[0].service_type,
	    configuration[0].allowed_ip_source_ranges,
	    configuration[0].database_configuration,
	  ]
	}

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
`, cloudProvider, cloudRegion, accountID, cloudRegion, cloudProvider)

	config_update := provider + fmt.Sprintf(`
data "qdrant-cloud_booking_packages" "test" {
	cloud_provider = "%s"
	cloud_region   = "%s"
}
locals {
  resource_data = data.qdrant-cloud_booking_packages.test.packages
  // Filter on gpx1 
  gpx1_packages = [
    for resource in local.resource_data : resource if resource.name == "gpx1"
  ]
  // Get the first gpx1 package
  gpx1_package = local.gpx1_packages[0]
}

resource "qdrant-cloud_accounts_cluster" "test" {
	name           = "test-cluster"
	account_id     = "%s"
	cloud_region   = "%s"
	cloud_provider = "%s"

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
`, cloudProvider, cloudRegion, accountID, cloudRegion, cloudProvider)

	t.Run("creates a cluster", func(t *testing.T) {
		testCase := func(t *testing.T, _ string) {
			resource.Test(t, resource.TestCase{
				ProviderFactories: map[string]func() (*schema.Provider, error){
					//nolint:unparam // Ignoring unparam as we know error will always be nil
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

func TestAccResourceClusterCreateWithExtraDisk(t *testing.T) {
	apiKey := os.Getenv("QDRANT_CLOUD_API_KEY")
	accountID := os.Getenv("QDRANT_CLOUD_ACCOUNT_ID")
	cloudProvider := getEnvDefault("QDRANT_CLOUD_CLOUD_PROVIDER", "gcp")
	cloudRegion := getEnvDefault("QDRANT_CLOUD_REGION", "europe-west3")

	provider := fmt.Sprintf(`
provider "qdrant-cloud" {
  api_key = "%s"
}
	`, apiKey)

	config := provider + fmt.Sprintf(`
data "qdrant-cloud_booking_packages" "test" {
  cloud_provider = "%s"
  cloud_region   = "%s"
}
locals {
  resource_data = data.qdrant-cloud_booking_packages.test.packages
  // Filter on gpx1
  gpx1_packages = [
    for resource in local.resource_data : resource if resource.name == "gpx1"
  ]
  // Get the first gpx1 package
  gpx1_package = local.gpx1_packages[0]
}

resource "qdrant-cloud_accounts_cluster" "test" {
  name           = "test-cluster-with-extra-disk"
  account_id     = "%s"
  cloud_region   = "%s"
  cloud_provider = "%s"

  lifecycle {
    ignore_changes = [
      configuration[0].gpu_type,
      configuration[0].rebalance_strategy,
      configuration[0].restart_policy,
      configuration[0].service_type,
      configuration[0].allowed_ip_source_ranges,
      configuration[0].database_configuration,
    ]
  }

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
`, cloudProvider, cloudRegion, accountID, cloudRegion, cloudProvider)

	t.Run("creates a cluster with extra disk", func(t *testing.T) {
		testCase := func(t *testing.T, _ string) {
			resource.Test(t, resource.TestCase{
				ProviderFactories: map[string]func() (*schema.Provider, error){
					//nolint:unparam // Ignoring unparam as we know error will always be nil
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

func TestAccResourceClusterDeleteWithoutBackups(t *testing.T) {
	apiKey := os.Getenv("QDRANT_CLOUD_API_KEY")
	accountID := os.Getenv("QDRANT_CLOUD_ACCOUNT_ID")
	cloudProvider := getEnvDefault("QDRANT_CLOUD_CLOUD_PROVIDER", "gcp")
	cloudRegion := getEnvDefault("QDRANT_CLOUD_REGION", "europe-west3")

	provider := fmt.Sprintf(`
provider "qdrant-cloud" {
  api_key = "%s"
}
	`, apiKey)

	config := provider + fmt.Sprintf(`
data "qdrant-cloud_booking_packages" "test" {
  cloud_provider = "%s"
  cloud_region   = "%s"
}
locals {
  resource_data = data.qdrant-cloud_booking_packages.test.packages
  // Filter on gpx1
  gpx1_packages = [
    for resource in local.resource_data : resource if resource.name == "gpx1"
  ]
  // Get the first gpx1 package
  gpx1_package = local.gpx1_packages[0]
}

resource "qdrant-cloud_accounts_cluster" "test" {
  name                      = "test-cluster-delete-no-backups"
  account_id                = "%s"
  cloud_region              = "%s"
  cloud_provider            = "%s"
  delete_backups_on_destroy = false

  configuration {
    number_of_nodes = 1

    node_configuration {
      package_id = local.gpx1_package.id
    }
  }
}

output "cluster_name_del" {
  value = qdrant-cloud_accounts_cluster.test.name
}
`, cloudProvider, cloudRegion, accountID, cloudRegion, cloudProvider)

	t.Run("creates and deletes a cluster without deleting backups", func(t *testing.T) {
		resource.Test(t, resource.TestCase{
			ProviderFactories: map[string]func() (*schema.Provider, error){
				//nolint:unparam // Ignoring unparam as we know error will always be nil
				"qdrant-cloud": func() (*schema.Provider, error) { return Provider(), nil },
			},
			Steps: []resource.TestStep{
				{
					Config:  config,
					Destroy: true, // Destroy will be called, testing delete_backups_on_destroy
					Check:   resource.TestCheckOutput("cluster_name_del", "test-cluster-delete-no-backups"),
				},
			},
		})
	})
}
