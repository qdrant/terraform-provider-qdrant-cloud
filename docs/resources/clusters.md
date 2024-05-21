An example of how to use the `qdrant_cluster` resource in a Terraform configuration:

```terraform
resource "qdrant_cluster" "example" {
  account_id = "<your_account_id>"
  name = "example-cluster"
  cloud_provider = "aws"
  cloud_region = "us-west-2"
  configuration {
    num_nodes = 3
    num_nodes_max = 5
    node_configuration {
      package_id = "example-package-id"
    }
  }
}
```

In this example, replace `<your_account_id>` with your actual account ID. This configuration creates a Qdrant Cloud cluster in the AWS `us-west-2` region with the name `example-cluster`. The cluster configuration specifies that the cluster should have 3 nodes initially and can scale up to a maximum of 5 nodes. The `package_id` for the node configuration is set to `example-package-id`.

Please note that this is a basic example and the actual values will depend on your specific use case and the options available in your Qdrant Cloud account.
