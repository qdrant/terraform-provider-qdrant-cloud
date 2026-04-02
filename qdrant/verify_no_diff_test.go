package qdrant

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	qcCluster "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/v1"
)

// TestVerifyNoPerpetualDiffWithRealCluster connects to a real Qdrant Cloud cluster,
// fetches its configuration, flattens it, and verifies that no phantom diffs would
// occur when the user hasn't set optional enum fields.
//
// This is the end-to-end proof that CP-393 is fixed:
//  1. Fetch real cluster from the API (which returns UNSPECIFIED for unset enums)
//  2. Flatten the response (simulating what the provider Read function does)
//  3. Simulate an empty user config (user didn't set the enum fields)
//  4. Compare: if the flatten output doesn't include UNSPECIFIED values,
//     AND the schema has Computed: true, Terraform won't show a diff.
func TestVerifyNoPerpetualDiffWithRealCluster(t *testing.T) {
	apiKey := os.Getenv("QDRANT_CLOUD_API_KEY")
	accountID := os.Getenv("QDRANT_CLOUD_ACCOUNT_ID")
	apiURL := os.Getenv("QDRANT_CLOUD_API_URL")
	if apiKey == "" || accountID == "" {
		t.Skip("QDRANT_CLOUD_API_KEY and QDRANT_CLOUD_ACCOUNT_ID required")
	}
	if apiURL == "" {
		apiURL = "grpc.development-cloud.qdrant.io"
	}

	creds := credentials.NewTLS(&tls.Config{})
	conn, err := grpc.NewClient(apiURL, grpc.WithTransportCredentials(creds))
	require.NoError(t, err)
	defer func() { _ = conn.Close() }()

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", fmt.Sprintf("apikey %s", apiKey))
	client := qcCluster.NewClusterServiceClient(conn)

	resp, err := client.ListClusters(ctx, &qcCluster.ListClustersRequest{AccountId: accountID})
	require.NoError(t, err)
	require.NotEmpty(t, resp.GetItems(), "expected at least one cluster in the account")

	cluster := resp.GetItems()[0]
	cfg := cluster.GetConfiguration()
	require.NotNil(t, cfg, "cluster should have a configuration")

	t.Logf("Cluster: %s (%s)", cluster.GetName(), cluster.GetId())
	t.Logf("Raw API values:")
	t.Logf("  ServiceType:       %s (val=%d)", cfg.GetServiceType(), int(cfg.GetServiceType()))
	t.Logf("  GpuType:           %s (val=%d)", cfg.GetGpuType(), int(cfg.GetGpuType()))
	t.Logf("  RestartPolicy:     %s (val=%d)", cfg.GetRestartPolicy(), int(cfg.GetRestartPolicy()))
	t.Logf("  RebalanceStrategy: %s (val=%d)", cfg.GetRebalanceStrategy(), int(cfg.GetRebalanceStrategy()))
	t.Logf("  ReservedCpuPct:    %v", cfg.ReservedCpuPercentage)
	t.Logf("  ReservedMemPct:    %v", cfg.ReservedMemoryPercentage)

	flattened := flattenCluster(cluster)
	configList := flattened[configurationFieldName].([]interface{})
	require.Len(t, configList, 1)
	configMap := configList[0].(map[string]interface{})

	enumFields := []struct {
		fieldName string
		apiValue  int
	}{
		{serviceTypeFieldName, int(cfg.GetServiceType())},
		{dbConfigGpuTypeFieldName, int(cfg.GetGpuType())},
		{dbConfigRestartPolicyFieldName, int(cfg.GetRestartPolicy())},
		{dbConfigRebalanceStrategyFieldName, int(cfg.GetRebalanceStrategy())},
	}

	for _, ef := range enumFields {
		val, present := configMap[ef.fieldName]
		if ef.apiValue == 0 {
			// UNSPECIFIED (val=0) — must NOT be in flattened output
			assert.False(t, present,
				"field %s: backend returned UNSPECIFIED but flatten included it in state — this would cause a perpetual diff",
				ef.fieldName)
			t.Logf("  ✓ %s: UNSPECIFIED correctly excluded from state", ef.fieldName)
		} else {
			// Non-UNSPECIFIED — must be present with correct value
			assert.True(t, present,
				"field %s: backend returned a real value but flatten excluded it",
				ef.fieldName)
			t.Logf("  ✓ %s: value %q correctly included in state", ef.fieldName, val)
		}
	}

	ptrFields := []struct {
		fieldName string
		isNil     bool
	}{
		{dbConfigReservedCpuPercentageFieldName, cfg.ReservedCpuPercentage == nil},
		{dbConfigReservedMemoryPercentageFieldName, cfg.ReservedMemoryPercentage == nil},
	}

	for _, pf := range ptrFields {
		_, present := configMap[pf.fieldName]
		if pf.isNil {
			assert.False(t, present,
				"field %s: backend returned nil but flatten included it — this would cause a perpetual diff",
				pf.fieldName)
			t.Logf("  ✓ %s: nil correctly excluded from state", pf.fieldName)
		} else {
			assert.True(t, present,
				"field %s: backend returned a value but flatten excluded it",
				pf.fieldName)
			t.Logf("  ✓ %s: value correctly included in state", pf.fieldName)
		}
	}

	schemaMap := accountsClusterConfigurationSchema(false) // resource mode
	computedFields := []string{
		serviceTypeFieldName,
		dbConfigGpuTypeFieldName,
		dbConfigRestartPolicyFieldName,
		dbConfigRebalanceStrategyFieldName,
		dbConfigReservedCpuPercentageFieldName,
		dbConfigReservedMemoryPercentageFieldName,
	}
	for _, field := range computedFields {
		s := schemaMap[field]
		require.NotNil(t, s, "field %s missing from schema", field)
		assert.True(t, s.Optional, "field %s must be Optional", field)
		assert.True(t, s.Computed, "field %s must be Computed — without this, Terraform shows perpetual diffs", field)
	}
	t.Log("  ✓ All 6 fields are Optional+Computed in resource schema")

	t.Log("")
	t.Log("No perpetual diffs: UNSPECIFIED values excluded from state, schema marks fields as Computed.")
}
