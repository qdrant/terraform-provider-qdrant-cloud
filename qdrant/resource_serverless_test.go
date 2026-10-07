package qdrant

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// testAccServerlessCloudRegionID returns the cloud region to create serverless spaces in.
func testAccServerlessCloudRegionID() string {
	return getEnvDefault("QDRANT_CLOUD_SERVERLESS_REGION_ID", "eu-north-1")
}

func testAccServerlessSpaceConfig(regionID, name, extra string) string {
	return fmt.Sprintf(`
resource "qdrant-cloud_serverless_space" "test" {
  name            = "%s"
  cloud_region_id = "%s"
  delete_backups  = true
%s
}
`, name, regionID, extra)
}

// testAccImportSpaceChildID returns the `<space_id>/<id>` import ID of a resource that belongs to a space.
func testAccImportSpaceChildID(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource %s not found in state", resourceName)
		}
		return fmt.Sprintf("%s/%s", rs.Primary.Attributes["space_id"], rs.Primary.ID), nil
	}
}

func TestAccResourceServerlessSpace(t *testing.T) {
	regionID := testAccServerlessCloudRegionID()
	name := "tf-acc-test-space-" + acctest.RandString(6)
	resourceName := "qdrant-cloud_serverless_space.test"

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServerlessSpaceConfig(regionID, name, ""),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "cloud_region_id", regionID),
					resource.TestCheckResourceAttr(resourceName, "state.0.phase", "SPACE_STATE_PHASE_READY"),
					resource.TestCheckResourceAttrSet(resourceName, "url"),
				),
			},
			{
				Config: testAccServerlessSpaceConfig(regionID, name+"-upd", `
  labels {
    key   = "env"
    value = "tf-acc-test"
  }
  configuration {
    allowed_ip_source_ranges = ["10.0.0.0/8"]
    allowed_origins          = ["https://example.com"]
    searcher_settings {
      idle_timeout = "10m"
    }
  }
`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name+"-upd"),
					resource.TestCheckResourceAttr(resourceName, "labels.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.allowed_ip_source_ranges.0", "10.0.0.0/8"),
					resource.TestCheckResourceAttr(resourceName, "configuration.0.allowed_origins.0", "https://example.com"),
					resource.TestCheckResourceAttr(resourceName, "state.0.phase", "SPACE_STATE_PHASE_READY"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_backups"},
			},
		},
	})
}

func TestAccResourceServerlessSpaceChildren(t *testing.T) {
	regionID := testAccServerlessCloudRegionID()
	name := "tf-acc-test-space-" + acctest.RandString(6)
	apiKeyName := "qdrant-cloud_serverless_space_api_key.test"
	scheduleName := "qdrant-cloud_serverless_backup_schedule.test"
	backupName := "qdrant-cloud_serverless_backup.test"
	restoredName := "qdrant-cloud_serverless_space.restored"

	children := func(schedule string, paused bool) string {
		return testAccServerlessSpaceConfig(regionID, name, "") + fmt.Sprintf(`
resource "qdrant-cloud_serverless_space_api_key" "test" {
  space_id = qdrant-cloud_serverless_space.test.id
  name     = "tf-acc-test-key"
  global_access_rule {
    access_type = "GLOBAL_ACCESS_RULE_ACCESS_TYPE_READ_ONLY"
  }
}

resource "qdrant-cloud_serverless_backup_schedule" "test" {
  space_id         = qdrant-cloud_serverless_space.test.id
  name             = "tf-acc-test-schedule"
  schedule         = "%s"
  retention_period = "168h"
  paused           = %t
  delete_backups   = true
}

resource "qdrant-cloud_serverless_backup" "test" {
  space_id         = qdrant-cloud_serverless_space.test.id
  retention_period = "72h"
}
`, schedule, paused)
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: children("0 1 * * *", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(apiKeyName, "key"),
					resource.TestCheckResourceAttrSet(apiKeyName, "postfix"),
					resource.TestCheckResourceAttr(apiKeyName, "state.0.phase", "SPACE_API_KEY_STATE_PHASE_READY"),
					resource.TestCheckResourceAttr(scheduleName, "schedule", "0 1 * * *"),
					resource.TestCheckResourceAttr(scheduleName, "retention_period", "168h0m0s"),
					resource.TestCheckResourceAttr(scheduleName, "paused", "false"),
					resource.TestCheckResourceAttr(backupName, "status", "BACKUP_STATUS_SUCCEEDED"),
					resource.TestCheckResourceAttrSet(backupName, "name"),
				),
			},
			{
				Config: children("0 3 * * *", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(scheduleName, "schedule", "0 3 * * *"),
					resource.TestCheckResourceAttr(scheduleName, "paused", "true"),
					resource.TestCheckResourceAttrSet(scheduleName, "paused_at"),
				),
			},
			{
				ResourceName:            apiKeyName,
				ImportState:             true,
				ImportStateIdFunc:       testAccImportSpaceChildID(apiKeyName),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"key"},
			},
			{
				ResourceName:            scheduleName,
				ImportState:             true,
				ImportStateIdFunc:       testAccImportSpaceChildID(scheduleName),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_backups"},
			},
			{
				ResourceName:      backupName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Restore the manual backup into a new space.
				Config: children("0 3 * * *", true) + fmt.Sprintf(`
resource "qdrant-cloud_serverless_space" "restored" {
  name           = "%s-restored"
  from_backup_id = qdrant-cloud_serverless_backup.test.id
}
`, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(restoredName, "state.0.phase", "SPACE_STATE_PHASE_READY"),
					resource.TestCheckResourceAttrPair(restoredName, "cloud_region_id", "qdrant-cloud_serverless_space.test", "cloud_region_id"),
				),
			},
		},
	})
}
