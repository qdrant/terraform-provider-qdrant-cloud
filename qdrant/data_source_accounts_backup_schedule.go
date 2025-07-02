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

func dataSourceAccountsBackupSchedule() *schema.Resource {
	return &schema.Resource{
		Description: "Backup Schedule Data Source",
		ReadContext: dataAccountsBackupScheduleRead,
		Schema:      accountsBackupScheduleResourceSchema(true),
	}
}

func dataAccountsBackupScheduleRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
