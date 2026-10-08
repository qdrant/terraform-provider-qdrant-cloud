package qdrant

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceServerless(t *testing.T) {
	regionID := testAccServerlessCloudRegionID()
	name := "tf-acc-test-space-" + acctest.RandString(6)

	config := testAccServerlessSpaceConfig(regionID, name, "") + `
resource "qdrant-cloud_serverless_space_api_key" "test" {
  space_id = qdrant-cloud_serverless_space.test.id
  name     = "tf-acc-test-key"
  global_access_rule {
    access_type = "GLOBAL_ACCESS_RULE_ACCESS_TYPE_READ_ONLY"
  }
}

resource "qdrant-cloud_serverless_backup_schedule" "test" {
  space_id       = qdrant-cloud_serverless_space.test.id
  name           = "tf-acc-test-schedule"
  schedule       = "0 1 * * *"
  delete_backups = true
}

resource "qdrant-cloud_serverless_backup" "test" {
  space_id = qdrant-cloud_serverless_space.test.id
}

data "qdrant-cloud_serverless_cloud_regions" "test" {}

data "qdrant-cloud_serverless_space" "test" {
  id = qdrant-cloud_serverless_space.test.id
}

data "qdrant-cloud_serverless_spaces" "test" {
  cloud_region_id = qdrant-cloud_serverless_space.test.cloud_region_id
  depends_on      = [qdrant-cloud_serverless_space.test]
}

data "qdrant-cloud_serverless_space_api_keys" "test" {
  space_id   = qdrant-cloud_serverless_space.test.id
  depends_on = [qdrant-cloud_serverless_space_api_key.test]
}

data "qdrant-cloud_serverless_backup_schedule" "test" {
  space_id = qdrant-cloud_serverless_space.test.id
  id       = qdrant-cloud_serverless_backup_schedule.test.id
}

data "qdrant-cloud_serverless_backup_schedules" "test" {
  space_id   = qdrant-cloud_serverless_space.test.id
  depends_on = [qdrant-cloud_serverless_backup_schedule.test]
}

data "qdrant-cloud_serverless_backups" "test" {
  space_id   = qdrant-cloud_serverless_space.test.id
  depends_on = [qdrant-cloud_serverless_backup.test]
}
`

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		ErrorCheck:        testAccServerlessErrorCheck(t),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckTypeSetElemNestedAttrs("data.qdrant-cloud_serverless_cloud_regions.test", "regions.*", map[string]string{
						"id": regionID,
					}),
					resource.TestCheckResourceAttrPair("data.qdrant-cloud_serverless_space.test", "name", "qdrant-cloud_serverless_space.test", "name"),
					resource.TestCheckResourceAttrPair("data.qdrant-cloud_serverless_space.test", "url", "qdrant-cloud_serverless_space.test", "url"),
					resource.TestCheckTypeSetElemNestedAttrs("data.qdrant-cloud_serverless_spaces.test", "spaces.*", map[string]string{
						"name": name,
					}),
					resource.TestCheckResourceAttr("data.qdrant-cloud_serverless_space_api_keys.test", "keys.#", "1"),
					resource.TestCheckNoResourceAttr("data.qdrant-cloud_serverless_space_api_keys.test", "keys.0.key"),
					resource.TestCheckResourceAttr("data.qdrant-cloud_serverless_backup_schedule.test", "schedule", "0 1 * * *"),
					resource.TestCheckResourceAttr("data.qdrant-cloud_serverless_backup_schedules.test", "schedules.#", "1"),
					resource.TestCheckResourceAttr("data.qdrant-cloud_serverless_backups.test", "backups.#", "1"),
				),
			},
		},
	})
}
