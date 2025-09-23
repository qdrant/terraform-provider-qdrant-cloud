package qdrant

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	qch "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/hybrid/v1"
)

func TestFlattenHCEnv(t *testing.T) {
	env := &qch.HybridCloudEnvironment{
		AccountId:      "00000000-1000-0000-0000-000000000001",
		Id:             "00000000-0000-0000-0000-000000000001",
		CreatedAt:      timestamppb.New(time.Date(2025, 9, 16, 9, 41, 5, 0, time.UTC)),
		LastModifiedAt: timestamppb.New(time.Date(2025, 9, 16, 9, 41, 8, 0, time.UTC)),
		Name:           "local-test-2",
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
		},
	}
	got := flattenHCEnv(env)
	want := map[string]interface{}{
		hcEnvIdFieldName:             env.GetId(),
		hcEnvAccountIdFieldName:      env.GetAccountId(),
		hcEnvNameFieldName:           env.GetName(),
		hcEnvCreatedAtFieldName:      formatTime(env.GetCreatedAt()),
		hcEnvLastModifiedAtFieldName: formatTime(env.GetLastModifiedAt()),
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
	}

	d := schema.TestResourceDataRaw(t, accountsHybridCloudEnvironmentSchema(), map[string]interface{}{
		hcEnvNameFieldName:          "local-test-new",
		hcEnvConfigurationFieldName: []interface{}{configMap},
	})
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
