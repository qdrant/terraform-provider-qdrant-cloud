package qdrant

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qcCluster "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/v1"
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
	clusterStatusFieldName                     = "status"
	clusterStatusVersionFieldName              = "version"
	clusterStatusNodesUpFieldName              = "nodes_up"
	clusterStatusRestartedAtFieldName          = "restarted_at"
	clusterStatusPhaseFieldName                = "phase"
	clusterStatusReasonFieldName               = "reason"
	clusterStatusResourcesFieldName            = "resources"
	clusterStatusScalabilityInfoFieldName      = "scalability_info"
	clusterNodeResourcesSummaryDiskFieldName   = "disk"
	clusterNodeResourcesSummaryRamFieldName    = "ram"
	clusterNodeResourcesSummaryCpuFieldName    = "cpu"
	clusterNodeResourcesBaseFieldName          = "base"
	clusterNodeResourcesComplimentaryFieldName = "complimentary"
	clusterNodeResourcesAdditionalFieldName    = "additional"
	clusterNodeResourcesReservedFieldName      = "reserved"
	clusterNodeResourcesAvailableFieldName     = "available"
	clusterScalabilityInfoStatusFieldName      = "status"
	clusterScalabilityInfoReasonFieldName      = "reason"
	configurationFieldName                     = "configuration"
	nodeConfigurationFieldName                 = "node_configuration"
	numberOfNodesFieldName                     = "number_of_nodes"
	packageIDFieldName                         = "package_id"
	resourceConfigurationsFieldName            = "resource_configurations"
	resourceConfigurationAmountFieldName       = "amount"
	resourceConfigurationResourceTypeFieldName = "resource_type"
	resourceConfigurationResourceUnitFieldName = "resource_unit"

	// Backward compatibility.
	fieldAmount       = "amount"
	fieldResourceType = "resource_type"
	fieldResourceUnit = "resource_unit"

	descriptionAmount       = "The amount of the resource"
	descriptionResourceType = "The type of the resource"
	descriptionResourceUnit = "The unit of the resource"
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
		clusterStatusFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "The status of the cluster"),
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: accountsClusterStatusSchema(),
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
			Description: descriptionResourceConfiguration,
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: resourceConfigurationsSchema(asDataSource),
			},
		},
	}
}

// resourceConfigurationsSchema defines the schema structure for resource configurations.
// This is for backwards compatibility.
func resourceConfigurationsSchema(asDataSource bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		fieldAmount: {
			Description: descriptionAmount,
			Type:        schema.TypeInt,
			Required:    !asDataSource,
			Computed:    asDataSource,
		},
		fieldResourceType: {
			Description: descriptionResourceType,
			Type:        schema.TypeString,
			Required:    !asDataSource,
			Computed:    asDataSource,
		},
		fieldResourceUnit: {
			Description: descriptionResourceUnit,
			Type:        schema.TypeString,
			Required:    !asDataSource,
			Computed:    asDataSource,
		},
	}
}

// accountsClusterStatusSchema defines the schema for a cluster status.
func accountsClusterStatusSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		clusterStatusVersionFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Version of the cluster software"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		clusterStatusNodesUpFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Number of cluster nodes that are up and running"),
			Type:        schema.TypeInt,
			Computed:    true,
		},
		clusterStatusRestartedAtFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "The date and time when the cluster was restarted"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		clusterStatusPhaseFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Current phase of the cluster"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		clusterStatusReasonFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Reason for the current phase of the cluster"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		clusterStatusResourcesFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "The resources used by the cluster per node"),
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterNodeResourcesSummarySchema(),
			},
		},
		clusterStatusScalabilityInfoFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Whether the cluster can be scaled up or down"),
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterScalabilityInfoSchema(),
			},
		},
	}
}

// clusterNodeResourcesSummarySchema defines the schema for a cluster node resources summary.
func clusterNodeResourcesSummarySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		clusterNodeResourcesSummaryDiskFieldName: {
			Description: "Disk resources",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterNodeResourcesSchema(),
			},
		},
		clusterNodeResourcesSummaryRamFieldName: {
			Description: "Memory resources",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterNodeResourcesSchema(),
			},
		},
		clusterNodeResourcesSummaryCpuFieldName: {
			Description: "CPU resources",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: clusterNodeResourcesSchema(),
			},
		},
	}
}

// clusterNodeResourcesSchema defines the schema for a cluster node resources.
func clusterNodeResourcesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		clusterNodeResourcesBaseFieldName: {
			Description: "Base resources that are part of the standard allocation for the cluster per node.",
			Type:        schema.TypeFloat,
			Computed:    true,
		},
		clusterNodeResourcesComplimentaryFieldName: {
			Description: "Complimentary resources provided to the cluster at no additional cost.",
			Type:        schema.TypeFloat,
			Computed:    true,
		},
		clusterNodeResourcesAdditionalFieldName: {
			Description: "Additional resources allocated to the cluster.",
			Type:        schema.TypeFloat,
			Computed:    true,
		},
		clusterNodeResourcesReservedFieldName: {
			Description: "The reserved is the amount used by the system, which cannot be used by the database itself.",
			Type:        schema.TypeFloat,
			Computed:    true,
		},
		clusterNodeResourcesAvailableFieldName: {
			Description: "The available is the total (base+complimentary+additional) - reserved",
			Type:        schema.TypeFloat,
			Computed:    true,
		},
	}
}

// clusterScalabilityInfoSchema defines the schema for a cluster scalability info.
func clusterScalabilityInfoSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		clusterScalabilityInfoStatusFieldName: {
			Description: "The current scalability status of the cluster.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		clusterScalabilityInfoReasonFieldName: {
			Description: "Optional human-readable reason providing more context about the scalability status.",
			Type:        schema.TypeString,
			Computed:    true,
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
		Id:                    id.(string),
		Name:                  name.(string),
		CloudProviderId:       cloudProvider.(string),
		CloudProviderRegionId: cloudRegion.(string),
		AccountId:             accountID,
	}
	if v, ok := d.GetOk(clusterMarkedForDeletionAtFieldName); ok {
		cluster.DeletedAt = parseTime(v.(string))
	}
	if v, ok := d.GetOk(clusterCreatedAtFieldName); ok {
		cluster.CreatedAt = parseTime(v.(string))
	}
	if v, ok := d.GetOk(clusterPrivateRegionIDFieldName); ok {
		// Note this field has been merged with cloud-region (if provider indicates hybrid cloud)
		if cluster.CloudProviderId == hybridCloudClusterID {
			cluster.CloudProviderRegionId = v.(string)
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
	// The status is a read-only object, so no need to expand it.
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
			version := v.(string)
			if version != "" {
				config.Version = newPointer(v.(string))
			}
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
	if cluster.CloudProviderId == hybridCloudClusterID {
		// For backewards compatibility extract the region ID into separate field.
		privateRegionIdStr = cluster.GetCloudProviderRegionId()
	}
	return map[string]interface{}{
		clusterIdentifierFieldName:          cluster.GetId(),
		clusterCreatedAtFieldName:           formatTime(cluster.GetCreatedAt()),
		clusterAccountIDFieldName:           cluster.GetAccountId(),
		clusterNameFieldName:                cluster.GetName(),
		clusterCloudProviderFieldName:       cluster.GetCloudProviderId(),
		clusterCloudRegionFieldName:         cluster.GetCloudProviderRegionId(),
		clusterPrivateRegionIDFieldName:     privateRegionIdStr,
		clusterMarkedForDeletionAtFieldName: formatTime(cluster.GetDeletedAt()),
		clusterURLFieldName:                 cluster.GetState().GetEndpoint().GetUrl(),
		configurationFieldName:              flattenClusterConfiguration(cluster.GetConfiguration()),
		clusterStatusFieldName:              flattenClusterState(cluster.GetState()),
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

// flattenClusterState creates a map from a cluster state for easy storage in Terraform.
func flattenClusterState(state *qcCluster.ClusterState) []interface{} {
	if state == nil {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			clusterStatusVersionFieldName:         state.GetVersion(),
			clusterStatusNodesUpFieldName:         int(state.GetNodesUp()),
			clusterStatusRestartedAtFieldName:     formatTime(state.GetRestartedAt()),
			clusterStatusPhaseFieldName:           state.GetPhase().String(),
			clusterStatusReasonFieldName:          state.GetReason(),
			clusterStatusResourcesFieldName:       flattenClusterNodeResourcesSummary(state.GetResources()),
			clusterStatusScalabilityInfoFieldName: flattenClusterScalabilityInfo(state.GetScalabilityInfo()),
		},
	}
}

// flattenClusterNodeResourcesSummary creates a map from a cluster node resources summary for easy storage in Terraform.
func flattenClusterNodeResourcesSummary(summary *qcCluster.ClusterNodeResourcesSummary) []interface{} {
	if summary == nil {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			clusterNodeResourcesSummaryDiskFieldName: flattenClusterNodeResources(summary.GetDisk()),
			clusterNodeResourcesSummaryRamFieldName:  flattenClusterNodeResources(summary.GetRam()),
			clusterNodeResourcesSummaryCpuFieldName:  flattenClusterNodeResources(summary.GetCpu()),
		},
	}
}

// flattenClusterNodeResources creates a map from a cluster node resources for easy storage in Terraform.
func flattenClusterNodeResources(resources *qcCluster.ClusterNodeResources) []interface{} {
	if resources == nil {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			clusterNodeResourcesBaseFieldName:          resources.GetBase(),
			clusterNodeResourcesComplimentaryFieldName: resources.GetComplimentary(),
			clusterNodeResourcesAdditionalFieldName:    resources.GetAdditional(),
			clusterNodeResourcesReservedFieldName:      resources.GetReserved(),
			clusterNodeResourcesAvailableFieldName:     resources.GetAvailable(),
		},
	}
}

// flattenClusterScalabilityInfo creates a map from a cluster scalability info for easy storage in Terraform.
func flattenClusterScalabilityInfo(info *qcCluster.ClusterScalabilityInfo) []interface{} {
	if info == nil {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			clusterScalabilityInfoStatusFieldName: info.GetStatus().String(),
			clusterScalabilityInfoReasonFieldName: info.GetReason(),
		},
	}
}
