package qdrant

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	qc "terraform-provider-qdrant-cloud/v1/internal/client"
	"time"

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
	accountID, err := uuid.Parse(d.Get("account_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

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
	requestBody := qc.ApiKeyIn{
		ClusterIdList: &clusterIDs,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return diag.FromErr(err)
	}

	body := bytes.NewReader(jsonData)

	apiClient, err, diagnostics, done := GetClient(m)
	if done {
		return diagnostics
	}

	resp, err := apiClient.CreateApiKeyWithBodyWithResponse(ctx, accountID, "application/json", body)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("error creating API key: %v", resp.JSON422))
	}

	if err := d.Set("keys", resp.JSON200); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(time.Now().Format(time.RFC3339))
	return nil
}

// resourceAPIKeyRead performs a read operation to fetch a specific API key associated with an account.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
func resourceAPIKeyRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountID, err := uuid.Parse(d.Get("account_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	apiClient, err, diagnostics, done := GetClient(m)
	if done {
		return diagnostics
	}

	// Execute the request and handle the response
	resp, err := apiClient.ListApiKeysWithResponse(ctx, accountID)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("error listing API keys: %v", resp.JSON422))
	}

	err = d.Set("keys", resp.JSON200)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(time.Now().Format(time.RFC3339))

	return nil
}

// resourceAPIKeyDelete performs a delete operation to remove an API key associated with an account.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
func resourceAPIKeyDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountID, err := uuid.Parse(d.Get("account_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	apiKeyID := d.Get("key_id").(string)

	apiClient, err, diagnostics, done := GetClient(m)
	if done {
		return diagnostics
	}

	resp, err := apiClient.DeleteApiKeyWithResponse(ctx, accountID, apiKeyID)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("error deleting API key: %v", resp.JSON422))
	}

	// Clear the resource ID to mark as deleted
	d.SetId("")

	return nil
}
