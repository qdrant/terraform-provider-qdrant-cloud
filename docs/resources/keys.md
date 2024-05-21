# Qdrant Cloud Accounts Auth Keys Resource

This resource allows you to manage API keys in your Qdrant Cloud account.

## Example Usage

```terraform
resource "qdrant_accounts_auth_keys" "mykey" {
  account_id = "<your_account_id>"
  cluster_id_list = ["<cluster_id_1>", "<cluster_id_2>"]
}
```

In this example, replace `<your_account_id>` with your actual account ID and `<cluster_id_1>`, `<cluster_id_2>` with your actual cluster IDs. This configuration creates an API key in your Qdrant Cloud account that has access to the specified clusters.

## Argument Reference

The following arguments are supported:

- `account_id` - (Required) The ID of your Qdrant Cloud account.
- `cluster_id_list` - (Optional) A list of cluster IDs that the API key will have access to.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the API key.
- `created_at` - The creation timestamp of the API key.
- `prefix` - The prefix of the API key.
- `token` - The token of the API key. This is marked as sensitive and won't be displayed in the console.

Please note that this is a basic example and the actual values will depend on your specific use case and the options available in your Qdrant Cloud account.
