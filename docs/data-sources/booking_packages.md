---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "qdrant-cloud_booking_packages Data Source - terraform-provider-qdrant-cloud"
subcategory: ""
description: |-
  Booking packages Data Source
---

# qdrant-cloud_booking_packages (Data Source)

Booking packages Data Source



<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `packages` (List of Object) TODO (see [below for nested schema](#nestedatt--packages))

<a id="nestedatt--packages"></a>
### Nested Schema for `packages`

Read-Only:

- `currency` (String)
- `id` (String)
- `name` (String)
- `regional_mapping_id` (String)
- `resource_configuration` (List of Object) (see [below for nested schema](#nestedobjatt--packages--resource_configuration))
- `status` (Number)
- `unit_int_price_per_day` (Number)
- `unit_int_price_per_hour` (Number)
- `unit_int_price_per_month` (Number)
- `unit_int_price_per_year` (Number)

<a id="nestedobjatt--packages--resource_configuration"></a>
### Nested Schema for `packages.resource_configuration`

Read-Only:

- `resource_option` (List of Object) (see [below for nested schema](#nestedobjatt--packages--resource_configuration--resource_option))
- `resource_option_id` (String)

<a id="nestedobjatt--packages--resource_configuration--resource_option"></a>
### Nested Schema for `packages.resource_configuration.resource_option`

Read-Only:

- `currency` (String)
- `id` (String)
- `name` (String)
- `resource_type` (String)
- `resource_unit` (String)
- `status` (Number)
- `unit_int_price_per_day` (Number)
- `unit_int_price_per_hour` (Number)
- `unit_int_price_per_month` (Number)
- `unit_int_price_per_year` (Number)