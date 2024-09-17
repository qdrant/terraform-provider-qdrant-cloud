terraform {
  required_version = ">= 1.7.0"
  required_providers {
    qdrant-cloud = {
      source  = "qdrant/qdrant-cloud"
      version = ">=1.1.0"
    }
  }
}

provider "qdrant-cloud" {
  api_key    = "" // API Key generated in Qdrant Cloud (required)
  api_url    = "" // URL where the public API of Qdrant cloud can be found (optional: defaults to production URL).
  account_id = "" // The default account ID you want to use in Qdrant Cloud (can be overriden on resource level)
}

resource "qdrant-cloud_accounts_cluster" "example" {
  name           = "tf-example-cluster"
  cloud_provider = "gcp"
  cloud_region   = "us-east4"
  configuration {
    number_of_nodes = 1
    node_configuration {
      package_id = "7c939d96-d671-4051-aa16-3b8b7130fa42" # gpx1
    }
  }
}

resource "qdrant-cloud_accounts_auth_key" "example-key" {
  cluster_ids = [qdrant-cloud_accounts_cluster.example.id]
}

output "cluster_id" {
  value = qdrant-cloud_accounts_cluster.example.id
}

output "url" {
  value = qdrant-cloud_accounts_cluster.example.url
}

output "token" {
  value       = qdrant-cloud_accounts_auth_key.example-key.token
  description = "Token is available only once, after creation."
}

output "curl_command" {
  value       = "curl \\\n    -X GET '${qdrant-cloud_accounts_cluster.example.url}' \\\n    --header 'api-key: ${qdrant-cloud_accounts_auth_key.example-key.token}'"
  description = "Generating a curl command test cluster access using the API key."
}