package qdrant

import (
	"testing"
	"time"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	spacebackupv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/serverless/space/backup/v1"
)

func TestServerlessRetentionPeriodValidation(t *testing.T) {
	for _, v := range []string{"24h", "72h", "8760h"} {
		assert.Empty(t, serverlessRetentionPeriodValidation(v, cty.Path{}), v)
	}
	for _, v := range []string{"1h", "23h59m", "8761h", "not-a-duration"} {
		assert.NotEmpty(t, serverlessRetentionPeriodValidation(v, cty.Path{}), v)
	}
}

func TestFlattenServerlessBackupSchedule(t *testing.T) {
	paused := timestamppb.New(time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC))
	schedule := &spacebackupv1.BackupSchedule{
		Id:              "schedule-id",
		AccountId:       "account-id",
		SpaceId:         "space-id",
		Name:            "nightly",
		Schedule:        "0 2 * * *",
		RetentionPeriod: durationpb.New(72 * time.Hour),
		CollectionName:  newPointer("c1"),
		PausedAt:        paused,
		Status:          spacebackupv1.BackupScheduleStatus_BACKUP_SCHEDULE_STATUS_ACTIVE,
	}

	flat := flattenServerlessBackupSchedule(schedule)

	assert.Equal(t, "schedule-id", flat[slBackupScheduleIDFieldName])
	assert.Equal(t, "space-id", flat[slBackupScheduleSpaceIDFieldName])
	assert.Equal(t, "nightly", flat[slBackupScheduleNameFieldName])
	assert.Equal(t, "0 2 * * *", flat[slBackupScheduleScheduleFieldName])
	assert.Equal(t, "72h0m0s", flat[slBackupScheduleRetentionPeriodFieldName])
	assert.Equal(t, "c1", flat[slBackupScheduleCollectionNameFieldName])
	assert.Equal(t, true, flat[slBackupSchedulePausedFieldName])
	assert.Equal(t, formatTime(paused), flat[slBackupSchedulePausedAtFieldName])
	assert.Equal(t, "BACKUP_SCHEDULE_STATUS_ACTIVE", flat[slBackupScheduleStatusFieldName])
	assert.Empty(t, flat[slBackupScheduleLastFiredAtFieldName])

	d := schema.TestResourceDataRaw(t, serverlessBackupScheduleSchema(false), map[string]interface{}{})
	require.NoError(t, setFlattened(d, flat))

	t.Run("when the schedule is not paused and has no status, it flattens to empty values", func(t *testing.T) {
		flat := flattenServerlessBackupSchedule(&spacebackupv1.BackupSchedule{Id: "x"})
		assert.Equal(t, false, flat[slBackupSchedulePausedFieldName])
		assert.Empty(t, flat[slBackupScheduleStatusFieldName])
		assert.Empty(t, flat[slBackupScheduleRetentionPeriodFieldName])
	})
}

func TestExpandServerlessBackupSchedule(t *testing.T) {
	t.Run("when all fields are set, it builds the full schedule", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, serverlessBackupScheduleSchema(false), map[string]interface{}{
			slBackupScheduleSpaceIDFieldName:         "space-id",
			slBackupScheduleNameFieldName:            "nightly",
			slBackupScheduleScheduleFieldName:        "0 2 * * *",
			slBackupScheduleRetentionPeriodFieldName: "72h",
			slBackupScheduleCollectionNameFieldName:  "c1",
			slBackupSchedulePausedFieldName:          true,
		})
		s := expandServerlessBackupSchedule(d, "account-id")
		assert.Equal(t, "account-id", s.GetAccountId())
		assert.Equal(t, "space-id", s.GetSpaceId())
		assert.Equal(t, "nightly", s.GetName())
		assert.Equal(t, "0 2 * * *", s.GetSchedule())
		assert.Equal(t, 72*time.Hour, s.GetRetentionPeriod().AsDuration())
		assert.Equal(t, "c1", s.GetCollectionName())
		assert.NotNil(t, s.PausedAt)
	})
	t.Run("when optional fields are not set, it leaves them unset", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, serverlessBackupScheduleSchema(false), map[string]interface{}{
			slBackupScheduleSpaceIDFieldName:  "space-id",
			slBackupScheduleNameFieldName:     "nightly",
			slBackupScheduleScheduleFieldName: "0 2 * * *",
		})
		s := expandServerlessBackupSchedule(d, "account-id")
		assert.Nil(t, s.RetentionPeriod)
		assert.Nil(t, s.CollectionName)
		assert.Nil(t, s.PausedAt)
	})
}

func TestServerlessBackupScheduleUpdateMask(t *testing.T) {
	assert.Empty(t, serverlessBackupScheduleUpdateMask(func(string) bool { return false }))
	assert.Equal(t, []string{"paused_at"}, serverlessBackupScheduleUpdateMask(func(a string) bool {
		return a == slBackupSchedulePausedFieldName
	}))
	assert.Equal(t, []string{"name", "schedule", "retention_period", "collection_name", "paused_at"},
		serverlessBackupScheduleUpdateMask(func(string) bool { return true }))
}

func TestFlattenServerlessBackup(t *testing.T) {
	backup := &spacebackupv1.Backup{
		Id:               "backup-id",
		AccountId:        "account-id",
		SpaceId:          "space-id",
		Name:             "backup-1",
		Status:           spacebackupv1.BackupStatus_BACKUP_STATUS_SUCCEEDED,
		BackupScheduleId: newPointer("schedule-id"),
		RetentionPeriod:  durationpb.New(48 * time.Hour),
		Stats: &spacebackupv1.BackupStats{
			CollectionCount: newPointer(uint32(2)),
			SizeBytes:       newPointer(int64(1024)),
			TotalPoints:     newPointer(uint64(500)),
			Duration:        durationpb.New(30 * time.Second),
		},
	}

	flat := flattenServerlessBackup(backup)

	assert.Equal(t, "backup-id", flat[slBackupIDFieldName])
	assert.Equal(t, "space-id", flat[slBackupSpaceIDFieldName])
	assert.Equal(t, "backup-1", flat[slBackupNameFieldName])
	assert.Equal(t, "BACKUP_STATUS_SUCCEEDED", flat[slBackupStatusFieldName])
	assert.Equal(t, "schedule-id", flat[slBackupBackupScheduleIDFieldName])
	assert.Equal(t, "48h0m0s", flat[slBackupRetentionPeriodFieldName])
	assert.Empty(t, flat[slBackupCollectionNameFieldName])
	assert.Equal(t, []interface{}{map[string]interface{}{
		slBackupStatsCollectionCountFieldName: 2,
		slBackupStatsSizeBytesFieldName:       1024,
		slBackupStatsTotalPointsFieldName:     500,
		slBackupStatsDurationFieldName:        "30s",
	}}, flat[slBackupStatsFieldName])

	d := schema.TestResourceDataRaw(t, serverlessBackupSchema(false), map[string]interface{}{})
	require.NoError(t, setFlattened(d, flat))

	t.Run("when the backup has no stats or status, it flattens to empty values", func(t *testing.T) {
		flat := flattenServerlessBackup(&spacebackupv1.Backup{Id: "x"})
		assert.Empty(t, flat[slBackupStatusFieldName])
		assert.Equal(t, []interface{}{}, flat[slBackupStatsFieldName])
	})
}

func TestExpandServerlessBackup(t *testing.T) {
	d := schema.TestResourceDataRaw(t, serverlessBackupSchema(false), map[string]interface{}{
		slBackupSpaceIDFieldName:         "space-id",
		slBackupCollectionNameFieldName:  "c1",
		slBackupRetentionPeriodFieldName: "48h",
	})
	b := expandServerlessBackup(d, "account-id")
	assert.Equal(t, "account-id", b.GetAccountId())
	assert.Equal(t, "space-id", b.GetSpaceId())
	assert.Equal(t, "c1", b.GetCollectionName())
	assert.Equal(t, 48*time.Hour, b.GetRetentionPeriod().AsDuration())

	d = schema.TestResourceDataRaw(t, serverlessBackupSchema(false), map[string]interface{}{
		slBackupSpaceIDFieldName: "space-id",
	})
	b = expandServerlessBackup(d, "account-id")
	assert.Nil(t, b.CollectionName)
	assert.Nil(t, b.RetentionPeriod)
}
