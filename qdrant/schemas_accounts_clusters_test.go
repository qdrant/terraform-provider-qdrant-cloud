package qdrant

import (
	"testing"

	"github.com/stretchr/testify/assert"

	qc "terraform-provider-qdrant-cloud/v1/internal/client"
)

func TestResourceClusterFlatten(t *testing.T) {
	cluster := &qc.ClusterOut{
		AccountId:     newString("accountID"),
		Name:          "testName",
		CloudProvider: qc.ClusterOutCloudProviderAzure,
		CloudRegion:   qc.ClusterOutCloudRegionUksouth,
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
		"id":             "", // ClusterOut contains an ID
		"account_id":     cluster.AccountId,
		"name":           cluster.Name,
		"cloud_provider": cluster.CloudProvider,
		"cloud_region":   cluster.CloudRegion,
		"configuration": []interface{}{
			map[string]interface{}{
				//"id":            "", // ConfigurationOut contains an ID
				"num_nodes":     cluster.Configuration.NumNodes,
				"num_nodes_max": cluster.Configuration.NumNodesMax,
				"node_configuration": []interface{}{
					map[string]interface{}{
						"package_id": cluster.Configuration.NodeConfiguration.PackageId,
					},
				},
			},
		},
	}

	assert.Equal(t, expected, flattened)
}

func newString(s string) *string {
	return &s
}
