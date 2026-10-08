// Setup Terraform, including the qdrant-cloud providers
terraform {
  required_version = ">= 1.7.0"
  required_providers {
    qdrant-cloud = {
      source  = "qdrant/qdrant-cloud"
      version = ">=1.30.0"
    }
  }
}

// Add the provider to specify some provider wide settings
provider "qdrant-cloud" {
  api_key    = "" // API Key generated in Qdrant Cloud (required)
  account_id = "" // The default account ID you want to use in Qdrant Cloud (can be overriden on resource level)
}

// Pick the first available serverless cloud region
data "qdrant-cloud_serverless_cloud_regions" "all" {}
locals {
  cloud_region_id = [for r in data.qdrant-cloud_serverless_cloud_regions.all.regions : r.id if r.available][0]
}

// Create a serverless space
resource "qdrant-cloud_serverless_space" "example" {
  name            = "example-space"
  cloud_region_id = local.cloud_region_id
}

// Create a manual backup of the entire space
resource "qdrant-cloud_serverless_backup" "example" {
  space_id         = qdrant-cloud_serverless_space.example.id
  retention_period = "72h" // Retain for 3 days
}

// Restore the backup into a new space
resource "qdrant-cloud_serverless_space" "restored" {
  name           = "example-space-restored"
  from_backup_id = qdrant-cloud_serverless_backup.example.id
}

// Output the Backup ID
output "backup_id" {
  value = qdrant-cloud_serverless_backup.example.id
}
