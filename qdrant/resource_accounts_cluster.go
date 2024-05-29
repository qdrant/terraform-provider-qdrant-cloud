package qdrant

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qc "terraform-provider-qdrant-cloud/v1/internal/client"
)

// resourceAccountsCluster constructs a Terraform resource for
// managing a cluster associated with an account.
// Returns a schema.Resource pointer configured with schema definitions and the CRUD functions.
func resourceAccountsCluster() *schema.Resource {
	return &schema.Resource{
		Description:   "Account Cluster Resource",
		ReadContext:   resourceClusterRead,
		CreateContext: resourceClusterCreate,
		UpdateContext: resourceClusterUpdate,
		DeleteContext: resourceClusterDelete,
		Schema:        accountsClusterSchema(),
	}
}

// resourceClusterRead reads the specific cluster's data from the API.
func resourceClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get an authenticated client
	apiClient, diagnostics := getClient(m)
	if diagnostics.HasError() {
		return diagnostics
	}

	accountID, err := uuid.Parse(d.Get("account_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	clusterID, err := uuid.Parse(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := apiClient.GetClusterWithResponse(ctx, accountID, clusterID)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("error reading cluster: %s", getError(resp.JSON422)))
	}
	if resp.StatusCode() != 200 {
		return diag.FromErr(fmt.Errorf("error reading cluster: [%d] - %s", resp.StatusCode(), resp.Status()))
	}
	clusterOut := resp.JSON200
	if clusterOut == nil {
		return diag.FromErr(fmt.Errorf("error reading cluster: no cluster returned"))
	}
	// Flatten cluster and store in Terraform state
	for k, v := range flattenCluster(clusterOut) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func resourceClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get an authenticated client
	apiClient, diagnostics := getClient(m)
	if diagnostics.HasError() {
		return diagnostics
	}
	// Expand the cluster
	cluster, err := expandClusterIn(d, getDefaultAccountID(m))
	if err != nil {
		return diag.FromErr(err)
	}
	// Create an account UUID from the account ID
	accountUUID, err := uuid.Parse(*cluster.AccountId)
	if err != nil {
		return diag.FromErr(err)
	}
	// Set required version to 'latest' if not provided
	if cluster.Version == nil {
		version := "latest"
		cluster.Version = &version
	}
	// Create the cluster
	resp, err := apiClient.CreateClusterWithResponse(ctx, accountUUID, cluster)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("error creating cluster: %s", getError(resp.JSON422)))
	}
	if resp.StatusCode() != 200 {
		return diag.FromErr(fmt.Errorf("error creating cluster: [%d] - %s", resp.StatusCode(), resp.Status()))
	}

	clusterOut := resp.JSON200
	if clusterOut == nil {
		return diag.FromErr(fmt.Errorf("error creating cluster: no cluster returned"))
	}

	// Flatten cluster and store in Terraform state
	for k, v := range flattenCluster(clusterOut) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(clusterOut.Id)
	return nil
}

func resourceClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get an authenticated client
	apiClient, diagnostics := getClient(m)
	if diagnostics.HasError() {
		return diagnostics
	}
	accountUUID, err := uuid.Parse(d.Get("account_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	clusterUUID, err := uuid.Parse(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("configuration") {
		conf := expandClusterConfigurationIn(d.Get("configuration").([]interface{}))

		resp, err := apiClient.UpdateClusterWithResponse(ctx, accountUUID, clusterUUID, qc.PydanticClusterPatchIn{
			Configuration: &qc.PydanticClusterConfigurationPatchIn{
				NumNodes: &conf.NumNodes,
			},
		})
		if err != nil {
			return diag.FromErr(err)
		}
		if resp.JSON422 != nil {
			return diag.FromErr(fmt.Errorf("error updating cluster: %s", getError(resp.JSON422)))
		}
		if resp.StatusCode() != 200 {
			return diag.FromErr(fmt.Errorf("error updating cluster: [%d] - %s", resp.StatusCode(), resp.Status()))
		}
		clusterOut := resp.JSON200
		if clusterOut == nil {
			return diag.FromErr(fmt.Errorf("error updating cluster: no cluster returned"))
		}
		// Flatten cluster and store in Terraform state
		for k, v := range flattenCluster(clusterOut) {
			if err := d.Set(k, v); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return nil

}

// resourceClusterDelete performs a delete operation to remove a cluster associated with an account.
func resourceClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get an authenticated client
	apiClient, diagnostics := getClient(m)
	if diagnostics.HasError() {
		return diagnostics
	}
	accountUUID, err := uuid.Parse(d.Get("account_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	clusterUUID, err := uuid.Parse(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	deleteBackups := true
	params := &qc.DeleteClusterParams{
		DeleteBackups: &deleteBackups,
	}
	resp, err := apiClient.DeleteClusterWithResponse(ctx, accountUUID, clusterUUID, params)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("error deleting cluster: %s", getError(resp.JSON422)))
	}
	if resp.StatusCode() != 204 {
		return diag.FromErr(fmt.Errorf("error deleting cluster: [%d] - %s", resp.StatusCode(), resp.Status()))
	}
	d.SetId("")
	return nil
}
