// Setup Terraform, including the qdrant-cloud providers
terraform {
  required_version = ">= 1.7.0"
  required_providers {
    qdrant-cloud = {
      source  = "local/qdrant-cloud/qdrant-cloud"
      version = ">=1.0"
    }
  }
}

// Add the provider to specify some provider wide settings
provider "qdrant-cloud" {
  api_key    = "" // API Key generated in Qdrant Cloud (required)
  api_url    = "" // URL where the public API of Qdrant cloud can be found (can be left empty if the production URL need to be used)
  account_id = "" // The default account ID you want to use in Qdrant Cloud (can be overriden on resource level)
}

// Reference the booking packages
data "qdrant-cloud_booking_packages" "test" {
  // No filter here
}

// Output the packages (Note that the package id is relevent for creating a cluster)
output "packages" {
  value = data.qdrant-cloud_booking_packages.test.packages
}

