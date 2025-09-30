package qdrant

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	qcb "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/backup/v1"
)

func TestFlattenBackup(t *testing.T) {
	createdAt := timestamppb.New(time.Date(2025, 9, 24, 13, 20, 23, 663911000, time.UTC))
	deletedAt := timestamppb.New(time.Date(2025, 9, 25, 7, 5, 0, 0, time.UTC))

	b := &qcb.Backup{
		Id:             "b088d7d3-2ba9-4839-8d6d-f04db6ec14dd",
		CreatedAt:      createdAt,
		DeletedAt:      deletedAt,
		AccountId:      "222cda33-2c7a-4046-b2cc-0807170aed49",
		ClusterId:      "604fa4fc-fdd2-4ae9-ac36-d166a114d52f",
		Name:           "qdrant-604fa4fc-fdd2-4ae9-ac36-d166a114d52f-snapshot-1758720023",
		Status:         qcb.BackupStatus_BACKUP_STATUS_SUCCEEDED,
		BackupDuration: durationpb.New(68 * time.Second),
		ClusterInfo: &qcb.ClusterInfo{
			Name:                  "weary-wombat-bronze",
			CloudProviderId:       "aws",
			CloudProviderRegionId: "eu-central-1",
		},
	}

	got := flattenBackup(b)

	// got is map[string]interface{}
	// ignoring cluster config
	if clusterInfos, ok := got[backupClusterInfoFieldName].([]interface{}); ok && len(clusterInfos) > 0 {
		if clusterInfo, ok := clusterInfos[0].(map[string]interface{}); ok {
			delete(clusterInfo, bClusterCfgField)
		}
	}

	want := map[string]interface{}{
		backupIdFieldName:         b.GetId(),
		backupAccountIdFieldName:  b.GetAccountId(),
		backupClusterIdFieldName:  b.GetClusterId(),
		backupNameFieldName:       b.GetName(),
		backupStatusFieldName:     b.GetStatus().String(),
		backupDurationFieldName:   formatDuration(b.GetBackupDuration()),
		backupScheduleIdFieldName: b.GetBackupScheduleId(),
		backupCreatedAtFieldName:  formatTime(b.GetCreatedAt()),
		backupDeletedAtFieldName:  formatTime(b.GetDeletedAt()),
		backupClusterInfoFieldName: []interface{}{
			map[string]interface{}{
				bClusterInfoNameField:          b.GetClusterInfo().GetName(),
				bClusterInfoCloudProviderField: b.GetClusterInfo().GetCloudProviderId(),
				bClusterInfoRegionField:        b.GetClusterInfo().GetCloudProviderRegionId(),
			},
		},
	}

	assert.Equal(t, want, got)
}

func TestExpandBackupForCreate_UsesDefaultAccountID(t *testing.T) {
	const defaultAcct = "00000000-1000-0000-0000-000000000001"
	d := schema.TestResourceDataRaw(t, accountsBackupSchema(), map[string]interface{}{
		backupClusterIdFieldName: "604fa4fc-fdd2-4ae9-ac36-d166a114d52f",
	})

	backup, err := expandBackup(d, defaultAcct)
	require.NoError(t, err)
	assert.Equal(t, defaultAcct, backup.GetAccountId())
	assert.Equal(t, "604fa4fc-fdd2-4ae9-ac36-d166a114d52f", backup.GetClusterId())

	// read-only fields should not be set by expand
	assert.Nil(t, backup.GetCreatedAt())
	assert.Empty(t, backup.GetName())
	assert.Equal(t, qcb.BackupStatus_BACKUP_STATUS_UNSPECIFIED.String(), backup.GetStatus().String())
	assert.Nil(t, backup.GetBackupDuration())
	assert.Empty(t, backup.GetBackupScheduleId())
	assert.Nil(t, backup.GetClusterInfo())
}

func TestExpandBackupForCreate_WithOverrideAccountID(t *testing.T) {
	overrideAcct := "00000000-2000-0000-0000-000000000002"
	d := schema.TestResourceDataRaw(t, accountsBackupSchema(), map[string]interface{}{
		backupAccountIdFieldName: overrideAcct,
		backupClusterIdFieldName: "604fa4fc-fdd2-4ae9-ac36-d166a114d52f",
	})

	backup, err := expandBackup(d, "ignored-default")
	require.NoError(t, err)
	assert.Equal(t, overrideAcct, backup.GetAccountId())
	assert.Equal(t, "604fa4fc-fdd2-4ae9-ac36-d166a114d52f", backup.GetClusterId())
}

func TestExpandBackupForCreate_MissingClusterIDErrors(t *testing.T) {
	d := schema.TestResourceDataRaw(t, accountsBackupSchema(), map[string]interface{}{})
	_, err := expandBackup(d, "00000000-1000-0000-0000-000000000001")
	require.Error(t, err)
}

func TestSchema_BackupAccountIdOptionalComputed(t *testing.T) {
	s := accountsBackupSchema()
	field := s[backupAccountIdFieldName]
	require.NotNil(t, field)

	assert.True(t, field.Optional, "account_id should be Optional")
	assert.True(t, field.Computed, "account_id should be Computed")
	assert.Nil(t, field.Default, "account_id must NOT set a Default")
}

func TestSchema_BackupClusterIdForceNew(t *testing.T) {
	s := accountsBackupSchema()
	field := s[backupClusterIdFieldName]
	require.NotNil(t, field)

	assert.True(t, field.ForceNew, "cluster_id must be ForceNew")
}

func TestSchema_BackupRetentionPeriodField(t *testing.T) {
	s := accountsBackupSchema()
	f := s[backupRetentionPeriodFieldName]
	require.NotNil(t, f, "retention_period schema must exist")

	assert.Equal(t, schema.TypeString, f.Type, "retention_period should be a string (Go duration)")
	assert.True(t, f.Optional, "retention_period should be Optional")
	assert.True(t, f.ForceNew, "retention_period should be ForceNew on backups")
	require.NotNil(t, f.DiffSuppressFunc, "retention_period should have DiffSuppressFunc")

	// Prove DiffSuppress collapses semantically equal durations.
	d := schema.TestResourceDataRaw(t, accountsBackupSchema(), nil)
	assert.True(t, f.DiffSuppressFunc(backupRetentionPeriodFieldName, "60m0s", "1h0m0s", d))
	assert.True(t, f.DiffSuppressFunc(backupRetentionPeriodFieldName, "3600s", "1h0m0s", d))
	assert.True(t, f.DiffSuppressFunc(backupRetentionPeriodFieldName, "24h0m0s", "1440m0s", d))

	// And different durations should not be suppressed.
	assert.False(t, f.DiffSuppressFunc(backupRetentionPeriodFieldName, "30m0s", "1h0m0s", d))
}

func TestExpandBackupForCreate_WithRetentionPeriod(t *testing.T) {
	const defaultAcct = "00000000-1000-0000-0000-000000000001"
	d := schema.TestResourceDataRaw(t, accountsBackupSchema(), map[string]interface{}{
		backupClusterIdFieldName:       "604fa4fc-fdd2-4ae9-ac36-d166a114d52f",
		backupRetentionPeriodFieldName: "24h",
	})

	bk, err := expandBackup(d, defaultAcct)
	require.NoError(t, err)

	rp := bk.GetRetentionPeriod()
	require.NotNil(t, rp, "retention period should be set on the request")
	assert.Equal(t, 24*time.Hour, rp.AsDuration())
}

func TestExpandBackupForCreate_InvalidRetentionPeriodErrors(t *testing.T) {
	const defaultAcct = "00000000-1000-0000-0000-000000000001"
	d := schema.TestResourceDataRaw(t, accountsBackupSchema(), map[string]interface{}{
		backupClusterIdFieldName:       "604fa4fc-fdd2-4ae9-ac36-d166a114d52f",
		backupRetentionPeriodFieldName: "definitely-not-a-duration",
	})

	_, err := expandBackup(d, defaultAcct)
	require.Error(t, err, "invalid retention_period should error")
}

func TestExpandBackupForCreate_NoRetentionPeriod_Unset(t *testing.T) {
	const defaultAcct = "00000000-1000-0000-0000-000000000001"
	d := schema.TestResourceDataRaw(t, accountsBackupSchema(), map[string]interface{}{
		backupClusterIdFieldName: "604fa4fc-fdd2-4ae9-ac36-d166a114d52f",
	})

	bk, err := expandBackup(d, defaultAcct)
	require.NoError(t, err)
	assert.Nil(t, bk.GetRetentionPeriod(), "retention period should be nil when not provided")
}

func TestFlattenBackup_WithRetentionPeriod(t *testing.T) {
	createdAt := timestamppb.New(time.Date(2025, 9, 24, 13, 20, 23, 0, time.UTC))

	b := &qcb.Backup{
		Id:              "b088d7d3-2ba9-4839-8d6d-f04db6ec14dd",
		AccountId:       "222cda33-2c7a-4046-b2cc-0807170aed49",
		ClusterId:       "604fa4fc-fdd2-4ae9-ac36-d166a114d52f",
		Name:            "snapshot-with-retention",
		Status:          qcb.BackupStatus_BACKUP_STATUS_SUCCEEDED,
		BackupDuration:  durationpb.New(10 * time.Second),
		CreatedAt:       createdAt,
		RetentionPeriod: durationpb.New(36 * time.Hour),
		ClusterInfo: &qcb.ClusterInfo{
			Name:                  "weary-wombat-bronze",
			CloudProviderId:       "aws",
			CloudProviderRegionId: "eu-central-1",
		},
	}

	got := flattenBackup(b)

	// ignore cluster configuration noise for equality checks (same as your other test)
	if clusterInfos, ok := got[backupClusterInfoFieldName].([]interface{}); ok && len(clusterInfos) > 0 {
		if clusterInfo, ok := clusterInfos[0].(map[string]interface{}); ok {
			delete(clusterInfo, bClusterCfgField)
		}
	}

	require.Contains(t, got, backupRetentionPeriodFieldName)
	// formatDuration prints zero components, so expect "36h0m0s" rather than just "36h".
	assert.Equal(t, "36h0m0s", got[backupRetentionPeriodFieldName])
}
