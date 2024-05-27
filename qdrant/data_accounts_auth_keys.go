package qdrant

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qc "terraform-provider-qdrant-cloud/v1/internal/client"
)

// dataAccountsAuthKeys constructs a Terraform resource for managing the reading of API keys associated with an account.
// Returns a schema.Resource pointer configured with schema definitions and the read context function.
func dataAccountsAuthKeys() *schema.Resource {
	return &schema.Resource{
		Description: "Account AuthKey Data Source",
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
	// Get an authenticated client
	apiClient, diagnostics := getClient(m)
	if diagnostics.HasError() {
		return diagnostics
	}

	accountID, err := uuid.Parse(d.Get("account_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	response, err := apiClient.ListApiKeysWithResponse(ctx, accountID)
	if err != nil {
		d := diag.FromErr(fmt.Errorf("error listing packages: %v", err))
		if d.HasError() {
			return d
		}
	}

	if response.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("error listing packages: %v", response.JSON422))
	}

	// Set the keys attribute in the Terraform state
	if err := d.Set("keys", flattenAuthKeys(*response.JSON200)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(time.Now().Format(time.RFC3339))
	return nil
}

// flattenAuthKeys flattens the API keys response into a slice of map[string]interface{}.
func flattenAuthKeys(keys []qc.GetApiKeyOut) []interface{} {
	var flattenedKeys []interface{}
	for _, key := range keys {
		flattenedKeys = append(flattenedKeys, map[string]interface{}{
			"id":              key.Id,
			"account_id":      key.AccountId,
			"cluster_id_list": key.ClusterIdList,
		})
	}
	return flattenedKeys
}
