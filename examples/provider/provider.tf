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
  api_key = "" // API Key generated in Qdrant Cloud
  api_url = "" // URL where the public API of Qdrant cloud can be found.
}