package qdrant

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qc "github.com/qdrant/terraform-provider-qdrant-cloud/v1/internal/client"
)

// resourceAccountsAuthKey constructs a Terraform resource for managing an API keys associated with an account.
// Returns a schema.Resource pointer configured with schema definitions and the CRUD functions.
func resourceAccountsAuthKey() *schema.Resource {
	return &schema.Resource{
		Description:   "Account AuthKey Resource",
		ReadContext:   resourceAPIKeyRead,
		CreateContext: resourceAPIKeyCreate,
		UpdateContext: nil, // Not available in the public API
		DeleteContext: resourceAPIKeyDelete,
		Schema:        accountsAuthKeySchema(),
	}
}

// resourceAPIKeyRead performs a read operation to fetch a specific API key associated with an account.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
func resourceAPIKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error getting API key"
	// Get an authenticated client
	apiClient, diagnostics := getClient(m)
	if diagnostics.HasError() {
		return diagnostics
	}
	// Get The account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	apiKeyID := d.Get(authKeysKeysIDFieldName).(string)
	// convert to UUID
	apiKeyUUID, err := uuid.Parse(apiKeyID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Execute the request and handle the response
	resp, err := apiClient.ListApiKeysWithResponse(ctx, accountUUID, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	if resp.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("%s: %s", errorPrefix, getError(resp.JSON422)))
	}
	if resp.StatusCode() != 200 {
		return diag.FromErr(fmt.Errorf("%s", getErrorMessage(errorPrefix, resp.HTTPResponse)))
	}
	// Get the actual response
	apiKeys := resp.JSON200
	if apiKeys == nil {
		return diag.FromErr(fmt.Errorf("%s: no keys returned", errorPrefix))
	}
	for _, apiKey := range *apiKeys {
		if apiKey.Id == nil || *apiKey.Id != apiKeyUUID {
			// Skip unknown or incorrect ones
			continue
		}
		// Process the correct one,
		for k, v := range flattenAuthKey(apiKey) {
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
	// Get an authenticated client
	apiClient, diagnostics := getClient(m)
	if diagnostics.HasError() {
		return diagnostics
	}
	// Get The account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Prepare the payload for the API request
	var clusterIDs []uuid.UUID
	if clusterIDList, ok := d.GetOk(authKeysKeysClusterIDsFieldName); ok {
		// Prepare the payload for the API request
		clusterIDList := clusterIDList.([]interface{})
		clusterIDs = make([]uuid.UUID, len(clusterIDList))
		for i, v := range clusterIDList {
			clusterIDs[i] = uuid.MustParse(v.(string))
		}
	}
	// Create the request body
	resp, err := apiClient.CreateApiKeyWithResponse(ctx, accountUUID, qc.ApiKeySchema{
		ClusterIds: clusterIDs,
	})
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	if resp.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("%s: %s", errorPrefix, getError(resp.JSON422)))
	}
	if resp.StatusCode() != 200 {
		return diag.FromErr(fmt.Errorf("%s", getErrorMessage(errorPrefix, resp.HTTPResponse)))
	}
	// Get the actual response
	apiKey := resp.JSON200
	if apiKey == nil || apiKey.Id == nil {
		return diag.FromErr(fmt.Errorf("%s: no keys returned", errorPrefix))
	}
	// Flatten cluster and store in Terraform state
	for k, v := range flattenAuthKey(*apiKey) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
		}
	}
	// Set the ID
	d.SetId(apiKey.Id.String())
	return nil
}

// resourceAPIKeyDelete performs a delete operation to remove an API key associated with an account.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
func resourceAPIKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error deleting API Key"
	// Get an authenticated client
	apiClient, diagnostics := getClient(m)
	if diagnostics.HasError() {
		return diagnostics
	}
	// Get The account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	apiKeyID := d.Get(authKeysKeysIDFieldName).(string)

	resp, err := apiClient.DeleteApiKeyWithResponse(ctx, accountUUID, apiKeyID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	if resp.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("%s: %s", errorPrefix, getError(resp.JSON422)))
	}
	if resp.StatusCode() != 204 {
		return diag.FromErr(fmt.Errorf("%s", getErrorMessage(errorPrefix, resp.HTTPResponse)))
	}
	// Clear the resource ID to mark as deleted
	d.SetId("")
	return nil
}
