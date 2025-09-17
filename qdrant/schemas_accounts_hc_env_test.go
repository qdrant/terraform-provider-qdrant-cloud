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
			Namespace: "qdrant-hc",
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
				hcEnvCfgNamespaceFieldName: env.GetConfiguration().GetNamespace(),
			},
		},
	}

	assert.Equal(t, want, got)
}

func TestExpandHCEnvForCreate_UsesDefaultAccountID(t *testing.T) {
	defaultAcct := "00000000-1000-0000-0000-000000000001"

	d := schema.TestResourceDataRaw(t, accountsHybridCloudEnvironmentSchema(), map[string]interface{}{
		hcEnvNameFieldName: "local-test-new",
		hcEnvConfigurationFieldName: []interface{}{
			map[string]interface{}{
				hcEnvCfgNamespaceFieldName: "qdrant-hc-new",
			},
		},
	})
	env, err := expandHCEnvForCreate(d, defaultAcct)
	require.NoError(t, err)

	assert.Equal(t, defaultAcct, env.GetAccountId())
	assert.Equal(t, "local-test-new", env.GetName())
	require.NotNil(t, env.GetConfiguration())
	assert.Equal(t, "qdrant-hc-new", env.GetConfiguration().GetNamespace())
}

func TestExpandHCEnvForCreate_MissingNamespaceErrors(t *testing.T) {
	d := schema.TestResourceDataRaw(t, accountsHybridCloudEnvironmentSchema(), map[string]interface{}{
		hcEnvNameFieldName:          "local-test-new",
		hcEnvConfigurationFieldName: []interface{}{map[string]interface{}{}}, // no namespace
	})
	_, err := expandHCEnvForCreate(d, "00000000-1000-0000-0000-000000000001")
	require.Error(t, err)
}

func TestExpandHCEnvForUpdate_WithOverrideAccountID(t *testing.T) {
	overrideAcct := "00000000-2000-0000-0000-000000000002"

	d := schema.TestResourceDataRaw(t, accountsHybridCloudEnvironmentSchema(), map[string]interface{}{
		hcEnvAccountIdFieldName: overrideAcct,
		hcEnvNameFieldName:      "local-test-update",
		hcEnvConfigurationFieldName: []interface{}{
			map[string]interface{}{
				hcEnvCfgNamespaceFieldName: "qdrant-hc-update",
			},
		},
	})
	d.SetId("00000000-0000-0000-0000-0000000000AA")

	env, err := expandHCEnvForUpdate(d, "ignored-default")
	require.NoError(t, err)

	assert.Equal(t, overrideAcct, env.GetAccountId())
	assert.Equal(t, "00000000-0000-0000-0000-0000000000AA", env.GetId())
	assert.Equal(t, "local-test-update", env.GetName())
	require.NotNil(t, env.GetConfiguration())
	assert.Equal(t, "qdrant-hc-update", env.GetConfiguration().GetNamespace())
}
