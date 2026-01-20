package qdrant

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	qcCluster "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/v1"
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
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// resourceClusterRead reads the specific cluster's data from the API.
func resourceClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error reading cluster"
	client, clientCtx, diags := getServiceClient(ctx, m, qcCluster.NewClusterServiceClient)
	if diags.HasError() {
		return diags
	}
	// Get The account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Fetch the cluster
	var trailer metadata.MD
	resp, err := client.GetCluster(clientCtx, &qcCluster.GetClusterRequest{
		AccountId: accountUUID.String(),
		ClusterId: d.Id(),
	}, grpc.Trailer(&trailer))
	// enrich prefix with request ID
	errorPrefix += getRequestID(trailer)
	// Inspect the results
	if err != nil {
		// If the cluster is not found, it might have been deleted manually.
		// Remove it from the state.
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Flatten cluster and store in Terraform state
	for k, v := range flattenCluster(resp.GetCluster()) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
		}
	}
	return nil
}

func resourceClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error creating cluster"
	client, clientCtx, diags := getServiceClient(ctx, m, qcCluster.NewClusterServiceClient)
	if diags.HasError() {
		return diags
	}
	// Expand the cluster
	cluster, err := expandCluster(d, getDefaultAccountID(m))
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Create the cluster
	var trailer metadata.MD
	resp, err := client.CreateCluster(clientCtx, &qcCluster.CreateClusterRequest{
		Cluster: cluster,
	}, grpc.Trailer(&trailer))
	reqID := getRequestID(trailer)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.InvalidArgument {
			return diag.Errorf("Invalid argument for cluster creation%s: %s", reqID, st.Message())
		}
		return diag.FromErr(fmt.Errorf("%s%s: %w", errorPrefix, reqID, err))
	}
	// Flatten cluster and store in Terraform state
	for k, v := range flattenCluster(resp.GetCluster()) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
		}
	}
	// Set the ID
	d.SetId(resp.GetCluster().GetId())
	return nil
}

func resourceClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error updating cluster"
	client, clientCtx, diags := getServiceClient(ctx, m, qcCluster.NewClusterServiceClient)
	if diags.HasError() {
		return diags
	}
	// Expand the cluster
	cluster, err := expandCluster(d, getDefaultAccountID(m))
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Do not provide state in update
	cluster.State = nil
	// Update the cluster
	var trailer metadata.MD
	resp, err := client.UpdateCluster(clientCtx, &qcCluster.UpdateClusterRequest{
		Cluster: cluster,
	}, grpc.Trailer(&trailer))
	reqID := getRequestID(trailer)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.InvalidArgument {
			return diag.Errorf("Invalid argument for cluster update%s: %s", reqID, st.Message())
		}
		return diag.FromErr(fmt.Errorf("%s%s: %w", errorPrefix, reqID, err))
	}
	// Flatten cluster and store in Terraform state
	for k, v := range flattenCluster(resp.GetCluster()) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
		}
	}
	return nil

}

// resourceClusterDelete performs a delete operation to remove a cluster associated with an account.
func resourceClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error deleting cluster"
	client, clientCtx, diags := getServiceClient(ctx, m, qcCluster.NewClusterServiceClient)
	if diags.HasError() {
		return diags
	}
	// Get The account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Do we need to delete the backups?
	deleteBackups := d.Get(clusterDeleteBackupsOnDestroyFieldName).(bool)
	// Delete the cluster
	var trailer metadata.MD
	_, err = client.DeleteCluster(clientCtx, &qcCluster.DeleteClusterRequest{
		AccountId:     accountUUID.String(),
		ClusterId:     d.Id(),
		DeleteBackups: newPointer(deleteBackups),
	}, grpc.Trailer(&trailer))
	// enrich prefix with request ID
	errorPrefix += getRequestID(trailer)
	if err != nil {
		// If the cluster is not found, it has been deleted already.
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	d.SetId("")
	return nil
}
