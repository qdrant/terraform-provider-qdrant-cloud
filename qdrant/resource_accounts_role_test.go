package qdrant

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceAccountsRole_Basic(t *testing.T) {
	precheckAccRole(t)

	apiKey := os.Getenv("QDRANT_CLOUD_API_KEY")
	accountID := os.Getenv("QDRANT_CLOUD_ACCOUNT_ID")

	provider := fmt.Sprintf(`
provider "qdrant-cloud" {
  api_key    = "%s"
  account_id = "%s"
}
`, apiKey, accountID)

	const roleRes = "qdrant-cloud_accounts_role.test"

	config := provider + `
resource "qdrant-cloud_accounts_role" "test" {
  name        = "tf-acc-role-basic"
  description = "basic acceptance role"

  permissions {
    value = "read:clusters"
  }

  permissions {
    value = "read:backups"
  }
}
`

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			//nolint:unparam
			"qdrant-cloud": func() (*schema.Provider, error) { return Provider(), nil },
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(roleRes, "name", "tf-acc-role-basic"),
					resource.TestCheckResourceAttr(roleRes, "role_type", "ROLE_TYPE_CUSTOM"),

					resource.TestCheckResourceAttrSet(roleRes, "id"),
					resource.TestCheckResourceAttrSet(roleRes, "account_id"),
					resource.TestCheckResourceAttrSet(roleRes, "created_at"),

					// We only assert the count because permissions is a set (unordered).
					resource.TestCheckResourceAttr(roleRes, "permissions.#", "2"),
				),
			},
			{
				ResourceName:      roleRes,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[roleRes]
					if !ok {
						return "", fmt.Errorf("resource not found in state: %s", roleRes)
					}
					return rs.Primary.ID, nil
				},
			},
		},
	})
}

func TestAccResourceAccountsRole_Update(t *testing.T) {
	precheckAccRole(t)

	apiKey := os.Getenv("QDRANT_CLOUD_API_KEY")
	accountID := os.Getenv("QDRANT_CLOUD_ACCOUNT_ID")

	provider := fmt.Sprintf(`
provider "qdrant-cloud" {
  api_key    = "%s"
  account_id = "%s"
}
`, apiKey, accountID)

	const roleRes = "qdrant-cloud_accounts_role.test"

	configCreate := provider + `
resource "qdrant-cloud_accounts_role" "test" {
  name        = "tf-acc-role-update"
  description = "before update"

  permissions {
    value = "read:clusters"
  }
}
`

	configUpdate := provider + `
resource "qdrant-cloud_accounts_role" "test" {
  name        = "tf-acc-role-update"
  description = "after update"

  permissions {
    value = "read:clusters"
  }
  permissions {
    value = "write:backups"
  }
  permissions {
    value = "restore:backups"
  }
}
`

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			//nolint:unparam
			"qdrant-cloud": func() (*schema.Provider, error) { return Provider(), nil },
		},
		Steps: []resource.TestStep{
			{
				Config: configCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(roleRes, "name", "tf-acc-role-update"),
					resource.TestCheckResourceAttr(roleRes, "description", "before update"),
					resource.TestCheckResourceAttr(roleRes, "role_type", "ROLE_TYPE_CUSTOM"),
					resource.TestCheckResourceAttr(roleRes, "permissions.#", "1"),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(roleRes, "name", "tf-acc-role-update"),
					resource.TestCheckResourceAttr(roleRes, "description", "after update"),
					resource.TestCheckResourceAttr(roleRes, "role_type", "ROLE_TYPE_CUSTOM"),
					resource.TestCheckResourceAttr(roleRes, "permissions.#", "3"),
				),
			},
		},
	})
}

func TestAccResourceAccountsRole_Import(t *testing.T) {
	precheckAccRole(t)

	apiKey := os.Getenv("QDRANT_CLOUD_API_KEY")
	accountID := os.Getenv("QDRANT_CLOUD_ACCOUNT_ID")

	provider := fmt.Sprintf(`
provider "qdrant-cloud" {
  api_key    = "%s"
  account_id = "%s"
}
`, apiKey, accountID)

	const roleRes = "qdrant-cloud_accounts_role.test"

	config := provider + `
resource "qdrant-cloud_accounts_role" "test" {
  name        = "tf-acc-role-import"
  description = "import me"

  permissions {
    value = "read:clusters"
  }
  permissions {
    value = "read:backups"
  }
}
`

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			//nolint:unparam
			"qdrant-cloud": func() (*schema.Provider, error) { return Provider(), nil },
		},
		Steps: []resource.TestStep{
			{Config: config},
			{
				ResourceName:      roleRes,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: func(s *terraform.State) (string, error) {
					rs, ok := s.RootModule().Resources[roleRes]
					if !ok {
						return "", fmt.Errorf("resource not found in state: %s", roleRes)
					}
					return rs.Primary.ID, nil
				},
			},
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(roleRes, "name", "tf-acc-role-import"),
					resource.TestCheckResourceAttr(roleRes, "role_type", "ROLE_TYPE_CUSTOM"),
					resource.TestCheckResourceAttr(roleRes, "permissions.#", "2"),
				),
			},
		},
	})
}

// ------------------ helpers ------------------

func precheckAccRole(t *testing.T) {
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
