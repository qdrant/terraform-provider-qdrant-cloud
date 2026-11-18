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
    collection_settings {
      max_size                       = 4294967296 // 4 GiB per collection
      max_indexed_size               = 2147483648 // 2 GiB indexed per collection
      size_warning_threshold_percent = 90         // Warn at 90% of the maximum size
    }
    searcher_settings {
      max_workers = 2
    }
    // Cap the total storage of the space (leave out for unlimited)
    storage_limit {
      max_size_mib            = 1024              // 1 GiB (converted to bytes for the API)
      type                    = "LIMIT_TYPE_HARD" // Reject upserts at the limit (default)
      alert_threshold_percent = 80                // Alert at 80% of the limit (default)
    }
    // Cap the monthly spend on searches and writes (leave out for unlimited)
    monthly_spend_limit {
      max_amount_dollar       = 50                // 50 USD (converted to millicents for the API)
      type                    = "LIMIT_TYPE_SOFT" // Only alert at the limit
      alert_threshold_percent = 90
    }
  }

  delete_backups = true // Delete the backups of the space when it is destroyed
}

// Output the endpoint of the space
output "space_url" {
  value = qdrant-cloud_serverless_space.example.url
}
