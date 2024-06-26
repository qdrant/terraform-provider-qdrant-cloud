---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "qdrant-cloud_accounts_cluster Resource - terraform-provider-qdrant-cloud"
subcategory: ""
description: |-
  Account Cluster Resource
---

# qdrant-cloud_accounts_cluster (Resource)

Account Cluster Resource

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cloud_provider` (String) Cluster Schema Cloud provider where the cluster resides field
- `cloud_region` (String) Cluster Schema Cloud region where the cluster resides field
- `configuration` (Block List, Min: 1, Max: 1) Cluster Schema The configuration options of a cluster field (see [below for nested schema](#nestedblock--configuration))
- `name` (String) Cluster Schema Name of the cluster field

### Optional

- `account_id` (String) Cluster Schema Identifier of the account field
- `cloud_region_az` (String) Cluster Schema Cloud region availability zone where the cluster resides field
- `cloud_region_setup` (String) Cluster Schema Cloud region setup of the cluster field
- `encryption_key_id` (String) Cluster Schema Identifier of the encryption key field
- `private_region_id` (String) Cluster Schema Identifier of the Private Region field
- `total_extra_disk` (Number) Cluster Schema The total amount of extra disk in relation to the chosen package (in GiB) field
- `version` (String) Cluster Schema Version of the Qdrant cluster field

### Read-Only

- `created_at` (String) Cluster Schema Timestamp when the cluster is created field
- `current_configuration_id` (String) Cluster Schema Identifier of the current configuration field
- `id` (String) Cluster Schema Identifier of the cluster field
- `marked_for_deletion_at` (String) Cluster Schema Timestamp when this cluster was marked for deletion field
- `url` (String) Cluster Schema The URL of the endpoint of the Qdrant cluster field

<a id="nestedblock--configuration"></a>
### Nested Schema for `configuration`

Required:

- `node_configuration` (Block List, Min: 1, Max: 1) Cluster Schema The node configuration options of a cluster field (see [below for nested schema](#nestedblock--configuration--node_configuration))
- `num_nodes` (Number) Cluster Schema The number of nodes in the cluster field
- `num_nodes_max` (Number) Cluster Schema The maximum number of nodes in the cluster field

<a id="nestedblock--configuration--node_configuration"></a>
### Nested Schema for `configuration.node_configuration`

Required:

- `package_id` (String) Cluster Schema The package identifier (specifying: CPU, Memory, and disk size) field
