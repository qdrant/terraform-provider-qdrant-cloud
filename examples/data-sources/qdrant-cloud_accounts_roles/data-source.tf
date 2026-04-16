terraform {
  required_version = ">= 1.7.0"
  required_providers {
    qdrant-cloud = {
      source  = "qdrant/qdrant-cloud"
      version = ">=1.13.0"
    }
  }
}

provider "qdrant-cloud" {
  api_key    = "" # API Key generated in Qdrant Cloud (required)
  account_id = "" # Default account ID (can be overridden per data source)
}

# List all roles (system + custom) in the account
data "qdrant-cloud_accounts_roles" "all" {}

# Output all role names and their types
output "roles" {
  value = {
    for role in data.qdrant-cloud_accounts_roles.all.roles :
    role.name => role.role_type
  }
}

# Look up the admin system role ID
output "admin_role_id" {
  value = [
    for role in data.qdrant-cloud_accounts_roles.all.roles :
    role.id if role.name == "Admin"
  ]
}
