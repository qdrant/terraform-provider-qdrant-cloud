# Qdrant Cloud Terraform Provider

[![Terraform Registry](https://img.shields.io/badge/Terraform-Registry-blue.svg)](https://registry.terraform.io/providers/qdrant/qdrant-cloud)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Welcome to the official Terraform provider for Qdrant Cloud!

[Qdrant](https://qdrant.tech/) is a vector similarity search engine that provides a production-ready service with a convenient API to store, search, and manage points—vectors with an additional payload. You can use it to extract meaningful information from unstructured data.

[Qdrant Cloud](https://qdrant.tech/cloud/) is the DBaaS solution for Qdrant, offering a fully managed service that lets you focus on your applications without worrying about the underlying infrastructure.

This Terraform provider allows you to manage your Qdrant Cloud resources programmatically, making it easy to integrate Qdrant into your infrastructure as code (IaC) workflows.

## Getting Started

To get started with the Qdrant Cloud Terraform provider, you'll need to have a Qdrant Cloud account. If you don't have one, you can sign up for a free account at [cloud.qdrant.io](https://cloud.qdrant.io/).

Next, you'll need to create an API key. You can find instructions on how to do that in the [Qdrant Cloud documentation](https://qdrant.tech/documentation/cloud/authentication/).

Once you have your API key, you can configure the provider in your Terraform code:

```hcl
# see: https://registry.terraform.io/providers/qdrant/qdrant-cloud/latest/docs/guides/getting-started
# Setup Terraform, including the qdrant-cloud providers
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
  api_key    = "<QDRANT_CLOUD_MANAGEMENT_KEY>"  // API Key generated in Qdrant Cloud (required)
  account_id = "<QDRANT_CLOUD_ACCOUNT_ID>"      // The default account ID you want to use in Qdrant Cloud (can be overriden on resource level)
}
```

## Example Usage

Here is an example of how you can create a Qdrant Cloud cluster and some common steps to follow:

```hcl
// Add the provider to specify some provider wide settings
provider "qdrant-cloud" {
  api_key    = "<QDRANT_CLOUD_MANAGEMENT_KEY>"  // API Key generated in Qdrant Cloud (required)
  account_id = "<QDRANT_CLOUD_ACCOUNT_ID>"      // The default account ID you want to use in Qdrant Cloud (can be overriden on resource level)
}

// Get the cluster package
// see https://registry.terraform.io/providers/qdrant/qdrant-cloud/latest/docs/guides/getting-started#available-cloud-providers-and-regions
data "qdrant-cloud_booking_packages" "all_packages" {
  cloud_provider = "aws"
  cloud_region   = "us-west-2"
}

locals {
  desired_package = [
    for pkg in data.qdrant-cloud_booking_packages.all_packages.packages : pkg
    if pkg.resource_configuration[0].cpu == "500m" && pkg.resource_configuration[0].ram == "2Gi"
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
resource "qdrant-cloud_accounts_database_api_key_v2" "example" {
  cluster_id   = qdrant-cloud_accounts_cluster.example.id
  name         = "example-key"
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

// Output the Database API Key (which can be used to access the database cluster)
output "key" {
  value       = qdrant-cloud_accounts_database_api_key_v2.example.key
  description = "Key is available only once, after creation."
}
```
You will find all the documentation available, including more usage examples in the [provider documentation](https://qdrant.tech/documentation/cloud-tools/terraform/). Feel free to check it out and send us feedback so we can keep improving it.


## Schema Reference

The following resources are available in the Qdrant Cloud Terraform provider:

*   `qdrant-cloud_accounts_auth_key`
*   `qdrant-cloud_accounts_backup_schedule`
*   `qdrant-cloud_accounts_cluster`
*   `qdrant-cloud_accounts_database_api_key_v2`
*   `qdrant-cloud_accounts_hybrid_cloud_environment`
*   `qdrant-cloud_accounts_manual_backup`
*   `qdrant-cloud_accounts_role`
*   `qdrant-cloud_accounts_user_roles`

Additionally, the following data sources are also available:

*   `qdrant-cloud_accounts_auth_keys`
*   `qdrant-cloud_accounts_backup_schedule`
*   `qdrant-cloud_accounts_backup_schedules`
*   `qdrant-cloud_accounts_cluster`
*   `qdrant-cloud_accounts_clusters`
*   `qdrant-cloud_accounts_database_api_keys_v2`
*   `qdrant-cloud_booking_packages`

For more information on the available resources and their configuration options, please refer to the [provider documentation](https://qdrant.tech/documentation/cloud-tools/terraform/).

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md).

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.
