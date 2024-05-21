package qdrant

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataAccountsAuthKeys constructs a Terraform resource for managing the reading of API keys associated with an account.
// Returns a schema.Resource pointer configured with schema definitions and the read context function.
func dataAccountsAuthKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataAccountsAuthKeysRead,
		Schema:      accountsAuthKeysSchema(),
	}
}

// dataAccountsAuthKeysRead performs a read operation to fetch all API keys associated with a specific account.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func dataAccountsAuthKeysRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(ClientConfig)
	accountID := d.Get("account_id").(string)

	// Using newQdrantCloudRequest to handle HTTP request creation and error checking.
	req, diags := newQdrantCloudRequest(client, "GET", fmt.Sprintf("/accounts/%s/auth/api-keys", accountID), nil)
	if diags.HasError() {
		return diags
	}

	// Execute the request and handle the response
	resp, diags := ExecuteRequest(client, req.WithContext(ctx))
	if diags.HasError() {
		return diags
	}

	// Decode JSON response into apiKeys slice
	var apiKeys []APIKey
	if err := json.NewDecoder(resp.Body).Decode(&apiKeys); err != nil {
		return diag.FromErr(fmt.Errorf("failed to decode response body: %s", err))
	}

	keys := make([]map[string]interface{}, len(apiKeys))

	for i, key := range apiKeys {
		// Create a map to hold the attributes for the current key
		keyMap := make(map[string]interface{})

		// Set the attributes for the current key
		keyMap["created_at"] = key.CreatedAt
		keyMap["prefix"] = key.Prefix
		//TODO: keyMap["token"] = key.Token
		keyMap["user_id"] = key.UserID
		keyMap["account_id"] = key.AccountID

		// Add the map to the keys slice
		keys[i] = keyMap
	}

	// Set the keys attribute in the Terraform state
	if err := d.Set("keys", keys); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(accountID)
	return nil
}
