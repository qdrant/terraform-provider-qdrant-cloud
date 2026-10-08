package qdrant

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	spacebackupv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/serverless/space/backup/v1"
)

const (
	slBackupFieldTemplate = "Serverless Backup Schema %s field"

	slBackupIDFieldName               = "id"
	slBackupAccountIDFieldName        = "account_id"
	slBackupSpaceIDFieldName          = "space_id"
	slBackupNameFieldName             = "name"
	slBackupCollectionNameFieldName   = "collection_name"
	slBackupRetentionPeriodFieldName  = "retention_period"
	slBackupBackupScheduleIDFieldName = "backup_schedule_id"
	slBackupCreatedAtFieldName        = "created_at"
	slBackupStatusFieldName           = "status"
	slBackupStatsFieldName            = "stats"
	slBackupsFieldName                = "backups"

	slBackupStatsCollectionCountFieldName = "collection_count"
	slBackupStatsSizeBytesFieldName       = "size_bytes"
	slBackupStatsTotalPointsFieldName     = "total_points"
	slBackupStatsDurationFieldName        = "duration"
)

// serverlessBackupSchema returns the schema of a serverless backup.
// Backups are immutable, so all configurable fields force a new backup.
func serverlessBackupSchema(asDataSource bool) map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		slBackupIDFieldName: {
			Description: fmt.Sprintf(slBackupFieldTemplate, "Backup ID, can be used as `from_backup_id` of a serverless space"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		slBackupAccountIDFieldName: {
			Description: fmt.Sprintf(slBackupFieldTemplate, "Account ID (defaults to the provider account ID)"),
			Type:        schema.TypeString,
			Optional:    !asDataSource,
			Computed:    true,
			ForceNew:    !asDataSource,
		},
		slBackupSpaceIDFieldName: {
			Description: fmt.Sprintf(slBackupFieldTemplate, "ID of the serverless space to back up"),
			Type:        schema.TypeString,
			Required:    !asDataSource,
			Computed:    asDataSource,
			ForceNew:    !asDataSource,
		},
		slBackupCollectionNameFieldName: {
			Description: fmt.Sprintf(slBackupFieldTemplate, "Name of the collection to back up. If unset, the entire space is backed up"),
			Type:        schema.TypeString,
			Optional:    !asDataSource,
			Computed:    true,
			ForceNew:    !asDataSource,
		},
		slBackupRetentionPeriodFieldName: {
			Description: fmt.Sprintf(slBackupFieldTemplate, "Retention period of the backup, as a Go duration string between `24h` and `8760h`. If unset, the backup is retained indefinitely"),
			Type:        schema.TypeString,
			Optional:    !asDataSource,
			Computed:    true,
			ForceNew:    !asDataSource,
		},
		slBackupNameFieldName: {
			Description: fmt.Sprintf(slBackupFieldTemplate, "Name of the backup (generated)"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		slBackupBackupScheduleIDFieldName: {
			Description: fmt.Sprintf(slBackupFieldTemplate, "ID of the backup schedule that created the backup (empty for manual backups)"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		slBackupCreatedAtFieldName: {
			Description: fmt.Sprintf(slBackupFieldTemplate, "Creation time"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		slBackupStatusFieldName: {
			Description: fmt.Sprintf(slBackupFieldTemplate, "Current status of the backup"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		slBackupStatsFieldName: {
			Description: fmt.Sprintf(slBackupFieldTemplate, "Statistics of the backup"),
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				slBackupStatsCollectionCountFieldName: {
					Description: "Number of collections in the backup.",
					Type:        schema.TypeInt,
					Computed:    true,
				},
				slBackupStatsSizeBytesFieldName: {
					Description: "Size of the backup in bytes.",
					Type:        schema.TypeInt,
					Computed:    true,
				},
				slBackupStatsTotalPointsFieldName: {
					Description: "Total number of points in the backup.",
					Type:        schema.TypeInt,
					Computed:    true,
				},
				slBackupStatsDurationFieldName: {
					Description: "Time it took to create the backup.",
					Type:        schema.TypeString,
					Computed:    true,
				},
			}},
		},
	}
	if !asDataSource {
		s[slBackupCollectionNameFieldName].ValidateDiagFunc = serverlessCollectionNameValidation()
		s[slBackupRetentionPeriodFieldName].ValidateDiagFunc = serverlessRetentionPeriodValidation
		s[slBackupRetentionPeriodFieldName].DiffSuppressFunc = suppressDurationDiff
	}
	return s
}

func serverlessBackupsDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		slBackupAccountIDFieldName: {
			Description: fmt.Sprintf(slBackupFieldTemplate, "Account ID (defaults to the provider account ID)"),
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
		},
		slBackupSpaceIDFieldName: {
			Description: fmt.Sprintf(slBackupFieldTemplate, "Only list the backups of this serverless space"),
			Type:        schema.TypeString,
			Optional:    true,
		},
		slBackupBackupScheduleIDFieldName: {
			Description: fmt.Sprintf(slBackupFieldTemplate, "Only list the backups created by this backup schedule"),
			Type:        schema.TypeString,
			Optional:    true,
		},
		slBackupCollectionNameFieldName: {
			Description: fmt.Sprintf(slBackupFieldTemplate, "Only list the backups of this collection"),
			Type:        schema.TypeString,
			Optional:    true,
		},
		slBackupsFieldName: {
			Description: "List of serverless backups",
			Type:        schema.TypeList,
			Computed:    true,
			Elem:        &schema.Resource{Schema: serverlessBackupSchema(true)},
		},
	}
}

func flattenServerlessBackup(b *spacebackupv1.Backup) map[string]interface{} {
	stats := []interface{}{}
	if st := b.GetStats(); st != nil {
		stats = append(stats, map[string]interface{}{
			slBackupStatsCollectionCountFieldName: int(st.GetCollectionCount()),
			slBackupStatsSizeBytesFieldName:       int(st.GetSizeBytes()),
			slBackupStatsTotalPointsFieldName:     int(st.GetTotalPoints()),
			slBackupStatsDurationFieldName:        formatDuration(st.GetDuration()),
		})
	}
	return map[string]interface{}{
		slBackupIDFieldName:               b.GetId(),
		slBackupAccountIDFieldName:        b.GetAccountId(),
		slBackupSpaceIDFieldName:          b.GetSpaceId(),
		slBackupNameFieldName:             b.GetName(),
		slBackupCollectionNameFieldName:   b.GetCollectionName(),
		slBackupRetentionPeriodFieldName:  formatDuration(b.GetRetentionPeriod()),
		slBackupBackupScheduleIDFieldName: b.GetBackupScheduleId(),
		slBackupCreatedAtFieldName:        formatTime(b.GetCreatedAt()),
		slBackupStatusFieldName:           flattenProtoEnum(b.GetStatus()),
		slBackupStatsFieldName:            stats,
	}
}

func flattenServerlessBackups(backups []*spacebackupv1.Backup) []interface{} {
	result := make([]interface{}, 0, len(backups))
	for _, b := range backups {
		result = append(result, flattenServerlessBackup(b))
	}
	return result
}

// expandServerlessBackup builds the backup to create from the resource data.
func expandServerlessBackup(d *schema.ResourceData, accountID string) *spacebackupv1.Backup {
	b := &spacebackupv1.Backup{
		AccountId: accountID,
		SpaceId:   d.Get(slBackupSpaceIDFieldName).(string),
	}
	if v, ok := d.GetOk(slBackupCollectionNameFieldName); ok {
		b.CollectionName = newPointer(v.(string))
	}
	if v, ok := d.GetOk(slBackupRetentionPeriodFieldName); ok {
		b.RetentionPeriod = parseDuration(v.(string))
	}
	return b
}
