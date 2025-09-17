terraform {
  required_version = ">= 1.7.0"
  required_providers {
    qdrant-cloud = {
      source  = "qdrant/qdrant-cloud"
      version = ">=1.1.0"
    }
  }
}

// Provider configuration
// Tip: you can also supply these via environment variables or terraform.tfvars.
provider "qdrant-cloud" {
  api_key    = "" // API Key generated in Qdrant Cloud (required)
  account_id = "" // Default account ID (can be overridden per resource)
}

// Create a Hybrid Cloud Environment
resource "qdrant-cloud_accounts_hybrid_cloud_environment" "example" {
  name = "example-hc-env"

  configuration {
    // Namespace for Qdrant components in your Kubernetes cluster
    namespace = "qdrant-hc"
  }
}

// Outputs
output "hc_environment_id" {
  value = qdrant-cloud_accounts_hybrid_cloud_environment.example.id
}

output "hc_environment_name" {
  value = qdrant-cloud_accounts_hybrid_cloud_environment.example.name
}

output "hc_environment_namespace" {
  value = qdrant-cloud_accounts_hybrid_cloud_environment.example.configuration[0].namespace
}

output "hc_environment_created_at" {
  value = qdrant-cloud_accounts_hybrid_cloud_environment.example.created_at
}

output "hc_environment_last_modified_at" {
  value = qdrant-cloud_accounts_hybrid_cloud_environment.example.last_modified_at
}
