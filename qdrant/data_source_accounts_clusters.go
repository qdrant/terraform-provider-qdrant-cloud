package qdrant

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	qcCluster "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/v1"
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
		Schema:      accountsClusterSchema(true),
	}
}

// dataSourceAccountsClustersRead performs a read operation to fetch all clusters associated with a specific account.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func dataSourceAccountsClustersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error listing clusters"
	// Get a client connection and context
	apiClientConn, clientCtx, diagnostics := getClientConnection(ctx, m)
	if diagnostics.HasError() {
		return diagnostics
	}
	// Get a client
	client := qcCluster.NewClusterServiceClient(apiClientConn)
	// Get The account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// List all clusters for the provided account
	var header metadata.MD
	resp, err := client.ListClusters(clientCtx, &qcCluster.ListClustersRequest{
		AccountId: accountUUID.String(),
	}, grpc.Header(&header))
	// enrich prefix with request ID
	errorPrefix += getRequestID(header)
	if err != nil {
		d := diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
		if d.HasError() {
			return d
		}
	}
	// Update the Terraform state
	if err := d.Set("clusters", flattenClusters(resp.GetItems())); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	d.SetId(time.Now().UTC().Format(time.RFC3339))
	return nil
}

// dataSourceAccountsClusterRead performs a read operation to fetch a specific cluster associated with a given account and cluster ID.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func dataSourceAccountsClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error getting cluster"
	// Get a client connection and context
	apiClientConn, clientCtx, diagnostics := getClientConnection(ctx, m)
	if diagnostics.HasError() {
		return diagnostics
	}
	// Get a client
	client := qcCluster.NewClusterServiceClient(apiClientConn)
	// Get The account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Get the cluster ID
	clusterID := d.Get("id").(string)
	// Fetch the cluster
	var header metadata.MD
	resp, err := client.GetCluster(clientCtx, &qcCluster.GetClusterRequest{
		AccountId: accountUUID.String(),
		ClusterId: clusterID,
	}, grpc.Header(&header))
	// enrich prefix with request ID
	errorPrefix += getRequestID(header)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Flatten cluster and store in Terraform state
	for k, v := range flattenCluster(resp.GetCluster()) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
		}
	}
	d.SetId(clusterID)
	return nil
}
