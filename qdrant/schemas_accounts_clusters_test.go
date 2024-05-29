package qdrant

import (
	"testing"

	"github.com/stretchr/testify/assert"

	qc "terraform-provider-qdrant-cloud/v1/internal/client"
)

func TestResourceClusterFlatten(t *testing.T) {
	cluster := &qc.ClusterOut{
		Name:          "testName",
		CloudProvider: qc.ClusterOutCloudProviderAzure,
		CloudRegion:   qc.ClusterOutCloudRegionUksouth,
		Configuration: &qc.ClusterConfigurationOut{
			NumNodes:    5,
			NumNodesMax: 10,
		},
	}
	flattened := flattenCluster(cluster)

	expected := map[string]interface{}{
		"id":             "", // ClusterOut contains an ID
		"name":           cluster.Name,
		"cloud_provider": cluster.CloudProvider,
		"cloud_region":   cluster.CloudRegion,
		"configuration": []interface{}{
			map[string]interface{}{
				"id":            "", // ConfigurationOut contains an ID
				"num_nodes":     cluster.Configuration.NumNodes,
				"num_nodes_max": cluster.Configuration.NumNodesMax,
			},
		},
	}

	assert.Equal(t, expected, flattened)
}
