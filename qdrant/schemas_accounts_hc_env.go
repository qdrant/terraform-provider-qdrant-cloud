package qdrant

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qch "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/hybrid/v1"
)

const (
	hcEnvFieldTemplate                          = "Hybrid cloud environment Schema %s field"
	hcEnvIdFieldName                            = "id"
	hcEnvAccountIdFieldName                     = "account_id"
	hcEnvNameFieldName                          = "name"
	hcEnvConfigurationFieldName                 = "configuration"
	hcEnvCfgNamespaceFieldName                  = "namespace"
	hcEnvCreatedAtFieldName                     = "created_at"
	hcEnvLastModifiedAtFieldName                = "last_modified_at"
	hcEnvBootstrapCommandsFieldName             = "bootstrap_commands"
	hcEnvCfgHttpProxyUrlFieldName               = "http_proxy_url"
	hcEnvCfgHttpsProxyUrlFieldName              = "https_proxy_url"
	hcEnvCfgNoProxyConfigsFieldName             = "no_proxy_configs"
	hcEnvCfgContainerRegistryUrlFieldName       = "container_registry_url"
	hcEnvCfgChartRepositoryUrlFieldName         = "chart_repository_url"
	hcEnvCfgRegistrySecretNameFieldName         = "registry_secret_name"
	hcEnvCfgCaCertificatesFieldName             = "ca_certificates"
	hcEnvCfgDatabaseStorageClassFieldName       = "database_storage_class"
	hcEnvCfgSnapshotStorageClassFieldName       = "snapshot_storage_class"
	hcEnvCfgVolumeSnapshotStorageClassFieldName = "volume_snapshot_storage_class"
	hcEnvCfgLogLevelFieldName                   = "log_level"
)

func accountsHybridCloudEnvironmentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		hcEnvIdFieldName: {
			Description: fmt.Sprintf(hcEnvFieldTemplate, "ID"),
			Type:        schema.TypeString,
			Computed:    true, // mirror of d.Id()
		},
		hcEnvAccountIdFieldName: {
			Description: fmt.Sprintf(hcEnvFieldTemplate, "Account ID"),
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
		},
		hcEnvNameFieldName: {
			Description: fmt.Sprintf(hcEnvFieldTemplate, "Name"),
			Type:        schema.TypeString,
			Required:    true,
		},
		hcEnvConfigurationFieldName: {
			Description: fmt.Sprintf(hcEnvFieldTemplate, "Configuration"),
			Type:        schema.TypeList,
			Required:    true,
			MaxItems:    1,
			Elem:        &schema.Resource{Schema: accountsHybridCloudEnvironmentConfigurationSchema()},
		},
		hcEnvCreatedAtFieldName: {
			Description: fmt.Sprintf(hcEnvFieldTemplate, "Creation timestamp"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		hcEnvLastModifiedAtFieldName: {
			Description: fmt.Sprintf(hcEnvFieldTemplate, "Last modification timestamp"),
			Type:        schema.TypeString,
			Computed:    true,
		},

		// Intentionally no `status` in schema.

		hcEnvBootstrapCommandsFieldName: {
			Description: fmt.Sprintf(hcEnvFieldTemplate, "Commands to bootstrap a Kubernetes cluster into this environment"),
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Computed:    true,
			Sensitive:   true,
		},
	}
}

// accountsHybridCloudEnvironmentConfigurationSchema defines the schema for the configuration of a hybrid cloud environment.
func accountsHybridCloudEnvironmentConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		hcEnvCfgNamespaceFieldName: {
			Description: "The Kubernetes namespace where the Qdrant hybrid cloud components will be deployed.",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		hcEnvCfgHttpProxyUrlFieldName: {
			Description: "Optional HTTP proxy URL.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		hcEnvCfgHttpsProxyUrlFieldName: {
			Description: "Optional HTTPS proxy URL.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		hcEnvCfgNoProxyConfigsFieldName: {
			Description: "List of hosts that should not be proxied.",
			Type:        schema.TypeList,
			Optional:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		hcEnvCfgContainerRegistryUrlFieldName: {
			Description: "Container registry URL.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		hcEnvCfgChartRepositoryUrlFieldName: {
			Description: "Chart registry URL.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		hcEnvCfgRegistrySecretNameFieldName: {
			Description: "Kubernetes secret name containing registry credentials.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		hcEnvCfgCaCertificatesFieldName: {
			Description: "CA certificates for custom certificate authorities.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		hcEnvCfgDatabaseStorageClassFieldName: {
			Description: "Default database storage class.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		hcEnvCfgSnapshotStorageClassFieldName: {
			Description: "Default snapshot storage class.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		hcEnvCfgVolumeSnapshotStorageClassFieldName: {
			Description: "Default volume snapshot storage class.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		hcEnvCfgLogLevelFieldName: {
			Description: "Log level for deployed components.",
			Type:        schema.TypeString,
			Optional:    true,
		},
	}
}

// flattenHCEnv maps API object -> Terraform state fields (excluding status).
func flattenHCEnv(env *qch.HybridCloudEnvironment) map[string]interface{} {
	out := map[string]interface{}{
		hcEnvIdFieldName:             env.GetId(),
		hcEnvAccountIdFieldName:      env.GetAccountId(),
		hcEnvNameFieldName:           env.GetName(),
		hcEnvCreatedAtFieldName:      formatTime(env.GetCreatedAt()),
		hcEnvLastModifiedAtFieldName: formatTime(env.GetLastModifiedAt()),
	}

	out[hcEnvConfigurationFieldName] = flattenHCEnvConfiguration(env.GetConfiguration())

	return out
}

// flattenHCEnvConfiguration maps the configuration API object to a Terraform state list.
func flattenHCEnvConfiguration(cfg *qch.HybridCloudEnvironmentConfiguration) []interface{} {
	if cfg == nil {
		return []interface{}{}
	}
	configMap := map[string]interface{}{
		hcEnvCfgNamespaceFieldName:                  cfg.GetNamespace(),
		hcEnvCfgHttpProxyUrlFieldName:               cfg.GetHttpProxyUrl(),
		hcEnvCfgHttpsProxyUrlFieldName:              cfg.GetHttpsProxyUrl(),
		hcEnvCfgNoProxyConfigsFieldName:             cfg.GetNoProxyConfigs(),
		hcEnvCfgContainerRegistryUrlFieldName:       cfg.GetContainerRegistryUrl(),
		hcEnvCfgChartRepositoryUrlFieldName:         cfg.GetChartRepositoryUrl(),
		hcEnvCfgRegistrySecretNameFieldName:         cfg.GetRegistrySecretName(),
		hcEnvCfgCaCertificatesFieldName:             cfg.GetCaCertificates(),
		hcEnvCfgDatabaseStorageClassFieldName:       cfg.GetDatabaseStorageClass(),
		hcEnvCfgSnapshotStorageClassFieldName:       cfg.GetSnapshotStorageClass(),
		hcEnvCfgVolumeSnapshotStorageClassFieldName: cfg.GetVolumeSnapshotStorageClass(),
	}
	// Handle enum to avoid panic on nil
	if logLevel := cfg.GetLogLevel(); logLevel != qch.HybridCloudEnvironmentConfigurationLogLevel_HYBRID_CLOUD_ENVIRONMENT_CONFIGURATION_LOG_LEVEL_UNSPECIFIED {
		configMap[hcEnvCfgLogLevelFieldName] = logLevel.String()
	}
	return []interface{}{configMap}
}

// expandHCEnvConfiguration maps the configuration from Terraform state to an API object.
func expandHCEnvConfiguration(v []interface{}) *qch.HybridCloudEnvironmentConfiguration {
	if len(v) == 0 || v[0] == nil {
		return &qch.HybridCloudEnvironmentConfiguration{}
	}

	m := v[0].(map[string]interface{})
	config := &qch.HybridCloudEnvironmentConfiguration{}

	if ns, ok := m[hcEnvCfgNamespaceFieldName]; ok {
		config.Namespace = ns.(string)
	}
	if val, ok := m[hcEnvCfgHttpProxyUrlFieldName]; ok && val.(string) != "" {
		config.HttpProxyUrl = newPointer(val.(string))
	}
	if val, ok := m[hcEnvCfgHttpsProxyUrlFieldName]; ok && val.(string) != "" {
		config.HttpsProxyUrl = newPointer(val.(string))
	}
	if val, ok := m[hcEnvCfgNoProxyConfigsFieldName]; ok {
		config.NoProxyConfigs = interfaceSliceToStringSlice(val.([]interface{}))
	}
	if val, ok := m[hcEnvCfgContainerRegistryUrlFieldName]; ok && val.(string) != "" {
		config.ContainerRegistryUrl = newPointer(val.(string))
	}
	if val, ok := m[hcEnvCfgChartRepositoryUrlFieldName]; ok && val.(string) != "" {
		config.ChartRepositoryUrl = newPointer(val.(string))
	}
	if val, ok := m[hcEnvCfgRegistrySecretNameFieldName]; ok && val.(string) != "" {
		config.RegistrySecretName = newPointer(val.(string))
	}
	if val, ok := m[hcEnvCfgCaCertificatesFieldName]; ok && val.(string) != "" {
		config.CaCertificates = newPointer(val.(string))
	}
	if val, ok := m[hcEnvCfgDatabaseStorageClassFieldName]; ok && val.(string) != "" {
		config.DatabaseStorageClass = newPointer(val.(string))
	}
	if val, ok := m[hcEnvCfgSnapshotStorageClassFieldName]; ok && val.(string) != "" {
		config.SnapshotStorageClass = newPointer(val.(string))
	}
	if val, ok := m[hcEnvCfgVolumeSnapshotStorageClassFieldName]; ok && val.(string) != "" {
		config.VolumeSnapshotStorageClass = newPointer(val.(string))
	}
	if val, ok := m[hcEnvCfgLogLevelFieldName]; ok && val.(string) != "" {
		logLevel, llOK := qch.HybridCloudEnvironmentConfigurationLogLevel_value[val.(string)]
		if llOK {
			config.LogLevel = newPointer(qch.HybridCloudEnvironmentConfigurationLogLevel(logLevel))
		}
	}

	return config
}

// expandHCEnv builds the payload from config; validates required fields.
func expandHCEnv(d *schema.ResourceData, defaultAccountID string) (*qch.HybridCloudEnvironment, error) {
	accountID := defaultAccountID
	if v, ok := d.GetOk(hcEnvAccountIdFieldName); ok && v.(string) != "" {
		accountID = v.(string)
	}
	if accountID == "" {
		return nil, fmt.Errorf("account ID not specified")
	}

	env := &qch.HybridCloudEnvironment{
		AccountId:     accountID,
		Id:            d.Id(),
		Name:          d.Get(hcEnvNameFieldName).(string),
		Configuration: expandHCEnvConfiguration(d.Get(hcEnvConfigurationFieldName).([]interface{})),
	}

	// For create, namespace is required. For update, it's optional but if provided, it must not be empty.
	if env.Configuration.GetNamespace() == "" && d.IsNewResource() {
		return nil, fmt.Errorf("configuration.namespace must be set for a new environment")
	}

	return env, nil
}

// interfaceSliceToStringSlice converts a slice of interface{} to a slice of string.
func interfaceSliceToStringSlice(s []interface{}) []string {
	res := make([]string, len(s))
	for i, v := range s {
		res[i] = v.(string)
	}
	return res
}
