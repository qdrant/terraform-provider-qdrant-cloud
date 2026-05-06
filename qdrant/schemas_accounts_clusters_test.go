package qdrant

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	qcCluster "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/v1"
	commonv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/common/v1"
)

func TestResourceClusterFlatten(t *testing.T) {
	cluster := &qcCluster.Cluster{
		Id:        "00000000-0000-0000-0000-000000000001",
		CreatedAt: timestamppb.New(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
		AccountId: "00000000-1000-0000-0000-000000000001",
		Name:      "testName",
		Labels: []*commonv1.KeyValue{
			{Key: "key1", Value: "value1"},
		},
		CloudProviderId:       "Azure",
		CloudProviderRegionId: "Uksouth",
		DeletedAt:             timestamppb.New(time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)),
		Configuration: &qcCluster.ClusterConfiguration{
			LastModifiedAt: timestamppb.New(time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
			Version:        newPointer("v1.0"),
			NumberOfNodes:  5,
			PackageId:      "00000009-1000-0000-0000-000000000001",
			DatabaseConfiguration: &qcCluster.DatabaseConfiguration{
				Collection: &qcCluster.DatabaseConfigurationCollection{
					ReplicationFactor:      newPointer(uint32(2)),
					WriteConsistencyFactor: newPointer(int32(1)),
				},
				Storage: &qcCluster.DatabaseConfigurationStorage{
					Performance: &qcCluster.DatabaseConfigurationStoragePerformance{
						OptimizerCpuBudget: newPointer(int32(10)),
						AsyncScorer:        newPointer(true),
					},
				},
				Service: &qcCluster.DatabaseConfigurationService{
					ApiKey:         &commonv1.SecretKeyRef{Name: "api-key-secret", Key: "api-key"},
					ReadOnlyApiKey: &commonv1.SecretKeyRef{Name: "ro-api-key-secret", Key: "ro-api-key"},
				},
				LogLevel: newPointer(qcCluster.DatabaseConfigurationLogLevel_DATABASE_CONFIGURATION_LOG_LEVEL_DEBUG),
				Tls: &qcCluster.DatabaseConfigurationTls{
					Cert: &commonv1.SecretKeyRef{Name: "cert-secret", Key: "cert.pem"},
					Key:  &commonv1.SecretKeyRef{Name: "key-secret", Key: "key.pem"},
				},
				Inference: &qcCluster.DatabaseConfigurationInference{Enabled: true},
				AuditLogging: &qcCluster.DatabaseConfigurationAuditLogging{
					Enabled:               true,
					Rotation:              newPointer(qcCluster.AuditLogRotation_AUDIT_LOG_ROTATION_DAILY),
					MaxLogFiles:           newPointer(uint32(7)),
					TrustForwardedHeaders: newPointer(true),
				},
			},
			AdditionalResources: &qcCluster.AdditionalResources{
				Disk: 8,
			},
			NodeSelector: []*commonv1.KeyValue{
				{Key: "key1", Value: "value1"},
			},
			Tolerations: []*qcCluster.Toleration{
				{
					Key:      newPointer("key1"),
					Operator: newPointer(qcCluster.TolerationOperator_TOLERATION_OPERATOR_EQUAL),
					Value:    newPointer("value1"),
					Effect:   newPointer(qcCluster.TolerationEffect_TOLERATION_EFFECT_NO_SCHEDULE),
				},
			},
			TopologySpreadConstraints: []*commonv1.TopologySpreadConstraint{
				{
					MaxSkew:           newPointer(int32(1)),
					TopologyKey:       "topology.kubernetes.io/zone",
					WhenUnsatisfiable: newPointer(commonv1.TopologySpreadConstraintWhenUnsatisfiable_TOPOLOGY_SPREAD_CONSTRAINT_WHEN_UNSATISFIABLE_DO_NOT_SCHEDULE),
				},
			},
			Annotations: []*commonv1.KeyValue{
				{Key: "anno1", Value: "annoval1"},
			},
			AllowedIpSourceRanges: []string{"192.168.1.0/24"},
			ServiceType:           newPointer(qcCluster.ClusterServiceType_CLUSTER_SERVICE_TYPE_LOAD_BALANCER),
			ServiceAnnotations: []*commonv1.KeyValue{
				{Key: "serviceanno1", Value: "serviceannoval1"},
			},
			PodLabels: []*commonv1.KeyValue{
				{Key: "podlabel1", Value: "podlabelval1"},
			},
			ReservedCpuPercentage:    newPointer(uint32(10)),
			ReservedMemoryPercentage: newPointer(uint32(20)),
			GpuType:                  newPointer(qcCluster.ClusterConfigurationGpuType_CLUSTER_CONFIGURATION_GPU_TYPE_NVIDIA),
			RestartPolicy:            newPointer(qcCluster.ClusterConfigurationRestartPolicy_CLUSTER_CONFIGURATION_RESTART_POLICY_ROLLING),
			RebalanceStrategy:        newPointer(qcCluster.ClusterConfigurationRebalanceStrategy_CLUSTER_CONFIGURATION_REBALANCE_STRATEGY_BY_COUNT),
			ClusterStorageConfiguration: &qcCluster.ClusterStorageConfiguration{
				StorageTierType: commonv1.StorageTierType_STORAGE_TIER_TYPE_BALANCED,
			},
		},
		State: &qcCluster.ClusterState{
			Version:     "v1.1.1",
			NodesUp:     5,
			RestartedAt: timestamppb.New(time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC)),
			Phase:       qcCluster.ClusterPhase_CLUSTER_PHASE_HEALTHY,
			Reason:      "All good",
			Endpoint: &qcCluster.ClusterEndpoint{
				Url: "http://example.com",
			},
			Resources: &qcCluster.ClusterNodeResourcesSummary{
				Disk: &qcCluster.ClusterNodeResources{Base: 100, Complimentary: 10, Additional: 8, Reserved: 5, Available: 113},
				Ram:  &qcCluster.ClusterNodeResources{Base: 8, Complimentary: 1, Additional: 0, Reserved: 0.5, Available: 8.5},
				Cpu:  &qcCluster.ClusterNodeResources{Base: 4, Complimentary: 0, Additional: 0, Reserved: 0.2, Available: 3.8},
			},
			// ScalabilityInfo is deprecated, use Capabilities instead.
			// Note: The old ScalabilityInfo is still populated by the API for backward compatibility.
			ScalabilityInfo: &qcCluster.ClusterScalabilityInfo{
				Status: qcCluster.ClusterScalabilityStatus_CLUSTER_SCALABILITY_STATUS_SCALABLE,
				Reason: newPointer("Can be scaled"),
			},
			Capabilities: &qcCluster.ClusterCapabilities{
				DiskExpansion: &qcCluster.ClusterDiskExpansionSupportInfo{
					Status: qcCluster.ClusterDiskExpansionSupportStatus_CLUSTER_DISK_EXPANSION_SUPPORT_STATUS_SUPPORTED,
					Reason: newPointer("Disk can be expanded"),
				},
				Backup: &qcCluster.ClusterBackupSupportInfo{
					Status: qcCluster.ClusterBackupSupportStatus_CLUSTER_BACKUP_SUPPORT_STATUS_NOT_SUPPORTED,
					Reason: newPointer("Backup not supported on this tier"),
				},
				ScalabilityInfo: &qcCluster.ClusterScalabilityInfo{
					Status: qcCluster.ClusterScalabilityStatus_CLUSTER_SCALABILITY_STATUS_SCALABLE,
					Reason: newPointer("Can be scaled"),
				},
			},
		},
	}

	flattened := flattenCluster(cluster)

	expected := map[string]interface{}{
		clusterIdentifierFieldName: cluster.GetId(),
		clusterCreatedAtFieldName:  formatTime(cluster.GetCreatedAt()),
		clusterAccountIDFieldName:  cluster.GetAccountId(),
		clusterNameFieldName:       cluster.GetName(),
		clusterLabelsFieldName: []interface{}{
			map[string]interface{}{"key": "key1", "value": "value1"},
		},
		clusterCloudProviderFieldName:       cluster.GetCloudProviderId(),
		clusterCloudRegionFieldName:         cluster.GetCloudProviderRegionId(),
		clusterPrivateRegionIDFieldName:     "",
		clusterMarkedForDeletionAtFieldName: formatTime(cluster.GetDeletedAt()),
		clusterURLFieldName:                 cluster.GetState().GetEndpoint().GetUrl(),
		clusterStatusFieldName: []interface{}{
			map[string]interface{}{
				clusterStatusVersionFieldName:     cluster.GetState().GetVersion(),
				clusterStatusNodesUpFieldName:     int(cluster.GetState().GetNodesUp()),
				clusterStatusRestartedAtFieldName: formatTime(cluster.GetState().GetRestartedAt()),
				clusterStatusPhaseFieldName:       cluster.GetState().GetPhase().String(),
				clusterStatusReasonFieldName:      cluster.GetState().GetReason(),
				clusterStatusResourcesFieldName: []interface{}{
					map[string]interface{}{
						clusterNodeResourcesSummaryDiskFieldName: []interface{}{
							map[string]interface{}{
								clusterNodeResourcesBaseFieldName:          100.0,
								clusterNodeResourcesComplimentaryFieldName: 10.0,
								clusterNodeResourcesAdditionalFieldName:    8.0,
								clusterNodeResourcesReservedFieldName:      5.0,
								clusterNodeResourcesAvailableFieldName:     113.0,
							},
						},
						clusterNodeResourcesSummaryRamFieldName: []interface{}{
							map[string]interface{}{
								clusterNodeResourcesBaseFieldName:          8.0,
								clusterNodeResourcesComplimentaryFieldName: 1.0,
								clusterNodeResourcesAdditionalFieldName:    0.0,
								clusterNodeResourcesReservedFieldName:      0.5,
								clusterNodeResourcesAvailableFieldName:     8.5,
							},
						},
						clusterNodeResourcesSummaryCpuFieldName: []interface{}{
							map[string]interface{}{
								clusterNodeResourcesBaseFieldName:          4.0,
								clusterNodeResourcesComplimentaryFieldName: 0.0,
								clusterNodeResourcesAdditionalFieldName:    0.0,
								clusterNodeResourcesReservedFieldName:      0.2,
								clusterNodeResourcesAvailableFieldName:     3.8,
							},
						},
					},
				},
				clusterStatusScalabilityInfoFieldName: []interface{}{
					map[string]interface{}{
						clusterScalabilityInfoStatusFieldName: cluster.GetState().GetCapabilities().GetScalabilityInfo().GetStatus().String(),
						clusterScalabilityInfoReasonFieldName: cluster.GetState().GetCapabilities().GetScalabilityInfo().GetReason(),
					},
				},
				clusterStatusCapabilitiesFieldName: []interface{}{
					map[string]interface{}{
						clusterCapabilitiesDiskExpansionFieldName: []interface{}{
							map[string]interface{}{
								clusterCapabilityStatusFieldName: cluster.GetState().GetCapabilities().GetDiskExpansion().GetStatus().String(),
								clusterCapabilityReasonFieldName: cluster.GetState().GetCapabilities().GetDiskExpansion().GetReason(),
							},
						},
						clusterCapabilitiesBackupFieldName: []interface{}{
							map[string]interface{}{
								clusterCapabilityStatusFieldName: cluster.GetState().GetCapabilities().GetBackup().GetStatus().String(),
								clusterCapabilityReasonFieldName: cluster.GetState().GetCapabilities().GetBackup().GetReason(),
							},
						},
						clusterStatusScalabilityInfoFieldName: []interface{}{
							map[string]interface{}{
								clusterScalabilityInfoStatusFieldName: cluster.GetState().GetCapabilities().GetScalabilityInfo().GetStatus().String(),
								clusterScalabilityInfoReasonFieldName: cluster.GetState().GetCapabilities().GetScalabilityInfo().GetReason(),
							},
						},
					},
				},
			},
		},
		configurationFieldName: []interface{}{
			map[string]interface{}{
				clusterVersionFieldName:        cluster.GetConfiguration().GetVersion(),
				clusterLastModifiedAtFieldName: formatTime(cluster.GetConfiguration().GetLastModifiedAt()),
				numberOfNodesFieldName:         int(cluster.GetConfiguration().GetNumberOfNodes()),
				nodeConfigurationFieldName: []interface{}{
					map[string]interface{}{
						packageIDFieldName: cluster.GetConfiguration().GetPackageId(),
						resourceConfigurationsFieldName: []interface{}{
							map[string]interface{}{
								fieldAmount:       int(cluster.GetConfiguration().GetAdditionalResources().GetDisk()),
								fieldResourceUnit: string(ResourceUnitGi),
								fieldResourceType: string(ResourceTypeDisk),
							},
						},
					},
				},
				databaseConfigurationFieldName: []interface{}{
					map[string]interface{}{
						dbConfigCollectionFieldName: []interface{}{
							map[string]interface{}{
								dbConfigCollectionReplicationFactor:      2,
								dbConfigCollectionWriteConsistencyFactor: 1,
							},
						},
						dbConfigStorageFieldName: []interface{}{
							map[string]interface{}{
								dbConfigStoragePerformanceFieldName: []interface{}{
									map[string]interface{}{
										dbConfigStoragePerfOptimizerCpuBudget: 10,
										dbConfigStoragePerfAsyncScorer:        true,
									},
								},
							},
						},
						dbConfigServiceFieldName: []interface{}{
							map[string]interface{}{
								dbConfigServiceApiKeyFieldName: []interface{}{
									map[string]interface{}{
										dbConfigSecretKeyRefSecretNameFieldName: "api-key-secret",
										dbConfigSecretKeyRefSecretKeyFieldName:  "api-key",
									},
								},
								dbConfigServiceReadOnlyApiKeyFieldName: []interface{}{
									map[string]interface{}{
										dbConfigSecretKeyRefSecretNameFieldName: "ro-api-key-secret",
										dbConfigSecretKeyRefSecretKeyFieldName:  "ro-api-key",
									},
								},
								dbConfigServiceJwtRbacFieldName: false,
							},
						},
						dbConfigLogLevelFieldName: "DATABASE_CONFIGURATION_LOG_LEVEL_DEBUG",
						dbConfigTlsFieldName: []interface{}{
							map[string]interface{}{
								dbConfigTlsCertFieldName: []interface{}{
									map[string]interface{}{
										dbConfigSecretKeyRefSecretNameFieldName: "cert-secret",
										dbConfigSecretKeyRefSecretKeyFieldName:  "cert.pem",
									},
								},
								dbConfigTlsKeyFieldName: []interface{}{
									map[string]interface{}{
										dbConfigSecretKeyRefSecretNameFieldName: "key-secret",
										dbConfigSecretKeyRefSecretKeyFieldName:  "key.pem",
									},
								},
							},
						},
						dbConfigInferenceFieldName: []interface{}{
							map[string]interface{}{
								dbConfigInferenceEnabledFieldName: true,
							},
						},
						dbConfigAuditLoggingFieldName: []interface{}{
							map[string]interface{}{
								dbConfigAuditLoggingEnabledFieldName:               true,
								dbConfigAuditLoggingRotationFieldName:              "AUDIT_LOG_ROTATION_DAILY",
								dbConfigAuditLoggingMaxLogFilesFieldName:           7,
								dbConfigAuditLoggingTrustForwardedHeadersFieldName: true,
							},
						},
					},
				},
				nodeSelectorFieldName: []interface{}{
					map[string]interface{}{"key": "key1", "value": "value1"},
				},
				tolerationsFieldName: []interface{}{
					map[string]interface{}{
						tolerationKeyFieldName:      "key1",
						tolerationOperatorFieldName: "TOLERATION_OPERATOR_EQUAL",
						tolerationValueFieldName:    "value1",
						tolerationEffectFieldName:   "TOLERATION_EFFECT_NO_SCHEDULE",
					},
				},
				topologySpreadConstraintsFieldName: []interface{}{
					map[string]interface{}{
						topologySpreadConstraintMaxSkewFieldName:           1,
						topologySpreadConstraintTopologyKeyFieldName:       "topology.kubernetes.io/zone",
						topologySpreadConstraintWhenUnsatisfiableFieldName: "TOPOLOGY_SPREAD_CONSTRAINT_WHEN_UNSATISFIABLE_DO_NOT_SCHEDULE",
					},
				},
				annotationsFieldName: []interface{}{
					map[string]interface{}{"key": "anno1", "value": "annoval1"},
				},
				allowedIpSourceRangesFieldName: []string{"192.168.1.0/24"},
				serviceTypeFieldName:           "CLUSTER_SERVICE_TYPE_LOAD_BALANCER",
				serviceAnnotationsFieldName: []interface{}{
					map[string]interface{}{"key": "serviceanno1", "value": "serviceannoval1"},
				},
				podLabelsFieldName: []interface{}{
					map[string]interface{}{"key": "podlabel1", "value": "podlabelval1"},
				},
				dbConfigReservedCpuPercentageFieldName:    10,
				dbConfigReservedMemoryPercentageFieldName: 20,
				dbConfigGpuTypeFieldName:                  "CLUSTER_CONFIGURATION_GPU_TYPE_NVIDIA",
				dbConfigRestartPolicyFieldName:            "CLUSTER_CONFIGURATION_RESTART_POLICY_ROLLING",
				dbConfigRebalanceStrategyFieldName:        "CLUSTER_CONFIGURATION_REBALANCE_STRATEGY_BY_COUNT",
				clusterStorageConfigurationFieldName: []interface{}{
					map[string]interface{}{
						clusterStorageTierTypeFieldName: "STORAGE_TIER_TYPE_BALANCED",
					},
				},
			},
		},
	}

	assert.Equal(t, expected, flattened)
}

func TestExpandCluster(t *testing.T) {
	expected := &qcCluster.Cluster{
		Id:        "00000000-0000-0000-0000-000000000001",
		CreatedAt: timestamppb.New(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
		AccountId: "00000000-1000-0000-0000-000000000001",
		Name:      "testName",
		Labels: []*commonv1.KeyValue{
			{Key: "key1", Value: "value1"},
		},
		CloudProviderId:       "Azure",
		CloudProviderRegionId: "Uksouth",
		DeletedAt:             timestamppb.New(time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)),
		Configuration: &qcCluster.ClusterConfiguration{
			LastModifiedAt: timestamppb.New(time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
			Version:        newPointer("v1.0"),
			NumberOfNodes:  5,
			PackageId:      "00000009-1000-0000-0000-000000000001",
			AdditionalResources: &qcCluster.AdditionalResources{
				Disk: 10,
			},
			DatabaseConfiguration: &qcCluster.DatabaseConfiguration{
				Collection: &qcCluster.DatabaseConfigurationCollection{
					ReplicationFactor:      newPointer(uint32(3)),
					WriteConsistencyFactor: newPointer(int32(2)),
				},
				Service: &qcCluster.DatabaseConfigurationService{
					ApiKey:    &commonv1.SecretKeyRef{Name: "api-key-secret-expand", Key: "api-key-expand"},
					EnableTls: newPointer(false)},
				Tls: &qcCluster.DatabaseConfigurationTls{
					Cert: &commonv1.SecretKeyRef{Name: "cert-secret-expand", Key: "cert.pem-expand"},
				},
				AuditLogging: &qcCluster.DatabaseConfigurationAuditLogging{
					Enabled:               true,
					Rotation:              newPointer(qcCluster.AuditLogRotation_AUDIT_LOG_ROTATION_HOURLY),
					MaxLogFiles:           newPointer(uint32(14)),
					TrustForwardedHeaders: newPointer(false),
				},
			},
			NodeSelector: []*commonv1.KeyValue{
				{Key: "key1", Value: "value1"},
			},
			Tolerations: []*qcCluster.Toleration{
				{
					Key:      newPointer("key1"),
					Operator: newPointer(qcCluster.TolerationOperator_TOLERATION_OPERATOR_EQUAL),
					Value:    newPointer("value1"),
					Effect:   newPointer(qcCluster.TolerationEffect_TOLERATION_EFFECT_NO_SCHEDULE),
				},
			},
			TopologySpreadConstraints: []*commonv1.TopologySpreadConstraint{
				{
					MaxSkew:           newPointer(int32(1)),
					TopologyKey:       "topology.kubernetes.io/zone",
					WhenUnsatisfiable: newPointer(commonv1.TopologySpreadConstraintWhenUnsatisfiable_TOPOLOGY_SPREAD_CONSTRAINT_WHEN_UNSATISFIABLE_DO_NOT_SCHEDULE),
				},
			},
			ServiceType: newPointer(qcCluster.ClusterServiceType_CLUSTER_SERVICE_TYPE_LOAD_BALANCER),
			GpuType:     newPointer(qcCluster.ClusterConfigurationGpuType_CLUSTER_CONFIGURATION_GPU_TYPE_NVIDIA),
			ClusterStorageConfiguration: &qcCluster.ClusterStorageConfiguration{
				StorageTierType: commonv1.StorageTierType_STORAGE_TIER_TYPE_PERFORMANCE,
			},
		},
		State: &qcCluster.ClusterState{
			Endpoint: &qcCluster.ClusterEndpoint{
				Url: "http://example.com",
			},
		},
	}

	d := schema.TestResourceDataRaw(t, accountsClusterSchema(false), map[string]interface{}{
		clusterIdentifierFieldName: expected.GetId(),
		clusterCreatedAtFieldName:  formatTime(expected.GetCreatedAt()),
		clusterAccountIDFieldName:  expected.GetAccountId(),
		clusterNameFieldName:       expected.GetName(),
		clusterLabelsFieldName: []interface{}{
			map[string]interface{}{"key": "key1", "value": "value1"},
		},
		clusterCloudProviderFieldName:       expected.GetCloudProviderId(),
		clusterCloudRegionFieldName:         expected.GetCloudProviderRegionId(),
		clusterPrivateRegionIDFieldName:     "",
		clusterMarkedForDeletionAtFieldName: formatTime(expected.GetDeletedAt()),
		clusterURLFieldName:                 expected.GetState().GetEndpoint().GetUrl(),
		configurationFieldName: []interface{}{
			map[string]interface{}{
				clusterVersionFieldName:        expected.GetConfiguration().GetVersion(),
				clusterLastModifiedAtFieldName: formatTime(expected.GetConfiguration().GetLastModifiedAt()),
				numberOfNodesFieldName:         int(expected.GetConfiguration().GetNumberOfNodes()),
				nodeConfigurationFieldName: []interface{}{
					map[string]interface{}{
						packageIDFieldName: expected.GetConfiguration().GetPackageId(),
						resourceConfigurationsFieldName: []interface{}{
							map[string]interface{}{
								fieldAmount:       int(expected.GetConfiguration().GetAdditionalResources().GetDisk()),
								fieldResourceUnit: string(ResourceUnitGi),
								fieldResourceType: string(ResourceTypeDisk),
							},
						},
					},
				},
				databaseConfigurationFieldName: []interface{}{
					map[string]interface{}{
						dbConfigCollectionFieldName: []interface{}{
							map[string]interface{}{
								dbConfigCollectionReplicationFactor:      3,
								dbConfigCollectionWriteConsistencyFactor: 2,
							},
						},
						dbConfigServiceFieldName: []interface{}{
							map[string]interface{}{
								dbConfigServiceApiKeyFieldName: []interface{}{
									map[string]interface{}{
										dbConfigSecretKeyRefSecretNameFieldName: "api-key-secret-expand",
										dbConfigSecretKeyRefSecretKeyFieldName:  "api-key-expand",
									},
								},
							},
						},
						dbConfigTlsFieldName: []interface{}{
							map[string]interface{}{
								dbConfigTlsCertFieldName: []interface{}{
									map[string]interface{}{
										dbConfigSecretKeyRefSecretNameFieldName: "cert-secret-expand",
										dbConfigSecretKeyRefSecretKeyFieldName:  "cert.pem-expand",
									},
								},
							},
						},
						dbConfigAuditLoggingFieldName: []interface{}{
							map[string]interface{}{
								dbConfigAuditLoggingEnabledFieldName:               true,
								dbConfigAuditLoggingRotationFieldName:              "AUDIT_LOG_ROTATION_HOURLY",
								dbConfigAuditLoggingMaxLogFilesFieldName:           14,
								dbConfigAuditLoggingTrustForwardedHeadersFieldName: false,
							},
						},
					},
				},
				nodeSelectorFieldName: []interface{}{
					map[string]interface{}{"key": "key1", "value": "value1"},
				},
				tolerationsFieldName: []interface{}{
					map[string]interface{}{
						tolerationKeyFieldName:      "key1",
						tolerationOperatorFieldName: "TOLERATION_OPERATOR_EQUAL",
						tolerationValueFieldName:    "value1",
						tolerationEffectFieldName:   "TOLERATION_EFFECT_NO_SCHEDULE",
					},
				},
				topologySpreadConstraintsFieldName: []interface{}{
					map[string]interface{}{
						topologySpreadConstraintMaxSkewFieldName:           1,
						topologySpreadConstraintTopologyKeyFieldName:       "topology.kubernetes.io/zone",
						topologySpreadConstraintWhenUnsatisfiableFieldName: "TOPOLOGY_SPREAD_CONSTRAINT_WHEN_UNSATISFIABLE_DO_NOT_SCHEDULE",
					},
				},
				serviceTypeFieldName:     "CLUSTER_SERVICE_TYPE_LOAD_BALANCER",
				dbConfigGpuTypeFieldName: "CLUSTER_CONFIGURATION_GPU_TYPE_NVIDIA",
				clusterStorageConfigurationFieldName: []interface{}{
					map[string]interface{}{
						clusterStorageTierTypeFieldName: "STORAGE_TIER_TYPE_PERFORMANCE",
					},
				},
			},
		},
	})

	result, jwtRbac, err := expandCluster(d, expected.GetAccountId())
	require.NoError(t, err)
	assert.Equal(t, expected, result)
	assert.Nil(t, jwtRbac)
}

// TestFlattenClusterConfigurationUnspecifiedEnums verifies that UNSPECIFIED enum values
// are NOT included in the flattened configuration, preventing perpetual Terraform diffs.
func TestFlattenClusterConfigurationUnspecifiedEnums(t *testing.T) {
	// Create a minimal cluster configuration with UNSPECIFIED enum values (or nil)
	clusterConfig := &qcCluster.ClusterConfiguration{
		NumberOfNodes: 1,
		PackageId:     "test-package-id",
	}

	flattened := flattenClusterConfiguration(clusterConfig, nil)

	require.Len(t, flattened, 1)
	configMap := flattened[0].(map[string]interface{})

	// Verify that UNSPECIFIED enum fields are NOT present in the output
	_, hasServiceType := configMap[serviceTypeFieldName]
	_, hasGpuType := configMap[dbConfigGpuTypeFieldName]
	_, hasRestartPolicy := configMap[dbConfigRestartPolicyFieldName]
	_, hasRebalanceStrategy := configMap[dbConfigRebalanceStrategyFieldName]

	assert.False(t, hasServiceType, "service_type should not be present when UNSPECIFIED")
	assert.False(t, hasGpuType, "gpu_type should not be present when UNSPECIFIED")
	assert.False(t, hasRestartPolicy, "restart_policy should not be present when UNSPECIFIED")
	assert.False(t, hasRebalanceStrategy, "rebalance_strategy should not be present when UNSPECIFIED")

	_, hasClusterStorageConfig := configMap[clusterStorageConfigurationFieldName]
	assert.False(t, hasClusterStorageConfig, "cluster_storage_configuration should not be present when UNSPECIFIED")

	assert.Equal(t, 1, configMap[numberOfNodesFieldName])
}

// TestExpandKeyVal_EmptyInputReturnsNil verifies that expandKeyVal returns nil for empty input.
func TestExpandKeyVal_EmptyInputReturnsNil(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		result := expandKeyVal([]interface{}{})
		assert.Nil(t, result)
	})

	t.Run("nil slice", func(t *testing.T) {
		result := expandKeyVal(nil)
		assert.Nil(t, result)
	})
}

// TestExpandTolerations_EmptyInputReturnsNil verifies that expandTolerations returns nil for empty input.
func TestExpandTolerations_EmptyInputReturnsNil(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		result := expandTolerations([]interface{}{})
		assert.Nil(t, result)
	})

	t.Run("nil slice", func(t *testing.T) {
		result := expandTolerations(nil)
		assert.Nil(t, result)
	})
}

// TestExpandTopologySpreadConstraints_EmptyInputReturnsNil verifies that expandTopologySpreadConstraints returns nil for empty input.
func TestExpandTopologySpreadConstraints_EmptyInputReturnsNil(t *testing.T) {
	t.Run("empty slice", func(t *testing.T) {
		result := expandTopologySpreadConstraints([]interface{}{})
		assert.Nil(t, result)
	})

	t.Run("nil slice", func(t *testing.T) {
		result := expandTopologySpreadConstraints(nil)
		assert.Nil(t, result)
	})
}

// TestExpandClusterConfiguration_EmptyAllowedIpSourceRangesReturnsNil verifies that
// AllowedIpSourceRanges is nil when not provided or empty in Terraform config.
func TestExpandClusterConfiguration_EmptyAllowedIpSourceRangesReturnsNil(t *testing.T) {
	t.Run("allowed_ip_source_ranges omitted", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, accountsClusterSchema(false), map[string]interface{}{
			configurationFieldName: []interface{}{
				map[string]interface{}{
					numberOfNodesFieldName: 1,
					packageIDFieldName:     "some-package-id",
				},
			},
		})

		// Extract the configuration block for expandClusterConfiguration
		configBlock := d.Get(configurationFieldName).([]interface{})
		config, _ := expandClusterConfiguration(configBlock)

		require.NotNil(t, config)
		assert.Nil(t, config.AllowedIpSourceRanges)
	})

	t.Run("allowed_ip_source_ranges explicitly empty", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, accountsClusterSchema(false), map[string]interface{}{
			configurationFieldName: []interface{}{
				map[string]interface{}{
					numberOfNodesFieldName:         1,
					packageIDFieldName:             "some-package-id",
					allowedIpSourceRangesFieldName: []interface{}{}, // Explicitly empty list
				},
			},
		})

		// Extract the configuration block for expandClusterConfiguration
		configBlock := d.Get(configurationFieldName).([]interface{})
		config, _ := expandClusterConfiguration(configBlock)

		require.NotNil(t, config)
		assert.Nil(t, config.AllowedIpSourceRanges)
	})
}

// TestClusterConfigSchemaEnumFieldsAreOptionalAndComputed verifies that enum and
// backend-managed fields are both Optional and Computed in resource mode.
// This prevents perpetual diffs when the backend returns UNSPECIFIED or default values
// for fields the user didn't set.
func TestClusterConfigSchemaEnumFieldsAreOptionalAndComputed(t *testing.T) {
	schemaMap := accountsClusterConfigurationSchema(false) // resource mode
	fieldsToCheck := []string{
		serviceTypeFieldName,
		dbConfigGpuTypeFieldName,
		dbConfigRestartPolicyFieldName,
		dbConfigRebalanceStrategyFieldName,
		clusterStorageConfigurationFieldName,
		dbConfigReservedCpuPercentageFieldName,
		dbConfigReservedMemoryPercentageFieldName,
		databaseConfigurationFieldName,
	}
	for _, field := range fieldsToCheck {
		t.Run(field, func(t *testing.T) {
			s, ok := schemaMap[field]
			require.True(t, ok, "field %s should exist in schema", field)
			assert.True(t, s.Optional, "field %s should be Optional", field)
			assert.True(t, s.Computed, "field %s should be Computed to prevent perpetual diffs with UNSPECIFIED backend values", field)
		})
	}
}

// TestClusterConfigSchemaDataSourceFieldsAreComputed verifies that the same fields
// are Computed (and not Optional) in data source mode.
func TestDatabaseConfigSchemaFieldsAreOptionalAndComputed(t *testing.T) {
	schemaMap := databaseConfigurationSchema(false)
	fieldsToCheck := []string{
		dbConfigCollectionFieldName,
		dbConfigStorageFieldName,
		dbConfigServiceFieldName,
		dbConfigLogLevelFieldName,
		dbConfigTlsFieldName,
	}
	for _, field := range fieldsToCheck {
		t.Run(field, func(t *testing.T) {
			s, ok := schemaMap[field]
			require.True(t, ok, "field %s should exist in schema", field)
			assert.True(t, s.Optional, "field %s should be Optional", field)
			assert.True(t, s.Computed, "field %s should be Computed to prevent perpetual diffs", field)
		})
	}
}

// TestOptionalFieldsMustBeComputed ensures that every Optional field in the
// cluster configuration schema (resource mode) also has Computed: true.
// If the backend can return a value for a field the user didn't set,
// Terraform needs Computed: true to avoid perpetual diffs.
// This test prevents regressions when new Optional fields are added.
func TestOptionalFieldsMustBeComputed(t *testing.T) {
	schemas := map[string]map[string]*schema.Schema{
		"clusterConfiguration":        accountsClusterConfigurationSchema(false),
		"clusterStorageConfiguration": clusterStorageConfigurationSchema(false),
		"databaseConfiguration":       databaseConfigurationSchema(false),
		"databaseCollection":          databaseConfigurationCollectionSchema(false),
		"databaseStorage":             databaseConfigurationStorageSchema(false),
		"databaseService":             databaseConfigurationServiceSchema(false),
		"databaseTls":                 databaseConfigurationTlsSchema(false),
		"databaseAuditLogging":        databaseConfigurationAuditLoggingSchema(false),
	}

	for schemaName, schemaMap := range schemas {
		for fieldName, fieldSchema := range schemaMap {
			if fieldSchema.Optional && fieldSchema.Deprecated == "" {
				t.Run(schemaName+"/"+fieldName, func(t *testing.T) {
					assert.True(t, fieldSchema.Computed,
						"field %s in %s is Optional but not Computed — the backend may populate this field, "+
							"causing perpetual Terraform diffs. Change Computed to true.",
						fieldName, schemaName)
				})
			}
		}
	}
}

func TestClusterConfigSchemaDataSourceFieldsAreComputed(t *testing.T) {
	schemaMap := accountsClusterConfigurationSchema(true) // data source mode
	fieldsToCheck := []string{
		serviceTypeFieldName,
		dbConfigGpuTypeFieldName,
		dbConfigRestartPolicyFieldName,
		dbConfigRebalanceStrategyFieldName,
		clusterStorageConfigurationFieldName,
		dbConfigReservedCpuPercentageFieldName,
		dbConfigReservedMemoryPercentageFieldName,
	}
	for _, field := range fieldsToCheck {
		t.Run(field, func(t *testing.T) {
			s, ok := schemaMap[field]
			require.True(t, ok, "field %s should exist in schema", field)
			assert.False(t, s.Optional, "field %s should not be Optional in data source mode", field)
			assert.True(t, s.Computed, "field %s should be Computed in data source mode", field)
		})
	}
}

// TestFlattenClusterConfigurationExplicitUnspecifiedPointers verifies that even when
// the backend explicitly sends UNSPECIFIED enum pointers (not just nil), they are
// excluded from the flattened output.
func TestFlattenClusterConfigurationExplicitUnspecifiedPointers(t *testing.T) {
	clusterConfig := &qcCluster.ClusterConfiguration{
		NumberOfNodes:     1,
		PackageId:         "test-package-id",
		ServiceType:       newPointer(qcCluster.ClusterServiceType_CLUSTER_SERVICE_TYPE_UNSPECIFIED),
		GpuType:           newPointer(qcCluster.ClusterConfigurationGpuType_CLUSTER_CONFIGURATION_GPU_TYPE_UNSPECIFIED),
		RestartPolicy:     newPointer(qcCluster.ClusterConfigurationRestartPolicy_CLUSTER_CONFIGURATION_RESTART_POLICY_UNSPECIFIED),
		RebalanceStrategy: newPointer(qcCluster.ClusterConfigurationRebalanceStrategy_CLUSTER_CONFIGURATION_REBALANCE_STRATEGY_UNSPECIFIED),
		ClusterStorageConfiguration: &qcCluster.ClusterStorageConfiguration{
			StorageTierType: commonv1.StorageTierType_STORAGE_TIER_TYPE_UNSPECIFIED,
		},
	}

	flattened := flattenClusterConfiguration(clusterConfig, nil)

	require.Len(t, flattened, 1)
	configMap := flattened[0].(map[string]interface{})

	_, hasServiceType := configMap[serviceTypeFieldName]
	_, hasGpuType := configMap[dbConfigGpuTypeFieldName]
	_, hasRestartPolicy := configMap[dbConfigRestartPolicyFieldName]
	_, hasRebalanceStrategy := configMap[dbConfigRebalanceStrategyFieldName]

	assert.False(t, hasServiceType, "service_type should not be present when explicitly UNSPECIFIED")
	assert.False(t, hasGpuType, "gpu_type should not be present when explicitly UNSPECIFIED")
	assert.False(t, hasRestartPolicy, "restart_policy should not be present when explicitly UNSPECIFIED")
	assert.False(t, hasRebalanceStrategy, "rebalance_strategy should not be present when explicitly UNSPECIFIED")

	_, hasClusterStorageConfig := configMap[clusterStorageConfigurationFieldName]
	assert.False(t, hasClusterStorageConfig, "cluster_storage_configuration should not be present when explicitly UNSPECIFIED")
}

// TestFlattenClusterConfigurationAllEnumValues verifies that every non-UNSPECIFIED
// enum value for each field is correctly included in the flattened output.
func TestFlattenClusterConfigurationAllEnumValues(t *testing.T) {
	t.Run("all ServiceType values", func(t *testing.T) {
		for val, name := range qcCluster.ClusterServiceType_name {
			if val == 0 {
				continue // skip UNSPECIFIED
			}
			t.Run(name, func(t *testing.T) {
				cfg := &qcCluster.ClusterConfiguration{
					ServiceType: newPointer(qcCluster.ClusterServiceType(val)),
				}
				flattened := flattenClusterConfiguration(cfg, nil)
				configMap := flattened[0].(map[string]interface{})
				assert.Equal(t, name, configMap[serviceTypeFieldName])
			})
		}
	})

	t.Run("all GpuType values", func(t *testing.T) {
		for val, name := range qcCluster.ClusterConfigurationGpuType_name {
			if val == 0 {
				continue
			}
			t.Run(name, func(t *testing.T) {
				cfg := &qcCluster.ClusterConfiguration{
					GpuType: newPointer(qcCluster.ClusterConfigurationGpuType(val)),
				}
				flattened := flattenClusterConfiguration(cfg, nil)
				configMap := flattened[0].(map[string]interface{})
				assert.Equal(t, name, configMap[dbConfigGpuTypeFieldName])
			})
		}
	})

	t.Run("all RestartPolicy values", func(t *testing.T) {
		for val, name := range qcCluster.ClusterConfigurationRestartPolicy_name {
			if val == 0 {
				continue
			}
			t.Run(name, func(t *testing.T) {
				cfg := &qcCluster.ClusterConfiguration{
					RestartPolicy: newPointer(qcCluster.ClusterConfigurationRestartPolicy(val)),
				}
				flattened := flattenClusterConfiguration(cfg, nil)
				configMap := flattened[0].(map[string]interface{})
				assert.Equal(t, name, configMap[dbConfigRestartPolicyFieldName])
			})
		}
	})

	t.Run("all RebalanceStrategy values", func(t *testing.T) {
		for val, name := range qcCluster.ClusterConfigurationRebalanceStrategy_name {
			if val == 0 {
				continue
			}
			t.Run(name, func(t *testing.T) {
				cfg := &qcCluster.ClusterConfiguration{
					RebalanceStrategy: newPointer(qcCluster.ClusterConfigurationRebalanceStrategy(val)),
				}
				flattened := flattenClusterConfiguration(cfg, nil)
				configMap := flattened[0].(map[string]interface{})
				assert.Equal(t, name, configMap[dbConfigRebalanceStrategyFieldName])
			})
		}
	})

	t.Run("all StorageTierType values", func(t *testing.T) {
		for val, name := range commonv1.StorageTierType_name {
			if val == 0 {
				continue // skip UNSPECIFIED
			}
			t.Run(name, func(t *testing.T) {
				cfg := &qcCluster.ClusterConfiguration{
					ClusterStorageConfiguration: &qcCluster.ClusterStorageConfiguration{
						StorageTierType: commonv1.StorageTierType(val),
					},
				}
				flattened := flattenClusterConfiguration(cfg, nil)
				configMap := flattened[0].(map[string]interface{})
				clusterStorageConfigList, ok := configMap[clusterStorageConfigurationFieldName].([]interface{})
				require.True(t, ok, "cluster_storage_configuration should be present")
				require.Len(t, clusterStorageConfigList, 1)
				clusterStorageConfig := clusterStorageConfigList[0].(map[string]interface{})
				assert.Equal(t, name, clusterStorageConfig[clusterStorageTierTypeFieldName])
			})
		}
	})
}

// TestExpandClusterConfigurationUnspecifiedEnumStrings verifies that if UNSPECIFIED
// string values somehow end up in the Terraform state, expand correctly ignores them
// and does not send them to the API.
func TestExpandClusterConfigurationUnspecifiedEnumStrings(t *testing.T) {
	configBlock := []interface{}{
		map[string]interface{}{
			numberOfNodesFieldName:             1,
			serviceTypeFieldName:               "CLUSTER_SERVICE_TYPE_UNSPECIFIED",
			dbConfigGpuTypeFieldName:           "CLUSTER_CONFIGURATION_GPU_TYPE_UNSPECIFIED",
			dbConfigRestartPolicyFieldName:     "CLUSTER_CONFIGURATION_RESTART_POLICY_UNSPECIFIED",
			dbConfigRebalanceStrategyFieldName: "CLUSTER_CONFIGURATION_REBALANCE_STRATEGY_UNSPECIFIED",
			clusterStorageConfigurationFieldName: []interface{}{
				map[string]interface{}{
					clusterStorageTierTypeFieldName: "STORAGE_TIER_TYPE_UNSPECIFIED",
				},
			},
		},
	}
	config, _ := expandClusterConfiguration(configBlock)

	require.NotNil(t, config)
	assert.Nil(t, config.ServiceType, "ServiceType should be nil when UNSPECIFIED string is provided")
	assert.Nil(t, config.GpuType, "GpuType should be nil when UNSPECIFIED string is provided")
	assert.Nil(t, config.RestartPolicy, "RestartPolicy should be nil when UNSPECIFIED string is provided")
	assert.Nil(t, config.RebalanceStrategy, "RebalanceStrategy should be nil when UNSPECIFIED string is provided")
	assert.Nil(t, config.ClusterStorageConfiguration, "ClusterStorageConfiguration should be nil when UNSPECIFIED string is provided")
}

// TestExpandClusterConfigurationValidEnumStrings verifies that expand correctly
// converts valid enum string values to their protobuf counterparts.
func TestExpandClusterConfigurationValidEnumStrings(t *testing.T) {
	configBlock := []interface{}{
		map[string]interface{}{
			numberOfNodesFieldName:             1,
			serviceTypeFieldName:               "CLUSTER_SERVICE_TYPE_LOAD_BALANCER",
			dbConfigGpuTypeFieldName:           "CLUSTER_CONFIGURATION_GPU_TYPE_NVIDIA",
			dbConfigRestartPolicyFieldName:     "CLUSTER_CONFIGURATION_RESTART_POLICY_ROLLING",
			dbConfigRebalanceStrategyFieldName: "CLUSTER_CONFIGURATION_REBALANCE_STRATEGY_BY_COUNT",
			clusterStorageConfigurationFieldName: []interface{}{
				map[string]interface{}{
					clusterStorageTierTypeFieldName: "STORAGE_TIER_TYPE_PERFORMANCE",
				},
			},
		},
	}
	config, _ := expandClusterConfiguration(configBlock)

	require.NotNil(t, config)
	assert.Equal(t, qcCluster.ClusterServiceType_CLUSTER_SERVICE_TYPE_LOAD_BALANCER, config.GetServiceType())
	assert.Equal(t, qcCluster.ClusterConfigurationGpuType_CLUSTER_CONFIGURATION_GPU_TYPE_NVIDIA, config.GetGpuType())
	assert.Equal(t, qcCluster.ClusterConfigurationRestartPolicy_CLUSTER_CONFIGURATION_RESTART_POLICY_ROLLING, config.GetRestartPolicy())
	assert.Equal(t, qcCluster.ClusterConfigurationRebalanceStrategy_CLUSTER_CONFIGURATION_REBALANCE_STRATEGY_BY_COUNT, config.GetRebalanceStrategy())
	assert.Equal(t, commonv1.StorageTierType_STORAGE_TIER_TYPE_PERFORMANCE, config.GetClusterStorageConfiguration().GetStorageTierType())
}

// TestExpandClusterConfigurationEmptyEnumStrings verifies that empty string values
// for enum fields (the Terraform zero value) do not get sent to the API.
func TestExpandClusterConfigurationEmptyEnumStrings(t *testing.T) {
	configBlock := []interface{}{
		map[string]interface{}{
			numberOfNodesFieldName:             1,
			serviceTypeFieldName:               "",
			dbConfigGpuTypeFieldName:           "",
			dbConfigRestartPolicyFieldName:     "",
			dbConfigRebalanceStrategyFieldName: "",
			clusterStorageConfigurationFieldName: []interface{}{
				map[string]interface{}{
					clusterStorageTierTypeFieldName: "",
				},
			},
		},
	}
	config, _ := expandClusterConfiguration(configBlock)

	require.NotNil(t, config)
	assert.Nil(t, config.ServiceType, "ServiceType should be nil for empty string")
	assert.Nil(t, config.GpuType, "GpuType should be nil for empty string")
	assert.Nil(t, config.RestartPolicy, "RestartPolicy should be nil for empty string")
	assert.Nil(t, config.RebalanceStrategy, "RebalanceStrategy should be nil for empty string")
	assert.Nil(t, config.ClusterStorageConfiguration, "ClusterStorageConfiguration should be nil for empty string")
}

// TestFlattenClusterConfigurationNilReservedPercentages verifies that nil pointer
// values for reserved percentage fields are not included in the flattened output.
func TestFlattenClusterConfigurationNilReservedPercentages(t *testing.T) {
	clusterConfig := &qcCluster.ClusterConfiguration{
		NumberOfNodes:            1,
		PackageId:                "test-package-id",
		ReservedCpuPercentage:    nil,
		ReservedMemoryPercentage: nil,
	}

	flattened := flattenClusterConfiguration(clusterConfig, nil)

	require.Len(t, flattened, 1)
	configMap := flattened[0].(map[string]interface{})

	_, hasCpu := configMap[dbConfigReservedCpuPercentageFieldName]
	_, hasMem := configMap[dbConfigReservedMemoryPercentageFieldName]

	assert.False(t, hasCpu, "reserved_cpu_percentage should not be present when nil")
	assert.False(t, hasMem, "reserved_memory_percentage should not be present when nil")
}

// TestFlattenClusterConfigurationWithReservedPercentages verifies that actual
// reserved percentage values are correctly included in the flattened output.
func TestFlattenClusterConfigurationWithReservedPercentages(t *testing.T) {
	clusterConfig := &qcCluster.ClusterConfiguration{
		NumberOfNodes:            1,
		PackageId:                "test-package-id",
		ReservedCpuPercentage:    newPointer(uint32(10)),
		ReservedMemoryPercentage: newPointer(uint32(20)),
	}

	flattened := flattenClusterConfiguration(clusterConfig, nil)

	require.Len(t, flattened, 1)
	configMap := flattened[0].(map[string]interface{})

	assert.Equal(t, 10, configMap[dbConfigReservedCpuPercentageFieldName])
	assert.Equal(t, 20, configMap[dbConfigReservedMemoryPercentageFieldName])
}

// TestFlattenExpandRoundTripEnumFields verifies that flatten → expand → flatten
// produces consistent results for enum fields, ensuring no data is lost or corrupted.
func TestFlattenExpandRoundTripEnumFields(t *testing.T) {
	original := &qcCluster.ClusterConfiguration{
		NumberOfNodes:     3,
		PackageId:         "test-package-id",
		ServiceType:       newPointer(qcCluster.ClusterServiceType_CLUSTER_SERVICE_TYPE_LOAD_BALANCER),
		GpuType:           newPointer(qcCluster.ClusterConfigurationGpuType_CLUSTER_CONFIGURATION_GPU_TYPE_NVIDIA),
		RestartPolicy:     newPointer(qcCluster.ClusterConfigurationRestartPolicy_CLUSTER_CONFIGURATION_RESTART_POLICY_ROLLING),
		RebalanceStrategy: newPointer(qcCluster.ClusterConfigurationRebalanceStrategy_CLUSTER_CONFIGURATION_REBALANCE_STRATEGY_BY_COUNT),
		ClusterStorageConfiguration: &qcCluster.ClusterStorageConfiguration{
			StorageTierType: commonv1.StorageTierType_STORAGE_TIER_TYPE_BALANCED,
		},
	}

	// First flatten
	flattened1 := flattenClusterConfiguration(original, nil)
	require.Len(t, flattened1, 1)

	// Expand back
	expanded, _ := expandClusterConfiguration(flattened1)
	require.NotNil(t, expanded)

	// Second flatten
	flattened2 := flattenClusterConfiguration(expanded, nil)
	require.Len(t, flattened2, 1)

	// Compare enum fields between first and second flatten
	map1 := flattened1[0].(map[string]interface{})
	map2 := flattened2[0].(map[string]interface{})

	assert.Equal(t, map1[serviceTypeFieldName], map2[serviceTypeFieldName], "service_type should be consistent across round-trip")
	assert.Equal(t, map1[dbConfigGpuTypeFieldName], map2[dbConfigGpuTypeFieldName], "gpu_type should be consistent across round-trip")
	assert.Equal(t, map1[dbConfigRestartPolicyFieldName], map2[dbConfigRestartPolicyFieldName], "restart_policy should be consistent across round-trip")
	assert.Equal(t, map1[dbConfigRebalanceStrategyFieldName], map2[dbConfigRebalanceStrategyFieldName], "rebalance_strategy should be consistent across round-trip")

	// Check cluster_storage_configuration nested structure
	clusterStorageConfig1, ok1 := map1[clusterStorageConfigurationFieldName].([]interface{})
	clusterStorageConfig2, ok2 := map2[clusterStorageConfigurationFieldName].([]interface{})
	require.True(t, ok1, "cluster_storage_configuration should be present in first flatten")
	require.True(t, ok2, "cluster_storage_configuration should be present in second flatten")
	require.Len(t, clusterStorageConfig1, 1)
	require.Len(t, clusterStorageConfig2, 1)
	clusterStorageMap1 := clusterStorageConfig1[0].(map[string]interface{})
	clusterStorageMap2 := clusterStorageConfig2[0].(map[string]interface{})
	assert.Equal(t, clusterStorageMap1[clusterStorageTierTypeFieldName], clusterStorageMap2[clusterStorageTierTypeFieldName], "storage_tier_type should be consistent across round-trip")
}

// TestFlattenExpandRoundTripUnspecifiedEnumFields verifies that the round-trip
// for UNSPECIFIED values is also consistent — they should remain absent.
func TestFlattenExpandRoundTripUnspecifiedEnumFields(t *testing.T) {
	original := &qcCluster.ClusterConfiguration{
		NumberOfNodes: 1,
		PackageId:     "test-package-id",
		// All enum fields default to UNSPECIFIED (nil pointers)
	}

	// First flatten — UNSPECIFIED values should be absent
	flattened1 := flattenClusterConfiguration(original, nil)
	require.Len(t, flattened1, 1)

	// Expand back — should produce nil enum pointers
	expanded, _ := expandClusterConfiguration(flattened1)
	require.NotNil(t, expanded)
	assert.Nil(t, expanded.ServiceType)
	assert.Nil(t, expanded.GpuType)
	assert.Nil(t, expanded.RestartPolicy)
	assert.Nil(t, expanded.RebalanceStrategy)
	assert.Nil(t, expanded.ClusterStorageConfiguration, "ClusterStorageConfiguration should be nil when no storage_tier_type is provided")

	// Second flatten — should still be absent
	flattened2 := flattenClusterConfiguration(expanded, nil)
	require.Len(t, flattened2, 1)
	map2 := flattened2[0].(map[string]interface{})

	_, hasServiceType := map2[serviceTypeFieldName]
	_, hasGpuType := map2[dbConfigGpuTypeFieldName]
	_, hasRestartPolicy := map2[dbConfigRestartPolicyFieldName]
	_, hasRebalanceStrategy := map2[dbConfigRebalanceStrategyFieldName]

	assert.False(t, hasServiceType, "service_type should remain absent across round-trip")
	assert.False(t, hasGpuType, "gpu_type should remain absent across round-trip")
	assert.False(t, hasRestartPolicy, "restart_policy should remain absent across round-trip")
	assert.False(t, hasRebalanceStrategy, "rebalance_strategy should remain absent across round-trip")

	_, hasClusterStorageConfig := map2[clusterStorageConfigurationFieldName]
	assert.False(t, hasClusterStorageConfig, "cluster_storage_configuration should remain absent across round-trip")
}

// TestFlattenClusterConfigurationSpecifiedEnums verifies that specified enum values
// ARE included in the flattened configuration.
func TestFlattenClusterConfigurationSpecifiedEnums(t *testing.T) {
	clusterConfig := &qcCluster.ClusterConfiguration{
		NumberOfNodes:     1,
		PackageId:         "test-package-id",
		ServiceType:       newPointer(qcCluster.ClusterServiceType_CLUSTER_SERVICE_TYPE_LOAD_BALANCER),
		GpuType:           newPointer(qcCluster.ClusterConfigurationGpuType_CLUSTER_CONFIGURATION_GPU_TYPE_NVIDIA),
		RestartPolicy:     newPointer(qcCluster.ClusterConfigurationRestartPolicy_CLUSTER_CONFIGURATION_RESTART_POLICY_ROLLING),
		RebalanceStrategy: newPointer(qcCluster.ClusterConfigurationRebalanceStrategy_CLUSTER_CONFIGURATION_REBALANCE_STRATEGY_BY_COUNT),
		ClusterStorageConfiguration: &qcCluster.ClusterStorageConfiguration{
			StorageTierType: commonv1.StorageTierType_STORAGE_TIER_TYPE_BALANCED,
		},
	}

	flattened := flattenClusterConfiguration(clusterConfig, nil)

	require.Len(t, flattened, 1)
	configMap := flattened[0].(map[string]interface{})

	// Verify that specified enum fields ARE present with correct values
	assert.Equal(t, "CLUSTER_SERVICE_TYPE_LOAD_BALANCER", configMap[serviceTypeFieldName])
	assert.Equal(t, "CLUSTER_CONFIGURATION_GPU_TYPE_NVIDIA", configMap[dbConfigGpuTypeFieldName])
	assert.Equal(t, "CLUSTER_CONFIGURATION_RESTART_POLICY_ROLLING", configMap[dbConfigRestartPolicyFieldName])
	assert.Equal(t, "CLUSTER_CONFIGURATION_REBALANCE_STRATEGY_BY_COUNT", configMap[dbConfigRebalanceStrategyFieldName])

	// Verify cluster_storage_configuration nested structure
	clusterStorageConfigList, ok := configMap[clusterStorageConfigurationFieldName].([]interface{})
	require.True(t, ok, "cluster_storage_configuration should be present")
	require.Len(t, clusterStorageConfigList, 1)
	clusterStorageConfig := clusterStorageConfigList[0].(map[string]interface{})
	assert.Equal(t, "STORAGE_TIER_TYPE_BALANCED", clusterStorageConfig[clusterStorageTierTypeFieldName])
}
