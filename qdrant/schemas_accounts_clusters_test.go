package qdrant

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	qc "github.com/qdrant/terraform-provider-qdrant-cloud/v1/internal/client"
)

func TestResourceClusterFlatten(t *testing.T) {
	configs := []qc.ResourceConfigurationSchema{
		{
			Amount:       8,
			ResourceUnit: "m",
			ResourceType: "cpu",
		},
	}
	cluster := &qc.ClusterSchema{
		Id:                  newPointer(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		CreatedAt:           newPointer(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
		AccountId:           uuid.MustParse("00000000-1000-0000-0000-000000000001"),
		Name:                "testName",
		CloudProvider:       newPointer(qc.ClusterSchemaCloudProviderAzure),
		CloudRegion:         newPointer(qc.ClusterSchemaCloudRegionUksouth),
		PrivateRegionId:     newPointer(uuid.MustParse("00000003-0000-0000-0000-000000000001")),
		MarkedForDeletionAt: newPointer(time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)),
		Url:                 newPointer("http://example.com"),
		Configuration: qc.ClusterConfigurationSchema{
			Version:       newPointer("v1.0"),
			NumberOfNodes: 5,
			NodeConfiguration: qc.NodeConfigurationSchema{
				PackageId:              uuid.MustParse("00000009-1000-0000-0000-000000000001"),
				ResourceConfigurations: &configs,
			},
		},
	}

	flattened := flattenCluster(cluster)

	expected := map[string]interface{}{
		clusterIdentifierFieldName:          cluster.Id.String(),
		clusterCreatedAtFieldName:           formatTime(cluster.CreatedAt),
		clusterAccountIDFieldName:           cluster.AccountId.String(),
		clusterNameFieldName:                cluster.Name,
		clusterCloudProviderFieldName:       string(derefPointer(cluster.CloudProvider)),
		clusterCloudRegionFieldName:         string(derefPointer(cluster.CloudRegion)),
		clusterPrivateRegionIDFieldName:     cluster.PrivateRegionId.String(),
		clusterMarkedForDeletionAtFieldName: formatTime(cluster.MarkedForDeletionAt),
		clusterURLFieldName:                 derefPointer(cluster.Url),
		configurationFieldName: []interface{}{
			map[string]interface{}{
				clusterVersionFieldName: derefPointer(cluster.Configuration.Version),
				numberOfNodesFieldName:  cluster.Configuration.NumberOfNodes,
				nodeConfigurationFieldName: []interface{}{
					map[string]interface{}{
						packageIDFieldName: cluster.Configuration.NodeConfiguration.PackageId.String(),
						resourceConfigurationsFieldName: []interface{}{
							map[string]interface{}{
								fieldAmount:       configs[0].Amount,
								fieldResourceUnit: string(configs[0].ResourceUnit),
								fieldResourceType: string(configs[0].ResourceType),
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
	configs := []qc.ResourceConfigurationSchema{
		{
			Amount:       10,
			ResourceUnit: "m",
			ResourceType: "cpu",
		},
	}
	expected := qc.ClusterSchema{
		Id:                  newPointer(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
		CreatedAt:           newPointer(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
		AccountId:           uuid.MustParse("00000000-1000-0000-0000-000000000001"),
		Name:                "testName",
		CloudProvider:       newPointer(qc.ClusterSchemaCloudProviderAzure),
		CloudRegion:         newPointer(qc.ClusterSchemaCloudRegionUksouth),
		PrivateRegionId:     newPointer(uuid.MustParse("00000003-0000-0000-0000-000000000001")),
		MarkedForDeletionAt: newPointer(time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)),
		Url:                 newPointer("http://example.com"),
		Configuration: qc.ClusterConfigurationSchema{
			Version:       newPointer("v1.0"),
			NumberOfNodes: 5,
			NodeConfiguration: qc.NodeConfigurationSchema{
				PackageId:              uuid.MustParse("00000009-1000-0000-0000-000000000001"),
				ResourceConfigurations: &configs,
			},
		},
	}

	d := schema.TestResourceDataRaw(t, accountsClusterSchema(false), map[string]interface{}{
		clusterIdentifierFieldName:          expected.Id.String(),
		clusterCreatedAtFieldName:           formatTime(expected.CreatedAt),
		clusterAccountIDFieldName:           expected.AccountId.String(),
		clusterNameFieldName:                expected.Name,
		clusterCloudProviderFieldName:       string(derefPointer(expected.CloudProvider)),
		clusterCloudRegionFieldName:         string(derefPointer(expected.CloudRegion)),
		clusterPrivateRegionIDFieldName:     expected.PrivateRegionId.String(),
		clusterMarkedForDeletionAtFieldName: formatTime(expected.MarkedForDeletionAt),
		clusterURLFieldName:                 derefPointer(expected.Url),
		configurationFieldName: []interface{}{
			map[string]interface{}{
				clusterVersionFieldName: derefPointer(expected.Configuration.Version),
				numberOfNodesFieldName:  expected.Configuration.NumberOfNodes,
				nodeConfigurationFieldName: []interface{}{
					map[string]interface{}{
						packageIDFieldName: expected.Configuration.NodeConfiguration.PackageId.String(),
						resourceConfigurationsFieldName: []interface{}{
							map[string]interface{}{
								fieldAmount:       configs[0].Amount,
								fieldResourceUnit: string(configs[0].ResourceUnit),
								fieldResourceType: string(configs[0].ResourceType),
							},
						},
					},
				},
			},
		},
	})

	result, err := expandCluster(d, expected.AccountId.String())
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}
