package qdrant

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/fieldmaskpb"

	spacebackupv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/serverless/space/backup/v1"
)

// resourceServerlessBackupSchedule constructs a Terraform resource for managing backup schedules of a serverless space.
func resourceServerlessBackupSchedule() *schema.Resource {
	return &schema.Resource{
		Description:   "Serverless Backup Schedule Resource",
		ReadContext:   resourceServerlessBackupScheduleRead,
		CreateContext: resourceServerlessBackupScheduleCreate,
		UpdateContext: resourceServerlessBackupScheduleUpdate,
		DeleteContext: resourceServerlessBackupScheduleDelete,
		Schema:        serverlessBackupScheduleSchema(false),
		Importer: &schema.ResourceImporter{
			StateContext: importSpaceChildState(slBackupScheduleSpaceIDFieldName, "backup_schedule_id"),
		},
	}
}

func resourceServerlessBackupScheduleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error reading serverless backup schedule"
	client, clientCtx, diags := getServiceClient(ctx, m, spacebackupv1.NewBackupServiceClient)
	if diags.HasError() {
		return diags
	}
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	var trailer metadata.MD
	resp, err := client.GetBackupSchedule(clientCtx, &spacebackupv1.GetBackupScheduleRequest{
		AccountId:        accountUUID.String(),
		SpaceId:          d.Get(slBackupScheduleSpaceIDFieldName).(string),
		BackupScheduleId: d.Id(),
	}, grpc.Trailer(&trailer))
	if err != nil {
		if isNotFound(err) {
			d.SetId("")
			return nil
		}
		return grpcErrorDiag(errorPrefix, getRequestID(trailer), err)
	}
	schedule := resp.GetBackupSchedule()
	if schedule.GetDeletedAt() != nil {
		d.SetId("")
		return nil
	}
	if err := setFlattened(d, flattenServerlessBackupSchedule(schedule)); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	return nil
}

func resourceServerlessBackupScheduleCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error creating serverless backup schedule"
	client, clientCtx, diags := getServiceClient(ctx, m, spacebackupv1.NewBackupServiceClient)
	if diags.HasError() {
		return diags
	}
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	var trailer metadata.MD
	resp, err := client.CreateBackupSchedule(clientCtx, &spacebackupv1.CreateBackupScheduleRequest{
		BackupSchedule: expandServerlessBackupSchedule(d, accountUUID.String()),
	}, grpc.Trailer(&trailer))
	if err != nil {
		return grpcErrorDiag(errorPrefix, getRequestID(trailer), err)
	}
	d.SetId(resp.GetBackupSchedule().GetId())
	if err := setFlattened(d, flattenServerlessBackupSchedule(resp.GetBackupSchedule())); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	return nil
}

func resourceServerlessBackupScheduleUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error updating serverless backup schedule"
	client, clientCtx, diags := getServiceClient(ctx, m, spacebackupv1.NewBackupServiceClient)
	if diags.HasError() {
		return diags
	}
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	paths := serverlessBackupScheduleUpdateMask(d.HasChange)
	if len(paths) == 0 {
		return resourceServerlessBackupScheduleRead(ctx, d, m)
	}
	schedule := expandServerlessBackupSchedule(d, accountUUID.String())
	schedule.Id = d.Id()
	var trailer metadata.MD
	resp, err := client.UpdateBackupSchedule(clientCtx, &spacebackupv1.UpdateBackupScheduleRequest{
		BackupSchedule: schedule,
		UpdateMask:     &fieldmaskpb.FieldMask{Paths: paths},
	}, grpc.Trailer(&trailer))
	if err != nil {
		return grpcErrorDiag(errorPrefix, getRequestID(trailer), err)
	}
	if err := setFlattened(d, flattenServerlessBackupSchedule(resp.GetBackupSchedule())); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	return nil
}

func resourceServerlessBackupScheduleDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error deleting serverless backup schedule"
	client, clientCtx, diags := getServiceClient(ctx, m, spacebackupv1.NewBackupServiceClient)
	if diags.HasError() {
		return diags
	}
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	var trailer metadata.MD
	_, err = client.DeleteBackupSchedule(clientCtx, &spacebackupv1.DeleteBackupScheduleRequest{
		AccountId:        accountUUID.String(),
		BackupScheduleId: d.Id(),
		DeleteBackups:    newPointer(d.Get(slBackupScheduleDeleteBackupsFieldName).(bool)),
	}, grpc.Trailer(&trailer))
	if err != nil && !isNotFound(err) {
		return grpcErrorDiag(errorPrefix, getRequestID(trailer), err)
	}
	d.SetId("")
	return nil
}
