package qdrant

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataClusterAccountsCluster constructs a Terraform resource for
// managing the reading of a specific cluster associated with an account.
func dataClusterAccountsCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataClusterAccountsClusterRead,
		Schema:      AccountsClustersSchema(),
	}
}

// dataClusterAccountsClusters constructs a Terraform resource for
// managing the reading of all clusters associated with an account.
func dataClusterAccountsClusters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataClusterAccountsClustersRead,
		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"clusters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: AccountsClustersSchema(),
				},
			},
		},
	}
}

// dataClusterAccountsClusterRead performs a read operation to fetch a specific cluster associated with a given account and cluster ID.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func dataClusterAccountsClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(ClientConfig)
	accountID := d.Get("account_id").(string)
	clusterID := d.Get("cluster_id").(string)

	// Using newQdrantCloudRequest to handle HTTP request creation and error checking.
	req, diags := newQdrantCloudRequest(client, "GET", fmt.Sprintf("/accounts/%s/clusters/%s", accountID, clusterID), nil)
	if diags.HasError() {
		return diags
	}

	// Execute the request and handle the response
	resp, diags := ExecuteRequest(client, req.WithContext(ctx))
	if diags.HasError() {
		return diags
	}

	// Decode JSON response into cluster struct
	var cluster Cluster
	if err := json.NewDecoder(resp.Body).Decode(&cluster); err != nil {
		return diag.FromErr(fmt.Errorf("failed to decode response body: %s", err))
	}

	// Update the Terraform state
	if err := d.Set("cluster", cluster); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(clusterID)
	return nil
}

// dataClusterAccountsClustersRead performs a read operation to fetch all clusters associated with a specific account.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func dataClusterAccountsClustersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(ClientConfig)
	accountID := d.Get("account_id").(string)

	// Using newQdrantCloudRequest to handle HTTP request creation and error checking.
	req, diags := newQdrantCloudRequest(client, "GET", fmt.Sprintf("/accounts/%s/clusters", accountID), nil)
	if diags.HasError() {
		return diags
	}

	// Execute the request and handle the response
	resp, diags := ExecuteRequest(client, req.WithContext(ctx))
	if diags.HasError() {
		return diags
	}

	// Decode JSON response into clusters slice
	var clusters []Cluster
	if err := json.NewDecoder(resp.Body).Decode(&clusters); err != nil {
		return diag.FromErr(fmt.Errorf("failed to decode response body: %s", err))
	}

	// Update the Terraform state
	if err := d.Set("clusters", clusters); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(time.Now().UTC().String())
	return nil
}
