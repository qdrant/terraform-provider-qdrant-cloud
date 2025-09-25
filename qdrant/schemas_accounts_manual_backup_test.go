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
	qccl "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/v1"
	qcco "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/common/v1"
)

func TestFlattenBackup(t *testing.T) {
	createdAt := timestamppb.New(time.Date(2025, 9, 24, 13, 20, 23, 663911000, time.UTC))
	deletedAt := timestamppb.New(time.Date(2025, 9, 25, 7, 5, 0, 0, time.UTC))
	cfgLastMod := timestamppb.New(time.Date(2025, 9, 24, 12, 34, 37, 505702000, time.UTC))
	version := "v1.15.4"

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
			Configuration: &qccl.ClusterConfiguration{
				LastModifiedAt: cfgLastMod,
				NumberOfNodes:  1,
				Version:        &version,
				PackageId:      "7c939d96-d671-4051-aa16-3b8b7130fa42",
				DatabaseConfiguration: &qccl.DatabaseConfiguration{
					Service: &qccl.DatabaseConfigurationService{
						ApiKey:  &qcco.SecretKeyRef{Name: "qdrant-api-key-604fa4fc-fdd2-4ae9-ac36-d166a114d52f", Key: "api-key"},
						JwtRbac: true,
					},
					Inference: &qccl.DatabaseConfigurationInference{
						Enabled: true,
					},
				},
			},
		},
	}

	got := flattenBackup(b)
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
				bClusterCfgField: []interface{}{
					map[string]interface{}{
						bClusterCfgLastModifiedAtField: formatTime(cfgLastMod),
						bClusterCfgNumberOfNodesField:  int(b.GetClusterInfo().GetConfiguration().GetNumberOfNodes()),
						bClusterCfgVersionField:        b.GetClusterInfo().GetConfiguration().GetVersion(),
						bClusterCfgPackageIdField:      b.GetClusterInfo().GetConfiguration().GetPackageId(),
						bClusterCfgServiceTypeField:    b.GetClusterInfo().GetConfiguration().GetServiceType().String(),
						bClusterCfgRebalanceStratField: b.GetClusterInfo().GetConfiguration().GetRebalanceStrategy().String(),
						bClusterCfgDbConfigField: []interface{}{
							map[string]interface{}{
								bDbCfgServiceField: []interface{}{
									map[string]interface{}{
										bDbCfgServiceApiKeyField: []interface{}{
											map[string]interface{}{
												bDbCfgServiceApiKeyNameField: "qdrant-api-key-604fa4fc-fdd2-4ae9-ac36-d166a114d52f",
												bDbCfgServiceApiKeyKeyField:  "api-key",
											},
										},
										bDbCfgServiceJwtRbacField: true,
									},
								},
								bDbCfgInferenceField: []interface{}{
									map[string]interface{}{
										bDbCfgInferenceEnabledField: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, want, got)
}

func TestExpandBackupForCreate_UsesDefaultAccountID(t *testing.T) {
	defaultAcct := "00000000-1000-0000-0000-000000000001"
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
