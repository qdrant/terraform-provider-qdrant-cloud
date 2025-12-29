package qdrant

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	k8v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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
				Service: &qcCluster.DatabaseConfigurationService{
					ApiKey:         &commonv1.SecretKeyRef{Name: "api-key-secret", Key: "api-key"},
					ReadOnlyApiKey: &commonv1.SecretKeyRef{Name: "ro-api-key-secret", Key: "ro-api-key"},
					JwtRbac:        newPointer(true),
				},
				LogLevel: newPointer(qcCluster.DatabaseConfigurationLogLevel_DATABASE_CONFIGURATION_LOG_LEVEL_DEBUG),
				Tls: &qcCluster.DatabaseConfigurationTls{
					Cert: &commonv1.SecretKeyRef{Name: "cert-secret", Key: "cert.pem"},
					Key:  &commonv1.SecretKeyRef{Name: "key-secret", Key: "key.pem"},
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
					Key:      "key1",
					Operator: newPointer(qcCluster.TolerationOperator_TOLERATION_OPERATOR_EQUAL),
					Value:    "value1",
					Effect:   newPointer(qcCluster.TolerationEffect_TOLERATION_EFFECT_NO_SCHEDULE),
				},
			},
			TopologySpreadConstraints: []*k8v1.TopologySpreadConstraint{
				{
					MaxSkew:           1,
					TopologyKey:       "topology.kubernetes.io/zone",
					WhenUnsatisfiable: k8v1.UnsatisfiableConstraintAction("DoNotSchedule"),
					LabelSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"key1": "value1"},
						MatchExpressions: []metav1.LabelSelectorRequirement{
							{
								Key:      "key2",
								Operator: metav1.LabelSelectorOperator("In"),
								Values:   []string{"val1", "val2"},
							},
						},
					},
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
			ScalabilityInfo: &qcCluster.ClusterScalabilityInfo{
				Status: qcCluster.ClusterScalabilityStatus_CLUSTER_SCALABILITY_STATUS_SCALABLE,
				Reason: newPointer("Can be scaled"),
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
						clusterScalabilityInfoStatusFieldName: cluster.GetState().GetScalabilityInfo().GetStatus().String(),
						clusterScalabilityInfoReasonFieldName: cluster.GetState().GetScalabilityInfo().GetReason(),
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
								dbConfigServiceJwtRbacFieldName:   cluster.GetConfiguration().GetDatabaseConfiguration().GetService().GetJwtRbac(),
								dbConfigServiceEnableTlsFieldName: false,
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
						topologySpreadConstraintWhenUnsatisfiableFieldName: "DoNotSchedule",
						topologySpreadConstraintLabelSelectorFieldName: []interface{}{
							map[string]interface{}{
								matchLabelsFieldName: []interface{}{
									map[string]interface{}{"key": "key1", "value": "value1"},
								},
								matchExpressionsFieldName: []interface{}{
									map[string]interface{}{
										"key":      "key2",
										"operator": "In",
										"values":   []interface{}{"val1", "val2"},
									},
								},
							},
						},
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
					ApiKey:  &commonv1.SecretKeyRef{Name: "api-key-secret-expand", Key: "api-key-expand"},
					JwtRbac: newPointer(false),
				},
				Tls: &qcCluster.DatabaseConfigurationTls{
					Cert: &commonv1.SecretKeyRef{Name: "cert-secret-expand", Key: "cert.pem-expand"},
				},
			},
			NodeSelector: []*commonv1.KeyValue{
				{Key: "key1", Value: "value1"},
			},
			Tolerations: []*qcCluster.Toleration{
				{
					Key:      "key1",
					Operator: newPointer(qcCluster.TolerationOperator_TOLERATION_OPERATOR_EQUAL),
					Value:    "value1",
					Effect:   newPointer(qcCluster.TolerationEffect_TOLERATION_EFFECT_NO_SCHEDULE),
				},
			},
			TopologySpreadConstraints: []*k8v1.TopologySpreadConstraint{
				{
					MaxSkew:           1,
					TopologyKey:       "topology.kubernetes.io/zone",
					WhenUnsatisfiable: k8v1.UnsatisfiableConstraintAction("DoNotSchedule"),
					LabelSelector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"key1": "value1"},
						MatchExpressions: []metav1.LabelSelectorRequirement{
							{
								Key:      "key2",
								Operator: metav1.LabelSelectorOperator("In"),
								Values:   []string{"val1", "val2"},
							},
						},
					},
				},
			},
			ServiceType: newPointer(qcCluster.ClusterServiceType_CLUSTER_SERVICE_TYPE_LOAD_BALANCER),
			GpuType:     newPointer(qcCluster.ClusterConfigurationGpuType_CLUSTER_CONFIGURATION_GPU_TYPE_NVIDIA),
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
						topologySpreadConstraintWhenUnsatisfiableFieldName: "DoNotSchedule",
						topologySpreadConstraintLabelSelectorFieldName: []interface{}{
							map[string]interface{}{
								matchLabelsFieldName: []interface{}{
									map[string]interface{}{"key": "key1", "value": "value1"},
								},
								matchExpressionsFieldName: []interface{}{
									map[string]interface{}{
										"key":      "key2",
										"operator": "In",
										"values":   []interface{}{"val1", "val2"},
									},
								},
							},
						},
					},
				},
				serviceTypeFieldName:     "CLUSTER_SERVICE_TYPE_LOAD_BALANCER",
				dbConfigGpuTypeFieldName: "CLUSTER_CONFIGURATION_GPU_TYPE_NVIDIA",
			},
		},
	})

	result, err := expandCluster(d, expected.GetAccountId())
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}
