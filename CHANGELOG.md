## 1.3.x 

FEATURES:

1. **Terraform Provider for Qdrant Cloud**: A Terraform provider has been created for Qdrant Cloud, which allows users to manage their Qdrant Cloud resources using Terraform. This is implemented in the `provider.go` file.

2. **API Key Management**: Users can create, read, and delete API keys associated with their Qdrant Cloud account. This is implemented in the `resource_accounts_auth_keys.go` file.

3. **Cluster Management**: Users can manage clusters associated with their Qdrant Cloud account. This is implemented in the `resource_accounts_clusters.go` file.

4. **HTTP Request Handling**: The provider includes functionality for creating and executing HTTP requests to the Qdrant API, including error handling. This is implemented in the `provider.go` file.

5. **Schema Definitions**: The provider includes schema definitions for the Qdrant Cloud resources it manages, including API keys and clusters. These are implemented in the `schemas_accounts_auth_keys.go` and `schemas_accounts_clusters.go` files.

6. **Testing**: The provider includes tests for its functionality, which are implemented in the `provider_test.go`, `resource_accounts_auth_keys_test.go`, and `resource_accounts_clusters_test.go` files.

7. **Time Formatting**: The provider includes a function for formatting time values, which is implemented in the `provider.go` file.

8. **Client Configuration**: The provider includes a function for configuring the HTTP client used to make requests to the Qdrant API, which is implemented in the `provider.go` file.

9. **Data Sources**: The provider includes data sources for retrieving Qdrant Cloud accounts' authorization keys, listing Qdrant Cloud clusters under an account, retrieving details of a specific Qdrant Cloud cluster, and Qdrant Cloud booking packages. These are implemented in the `provider.go` file.

10. **Resource Configuration**: The provider includes functionality for configuring the resources it manages, including setting the number of nodes and node configuration for clusters. This is implemented in the `resource_accounts_clusters.go` file.

## 1.4.x

FEATURES:

1. Replaced HTTP client with gRPC client, which leverage the new Qdrant public API

## 1.5.x

FEATURES:

1. **API Key Management V2**: Users can create, read, and delete v2 API keys with granular access controls for their Qdrant Cloud account. This is implemented in `resource_accounts_auth_key_v2.go` and `data_source_accounts_auth_keys_v2.go`.

TESTS:
1. **V2 API Key Acceptance Tests**: Added acceptance tests for the v2 API key resource and data source to ensure they function correctly against the live Qdrant Cloud API.
2. **V2 API Key Unit Tests**: Added unit tests for the v2 API key schema and flattening logic.
