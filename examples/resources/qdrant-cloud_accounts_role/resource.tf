terraform {
  required_version = ">= 1.7.0"
  required_providers {
    qdrant-cloud = {
      source  = "qdrant/qdrant-cloud"
      version = ">=1.13.0"
    }
  }
}

// Provider configuration
provider "qdrant-cloud" {
  api_key    = "" // API Key generated in Qdrant Cloud (required)
  account_id = "" // Default account ID (can be overridden per resource)
}

// Create a custom Role
resource "qdrant-cloud_accounts_role" "backup_operator" {
  name        = "backup-operator"
  description = "Can create and restore cluster backups."

  // Assign permissions for this role
  permissions {
    value = "read:backups"
  }

  permissions {
    value = "write:backups"
  }
}

// Outputs
output "role_id" {
  value = qdrant-cloud_accounts_role.backup_operator.id
}

output "role_name" {
  value = qdrant-cloud_accounts_role.backup_operator.name
}

output "role_created_at" {
  value = qdrant-cloud_accounts_role.backup_operator.created_at
}

output "role_last_modified_at" {
  value = qdrant-cloud_accounts_role.backup_operator.last_modified_at
}
