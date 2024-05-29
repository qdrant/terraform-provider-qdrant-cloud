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

// dataSourceAccountsClusters constructs a Terraform resource for
// managing the reading of all clusters associated with an account.
func dataSourceAccountsClusters() *schema.Resource {
	return &schema.Resource{
		Description: "Account Cluster List Data Source",
		ReadContext: dataSourceAccountsClustersRead,
		Schema:      accountsClustersSchema(),
	}
}

// dataSourceAccountsCluster constructs a Terraform resource for
// managing the reading of a specific cluster associated with an account.
func dataSourceAccountsCluster() *schema.Resource {
	return &schema.Resource{
		Description: "Account Cluster Data Source",
		ReadContext: dataSourceAccountsClusterRead,
		Schema:      accountsClusterSchema(),
	}
}

// dataSourceAccountsClustersRead performs a read operation to fetch all clusters associated with a specific account.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func dataSourceAccountsClustersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get an authenticated client
	apiClient, diagnostics := getClient(m)
	if diagnostics.HasError() {
		return diagnostics
	}
	// Get The account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error listing clusters: %v", err))
	}
	// List all clusters for the provided account
	resp, err := apiClient.ListClustersWithResponse(ctx, accountUUID, &qc.ListClustersParams{})
	if err != nil {
		d := diag.FromErr(fmt.Errorf("error listing clusters: %v", err))
		if d.HasError() {
			return d
		}
	}
	// Inspect result and get the resulting clusters
	if resp.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("error listing clusters: %s", getError(resp.JSON422)))
	}
	if resp.StatusCode() != 200 {
		return diag.FromErr(fmt.Errorf("error listing clusters: [%d] - %s", resp.StatusCode(), resp.Status()))
	}
	clustersOut := resp.JSON200
	if clustersOut == nil {
		return diag.FromErr(fmt.Errorf("error listing clusters: ListCluster didn't return clusters"))
	}
	// Update the Terraform state (TODO: Flatten)
	if err := d.Set("clusters", clustersOut); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(time.Now().UTC().String())
	return nil
}

// dataSourceAccountsClusterRead performs a read operation to fetch a specific cluster associated with a given account and cluster ID.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func dataSourceAccountsClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get an authenticated client
	apiClient, diagnostics := getClient(m)
	if diagnostics.HasError() {
		return diagnostics
	}
	// Get The account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting cluster: %v", err))
	}
	// Get the cluster ID as UUID
	clusterID := d.Get("cluster_id").(string)
	clusterUUID, err := uuid.Parse(clusterID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error parsing cluster ID: %v", err))
	}
	// Fetch the cluster
	resp, err := apiClient.GetClusterWithResponse(ctx, accountUUID, clusterUUID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting cluster: %v", err))
	}
	// Inspect result and get the resulting cluster
	if resp.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("error getting cluster: %s", getError(resp.JSON422)))
	}
	if resp.StatusCode() != 200 {
		return diag.FromErr(fmt.Errorf("error getting cluster: [%d] - %s", resp.StatusCode(), resp.Status()))
	}
	clusterOut := resp.JSON200
	if clusterOut == nil {
		return diag.FromErr(fmt.Errorf("ListCluster didn't return clusters"))
	}
	// Flatten cluster and store in Terraform state
	for k, v := range flattenCluster(clusterOut) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}
