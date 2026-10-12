terraform {
  required_version = ">= 1.7.0"
  required_providers {
    qdrant-cloud = {
      source  = "qdrant/qdrant-cloud"
      version = ">=1.1.0"
    }
  }
}

# Local interactive auth: run `qcloud auth login` first, then:
provider "qdrant-cloud" {
  auth       = "cli"
  account_id = "" // The default account ID you want to use in Qdrant Cloud
  # api_url  = "grpc.cloud.qdrant.io:443" // optional; forwarded to qcloud --endpoint
}
