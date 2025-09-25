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

	qcb "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/backup/v1"
)

// resourceAccountsManualBackup constructs a Terraform resource for managing a one-off
// cluster backup associated with an account. Returns a schema.Resource configured with
// schema definitions and CRUD functions.
func resourceAccountsManualBackup() *schema.Resource {
	return &schema.Resource{
		Description:   "Cluster Backup Resource (one-off backup).",
		CreateContext: resourceBackupCreate,
		ReadContext:   resourceBackupRead,
		UpdateContext: resourceBackupUpdate, // no-op (backups are immutable)
		DeleteContext: resourceBackupDelete,
		Schema:        accountsBackupSchema(),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

// resourceBackupCreate performs a create operation to trigger a new manual backup.
// ctx: Context to carry deadlines/cancellation across API calls.
// d: Resource data used to build the request and persist state.
// m: Provider meta containing client config/defaults.
// Returns diagnostics describing any runtime issues.
func resourceBackupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error creating backup"

	// Client
	apiClientConn, clientCtx, diags := getClientConnection(ctx, m)
	if diags.HasError() {
		return diags
	}
	client := qcb.NewBackupServiceClient(apiClientConn)

	// Build request
	backup, err := expandBackup(d, getDefaultAccountID(m))
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	var trailer metadata.MD
	resp, err := client.CreateBackup(
		clientCtx,
		&qcb.CreateBackupRequest{Backup: backup},
		grpc.Trailer(&trailer),
	)
	errorPrefix += getRequestID(trailer)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	created := resp.GetBackup()
	d.SetId(created.GetId())

	for k, v := range flattenBackup(created) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
		}
	}

	return nil
}

// resourceBackupRead performs a read operation to fetch the latest state of a manual backup.
// ctx: Context to carry deadlines/cancellation across API calls.
// d: Resource data providing the backup ID to read and where to persist state.
// m: Provider meta containing client config/defaults.
// Returns diagnostics describing any runtime issues (clears ID if NotFound).
func resourceBackupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error reading backup"

	// Client
	apiClientConn, clientCtx, diags := getClientConnection(ctx, m)
	if diags.HasError() {
		return diags
	}
	client := qcb.NewBackupServiceClient(apiClientConn)

	// Account from state or provider default
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	var trailer metadata.MD
	resp, err := client.GetBackup(
		clientCtx,
		&qcb.GetBackupRequest{
			AccountId: accountUUID.String(),
			BackupId:  d.Id(),
		},
		grpc.Trailer(&trailer),
	)
	errorPrefix += getRequestID(trailer)

	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			// Deleted outside Terraform: clear state.
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	got := resp.GetBackup()
	d.SetId(got.GetId())

	for k, v := range flattenBackup(got) {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
		}
	}

	return nil
}

// resourceBackupUpdate performs a no-op update since backups are immutable.
// ctx: Context to carry deadlines/cancellation across API calls.
// d: Resource data for the backup.
// m: Provider meta containing client config/defaults.
// Returns diagnostics from a state refresh (Read).
func resourceBackupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceBackupRead(ctx, d, m)
}

// resourceBackupDelete performs a delete operation to remove a manual backup.
// ctx: Context to carry deadlines/cancellation across API calls.
// d: Resource data providing the backup ID to delete.
// m: Provider meta containing client config/defaults.
// Returns diagnostics describing any runtime issues (treats NotFound as success).
func resourceBackupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error deleting backup"

	// Client
	apiClientConn, clientCtx, diags := getClientConnection(ctx, m)
	if diags.HasError() {
		return diags
	}
	client := qcb.NewBackupServiceClient(apiClientConn)

	// Account from state or provider default
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	var trailer metadata.MD
	_, err = client.DeleteBackup(
		clientCtx,
		&qcb.DeleteBackupRequest{
			AccountId: accountUUID.String(),
			BackupId:  d.Id(),
		},
		grpc.Trailer(&trailer),
	)
	errorPrefix += getRequestID(trailer)

	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	d.SetId("")
	return nil
}
