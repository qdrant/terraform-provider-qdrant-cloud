package qdrant

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccResourceAccountsUserRoles_Basic(t *testing.T) {
	precheckAccUserRoles(t)

	apiKey := os.Getenv("QDRANT_CLOUD_API_KEY")
	accountID := os.Getenv("QDRANT_CLOUD_ACCOUNT_ID")
	userEmail := os.Getenv("QDRANT_CLOUD_TEST_USER_EMAIL")

	provider := fmt.Sprintf(`
provider "qdrant-cloud" {
  api_key    = "%s"
  account_id = "%s"
}
`, apiKey, accountID)

	const resUR = "qdrant-cloud_accounts_user_roles.test"

	config := provider + `
resource "qdrant-cloud_accounts_role" "r1" {
  name        = "tf-acc-user-roles-basic-r1"
  description = "role one"

  permissions {
    value = "read:clusters"
  }
}

resource "qdrant-cloud_accounts_role" "r2" {
  name        = "tf-acc-user-roles-basic-r2"
  description = "role two"

  permissions {
    value = "read:backups"
  }
}

resource "qdrant-cloud_accounts_user_roles" "test" {
  user_email = "` + userEmail + `"

  role_ids = [
    qdrant-cloud_accounts_role.r1.id,
    qdrant-cloud_accounts_role.r2.id
  ]
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
					resource.TestCheckResourceAttr(resUR, "user_email", userEmail),

					// Computed/id fields present
					resource.TestCheckResourceAttrSet(resUR, "id"),
					resource.TestCheckResourceAttrSet(resUR, "account_id"),
					resource.TestCheckResourceAttrSet(resUR, "user_id"),

					// We can safely assert count; role_ids is an input set (unordered).
					resource.TestCheckResourceAttr(resUR, "role_ids.#", "2"),
				),
			},
		},
	})
}

func TestAccResourceAccountsUserRoles_Update(t *testing.T) {
	precheckAccUserRoles(t)

	apiKey := os.Getenv("QDRANT_CLOUD_API_KEY")
	accountID := os.Getenv("QDRANT_CLOUD_ACCOUNT_ID")
	userEmail := os.Getenv("QDRANT_CLOUD_TEST_USER_EMAIL")

	provider := fmt.Sprintf(`
provider "qdrant-cloud" {
  api_key    = "%s"
  account_id = "%s"
}
`, apiKey, accountID)

	const resUR = "qdrant-cloud_accounts_user_roles.test"

	// Create with two roles
	configCreate := provider + `
resource "qdrant-cloud_accounts_role" "r1" {
  name        = "tf-acc-user-roles-update-r1"
  description = "role one"

  permissions {
    value = "read:clusters"
  }
}

resource "qdrant-cloud_accounts_role" "r2" {
  name        = "tf-acc-user-roles-update-r2"
  description = "role two"

  permissions {
    value = "read:backups"
  }
}

resource "qdrant-cloud_accounts_user_roles" "test" {
  user_email = "` + userEmail + `"

  role_ids = [
    qdrant-cloud_accounts_role.r1.id,
    qdrant-cloud_accounts_role.r2.id
  ]
}
`

	// Update to a single role (removes r2; provider should call AssignUserRoles with RoleIdsToDelete)
	configUpdate := provider + `
resource "qdrant-cloud_accounts_role" "r1" {
  name        = "tf-acc-user-roles-update-r1"
  description = "role one"

  permissions {
    value = "read:clusters"
  }
}

resource "qdrant-cloud_accounts_role" "r2" {
  name        = "tf-acc-user-roles-update-r2"
  description = "role two"

  permissions {
    value = "read:backups"
  }
}

resource "qdrant-cloud_accounts_user_roles" "test" {
  user_email = "` + userEmail + `"

  role_ids = [
    qdrant-cloud_accounts_role.r1.id
  ]
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
					resource.TestCheckResourceAttr(resUR, "user_email", userEmail),
					resource.TestCheckResourceAttr(resUR, "role_ids.#", "2"),
					resource.TestCheckResourceAttrSet(resUR, "user_id"),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resUR, "user_email", userEmail),
					resource.TestCheckResourceAttr(resUR, "role_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resUR, "user_id"),
				),
			},
		},
	})
}

func TestAccResourceAccountsUserRoles_KeepOnDestroy(t *testing.T) {
	precheckAccUserRoles(t)

	apiKey := os.Getenv("QDRANT_CLOUD_API_KEY")
	accountID := os.Getenv("QDRANT_CLOUD_ACCOUNT_ID")
	userEmail := os.Getenv("QDRANT_CLOUD_TEST_USER_EMAIL")

	provider := fmt.Sprintf(`
provider "qdrant-cloud" {
  api_key    = "%s"
  account_id = "%s"
}
`, apiKey, accountID)

	const resUR = "qdrant-cloud_accounts_user_roles.test"

	// Create one role and assign it with keep_on_destroy = true.
	config := provider + `
resource "qdrant-cloud_accounts_role" "r1" {
  name        = "tf-acc-user-roles-keep-r1"
  description = "role one"

  permissions {
    value = "read:clusters"
  }
}

resource "qdrant-cloud_accounts_user_roles" "test" {
  user_email      = "` + userEmail + `"
  keep_on_destroy = true

  role_ids = [
    qdrant-cloud_accounts_role.r1.id
  ]
}
`

	// Provider-only config used to destroy the resource in the next step.
	configDestroy := provider

	resource.Test(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			//nolint:unparam
			"qdrant-cloud": func() (*schema.Provider, error) { return Provider(), nil },
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resUR, "user_email", userEmail),
					resource.TestCheckResourceAttr(resUR, "keep_on_destroy", "true"),
					resource.TestCheckResourceAttr(resUR, "role_ids.#", "1"),
					resource.TestCheckResourceAttrSet(resUR, "user_id"),
				),
			},
			{
				// Destroy the resource; since keep_on_destroy=true, provider should skip revocation.
				Config:  configDestroy,
				Destroy: true,
			},
		},
	})
}

// ------------------ helpers ------------------

func precheckAccUserRoles(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set, skipping acceptance tests.")
	}
	if os.Getenv("QDRANT_CLOUD_API_KEY") == "" {
		t.Skip("QDRANT_CLOUD_API_KEY not set, skipping acceptance tests.")
	}
	if os.Getenv("QDRANT_CLOUD_ACCOUNT_ID") == "" {
		t.Skip("QDRANT_CLOUD_ACCOUNT_ID not set, skipping acceptance tests.")
	}
	if os.Getenv("QDRANT_CLOUD_TEST_USER_EMAIL") == "" {
		t.Skip("QDRANT_CLOUD_TEST_USER_EMAIL not set, skipping acceptance tests.")
	}
}
