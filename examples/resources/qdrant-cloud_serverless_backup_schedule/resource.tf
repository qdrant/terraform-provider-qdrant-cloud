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

// Create a backup schedule for the entire space
resource "qdrant-cloud_serverless_backup_schedule" "example" {
  space_id         = qdrant-cloud_serverless_space.example.id
  name             = "nightly"
  schedule         = "0 2 * * *" // Daily at 02:00
  retention_period = "168h"      // Retain for 7 days
}

// Output the Backup Schedule ID
output "backup_schedule_id" {
  value = qdrant-cloud_serverless_backup_schedule.example.id
}
