package qdrant

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qc "terraform-provider-qdrant-cloud/v1/internal/client"
)

const (
	clustersFieldTemplate      = "Clusters Schema %s field"
	clustersAccountIDFieldName = "account_id"
	clustersClustersFieldName  = "clusters"

	clusterFieldTemplate                        = "Cluster Schema %s field"
	clusterIdentifierFieldName                  = "id"
	clusterCreatedAtFieldName                   = "created_at"
	clusterAccountIDFieldName                   = "account_id"
	clusterNameFieldName                        = "name"
	clusterCloudProviderFieldName               = "cloud_provider"
	clusterCloudRegionFieldName                 = "cloud_region"
	clusterCloudRegionAvailabilityZoneFieldName = "cloud_region_az"
	clusterVersionFieldName                     = "version"
	clusterCloudRegionSetupFieldName            = "cloud_region_setup"
	clusterPrivateRegionIDFieldName             = "private_region_id"
	clusterCurrentConfigurationIDFieldName      = "current_configuration_id"
	clusterEncryptionKeyIDFieldName             = "encryption_key_id"
	clusterMarkedForDeletionAtFieldName         = "marked_for_deletion_at"
	clusterURLFieldName                         = "url"
	clusterTotalExtraDiskFieldName              = "total_extra_disk"
	configurationFieldName                      = "configuration"
	nodeConfigurationFieldName                  = "node_configuration"
	numNodesMaxFieldName                        = "num_nodes_max"
	numNodesFieldName                           = "num_nodes"
	packageIDFieldName                          = "package_id"
)

// accountsClustersSchema defines the schema for a cluster list resource.
func accountsClustersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		clustersAccountIDFieldName: {
			Description: fmt.Sprintf(clustersFieldTemplate, "Identifier of the account"),
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		clustersClustersFieldName: {
			Description: fmt.Sprintf(clustersFieldTemplate, "List of clusters"),
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Description: fmt.Sprintf(clustersFieldTemplate, "Individual cluster"),
				Schema:      accountsClusterSchema(),
			},
		},
	}
}

// accountsClusterSchema defines the schema for a cluster resource.
func accountsClusterSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		clusterIdentifierFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Identifier of the cluster"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		clusterCreatedAtFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Timestamp when the cluster is created"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		clusterAccountIDFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Identifier of the account"),
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
		},
		clusterNameFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Name of the cluster"),
			Type:        schema.TypeString,
			Required:    true,
		},
		clusterCloudProviderFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Cloud provider where the cluster resides"),
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true, // Cross provider migration isn't supported
		},
		clusterCloudRegionFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Cloud region where the cluster resides"),
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true, // Cross region migration isn't supported
		},
		clusterCloudRegionAvailabilityZoneFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Cloud region availability zone where the cluster resides"),
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		clusterCloudRegionSetupFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Cloud region setup of the cluster"),
			Type:        schema.TypeString,
			Optional:    true,
		},
		clusterPrivateRegionIDFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Identifier of the Private Region"),
			Type:        schema.TypeString,
			Optional:    true,
		},
		clusterCurrentConfigurationIDFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Identifier of the current configuration"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		clusterEncryptionKeyIDFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Identifier of the encryption key"),
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		clusterMarkedForDeletionAtFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Timestamp when this cluster was marked for deletion"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		clusterVersionFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Version of the Qdrant cluster"),
			Type:        schema.TypeString,
			Optional:    true,
		},
		clusterURLFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "The URL of the endpoint of the Qdrant cluster"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		configurationFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "The configuration options of a cluster"),
			Type:        schema.TypeList, // There is a single required item only, no need for a set.
			Required:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					numNodesMaxFieldName: {
						Description: fmt.Sprintf(clusterFieldTemplate, "The maximum number of nodes in the cluster"),
						Type:        schema.TypeInt,
						Required:    true,
					},
					numNodesFieldName: {
						Description: fmt.Sprintf(clusterFieldTemplate, "The number of nodes in the cluster"),
						Type:        schema.TypeInt,
						Required:    true,
					},
					nodeConfigurationFieldName: {
						Description: fmt.Sprintf(clusterFieldTemplate, "The node configuration options of a cluster"),
						Type:        schema.TypeList,
						Required:    true,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								packageIDFieldName: {
									Description: fmt.Sprintf(clusterFieldTemplate, "The package identifier (specifying: CPU, Memory, and disk size)"),
									Type:        schema.TypeString,
									Required:    true,
								},
							},
						},
					},
				},
			},
		},
		clusterTotalExtraDiskFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "The total amount of extra disk in relation to the chosen package (in GiB)"),
			Type:        schema.TypeInt,
			Optional:    true,
		},
	}
}

func expandClusterIn(d *schema.ResourceData, accountID string) (qc.ClusterIn, error) {
	// Check if we need to override the default
	if v, ok := d.GetOk(clusterAccountIDFieldName); ok {
		accountID = v.(string)
	}
	if accountID == "" {
		return qc.ClusterIn{}, fmt.Errorf("account ID not specified")
	}
	name := d.Get(clusterNameFieldName)
	cloudProvider := d.Get(clusterCloudProviderFieldName)
	cloudRegion := d.Get(clusterCloudRegionFieldName)

	cluster := qc.ClusterIn{
		Name:          name.(string),
		CloudProvider: qc.ClusterInCloudProvider(cloudProvider.(string)),
		CloudRegion:   qc.ClusterInCloudRegion(cloudRegion.(string)),
		AccountId:     &accountID,
	}

	if v, ok := d.GetOk(clusterVersionFieldName); ok {
		val := v.(string)
		cluster.Version = &val
	}
	if v, ok := d.GetOk(clusterCloudRegionAvailabilityZoneFieldName); ok {
		val := v.(string)
		cluster.CloudRegionAz = &val
	}
	if v, ok := d.GetOk(clusterCloudRegionSetupFieldName); ok {
		val := v.(string)
		cluster.CloudRegionSetup = &val
	}
	if v, ok := d.GetOk(clusterPrivateRegionIDFieldName); ok {
		val := v.(string)
		cluster.PrivateRegionId = &val
	}
	/*if v, ok := d.GetOk(clusterEncryptionKeyIDFieldName); ok {
		val := v.(string)
		cluster.EncryptionKeyId = &val
	}
	if v, ok := d.GetOk(clusterTotalExtraDiskFieldName); ok {
		extraDisk := v.(int)
		cluster.TotalExtraDisk = &extraDisk
	}*/
	if v, ok := d.GetOk(configurationFieldName); ok {
		configuration := expandClusterConfigurationIn(v.([]interface{}))
		cluster.Configuration = *configuration
	}
	return cluster, nil
}

func expandClusterConfigurationIn(v []interface{}) *qc.ClusterConfigurationIn {
	config := qc.ClusterConfigurationIn{}
	for _, m := range v {
		item := m.(map[string]interface{})
		if v, ok := item[numNodesMaxFieldName]; ok {
			config.NumNodesMax = v.(int)
		}
		if v, ok := item[numNodesFieldName]; ok {
			config.NumNodes = v.(int)
		}
		if v, ok := item[nodeConfigurationFieldName]; ok {
			nodeConfig := expandNodeConfigurationIn(v.([]interface{}))
			if nodeConfig != nil {
				config.NodeConfiguration = *nodeConfig
			}
		}
	}
	return &config
}

func expandNodeConfigurationIn(v []interface{}) *qc.NodeConfiguration {
	config := qc.NodeConfiguration{}
	for _, m := range v {
		item := m.(map[string]interface{})
		if v, ok := item[packageIDFieldName]; ok {
			config.PackageId = v.(string)
		}
	}
	return &config
}

// flattenClusters creates an interface from a list of clusters for easy storage in Terraform.
func flattenClusters(clusters []qc.ClusterOut) []interface{} {
	var flattenedClusters []interface{}
	for _, cluster := range clusters {
		flattenedClusters = append(flattenedClusters, flattenCluster(&cluster))
	}
	return flattenedClusters
}

// flattenCluster creates a map from a cluster for easy storage in Terraform.
func flattenCluster(cluster *qc.ClusterOut) map[string]interface{} {
	return map[string]interface{}{
		clusterIdentifierFieldName:                  cluster.Id,
		clusterCreatedAtFieldName:                   formatTime(cluster.CreatedAt),
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
		configurationFieldName:                      flattenClusterConfiguration(cluster.Configuration),
	}
}

// flattenClusterConfiguration creates a map from a cluster configuration for easy storage in Terraform.
func flattenClusterConfiguration(clusterConfig *qc.ClusterConfigurationOut) []interface{} {
	return []interface{}{
		map[string]interface{}{
			numNodesFieldName:          clusterConfig.NumNodes,
			numNodesMaxFieldName:       clusterConfig.NumNodesMax,
			nodeConfigurationFieldName: flattenNodeConfiguration(clusterConfig.NodeConfiguration),
		},
	}
}

// flattenNodeConfiguration creates a map from a node configuration for easy storage in Terraform.
func flattenNodeConfiguration(nodeConfig qc.NodeConfiguration) []interface{} {
	return []interface{}{
		map[string]interface{}{
			packageIDFieldName: nodeConfig.PackageId,
		},
	}
}
