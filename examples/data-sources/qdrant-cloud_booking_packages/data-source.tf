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

// Reference the booking packages
data "qdrant-cloud_booking_packages" "test" {
  cloud_provider = "aws"       // Required. Please refer to the documentation (https://registry.terraform.io/providers/qdrant/qdrant-cloud/latest/docs/guides/getting-started) for the available options.
  cloud_region   = "us-west-2" // Required. Please refer to the documentation (https://registry.terraform.io/providers/qdrant/qdrant-cloud/latest/docs/guides/getting-started) for the available options.
}

// Output the packages (Note that the package id is relevent for creating a cluster)
output "packages" {
  value = data.qdrant-cloud_booking_packages.test.packages
}

// Get a package with a specific configuration
locals {
  desired_package = [
    for pkg in data.qdrant-cloud_booking_packages.test.packages : pkg
    if pkg.resource_configuration[0].cpu == "16000m" && pkg.resource_configuration[0].ram == "64Gi"
  ]
}
