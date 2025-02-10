package qdrant

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	qcCluster "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/v2"
)

func TestResourceClusterFlatten(t *testing.T) {
	cluster := &qcCluster.Cluster{
		Id:            "00000000-0000-0000-0000-000000000001",
		CreatedAt:     timestamppb.New(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
		AccountId:     "00000000-1000-0000-0000-000000000001",
		Name:          "testName",
		CloudProvider: "Azure",
		CloudRegion:   "Uksouth",
		// TODO: MarkedForDeletionAt: timestamppb.New(time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)),
		Configuration: &qcCluster.ClusterConfiguration{
			Version:       "v1.0",
			NumberOfNodes: 5,
			PackageId:     "00000009-1000-0000-0000-000000000001",
			AdditionalResources: &qcCluster.AdditionalResources{
				Disk: 8,
			},
		},
		State: &qcCluster.ClusterState{
			Endpoint: &qcCluster.ClusterEndpoint{
				Url: "http://example.com",
			},
		},
	}

	flattened := flattenCluster(cluster)

	expected := map[string]interface{}{
		clusterIdentifierFieldName:      cluster.GetId(),
		clusterCreatedAtFieldName:       formatTime(cluster.GetCreatedAt()),
		clusterAccountIDFieldName:       cluster.GetAccountId(),
		clusterNameFieldName:            cluster.GetName(),
		clusterCloudProviderFieldName:   cluster.GetCloudProvider(),
		clusterCloudRegionFieldName:     cluster.GetCloudRegion(),
		clusterPrivateRegionIDFieldName: "",
		//TODO: clusterMarkedForDeletionAtFieldName: formatTime(cluster.MarkedForDeletionAt),
		clusterURLFieldName: cluster.GetState().GetEndpoint().GetUrl(),
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
			},
		},
	}

	assert.Equal(t, expected, flattened)
}

func TestExpandCluster(t *testing.T) {
	expected := &qcCluster.Cluster{
		Id:            "00000000-0000-0000-0000-000000000001",
		CreatedAt:     timestamppb.New(time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)),
		AccountId:     "00000000-1000-0000-0000-000000000001",
		Name:          "testName",
		CloudProvider: "Azure",
		CloudRegion:   "Uksouth",
		// TODO: MarkedForDeletionAt: timestamppb.New(time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)),
		Configuration: &qcCluster.ClusterConfiguration{
			Version:       "v1.0",
			NumberOfNodes: 5,
			PackageId:     "00000009-1000-0000-0000-000000000001",
			AdditionalResources: &qcCluster.AdditionalResources{
				Disk: 10,
			},
		},
		State: &qcCluster.ClusterState{
			Endpoint: &qcCluster.ClusterEndpoint{
				Url: "http://example.com",
			},
		},
	}

	d := schema.TestResourceDataRaw(t, accountsClusterSchema(false), map[string]interface{}{
		clusterIdentifierFieldName:      expected.GetId(),
		clusterCreatedAtFieldName:       formatTime(expected.GetCreatedAt()),
		clusterAccountIDFieldName:       expected.GetAccountId(),
		clusterNameFieldName:            expected.GetName(),
		clusterCloudProviderFieldName:   expected.GetCloudProvider(),
		clusterCloudRegionFieldName:     expected.GetCloudRegion(),
		clusterPrivateRegionIDFieldName: "",
		// TODO: clusterMarkedForDeletionAtFieldName: formatTime(expected.MarkedForDeletionAt),
		clusterURLFieldName: expected.GetState().GetEndpoint().GetUrl(),
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
			},
		},
	})

	result, err := expandCluster(d, expected.GetAccountId())
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}
