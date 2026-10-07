package qdrant

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	spacebackupv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/serverless/space/backup/v1"
)

// dataSourceServerlessBackupSchedule retrieves a specific backup schedule of a serverless space.
func dataSourceServerlessBackupSchedule() *schema.Resource {
	return &schema.Resource{
		Description: "Serverless Backup Schedule Data Source",
		ReadContext: dataSourceServerlessBackupScheduleRead,
		Schema:      serverlessBackupScheduleSchema(true),
	}
}

func dataSourceServerlessBackupScheduleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
		BackupScheduleId: d.Get(slBackupScheduleIDFieldName).(string),
	}, grpc.Trailer(&trailer))
	if err != nil {
		return grpcErrorDiag(errorPrefix, getRequestID(trailer), err)
	}
	if err := setFlattened(d, flattenServerlessBackupSchedule(resp.GetBackupSchedule())); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	d.SetId(resp.GetBackupSchedule().GetId())
	return nil
}
