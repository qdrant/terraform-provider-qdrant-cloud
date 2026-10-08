package qdrant

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	spacebackupv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/serverless/space/backup/v1"
)

// resourceServerlessBackup constructs a Terraform resource for managing a one-off backup of a serverless space.
// Backups are immutable, so all configurable fields force a new backup.
func resourceServerlessBackup() *schema.Resource {
	return &schema.Resource{
		Description:   "Serverless Backup Resource (one-off backup of a serverless space)",
		ReadContext:   resourceServerlessBackupRead,
		CreateContext: resourceServerlessBackupCreate,
		UpdateContext: nil, // Backups are immutable
		DeleteContext: resourceServerlessBackupDelete,
		Schema:        serverlessBackupSchema(false),
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(serverlessDefaultTimeout),
		},
	}
}

func resourceServerlessBackupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error reading serverless backup"
	client, clientCtx, diags := getServiceClient(ctx, m, spacebackupv1.NewBackupServiceClient)
	if diags.HasError() {
		return diags
	}
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	var trailer metadata.MD
	resp, err := client.GetBackup(clientCtx, &spacebackupv1.GetBackupRequest{
		AccountId: accountUUID.String(),
		BackupId:  d.Id(),
	}, grpc.Trailer(&trailer))
	if err != nil {
		if isNotFound(err) {
			d.SetId("")
			return nil
		}
		return grpcErrorDiag(errorPrefix, getRequestID(trailer), err)
	}
	backup := resp.GetBackup()
	if backup.GetDeletedAt() != nil {
		d.SetId("")
		return nil
	}
	if err := setFlattened(d, flattenServerlessBackup(backup)); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	return nil
}

func resourceServerlessBackupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error creating serverless backup"
	client, clientCtx, diags := getServiceClient(ctx, m, spacebackupv1.NewBackupServiceClient)
	if diags.HasError() {
		return diags
	}
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	var trailer metadata.MD
	resp, err := client.CreateBackup(clientCtx, &spacebackupv1.CreateBackupRequest{
		Backup: expandServerlessBackup(d, accountUUID.String()),
	}, grpc.Trailer(&trailer))
	if err != nil {
		return grpcErrorDiag(errorPrefix, getRequestID(trailer), err)
	}
	d.SetId(resp.GetBackup().GetId())

	stateConf := &retry.StateChangeConf{
		Pending:      []string{serverlessWaitPending},
		Target:       []string{serverlessWaitReady},
		Refresh:      serverlessBackupRefreshFunc(client, clientCtx, accountUUID.String(), d.Id()),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		PollInterval: serverlessPollInterval,
	}
	result, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	if err := setFlattened(d, flattenServerlessBackup(result.(*spacebackupv1.Backup))); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	return nil
}

func resourceServerlessBackupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error deleting serverless backup"
	client, clientCtx, diags := getServiceClient(ctx, m, spacebackupv1.NewBackupServiceClient)
	if diags.HasError() {
		return diags
	}
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	var trailer metadata.MD
	_, err = client.DeleteBackup(clientCtx, &spacebackupv1.DeleteBackupRequest{
		AccountId: accountUUID.String(),
		BackupId:  d.Id(),
	}, grpc.Trailer(&trailer))
	if err != nil && !isNotFound(err) {
		return grpcErrorDiag(errorPrefix, getRequestID(trailer), err)
	}
	d.SetId("")
	return nil
}

// serverlessBackupRefreshFunc returns a StateRefreshFunc that polls GetBackup until the backup succeeded or failed.
func serverlessBackupRefreshFunc(
	client spacebackupv1.BackupServiceClient,
	ctx context.Context,
	accountID, backupID string,
) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.GetBackup(ctx, &spacebackupv1.GetBackupRequest{
			AccountId: accountID,
			BackupId:  backupID,
		})
		if err != nil {
			return nil, "", err
		}
		backup := resp.GetBackup()
		switch backup.GetStatus() {
		case spacebackupv1.BackupStatus_BACKUP_STATUS_SUCCEEDED:
			return backup, serverlessWaitReady, nil
		case spacebackupv1.BackupStatus_BACKUP_STATUS_FAILED,
			spacebackupv1.BackupStatus_BACKUP_STATUS_FAILED_TO_SYNC,
			spacebackupv1.BackupStatus_BACKUP_STATUS_SKIPPED,
			spacebackupv1.BackupStatus_BACKUP_STATUS_NOT_FOUND:
			return nil, "", fmt.Errorf("backup did not succeed (status=%q)", backup.GetStatus().String())
		default:
			return backup, serverlessWaitPending, nil
		}
	}
}
