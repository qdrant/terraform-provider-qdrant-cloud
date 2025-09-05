---
page_title: "Getting started with Qdrant Cloud Terraform Provider"
description: |-
    Guide to getting started with the Qdrant Cloud Terraform Provider
weight: 200
---

# Getting started with Qdrant Cloud Terraform Provider
Qdrant Cloud, provides Qdrant databases as a Service (DBaaS).
It enables you to use the entire functionality of a Qdrant database without the need to run or manage the system yourself.

Terraform Provider Qdrant Cloud is a plugin for Terraform that allows for the full lifecycle management of Qdant Cloud resources.

## Provider Setup

You need to supply proper credentials to the provider before it can be used.
API keys serve as the credentials to the provider. You can obtain the keys from [Qdrant Cloud console](https://cloud.qdrant.io/).

```hcl
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

```

## Example Usage

### Create your first Cluster (including Database Key):

```hcl
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
  name           = "tf-example-cluster"
  cloud_provider = data.qdrant-cloud_booking_packages.all_packages.cloud_provider
  cloud_region   = data.qdrant-cloud_booking_packages.all_packages.cloud_region
  configuration {
    number_of_nodes = 1
    database_configuration {
      service {
        jwt_rbac = true
      }
    }
    node_configuration {
      package_id = local.desired_package[0].id
    }
  }
}

// Create an V2 Database Key, which refers to the cluster provided above
resource "qdrant-cloud_accounts_database_api_key_v2" "example-key" {
  cluster_id   = qdrant-cloud_accounts_cluster.example.id
  name         = "example-key"
}

```

### Show the output of the 2 resources created

```hcl
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

// Output the Database API Key (which can be used to access the database cluster)
output "key" {
  value       = qdrant-cloud_accounts_database_api_key_v2.example-key.key
}
```

### Available cloud providers and regions
```yaml
aws:
- us-east-1
- us-west-1
- us-west-2
- ap-northeast-1
- ap-southeast-1
- ap-southeast-2
- ap-south-1
- eu-central-1
- eu-west-1
- eu-west-2
gcp:
- europe-west3
- us-east4
private:
- private
azure:
- eastus
- germanywestcentral
- southeastasia
- uksouth

```

