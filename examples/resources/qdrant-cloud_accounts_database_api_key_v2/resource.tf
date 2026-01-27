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

// Get the cluster package
data "qdrant-cloud_booking_packages" "all_packages" {
  cloud_provider = "aws"       // Required. Please refer to the documentation (https://registry.terraform.io/providers/qdrant/qdrant-cloud/latest/docs/guides/getting-started) for the available options.
  cloud_region   = "us-west-2" // Required. Please refer to the documentation (https://registry.terraform.io/providers/qdrant/qdrant-cloud/latest/docs/guides/getting-started) for the available options.
}
locals {
  desired_package = [
    for pkg in data.qdrant-cloud_booking_packages.all_packages.packages : pkg
    if pkg.resource_configuration[0].cpu == "16000m" && pkg.resource_configuration[0].ram == "64Gi"
  ]
}

// Create a cluster (for the sake of having an ID, see below)
resource "qdrant-cloud_accounts_cluster" "example" {
  name           = "example-cluster"
  cloud_provider = data.qdrant-cloud_booking_packages.all_packages.cloud_provider
  cloud_region   = data.qdrant-cloud_booking_packages.all_packages.cloud_region
  configuration {
    number_of_nodes = 1
    node_configuration {
      package_id = local.desired_package[0].id
    }
  }
}

// Create an V2 Auth Key, which refers to the cluster provided above
resource "qdrant-cloud_accounts_database_api_key_v2" "example-key" {
  cluster_id = qdrant-cloud_accounts_cluster.example.id
  name       = "example-key"
}

// Output the Database API Key (which can be used to access the database cluster)
output "key" {
  value = qdrant-cloud_accounts_database_api_key_v2.example-key.key
}