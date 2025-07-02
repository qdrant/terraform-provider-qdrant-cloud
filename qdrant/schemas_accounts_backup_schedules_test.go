package qdrant

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	backupv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/backup/v1"
)

func TestFlattenBackupSchedule(t *testing.T) {
	createdAt := timestamppb.New(time.Now())
	deletedAt := timestamppb.New(time.Now().Add(time.Hour))
	retentionPeriod := durationpb.New(7 * 24 * time.Hour)

	schedule := &backupv1.BackupSchedule{
		Id:              "schedule-id-1",
		AccountId:       "account-id-1",
		ClusterId:       "cluster-id-1",
		Schedule:        "0 0 * * *",
		RetentionPeriod: retentionPeriod,
		CreatedAt:       createdAt,
		DeletedAt:       deletedAt,
		Status:          backupv1.BackupScheduleStatus_BACKUP_SCHEDULE_STATUS_ACTIVE,
	}

	flattened := flattenBackupSchedule(schedule)

	expected := map[string]interface{}{
		backupScheduleIDFieldName:              "schedule-id-1",
		backupScheduleAccountIDFieldName:       "account-id-1",
		backupScheduleClusterIDFieldName:       "cluster-id-1",
		backupScheduleCronExpressionFieldName:  "0 0 * * *",
		backupScheduleRetentionPeriodFieldName: formatDuration(retentionPeriod),
		backupScheduleCreatedAtFieldName:       formatTime(createdAt),
		backupScheduleDeletedAtFieldName:       formatTime(deletedAt),
		backupScheduleStatusFieldName:          "BACKUP_SCHEDULE_STATUS_ACTIVE",
	}

	assert.Equal(t, expected, flattened)
}

func TestExpandBackupSchedule(t *testing.T) {
	retentionPeriod := 7 * 24 * time.Hour
	expected := &backupv1.BackupSchedule{
		Id: "schedule-id-2",
		// AccountId is set in the resource function
		ClusterId:       "cluster-id-2",
		Schedule:        "0 12 * * *",
		RetentionPeriod: durationpb.New(retentionPeriod),
	}

	d := schema.TestResourceDataRaw(t, accountsBackupScheduleResourceSchema(false), map[string]interface{}{
		// backupScheduleAccountIDFieldName is not set here, as it's handled by the resource function
		backupScheduleClusterIDFieldName:       expected.GetClusterId(),
		backupScheduleCronExpressionFieldName:  expected.GetSchedule(),
		backupScheduleRetentionPeriodFieldName: retentionPeriod.String(),
	})
	d.SetId(expected.GetId())

	result := expandBackupSchedule(d)
	assert.Equal(t, expected, result)
}

func TestFlattenBackupSchedules(t *testing.T) {
	createdAt := timestamppb.New(time.Now())
	retentionPeriod1 := durationpb.New(3 * 24 * time.Hour)
	retentionPeriod2 := durationpb.New(5 * 24 * time.Hour)

	schedules := []*backupv1.BackupSchedule{
		{
			Id:              "schedule-id-1",
			AccountId:       "account-id-1",
			ClusterId:       "cluster-id-1",
			Schedule:        "0 0 * * *",
			RetentionPeriod: retentionPeriod1,
			CreatedAt:       createdAt,
			Status:          backupv1.BackupScheduleStatus_BACKUP_SCHEDULE_STATUS_ACTIVE,
		},
		{
			Id:              "schedule-id-2",
			AccountId:       "account-id-1",
			ClusterId:       "cluster-id-1",
			Schedule:        "0 12 * * *",
			RetentionPeriod: retentionPeriod2,
			CreatedAt:       createdAt,
			Status:          backupv1.BackupScheduleStatus_BACKUP_SCHEDULE_STATUS_DISABLED,
		},
	}

	flattened := flattenBackupSchedules(schedules)

	expected := []interface{}{
		map[string]interface{}{
			backupScheduleIDFieldName:              "schedule-id-1",
			backupScheduleAccountIDFieldName:       "account-id-1",
			backupScheduleClusterIDFieldName:       "cluster-id-1",
			backupScheduleCronExpressionFieldName:  "0 0 * * *",
			backupScheduleRetentionPeriodFieldName: formatDuration(retentionPeriod1),
			backupScheduleCreatedAtFieldName:       formatTime(createdAt),
			backupScheduleDeletedAtFieldName:       "",
			backupScheduleStatusFieldName:          "BACKUP_SCHEDULE_STATUS_ACTIVE",
		},
		map[string]interface{}{
			backupScheduleIDFieldName:              "schedule-id-2",
			backupScheduleAccountIDFieldName:       "account-id-1",
			backupScheduleClusterIDFieldName:       "cluster-id-1",
			backupScheduleCronExpressionFieldName:  "0 12 * * *",
			backupScheduleRetentionPeriodFieldName: formatDuration(retentionPeriod2),
			backupScheduleCreatedAtFieldName:       formatTime(createdAt),
			backupScheduleDeletedAtFieldName:       "",
			backupScheduleStatusFieldName:          "BACKUP_SCHEDULE_STATUS_DISABLED",
		},
	}

	assert.Equal(t, expected, flattened)
}
