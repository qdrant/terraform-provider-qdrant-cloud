package qdrant

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qcCluster "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/v2"
)

const (
	hybridCloudClusterID = "hybrid"

	clustersFieldTemplate      = "Clusters Schema %s field"
	clustersAccountIDFieldName = "account_id"
	clustersClustersFieldName  = "clusters"

	clusterFieldTemplate                       = "Cluster Schema %s field"
	clusterIdentifierFieldName                 = "id"
	clusterCreatedAtFieldName                  = "created_at"
	clusterAccountIDFieldName                  = "account_id"
	clusterNameFieldName                       = "name"
	clusterCloudProviderFieldName              = "cloud_provider"
	clusterCloudRegionFieldName                = "cloud_region"
	clusterVersionFieldName                    = "version"
	clusterPrivateRegionIDFieldName            = "private_region_id"
	clusterMarkedForDeletionAtFieldName        = "marked_for_deletion_at"
	clusterURLFieldName                        = "url"
	configurationFieldName                     = "configuration"
	nodeConfigurationFieldName                 = "node_configuration"
	numberOfNodesFieldName                     = "number_of_nodes"
	packageIDFieldName                         = "package_id"
	resourceConfigurationsFieldName            = "resource_configurations"
	resourceConfigurationAmountFieldName       = "amount"
	resourceConfigurationResourceTypeFieldName = "resource_type"
	resourceConfigurationResourceUnitFieldName = "resource_unit"
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
		clusterVersionFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Version of the Qdrant cluster"),
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    !asDataSource,
		},
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
		resourceConfigurationsFieldName: {
			Description: descriptionResourceConfigurations,
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: resourceConfigurationsSchema(asDataSource),
			},
		},
	}
}

func expandCluster(d *schema.ResourceData, accountID string) (*qcCluster.Cluster, error) {
	// Check if we need to override the default
	if v, ok := d.GetOk(clusterAccountIDFieldName); ok {
		accountID = v.(string)
	}
	if accountID == "" {
		return nil, fmt.Errorf("account ID not specified")
	}
	id := d.Get(clusterIdentifierFieldName)
	name := d.Get(clusterNameFieldName)
	cloudProvider := d.Get(clusterCloudProviderFieldName)
	cloudRegion := d.Get(clusterCloudRegionFieldName)

	cluster := &qcCluster.Cluster{
		Id:            id.(string),
		Name:          name.(string),
		CloudProvider: cloudProvider.(string),
		CloudRegion:   cloudRegion.(string),
		AccountId:     accountID,
	}
	if _, ok := d.GetOk(clusterMarkedForDeletionAtFieldName); ok {
		// TODO: cluster.MarkedForDeletionAt = parseTime(v.(string))
	}
	if v, ok := d.GetOk(clusterCreatedAtFieldName); ok {
		cluster.CreatedAt = parseTime(v.(string))
	}
	if v, ok := d.GetOk(clusterPrivateRegionIDFieldName); ok {
		// Note this field has been merged with cloud-region (if provider indicates hybrid cloud)
		if cluster.CloudProvider == hybridCloudClusterID {
			cluster.CloudRegion = v.(string)
		}
	}
	if v, ok := d.GetOk(configurationFieldName); ok {
		configuration := expandClusterConfiguration(v.([]interface{}))
		cluster.Configuration = configuration
	}
	if v, ok := d.GetOk(clusterURLFieldName); ok {
		cluster.State = &qcCluster.ClusterState{
			Endpoint: &qcCluster.ClusterEndpoint{
				Url: v.(string),
			},
		}
	}
	return cluster, nil
}

func expandClusterConfiguration(v []interface{}) *qcCluster.ClusterConfiguration {
	config := &qcCluster.ClusterConfiguration{}
	for _, m := range v {
		item := m.(map[string]interface{})
		if v, ok := item[numberOfNodesFieldName]; ok {
			config.NumberOfNodes = uint32(v.(int))
		}
		if v, ok := item[clusterVersionFieldName]; ok {
			config.Version = v.(string)
		}
		if v, ok := item[nodeConfigurationFieldName]; ok {
			packageId, additionalResources := expandNodeConfiguration(v.([]interface{}))
			config.PackageId = packageId
			config.AdditionalResources = additionalResources
		}
	}
	return config
}

func expandNodeConfiguration(v []interface{}) (string, *qcCluster.AdditionalResources) {
	var packageId string
	var additionalResources *qcCluster.AdditionalResources
	for _, m := range v {
		item := m.(map[string]interface{})
		if v, ok := item[packageIDFieldName]; ok {
			packageId = v.(string)
		}
		if v, ok := item[resourceConfigurationsFieldName]; ok {
			additionalResources = expandClusterNodeResourceConfigurationsToAdditionalResources(v.([]interface{}))
		}
	}
	return packageId, additionalResources
}

func expandClusterNodeResourceConfigurationsToAdditionalResources(v []interface{}) *qcCluster.AdditionalResources {
	var result *qcCluster.AdditionalResources
	for _, m := range v {
		if result == nil {
			result = &qcCluster.AdditionalResources{}
		}
		var amount int
		var resourceType, resourceUnit string
		item := m.(map[string]interface{})
		if v, ok := item[resourceConfigurationAmountFieldName]; ok {
			amount = v.(int)
		}
		if v, ok := item[resourceConfigurationResourceTypeFieldName]; ok {
			resourceType = v.(string)
		}
		if v, ok := item[resourceConfigurationResourceUnitFieldName]; ok {
			resourceUnit = v.(string)
		}
		switch ResourceType(resourceType) {
		case ResourceTypeCpu:
			// Not supported
		case ResourceTypeRam:
			// Not supported
		case ResourceTypeComplimentaryDisk:
			// Not supported
		case ResourceTypeSnapshot:
			// Not supported
		case ResourceTypeDisk:
			// Not supported:
			if ResourceUnit(resourceUnit) == ResourceUnitGi {
				result.Disk += uint32(amount)
			}
		}
	}
	return result
}

// flattenClusters creates an interface from a list of clusters for easy storage in Terraform.
func flattenClusters(clusters []*qcCluster.Cluster) []interface{} {
	var flattenedClusters []interface{}
	for _, cluster := range clusters {
		flattenedClusters = append(flattenedClusters, flattenCluster(cluster))
	}
	return flattenedClusters
}

// flattenCluster creates a map from a cluster for easy storage in Terraform.
func flattenCluster(cluster *qcCluster.Cluster) map[string]interface{} {
	var privateRegionIdStr string
	if cluster.CloudProvider == hybridCloudClusterID {
		// For backewards compatibility extract the region ID into separate field.
		privateRegionIdStr = cluster.GetCloudRegion()
	}
	return map[string]interface{}{
		clusterIdentifierFieldName:      cluster.GetId(),
		clusterCreatedAtFieldName:       formatTime(cluster.GetCreatedAt()),
		clusterAccountIDFieldName:       cluster.GetAccountId(),
		clusterNameFieldName:            cluster.GetName(),
		clusterCloudProviderFieldName:   cluster.GetCloudProvider(),
		clusterCloudRegionFieldName:     cluster.GetCloudRegion(),
		clusterPrivateRegionIDFieldName: privateRegionIdStr,
		// TODO: clusterMarkedForDeletionAtFieldName: formatTime(cluster.GetMarkedForDeletionAt()),
		clusterURLFieldName:    cluster.GetState().GetEndpoint().GetUrl(),
		configurationFieldName: flattenClusterConfiguration(cluster.GetConfiguration()),
	}
}

// flattenClusterConfiguration creates a map from a cluster configuration for easy storage in Terraform.
func flattenClusterConfiguration(clusterConfig *qcCluster.ClusterConfiguration) []interface{} {
	return []interface{}{
		map[string]interface{}{
			clusterVersionFieldName:    clusterConfig.GetVersion(),
			numberOfNodesFieldName:     int(clusterConfig.GetNumberOfNodes()),
			nodeConfigurationFieldName: flattenNodeConfiguration(clusterConfig.GetPackageId(), clusterConfig.GetAdditionalResources()),
		},
	}
}

// flattenNodeConfiguration creates a map from a packageID and additional resources for easy storage in Terraform.
// Note the TF structure is kept backwards compatible with the OpenAPI v1, so we need to map a bit here.
func flattenNodeConfiguration(packageID string, additionalResources *qcCluster.AdditionalResources) []interface{} {
	return []interface{}{
		map[string]interface{}{
			packageIDFieldName:              packageID,
			resourceConfigurationsFieldName: flattenResourceConfigurationsFromAdditionalResources(additionalResources),
		},
	}
}

// flattenResourceConfigurations flattens the resource configurations data into a format that Terraform can understand.
func flattenResourceConfigurationsFromAdditionalResources(additionalResources *qcCluster.AdditionalResources) []interface{} {
	var flattenedResourceConfigurations []interface{}
	if additionalResources.GetDisk() > 0 {
		flattenedResourceConfigurations = append(flattenedResourceConfigurations, map[string]interface{}{
			fieldAmount:       int(additionalResources.GetDisk()),
			fieldResourceType: string(ResourceTypeDisk),
			fieldResourceUnit: string(ResourceUnitGi),
		})
	}
	return flattenedResourceConfigurations
}
