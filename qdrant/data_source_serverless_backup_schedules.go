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

// dataSourceServerlessBackupSchedules lists the serverless backup schedules of an account, optionally for a single space.
func dataSourceServerlessBackupSchedules() *schema.Resource {
	return &schema.Resource{
		Description: "Serverless Backup Schedules Data Source",
		ReadContext: dataSourceServerlessBackupSchedulesRead,
		Schema:      serverlessBackupSchedulesDataSourceSchema(),
	}
}

func dataSourceServerlessBackupSchedulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error listing serverless backup schedules"
	client, clientCtx, diags := getServiceClient(ctx, m, spacebackupv1.NewBackupServiceClient)
	if diags.HasError() {
		return diags
	}
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	var spaceID *string
	if v, ok := d.GetOk(slBackupScheduleSpaceIDFieldName); ok {
		spaceID = newPointer(v.(string))
	}
	var trailer metadata.MD
	schedules, err := listAllPages(func(pageToken *string) ([]*spacebackupv1.BackupSchedule, string, error) {
		resp, err := client.ListBackupSchedules(clientCtx, &spacebackupv1.ListBackupSchedulesRequest{
			AccountId: accountUUID.String(),
			SpaceId:   spaceID,
			PageSize:  newPointer(serverlessListPageSize),
			PageToken: pageToken,
		}, grpc.Trailer(&trailer))
		if err != nil {
			return nil, "", err
		}
		return resp.GetItems(), resp.GetNextPageToken(), nil
	})
	if err != nil {
		return grpcErrorDiag(errorPrefix, getRequestID(trailer), err)
	}
	if err := d.Set(slBackupSchedulesFieldName, flattenServerlessBackupSchedules(schedules)); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	if err := d.Set(slBackupScheduleAccountIDFieldName, accountUUID.String()); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	d.SetId(fmt.Sprintf("%s/%s", accountUUID.String(), d.Get(slBackupScheduleSpaceIDFieldName).(string)))
	return nil
}
