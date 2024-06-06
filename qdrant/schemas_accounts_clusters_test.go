package qdrant

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	qc "terraform-provider-qdrant-cloud/v1/internal/client"
)

func TestResourceClusterFlatten(t *testing.T) {
	cluster := &qc.ClusterOut{
		Id:                     "testID",
		CreatedAt:              time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
		AccountId:              newString("accountID"),
		Name:                   "testName",
		CloudProvider:          qc.ClusterOutCloudProviderAzure,
		CloudRegion:            qc.ClusterOutCloudRegionUksouth,
		CloudRegionAz:          newString("1"),
		Version:                newString("v1.0"),
		CloudRegionSetup:       newString("Standard"),
		PrivateRegionId:        newString("privateRegionID"),
		CurrentConfigurationId: "currentConfigID",
		EncryptionKeyId:        newString("encryptionKeyID"),
		MarkedForDeletionAt:    newTime(time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC)),
		Url:                    "http://example.com",
		TotalExtraDisk:         newInt(100),
		Configuration: &qc.ClusterConfigurationOut{
			NumNodes:    5,
			NumNodesMax: 10,
			NodeConfiguration: qc.NodeConfiguration{
				PackageId: "testPackageID",
			},
		},
	}

	flattened := flattenCluster(cluster)

	expected := map[string]interface{}{
		clusterIdentifierFieldName:                  cluster.Id,
		clusterCreatedAtFieldName:                   formatTime(&cluster.CreatedAt),
		clusterAccountIDFieldName:                   derefString(cluster.AccountId),
		clusterNameFieldName:                        cluster.Name,
		clusterCloudProviderFieldName:               cluster.CloudProvider,
		clusterCloudRegionFieldName:                 cluster.CloudRegion,
		clusterCloudRegionAvailabilityZoneFieldName: derefString(cluster.CloudRegionAz),
		clusterVersionFieldName:                     derefString(cluster.Version),
		clusterCloudRegionSetupFieldName:            derefString(cluster.CloudRegionSetup),
		clusterPrivateRegionIDFieldName:             derefString(cluster.PrivateRegionId),
		clusterCurrentConfigurationIDFieldName:      cluster.CurrentConfigurationId,
		clusterEncryptionKeyIDFieldName:             derefString(cluster.EncryptionKeyId),
		clusterMarkedForDeletionAtFieldName:         formatTime(cluster.MarkedForDeletionAt),
		clusterURLFieldName:                         cluster.Url,
		clusterTotalExtraDiskFieldName:              derefInt(cluster.TotalExtraDisk),
		configurationFieldName: []interface{}{
			map[string]interface{}{
				numNodesFieldName:    cluster.Configuration.NumNodes,
				numNodesMaxFieldName: cluster.Configuration.NumNodesMax,
				nodeConfigurationFieldName: []interface{}{
					map[string]interface{}{
						packageIDFieldName: cluster.Configuration.NodeConfiguration.PackageId,
					},
				},
			},
		},
	}

	assert.Equal(t, expected, flattened)
}

func TestExpandClusterIn(t *testing.T) {
	expected := qc.ClusterIn{
		AccountId:        newString("accountID"),
		Name:             "testName",
		CloudProvider:    qc.ClusterInCloudProviderAzure,
		CloudRegion:      qc.ClusterInCloudRegionUksouth,
		CloudRegionAz:    newString("1"),
		Version:          newString("v1.0"),
		CloudRegionSetup: newString("Standard"),
		PrivateRegionId:  newString("privateRegionID"),
		Configuration: qc.ClusterConfigurationIn{
			NumNodes:    5,
			NumNodesMax: 10,
			NodeConfiguration: qc.NodeConfiguration{
				PackageId: "testPackageID",
			},
		},
		EncryptionConfig: &qc.EncryptionConfigIn{
			AwsEncryptionConfig: &qc.AWSEncryptionConfig{
				EncryptionKeyId: newString("encryptionKey"),
			},
		},
	}

	d := schema.TestResourceDataRaw(t, accountsClusterSchema(false), map[string]interface{}{
		clusterAccountIDFieldName:                   derefString(expected.AccountId),
		clusterNameFieldName:                        expected.Name,
		clusterCloudProviderFieldName:               string(expected.CloudProvider),
		clusterCloudRegionFieldName:                 string(expected.CloudRegion),
		clusterCloudRegionAvailabilityZoneFieldName: derefString(expected.CloudRegionAz),
		clusterVersionFieldName:                     derefString(expected.Version),
		clusterCloudRegionSetupFieldName:            derefString(expected.CloudRegionSetup),
		clusterPrivateRegionIDFieldName:             derefString(expected.PrivateRegionId),
		configurationFieldName: []interface{}{
			map[string]interface{}{
				numNodesFieldName:    expected.Configuration.NumNodes,
				numNodesMaxFieldName: expected.Configuration.NumNodesMax,
				nodeConfigurationFieldName: []interface{}{
					map[string]interface{}{
						packageIDFieldName: expected.Configuration.NodeConfiguration.PackageId,
					},
				},
			},
		},
		clusterEncryptionKeyIDFieldName: derefString(expected.EncryptionConfig.AwsEncryptionConfig.EncryptionKeyId),
	})

	result, err := expandClusterIn(d, derefString(expected.AccountId))
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
