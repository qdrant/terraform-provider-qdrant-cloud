---
page_title: "Getting started with Qdrant Cloud Terraform Provider"
description: |-
    Guide to getting started with the Qdrant Cloud Terraform Provider
---

# Getting started with Qdrant Cloud Terraform Provider
Qdrant Cloud, provides Qdrant databases as a Service (DBaaS). 
It enables you to use the entire functionality of a Qdrant database without the need to run or manage the system yourself.

Terraform Provider Qdrant Cloud is a plugin for Terraform that allows for the full lifecycle management of Qdant Cloud resources.

## Provider Setup

You need to supply proper credentials to the provider before it can be used. 
API keys serve as the credentials to the provider. You can obtain the keys from [Qdrant Cloud console](https://cloud.qdrant.io/).
_Since the provider is in closed beta, you need to [contact Qdrant support](https://support.qdrant.io/support/tickets/new) to get access to the API keys._

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

### Creating your first Cluster (including access token):

```hcl
// Create a cluster (for the sake of having an ID, see below)
resource "qdrant-cloud_accounts_cluster" "example" {
  name           = "tf-example-cluster"
  cloud_provider = "gcp"
  cloud_region   = "us-east4"
  configuration {
    number_of_nodes = 1
    node_configuration {
      package_id = "39b48a76-2a60-4ee0-9266-6d1e0f91ea14" # free
    }
  }
}

// Create an Auth Key, which refers to the cluster provided above
resource "qdrant-cloud_accounts_auth_key" "example-key" {
  cluster_ids = [qdrant-cloud_accounts_cluster.example.id]
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

// Output the token (which can be used to access the database cluster)
output "token" {
  value = qdrant-cloud_accounts_auth_key.example-key.token
}
```


### Available cloud providers and regions
```yaml
aws:
- us-east-1
- us-west-1
- eu-central_lega-y
- eu-central-1
- ap-northeast-1
- ap-southeast-1
- ap-southeast-2
- eu-west-2
- us-west-2
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


