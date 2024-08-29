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
		Schema:        accountsClusterSchema(false),
	}
}

// resourceClusterRead reads the specific cluster's data from the API.
func resourceClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error reading cluster"
	// Get an authenticated client
	apiClient, diagnostics := getClient(m)
	if diagnostics.HasError() {
		return diagnostics
	}
	// Get The account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %v", errorPrefix, err))
	}
	// Get the cluster ID as UUID
	clusterUUID, err := uuid.Parse(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %v", errorPrefix, err))
	}
	// Fetch the cluster
	resp, err := apiClient.GetClusterWithResponse(ctx, accountUUID, clusterUUID)
	// Inspect the results
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %v", errorPrefix, err))
	}
	if resp.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("%s: %s", errorPrefix, getError(resp.JSON422)))
	}
	if resp.StatusCode() != 200 {
		return diag.FromErr(fmt.Errorf(getErrorMessage(errorPrefix, resp.HTTPResponse)))
	}
	clusterOut := resp.JSON200
	if clusterOut == nil {
		return diag.FromErr(fmt.Errorf("%s: no cluster returned", errorPrefix))
	}
	// Flatten cluster and store in Terraform state
	for k, v := range flattenCluster(clusterOut) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("%s: %v", errorPrefix, err))
		}
	}

	return nil
}

func resourceClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error creating cluster"
	// Get an authenticated client
	apiClient, diagnostics := getClient(m)
	if diagnostics.HasError() {
		return diagnostics
	}
	// Expand the cluster
	cluster, err := expandCluster(d, getDefaultAccountID(m))
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %v", errorPrefix, err))
	}
	// Set required version to 'latest' if not provided
	if cluster.Configuration.Version == nil || *cluster.Configuration.Version == "" {
		version := "latest"
		cluster.Configuration.Version = &version
	}
	// Create the cluster
	resp, err := apiClient.CreateClusterWithResponse(ctx, cluster.AccountId, nil, cluster)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %v", errorPrefix, err))
	}
	if resp.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("%s: %s", errorPrefix, getError(resp.JSON422)))
	}
	if resp.StatusCode() != 200 {
		return diag.FromErr(fmt.Errorf(getErrorMessage(errorPrefix, resp.HTTPResponse)))
	}

	clusterSchema := resp.JSON200
	if clusterSchema == nil {
		return diag.FromErr(fmt.Errorf("%s: no cluster returned", errorPrefix))
	}
	// Flatten cluster and store in Terraform state
	for k, v := range flattenCluster(clusterSchema) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("%s: %v", errorPrefix, err))
		}
	}
	// Set the ID
	d.SetId(clusterSchema.Id.String())
	return nil
}

func resourceClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error updating cluster"
	// Get an authenticated client
	apiClient, diagnostics := getClient(m)
	if diagnostics.HasError() {
		return diagnostics
	}
	// Get The account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %v", errorPrefix, err))
	}
	// Get the cluster ID as UUID
	clusterUUID, err := uuid.Parse(d.Get(clusterIdentifierFieldName).(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %v", errorPrefix, err))
	}
	// Expand the cluster
	cluster, err := expandCluster(d, getDefaultAccountID(m))
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %v", errorPrefix, err))
	}
	// Update the cluster
	resp, err := apiClient.UpdateClusterWithResponse(ctx, accountUUID, clusterUUID, cluster)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %v", errorPrefix, err))
	}
	if resp.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("%s: %s", errorPrefix, getError(resp.JSON422)))
	}
	if resp.StatusCode() != 200 {
		return diag.FromErr(fmt.Errorf(getErrorMessage(errorPrefix, resp.HTTPResponse)))
	}
	clusterOut := resp.JSON200
	if clusterOut == nil {
		return diag.FromErr(fmt.Errorf("%s: no cluster returned", errorPrefix))
	}
	// Flatten cluster and store in Terraform state
	for k, v := range flattenCluster(clusterOut) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("%s: %v", errorPrefix, err))
		}
	}
	return nil

}

// resourceClusterDelete performs a delete operation to remove a cluster associated with an account.
func resourceClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error deleting cluster"
	// Get an authenticated client
	apiClient, diagnostics := getClient(m)
	if diagnostics.HasError() {
		return diagnostics
	}
	// Get The account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %v", errorPrefix, err))
	}
	// Get the cluster ID as UUID
	clusterUUID, err := uuid.Parse(d.Get(clusterIdentifierFieldName).(string))
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %v", errorPrefix, err))
	}
	// Delete with all backups as well.
	deleteBackups := true
	params := &qc.DeleteClusterParams{
		DeleteBackups: &deleteBackups,
	}
	// Delete the cluster
	resp, err := apiClient.DeleteClusterWithResponse(ctx, accountUUID, clusterUUID, params)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %v", errorPrefix, err))
	}
	if resp.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("%s: %s", errorPrefix, getError(resp.JSON422)))
	}
	if resp.StatusCode() != 204 {
		return diag.FromErr(fmt.Errorf(getErrorMessage(errorPrefix, resp.HTTPResponse)))
	}
	d.SetId("")
	return nil
}
