---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "qdrant-cloud_accounts_database_api_keys_v2 Data Source - terraform-provider-qdrant-cloud"
subcategory: ""
description: |-
  Account Database API Keys Data Source (v2)
---

# qdrant-cloud_accounts_database_api_keys_v2 (Data Source)

Account Database API Keys Data Source (v2)



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cluster_id` (String) Database API Keys V2 Schema Cluster Identifier for which this Database API Key is attached field

### Optional

- `account_id` (String) Database API Keys V2 Schema Account Identifier where all those Database API Keys belongs to field

### Read-Only

- `id` (String) The ID of this resource.
- `keys` (List of Object) Database API Keys V2 Schema List of Database API Keys field (see [below for nested schema](#nestedatt--keys))

<a id="nestedatt--keys"></a>
### Nested Schema for `keys`

Read-Only:

- `account_id` (String)
- `cluster_id` (String)
- `collection_access_rules` (List of Object) (see [below for nested schema](#nestedobjatt--keys--collection_access_rules))
- `created_at` (String)
- `created_by_email` (String)
- `expires_at` (String)
- `global_access_rule` (List of Object) (see [below for nested schema](#nestedobjatt--keys--global_access_rule))
- `id` (String)
- `key` (String)
- `name` (String)
- `postfix` (String)

<a id="nestedobjatt--keys--collection_access_rules"></a>
### Nested Schema for `keys.collection_access_rules`

Read-Only:

- `access_type` (String)
- `collection_name` (String)
- `payload` (Map of String)


<a id="nestedobjatt--keys--global_access_rule"></a>
### Nested Schema for `keys.global_access_rule`

Read-Only:

- `access_type` (String)
