// Setup Terraform, including the qdrant-cloud providers
terraform {
  required_version = ">= 1.7.0"
  required_providers {
    qdrant-cloud = {
      source  = "qdrant/qdrant-cloud"
      version = ">=1.1.0"
    }
  }
}

// Add the provider to specify some provider wide settings
provider "qdrant-cloud" {
  api_key    = "" // API Key generated in Qdrant Cloud (required)
  account_id = "" // The default account ID you want to use in Qdrant Cloud (can be overriden on resource level)
}


data "qdrant-cloud_accounts_auth_keys" "test" {
  // No keys needed here, the account ID is specified on provider level
}

output "keys" {
  value = data.qdrant-cloud_accounts_auth_keys.test.keys
}