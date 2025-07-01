package qdrant

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	commonv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/common/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	qcCluster "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/v1"
)

func TestResourceClusterFlatten(t *testing.T) {
	cluster := &qcCluster.Cluster{
		Id:                    "00000000-0000-0000-0000-000000000001",
		CreatedAt:             timestamppb.New(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
		AccountId:             "00000000-1000-0000-0000-000000000001",
		Name:                  "testName",
		CloudProviderId:       "Azure",
		CloudProviderRegionId: "Uksouth",
		DeletedAt:             timestamppb.New(time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)),
		Configuration: &qcCluster.ClusterConfiguration{
			Version:       newPointer("v1.0"),
			NumberOfNodes: 5,
			PackageId:     "00000009-1000-0000-0000-000000000001",
			AdditionalResources: &qcCluster.AdditionalResources{
				Disk: 8,
			},
			DatabaseConfiguration: &qcCluster.DatabaseConfiguration{
				Collection: &qcCluster.DatabaseConfigurationCollection{
					ReplicationFactor:      newPointer(uint32(2)),
					WriteConsistencyFactor: 1,
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
			ScalabilityInfo: &qcCluster.ClusterScalabilityInfo{
				Status: qcCluster.ClusterScalabilityStatus_CLUSTER_SCALABILITY_STATUS_SCALABLE,
				Reason: newPointer("Can be scaled"),
			},
		},
	}

	flattened := flattenCluster(cluster)

	expected := map[string]interface{}{
		clusterIdentifierFieldName:          cluster.GetId(),
		clusterCreatedAtFieldName:           formatTime(cluster.GetCreatedAt()),
		clusterAccountIDFieldName:           cluster.GetAccountId(),
		clusterNameFieldName:                cluster.GetName(),
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
				clusterVersionFieldName: cluster.GetConfiguration().GetVersion(),
				numberOfNodesFieldName:  int(cluster.GetConfiguration().GetNumberOfNodes()),
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
								dbConfigServiceJwtRbacFieldName:   false,
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
			},
		},
	}

	assert.Equal(t, expected, flattened)
}

func TestExpandCluster(t *testing.T) {
	expected := &qcCluster.Cluster{
		Id:                    "00000000-0000-0000-0000-000000000001",
		CreatedAt:             timestamppb.New(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
		AccountId:             "00000000-1000-0000-0000-000000000001",
		Name:                  "testName",
		CloudProviderId:       "Azure",
		CloudProviderRegionId: "Uksouth",
		DeletedAt:             timestamppb.New(time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)),
		Configuration: &qcCluster.ClusterConfiguration{
			Version:       newPointer("v1.0"),
			NumberOfNodes: 5,
			PackageId:     "00000009-1000-0000-0000-000000000001",
			AdditionalResources: &qcCluster.AdditionalResources{
				Disk: 10,
			},
			DatabaseConfiguration: &qcCluster.DatabaseConfiguration{
				Collection: &qcCluster.DatabaseConfigurationCollection{
					ReplicationFactor: newPointer(uint32(3)),
				},
				Service: &qcCluster.DatabaseConfigurationService{
					ApiKey: &commonv1.SecretKeyRef{Name: "api-key-secret-expand", Key: "api-key-expand"},
				},
				Tls: &qcCluster.DatabaseConfigurationTls{
					Cert: &commonv1.SecretKeyRef{Name: "cert-secret-expand", Key: "cert.pem-expand"},
				},
			},
		},
		State: &qcCluster.ClusterState{
			Endpoint: &qcCluster.ClusterEndpoint{
				Url: "http://example.com",
			},
		},
	}

	d := schema.TestResourceDataRaw(t, accountsClusterSchema(false), map[string]interface{}{
		clusterIdentifierFieldName:          expected.GetId(),
		clusterCreatedAtFieldName:           formatTime(expected.GetCreatedAt()),
		clusterAccountIDFieldName:           expected.GetAccountId(),
		clusterNameFieldName:                expected.GetName(),
		clusterCloudProviderFieldName:       expected.GetCloudProviderId(),
		clusterCloudRegionFieldName:         expected.GetCloudProviderRegionId(),
		clusterPrivateRegionIDFieldName:     "",
		clusterMarkedForDeletionAtFieldName: formatTime(expected.GetDeletedAt()),
		clusterURLFieldName:                 expected.GetState().GetEndpoint().GetUrl(),
		configurationFieldName: []interface{}{
			map[string]interface{}{
				clusterVersionFieldName: expected.GetConfiguration().GetVersion(),
				numberOfNodesFieldName:  int(expected.GetConfiguration().GetNumberOfNodes()),
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
								dbConfigCollectionReplicationFactor: 3,
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
			},
		},
	})

	result, err := expandCluster(d, expected.GetAccountId())
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}
