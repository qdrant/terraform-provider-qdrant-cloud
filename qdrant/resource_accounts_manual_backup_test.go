package qdrant

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Env vars used by these tests:
//   TF_ACC=1
//   QDRANT_CLOUD_API_KEY          (required)
//   QDRANT_CLOUD_ACCOUNT_ID       (required)
//   QDRANT_CLOUD_CLOUD_PROVIDER   (optional; default "aws")
//   QDRANT_CLOUD_REGION           (optional; default "eu-central-1")
//   QDRANT_CLOUD_PACKAGE_ID       (optional; when set, skips the booking_packages data source and uses this ID)

func TestAccResourceAccountsManualBackup_Basic(t *testing.T) {
	precheckAccManualBackup(t)

	apiKey := os.Getenv("QDRANT_CLOUD_API_KEY")
	accountID := os.Getenv("QDRANT_CLOUD_ACCOUNT_ID")
	cloudProvider := getEnvDefault("QDRANT_CLOUD_CLOUD_PROVIDER", "aws")
	cloudRegion := getEnvDefault("QDRANT_CLOUD_REGION", "eu-central-1")
	packageID := os.Getenv("QDRANT_CLOUD_PACKAGE_ID")

	provider := fmt.Sprintf(`
provider "qdrant-cloud" {
  api_key    = "%s"
  account_id = "%s"
}
`, apiKey, accountID)

	// Step A: Create cluster only (lets the service provision the cluster)
	configClusterOnly := provider + buildAccManualBackupClusterConfig(cloudProvider, cloudRegion, packageID, "tf-acc-test-cluster")

	// Step B: Add the manual backup resource referencing the cluster
	configClusterAndBackup := provider + buildAccManualBackupClusterAndBackupConfig(cloudProvider, cloudRegion, packageID, "tf-acc-test-cluster")

	const clusterRes = "qdrant-cloud_accounts_cluster.test"
	const backupRes = "qdrant-cloud_accounts_manual_backup.test"

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			//nolint:unparam
			"qdrant-cloud": func() (*schema.Provider, error) { return Provider(), nil },
		},
		Steps: []resource.TestStep{
			// Step 1: Create cluster only
			{
				Config: configClusterOnly,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(clusterRes, "name", "tf-acc-test-cluster"),
					resource.TestCheckResourceAttrSet(clusterRes, "id"),
					resource.TestCheckResourceAttrSet(clusterRes, "account_id"),
					resource.TestCheckResourceAttrSet(clusterRes, "created_at"),
					resource.TestCheckResourceAttrSet(clusterRes, "configuration.0.number_of_nodes"),
					testCheckHasListAttr(clusterRes, "status"),
				),
			},

			// Step 2: Create manual backup (don’t expect duration yet)
			{
				Config: configClusterAndBackup,
				Check: resource.ComposeTestCheckFunc(
					// cluster still present
					resource.TestCheckResourceAttr(clusterRes, "name", "tf-acc-test-cluster"),
					resource.TestCheckResourceAttrSet(clusterRes, "id"),

					// backup basics only (status may be RUNNING/PENDING right after create)
					resource.TestCheckResourceAttrSet(backupRes, "id"),
					resource.TestCheckResourceAttrSet(backupRes, "account_id"),
					resource.TestCheckResourceAttrSet(backupRes, "cluster_id"),
					resource.TestCheckResourceAttrSet(backupRes, "created_at"),
					resource.TestCheckResourceAttrSet(backupRes, "name"),
					resource.TestCheckResourceAttrSet(backupRes, "status"),

					// sanity: backup.cluster_id equals cluster.id
					testAccCheckAttrEqual(backupRes, "cluster_id", clusterRes, "id"),
				),
			},

			// Step 3: Wait then verify terminal fields
			// Backups can only be deleted once they are ready
			{
				PreConfig: func() { time.Sleep(240 * time.Second) },
				Config:    configClusterAndBackup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(backupRes, "status", "BACKUP_STATUS_SUCCEEDED"),
					resource.TestCheckResourceAttrSet(backupRes, "backup_duration"),

					// a couple of read-only nested fields to prove we captured snapshot metadata
					testAccCheckListNonEmpty(backupRes, "cluster_info"),
					resource.TestCheckResourceAttrSet(backupRes, "cluster_info.0.name"),
					resource.TestCheckResourceAttrSet(backupRes, "cluster_info.0.configuration.0.version"),
				),
			},

			// Step 4: Destroy both (backup and cluster)
			{
				Config:  configClusterAndBackup,
				Destroy: true,
			},
		},
	})
}

func TestAccResourceAccountsManualBackup_RetentionPeriod(t *testing.T) {
	precheckAccManualBackup(t)

	apiKey := os.Getenv("QDRANT_CLOUD_API_KEY")
	accountID := os.Getenv("QDRANT_CLOUD_ACCOUNT_ID")
	cloudProvider := getEnvDefault("QDRANT_CLOUD_CLOUD_PROVIDER", "aws")
	cloudRegion := getEnvDefault("QDRANT_CLOUD_REGION", "eu-central-1")
	packageID := os.Getenv("QDRANT_CLOUD_PACKAGE_ID")

	provider := fmt.Sprintf(`
provider "qdrant-cloud" {
  api_key    = "%s"
  account_id = "%s"
}
`, apiKey, accountID)

	const clusterRes = "qdrant-cloud_accounts_cluster.test"
	const backupRes = "qdrant-cloud_accounts_manual_backup.test"

	// Step A: Create cluster only
	configClusterOnly := provider + buildAccManualBackupClusterConfig(cloudProvider, cloudRegion, packageID, "tf-acc-test-cluster-retention")

	// Step B: Create cluster + backup with retention
	configWithRetention := provider + buildAccManualBackupClusterAndBackupWithRetentionConfig(
		cloudProvider, cloudRegion, packageID, "tf-acc-test-cluster-retention", "24h",
	)

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			//nolint:unparam
			"qdrant-cloud": func() (*schema.Provider, error) { return Provider(), nil },
		},
		Steps: []resource.TestStep{
			// Step 1: Create cluster only
			{
				Config: configClusterOnly,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(clusterRes, "name", "tf-acc-test-cluster-retention"),
					resource.TestCheckResourceAttrSet(clusterRes, "id"),
				),
			},

			// Step 2: Create manual backup with retention period
			{
				Config: configWithRetention,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(backupRes, "id"),
					resource.TestCheckResourceAttrSet(backupRes, "account_id"),
					resource.TestCheckResourceAttrSet(backupRes, "cluster_id"),
					resource.TestCheckResourceAttrSet(backupRes, "created_at"),
					resource.TestCheckResourceAttrSet(backupRes, "name"),
					resource.TestCheckResourceAttrSet(backupRes, "status"),
					resource.TestCheckResourceAttr(backupRes, "retention_period", "24h0m0s"),
					testAccCheckAttrEqual(backupRes, "cluster_id", clusterRes, "id"),
				),
			},

			// Step 3: Wait then verify terminal fields
			// Backups can only be deleted once they are ready
			{
				PreConfig: func() { time.Sleep(240 * time.Second) },
				Config:    configWithRetention,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(backupRes, "status", "BACKUP_STATUS_SUCCEEDED"),
					resource.TestCheckResourceAttrSet(backupRes, "backup_duration"),
					resource.TestCheckResourceAttr(backupRes, "retention_period", "24h0m0s"),
				),
			},

			// Step 4: Destroy both (backup and cluster)
			{
				Config:  configWithRetention,
				Destroy: true,
			},
		},
	})
}

func TestAccResourceAccountsManualBackup_Import(t *testing.T) {
	precheckAccManualBackup(t)

	apiKey := os.Getenv("QDRANT_CLOUD_API_KEY")
	accountID := os.Getenv("QDRANT_CLOUD_ACCOUNT_ID")
	cloudProvider := getEnvDefault("QDRANT_CLOUD_CLOUD_PROVIDER", "aws")
	cloudRegion := getEnvDefault("QDRANT_CLOUD_REGION", "eu-central-1")
	packageID := os.Getenv("QDRANT_CLOUD_PACKAGE_ID")

	provider := fmt.Sprintf(`
provider "qdrant-cloud" {
  api_key    = "%s"
  account_id = "%s"
}
`, apiKey, accountID)

	config := provider + buildAccManualBackupClusterAndBackupConfig(
		cloudProvider, cloudRegion, packageID, "tf-acc-test-cluster-import",
	)

	const backupRes = "qdrant-cloud_accounts_manual_backup.test"

	resource.Test(t, resource.TestCase{
		//nolint:unparam
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"qdrant-cloud": func() (*schema.Provider, error) { return Provider(), nil },
		},
		Steps: []resource.TestStep{
			// 1) Create cluster + backup
			{Config: config},

			// 2) Import the backup
			{
				ResourceName:      backupRes,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[backupRes]
					if !ok {
						return "", fmt.Errorf("resource not found in state: %s", backupRes)
					}
					return rs.Primary.ID, nil
				},
			},

			// 3) Wait for backup to finish, then assert terminal fields
			{
				PreConfig: func() { time.Sleep(240 * time.Second) },
				Config:    config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(backupRes, "status", "BACKUP_STATUS_SUCCEEDED"),
					resource.TestCheckResourceAttrSet(backupRes, "backup_duration"),
				),
			},

			// 4) Destroy (backup then cluster)
			{Config: config, Destroy: true},
		},
	})
}

// ------------------ helpers ------------------

func precheckAccManualBackup(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set, skipping acceptance tests.")
	}
	if os.Getenv("QDRANT_CLOUD_API_KEY") == "" {
		t.Skip("QDRANT_CLOUD_API_KEY not set, skipping acceptance tests.")
	}
	if os.Getenv("QDRANT_CLOUD_ACCOUNT_ID") == "" {
		t.Skip("QDRANT_CLOUD_ACCOUNT_ID not set, skipping acceptance tests.")
	}
}

func getEnvDefault(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func buildAccManualBackupClusterConfig(cloudProvider, cloudRegion, packageID, name string) string {
	if packageID != "" {
		// Use explicit package id (skip data source)
		return fmt.Sprintf(`
resource "qdrant-cloud_accounts_cluster" "test" {
  name           = "%s"
  cloud_provider = "%s"
  cloud_region   = "%s"

  # Ignore enum defaults the API returns but rejects on create (enum 0)
  lifecycle {
    ignore_changes = [
      configuration[0].gpu_type,
      configuration[0].restart_policy,
    ]
  }

  configuration {
    number_of_nodes = 1

    database_configuration {
      service { jwt_rbac = true }
    }

    node_configuration {
      package_id = "%s"
    }

    # keep plans stable (avoid drift)
    allowed_ip_source_ranges = []
    service_type             = "CLUSTER_SERVICE_TYPE_CLUSTER_IP"
    rebalance_strategy       = "CLUSTER_CONFIGURATION_REBALANCE_STRATEGY_BY_COUNT_AND_SIZE"
  }
}
`, name, cloudProvider, cloudRegion, packageID)
	}

	// Discover a package via data source (pick the first)
	return fmt.Sprintf(`
data "qdrant-cloud_booking_packages" "all" {
  cloud_provider = "%s"
  cloud_region   = "%s"
}

locals {
  selected = data.qdrant-cloud_booking_packages.all.packages[0]
}

resource "qdrant-cloud_accounts_cluster" "test" {
  name           = "%s"
  cloud_provider = data.qdrant-cloud_booking_packages.all.cloud_provider
  cloud_region   = data.qdrant-cloud_booking_packages.all.cloud_region

  lifecycle {
    ignore_changes = [
      configuration[0].gpu_type,
      configuration[0].restart_policy,
    ]
  }

  configuration {
    number_of_nodes = 1

    database_configuration {
      service { jwt_rbac = true }
    }

    node_configuration {
      package_id = local.selected.id
    }

    allowed_ip_source_ranges = []
    service_type             = "CLUSTER_SERVICE_TYPE_CLUSTER_IP"
    rebalance_strategy       = "CLUSTER_CONFIGURATION_REBALANCE_STRATEGY_BY_COUNT_AND_SIZE"
  }
}
`, cloudProvider, cloudRegion, name)
}

// Build config that creates a cluster and a manual backup referencing it.
func buildAccManualBackupClusterAndBackupConfig(cloudProvider, cloudRegion, packageID, name string) string {
	return buildAccManualBackupClusterConfig(cloudProvider, cloudRegion, packageID, name) + `
resource "qdrant-cloud_accounts_manual_backup" "test" {
  cluster_id = qdrant-cloud_accounts_cluster.test.id
}
`
}

// Build config that creates a cluster and a manual backup with a retention period.
func buildAccManualBackupClusterAndBackupWithRetentionConfig(cloudProvider, cloudRegion, packageID, name, retention string) string {
	return buildAccManualBackupClusterConfig(cloudProvider, cloudRegion, packageID, name) + fmt.Sprintf(`
resource "qdrant-cloud_accounts_manual_backup" "test" {
  cluster_id       = qdrant-cloud_accounts_cluster.test.id
  retention_period = "%s"
}
`, retention)
}
