package qdrant

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	qcCluster "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/v1"
	commonv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/common/v1"
)

func TestExpandClusterStorageConfigurationStorageClasses(t *testing.T) {
	config := expandClusterStorageConfiguration([]interface{}{
		map[string]interface{}{
			clusterStorageTierTypeFieldName:      "STORAGE_TIER_TYPE_BALANCED",
			clusterDatabaseStorageClassFieldName: "premium-perf1-stackit",
			clusterSnapshotStorageClassFieldName: "premium-perf1-stackit",
			clusterVolumeSnapshotClassFieldName:  "csi-snapclass",
		},
	})
	require.NotNil(t, config)
	assert.Equal(t, commonv1.StorageTierType_STORAGE_TIER_TYPE_BALANCED, config.GetStorageTierType())
	assert.Equal(t, "premium-perf1-stackit", config.GetDatabaseStorageClass())
	assert.Equal(t, "premium-perf1-stackit", config.GetSnapshotStorageClass())
	assert.Equal(t, "csi-snapclass", config.GetVolumeSnapshotClass())
}

func TestExpandClusterStorageConfigurationClassesWithoutTier(t *testing.T) {
	// Storage classes alone must produce a non-nil configuration, otherwise the
	// backend clears volume_snapshot_class when the message is omitted.
	config := expandClusterStorageConfiguration([]interface{}{
		map[string]interface{}{
			clusterVolumeSnapshotClassFieldName: "csi-snapclass",
		},
	})
	require.NotNil(t, config)
	assert.Nil(t, config.DatabaseStorageClass)
	assert.Nil(t, config.SnapshotStorageClass)
	assert.Equal(t, "csi-snapclass", config.GetVolumeSnapshotClass())
}

func TestExpandClusterStorageConfigurationEmptyReturnsNil(t *testing.T) {
	config := expandClusterStorageConfiguration([]interface{}{
		map[string]interface{}{
			clusterStorageTierTypeFieldName:      "",
			clusterDatabaseStorageClassFieldName: "",
			clusterSnapshotStorageClassFieldName: "",
			clusterVolumeSnapshotClassFieldName:  "",
		},
	})
	assert.Nil(t, config)
}

func TestFlattenClusterStorageConfigurationStorageClasses(t *testing.T) {
	flattened := flattenClusterStorageConfiguration(&qcCluster.ClusterStorageConfiguration{
		StorageTierType:      commonv1.StorageTierType_STORAGE_TIER_TYPE_COST_OPTIMISED,
		DatabaseStorageClass: newPointer("premium-perf1-stackit"),
		SnapshotStorageClass: newPointer("premium-perf1-stackit"),
		VolumeSnapshotClass:  newPointer("csi-snapclass"),
	})
	expected := []interface{}{
		map[string]interface{}{
			clusterStorageTierTypeFieldName:      "STORAGE_TIER_TYPE_COST_OPTIMISED",
			clusterDatabaseStorageClassFieldName: "premium-perf1-stackit",
			clusterSnapshotStorageClassFieldName: "premium-perf1-stackit",
			clusterVolumeSnapshotClassFieldName:  "csi-snapclass",
		},
	}
	assert.Equal(t, expected, flattened)
}

// TestClusterVolumeSnapshotClassSurvivesUpdate reproduces the UI-vs-Terraform
// wipe: volume_snapshot_class was set outside Terraform (UI) and refreshed into
// state, while the Terraform config does not manage it. The backend clears
// volume_snapshot_class when an update omits it, so the diff engine must
// preserve the state value and expandCluster must echo it in the request.
func TestClusterVolumeSnapshotClassSurvivesUpdate(t *testing.T) {
	r := resourceAccountsCluster()

	state := &terraform.InstanceState{
		ID: "cluster-1",
		Attributes: map[string]string{
			"id":                                   "cluster-1",
			"account_id":                           "acc-1",
			"name":                                 "test-shared-qdrant",
			"cloud_provider":                       "private",
			"cloud_region":                         "private",
			"configuration.#":                      "1",
			"configuration.0.number_of_nodes":      "3",
			"configuration.0.node_configuration.#": "1",
			"configuration.0.node_configuration.0.package_id":                        "pkg-1",
			"configuration.0.cluster_storage_configuration.#":                        "1",
			"configuration.0.cluster_storage_configuration.0.storage_tier_type":      "STORAGE_TIER_TYPE_COST_OPTIMISED",
			"configuration.0.cluster_storage_configuration.0.database_storage_class": "premium-perf1-stackit",
			"configuration.0.cluster_storage_configuration.0.snapshot_storage_class": "premium-perf1-stackit",
			"configuration.0.cluster_storage_configuration.0.volume_snapshot_class":  "csi-snapclass",
		},
	}

	for name, rawConfig := range map[string]map[string]interface{}{
		// The customer's shape: block present, storage classes not managed.
		"block present without classes": {
			"name":           "test-shared-qdrant",
			"cloud_provider": "private",
			"cloud_region":   "private",
			"configuration": []interface{}{
				map[string]interface{}{
					"number_of_nodes": 3,
					"node_configuration": []interface{}{
						map[string]interface{}{"package_id": "pkg-1"},
					},
					"cluster_storage_configuration": []interface{}{
						map[string]interface{}{"storage_tier_type": "STORAGE_TIER_TYPE_COST_OPTIMISED"},
					},
				},
			},
		},
		// The whole storage block omitted from the Terraform config.
		"storage block omitted": {
			"name":           "test-shared-qdrant",
			"cloud_provider": "private",
			"cloud_region":   "private",
			"configuration": []interface{}{
				map[string]interface{}{
					"number_of_nodes": 3,
					"node_configuration": []interface{}{
						map[string]interface{}{"package_id": "pkg-1"},
					},
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(rawConfig), nil)
			require.NoError(t, err)

			data, err := schema.InternalMap(r.SchemaMap()).Data(state, diff)
			require.NoError(t, err)

			cluster, _, err := expandCluster(data, "acc-1")
			require.NoError(t, err)

			storageConfig := cluster.GetConfiguration().GetClusterStorageConfiguration()
			require.NotNil(t, storageConfig)
			assert.Equal(t, "csi-snapclass", storageConfig.GetVolumeSnapshotClass())
			assert.Equal(t, "premium-perf1-stackit", storageConfig.GetDatabaseStorageClass())
			assert.Equal(t, "premium-perf1-stackit", storageConfig.GetSnapshotStorageClass())
		})
	}
}
