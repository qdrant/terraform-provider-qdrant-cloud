package qdrant

import (
	"fmt"
	"reflect"

	"github.com/goccy/go-yaml"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/protobuf/types/known/structpb"

	qch "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/hybrid/v1"
)

const (
	hcEnvFieldTemplate                       = "Hybrid cloud environment Schema %s field"
	hcEnvIdFieldName                         = "id"
	hcEnvAccountIdFieldName                  = "account_id"
	hcEnvNameFieldName                       = "name"
	hcEnvConfigurationFieldName              = "configuration"
	hcEnvCfgNamespaceFieldName               = "namespace"
	hcEnvCreatedAtFieldName                  = "created_at"
	hcEnvLastModifiedAtFieldName             = "last_modified_at"
	hcEnvCreatedByEmailFieldName             = "created_by_email"
	hcEnvBootstrapCommandsGeneratedFieldName = "bootstrap_commands_generated"
	hcEnvBootstrapCommandsFieldName          = "bootstrap_commands"
	hcEnvBootstrapCommandsVersionFieldName   = "bootstrap_commands_version"
	hcEnvStatusFieldName                     = "status"

	hcEnvCfgLastModifiedAtFieldName             = "last_modified_at"
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
	hcEnvCfgAdvancedOperatorSettingsFieldName   = "advanced_operator_settings"

	hcEnvStatusLastModifiedAtFieldName           = "last_modified_at"
	hcEnvStatusPhaseFieldName                    = "phase"
	hcEnvStatusKubernetesVersionFieldName        = "kubernetes_version"
	hcEnvStatusNumberOfNodesFieldName            = "number_of_nodes"
	hcEnvStatusClusterCreationReadinessFieldName = "cluster_creation_readiness"
	hcEnvStatusKubernetesDistributionFieldName   = "kubernetes_distribution"
	hcEnvStatusMessageFieldName                  = "message"
	hcEnvStatusCapabilitiesFieldName             = "capabilities"
	hcEnvStatusCapabilitiesVolumeSnapshotField   = "volume_snapshot"
	hcEnvStatusCapabilitiesVolumeExpansionField  = "volume_expansion"
	hcEnvStatusComponentStatusesFieldName        = "component_statuses"
	hcEnvStatusComponentNameField                = "name"
	hcEnvStatusComponentVersionField             = "version"
	hcEnvStatusComponentPhaseField               = "phase"
	hcEnvStatusComponentMessageField             = "message"
	hcEnvStatusComponentNamespaceField           = "namespace"
	hcEnvStatusStorageClassesFieldName           = "storage_classes"
	hcEnvStatusStorageClassNameField             = "name"
	hcEnvStatusStorageClassDefaultField          = "default"
	hcEnvStatusStorageClassProvisionerField      = "provisioner"
	hcEnvStatusStorageClassAllowExpansionField   = "allow_volume_expansion"
	hcEnvStatusStorageClassReclaimPolicyField    = "reclaim_policy"
	hcEnvStatusStorageClassParametersField       = "parameters"
	hcEnvStatusVolumeSnapshotClassesFieldName    = "volume_snapshot_classes"
	hcEnvStatusVSCNameField                      = "name"
	hcEnvStatusVSCDriverField                    = "driver"
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
		hcEnvCreatedByEmailFieldName: {
			Description: fmt.Sprintf(hcEnvFieldTemplate, "The email of the user who created the hybrid cloud environment"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		hcEnvBootstrapCommandsGeneratedFieldName: {
			Description: fmt.Sprintf(hcEnvFieldTemplate, "Set if the generate bootstrap commands has been called at least once"),
			Type:        schema.TypeBool,
			Computed:    true,
		},

		hcEnvStatusFieldName: {
			Description: "Current status of the hybrid cloud environment (read-only).",
			Type:        schema.TypeList,
			Computed:    true,
			Elem:        &schema.Resource{Schema: accountsHybridCloudEnvironmentStatusSchema()},
		},

		hcEnvBootstrapCommandsFieldName: {
			Description: fmt.Sprintf(hcEnvFieldTemplate, "Commands to bootstrap a Kubernetes cluster into this environment"),
			Type:        schema.TypeList,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Computed:    true,
			Sensitive:   true,
		},

		hcEnvBootstrapCommandsVersionFieldName: {
			Description: "Version knob to (re)generate bootstrap commands. -1 = never generate, 0 = idle/do not (re)generate, >0 = generate/rotate.",
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
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
		hcEnvCfgLastModifiedAtFieldName: {
			Description: "Last modification timestamp of the configuration.",
			Type:        schema.TypeString,
			Computed:    true,
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
		hcEnvCfgAdvancedOperatorSettingsFieldName: {
			Description: "Advanced operator settings as a YAML string.",
			Type:        schema.TypeString,
			Optional:    true,
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				var oldData, newData interface{}
				if err := yaml.Unmarshal([]byte(old), &oldData); err != nil {
					return false
				}
				if err := yaml.Unmarshal([]byte(new), &newData); err != nil {
					return false
				}
				return reflect.DeepEqual(oldData, newData)
			},
			ValidateFunc: func(v interface{}, k string) (ws []string, es []error) {
				if err := yaml.Unmarshal([]byte(v.(string)), new(interface{})); err != nil {
					es = append(es, fmt.Errorf("%q contains invalid YAML: %w", k, err))
				}
				return
			},
			// StateFunc normalizes the YAML on save, which is good practice with DiffSuppressFunc.
			// This ensures the state file has a consistent format, even if user input varies.
			StateFunc: func(v interface{}) string {
				var data interface{}
				if err := yaml.Unmarshal([]byte(v.(string)), &data); err != nil {
					return v.(string) // On error, keep original
				}
				out, _ := yaml.Marshal(data)
				return string(out)
			},
		},
	}
}

// accountsHybridCloudEnvironmentStatusSchema defines the read-only status shape for a hybrid cloud environment.
func accountsHybridCloudEnvironmentStatusSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		hcEnvStatusLastModifiedAtFieldName: {
			Description: "Timestamp when the status was last updated.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		hcEnvStatusPhaseFieldName: {
			Description: "Environment status phase.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		hcEnvStatusKubernetesVersionFieldName: {
			Description: "Kubernetes version.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		hcEnvStatusNumberOfNodesFieldName: {
			Description: "Number of Kubernetes nodes.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		hcEnvStatusClusterCreationReadinessFieldName: {
			Description: "Cluster creation readiness.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		hcEnvStatusKubernetesDistributionFieldName: {
			Description: "Kubernetes distribution.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		hcEnvStatusMessageFieldName: {
			Description: "Status message.",
			Type:        schema.TypeString,
			Computed:    true,
		},

		// capabilities (single block)
		hcEnvStatusCapabilitiesFieldName: {
			Description: "Environment capabilities.",
			Type:        schema.TypeList,
			Computed:    true,
			Elem:        &schema.Resource{Schema: accountsHybridCloudEnvironmentStatusCapabilitiesSchema()},
		},

		// component_statuses (list)
		hcEnvStatusComponentStatusesFieldName: {
			Description: "Status of deployed components.",
			Type:        schema.TypeList,
			Computed:    true,
			Elem:        &schema.Resource{Schema: accountsHybridCloudEnvironmentStatusComponentStatusesSchema()},
		},

		// storage_classes (list)
		hcEnvStatusStorageClassesFieldName: {
			Description: "Available storage classes.",
			Type:        schema.TypeList,
			Computed:    true,
			Elem:        &schema.Resource{Schema: accountsHybridCloudEnvironmentStatusStorageClassesSchema()},
		},

		// volume_snapshot_classes (list)
		hcEnvStatusVolumeSnapshotClassesFieldName: {
			Description: "Available volume snapshot classes.",
			Type:        schema.TypeList,
			Computed:    true,
			Elem:        &schema.Resource{Schema: accountsHybridCloudEnvironmentStatusVSCsSchema()},
		},
	}
}

func accountsHybridCloudEnvironmentStatusCapabilitiesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		hcEnvStatusCapabilitiesVolumeSnapshotField: {
			Description: "Whether the cluster supports Kubernetes VolumeSnapshot functionality.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		hcEnvStatusCapabilitiesVolumeExpansionField: {
			Description: "Whether PersistentVolumeClaims can be expanded on this cluster.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
	}
}

func accountsHybridCloudEnvironmentStatusComponentStatusesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		hcEnvStatusComponentNameField: {
			Description: "Component name (e.g., Helm release/controller).",
			Type:        schema.TypeString,
			Computed:    true,
		},
		hcEnvStatusComponentVersionField: {
			Description: "Reported component version (chart/app).",
			Type:        schema.TypeString,
			Computed:    true,
		},
		hcEnvStatusComponentPhaseField: {
			Description: "Component phase (e.g., PENDING, READY, ERROR).",
			Type:        schema.TypeString,
			Computed:    true,
		},
		hcEnvStatusComponentMessageField: {
			Description: "Human-readable message for the component state.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		hcEnvStatusComponentNamespaceField: {
			Description: "Kubernetes namespace where the component runs.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func accountsHybridCloudEnvironmentStatusStorageClassesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		hcEnvStatusStorageClassNameField: {
			Description: "StorageClass name.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		hcEnvStatusStorageClassDefaultField: {
			Description: "Whether this StorageClass is the default.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		hcEnvStatusStorageClassProvisionerField: {
			Description: "CSI provisioner for this StorageClass.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		hcEnvStatusStorageClassAllowExpansionField: {
			Description: "Whether PVCs using this StorageClass can be expanded.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		hcEnvStatusStorageClassReclaimPolicyField: {
			Description: "Reclaim policy (e.g., Delete, Retain).",
			Type:        schema.TypeString,
			Computed:    true,
		},
		hcEnvStatusStorageClassParametersField: {
			Description: "Key-value parameters for the provisioner.",
			Type:        schema.TypeMap,
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
	}
}

func accountsHybridCloudEnvironmentStatusVSCsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		hcEnvStatusVSCNameField: {
			Description: "VolumeSnapshotClass name.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		hcEnvStatusVSCDriverField: {
			Description: "CSI snapshot driver for this class.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

// flattenHCEnv maps API object -> Terraform state fields (excluding status).
func flattenHCEnv(env *qch.HybridCloudEnvironment) map[string]interface{} {
	out := map[string]interface{}{
		hcEnvIdFieldName:                         env.GetId(),
		hcEnvAccountIdFieldName:                  env.GetAccountId(),
		hcEnvNameFieldName:                       env.GetName(),
		hcEnvCreatedAtFieldName:                  formatTime(env.GetCreatedAt()),
		hcEnvLastModifiedAtFieldName:             formatTime(env.GetLastModifiedAt()),
		hcEnvCreatedByEmailFieldName:             env.GetCreatedByEmail(),
		hcEnvBootstrapCommandsGeneratedFieldName: env.GetBootstrapCommandsGenerated(),
	}

	out[hcEnvConfigurationFieldName] = flattenHCEnvConfiguration(env.GetConfiguration())
	out[hcEnvStatusFieldName] = flattenHCEnvStatus(env.GetStatus())

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
	if adv := cfg.GetAdvancedOperatorSettings(); adv != nil {
		advMap := adv.AsMap()
		if len(advMap) > 0 {
			if yamlBytes, err := yaml.Marshal(advMap); err == nil {
				configMap[hcEnvCfgAdvancedOperatorSettingsFieldName] = string(yamlBytes)
			}
		}
	}

	if ts := cfg.GetLastModifiedAt(); ts != nil {
		configMap[hcEnvCfgLastModifiedAtFieldName] = formatTime(ts)
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
	if val, ok := m[hcEnvCfgAdvancedOperatorSettingsFieldName]; ok && val.(string) != "" {
		var settingsMap map[string]interface{}
		if err := yaml.Unmarshal([]byte(val.(string)), &settingsMap); err == nil && len(settingsMap) > 0 {
			// Make sure the keys are strings, even for the nested ones (as NewStruct expects that)
			convertedMap := convertMapKeysToStrings(settingsMap)
			if m, ok := convertedMap.(map[string]interface{}); ok {
				if s, err := structpb.NewStruct(m); err == nil {
					config.AdvancedOperatorSettings = s
				}
			}
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

// convertMapKeysToStrings recursively converts map keys from interface{} to string.
// This is necessary because yaml.Unmarshal can create map[interface{}]interface{}.
func convertMapKeysToStrings(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			if ks, ok := k.(string); ok {
				m2[ks] = convertMapKeysToStrings(v)
			}
		}
		return m2
	case map[string]interface{}:
		for k, v := range x {
			x[k] = convertMapKeysToStrings(v)
		}
	case []interface{}:
		for i, v := range x {
			x[i] = convertMapKeysToStrings(v)
		}
	}
	return i
}

func flattenHCEnvStatus(st *qch.HybridCloudEnvironmentStatus) []interface{} {
	if st == nil {
		return []interface{}{}
	}
	m := map[string]interface{}{
		hcEnvStatusPhaseFieldName:                    st.GetPhase().String(),
		hcEnvStatusKubernetesVersionFieldName:        st.GetKubernetesVersion(),
		hcEnvStatusNumberOfNodesFieldName:            int(st.GetNumberOfNodes()),
		hcEnvStatusClusterCreationReadinessFieldName: st.GetClusterCreationReadiness().String(),
		hcEnvStatusKubernetesDistributionFieldName:   st.GetKubernetesDistribution().String(),
		hcEnvStatusMessageFieldName:                  st.GetMessage(),
		hcEnvStatusCapabilitiesFieldName:             flattenHCEnvCapabilities(st.GetCapabilities()),
		hcEnvStatusComponentStatusesFieldName:        flattenHCEnvComponentStatuses(st.GetComponentStatuses()),
		hcEnvStatusStorageClassesFieldName:           flattenHCEnvStorageClasses(st.GetStorageClasses()),
		hcEnvStatusVolumeSnapshotClassesFieldName:    flattenHCEnvVSCs(st.GetVolumeSnapshotClasses()),
	}
	if ts := st.GetLastModifiedAt(); ts != nil {
		m[hcEnvStatusLastModifiedAtFieldName] = formatTime(ts)
	}
	return []interface{}{m}
}

func flattenHCEnvCapabilities(c *qch.HybridCloudEnvironmentCapabilities) []interface{} {
	if c == nil {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			hcEnvStatusCapabilitiesVolumeSnapshotField:  c.GetVolumeSnapshot(),
			hcEnvStatusCapabilitiesVolumeExpansionField: c.GetVolumeExpansion(),
		},
	}
}

func flattenHCEnvComponentStatuses(cs []*qch.HybridCloudEnvironmentComponentStatus) []interface{} {
	if len(cs) == 0 {
		return []interface{}{}
	}
	out := make([]interface{}, 0, len(cs))
	for _, x := range cs {
		out = append(out, map[string]interface{}{
			hcEnvStatusComponentNameField:    x.GetName(),
			hcEnvStatusComponentVersionField: x.GetVersion(),
			hcEnvStatusComponentPhaseField:   x.GetPhase().String(),
			hcEnvStatusComponentMessageField: x.GetMessage(),
			// (namespace intentionally omitted to avoid duplication)
		})
	}
	return out
}

func flattenHCEnvStorageClasses(sc []*qch.HybridCloudEnvironmentStorageClass) []interface{} {
	if len(sc) == 0 {
		return []interface{}{}
	}
	out := make([]interface{}, 0, len(sc))
	for _, x := range sc {
		m := map[string]interface{}{
			hcEnvStatusStorageClassNameField:           x.GetName(),
			hcEnvStatusStorageClassDefaultField:        x.GetDefault(),
			hcEnvStatusStorageClassProvisionerField:    x.GetProvisioner(),
			hcEnvStatusStorageClassAllowExpansionField: x.GetAllowVolumeExpansion(),
			hcEnvStatusStorageClassReclaimPolicyField:  x.GetReclaimPolicy(),
		}
		if len(x.GetParameters()) > 0 {
			pm := map[string]string{}
			for _, kv := range x.GetParameters() {
				pm[kv.GetKey()] = kv.GetValue()
			}
			m[hcEnvStatusStorageClassParametersField] = pm
		}
		out = append(out, m)
	}
	return out
}

func flattenHCEnvVSCs(vsc []*qch.HybridCloudEnvironmentVolumeSnapshotClass) []interface{} {
	if len(vsc) == 0 {
		return []interface{}{}
	}
	out := make([]interface{}, 0, len(vsc))
	for _, x := range vsc {
		out = append(out, map[string]interface{}{
			hcEnvStatusVSCNameField:   x.GetName(),
			hcEnvStatusVSCDriverField: x.GetDriver(),
		})
	}
	return out
}
