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

// dataSourceServerlessBackups lists the serverless backups of an account, with optional filters.
func dataSourceServerlessBackups() *schema.Resource {
	return &schema.Resource{
		Description: "Serverless Backups Data Source",
		ReadContext: dataSourceServerlessBackupsRead,
		Schema:      serverlessBackupsDataSourceSchema(),
	}
}

func dataSourceServerlessBackupsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error listing serverless backups"
	client, clientCtx, diags := getServiceClient(ctx, m, spacebackupv1.NewBackupServiceClient)
	if diags.HasError() {
		return diags
	}
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	optional := func(field string) *string {
		if v, ok := d.GetOk(field); ok {
			return newPointer(v.(string))
		}
		return nil
	}
	spaceID := optional(slBackupSpaceIDFieldName)
	scheduleID := optional(slBackupBackupScheduleIDFieldName)
	collectionName := optional(slBackupCollectionNameFieldName)

	var trailer metadata.MD
	backups, err := listAllPages(func(pageToken *string) ([]*spacebackupv1.Backup, string, error) {
		resp, err := client.ListBackups(clientCtx, &spacebackupv1.ListBackupsRequest{
			AccountId:        accountUUID.String(),
			SpaceId:          spaceID,
			BackupScheduleId: scheduleID,
			CollectionName:   collectionName,
			PageSize:         newPointer(serverlessListPageSize),
			PageToken:        pageToken,
		}, grpc.Trailer(&trailer))
		if err != nil {
			return nil, "", err
		}
		return resp.GetItems(), resp.GetNextPageToken(), nil
	})
	if err != nil {
		return grpcErrorDiag(errorPrefix, getRequestID(trailer), err)
	}
	if err := d.Set(slBackupsFieldName, flattenServerlessBackups(backups)); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	if err := d.Set(slBackupAccountIDFieldName, accountUUID.String()); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	d.SetId(fmt.Sprintf("%s/%s/%s/%s", accountUUID.String(),
		d.Get(slBackupSpaceIDFieldName).(string),
		d.Get(slBackupBackupScheduleIDFieldName).(string),
		d.Get(slBackupCollectionNameFieldName).(string)))
	return nil
}
