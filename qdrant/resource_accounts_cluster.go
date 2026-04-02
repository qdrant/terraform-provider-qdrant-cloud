package qdrant

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	qcCluster "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/v1"
)

const (
	clusterCreatePollInterval = 10 * time.Second
	clusterCreateTimeout      = 20 * time.Minute
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(clusterCreateTimeout),
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
	cluster, jwtRbac, err := expandCluster(d, getDefaultAccountID(m))
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	if jwtRbac != nil {
		// If jwtRbac is not nil, we need to add the value ("true" or "false") to the gRPC context with key "qc-jwt-rbac"
		clientCtx = metadata.AppendToOutgoingContext(clientCtx, "qc-jwt-rbac", fmt.Sprintf("%t", *jwtRbac))
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

	createdCluster := resp.GetCluster()
	d.SetId(createdCluster.GetId())

	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	stateConf := &retry.StateChangeConf{
		Pending: []string{
			clusterWaitPending,
		},
		Target: []string{
			clusterWaitReady,
		},
		Refresh:      clusterEndpointRefreshFunc(client, clientCtx, accountUUID.String(), createdCluster.GetId()),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		PollInterval: clusterCreatePollInterval,
	}

	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	readyCluster := result.(*qcCluster.Cluster)
	for k, v := range flattenCluster(readyCluster) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
		}
	}
	return nil
}

const (
	clusterWaitPending = "waiting"
	clusterWaitReady   = "ready"
)

// clusterEndpointRefreshFunc returns a StateRefreshFunc that polls GetCluster
// until the endpoint URL is populated or creation fails.
func clusterEndpointRefreshFunc(
	client qcCluster.ClusterServiceClient,
	ctx context.Context,
	accountID, clusterID string,
) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.GetCluster(ctx, &qcCluster.GetClusterRequest{
			AccountId: accountID,
			ClusterId: clusterID,
		})
		if err != nil {
			return nil, "", err
		}

		cluster := resp.GetCluster()
		phase := cluster.GetState().GetPhase()

		if phase == qcCluster.ClusterPhase_CLUSTER_PHASE_FAILED_TO_CREATE {
			return nil, "", fmt.Errorf("cluster creation failed (phase=%q reason=%q)",
				phase.String(), cluster.GetState().GetReason())
		}

		url := strings.TrimSpace(cluster.GetState().GetEndpoint().GetUrl())
		if url != "" {
			return cluster, clusterWaitReady, nil
		}

		return cluster, clusterWaitPending, nil
	}
}

func resourceClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error updating cluster"
	client, clientCtx, diags := getServiceClient(ctx, m, qcCluster.NewClusterServiceClient)
	if diags.HasError() {
		return diags
	}
	// Expand the cluster
	cluster, jwtRbac, err := expandCluster(d, getDefaultAccountID(m))
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
	// Check if we need to enable JWT RBAC
	if jwtRbac != nil && *jwtRbac && !resp.GetCluster().GetState().GetJwtRbac() {
		_, err := client.EnableClusterJwtRbac(clientCtx, &qcCluster.EnableClusterJwtRbacRequest{
			AccountId: cluster.GetAccountId(),
			ClusterId: cluster.GetId(),
		}, grpc.Trailer(&trailer))
		// enrich prefix with request ID
		errorPrefix += getRequestID(trailer)
		if err != nil {
			return diag.FromErr(fmt.Errorf("%s (EnableJwtRbac): %w", errorPrefix, err))
		}
		// Update the cluster, so it's stored correctly in the state (flatten cluster)
		resp.Cluster.State.JwtRbac = true
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
