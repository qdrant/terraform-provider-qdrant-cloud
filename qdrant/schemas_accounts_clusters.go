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
)

// accountsClustersSchema defines the schema for a cluster list resource.
// Returns a pointer to the schema.Resource object.
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
// Returns a pointer to the schema.Resource object.
func accountsClusterSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		clusterIdentifierFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Identifier of the cluster"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		clusterCreatedAtFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Timestamp then the cluster is created"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		"owner_id": { // TODO: Remove?
			Description: fmt.Sprintf(clusterFieldTemplate, "Identifier of the owner"),
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		clusterAccountIDFieldName: { // If set here, overrides account ID in provider
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
		"cloud_region_setup": {
			Description: fmt.Sprintf(clusterFieldTemplate, "Cloud region setup of the cluster"),
			Type:        schema.TypeString,
			Optional:    true,
		},
		"private_region_id": {
			Description: fmt.Sprintf(clusterFieldTemplate, "Identifier of the Private Region"),
			Type:        schema.TypeString,
			Optional:    true,
		},
		"current_configuration_id": {
			Description: fmt.Sprintf(clusterFieldTemplate, "Identifier of the current configuration"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		"encryption_key_id": {
			Description: fmt.Sprintf(clusterFieldTemplate, "Identifier of the encrption key"),
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"marked_for_deletion_at": {
			Description: fmt.Sprintf(clusterFieldTemplate, "Timstamp when this cluster was marked for deletion"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		clusterVersionFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Version of the qdrant cluster"),
			Type:        schema.TypeString,
			Optional:    true,
		},
		"url": {
			Description: fmt.Sprintf(clusterFieldTemplate, "The URL of the endpoint of the qdrant cluster"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		"state": {
			Description: "TODO",
			Type:        schema.TypeMap,
			Computed:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"configuration": {
			Description: fmt.Sprintf(clusterFieldTemplate, "The configuration options of a cluster"),
			Type:        schema.TypeList, // There is a single required item only, no need for a set.
			Required:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"num_nodes_max": {
						Description: fmt.Sprintf(clusterFieldTemplate, "The maximum number of nodes from the cluster"),
						Type:        schema.TypeInt,
						Required:    true,
					},
					"num_nodes": {
						Description: fmt.Sprintf(clusterFieldTemplate, "The number of nodes from the cluster"),
						Type:        schema.TypeInt,
						Required:    true,
					},
					"node_configuration": {
						Description: fmt.Sprintf(clusterFieldTemplate, "The node configuration options of a cluster"),
						Type:        schema.TypeList,
						Required:    true,
						MaxItems:    1,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"package_id": {
									Description: fmt.Sprintf(clusterFieldTemplate, "The package identifier (specifying: CPU, Memory and disk size)"),
									Type:        schema.TypeString,
									Required:    true,
								},
							},
						},
					},
				},
			},
		},
		"resources": {
			Description: "TODO",
			Type:        schema.TypeMap,
			Computed:    true,
			Elem: &schema.Schema{
				Description: "TODO",
				Type:        schema.TypeString,
			},
		},
		"total_extra_disk": {
			Description: fmt.Sprintf(clusterFieldTemplate, "The total ammount of extra disk in relation to the choosen package (in Gib)"),
			Type:        schema.TypeInt,
			Optional:    true,
		},
	}
}

func expandClusterIn(d *schema.ResourceData, accountID string) (qc.ClusterIn, error) {
	// Check if we need to override the default
	if v, ok := d.GetOk("account_id"); ok {
		accountID = v.(string)
	}
	if accountID == "" {
		return qc.ClusterIn{}, fmt.Errorf("account ID not specified")
	}
	name := d.Get("name")
	cloudProvider := d.Get("cloud_provider")
	cloudRegion := d.Get("cloud_region")

	cluster := qc.ClusterIn{
		Name:          name.(string),
		CloudProvider: qc.ClusterInCloudProvider(cloudProvider.(string)),
		CloudRegion:   qc.ClusterInCloudRegion(cloudRegion.(string)),
		AccountId:     &accountID,
	}
	if v, ok := d.GetOk("configuration"); ok {
		configuration := expandClusterConfigurationIn(v.([]interface{}))
		cluster.Configuration = *configuration
	}
	return cluster, nil
}

func expandClusterConfigurationIn(v []interface{}) *qc.ClusterConfigurationIn {
	config := qc.ClusterConfigurationIn{}
	for _, m := range v {
		item := m.(map[string]interface{})
		if v, ok := item["num_nodes_max"]; ok {
			config.NumNodesMax = v.(int)
		}
		if v, ok := item["num_nodes"]; ok {
			config.NumNodes = v.(int)
		}
		if v, ok := item["node_configuration"]; ok {
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
		if v, ok := item["package_id"]; ok {
			config.PackageId = v.(string)
		}
	}
	return &config
}

// flattenCluster creates a map from a cluster for easy storage on terraform.
func flattenCluster(cluster *qc.ClusterOut) map[string]interface{} {
	result := map[string]interface{}{
		clusterIdentifierFieldName:    cluster.Id,
		clusterAccountIDFieldName:     cluster.AccountId,
		clusterNameFieldName:          cluster.Name,
		clusterCloudProviderFieldName: cluster.CloudProvider,
		clusterCloudRegionFieldName:   cluster.CloudRegion,
		"configuration":               flattenClusterConfiguration(cluster.Configuration),
	}
	return result
}

// flattenClusterConfiguration creates a map from a cluster for easy storage on terraform.
func flattenClusterConfiguration(clusterConfig *qc.ClusterConfigurationOut) []interface{} {
	result := []interface{}{
		map[string]interface{}{
			//"id":            clusterConfig.Id,
			"num_nodes":     clusterConfig.NumNodes,
			"num_nodes_max": clusterConfig.NumNodesMax,
		},
	}
	return result
}
