package qdrant

import (
	"fmt"
	"regexp"
	"strings"

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
	slSpacesFieldName                   = "spaces"

	slSpaceCfgLastModifiedAtFieldName         = "last_modified_at"
	slSpaceCfgAllowedIPSourceRangesFieldName  = "allowed_ip_source_ranges"
	slSpaceCfgAllowedOriginsFieldName         = "allowed_origins"
	slSpaceCfgMaxCollectionsPerSpaceFieldName = "max_collections_per_space"
	slSpaceCfgCollectionSettingsFieldName     = "collection_settings"
	slSpaceCfgSearcherSettingsFieldName       = "searcher_settings"

	slCollectionSettingsPlatformMaxSizeFieldName             = "platform_max_size"
	slCollectionSettingsMaxSizeFieldName                     = "max_size"
	slCollectionSettingsMaxIndexedSizeFieldName              = "max_indexed_size"
	slCollectionSettingsPlatformMaxIndexedSizeFieldName      = "platform_max_indexed_size"
	slCollectionSettingsSizeWarningThresholdPercentFieldName = "size_warning_threshold_percent"

	slSpaceCfgStorageLimitFieldName      = "storage_limit"
	slSpaceCfgMonthlySpendLimitFieldName = "monthly_spend_limit"

	slLimitMaxSizeFieldName               = "max_size"
	slLimitMaxAmountMillicentsFieldName   = "max_amount_millicents"
	slLimitTypeFieldName                  = "type"
	slLimitAlertThresholdPercentFieldName = "alert_threshold_percent"
	slLimitCurrencyFieldName              = "currency"

	slSpaceStateLimitsFieldName                = "limits"
	slLimitsStatusAccessModeFieldName          = "access_mode"
	slLimitsStatusStorageFieldName             = "storage"
	slLimitsStatusMonthlySpendFieldName        = "monthly_spend"
	slLimitStatusUsedSizeFieldName             = "used_size"
	slLimitStatusUsedAmountMillicentsFieldName = "used_amount_millicents"
	slLimitStatusLevelFieldName                = "level"
	slLimitStatusPeriodStartFieldName          = "period_start"
	slLimitStatusPeriodEndFieldName            = "period_end"

	// slLimitDefaultAlertThresholdPercent is the alert threshold the API applies when a limit is set without one.
	slLimitDefaultAlertThresholdPercent = 80

	slSearcherSettingsMinIdleTimeoutFieldName     = "min_idle_timeout"
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
			Elem:        &schema.Resource{Schema: serverlessSpaceStateSchema()},
		},
		slSpaceURLFieldName: {
			Description: fmt.Sprintf(slSpaceFieldTemplate, "URL (without port) to access the space"),
			Type:        schema.TypeString,
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
					Description: "Customer-defined maximum size in bytes per collection, used for cost control. Must not exceed `platform_max_size` (defaults to 2 GiB, capped by `platform_max_size`).",
					Type:        schema.TypeInt,
					Optional:    !asDataSource,
					Computed:    true,
				},
				slCollectionSettingsMaxIndexedSizeFieldName: {
					Description: "Customer-defined maximum indexed size in bytes per collection. Must not exceed `platform_max_indexed_size` or the effective maximum size (`max_size`, or `platform_max_size` when `max_size` is unset). Defaults to 1 GiB, capped by `platform_max_indexed_size`.",
					Type:        schema.TypeInt,
					Optional:    !asDataSource,
					Computed:    true,
				},
				slCollectionSettingsPlatformMaxIndexedSizeFieldName: {
					Description: "Platform-enforced maximum indexed size in bytes per collection (0 means unlimited).",
					Type:        schema.TypeInt,
					Computed:    true,
				},
				slCollectionSettingsSizeWarningThresholdPercentFieldName: {
					Description: "Percentage of the effective maximum size at which a size warning is sent for a collection (between 40 and 100, defaults to 80).",
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
				slSearcherSettingsMinIdleTimeoutFieldName: {
					Description: `Platform-defined minimum time a search worker stays idle before it may be scaled down, as a Go duration string (defaults to "5m").`,
					Type:        schema.TypeString,
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
	s[slSpaceCfgStorageLimitFieldName] = serverlessSpaceLimitSchema(asDataSource,
		"Limit on the total storage of the space. When not set, storage is unlimited (the default). "+
			"Removing the block removes the limit.",
		slLimitMaxSizeFieldName, &schema.Schema{
			Description: "Maximum total storage of the space, in bytes. It may be set below the current usage, in which case it takes effect immediately.",
			Type:        schema.TypeInt,
			Required:    !asDataSource,
			Computed:    asDataSource,
		},
		"With a hard limit, upserts are rejected once storage reaches the limit, while searches and deletes keep working. "+
			"Writes resume once storage drops below the limit, or when the limit is raised, switched to a soft limit, or removed. "+
			"With a soft limit, an alert is sent and all operations keep working.")
	spend := serverlessSpaceLimitSchema(asDataSource,
		"Limit on the monthly spend of the space for searches and writes combined (storage cost is not included). "+
			"The spend period is the calendar month in UTC. When not set, spend is unlimited (the default). "+
			"Removing the block removes the limit.",
		slLimitMaxAmountMillicentsFieldName, &schema.Schema{
			Description: "Maximum spend per calendar month, in millicents of `currency`. A millicent is 1/1000 of a cent " +
				"(0.00001 USD), so 1 USD is 100000 millicents and 50 USD is 5000000. " +
				"It may be set below the current spend, in which case it takes effect immediately.",
			Type:     schema.TypeInt,
			Required: !asDataSource,
			Computed: asDataSource,
		},
		"With a hard limit, all read and write operations (including deletes) are rejected once spend reaches the limit, "+
			"until the next month starts, or until the limit is raised, switched to a soft limit, or removed. "+
			"With a soft limit, an alert is sent and all operations keep working.")
	spend.Elem.(*schema.Resource).Schema[slLimitCurrencyFieldName] = &schema.Schema{
		Description: "Currency of `max_amount_millicents`, as a 3-letter ISO 4217 code (currently always `USD`).",
		Type:        schema.TypeString,
		Computed:    true,
	}
	s[slSpaceCfgMonthlySpendLimitFieldName] = spend

	if !asDataSource {
		s[slSpaceCfgAllowedIPSourceRangesFieldName].Elem = &schema.Schema{
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
		}
		searcher := s[slSpaceCfgSearcherSettingsFieldName].Elem.(*schema.Resource).Schema
		searcher[slSearcherSettingsMaxWorkersFieldName].ValidateDiagFunc = validation.ToDiagFunc(validation.IntAtLeast(1))
		collection := s[slSpaceCfgCollectionSettingsFieldName].Elem.(*schema.Resource).Schema
		collection[slCollectionSettingsMaxSizeFieldName].ValidateDiagFunc = validation.ToDiagFunc(validation.IntAtLeast(0))
		collection[slCollectionSettingsMaxIndexedSizeFieldName].ValidateDiagFunc = validation.ToDiagFunc(validation.IntAtLeast(0))
		collection[slCollectionSettingsSizeWarningThresholdPercentFieldName].ValidateDiagFunc = validation.ToDiagFunc(validation.IntBetween(40, 100))
	}
	return s
}

// serverlessSpaceLimitSchema returns the schema of a storage or monthly spend limit block. The block is optional and
// not computed, so leaving it out of the configuration means unlimited, and removing it removes the limit.
func serverlessSpaceLimitSchema(asDataSource bool, description, maxFieldName string, maxField *schema.Schema, typeDescription string) *schema.Schema {
	validTypes := protoEnumNames(spacev1.LimitType_name)
	limitType := &schema.Schema{
		Description: fmt.Sprintf("What happens when usage reaches the limit. Must be one of: %s (defaults to `%s`). %s",
			strings.Join(validTypes, ", "), spacev1.LimitType_LIMIT_TYPE_HARD.String(), typeDescription),
		Type:     schema.TypeString,
		Optional: !asDataSource,
		Computed: asDataSource,
	}
	threshold := &schema.Schema{
		Description: fmt.Sprintf("Percentage of the limit at which an alert is sent to the account owner (between 1 and 99, defaults to %d).",
			slLimitDefaultAlertThresholdPercent),
		Type:     schema.TypeInt,
		Optional: !asDataSource,
		Computed: asDataSource,
	}
	if !asDataSource {
		// The API always returns explicit values, so the defaults match what it applies to keep plans clean.
		limitType.Default = spacev1.LimitType_LIMIT_TYPE_HARD.String()
		limitType.ValidateDiagFunc = validation.ToDiagFunc(validation.StringInSlice(validTypes, false))
		threshold.Default = slLimitDefaultAlertThresholdPercent
		threshold.ValidateDiagFunc = validation.ToDiagFunc(validation.IntBetween(1, 99))
		maxField.ValidateDiagFunc = validation.ToDiagFunc(validation.IntAtLeast(1))
	}
	return &schema.Schema{
		Description: description,
		Type:        schema.TypeList,
		Optional:    !asDataSource,
		Computed:    asDataSource,
		MaxItems:    configurableMaxItems(1, asDataSource),
		Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			maxFieldName:                          maxField,
			slLimitTypeFieldName:                  limitType,
			slLimitAlertThresholdPercentFieldName: threshold,
		}},
	}
}

// serverlessSpaceStateSchema is the (computed) state of a space: the common phase and reason, plus the usage against
// its limits.
func serverlessSpaceStateSchema() map[string]*schema.Schema {
	s := serverlessStateSchema()
	usageLevel := func(limit string) *schema.Schema {
		return &schema.Schema{
			Description: fmt.Sprintf("How the usage relates to the %s limit: one of %s.", limit,
				strings.Join(protoEnumNames(spacev1.LimitUsageLevel_name), ", ")),
			Type:     schema.TypeString,
			Computed: true,
		}
	}
	s[slSpaceStateLimitsFieldName] = &schema.Schema{
		Description: "Usage of the space against its storage and monthly spend limits, and which operations it currently accepts. " +
			"Usage is refreshed about every 5 minutes. Empty until usage has been measured.",
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{Schema: map[string]*schema.Schema{
			slLimitsStatusAccessModeFieldName: {
				Description: fmt.Sprintf("Operations the space currently accepts: one of %s. "+
					"`SPACE_ACCESS_MODE_WRITES_BLOCKED` rejects upserts (a hard storage limit is reached); "+
					"`SPACE_ACCESS_MODE_BLOCKED` rejects all read and write operations (a hard monthly spend limit is reached).",
					strings.Join(protoEnumNames(spacev1.SpaceAccessMode_name), ", ")),
				Type:     schema.TypeString,
				Computed: true,
			},
			slLimitsStatusStorageFieldName: {
				Description: "Storage usage against the storage limit.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					slLimitStatusUsedSizeFieldName: {
						Description: "Current total storage of the space, in bytes (reported also without a storage limit).",
						Type:        schema.TypeInt,
						Computed:    true,
					},
					slLimitStatusLevelFieldName: usageLevel("storage"),
				}},
			},
			slLimitsStatusMonthlySpendFieldName: {
				Description: "Spend in the current month against the monthly spend limit.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					slLimitStatusUsedAmountMillicentsFieldName: {
						Description: "Search and write spend in the current month, in millicents of `currency` " +
							"(1/1000 of a cent, so 1 USD is 100000), reported also without a spend limit.",
						Type:     schema.TypeInt,
						Computed: true,
					},
					slLimitStatusLevelFieldName: usageLevel("monthly spend"),
					slLimitStatusPeriodStartFieldName: {
						Description: "Start of the current month (inclusive), in UTC.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					slLimitStatusPeriodEndFieldName: {
						Description: "End of the current month (exclusive), in UTC. A reached hard spend limit unblocks at this time.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					slLimitCurrencyFieldName: {
						Description: "Currency of `used_amount_millicents` (currently always `USD`).",
						Type:        schema.TypeString,
						Computed:    true,
					},
				}},
			},
		}},
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
		slSpaceStateFieldName:               flattenServerlessSpaceState(space.GetState()),
		slSpaceURLFieldName:                 endpoint.GetUrl(),
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
			slCollectionSettingsPlatformMaxSizeFieldName:             int(cfg.GetCollectionSettings().GetPlatformMaxSize()),
			slCollectionSettingsMaxSizeFieldName:                     int(cfg.GetCollectionSettings().GetMaxSize()),
			slCollectionSettingsMaxIndexedSizeFieldName:              int(cfg.GetCollectionSettings().GetMaxIndexedSize()),
			slCollectionSettingsPlatformMaxIndexedSizeFieldName:      int(cfg.GetCollectionSettings().GetPlatformMaxIndexedSize()),
			slCollectionSettingsSizeWarningThresholdPercentFieldName: int(cfg.GetCollectionSettings().GetSizeWarningThresholdPercent()),
		}},
		slSpaceCfgSearcherSettingsFieldName: []interface{}{map[string]interface{}{
			slSearcherSettingsMinIdleTimeoutFieldName:     formatDuration(cfg.GetSearcherSettings().GetMinIdleTimeout()),
			slSearcherSettingsPlatformMaxWorkersFieldName: int(cfg.GetSearcherSettings().GetPlatformMaxWorkers()),
			slSearcherSettingsMaxWorkersFieldName:         int(cfg.GetSearcherSettings().GetMaxWorkers()),
		}},
		slSpaceCfgStorageLimitFieldName:      flattenServerlessStorageLimit(cfg.GetStorageLimit()),
		slSpaceCfgMonthlySpendLimitFieldName: flattenServerlessSpendLimit(cfg.GetMonthlySpendLimit()),
	}}
}

// flattenServerlessStorageLimit flattens a storage limit; an unset limit (unlimited) is an empty list.
func flattenServerlessStorageLimit(l *spacev1.StorageLimit) []interface{} {
	if l == nil {
		return []interface{}{}
	}
	return []interface{}{map[string]interface{}{
		slLimitMaxSizeFieldName:               int(l.GetMaxSize()),
		slLimitTypeFieldName:                  flattenProtoEnum(l.GetType()),
		slLimitAlertThresholdPercentFieldName: int(l.GetAlertThresholdPercent()),
	}}
}

// flattenServerlessSpendLimit flattens a monthly spend limit; an unset limit (unlimited) is an empty list.
func flattenServerlessSpendLimit(l *spacev1.SpendLimit) []interface{} {
	if l == nil {
		return []interface{}{}
	}
	return []interface{}{map[string]interface{}{
		slLimitMaxAmountMillicentsFieldName:   int(l.GetMaxAmountMillicents()),
		slLimitTypeFieldName:                  flattenProtoEnum(l.GetType()),
		slLimitAlertThresholdPercentFieldName: int(l.GetAlertThresholdPercent()),
		slLimitCurrencyFieldName:              l.GetCurrency(),
	}}
}

// flattenServerlessSpaceState flattens the state of a space, including its usage against its limits.
func flattenServerlessSpaceState(state *spacev1.SpaceState) []interface{} {
	out := flattenServerlessState(flattenProtoEnum(state.GetPhase()), state.GetReason())
	limits := []interface{}{}
	if l := state.GetLimits(); l != nil {
		storage := l.GetStorage()
		spend := l.GetMonthlySpend()
		limits = []interface{}{map[string]interface{}{
			slLimitsStatusAccessModeFieldName: flattenProtoEnum(l.GetAccessMode()),
			slLimitsStatusStorageFieldName: []interface{}{map[string]interface{}{
				slLimitStatusUsedSizeFieldName: int(storage.GetUsedSize()),
				slLimitStatusLevelFieldName:    flattenProtoEnum(storage.GetLevel()),
			}},
			slLimitsStatusMonthlySpendFieldName: []interface{}{map[string]interface{}{
				slLimitStatusUsedAmountMillicentsFieldName: int(spend.GetUsedAmountMillicents()),
				slLimitStatusLevelFieldName:                flattenProtoEnum(spend.GetLevel()),
				slLimitStatusPeriodStartFieldName:          formatTime(spend.GetPeriodStart()),
				slLimitStatusPeriodEndFieldName:            formatTime(spend.GetPeriodEnd()),
				slLimitCurrencyFieldName:                   spend.GetCurrency(),
			}},
		}}
	}
	out[0].(map[string]interface{})[slSpaceStateLimitsFieldName] = limits
	return out
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
		if maxIndexedSize, ok := csItem[slCollectionSettingsMaxIndexedSizeFieldName].(int); ok && maxIndexedSize > 0 {
			cfg.CollectionSettings.MaxIndexedSize = newPointer(uint64(maxIndexedSize))
		}
		if threshold, ok := csItem[slCollectionSettingsSizeWarningThresholdPercentFieldName].(int); ok && threshold > 0 {
			cfg.CollectionSettings.SizeWarningThresholdPercent = newPointer(uint32(threshold))
		}
	}
	if ss, ok := item[slSpaceCfgSearcherSettingsFieldName].([]interface{}); ok && len(ss) > 0 && ss[0] != nil {
		ssItem := ss[0].(map[string]interface{})
		cfg.SearcherSettings = &spacev1.SearcherSettings{}
		if maxWorkers, ok := ssItem[slSearcherSettingsMaxWorkersFieldName].(int); ok && maxWorkers > 0 {
			cfg.SearcherSettings.MaxWorkers = newPointer(uint64(maxWorkers))
		}
	}
	if l, ok := firstBlock(item[slSpaceCfgStorageLimitFieldName]); ok {
		cfg.StorageLimit = &spacev1.StorageLimit{
			MaxSize:               uint64(intOf(l[slLimitMaxSizeFieldName])),
			Type:                  expandLimitType(l[slLimitTypeFieldName]),
			AlertThresholdPercent: expandAlertThresholdPercent(l[slLimitAlertThresholdPercentFieldName]),
		}
	}
	if l, ok := firstBlock(item[slSpaceCfgMonthlySpendLimitFieldName]); ok {
		cfg.MonthlySpendLimit = &spacev1.SpendLimit{
			MaxAmountMillicents:   uint64(intOf(l[slLimitMaxAmountMillicentsFieldName])),
			Type:                  expandLimitType(l[slLimitTypeFieldName]),
			AlertThresholdPercent: expandAlertThresholdPercent(l[slLimitAlertThresholdPercentFieldName]),
		}
	}
	return cfg
}

// firstBlock returns the attributes of the single element of a MaxItems=1 block, and false when the block is absent.
func firstBlock(v interface{}) (map[string]interface{}, bool) {
	list, ok := v.([]interface{})
	if !ok || len(list) == 0 || list[0] == nil {
		return nil, false
	}
	m, ok := list[0].(map[string]interface{})
	return m, ok
}

// intOf returns v as an int, or 0 when it is not one.
func intOf(v interface{}) int {
	i, _ := v.(int)
	return i
}

// expandLimitType maps the configured limit type to the API enum; empty leaves it to the API default (hard).
func expandLimitType(v interface{}) spacev1.LimitType {
	name, _ := v.(string)
	return spacev1.LimitType(spacev1.LimitType_value[name])
}

// expandAlertThresholdPercent returns the configured alert threshold, or nil to leave it to the API default.
func expandAlertThresholdPercent(v interface{}) *uint32 {
	if p := intOf(v); p > 0 {
		return newPointer(uint32(p))
	}
	return nil
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
	{slSpaceConfigurationFieldName + ".0." + slSpaceCfgCollectionSettingsFieldName + ".0." + slCollectionSettingsMaxIndexedSizeFieldName, "configuration.collection_settings.max_indexed_size"},
	{slSpaceConfigurationFieldName + ".0." + slSpaceCfgCollectionSettingsFieldName + ".0." + slCollectionSettingsSizeWarningThresholdPercentFieldName, "configuration.collection_settings.size_warning_threshold_percent"},
	{slSpaceConfigurationFieldName + ".0." + slSpaceCfgSearcherSettingsFieldName + ".0." + slSearcherSettingsMaxWorkersFieldName, "configuration.searcher_settings.max_workers"},
	// The limits are sent as a whole: an absent block sends an unset limit, which removes it (unlimited).
	{slSpaceConfigurationFieldName + ".0." + slSpaceCfgStorageLimitFieldName, "configuration.storage_limit"},
	{slSpaceConfigurationFieldName + ".0." + slSpaceCfgMonthlySpendLimitFieldName, "configuration.monthly_spend_limit"},
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
