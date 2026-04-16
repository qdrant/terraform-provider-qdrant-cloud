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

# List all members in the account
data "qdrant-cloud_accounts_members" "all" {}

# Output member emails
output "member_emails" {
  value = [for member in data.qdrant-cloud_accounts_members.all.members : member.email]
}

# Find the account owner
output "account_owner" {
  value = [for member in data.qdrant-cloud_accounts_members.all.members : member.email if member.is_owner]
}
