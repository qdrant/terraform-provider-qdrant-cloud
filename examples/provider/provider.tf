terraform {
  required_version = ">= 1.0.0"
  required_providers {
    qdrant = {
      source  = "qdrant/qdrant-cloud"
      version = ">=1.0.0"
    }
  }
}

provider "qdrant-cloud" {
  api_key    = "" // API Key generated in Qdrant Cloud (required)
  api_url    = "" // URL where the public API of Qdrant cloud can be found (optional: defaults to production URL).
  account_id = "" // The default account ID you want to use in Qdrant Cloud (can be overriden on resource level)
}

resource "qdrant-cloud_accounts_cluster" "example" {
  name           = "example-cluster"
  cloud_provider = "gcp"
  cloud_region   = "us-east4"
  configuration {
    num_nodes     = 1
    num_nodes_max = 1
    node_configuration {
      package_id = "39b48a76-2a60-4ee0-9266-6d1e0f91ea14" # free
      // package_id = "7c939d96-d671-4051-aa16-3b8b7130fa42" # gpx1
    }
  }
}

resource "qdrant-cloud_accounts_auth_key" "example-key" {
  cluster_ids = [qdrant-cloud_accounts_cluster.example.id]
}

output "cluster_id" {
  value = qdrant-cloud_accounts_cluster.example.id
}

output "cluster_version" {
  value = qdrant-cloud_accounts_cluster.example.version
}

output "url" {
  value = qdrant-cloud_accounts_cluster.example.url
}

output "token" {
  value = qdrant-cloud_accounts_auth_key.example-key.token
}