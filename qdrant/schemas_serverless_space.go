package qdrant

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	spacev1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/serverless/space/v1"
)

const (
	slSpaceFieldTemplate = "Serverless Space Schema %s field"

	slSpaceIDFieldName                  = "id"
	slSpaceAccountIDFieldName           = "account_id"
	slSpaceNameFieldName                = "name"
	slSpaceCloudRegionIDFieldName       = "cloud_region_id"
	slSpaceLabelsFieldName              = "labels"
	slSpaceCostAllocationLabelFieldName = "cost_allocation_label"
	slSpaceConfigurationFieldName       = "configuration"
	slSpaceFromBackupIDFieldName        = "from_backup_id"
	slSpaceDeleteBackupsFieldName       = "delete_backups"
	slSpaceCreatedAtFieldName           = "created_at"
	slSpaceStateFieldName               = "state"
	slSpaceURLFieldName                 = "url"
	slSpaceRestPortFieldName            = "rest_port"
	slSpaceGrpcPortFieldName            = "grpc_port"
	slSpacesFieldName                   = "spaces"

	slSpaceCfgLastModifiedAtFieldName         = "last_modified_at"
	slSpaceCfgAllowedIPSourceRangesFieldName  = "allowed_ip_source_ranges"
	slSpaceCfgAllowedOriginsFieldName         = "allowed_origins"
	slSpaceCfgMaxCollectionsPerSpaceFieldName = "max_collections_per_space"
	slSpaceCfgCollectionSettingsFieldName     = "collection_settings"
	slSpaceCfgSearcherSettingsFieldName       = "searcher_settings"

	slCollectionSettingsPlatformMaxSizeFieldName = "platform_max_size"
	slCollectionSettingsMaxSizeFieldName         = "max_size"

	slSearcherSettingsIdleTimeoutFieldName        = "idle_timeout"
	slSearcherSettingsPlatformMaxWorkersFieldName = "platform_max_workers"
	slSearcherSettingsMaxWorkersFieldName         = "max_workers"

	slStatePhaseFieldName  = "phase"
	slStateReasonFieldName = "reason"
)

var slSpaceNameRegexp = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)

// serverlessSpaceSchema returns the schema of a serverless space, used by both the resource and data sources.
func serverlessSpaceSchema(asDataSource bool) map[string]*schema.Schema {
	maxItems := configurableMaxItems(1, asDataSource)
	s := map[string]*schema.Schema{
		slSpaceIDFieldName: {
			Description: fmt.Sprintf(slSpaceFieldTemplate, "Space ID"),
			Type:        schema.TypeString,
			Required:    asDataSource,
			Computed:    !asDataSource,
		},
		slSpaceAccountIDFieldName: {
			Description: fmt.Sprintf(slSpaceFieldTemplate, "Account ID (defaults to the provider account ID)"),
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			ForceNew:    !asDataSource,
		},
		slSpaceNameFieldName: {
			Description: fmt.Sprintf(slSpaceFieldTemplate, "Name of the space (4-64 characters: letters, numbers, `-` and `_`). Unique within the account"),
			Type:        schema.TypeString,
			Required:    !asDataSource,
			Computed:    asDataSource,
		},
		slSpaceCloudRegionIDFieldName: {
			Description: fmt.Sprintf(slSpaceFieldTemplate, "Cloud region where the space is hosted (see the `qdrant-cloud_serverless_cloud_regions` data source). Required unless `from_backup_id` is set, in which case the region of the backup is used"),
			Type:        schema.TypeString,
			Optional:    !asDataSource,
			Computed:    true,
			ForceNew:    !asDataSource,
		},
		slSpaceLabelsFieldName: {
			Description: fmt.Sprintf(slSpaceFieldTemplate, "Labels of the space, used in the cloud UI and billing reports (max 10)"),
			Type:        schema.TypeSet,
			Optional:    !asDataSource,
			Computed:    true,
			MaxItems:    configurableMaxItems(10, asDataSource),
			Elem:        &schema.Resource{Schema: keyValSchema(asDataSource)},
			Set:         keyValHashFunc,
		},
		slSpaceCostAllocationLabelFieldName: {
			Description: fmt.Sprintf(slSpaceFieldTemplate, "Key of the label sent as metadata when metering usage through cloud marketplaces, for cost allocation"),
			Type:        schema.TypeString,
			Optional:    !asDataSource,
			Computed:    true,
		},
		slSpaceConfigurationFieldName: {
			Description: fmt.Sprintf(slSpaceFieldTemplate, "Configuration of the space"),
			Type:        schema.TypeList,
			Optional:    !asDataSource,
			Computed:    true,
			MaxItems:    maxItems,
			Elem:        &schema.Resource{Schema: serverlessSpaceConfigurationSchema(asDataSource)},
		},
		slSpaceCreatedAtFieldName: {
			Description: fmt.Sprintf(slSpaceFieldTemplate, "Creation time"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		slSpaceStateFieldName: {
			Description: fmt.Sprintf(slSpaceFieldTemplate, "Current state of the space"),
			Type:        schema.TypeList,
			Computed:    true,
			Elem:        &schema.Resource{Schema: serverlessStateSchema()},
		},
		slSpaceURLFieldName: {
			Description: fmt.Sprintf(slSpaceFieldTemplate, "URL (without port) to access the space"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		slSpaceRestPortFieldName: {
			Description: fmt.Sprintf(slSpaceFieldTemplate, "Port for HTTP REST calls"),
			Type:        schema.TypeInt,
			Computed:    true,
		},
		slSpaceGrpcPortFieldName: {
			Description: fmt.Sprintf(slSpaceFieldTemplate, "Port for gRPC calls"),
			Type:        schema.TypeInt,
			Computed:    true,
		},
	}
	if !asDataSource {
		s[slSpaceNameFieldName].ValidateDiagFunc = validation.ToDiagFunc(validation.All(
			validation.StringLenBetween(4, 64),
			validation.StringMatch(slSpaceNameRegexp, "must only contain letters, numbers, '-' and '_'"),
		))
		s[slSpaceFromBackupIDFieldName] = &schema.Schema{
			Description: "ID of a serverless backup to create the space from. Changing this forces a new space to be created.",
			Type:        schema.TypeString,
			Optional:    true,
			// Not returned by the API, the value from the configuration is kept in the state.
			Computed:         true,
			ForceNew:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
		}
		s[slSpaceDeleteBackupsFieldName] = &schema.Schema{
			Description: "Whether the backups of the space should be deleted when the space is destroyed. Defaults to false.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		}
		s[slSpaceCloudRegionIDFieldName].AtLeastOneOf = []string{slSpaceCloudRegionIDFieldName, slSpaceFromBackupIDFieldName}
	}
	return s
}

func serverlessSpaceConfigurationSchema(asDataSource bool) map[string]*schema.Schema {
	maxItems := configurableMaxItems(1, asDataSource)
	s := map[string]*schema.Schema{
		slSpaceCfgLastModifiedAtFieldName: {
			Description: "Time when the configuration was last updated.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		slSpaceCfgAllowedIPSourceRangesFieldName: {
			Description: "List of allowed IPv4 source ranges (CIDR) for this space (max 40).",
			Type:        schema.TypeList,
			Optional:    !asDataSource,
			Computed:    asDataSource,
			MaxItems:    configurableMaxItems(40, asDataSource),
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		slSpaceCfgAllowedOriginsFieldName: {
			Description: "List of origins (scheme + host + optional port, e.g. `https://app.example.com`) from which browser-based clients are allowed to call the space via CORS (max 10).",
			Type:        schema.TypeList,
			Optional:    !asDataSource,
			Computed:    asDataSource,
			MaxItems:    configurableMaxItems(10, asDataSource),
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		slSpaceCfgMaxCollectionsPerSpaceFieldName: {
			Description: "Platform-enforced limit on the number of collections for this space (0 means unlimited).",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		slSpaceCfgCollectionSettingsFieldName: {
			Description: "Per-collection size limits for this space.",
			Type:        schema.TypeList,
			Optional:    !asDataSource,
			Computed:    true,
			MaxItems:    maxItems,
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				slCollectionSettingsPlatformMaxSizeFieldName: {
					Description: "Platform-enforced maximum size in bytes per collection (0 means unlimited).",
					Type:        schema.TypeInt,
					Computed:    true,
				},
				slCollectionSettingsMaxSizeFieldName: {
					Description: "Customer-defined maximum size in bytes per collection, used for cost control. Must not exceed `platform_max_size`.",
					Type:        schema.TypeInt,
					Optional:    !asDataSource,
					Computed:    true,
				},
			}},
		},
		slSpaceCfgSearcherSettingsFieldName: {
			Description: "Search-worker settings for this space.",
			Type:        schema.TypeList,
			Optional:    !asDataSource,
			Computed:    true,
			MaxItems:    maxItems,
			Elem: &schema.Resource{Schema: map[string]*schema.Schema{
				slSearcherSettingsIdleTimeoutFieldName: {
					Description: `Idle timeout after which idle search workers may be scaled down, as a Go duration string (between "1m" and "15m", defaults to "5m").`,
					Type:        schema.TypeString,
					Optional:    !asDataSource,
					Computed:    true,
				},
				slSearcherSettingsPlatformMaxWorkersFieldName: {
					Description: "Platform-enforced maximum number of search workers per collection (0 means unlimited).",
					Type:        schema.TypeInt,
					Computed:    true,
				},
				slSearcherSettingsMaxWorkersFieldName: {
					Description: "Maximum number of search workers per collection in this space (defaults to 2).",
					Type:        schema.TypeInt,
					Optional:    !asDataSource,
					Computed:    true,
				},
			}},
		},
	}
	if !asDataSource {
		s[slSpaceCfgAllowedIPSourceRangesFieldName].Elem = &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
		}
		searcher := s[slSpaceCfgSearcherSettingsFieldName].Elem.(*schema.Resource).Schema
		searcher[slSearcherSettingsIdleTimeoutFieldName].DiffSuppressFunc = suppressDurationDiff
		searcher[slSearcherSettingsMaxWorkersFieldName].ValidateDiagFunc = validation.ToDiagFunc(validation.IntAtLeast(1))
		collection := s[slSpaceCfgCollectionSettingsFieldName].Elem.(*schema.Resource).Schema
		collection[slCollectionSettingsMaxSizeFieldName].ValidateDiagFunc = validation.ToDiagFunc(validation.IntAtLeast(0))
	}
	return s
}

// serverlessStateSchema is the (computed) state of a serverless object (phase + reason).
func serverlessStateSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		slStatePhaseFieldName: {
			Description: "Current phase.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		slStateReasonFieldName: {
			Description: "Description of any errors or issues (empty when operating normally).",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func serverlessSpacesDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		slSpaceAccountIDFieldName: {
			Description: fmt.Sprintf(slSpaceFieldTemplate, "Account ID (defaults to the provider account ID)"),
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
		},
		slSpaceCloudRegionIDFieldName: {
			Description: fmt.Sprintf(slSpaceFieldTemplate, "Only list spaces in this cloud region"),
			Type:        schema.TypeString,
			Optional:    true,
		},
		slSpacesFieldName: {
			Description: "List of serverless spaces",
			Type:        schema.TypeList,
			Computed:    true,
			Elem:        &schema.Resource{Schema: serverlessSpaceSchema(true)},
		},
	}
}

// flattenServerlessSpace flattens the space into a map, the keys match the space schema.
func flattenServerlessSpace(space *spacev1.Space) map[string]interface{} {
	endpoint := space.GetState().GetEndpoint()
	return map[string]interface{}{
		slSpaceIDFieldName:                  space.GetId(),
		slSpaceAccountIDFieldName:           space.GetAccountId(),
		slSpaceNameFieldName:                space.GetName(),
		slSpaceCloudRegionIDFieldName:       space.GetCloudRegionId(),
		slSpaceLabelsFieldName:              flattenKeyVal(space.GetLabels()),
		slSpaceCostAllocationLabelFieldName: space.GetCostAllocationLabel(),
		slSpaceConfigurationFieldName:       flattenServerlessSpaceConfiguration(space.GetConfiguration()),
		slSpaceCreatedAtFieldName:           formatTime(space.GetCreatedAt()),
		slSpaceStateFieldName:               flattenServerlessState(flattenProtoEnum(space.GetState().GetPhase()), space.GetState().GetReason()),
		slSpaceURLFieldName:                 endpoint.GetUrl(),
		slSpaceRestPortFieldName:            int(endpoint.GetRestPort()),
		slSpaceGrpcPortFieldName:            int(endpoint.GetGrpcPort()),
	}
}

func flattenServerlessSpaces(spaces []*spacev1.Space) []interface{} {
	result := make([]interface{}, 0, len(spaces))
	for _, space := range spaces {
		result = append(result, flattenServerlessSpace(space))
	}
	return result
}

func flattenServerlessSpaceConfiguration(cfg *spacev1.SpaceConfiguration) []interface{} {
	if cfg == nil {
		return []interface{}{}
	}
	return []interface{}{map[string]interface{}{
		slSpaceCfgLastModifiedAtFieldName:         formatTime(cfg.GetLastModifiedAt()),
		slSpaceCfgAllowedIPSourceRangesFieldName:  cfg.GetAllowedIpSourceRanges(),
		slSpaceCfgAllowedOriginsFieldName:         cfg.GetAllowedOrigins(),
		slSpaceCfgMaxCollectionsPerSpaceFieldName: int(cfg.GetMaxCollectionsPerSpace()),
		slSpaceCfgCollectionSettingsFieldName: []interface{}{map[string]interface{}{
			slCollectionSettingsPlatformMaxSizeFieldName: int(cfg.GetCollectionSettings().GetPlatformMaxSize()),
			slCollectionSettingsMaxSizeFieldName:         int(cfg.GetCollectionSettings().GetMaxSize()),
		}},
		slSpaceCfgSearcherSettingsFieldName: []interface{}{map[string]interface{}{
			slSearcherSettingsIdleTimeoutFieldName:        formatDuration(cfg.GetSearcherSettings().GetIdleTimeout()),
			slSearcherSettingsPlatformMaxWorkersFieldName: int(cfg.GetSearcherSettings().GetPlatformMaxWorkers()),
			slSearcherSettingsMaxWorkersFieldName:         int(cfg.GetSearcherSettings().GetMaxWorkers()),
		}},
	}}
}

func flattenServerlessState(phase, reason string) []interface{} {
	return []interface{}{map[string]interface{}{
		slStatePhaseFieldName:  phase,
		slStateReasonFieldName: reason,
	}}
}

// expandServerlessSpace builds the space from the resource data (without ID).
func expandServerlessSpace(d *schema.ResourceData, accountID string) *spacev1.Space {
	space := &spacev1.Space{
		AccountId:     accountID,
		Name:          d.Get(slSpaceNameFieldName).(string),
		CloudRegionId: d.Get(slSpaceCloudRegionIDFieldName).(string),
		Labels:        expandKeyVal(getInterfaceSliceFromSchemaValue(d.Get(slSpaceLabelsFieldName))),
		Configuration: expandServerlessSpaceConfiguration(d.Get(slSpaceConfigurationFieldName).([]interface{})),
	}
	if v := d.Get(slSpaceCostAllocationLabelFieldName).(string); v != "" {
		space.CostAllocationLabel = newPointer(v)
	}
	return space
}

func expandServerlessSpaceConfiguration(v []interface{}) *spacev1.SpaceConfiguration {
	cfg := &spacev1.SpaceConfiguration{}
	if len(v) == 0 || v[0] == nil {
		return cfg
	}
	item := v[0].(map[string]interface{})
	cfg.AllowedIpSourceRanges = setToStringSlice(item[slSpaceCfgAllowedIPSourceRangesFieldName])
	cfg.AllowedOrigins = setToStringSlice(item[slSpaceCfgAllowedOriginsFieldName])
	if cs, ok := item[slSpaceCfgCollectionSettingsFieldName].([]interface{}); ok && len(cs) > 0 && cs[0] != nil {
		csItem := cs[0].(map[string]interface{})
		cfg.CollectionSettings = &spacev1.CollectionSettings{}
		if maxSize, ok := csItem[slCollectionSettingsMaxSizeFieldName].(int); ok && maxSize > 0 {
			cfg.CollectionSettings.MaxSize = newPointer(uint64(maxSize))
		}
	}
	if ss, ok := item[slSpaceCfgSearcherSettingsFieldName].([]interface{}); ok && len(ss) > 0 && ss[0] != nil {
		ssItem := ss[0].(map[string]interface{})
		cfg.SearcherSettings = &spacev1.SearcherSettings{}
		if idle, ok := ssItem[slSearcherSettingsIdleTimeoutFieldName].(string); ok && idle != "" {
			cfg.SearcherSettings.IdleTimeout = parseDuration(idle)
		}
		if maxWorkers, ok := ssItem[slSearcherSettingsMaxWorkersFieldName].(int); ok && maxWorkers > 0 {
			cfg.SearcherSettings.MaxWorkers = newPointer(uint64(maxWorkers))
		}
	}
	return cfg
}

// serverlessSpaceUpdateMaskPaths maps Terraform attributes to the Space proto field paths used in the update mask.
var serverlessSpaceUpdateMaskPaths = []struct {
	attribute string
	path      string
}{
	{slSpaceNameFieldName, "name"},
	{slSpaceLabelsFieldName, "labels"},
	{slSpaceCostAllocationLabelFieldName, "cost_allocation_label"},
	{slSpaceConfigurationFieldName + ".0." + slSpaceCfgAllowedIPSourceRangesFieldName, "configuration.allowed_ip_source_ranges"},
	{slSpaceConfigurationFieldName + ".0." + slSpaceCfgAllowedOriginsFieldName, "configuration.allowed_origins"},
	{slSpaceConfigurationFieldName + ".0." + slSpaceCfgCollectionSettingsFieldName + ".0." + slCollectionSettingsMaxSizeFieldName, "configuration.collection_settings.max_size"},
	{slSpaceConfigurationFieldName + ".0." + slSpaceCfgSearcherSettingsFieldName + ".0." + slSearcherSettingsIdleTimeoutFieldName, "configuration.searcher_settings.idle_timeout"},
	{slSpaceConfigurationFieldName + ".0." + slSpaceCfgSearcherSettingsFieldName + ".0." + slSearcherSettingsMaxWorkersFieldName, "configuration.searcher_settings.max_workers"},
}

// serverlessSpaceUpdateMask returns the update mask paths for all attributes for which hasChange returns true.
func serverlessSpaceUpdateMask(hasChange func(string) bool) []string {
	var paths []string
	for _, p := range serverlessSpaceUpdateMaskPaths {
		if hasChange(p.attribute) {
			paths = append(paths, p.path)
		}
	}
	return paths
}
