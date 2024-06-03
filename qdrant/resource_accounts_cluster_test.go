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
resource "qdrant-cloud_accounts_cluster" "test" {
	name = "test-cluster"
	account_id = "%s"
	cloud_region = "us-east4"
	cloud_provider = "gcp"

	configuration {
		num_nodes_max = 1
		num_nodes = 1

		node_configuration {
			package_id = "39b48a76-2a60-4ee0-9266-6d1e0f91ea14"
		}
	}
}

output "cluster_id" {
	value = qdrant-cloud_accounts_cluster.test.id
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
						Config:             config,
						PlanOnly:           true,
						ExpectNonEmptyPlan: false,
						Destroy:            true,
						Check: resource.ComposeTestCheckFunc(
							resource.TestCheckOutput("cluster_name", "test-cluster"),
						),
					},
				},
			})
		}
		testCase(t, "apply")
	})
}
