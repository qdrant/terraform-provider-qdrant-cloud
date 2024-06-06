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

// Create a cluster
resource "qdrant-cloud_accounts_cluster" "example" {
  name           = "example-cluster"
  cloud_provider = "gcp"
  cloud_region   = "us-east4"
  configuration {
    num_nodes     = 1
    num_nodes_max = 1 // TODO: Remove
    node_configuration {
      package_id = "7c939d96-d671-4051-aa16-3b8b7130fa42" # gpx1
      // package_id = "39b48a76-2a60-4ee0-9266-6d1e0f91ea14" # free
    }
  }
}

// Output some of the cluster info
output "cluster_id" {
  value = qdrant-cloud_accounts_cluster.example.id
}

output "cluster_version" {
  value = qdrant-cloud_accounts_cluster.example.version
}

output "url" {
  value = qdrant-cloud_accounts_cluster.example.url
}

