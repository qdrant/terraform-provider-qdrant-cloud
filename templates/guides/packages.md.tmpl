---
page_title: "Configure desired resources for Qdrant Cloud cluster"
description: |-
    Guide to getting started with the Qdrant Cloud packages and resources
weight: 200
---

# Configure desired resources for Qdrant Cloud cluster

When creating a cluster, the `package_id` is required. Packages are set of resources. Packages can be queried using the `qdrant-cloud_booking_packages` data source.
Example response looks like:
```json
{
  "currency": "USD",
  "id":  "2b194353-30d4-4c76-8b4d-4d7dd04a82a6",
  "name": "mx6",
  "resource_configuration": [
    {
      "cpu": "16000m",
      "disk": "512Gi",
      "ram": "128Gi"
    }
  ],
  "status": "PACKAGE_STATUS_ACTIVE",
  "type": "paid",
  "unit_int_price_per_hour": 167168
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
data "qdrant-cloud_booking_packages" "all_aws_gcp_us_east4_packages" {
  cloud_provider = "gcp"  // Required. Please refer to the documentation (https://registry.terraform.io/providers/qdrant/qdrant-cloud/latest/docs/guides/getting-started) for the available options.
  cloud_region   = "us-east4"  // Required. Please refer to the documentation (https://registry.terraform.io/providers/qdrant/qdrant-cloud/latest/docs/guides/getting-started) for the available options.
}

locals {
  desired_package = [
    for pkg in data.qdrant-cloud_booking_packages.all_packages.packages : pkg
    if pkg.resource_configuration[0].cpu == "16000m" && pkg.resource_configuration[0].ram == "64Gi"
  ]
}

output "desired_package_id" {
  value = local.desired_package[0].id
}

```
