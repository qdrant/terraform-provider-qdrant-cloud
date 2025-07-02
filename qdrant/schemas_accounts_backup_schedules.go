package qdrant

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	backupv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/backup/v1"
)

const (
	backupScheduleFieldTemplate            = "Backup Schedule Schema %s field"
	backupScheduleIDFieldName              = "id"
	backupScheduleAccountIDFieldName       = "account_id"
	backupScheduleClusterIDFieldName       = "cluster_id"
	backupScheduleCronExpressionFieldName  = "cron_expression"
	backupScheduleRetentionPeriodFieldName = "retention_period"
	backupScheduleStatusFieldName          = "status"
	backupScheduleCreatedAtFieldName       = "created_at"
	backupScheduleDeletedAtFieldName       = "deleted_at"
	backupSchedulesFieldName               = "schedules"
)

func accountsBackupScheduleResourceSchema(asDataSource bool) map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		backupScheduleIDFieldName: {
			Description: fmt.Sprintf(backupScheduleFieldTemplate, "ID"),
			Type:        schema.TypeString,
			Required:    asDataSource,
			Computed:    !asDataSource,
		},
		backupScheduleAccountIDFieldName: {
			Description: fmt.Sprintf(backupScheduleFieldTemplate, "Account ID"),
			Type:        schema.TypeString,
			Optional:    !asDataSource,
			Computed:    true,
		},
		backupScheduleClusterIDFieldName: {
			Description: fmt.Sprintf(backupScheduleFieldTemplate, "Cluster ID"),
			Type:        schema.TypeString,
			Required:    !asDataSource,
			Computed:    asDataSource,
			ForceNew:    !asDataSource,
		},
		backupScheduleCronExpressionFieldName: {
			Description: fmt.Sprintf(backupScheduleFieldTemplate, "Cron expression for the schedule"),
			Type:        schema.TypeString,
			Required:    !asDataSource,
			Computed:    asDataSource,
		},
		backupScheduleRetentionPeriodFieldName: {
			Description: fmt.Sprintf(backupScheduleFieldTemplate, `Retention period as a Go duration string (e.g., "72h"). The "d" unit for days is not supported.`),
			Type:        schema.TypeString,
			Optional:    !asDataSource,
			Computed:    asDataSource,
		},
		backupScheduleCreatedAtFieldName: {
			Description: fmt.Sprintf(backupScheduleFieldTemplate, "Creation time"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		backupScheduleDeletedAtFieldName: {
			Description: fmt.Sprintf(backupScheduleFieldTemplate, "Deletion time"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		backupScheduleStatusFieldName: {
			Description: fmt.Sprintf(backupScheduleFieldTemplate, "Status"),
			Type:        schema.TypeString,
			Computed:    true,
		},
	}

	if !asDataSource {
		s[backupScheduleRetentionPeriodFieldName].DiffSuppressFunc = suppressDurationDiff
	}
	return s
}

func accountsBackupSchedulesDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		backupScheduleAccountIDFieldName: {
			Description: fmt.Sprintf(backupScheduleFieldTemplate, "Account ID"),
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
		},
		backupScheduleClusterIDFieldName: {
			Description: fmt.Sprintf(backupScheduleFieldTemplate, "Cluster ID"),
			Type:        schema.TypeString,
			Required:    true,
		},
		backupSchedulesFieldName: {
			Description: "List of backup schedules",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: accountsBackupScheduleResourceSchema(true),
			},
		},
	}
}

func flattenBackupSchedule(schedule *backupv1.BackupSchedule) map[string]interface{} {
	return map[string]interface{}{
		backupScheduleIDFieldName:              schedule.GetId(),
		backupScheduleAccountIDFieldName:       schedule.GetAccountId(),
		backupScheduleClusterIDFieldName:       schedule.GetClusterId(),
		backupScheduleCronExpressionFieldName:  schedule.GetSchedule(),
		backupScheduleRetentionPeriodFieldName: formatDuration(schedule.GetRetentionPeriod()),
		backupScheduleCreatedAtFieldName:       formatTime(schedule.GetCreatedAt()),
		backupScheduleDeletedAtFieldName:       formatTime(schedule.GetDeletedAt()),
		backupScheduleStatusFieldName:          schedule.GetStatus().String(),
	}
}

func flattenBackupSchedules(schedules []*backupv1.BackupSchedule) []interface{} {
	flattened := make([]interface{}, len(schedules))
	for i, schedule := range schedules {
		flattened[i] = flattenBackupSchedule(schedule)
	}
	return flattened
}

func expandBackupSchedule(d *schema.ResourceData) *backupv1.BackupSchedule {
	return &backupv1.BackupSchedule{
		Id:              d.Id(),
		ClusterId:       d.Get(backupScheduleClusterIDFieldName).(string),
		Schedule:        d.Get(backupScheduleCronExpressionFieldName).(string),
		RetentionPeriod: parseDuration(d.Get(backupScheduleRetentionPeriodFieldName).(string)),
	}
}
