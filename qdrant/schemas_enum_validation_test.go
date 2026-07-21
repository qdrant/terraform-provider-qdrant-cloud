package qdrant

import (
	"sort"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	authv2 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/auth/v2"
	qcCluster "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/v1"
	commonv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/common/v1"
)

func TestAccountsClusterConfigurationSchema(t *testing.T) {
	resource := accountsClusterConfigurationSchema(false)
	dataSource := accountsClusterConfigurationSchema(true)

	t.Run(serviceTypeFieldName, func(t *testing.T) {
		testGeneratedEnumValidation(t, serviceTypeFieldName, resource[serviceTypeFieldName], dataSource[serviceTypeFieldName], qcCluster.ClusterServiceType_name)
	})
	t.Run(dbConfigGpuTypeFieldName, func(t *testing.T) {
		testGeneratedEnumValidation(t, dbConfigGpuTypeFieldName, resource[dbConfigGpuTypeFieldName], dataSource[dbConfigGpuTypeFieldName], qcCluster.ClusterConfigurationGpuType_name)
	})
	t.Run(dbConfigRestartPolicyFieldName, func(t *testing.T) {
		testGeneratedEnumValidation(t, dbConfigRestartPolicyFieldName, resource[dbConfigRestartPolicyFieldName], dataSource[dbConfigRestartPolicyFieldName], qcCluster.ClusterConfigurationRestartPolicy_name)
	})
	t.Run(dbConfigRebalanceStrategyFieldName, func(t *testing.T) {
		testGeneratedEnumValidation(t, dbConfigRebalanceStrategyFieldName, resource[dbConfigRebalanceStrategyFieldName], dataSource[dbConfigRebalanceStrategyFieldName], qcCluster.ClusterConfigurationRebalanceStrategy_name)
	})
}

func TestDatabaseConfigurationSchema(t *testing.T) {
	resource := databaseConfigurationSchema(false)
	dataSource := databaseConfigurationSchema(true)

	t.Run(dbConfigLogLevelFieldName, func(t *testing.T) {
		testGeneratedEnumValidation(t, dbConfigLogLevelFieldName, resource[dbConfigLogLevelFieldName], dataSource[dbConfigLogLevelFieldName], qcCluster.DatabaseConfigurationLogLevel_name)
	})
}

func TestDatabaseConfigurationAuditLoggingSchema(t *testing.T) {
	resource := databaseConfigurationAuditLoggingSchema(false)
	dataSource := databaseConfigurationAuditLoggingSchema(true)

	t.Run(dbConfigAuditLoggingRotationFieldName, func(t *testing.T) {
		testGeneratedEnumValidation(t, dbConfigAuditLoggingRotationFieldName, resource[dbConfigAuditLoggingRotationFieldName], dataSource[dbConfigAuditLoggingRotationFieldName], qcCluster.AuditLogRotation_name)
	})
}

func TestClusterStorageConfigurationSchema(t *testing.T) {
	resource := clusterStorageConfigurationSchema(false)
	dataSource := clusterStorageConfigurationSchema(true)

	t.Run(clusterStorageTierTypeFieldName, func(t *testing.T) {
		testGeneratedEnumValidation(t, clusterStorageTierTypeFieldName, resource[clusterStorageTierTypeFieldName], dataSource[clusterStorageTierTypeFieldName], commonv1.StorageTierType_name)
	})
}

func TestTolerationSchema(t *testing.T) {
	resource := tolerationSchema(false)
	dataSource := tolerationSchema(true)

	t.Run(tolerationOperatorFieldName, func(t *testing.T) {
		testGeneratedEnumValidation(t, tolerationOperatorFieldName, resource[tolerationOperatorFieldName], dataSource[tolerationOperatorFieldName], qcCluster.TolerationOperator_name)
		t.Run("when a Kubernetes short form is provided, it rejects it", func(t *testing.T) {
			assert.NotEmpty(t, resource[tolerationOperatorFieldName].ValidateDiagFunc("Exists", cty.GetAttrPath(tolerationOperatorFieldName)))
		})
	})
	t.Run(tolerationEffectFieldName, func(t *testing.T) {
		testGeneratedEnumValidation(t, tolerationEffectFieldName, resource[tolerationEffectFieldName], dataSource[tolerationEffectFieldName], qcCluster.TolerationEffect_name)
		t.Run("when a Kubernetes short form is provided, it rejects it", func(t *testing.T) {
			assert.NotEmpty(t, resource[tolerationEffectFieldName].ValidateDiagFunc("NoSchedule", cty.GetAttrPath(tolerationEffectFieldName)))
		})
	})
}

func TestTopologySpreadConstraintSchema(t *testing.T) {
	resource := topologySpreadConstraintSchema(false)
	dataSource := topologySpreadConstraintSchema(true)

	t.Run(topologySpreadConstraintWhenUnsatisfiableFieldName, func(t *testing.T) {
		testGeneratedEnumValidation(t, topologySpreadConstraintWhenUnsatisfiableFieldName, resource[topologySpreadConstraintWhenUnsatisfiableFieldName], dataSource[topologySpreadConstraintWhenUnsatisfiableFieldName], commonv1.TopologySpreadConstraintWhenUnsatisfiable_name)
	})
}

func TestGlobalAccessRuleSchema(t *testing.T) {
	resource := globalAccessRuleSchema(false)
	dataSource := globalAccessRuleSchema(true)

	t.Run(authKeysV2AccessTypeFieldName, func(t *testing.T) {
		testGeneratedEnumValidation(t, authKeysV2AccessTypeFieldName, resource[authKeysV2AccessTypeFieldName], dataSource[authKeysV2AccessTypeFieldName], authv2.GlobalAccessRuleAccessType_name)
	})
}

func TestCollectionAccessRuleSchema(t *testing.T) {
	resource := collectionAccessRuleSchema(false)
	dataSource := collectionAccessRuleSchema(true)

	t.Run(authKeysV2AccessTypeFieldName, func(t *testing.T) {
		testGeneratedEnumValidation(t, authKeysV2AccessTypeFieldName, resource[authKeysV2AccessTypeFieldName], dataSource[authKeysV2AccessTypeFieldName], authv2.CollectionAccessRuleAccessType_name)
	})
}

func testGeneratedEnumValidation(t *testing.T, field string, resource, computedOnly *schema.Schema, generatedNames map[int32]string) {
	t.Helper()
	require.NotNil(t, resource)
	require.NotNil(t, resource.ValidateDiagFunc)

	t.Run("when used for resource input, it uses diagnostic validation", func(t *testing.T) {
		assert.Nil(t, resource.ValidateFunc)
		assert.NotNil(t, resource.ValidateDiagFunc)
	})
	t.Run("when enum names are derived, it excludes the sentinel and sorts the remaining names", func(t *testing.T) {
		validNames := protoEnumNames(generatedNames)
		assert.True(t, sort.StringsAreSorted(validNames))
		assert.Len(t, validNames, len(generatedNames)-1)
	})
	t.Run("when an enum value is provided, it accepts defined values and rejects invalid values", func(t *testing.T) {
		for number, name := range generatedNames {
			diagnostics := resource.ValidateDiagFunc(name, cty.GetAttrPath(field))
			if number == 0 {
				assert.NotEmpty(t, diagnostics, name)
			} else {
				assert.Empty(t, diagnostics, name)
			}
		}
		assert.NotEmpty(t, resource.ValidateDiagFunc("NOT_A_PROTO_ENUM_VALUE", cty.GetAttrPath(field)))
	})

	if computedOnly != nil {
		t.Run("when used for computed-only data, it omits validation", func(t *testing.T) {
			assert.Nil(t, computedOnly.ValidateFunc)
			assert.Nil(t, computedOnly.ValidateDiagFunc)
		})
	}
}
