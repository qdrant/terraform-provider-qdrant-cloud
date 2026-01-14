package qdrant

import (
	"testing"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	qcCluster "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/v1"
	commonv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/common/v1"
	qch "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/hybrid/v1"
)

func TestFlattenHCEnv(t *testing.T) {
	env := &qch.HybridCloudEnvironment{
		AccountId:                  "00000000-1000-0000-0000-000000000001",
		Id:                         "00000000-0000-0000-0000-000000000001",
		CreatedAt:                  timestamppb.New(time.Date(2025, 9, 16, 9, 41, 5, 0, time.UTC)),
		LastModifiedAt:             timestamppb.New(time.Date(2025, 9, 16, 9, 41, 8, 0, time.UTC)),
		Name:                       "local-test-2",
		CreatedByEmail:             "creator@example.com",
		BootstrapCommandsGenerated: true,
		Configuration: &qch.HybridCloudEnvironmentConfiguration{
			Namespace:                  "qdrant-hc",
			HttpProxyUrl:               newPointer("http://proxy.example.com"),
			HttpsProxyUrl:              newPointer("https://proxy.example.com"),
			NoProxyConfigs:             []string{"localhost", "127.0.0.1"},
			ContainerRegistryUrl:       newPointer("registry.example.com"),
			ChartRepositoryUrl:         newPointer("charts.example.com"),
			RegistrySecretName:         newPointer("reg-secret"),
			CaCertificates:             newPointer("ca-cert-data"),
			DatabaseStorageClass:       newPointer("db-storage"),
			SnapshotStorageClass:       newPointer("snap-storage"),
			VolumeSnapshotStorageClass: newPointer("vol-snap-storage"),
			LogLevel:                   newPointer(qch.HybridCloudEnvironmentConfigurationLogLevel_HYBRID_CLOUD_ENVIRONMENT_CONFIGURATION_LOG_LEVEL_DEBUG),
			AdvancedOperatorSettings: func() *structpb.Struct {
				s, err := structpb.NewStruct(map[string]interface{}{"key": "value", "nested": map[string]interface{}{"num": 1}})
				require.NoError(t, err)
				return s
			}(),
			NodeSelector: []*commonv1.KeyValue{
				{Key: "key1", Value: "value1"},
			},
			Tolerations: []*qcCluster.Toleration{
				{
					Key:      newPointer("key1"),
					Operator: newPointer(qcCluster.TolerationOperator_TOLERATION_OPERATOR_EQUAL),
					Value:    newPointer("value1"),
					Effect:   newPointer(qcCluster.TolerationEffect_TOLERATION_EFFECT_NO_SCHEDULE),
				},
			},
			ControlPlaneLabels: []*commonv1.KeyValue{
				{Key: "label2", Value: "value2"},
			},
		},
	}
	got := flattenHCEnv(env)
	want := map[string]interface{}{
		hcEnvIdFieldName:                         env.GetId(),
		hcEnvAccountIdFieldName:                  env.GetAccountId(),
		hcEnvNameFieldName:                       env.GetName(),
		hcEnvCreatedAtFieldName:                  formatTime(env.GetCreatedAt()),
		hcEnvLastModifiedAtFieldName:             formatTime(env.GetLastModifiedAt()),
		hcEnvCreatedByEmailFieldName:             env.GetCreatedByEmail(),
		hcEnvBootstrapCommandsGeneratedFieldName: env.GetBootstrapCommandsGenerated(),
		hcEnvConfigurationFieldName: []interface{}{
			map[string]interface{}{
				hcEnvCfgNamespaceFieldName:                  env.GetConfiguration().GetNamespace(),
				hcEnvCfgHttpProxyUrlFieldName:               env.GetConfiguration().GetHttpProxyUrl(),
				hcEnvCfgHttpsProxyUrlFieldName:              env.GetConfiguration().GetHttpsProxyUrl(),
				hcEnvCfgNoProxyConfigsFieldName:             env.GetConfiguration().GetNoProxyConfigs(),
				hcEnvCfgContainerRegistryUrlFieldName:       env.GetConfiguration().GetContainerRegistryUrl(),
				hcEnvCfgChartRepositoryUrlFieldName:         env.GetConfiguration().GetChartRepositoryUrl(),
				hcEnvCfgRegistrySecretNameFieldName:         env.GetConfiguration().GetRegistrySecretName(),
				hcEnvCfgCaCertificatesFieldName:             env.GetConfiguration().GetCaCertificates(),
				hcEnvCfgDatabaseStorageClassFieldName:       env.GetConfiguration().GetDatabaseStorageClass(),
				hcEnvCfgSnapshotStorageClassFieldName:       env.GetConfiguration().GetSnapshotStorageClass(),
				hcEnvCfgVolumeSnapshotStorageClassFieldName: env.GetConfiguration().GetVolumeSnapshotStorageClass(),
				hcEnvCfgLogLevelFieldName:                   env.GetConfiguration().GetLogLevel().String(),
				hcEnvCfgAdvancedOperatorSettingsFieldName: func() string {
					b, err := yaml.Marshal(env.GetConfiguration().GetAdvancedOperatorSettings().AsMap())
					require.NoError(t, err)
					return string(b)
				}(),
				hcEnvCfgNodeSelectorFieldName: []interface{}{
					map[string]interface{}{"key": "key1", "value": "value1"},
				},
				hcEnvCfgTolerationsFieldName: []interface{}{
					map[string]interface{}{
						tolerationKeyFieldName:      "key1",
						tolerationOperatorFieldName: "TOLERATION_OPERATOR_EQUAL",
						tolerationValueFieldName:    "value1",
						tolerationEffectFieldName:   "TOLERATION_EFFECT_NO_SCHEDULE",
					},
				},
				hcEnvCfgControlPlaneLabelsFieldName: []interface{}{
					map[string]interface{}{"key": "label2", "value": "value2"},
				},
			},
		},
		hcEnvStatusFieldName: []interface{}{},
	}

	assert.Equal(t, want, got)

	// sanity: bootstrap fields are NOT set by flattenHCEnv
	_, hasCmds := got[hcEnvBootstrapCommandsFieldName]
	_, hasVer := got[hcEnvBootstrapCommandsVersionFieldName]
	assert.False(t, hasCmds)
	assert.False(t, hasVer)
}

func TestExpandHCEnvForCreate_UsesDefaultAccountID(t *testing.T) {
	defaultAcct := "00000000-1000-0000-0000-000000000001"
	configMap := map[string]interface{}{
		hcEnvCfgNamespaceFieldName:                  "qdrant-hc-new",
		hcEnvCfgHttpProxyUrlFieldName:               "http://proxy.example.com",
		hcEnvCfgHttpsProxyUrlFieldName:              "https://proxy.example.com",
		hcEnvCfgNoProxyConfigsFieldName:             []interface{}{"localhost", "127.0.0.1"},
		hcEnvCfgContainerRegistryUrlFieldName:       "registry.example.com",
		hcEnvCfgChartRepositoryUrlFieldName:         "charts.example.com",
		hcEnvCfgRegistrySecretNameFieldName:         "reg-secret",
		hcEnvCfgCaCertificatesFieldName:             "ca-cert-data",
		hcEnvCfgDatabaseStorageClassFieldName:       "db-storage",
		hcEnvCfgSnapshotStorageClassFieldName:       "snap-storage",
		hcEnvCfgVolumeSnapshotStorageClassFieldName: "vol-snap-storage",
		hcEnvCfgLogLevelFieldName:                   "HYBRID_CLOUD_ENVIRONMENT_CONFIGURATION_LOG_LEVEL_INFO",
		hcEnvCfgAdvancedOperatorSettingsFieldName:   "key: value\nnested:\n  num: 1\n",
		hcEnvCfgNodeSelectorFieldName: []interface{}{
			map[string]interface{}{"key": "key1", "value": "value1"},
		},
		hcEnvCfgTolerationsFieldName: []interface{}{
			map[string]interface{}{
				tolerationKeyFieldName:      "key1",
				tolerationOperatorFieldName: "TOLERATION_OPERATOR_EQUAL",
				tolerationValueFieldName:    "value1",
				tolerationEffectFieldName:   "TOLERATION_EFFECT_NO_SCHEDULE",
			},
		},
		hcEnvCfgControlPlaneLabelsFieldName: []interface{}{
			map[string]interface{}{"key": "label2", "value": "value2"},
		},
	}

	d := schema.TestResourceDataRaw(t, accountsHybridCloudEnvironmentSchema(), map[string]interface{}{
		hcEnvNameFieldName:          "local-test-new",
		hcEnvConfigurationFieldName: []interface{}{configMap},
	})
	d.MarkNewResource()
	env, err := expandHCEnv(d, defaultAcct)
	require.NoError(t, err)

	assert.Equal(t, defaultAcct, env.GetAccountId())
	assert.Equal(t, "local-test-new", env.GetName())
	require.NotNil(t, env.GetConfiguration())
	assert.Equal(t, configMap[hcEnvCfgNamespaceFieldName], env.GetConfiguration().GetNamespace())
	assert.Equal(t, configMap[hcEnvCfgHttpProxyUrlFieldName], env.GetConfiguration().GetHttpProxyUrl())
	assert.Equal(t, configMap[hcEnvCfgHttpsProxyUrlFieldName], env.GetConfiguration().GetHttpsProxyUrl())
	assert.Equal(t, []string{"localhost", "127.0.0.1"}, env.GetConfiguration().GetNoProxyConfigs())
	assert.Equal(t, configMap[hcEnvCfgContainerRegistryUrlFieldName], env.GetConfiguration().GetContainerRegistryUrl())
	assert.Equal(t, configMap[hcEnvCfgChartRepositoryUrlFieldName], env.GetConfiguration().GetChartRepositoryUrl())
	assert.Equal(t, configMap[hcEnvCfgRegistrySecretNameFieldName], env.GetConfiguration().GetRegistrySecretName())
	assert.Equal(t, configMap[hcEnvCfgCaCertificatesFieldName], env.GetConfiguration().GetCaCertificates())
	assert.Equal(t, configMap[hcEnvCfgDatabaseStorageClassFieldName], env.GetConfiguration().GetDatabaseStorageClass())
	assert.Equal(t, configMap[hcEnvCfgSnapshotStorageClassFieldName], env.GetConfiguration().GetSnapshotStorageClass())
	assert.Equal(t, configMap[hcEnvCfgVolumeSnapshotStorageClassFieldName], env.GetConfiguration().GetVolumeSnapshotStorageClass())
	assert.Equal(t, qch.HybridCloudEnvironmentConfigurationLogLevel_HYBRID_CLOUD_ENVIRONMENT_CONFIGURATION_LOG_LEVEL_INFO, env.GetConfiguration().GetLogLevel())

	require.NotNil(t, env.GetConfiguration().GetAdvancedOperatorSettings())
	advMap := env.GetConfiguration().GetAdvancedOperatorSettings().AsMap()
	assert.Equal(t, "value", advMap["key"])
	nested, ok := advMap["nested"].(map[string]interface{})
	require.True(t, ok)
	assert.EqualValues(t, 1, nested["num"])
}

func TestExpandHCEnvForCreate_MissingNamespaceErrors(t *testing.T) {
	d := schema.TestResourceDataRaw(t, accountsHybridCloudEnvironmentSchema(), map[string]interface{}{
		hcEnvNameFieldName:          "local-test-new",
		hcEnvConfigurationFieldName: []interface{}{map[string]interface{}{}}, // no namespace
	})
	d.MarkNewResource() // This is a create operation
	_, err := expandHCEnv(d, "00000000-1000-0000-0000-000000000001")
	require.Error(t, err)
}

func TestExpandHCEnvForUpdate_WithOverrideAccountID(t *testing.T) {
	overrideAcct := "00000000-2000-0000-0000-000000000002"
	configMap := map[string]interface{}{
		hcEnvCfgNamespaceFieldName: "qdrant-hc-update",
	}
	d := schema.TestResourceDataRaw(t, accountsHybridCloudEnvironmentSchema(), nil)
	d.SetId("00000000-0000-0000-0000-0000000000AA") // This makes HasChange work

	// Set the new values
	require.NoError(t, d.Set(hcEnvAccountIdFieldName, overrideAcct))
	require.NoError(t, d.Set(hcEnvNameFieldName, "local-test-update"))
	require.NoError(t, d.Set(hcEnvConfigurationFieldName, []interface{}{configMap}))

	env, err := expandHCEnv(d, "ignored-default")
	require.NoError(t, err)

	assert.Equal(t, overrideAcct, env.GetAccountId())
	assert.Equal(t, "00000000-0000-0000-0000-0000000000AA", env.GetId())
	assert.Equal(t, "local-test-update", env.GetName(), "name should be in the payload")
	require.NotNil(t, env.GetConfiguration())
	assert.Equal(t, "qdrant-hc-update", env.GetConfiguration().GetNamespace(), "configuration should be in the payload")
}

func TestExpandHCEnvForUpdate_NoChanges(t *testing.T) {
	d := schema.TestResourceDataRaw(t, accountsHybridCloudEnvironmentSchema(), map[string]interface{}{})
	d.SetId("00000000-0000-0000-0000-0000000000BB")

	env, err := expandHCEnv(d, "default-account")
	require.NoError(t, err)

	assert.Empty(t, env.GetName(), "name should be empty in payload if not set")
	assert.NotNil(t, env.GetConfiguration(), "configuration should be an empty object if not set")
	assert.Empty(t, env.GetConfiguration().GetNamespace())
}

func TestSchema_BootstrapVersionFlags(t *testing.T) {
	s := accountsHybridCloudEnvironmentSchema()
	field := s[hcEnvBootstrapCommandsVersionFieldName]
	require.NotNil(t, field)

	assert.True(t, field.Optional, "bootstrap_commands_version should be Optional")
	assert.True(t, field.Computed, "bootstrap_commands_version should be Computed")
	assert.Nil(t, field.Default, "bootstrap_commands_version must NOT set a Default")
}

func TestSchema_ConfigurationNamespaceForceNew(t *testing.T) {
	cfg := accountsHybridCloudEnvironmentConfigurationSchema()
	ns := cfg[hcEnvCfgNamespaceFieldName]
	require.NotNil(t, ns)
	assert.True(t, ns.ForceNew, "configuration.namespace must be ForceNew")
}

func TestConvertMapKeysToStrings(t *testing.T) {
	testCases := []struct {
		name  string
		input interface{}
		want  interface{}
	}{
		{
			name:  "nil input",
			input: nil,
			want:  nil,
		},
		{
			name:  "simple map[interface{}]interface{}",
			input: map[interface{}]interface{}{"key": "value", "num": 1},
			want:  map[string]interface{}{"key": "value", "num": 1},
		},
		{
			name: "nested map[interface{}]interface{}",
			input: map[interface{}]interface{}{
				"level1": map[interface{}]interface{}{
					"level2": "value",
				},
			},
			want: map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": "value",
				},
			},
		},
		{
			name: "map[string]interface{} with nested map[interface{}]interface{}",
			input: map[string]interface{}{
				"level1_str": map[interface{}]interface{}{
					"level2": "value",
				},
			},
			want: map[string]interface{}{
				"level1_str": map[string]interface{}{
					"level2": "value",
				},
			},
		},
		{
			name: "slice with nested map[interface{}]interface{}",
			input: []interface{}{
				"item1",
				map[interface{}]interface{}{"key": "value"},
				"item2",
			},
			want: []interface{}{
				"item1",
				map[string]interface{}{"key": "value"},
				"item2",
			},
		},
		{
			name: "complex nested structure",
			input: map[string]interface{}{
				"config": map[interface{}]interface{}{
					"params": []interface{}{
						map[interface{}]interface{}{"id": 1, "enabled": true},
						map[interface{}]interface{}{"id": 2, "enabled": false},
					},
				},
			},
			want: map[string]interface{}{
				"config": map[string]interface{}{
					"params": []interface{}{
						map[string]interface{}{"id": 1, "enabled": true},
						map[string]interface{}{"id": 2, "enabled": false},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := convertMapKeysToStrings(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}
