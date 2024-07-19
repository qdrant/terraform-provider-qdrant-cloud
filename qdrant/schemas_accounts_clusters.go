package qdrant

import (
	"fmt"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qc "terraform-provider-qdrant-cloud/v1/internal/client"
)

const (
	clustersFieldTemplate      = "Clusters Schema %s field"
	clustersAccountIDFieldName = "account_id"
	clustersClustersFieldName  = "clusters"

	clusterFieldTemplate                = "Cluster Schema %s field"
	clusterIdentifierFieldName          = "id"
	clusterCreatedAtFieldName           = "created_at"
	clusterAccountIDFieldName           = "account_id"
	clusterNameFieldName                = "name"
	clusterCloudProviderFieldName       = "cloud_provider"
	clusterCloudRegionFieldName         = "cloud_region"
	clusterVersionFieldName             = "version"
	clusterPrivateRegionIDFieldName     = "private_region_id"
	clusterMarkedForDeletionAtFieldName = "marked_for_deletion_at"
	clusterURLFieldName                 = "url"
	configurationFieldName              = "configuration"
	nodeConfigurationFieldName          = "node_configuration"
	numberOfNodesFieldName              = "number_of_nodes"
	packageIDFieldName                  = "package_id"
)

// accountsClustersSchema defines the schema for a cluster list data-source.
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
				Schema:      accountsClusterSchema(true),
			},
		},
	}
}

// accountsClusterSchema defines the schema for a cluster resource or data-source.
func accountsClusterSchema(asDataSource bool) map[string]*schema.Schema {
	maxItems := 1
	if asDataSource {
		// We should not set Max Items
		maxItems = 0
	}
	return map[string]*schema.Schema{
		clusterIdentifierFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Identifier of the cluster"),
			Type:        schema.TypeString,
			Required:    asDataSource,
			Computed:    !asDataSource,
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
			Required:    !asDataSource,
			Computed:    asDataSource,
		},
		clusterCloudProviderFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Cloud provider where the cluster resides"),
			Type:        schema.TypeString,
			Required:    !asDataSource,
			ForceNew:    !asDataSource, // Cross provider migration isn't supported
			Computed:    asDataSource,
		},
		clusterCloudRegionFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Cloud region where the cluster resides"),
			Type:        schema.TypeString,
			Required:    !asDataSource,
			ForceNew:    !asDataSource, // Cross region migration isn't supported
			Computed:    asDataSource,
		},
		clusterPrivateRegionIDFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Identifier of the Private Region"),
			Type:        schema.TypeString,
			Computed:    asDataSource,
			Optional:    !asDataSource,
		},
		clusterMarkedForDeletionAtFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Timestamp when this cluster was marked for deletion"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		clusterVersionFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Version of the Qdrant cluster"),
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    !asDataSource,
		},
		clusterURLFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "The URL of the endpoint of the Qdrant cluster"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		configurationFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "The configuration options of a cluster"),
			Type:        schema.TypeList, // There is a single required item only, no need for a set.
			Required:    !asDataSource,
			Computed:    asDataSource,
			MaxItems:    maxItems,
			Elem: &schema.Resource{
				Schema: accountsClusterConfigurationSchema(asDataSource),
			},
		},
	}
}

// accountsClusterConfigurationSchema defines the schema for a cluster configuration resource or data-source.
func accountsClusterConfigurationSchema(asDataSource bool) map[string]*schema.Schema {
	maxItems := 1
	if asDataSource {
		// We should not set Max Items
		maxItems = 0
	}
	return map[string]*schema.Schema{
		numberOfNodesFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "The number of nodes in the cluster"),
			Type:        schema.TypeInt,
			Required:    !asDataSource,
			Computed:    asDataSource,
		},
		nodeConfigurationFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "The node configuration options of a cluster"),
			Type:        schema.TypeList,
			Required:    !asDataSource,
			Computed:    asDataSource,
			MaxItems:    maxItems,
			Elem: &schema.Resource{
				Schema: accountsClusterNodeConfigurationSchema(asDataSource),
			},
		},
	}
}

// accountsClusterNodeConfigurationSchema defines the schema for a cluster node configuration resource or data-source.
func accountsClusterNodeConfigurationSchema(asDataSource bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		packageIDFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "The package identifier (specifying: CPU, Memory, and disk size)"),
			Type:        schema.TypeString,
			Required:    !asDataSource,
			Computed:    asDataSource,
		},
	}
}

func expandCluster(d *schema.ResourceData, accountID string) (qc.ClusterSchema, error) {
	// Check if we need to override the default
	if v, ok := d.GetOk(clusterAccountIDFieldName); ok {
		accountID = v.(string)
	}
	if accountID == "" {
		return qc.ClusterSchema{}, fmt.Errorf("account ID not specified")
	}
	id := d.Get(clusterIdentifierFieldName)
	name := d.Get(clusterNameFieldName)
	cloudProvider := d.Get(clusterCloudProviderFieldName)
	cloudRegion := d.Get(clusterCloudRegionFieldName)

	var uuid_id openapi_types.UUID
	if id != nil && id.(string) != "" {
		uuid_id = uuid.MustParse(id.(string))
	}

	cluster := qc.ClusterSchema{
		Id:            newPointer(uuid_id),
		Name:          name.(string),
		CloudProvider: newPointer(qc.ClusterSchemaCloudProvider(cloudProvider.(string))),
		CloudRegion:   newPointer(qc.ClusterSchemaCloudRegion(cloudRegion.(string))),
		AccountId:     uuid.MustParse(accountID),
	}
	if v, ok := d.GetOk(clusterVersionFieldName); ok {
		cluster.Version = newPointer(v.(string))
	}
	if v, ok := d.GetOk(clusterMarkedForDeletionAtFieldName); ok {
		cluster.MarkedForDeletionAt = newPointer(parseTime(v.(string)))
	}
	if v, ok := d.GetOk(clusterCreatedAtFieldName); ok {
		cluster.CreatedAt = newPointer(parseTime(v.(string)))
	}
	if v, ok := d.GetOk(clusterURLFieldName); ok {
		cluster.Url = newPointer(v.(string))
	}
	if v, ok := d.GetOk(clusterPrivateRegionIDFieldName); ok {
		cluster.PrivateRegionId = newPointer(uuid.MustParse(v.(string)))
	}
	if v, ok := d.GetOk(configurationFieldName); ok {
		configuration := expandClusterConfiguration(v.([]interface{}))
		cluster.Configuration = *configuration
	}
	return cluster, nil
}

func expandClusterConfiguration(v []interface{}) *qc.ClusterConfigurationSchema {
	config := qc.ClusterConfigurationSchema{}
	for _, m := range v {
		item := m.(map[string]interface{})
		if v, ok := item[numberOfNodesFieldName]; ok {
			config.NumberOfNodes = v.(int)
		}
		if v, ok := item[nodeConfigurationFieldName]; ok {
			nodeConfig := expandNodeConfiguration(v.([]interface{}))
			if nodeConfig != nil {
				config.NodeConfiguration = *nodeConfig
			}
		}
	}
	return &config
}

func expandNodeConfiguration(v []interface{}) *qc.NodeConfigurationSchema {
	config := qc.NodeConfigurationSchema{}
	for _, m := range v {
		item := m.(map[string]interface{})
		if v, ok := item[packageIDFieldName]; ok {
			config.PackageId = uuid.MustParse(v.(string))
		}
	}
	return &config
}

// flattenClusters creates an interface from a list of clusters for easy storage in Terraform.
func flattenClusters(clusters []qc.ClusterSchema) []interface{} {
	var flattenedClusters []interface{}
	for _, cluster := range clusters {
		flattenedClusters = append(flattenedClusters, flattenCluster(&cluster))
	}
	return flattenedClusters
}

// flattenCluster creates a map from a cluster for easy storage in Terraform.
func flattenCluster(cluster *qc.ClusterSchema) map[string]interface{} {
	var privateRegionIdStr string
	if cluster.PrivateRegionId != nil {
		privateRegionIdStr = cluster.PrivateRegionId.String()
	}
	return map[string]interface{}{
		clusterIdentifierFieldName:          cluster.Id.String(),
		clusterCreatedAtFieldName:           formatTime(cluster.CreatedAt),
		clusterAccountIDFieldName:           cluster.AccountId.String(),
		clusterNameFieldName:                cluster.Name,
		clusterCloudProviderFieldName:       string(derefPointer(cluster.CloudProvider)),
		clusterCloudRegionFieldName:         string(derefPointer(cluster.CloudRegion)),
		clusterVersionFieldName:             derefPointer(cluster.Version),
		clusterPrivateRegionIDFieldName:     privateRegionIdStr,
		clusterMarkedForDeletionAtFieldName: formatTime(cluster.MarkedForDeletionAt),
		clusterURLFieldName:                 derefPointer(cluster.Url),
		configurationFieldName:              flattenClusterConfiguration(cluster.Configuration),
	}
}

// flattenClusterConfiguration creates a map from a cluster configuration for easy storage in Terraform.
func flattenClusterConfiguration(clusterConfig qc.ClusterConfigurationSchema) []interface{} {
	return []interface{}{
		map[string]interface{}{
			numberOfNodesFieldName:     clusterConfig.NumberOfNodes,
			nodeConfigurationFieldName: flattenNodeConfiguration(clusterConfig.NodeConfiguration),
		},
	}
}

// flattenNodeConfiguration creates a map from a node configuration for easy storage in Terraform.
func flattenNodeConfiguration(nodeConfig qc.NodeConfigurationSchema) []interface{} {
	return []interface{}{
		map[string]interface{}{
			packageIDFieldName: nodeConfig.PackageId.String(),
		},
	}
}
