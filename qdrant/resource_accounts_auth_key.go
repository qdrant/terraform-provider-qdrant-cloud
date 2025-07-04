package qdrant

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	qcAuth "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/auth/v1"
)

// resourceAccountsAuthKey constructs a Terraform resource for managing an API keys associated with an account.
// Returns a schema.Resource pointer configured with schema definitions and the CRUD functions.
func resourceAccountsAuthKey() *schema.Resource {
	return &schema.Resource{
		Description:        "Account AuthKey Resource [Deprecated, see `qdrant-cloud_accounts_database_api_key_v2` instead]",
		ReadContext:        resourceAPIKeyRead,
		CreateContext:      resourceAPIKeyCreate,
		UpdateContext:      nil, // Not available in the public API
		DeleteContext:      resourceAPIKeyDelete,
		Schema:             accountsAuthKeySchema(),
		DeprecationMessage: "The `qdrant-cloud_accounts_auth_key` resource is deprecated and will be removed in a future version. Please use the `qdrant-cloud_accounts_database_api_key_v2` resource instead.",
	}
}

// resourceAPIKeyRead performs a read operation to fetch a specific API key associated with an account.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
func resourceAPIKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error getting API key"
	// Get a client connection and context
	apiClientConn, clientCtx, diagnostics := getClientConnection(ctx, m)
	if diagnostics.HasError() {
		return diagnostics
	}
	// Get a client
	client := qcAuth.NewDatabaseApiKeyServiceClient(apiClientConn)
	// Get The account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Get API Key ID
	apiKeyID := d.Get(authKeysKeysIDFieldName).(string)
	// Execute the request and handle the response
	var trailer metadata.MD
	resp, err := client.ListDatabaseApiKeys(clientCtx, &qcAuth.ListDatabaseApiKeysRequest{
		AccountId: accountUUID.String(),
	}, grpc.Trailer(&trailer))
	// enrich prefix with request ID
	errorPrefix += getRequestID(trailer)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	for _, apiKey := range resp.GetItems() {
		if apiKey.GetId() != apiKeyID {
			// Skip unknown or incorrect ones
			continue
		}
		// Process the correct one,
		for k, v := range flattenAuthKey(apiKey, false) {
			if err := d.Set(k, v); err != nil {
				return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
			}
		}
		d.SetId(apiKeyID)
		return nil
	}
	return diag.Errorf("%s: API key ID cannot be found anymore", errorPrefix)
}

// resourceAPIKeyCreate performs a create operation to generate a new API key associated with an account.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func resourceAPIKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error creating API Key"
	// Get a client connection and context
	apiClientConn, clientCtx, diagnostics := getClientConnection(ctx, m)
	if diagnostics.HasError() {
		return diagnostics
	}
	// Get a client
	client := qcAuth.NewDatabaseApiKeyServiceClient(apiClientConn)
	// Get The account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Prepare the payload for the API request
	var clusterIDs []string
	if clusterIDList, ok := d.GetOk(authKeysKeysClusterIDsFieldName); ok {
		// Prepare the payload for the API request
		clusterIDList := clusterIDList.([]interface{})
		clusterIDs = make([]string, len(clusterIDList))
		for i, v := range clusterIDList {
			clusterIDs[i] = v.(string)
		}
	}
	// Create the request body
	var trailer metadata.MD
	resp, err := client.CreateDatabaseApiKey(clientCtx, &qcAuth.CreateDatabaseApiKeyRequest{
		DatabaseApiKey: &qcAuth.DatabaseApiKey{
			AccountId:  accountUUID.String(),
			ClusterIds: clusterIDs,
		},
	}, grpc.Trailer(&trailer))
	// enrich prefix with request ID
	errorPrefix += getRequestID(trailer)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Flatten cluster and store in Terraform state
	for k, v := range flattenAuthKey(resp.GetDatabaseApiKey(), true) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
		}
	}
	// Set the ID
	d.SetId(resp.GetDatabaseApiKey().GetId())
	return nil
}

// resourceAPIKeyDelete performs a delete operation to remove an API key associated with an account.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
func resourceAPIKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error deleting API Key"
	// Get a client connection and context
	apiClientConn, clientCtx, diagnostics := getClientConnection(ctx, m)
	if diagnostics.HasError() {
		return diagnostics
	}
	// Get a client
	client := qcAuth.NewDatabaseApiKeyServiceClient(apiClientConn)
	// Get The account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Get API Key ID
	apiKeyID := d.Get(authKeysKeysIDFieldName).(string)
	// Delete the key
	var trailer metadata.MD
	_, err = client.DeleteDatabaseApiKey(clientCtx, &qcAuth.DeleteDatabaseApiKeyRequest{
		AccountId:        accountUUID.String(),
		DatabaseApiKeyId: apiKeyID,
	}, grpc.Trailer(&trailer))
	// enrich prefix with request ID
	errorPrefix += getRequestID(trailer)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Clear the resource ID to mark as deleted
	d.SetId("")
	return nil
}
