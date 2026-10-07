package qdrant

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"google.golang.org/protobuf/types/known/timestamppb"

	spacebackupv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/serverless/space/backup/v1"
)

const (
	slBackupScheduleFieldTemplate = "Serverless Backup Schedule Schema %s field"

	slBackupScheduleIDFieldName              = "id"
	slBackupScheduleAccountIDFieldName       = "account_id"
	slBackupScheduleSpaceIDFieldName         = "space_id"
	slBackupScheduleNameFieldName            = "name"
	slBackupScheduleScheduleFieldName        = "schedule"
	slBackupScheduleRetentionPeriodFieldName = "retention_period"
	slBackupScheduleCollectionNameFieldName  = "collection_name"
	slBackupSchedulePausedFieldName          = "paused"
	slBackupSchedulePausedAtFieldName        = "paused_at"
	slBackupScheduleDeleteBackupsFieldName   = "delete_backups"
	slBackupScheduleCreatedAtFieldName       = "created_at"
	slBackupScheduleStatusFieldName          = "status"
	slBackupScheduleLastFiredAtFieldName     = "last_fired_at"
	slBackupSchedulesFieldName               = "schedules"
)

var (
	// slBackupScheduleNameRegexp matches the allowed names for backup schedules.
	slBackupScheduleNameRegexp = regexp.MustCompile(`^[\w\s-]+$`)
	// slCollectionNameRegexp matches a valid Qdrant collection name.
	slCollectionNameRegexp = regexp.MustCompile(`^[a-zA-Z0-9-_.]+$`)
)

// serverlessCollectionNameValidation validates a Qdrant collection name.
func serverlessCollectionNameValidation() schema.SchemaValidateDiagFunc {
	return validation.ToDiagFunc(validation.All(
		validation.StringLenBetween(1, 255),
		validation.StringMatch(slCollectionNameRegexp, "must only contain letters, numbers, '-', '_' and '.'"),
	))
}

// serverlessRetentionPeriodValidation validates a retention period (between 1 day and 365 days).
func serverlessRetentionPeriodValidation(v interface{}, _ cty.Path) diag.Diagnostics {
	d := parseDuration(v.(string))
	if d == nil {
		return diag.Errorf("invalid duration %q, expected a Go duration string (e.g. \"72h\")", v)
	}
	if secs := d.GetSeconds(); secs < 86400 || secs > 31536000 {
		return diag.Errorf("retention period must be between 24h and 8760h (365 days), got %q", v)
	}
	return nil
}

// serverlessBackupScheduleSchema returns the schema of a serverless backup schedule.
func serverlessBackupScheduleSchema(asDataSource bool) map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		slBackupScheduleIDFieldName: {
			Description: fmt.Sprintf(slBackupScheduleFieldTemplate, "Backup Schedule ID"),
			Type:        schema.TypeString,
			Required:    asDataSource,
			Computed:    !asDataSource,
		},
		slBackupScheduleAccountIDFieldName: {
			Description: fmt.Sprintf(slBackupScheduleFieldTemplate, "Account ID (defaults to the provider account ID)"),
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			ForceNew:    !asDataSource,
		},
		slBackupScheduleSpaceIDFieldName: {
			Description: fmt.Sprintf(slBackupScheduleFieldTemplate, "ID of the serverless space to back up"),
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    !asDataSource,
		},
		slBackupScheduleNameFieldName: {
			Description: fmt.Sprintf(slBackupScheduleFieldTemplate, "Name of the backup schedule (4-128 characters: letters, numbers, spaces, `_` and `-`)"),
			Type:        schema.TypeString,
			Required:    !asDataSource,
			Computed:    asDataSource,
		},
		slBackupScheduleScheduleFieldName: {
			Description: fmt.Sprintf(slBackupScheduleFieldTemplate, "Schedule in crontab format (e.g. `0 2 * * *`)"),
			Type:        schema.TypeString,
			Required:    !asDataSource,
			Computed:    asDataSource,
		},
		slBackupScheduleRetentionPeriodFieldName: {
			Description: fmt.Sprintf(slBackupScheduleFieldTemplate, "Retention period of the created backups, as a Go duration string between `24h` and `8760h`. If unset, backups are retained indefinitely"),
			Type:        schema.TypeString,
			Optional:    !asDataSource,
			Computed:    true,
		},
		slBackupScheduleCollectionNameFieldName: {
			Description: fmt.Sprintf(slBackupScheduleFieldTemplate, "Name of the collection to back up. If unset, the entire space is backed up"),
			Type:        schema.TypeString,
			Optional:    !asDataSource,
			Computed:    true,
		},
		slBackupSchedulePausedFieldName: {
			Description: fmt.Sprintf(slBackupScheduleFieldTemplate, "Whether the schedule is paused (no new backups are created while paused)"),
			Type:        schema.TypeBool,
			Optional:    !asDataSource,
			Computed:    asDataSource,
		},
		slBackupSchedulePausedAtFieldName: {
			Description: fmt.Sprintf(slBackupScheduleFieldTemplate, "Time at which the schedule was paused"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		slBackupScheduleCreatedAtFieldName: {
			Description: fmt.Sprintf(slBackupScheduleFieldTemplate, "Creation time"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		slBackupScheduleStatusFieldName: {
			Description: fmt.Sprintf(slBackupScheduleFieldTemplate, "Current status of the backup schedule"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		slBackupScheduleLastFiredAtFieldName: {
			Description: fmt.Sprintf(slBackupScheduleFieldTemplate, "Time when this schedule last created a backup"),
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
	if !asDataSource {
		s[slBackupSchedulePausedFieldName].Default = false
		s[slBackupScheduleNameFieldName].ValidateDiagFunc = validation.ToDiagFunc(validation.All(
			validation.StringLenBetween(4, 128),
			validation.StringMatch(slBackupScheduleNameRegexp, "must only contain letters, numbers, spaces, '_' and '-'"),
		))
		s[slBackupScheduleRetentionPeriodFieldName].ValidateDiagFunc = serverlessRetentionPeriodValidation
		s[slBackupScheduleRetentionPeriodFieldName].DiffSuppressFunc = suppressDurationDiff
		s[slBackupScheduleCollectionNameFieldName].ValidateDiagFunc = serverlessCollectionNameValidation()
		s[slBackupScheduleDeleteBackupsFieldName] = &schema.Schema{
			Description: "Whether the backups created by this schedule should be deleted when the schedule is destroyed. Defaults to false.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		}
	}
	return s
}

func serverlessBackupSchedulesDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		slBackupScheduleAccountIDFieldName: {
			Description: fmt.Sprintf(slBackupScheduleFieldTemplate, "Account ID (defaults to the provider account ID)"),
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
		},
		slBackupScheduleSpaceIDFieldName: {
			Description: fmt.Sprintf(slBackupScheduleFieldTemplate, "Only list the backup schedules of this serverless space"),
			Type:        schema.TypeString,
			Optional:    true,
		},
		slBackupSchedulesFieldName: {
			Description: "List of serverless backup schedules",
			Type:        schema.TypeList,
			Computed:    true,
			Elem:        &schema.Resource{Schema: serverlessBackupScheduleSchema(true)},
		},
	}
}

func flattenServerlessBackupSchedule(s *spacebackupv1.BackupSchedule) map[string]interface{} {
	return map[string]interface{}{
		slBackupScheduleIDFieldName:              s.GetId(),
		slBackupScheduleAccountIDFieldName:       s.GetAccountId(),
		slBackupScheduleSpaceIDFieldName:         s.GetSpaceId(),
		slBackupScheduleNameFieldName:            s.GetName(),
		slBackupScheduleScheduleFieldName:        s.GetSchedule(),
		slBackupScheduleRetentionPeriodFieldName: formatDuration(s.GetRetentionPeriod()),
		slBackupScheduleCollectionNameFieldName:  s.GetCollectionName(),
		slBackupSchedulePausedFieldName:          s.PausedAt != nil,
		slBackupSchedulePausedAtFieldName:        formatTime(s.GetPausedAt()),
		slBackupScheduleCreatedAtFieldName:       formatTime(s.GetCreatedAt()),
		slBackupScheduleStatusFieldName:          flattenProtoEnum(s.GetStatus()),
		slBackupScheduleLastFiredAtFieldName:     formatTime(s.GetLastFiredAt()),
	}
}

func flattenServerlessBackupSchedules(schedules []*spacebackupv1.BackupSchedule) []interface{} {
	result := make([]interface{}, 0, len(schedules))
	for _, s := range schedules {
		result = append(result, flattenServerlessBackupSchedule(s))
	}
	return result
}

// expandServerlessBackupSchedule builds the backup schedule from the resource data (without ID).
func expandServerlessBackupSchedule(d *schema.ResourceData, accountID string) *spacebackupv1.BackupSchedule {
	s := &spacebackupv1.BackupSchedule{
		AccountId: accountID,
		SpaceId:   d.Get(slBackupScheduleSpaceIDFieldName).(string),
		Name:      d.Get(slBackupScheduleNameFieldName).(string),
		Schedule:  d.Get(slBackupScheduleScheduleFieldName).(string),
	}
	if v, ok := d.GetOk(slBackupScheduleRetentionPeriodFieldName); ok {
		s.RetentionPeriod = parseDuration(v.(string))
	}
	if v, ok := d.GetOk(slBackupScheduleCollectionNameFieldName); ok {
		s.CollectionName = newPointer(v.(string))
	}
	if d.Get(slBackupSchedulePausedFieldName).(bool) {
		s.PausedAt = timestamppb.Now()
	}
	return s
}

// serverlessBackupScheduleUpdateMaskPaths maps Terraform attributes to the BackupSchedule proto field paths.
var serverlessBackupScheduleUpdateMaskPaths = []struct {
	attribute string
	path      string
}{
	{slBackupScheduleNameFieldName, "name"},
	{slBackupScheduleScheduleFieldName, "schedule"},
	{slBackupScheduleRetentionPeriodFieldName, "retention_period"},
	{slBackupScheduleCollectionNameFieldName, "collection_name"},
	{slBackupSchedulePausedFieldName, "paused_at"},
}

func serverlessBackupScheduleUpdateMask(hasChange func(string) bool) []string {
	var paths []string
	for _, p := range serverlessBackupScheduleUpdateMaskPaths {
		if hasChange(p.attribute) {
			paths = append(paths, p.path)
		}
	}
	return paths
}
