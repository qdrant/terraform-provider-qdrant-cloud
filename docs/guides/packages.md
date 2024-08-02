---
page_title: "Configure desired resources for Qdrant Cloud cluster"
description: |-
    Guide to getting started with the Qdrant Cloud packages and resources
---

# Configure desired resources for Qdrant Cloud cluster

When creating a cluster, the `package_id` is required. Packages are set of resources. Packages can be queried using the `qdrant-cloud_booking_packages` data source.
Example response looks like:
```json
{
  "currency": "usd",
  "id": "7e1d59b1-6982-4c3c-9cbd-725e6ec29a98",
  "name": "mx2",
  "resource_configurations": [
    {
      "amount": 8,
      "resource_type": "ram",
      "resource_unit": "Gi"
    },
    {
      "amount": 1000,
      "resource_type": "cpu",
      "resource_unit": "m"
    },
    {
      "amount": 32,
      "resource_type": "disk",
      "resource_unit": "Gi"
    }
  ],
  "unit_int_price_per_hour": 9424
}
```

## Example Usage

```terraform
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
  account_id = "" // The default account ID you want to use in Qdrant Cloud (can be overriden on resource level)
}

// Reference the booking packages
data "qdrant-cloud_booking_packages" "all_aws_us_west_2_packages" {
  cloud_provider = "gcp"  // Required. Please refer to the documentation (https://registry.terraform.io/providers/qdrant/qdrant-cloud/latest/docs/guides/getting-started) for the available options.
  cloud_region   = "us-east4"  // Required. Please refer to the documentation (https://registry.terraform.io/providers/qdrant/qdrant-cloud/latest/docs/guides/getting-started) for the available options.
}

locals {
  mx1_packages = [
    for pkg in data.qdrant-cloud_booking_packages.all_aws_us_west_2_packages.packages : pkg
    if pkg.name == "mx1"
  ]
   free_packages = [
    for pkg in data.qdrant-cloud_booking_packages.all_aws_us_west_2_packages.packages : pkg
    if pkg.unit_int_price_per_hour == 0
  ]
  # Find the maximum value
  max_available_ram_on_gcp = max([for pkg in data.qdrant-cloud_booking_packages.all_aws_us_west_2_packages.packages : pkg.resource_configurations[0].amount]...)

   max_memory_packages = [
    for pkg in data.qdrant-cloud_booking_packages.all_aws_us_west_2_packages.packages : pkg
    if pkg.resource_configurations[0].amount == local.max_available_ram_on_gcp
  ]
}

output "mx1_package_id" {
  value = local.mx1_packages[0].id
}

output "free_package_id" {
  value = local.free_packages[0].id
}

output "max_memory_on_gcp_package_id" {
  value = local.max_memory_packages[0].id
}

output "max_memory_on_gcp_package_name" {
  value = local.max_memory_packages[0].name
}

output "max_memory_on_gcp_package_ram_size_in_GiB" {
  value = local.max_memory_packages[0].resource_configurations[0].amount
}

```
