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

// Create an API key with read-only access to the entire space
resource "qdrant-cloud_serverless_space_api_key" "read_only" {
  space_id = qdrant-cloud_serverless_space.example.id
  name     = "read-only-key"
  global_access_rule {
    access_type = "GLOBAL_ACCESS_RULE_ACCESS_TYPE_READ_ONLY"
  }
}

// Create an API key with write access to a single collection, expiring at the end of the year
resource "qdrant-cloud_serverless_space_api_key" "collection_writer" {
  space_id   = qdrant-cloud_serverless_space.example.id
  name       = "collection-writer-key"
  expires_at = "2026-12-31T23:59:59Z"
  collection_access_rules {
    collection_name = "my-collection"
    access_type     = "COLLECTION_ACCESS_RULE_ACCESS_TYPE_READ_WRITE"
  }
}

// Output the secret key (only available after creation)
output "read_only_key" {
  value     = qdrant-cloud_serverless_space_api_key.read_only.key
  sensitive = true
}
