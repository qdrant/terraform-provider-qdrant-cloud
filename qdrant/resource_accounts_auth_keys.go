package qdrant

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceAccountsAuthKeys constructs a Terraform resource for managing the creation, reading, and deletion of API keys associated with an account.
// Returns a schema.Resource pointer configured with schema definitions and the CRuD functions.
func resourceAccountsAuthKeys() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAPIKeyCreate,
		ReadContext:   resourceAPIKeyRead,
		DeleteContext: resourceAPIKeyDelete,
		UpdateContext: nil,
		Schema:        accountsAuthKeysSchema(),
	}
}

// resourceAPIKeyCreate performs a create operation to generate a new API key associated with an account.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func resourceAPIKeyCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(ClientConfig)
	accountID := d.Get("account_id").(string)

	// Prepare the payload for the API request
	var clusterIDs []string
	if clusterIDList, ok := d.GetOk("cluster_id_list"); ok {
		// Prepare the payload for the API request
		clusterIDList := clusterIDList.([]interface{})
		clusterIDs = make([]string, len(clusterIDList))
		for i, v := range clusterIDList {
			clusterIDs[i] = v.(string)
		}
	}

	// Create the request body
	requestBody := ApiKeyCreateRequestBody{
		ClusterIDList: clusterIDs,
	}

	// Using newQdrantCloudRequest to handle HTTP request creation and error checking.
	req, diags := newQdrantCloudRequest(client, "POST", fmt.Sprintf("/accounts/%s/auth/api-keys", accountID), requestBody)
	if diags.HasError() {
		return diags
	}

	// Execute the request and handle the response
	resp, diags := ExecuteRequest(client, req.WithContext(ctx))
	if diags.HasError() {
		return diags
	}

	// Decode the response
	var apiKeysCreate ApiKeyCreate
	if err := json.NewDecoder(resp.Body).Decode(&apiKeysCreate); err != nil {
		return diag.FromErr(err)
	}

	keys := make([]map[string]interface{}, 1)
	keys[0] = map[string]interface{}{
		"id":              apiKeysCreate.ID,
		"created_at":      apiKeysCreate.CreatedAt,
		"user_id":         GetStringOrEmpty(apiKeysCreate.UserID),
		"prefix":          apiKeysCreate.Prefix,
		"cluster_id_list": apiKeysCreate.ClusterIDList,
		"account_id":      apiKeysCreate.AccountID,
		"token":           apiKeysCreate.Token,
	}

	if err := d.Set("keys", keys); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(apiKeysCreate.ID)
	return nil
}

// resourceAPIKeyRead performs a read operation to fetch a specific API key associated with an account.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
func resourceAPIKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(ClientConfig)
	accountID := d.Get("account_id").(string)

	// Using newQdrantCloudRequest to handle HTTP request creation and error checking.
	req, diags := newQdrantCloudRequest(client, "GET", fmt.Sprintf("/accounts/%s/auth/api-keys/%s", accountID, client.ApiKey), nil)
	if diags.HasError() {
		return diags
	}

	// Execute the request and handle the response
	var apiKeysCreate ApiKeyCreate
	resp, diags := ExecuteRequest(client, req.WithContext(ctx))
	if diags.HasError() {
		return diags
	}

	// Decode JSON response into apiKeysCreate struct
	if err := json.NewDecoder(resp.Body).Decode(&apiKeysCreate); err != nil {
		return diag.FromErr(fmt.Errorf("failed to decode response body: %s", err))
	}

	// Set the fields from response
	keys := make([]map[string]interface{}, 1)
	keys[0] = map[string]interface{}{
		"id":              apiKeysCreate.ID,
		"created_at":      apiKeysCreate.CreatedAt,
		"user_id":         GetStringOrEmpty(apiKeysCreate.UserID),
		"prefix":          apiKeysCreate.Prefix,
		"cluster_id_list": apiKeysCreate.ClusterIDList,
		"account_id":      apiKeysCreate.AccountID,
		"token":           apiKeysCreate.Token,
	}

	if err := d.Set("keys", keys); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(apiKeysCreate.ID)

	return nil
}

// resourceAPIKeyDelete performs a delete operation to remove an API key associated with an account.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
func resourceAPIKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(ClientConfig)
	accountID := d.Get("account_id").(string)
	apiKeyID := d.Get("key_id").(string)

	// Using newQdrantCloudRequest to handle HTTP request creation and error checking.
	req, diags := newQdrantCloudRequest(client, "DELETE", fmt.Sprintf("/accounts/%s/auth/api-keys/%s", accountID, apiKeyID), nil)
	if diags.HasError() {
		return diags
	}

	// Execute the request and handle the response
	_, diags = ExecuteRequest(client, req.WithContext(ctx))
	if diags.HasError() {
		return diags
	}

	// Clear the resource ID to mark as deleted
	d.SetId("")

	return nil
}
