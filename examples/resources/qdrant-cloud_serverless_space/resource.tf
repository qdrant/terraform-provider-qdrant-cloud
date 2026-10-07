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

  labels {
    key   = "team"
    value = "search"
  }

  configuration {
    allowed_ip_source_ranges = ["10.0.0.0/8"]
    searcher_settings {
      idle_timeout = "15m"
    }
  }

  delete_backups = true // Delete the backups of the space when it is destroyed
}

// Output the endpoint of the space
output "space_url" {
  value = qdrant-cloud_serverless_space.example.url
}
