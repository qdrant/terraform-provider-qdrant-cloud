package qdrant

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	backupv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/backup/v1"
)

func resourceAccountsBackupSchedule() *schema.Resource {
	return &schema.Resource{
		Description:   "Backup Schedule Resource",
		CreateContext: resourceBackupScheduleCreate,
		ReadContext:   resourceBackupScheduleRead,
		UpdateContext: resourceBackupScheduleUpdate,
		DeleteContext: resourceBackupScheduleDelete,
		Schema:        accountsBackupScheduleResourceSchema(false),
		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts := strings.Split(d.Id(), "/")
				if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
					return nil, fmt.Errorf("unexpected format of ID (%s), expected <cluster_id>/<backup_schedule_id>", d.Id())
				}
				if err := d.Set(backupScheduleClusterIDFieldName, parts[0]); err != nil {
					return nil, fmt.Errorf("error setting cluster_id: %w", err)
				}
				d.SetId(parts[1])
				return []*schema.ResourceData{d}, nil
			},
		},
	}
}

func resourceBackupScheduleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error creating backup schedule"
	apiClientConn, clientCtx, diagnostics := getClientConnection(ctx, m)
	if diagnostics.HasError() {
		return diagnostics
	}
	client := backupv1.NewBackupServiceClient(apiClientConn)

	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	schedule := expandBackupSchedule(d)
	schedule.AccountId = accountUUID.String()

	var trailer metadata.MD
	resp, err := client.CreateBackupSchedule(clientCtx, &backupv1.CreateBackupScheduleRequest{
		BackupSchedule: schedule,
	}, grpc.Trailer(&trailer))
	errorPrefix += getRequestID(trailer)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	d.SetId(resp.GetBackupSchedule().GetId())
	return resourceBackupScheduleRead(ctx, d, m)
}

func resourceBackupScheduleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error reading backup schedule"
	apiClientConn, clientCtx, diagnostics := getClientConnection(ctx, m)
	if diagnostics.HasError() {
		return diagnostics
	}
	client := backupv1.NewBackupServiceClient(apiClientConn)

	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	var trailer metadata.MD
	resp, err := client.GetBackupSchedule(clientCtx, &backupv1.GetBackupScheduleRequest{
		AccountId:        accountUUID.String(),
		ClusterId:        d.Get(backupScheduleClusterIDFieldName).(string),
		BackupScheduleId: d.Id(),
	}, grpc.Trailer(&trailer))
	errorPrefix += getRequestID(trailer)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	flattened := flattenBackupSchedule(resp.GetBackupSchedule())
	for k, v := range flattened {
		if err := d.Set(k, v); err != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
		}
	}

	return nil
}

func resourceBackupScheduleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error updating backup schedule"
	apiClientConn, clientCtx, diagnostics := getClientConnection(ctx, m)
	if diagnostics.HasError() {
		return diagnostics
	}
	client := backupv1.NewBackupServiceClient(apiClientConn)

	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	if d.HasChanges(backupScheduleCronExpressionFieldName, backupScheduleRetentionPeriodFieldName) {
		schedule := expandBackupSchedule(d)
		schedule.AccountId = accountUUID.String()

		var trailer metadata.MD
		_, err := client.UpdateBackupSchedule(clientCtx, &backupv1.UpdateBackupScheduleRequest{
			BackupSchedule: schedule,
		}, grpc.Trailer(&trailer))
		errorPrefix += getRequestID(trailer)
		if err != nil {
			return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
		}
	}

	return resourceBackupScheduleRead(ctx, d, m)
}

func resourceBackupScheduleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error deleting backup schedule"
	apiClientConn, clientCtx, diagnostics := getClientConnection(ctx, m)
	if diagnostics.HasError() {
		return diagnostics
	}
	client := backupv1.NewBackupServiceClient(apiClientConn)

	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Do we need to delete the backups, created with this schedule?
	deleteBackups := d.Get(backupScheduleDeleteBackupsOnDestroy).(bool)
	// Delete the cackup schedule
	var trailer metadata.MD
	_, err = client.DeleteBackupSchedule(clientCtx, &backupv1.DeleteBackupScheduleRequest{
		AccountId:        accountUUID.String(),
		BackupScheduleId: d.Id(),
		DeleteBackups:    newPointer(deleteBackups),
	}, grpc.Trailer(&trailer))
	errorPrefix += getRequestID(trailer)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	d.SetId("")
	return nil
}
