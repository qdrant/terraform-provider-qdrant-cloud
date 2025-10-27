package qdrant

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataAccountsBackupSchedule(t *testing.T) {
	apiKey := os.Getenv("QDRANT_CLOUD_API_KEY")
	accountID := os.Getenv("QDRANT_CLOUD_ACCOUNT_ID")
	cloudProvider := getEnvDefault("QDRANT_CLOUD_CLOUD_PROVIDER", "aws")
	cloudRegion := getEnvDefault("QDRANT_CLOUD_REGION", "eu-central-1")

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
  account_id    = "%s"
  resource_data = data.qdrant-cloud_booking_packages.test.packages

  # Filter only paid tariffs
  paid_tariffs = [
    for p in local.resource_data : p if try(p.type, "") == "paid"
  ]

  # Assign very high sentinel price if missing, so it won't be picked
  prices = [for p in local.paid_tariffs : try(p.unit_int_price_per_hour, 999999999)]

  # Find index of cheapest paid tariff
  min_price = try(min(local.prices...), null)
  min_idx   = can(local.min_price) ? index(local.prices, local.min_price) : 0

  # Final selection: cheapest paid if any; otherwise fall back to first package
  cheapest_paid_tariff = length(local.paid_tariffs) > 0 ? local.paid_tariffs[local.min_idx] : local.resource_data[0]
}

resource "qdrant-cloud_accounts_cluster" "test" {
	name           = "tf-acc-test-cluster-backup-ds-single"
	account_id     = local.account_id
	cloud_region   = "%s"
	cloud_provider = "%s"

	# Ignore enum defaults the API returns but rejects on create (enum 0)
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
			package_id = local.cheapest_paid_tariff.id
		}
	}
}

resource "qdrant-cloud_accounts_backup_schedule" "test" {
	cluster_id       = qdrant-cloud_accounts_cluster.test.id
	cron_expression  = "0 6 1 * *" // Run at 6am on the 1st of every month
	retention_period = "120h" // 5 days
}

data "qdrant-cloud_accounts_backup_schedule" "test" {
	id         = qdrant-cloud_accounts_backup_schedule.test.id
	cluster_id = qdrant-cloud_accounts_cluster.test.id
	depends_on = [qdrant-cloud_accounts_backup_schedule.test]
}
	`, cloudProvider, cloudRegion, accountID, cloudRegion, cloudProvider)

	check := resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrPair(
			"data.qdrant-cloud_accounts_backup_schedule.test", "id",
			"qdrant-cloud_accounts_backup_schedule.test", "id",
		),
		resource.TestCheckResourceAttr("data.qdrant-cloud_accounts_backup_schedule.test", "cron_expression", "0 6 1 * *"),
	)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			// Step 1: create cluster + schedule (no checks; allow eventual consistency)
			{Config: config},

			// Step 2: short wait, then re-apply same config and assert via data source
			{
				PreConfig: func() { time.Sleep(60 * time.Second) },
				Config:    config,
				Destroy:   true,
				Check:     check,
			},
		},
	})
}
