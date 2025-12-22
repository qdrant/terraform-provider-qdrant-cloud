package qdrant

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qcCluster "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/v1"
	commonv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/common/v1"
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
	clusterLabelsFieldName                     = "labels"
	clusterCloudProviderFieldName              = "cloud_provider"
	clusterCloudRegionFieldName                = "cloud_region"
	clusterVersionFieldName                    = "version"
	clusterLastModifiedAtFieldName             = "last_modified_at"
	clusterPrivateRegionIDFieldName            = "private_region_id"
	clusterMarkedForDeletionAtFieldName        = "marked_for_deletion_at"
	clusterURLFieldName                        = "url"
	clusterStatusFieldName                     = "status"
	clusterStatusVersionFieldName              = "version"
	clusterDeleteBackupsOnDestroyFieldName     = "delete_backups_on_destroy"
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
	nodeSelectorFieldName                      = "node_selector"
	tolerationsFieldName                       = "tolerations"
	tolerationKeyFieldName                     = "key"
	tolerationOperatorFieldName                = "operator"
	tolerationValueFieldName                   = "value"
	tolerationEffectFieldName                  = "effect"
	tolerationSecondsFieldName                 = "toleration_seconds"
	annotationsFieldName                       = "annotations"
	allowedIpSourceRangesFieldName             = "allowed_ip_source_ranges"
	serviceTypeFieldName                       = "service_type"
	serviceAnnotationsFieldName                = "service_annotations"
	podLabelsFieldName                         = "pod_labels"
	databaseConfigurationFieldName             = "database_configuration"
	dbConfigCollectionFieldName                = "collection"
	dbConfigStorageFieldName                   = "storage"
	dbConfigServiceFieldName                   = "service"
	dbConfigLogLevelFieldName                  = "log_level"
	dbConfigTlsFieldName                       = "tls"
	dbConfigInferenceFieldName                 = "inference"
	dbConfigReservedCpuPercentageFieldName     = "reserved_cpu_percentage"
	dbConfigReservedMemoryPercentageFieldName  = "reserved_memory_percentage"
	dbConfigGpuTypeFieldName                   = "gpu_type"
	dbConfigRestartPolicyFieldName             = "restart_policy"
	dbConfigRebalanceStrategyFieldName         = "rebalance_strategy"
	dbConfigCollectionReplicationFactor        = "replication_factor"
	dbConfigCollectionWriteConsistencyFactor   = "write_consistency_factor"
	dbConfigCollectionVectorsFieldName         = "vectors"
	dbConfigCollectionVectorsOnDiskFieldName   = "on_disk"
	dbConfigStoragePerformanceFieldName        = "performance"
	dbConfigStoragePerfOptimizerCpuBudget      = "optimizer_cpu_budget"
	dbConfigStoragePerfAsyncScorer             = "async_scorer"
	dbConfigServiceApiKeyFieldName             = "api_key"
	dbConfigServiceReadOnlyApiKeyFieldName     = "read_only_api_key"
	dbConfigServiceJwtRbacFieldName            = "jwt_rbac"
	dbConfigServiceEnableTlsFieldName          = "enable_tls"
	dbConfigSecretKeyRefSecretNameFieldName    = "secret_name"
	dbConfigSecretKeyRefSecretKeyFieldName     = "secret_key"
	dbConfigTlsCertFieldName                   = "cert"
	dbConfigTlsKeyFieldName                    = "key"
	dbConfigInferenceEnabledFieldName          = "enabled"

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
		clusterLabelsFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "List of labels associated with the cluster"),
			Type:        schema.TypeList,
			Optional:    !asDataSource,
			Computed:    asDataSource,
			Elem: &schema.Resource{
				Schema: keyValSchema(asDataSource),
			},
		},
		clusterCloudProviderFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, `Cloud provider where the cluster is hosted.
Must match one of the provider IDs returned by the "qdrant.cloud.platform.v1.PlatformService.ListCloudProviders" method.
For Hybrid cloud this should be "hybrid".`),
			Type:     schema.TypeString,
			Required: !asDataSource,
			ForceNew: !asDataSource, // Cross provider migration isn't supported
			Computed: asDataSource,
		},
		clusterCloudRegionFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, `Cloud provider region where the cluster is hosted.
Must match one of the region IDs returned by the "qdrant.cloud.platform.v1.PlatformService.ListCloudProviderRegions" method.
For hybrid this should be the hybrid cloud environment ID.`),
			Type:     schema.TypeString,
			Required: !asDataSource,
			ForceNew: !asDataSource, // Cross region migration isn't supported
			Computed: asDataSource,
		},
		clusterPrivateRegionIDFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Identifier of the Hybrid cloud region"),
			Type:        schema.TypeString,
			Computed:    asDataSource,
			Optional:    !asDataSource,
			Deprecated:  `Please use cloud_provider="hybrid" and the cloud_region field instead`,
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
		clusterDeleteBackupsOnDestroyFieldName: {
			Description: "Whether to delete backups when the cluster is destroyed.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
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
		clusterLastModifiedAtFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Timestamp when the cluster configuration was last updated"),
			Type:        schema.TypeString,
			Computed:    true,
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
		nodeSelectorFieldName: {
			Description: "The node selector for this cluster in a hybrid cloud environment.",
			Type:        schema.TypeList,
			Optional:    !asDataSource,
			Computed:    asDataSource,
			Elem: &schema.Resource{
				Schema: keyValSchema(asDataSource),
			},
		},
		tolerationsFieldName: {
			Description: "List of tolerations for this cluster in a hybrid cloud environment.",
			Type:        schema.TypeList,
			Optional:    !asDataSource,
			Computed:    asDataSource,
			Elem: &schema.Resource{
				Schema: tolerationSchema(asDataSource),
			},
		},
		annotationsFieldName: {
			Description: "List of annotations for this cluster in a hybrid cloud environment.",
			Type:        schema.TypeList,
			Optional:    !asDataSource,
			Computed:    asDataSource,
			Elem: &schema.Resource{
				Schema: keyValSchema(asDataSource),
			},
		},
		allowedIpSourceRangesFieldName: {
			Description: "List of allowed IP source ranges for this cluster.",
			Type:        schema.TypeList,
			Optional:    !asDataSource,
			Computed:    asDataSource,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		serviceTypeFieldName: {
			Description: "The type of service to use for this cluster in a hybrid cloud environment.",
			Type:        schema.TypeString,
			Optional:    !asDataSource,
			Computed:    asDataSource,
		},
		serviceAnnotationsFieldName: {
			Description: "List of annotations applied to the service of this cluster in a hybrid cloud environment.",
			Type:        schema.TypeList,
			Optional:    !asDataSource,
			Computed:    asDataSource,
			Elem: &schema.Resource{
				Schema: keyValSchema(asDataSource),
			},
		},
		podLabelsFieldName: {
			Description: "List of labels applied to the pods of this cluster in a hybrid cloud environment.",
			Type:        schema.TypeList,
			Optional:    !asDataSource,
			Computed:    asDataSource,
			Elem: &schema.Resource{
				Schema: keyValSchema(asDataSource),
			},
		},
		databaseConfigurationFieldName: {
			Description: "Configuration for the Qdrant database engine, primarily for hybrid cloud setups.",
			Type:        schema.TypeList,
			Optional:    !asDataSource,
			Computed:    asDataSource,
			MaxItems:    maxItems,
			Elem: &schema.Resource{
				Schema: databaseConfigurationSchema(asDataSource),
			},
		},
		dbConfigReservedCpuPercentageFieldName: {
			Description: "The percentage of CPU resources reserved for system components.",
			Type:        schema.TypeInt,
			Optional:    !asDataSource,
			Computed:    asDataSource,
		},
		dbConfigReservedMemoryPercentageFieldName: {
			Description: "The percentage of RAM resources reserved for system components.",
			Type:        schema.TypeInt,
			Optional:    !asDataSource,
			Computed:    asDataSource,
		},
		dbConfigGpuTypeFieldName: {
			Description: "The GPU type that should be used for the database.",
			Type:        schema.TypeString,
			Optional:    !asDataSource,
			Computed:    asDataSource,
		},
		dbConfigRestartPolicyFieldName: {
			Description: "The restart policy for the database.",
			Type:        schema.TypeString,
			Optional:    !asDataSource,
			Computed:    asDataSource,
		},
		dbConfigRebalanceStrategyFieldName: {
			Description: "The automatic shard rebalancing strategy for the database.",
			Type:        schema.TypeString,
			Optional:    !asDataSource,
			Computed:    asDataSource,
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

// databaseConfigurationSchema defines the schema for the database configuration.
func databaseConfigurationSchema(asDataSource bool) map[string]*schema.Schema {
	maxItems := 1
	if asDataSource {
		maxItems = 0
	}
	return map[string]*schema.Schema{
		dbConfigCollectionFieldName: {
			Description: "Default collection parameters.",
			Type:        schema.TypeList,
			Optional:    !asDataSource,
			Computed:    asDataSource,
			MaxItems:    maxItems,
			Elem: &schema.Resource{
				Schema: databaseConfigurationCollectionSchema(asDataSource),
			},
		},
		dbConfigStorageFieldName: {
			Description: "Storage-related configuration.",
			Type:        schema.TypeList,
			Optional:    !asDataSource,
			Computed:    asDataSource,
			MaxItems:    maxItems,
			Elem: &schema.Resource{
				Schema: databaseConfigurationStorageSchema(asDataSource),
			},
		},
		dbConfigServiceFieldName: {
			Description: "Service-related configuration.",
			Type:        schema.TypeList,
			Optional:    !asDataSource,
			Computed:    asDataSource,
			MaxItems:    maxItems,
			Elem: &schema.Resource{
				Schema: databaseConfigurationServiceSchema(asDataSource),
			},
		},
		dbConfigLogLevelFieldName: {
			Description: "Logging level for the database.",
			Type:        schema.TypeString,
			Optional:    !asDataSource,
			Computed:    asDataSource,
		},
		dbConfigTlsFieldName: {
			Description: "TLS configuration for the database.",
			Type:        schema.TypeList,
			Optional:    !asDataSource,
			Computed:    asDataSource,
			MaxItems:    maxItems,
			Elem: &schema.Resource{
				Schema: databaseConfigurationTlsSchema(asDataSource),
			},
		},
		dbConfigInferenceFieldName: {
			Description: "Inference service configuration.",
			Type:        schema.TypeList,
			Optional:    !asDataSource,
			Computed:    asDataSource,
			MaxItems:    maxItems,
			Elem: &schema.Resource{
				Schema: databaseConfigurationInferenceSchema(asDataSource),
			},
		},
	}
}

// databaseConfigurationCollectionSchema defines the schema for collection configuration.
func databaseConfigurationCollectionSchema(asDataSource bool) map[string]*schema.Schema {
	maxItems := 1
	if asDataSource {
		maxItems = 0
	}
	return map[string]*schema.Schema{
		dbConfigCollectionReplicationFactor: {
			Type:     schema.TypeInt,
			Optional: !asDataSource,
			Computed: asDataSource,
		},
		dbConfigCollectionWriteConsistencyFactor: {
			Type:     schema.TypeInt,
			Optional: !asDataSource,
			Computed: asDataSource,
		},
		dbConfigCollectionVectorsFieldName: {
			Type:     schema.TypeList,
			Optional: !asDataSource,
			Computed: asDataSource,
			MaxItems: maxItems,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					dbConfigCollectionVectorsOnDiskFieldName: {
						Type:     schema.TypeBool,
						Optional: !asDataSource,
						Computed: asDataSource,
					},
				},
			},
		},
	}
}

// databaseConfigurationStorageSchema defines the schema for storage configuration.
func databaseConfigurationStorageSchema(asDataSource bool) map[string]*schema.Schema {
	maxItems := 1
	if asDataSource {
		maxItems = 0
	}
	return map[string]*schema.Schema{
		dbConfigStoragePerformanceFieldName: {
			Type:     schema.TypeList,
			Optional: !asDataSource,
			Computed: asDataSource,
			MaxItems: maxItems,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					dbConfigStoragePerfOptimizerCpuBudget: {
						Type:     schema.TypeInt,
						Optional: !asDataSource,
						Computed: asDataSource,
					},
					dbConfigStoragePerfAsyncScorer: {
						Type:     schema.TypeBool,
						Optional: !asDataSource,
						Computed: asDataSource,
					},
				},
			},
		},
	}
}

// databaseConfigurationServiceSchema defines the schema for service configuration.
func databaseConfigurationServiceSchema(asDataSource bool) map[string]*schema.Schema {
	maxItems := 1
	if asDataSource {
		maxItems = 0
	}
	return map[string]*schema.Schema{
		dbConfigServiceApiKeyFieldName: {
			Description: "Secret to use for the main API key.",
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			MaxItems:    maxItems,
			Elem:        secretKeyRefSchema(asDataSource),
		},
		dbConfigServiceReadOnlyApiKeyFieldName: {
			Description: "Secret to use for the read-only API key.",
			Type:        schema.TypeList,
			Optional:    !asDataSource,
			Computed:    asDataSource,
			MaxItems:    maxItems,
			Elem:        secretKeyRefSchema(asDataSource),
		},
		dbConfigServiceJwtRbacFieldName: {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		dbConfigServiceEnableTlsFieldName: {
			Type:     schema.TypeBool,
			Optional: !asDataSource,
			Computed: asDataSource,
		},
	}
}

// databaseConfigurationTlsSchema defines the schema for TLS configuration.
func databaseConfigurationTlsSchema(asDataSource bool) map[string]*schema.Schema {
	maxItems := 1
	if asDataSource {
		maxItems = 0
	}
	return map[string]*schema.Schema{
		dbConfigTlsCertFieldName: {
			Description: "Secret to use for the certificate.",
			Type:        schema.TypeList,
			Optional:    !asDataSource,
			Computed:    asDataSource,
			MaxItems:    maxItems,
			Elem:        secretKeyRefSchema(asDataSource),
		},
		dbConfigTlsKeyFieldName: {
			Description: "Secret to use for the private key.",
			Type:        schema.TypeList,
			Optional:    !asDataSource,
			Computed:    asDataSource,
			MaxItems:    maxItems,
			Elem:        secretKeyRefSchema(asDataSource),
		},
	}
}

// databaseConfigurationInferenceSchema defines the schema for inference configuration.
func databaseConfigurationInferenceSchema(asDataSource bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		dbConfigInferenceEnabledFieldName: {
			Type:     schema.TypeBool,
			Optional: !asDataSource,
			Computed: asDataSource,
		},
	}
}

// secretKeyRefSchema defines the schema for a secret key reference.
func secretKeyRefSchema(asDataSource bool) *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			dbConfigSecretKeyRefSecretNameFieldName: {
				Description: "The name of the secret.",
				Type:        schema.TypeString,
				Required:    !asDataSource,
				Computed:    asDataSource,
			},
			dbConfigSecretKeyRefSecretKeyFieldName: {
				Description: "The key within the secret.",
				Type:        schema.TypeString,
				Required:    !asDataSource,
				Computed:    asDataSource,
			},
		},
	}
}

// keyValSchema defines the schema for a key-value pair.
func keyValSchema(asDataSource bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"key": {
			Type:     schema.TypeString,
			Required: !asDataSource,
			Computed: asDataSource,
		},
		"value": {
			Type:     schema.TypeString,
			Required: !asDataSource,
			Computed: asDataSource,
		},
	}
}

// tolerationSchema defines the schema for a toleration.
func tolerationSchema(asDataSource bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		tolerationKeyFieldName: {
			Type:     schema.TypeString,
			Optional: !asDataSource,
			Computed: asDataSource,
		},
		tolerationOperatorFieldName: {
			Type:     schema.TypeString,
			Optional: !asDataSource,
			Computed: asDataSource,
		},
		tolerationValueFieldName: {
			Type:     schema.TypeString,
			Optional: !asDataSource,
			Computed: asDataSource,
		},
		tolerationEffectFieldName: {
			Type:     schema.TypeString,
			Optional: !asDataSource,
			Computed: asDataSource,
		},
		tolerationSecondsFieldName: {
			Type:     schema.TypeInt,
			Optional: !asDataSource,
			Computed: asDataSource,
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
	if v, ok := d.GetOk(clusterLabelsFieldName); ok {
		cluster.Labels = expandKeyVal(v.([]interface{}))
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
		if v, ok := item[clusterLastModifiedAtFieldName]; ok {
			config.LastModifiedAt = parseTime(v.(string))
		}
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
		if v, ok := item[databaseConfigurationFieldName]; ok {
			config.DatabaseConfiguration = expandDatabaseConfiguration(v.([]interface{}))
		}
		if v, ok := item[nodeSelectorFieldName]; ok {
			config.NodeSelector = expandKeyVal(v.([]interface{}))
		}
		if v, ok := item[tolerationsFieldName]; ok {
			config.Tolerations = expandTolerations(v.([]interface{}))
		}
		if v, ok := item[annotationsFieldName]; ok {
			config.Annotations = expandKeyVal(v.([]interface{}))
		}
		if v, ok := item[allowedIpSourceRangesFieldName]; ok {
			if ranges, ok := v.([]interface{}); ok && len(ranges) > 0 {
				config.AllowedIpSourceRanges = make([]string, len(ranges))
				for i, r := range ranges {
					config.AllowedIpSourceRanges[i] = r.(string)
				}
			}
		}
		if v, ok := item[serviceTypeFieldName]; ok {
			st, stOK := qcCluster.ClusterServiceType_value[v.(string)]
			if stOK {
				config.ServiceType = newPointer(qcCluster.ClusterServiceType(st))
			}
		}
		if v, ok := item[serviceAnnotationsFieldName]; ok {
			config.ServiceAnnotations = expandKeyVal(v.([]interface{}))
		}
		if v, ok := item[podLabelsFieldName]; ok {
			config.PodLabels = expandKeyVal(v.([]interface{}))
		}
		if v, ok := item[dbConfigReservedCpuPercentageFieldName]; ok {
			if percent := v.(int); percent != 0 {
				config.ReservedCpuPercentage = newPointer(uint32(percent))
			}
		}
		if v, ok := item[dbConfigReservedMemoryPercentageFieldName]; ok {
			if percent := v.(int); percent != 0 {
				config.ReservedMemoryPercentage = newPointer(uint32(percent))
			}
		}
		if v, ok := item[dbConfigGpuTypeFieldName]; ok {
			gt, gtOK := qcCluster.ClusterConfigurationGpuType_value[v.(string)]
			if gtOK {
				config.GpuType = newPointer(qcCluster.ClusterConfigurationGpuType(gt))
			}
		}
		if v, ok := item[dbConfigRestartPolicyFieldName]; ok {
			rp, rpOK := qcCluster.ClusterConfigurationRestartPolicy_value[v.(string)]
			if rpOK {
				config.RestartPolicy = newPointer(qcCluster.ClusterConfigurationRestartPolicy(rp))
			}
		}
		if v, ok := item[dbConfigRebalanceStrategyFieldName]; ok {
			rs, rsOK := qcCluster.ClusterConfigurationRebalanceStrategy_value[v.(string)]
			if rsOK {
				config.RebalanceStrategy = newPointer(qcCluster.ClusterConfigurationRebalanceStrategy(rs))
			}
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

// expandKeyVal expands a key-value pair from Terraform data.
func expandKeyVal(v []interface{}) []*commonv1.KeyValue {
	var result []*commonv1.KeyValue
	for _, m := range v {
		item := m.(map[string]interface{})
		result = append(result, &commonv1.KeyValue{
			Key:   item["key"].(string),
			Value: item["value"].(string),
		})
	}
	return result
}

// expandTolerations expands tolerations from Terraform data.
func expandTolerations(v []interface{}) []*qcCluster.Toleration {
	var result []*qcCluster.Toleration
	for _, m := range v {
		item := m.(map[string]interface{})
		toleration := &qcCluster.Toleration{
			Key:   item[tolerationKeyFieldName].(string),
			Value: item[tolerationValueFieldName].(string),
		}
		if op, ok := item[tolerationOperatorFieldName]; ok {
			opVal, opOk := qcCluster.TolerationOperator_value[op.(string)]
			if opOk {
				toleration.Operator = newPointer(qcCluster.TolerationOperator(opVal))
			}
		}
		if eff, ok := item[tolerationEffectFieldName]; ok {
			effVal, effOk := qcCluster.TolerationEffect_value[eff.(string)]
			if effOk {
				toleration.Effect = newPointer(qcCluster.TolerationEffect(effVal))
			}
		}
		if sec, ok := item[tolerationSecondsFieldName]; ok {
			if secInt := sec.(int); secInt > 0 {
				toleration.TolerationSeconds = newPointer(uint64(secInt))
			}
		}
		result = append(result, toleration)
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
		clusterLabelsFieldName:              flattenKeyVal(cluster.GetLabels()),
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
	if clusterConfig == nil {
		return []interface{}{}
	}
	config := map[string]interface{}{
		clusterVersionFieldName:            clusterConfig.GetVersion(),
		clusterLastModifiedAtFieldName:     formatTime(clusterConfig.GetLastModifiedAt()),
		numberOfNodesFieldName:             int(clusterConfig.GetNumberOfNodes()),
		nodeConfigurationFieldName:         flattenNodeConfiguration(clusterConfig.GetPackageId(), clusterConfig.GetAdditionalResources()),
		databaseConfigurationFieldName:     flattenDatabaseConfiguration(clusterConfig.GetDatabaseConfiguration()),
		nodeSelectorFieldName:              flattenKeyVal(clusterConfig.GetNodeSelector()),
		tolerationsFieldName:               flattenTolerations(clusterConfig.GetTolerations()),
		annotationsFieldName:               flattenKeyVal(clusterConfig.GetAnnotations()),
		allowedIpSourceRangesFieldName:     clusterConfig.GetAllowedIpSourceRanges(),
		serviceTypeFieldName:               clusterConfig.GetServiceType().String(),
		serviceAnnotationsFieldName:        flattenKeyVal(clusterConfig.GetServiceAnnotations()),
		podLabelsFieldName:                 flattenKeyVal(clusterConfig.GetPodLabels()),
		dbConfigGpuTypeFieldName:           clusterConfig.GetGpuType().String(),
		dbConfigRestartPolicyFieldName:     clusterConfig.GetRestartPolicy().String(),
		dbConfigRebalanceStrategyFieldName: clusterConfig.GetRebalanceStrategy().String(),
	}

	if clusterConfig.ReservedCpuPercentage != nil {
		config[dbConfigReservedCpuPercentageFieldName] = int(clusterConfig.GetReservedCpuPercentage())
	}
	if clusterConfig.ReservedMemoryPercentage != nil {
		config[dbConfigReservedMemoryPercentageFieldName] = int(clusterConfig.GetReservedMemoryPercentage())
	}

	return []interface{}{config}
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

// flattenKeyVal flattens a key-value pair for storage in Terraform.
func flattenKeyVal(kvs []*commonv1.KeyValue) []interface{} {
	var result []interface{}
	for _, kv := range kvs {
		result = append(result, map[string]interface{}{
			"key":   kv.GetKey(),
			"value": kv.GetValue(),
		})
	}
	return result
}

// flattenTolerations flattens tolerations for storage in Terraform.
func flattenTolerations(tolerations []*qcCluster.Toleration) []interface{} {
	var result []interface{}
	for _, t := range tolerations {
		tolerationMap := map[string]interface{}{
			tolerationKeyFieldName:   t.GetKey(),
			tolerationValueFieldName: t.GetValue(),
		}
		if t.Operator != nil {
			tolerationMap[tolerationOperatorFieldName] = t.GetOperator().String()
		}
		if t.Effect != nil {
			tolerationMap[tolerationEffectFieldName] = t.GetEffect().String()
		}
		if t.TolerationSeconds != nil {
			tolerationMap[tolerationSecondsFieldName] = int(t.GetTolerationSeconds())
		}

		result = append(result, tolerationMap)
	}
	return result
}

// expandDatabaseConfiguration expands the Terraform resource data into a database configuration object.
func expandDatabaseConfiguration(v []interface{}) *qcCluster.DatabaseConfiguration {
	if len(v) == 0 || v[0] == nil {
		return nil
	}
	item := v[0].(map[string]interface{})
	config := &qcCluster.DatabaseConfiguration{}

	if val, ok := item[dbConfigCollectionFieldName]; ok && len(val.([]interface{})) > 0 {
		collItem := val.([]interface{})[0].(map[string]interface{})
		collConfig := &qcCluster.DatabaseConfigurationCollection{}
		if v, ok := collItem[dbConfigCollectionReplicationFactor]; ok {
			val := uint32(v.(int))
			if val > 0 {
				collConfig.ReplicationFactor = &val
			}
		}
		if v, ok := collItem[dbConfigCollectionWriteConsistencyFactor]; ok {
			val := int32(v.(int))
			if val > 0 {
				collConfig.WriteConsistencyFactor = &val
			}
		}
		if v, ok := collItem[dbConfigCollectionVectorsFieldName]; ok && len(v.([]interface{})) > 0 {
			vecItem := v.([]interface{})[0].(map[string]interface{})
			if onDisk, ok := vecItem[dbConfigCollectionVectorsOnDiskFieldName]; ok {
				collConfig.Vectors = &qcCluster.DatabaseConfigurationCollectionVectors{
					OnDisk: newPointer(onDisk.(bool)),
				}
			}
		}
		config.Collection = collConfig
	}

	if val, ok := item[dbConfigStorageFieldName]; ok && len(val.([]interface{})) > 0 {
		storageItem := val.([]interface{})[0].(map[string]interface{})
		storageConfig := &qcCluster.DatabaseConfigurationStorage{}
		if v, ok := storageItem[dbConfigStoragePerformanceFieldName]; ok && len(v.([]interface{})) > 0 {
			perfItem := v.([]interface{})[0].(map[string]interface{})
			perfConfig := &qcCluster.DatabaseConfigurationStoragePerformance{}
			if budget, ok := perfItem[dbConfigStoragePerfOptimizerCpuBudget]; ok {
				perfConfig.OptimizerCpuBudget = int32(budget.(int))
			}
			if scorer, ok := perfItem[dbConfigStoragePerfAsyncScorer]; ok {
				perfConfig.AsyncScorer = scorer.(bool)
			}
			storageConfig.Performance = perfConfig
		}
		config.Storage = storageConfig
	}

	if val, ok := item[dbConfigServiceFieldName]; ok && len(val.([]interface{})) > 0 {
		serviceItem := val.([]interface{})[0].(map[string]interface{})
		serviceConfig := &qcCluster.DatabaseConfigurationService{}
		if v, ok := serviceItem[dbConfigServiceApiKeyFieldName]; ok {
			serviceConfig.ApiKey = expandSecretKeyRef(v.([]interface{}))
		}
		if v, ok := serviceItem[dbConfigServiceReadOnlyApiKeyFieldName]; ok {
			serviceConfig.ReadOnlyApiKey = expandSecretKeyRef(v.([]interface{}))
		}
		if v, ok := serviceItem[dbConfigServiceJwtRbacFieldName]; ok {
			serviceConfig.JwtRbac = newPointer(v.(bool))
		}
		if v, ok := serviceItem[dbConfigServiceEnableTlsFieldName]; ok {
			serviceConfig.EnableTls = v.(bool)
		}
		config.Service = serviceConfig
	}

	if val, ok := item[dbConfigLogLevelFieldName]; ok {
		strVal := val.(string)
		if strVal != "" {
			config.LogLevel = newPointer(qcCluster.DatabaseConfigurationLogLevel(qcCluster.DatabaseConfigurationLogLevel_value[strVal]))
		}
	}

	// TLS and Inference are simple boolean flags in their respective objects
	if val, ok := item[dbConfigTlsFieldName]; ok && len(val.([]interface{})) > 0 {
		tlsItem := val.([]interface{})[0].(map[string]interface{})
		tlsConfig := &qcCluster.DatabaseConfigurationTls{}
		if v, ok := tlsItem[dbConfigTlsCertFieldName]; ok {
			tlsConfig.Cert = expandSecretKeyRef(v.([]interface{}))
		}
		if v, ok := tlsItem[dbConfigTlsKeyFieldName]; ok {
			tlsConfig.Key = expandSecretKeyRef(v.([]interface{}))
		}
		config.Tls = tlsConfig
	}
	if val, ok := item[dbConfigInferenceFieldName]; ok && len(val.([]interface{})) > 0 {
		infItem := val.([]interface{})[0].(map[string]interface{})
		if enabled, ok := infItem[dbConfigInferenceEnabledFieldName]; ok {
			config.Inference = &qcCluster.DatabaseConfigurationInference{Enabled: enabled.(bool)}
		}
	}

	return config
}

// expandSecretKeyRef expands a secret key reference from Terraform data.
func expandSecretKeyRef(v []interface{}) *commonv1.SecretKeyRef {
	if len(v) == 0 || v[0] == nil {
		return nil
	}
	item := v[0].(map[string]interface{})
	return &commonv1.SecretKeyRef{
		Name: item[dbConfigSecretKeyRefSecretNameFieldName].(string),
		Key:  item[dbConfigSecretKeyRefSecretKeyFieldName].(string),
	}
}

// flattenDatabaseConfiguration flattens the database configuration for storage in Terraform.
func flattenDatabaseConfiguration(config *qcCluster.DatabaseConfiguration) []interface{} {
	if config == nil {
		return []interface{}{}
	}

	m := make(map[string]interface{})

	if coll := config.GetCollection(); coll != nil {
		collMap := map[string]interface{}{}
		if coll.ReplicationFactor != nil {
			collMap[dbConfigCollectionReplicationFactor] = int(coll.GetReplicationFactor())
		}
		if coll.WriteConsistencyFactor != nil {
			collMap[dbConfigCollectionWriteConsistencyFactor] = int(coll.GetWriteConsistencyFactor())
		}
		if vec := coll.GetVectors(); vec != nil && vec.OnDisk != nil {
			collMap[dbConfigCollectionVectorsFieldName] = []interface{}{
				map[string]interface{}{
					dbConfigCollectionVectorsOnDiskFieldName: vec.GetOnDisk(),
				},
			}
		}
		m[dbConfigCollectionFieldName] = []interface{}{collMap}
	}

	if service := config.GetService(); service != nil {
		serviceMap := map[string]interface{}{
			dbConfigServiceApiKeyFieldName:         flattenSecretKeyRef(service.GetApiKey()),
			dbConfigServiceReadOnlyApiKeyFieldName: flattenSecretKeyRef(service.GetReadOnlyApiKey()),
			dbConfigServiceJwtRbacFieldName:        service.GetJwtRbac(),
			dbConfigServiceEnableTlsFieldName:      service.GetEnableTls(),
		}
		m[dbConfigServiceFieldName] = []interface{}{serviceMap}
	}

	if tls := config.GetTls(); tls != nil {
		m[dbConfigTlsFieldName] = []interface{}{
			map[string]interface{}{
				dbConfigTlsCertFieldName: flattenSecretKeyRef(tls.GetCert()),
				dbConfigTlsKeyFieldName:  flattenSecretKeyRef(tls.GetKey()),
			},
		}
	}

	if config.LogLevel != nil {
		m[dbConfigLogLevelFieldName] = config.GetLogLevel().String()
	}

	return []interface{}{m}
}

// flattenSecretKeyRef flattens a secret key reference for storage in Terraform.
func flattenSecretKeyRef(ref *commonv1.SecretKeyRef) []interface{} {
	if ref == nil {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			dbConfigSecretKeyRefSecretNameFieldName: ref.GetName(),
			dbConfigSecretKeyRefSecretKeyFieldName:  ref.GetKey(),
		},
	}
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
