package qdrant

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	backupv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/backup/v1"
)

func dataSourceAccountsBackupSchedules() *schema.Resource {
	return &schema.Resource{
		Description: "Backup Schedules Data Source",
		ReadContext: dataAccountsBackupSchedulesRead,
		Schema:      accountsBackupSchedulesDataSourceSchema(),
	}
}

func dataAccountsBackupSchedulesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error listing backup schedules"
	apiClientConn, clientCtx, diagnostics := getClientConnection(ctx, m)
	if diagnostics.HasError() {
		return diagnostics
	}
	client := backupv1.NewBackupServiceClient(apiClientConn)

	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	clusterID := d.Get(backupScheduleClusterIDFieldName).(string)

	var trailer metadata.MD
	resp, err := client.ListBackupSchedules(clientCtx, &backupv1.ListBackupSchedulesRequest{
		AccountId: accountUUID.String(),
		ClusterId: newPointer(clusterID),
	}, grpc.Trailer(&trailer))
	errorPrefix += getRequestID(trailer)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	if err := d.Set(backupSchedulesFieldName, flattenBackupSchedules(resp.GetItems())); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	d.SetId(fmt.Sprintf("%s/%s", accountUUID.String(), clusterID))
	return nil
}
