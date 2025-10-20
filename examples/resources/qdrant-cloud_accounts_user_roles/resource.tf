terraform {
  required_version = ">= 1.7.0"
  required_providers {
    qdrant-cloud = {
      source  = "qdrant/qdrant-cloud"
      version = ">=1.13.0"
    }
  }
}

# Provider configuration
provider "qdrant-cloud" {
  api_key    = "" # API Key generated in Qdrant Cloud (required)
  account_id = "" # Default account ID (can be overridden per resource)
}

# Manage role assignments for a single user
resource "qdrant-cloud_accounts_user_roles" "example" {
  # Optional; defaults to provider.account_id
  # account_id = ""

  user_email = "alice@example.com"

  # Role IDs to assign to this user (unordered)
  role_ids = [
    "11111111-1111-1111-1111-111111111111",
    "22222222-2222-2222-2222-222222222222",
  ]

  # Optional: retain roles after destroy (skip revocation)
  # keep_on_destroy = true
}

# Outputs
output "user_id" {
  value = qdrant-cloud_accounts_user_roles.example.user_id
}

output "ensured_role_ids" {
  value = qdrant-cloud_accounts_user_roles.example.role_ids
}
