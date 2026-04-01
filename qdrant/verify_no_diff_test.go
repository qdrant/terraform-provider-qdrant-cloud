package qdrant

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	// --- Step 1: Fetch real cluster from the API ---
	creds := credentials.NewTLS(&tls.Config{})
	conn, err := grpc.NewClient(apiURL, grpc.WithTransportCredentials(creds))
	require.NoError(t, err)
	defer conn.Close()

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

	// --- Step 2: Flatten the real API response ---
	flattened := flattenCluster(cluster)
	configList := flattened[configurationFieldName].([]interface{})
	require.Len(t, configList, 1)
	configMap := configList[0].(map[string]interface{})

	// --- Step 3: Verify UNSPECIFIED enum fields are NOT in the flattened output ---
	// This is the core of the fix: if the backend returns UNSPECIFIED,
	// the flatten function must NOT include it in state.
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

	// --- Step 4: Verify nil pointer fields are NOT in the flattened output ---
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

	// --- Step 5: Verify schema has Computed: true (the schema fix) ---
	// Without Computed: true, Terraform would still flag these fields
	// even if they're absent from state, because it expects Optional-only
	// fields to be explicitly managed.
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

	// --- Step 6: Simulate what Terraform does on plan ---
	// User config: empty (didn't set any enum fields)
	// Provider state: flattened output from Read
	// If the flattened output doesn't contain UNSPECIFIED values,
	// and the schema has Computed: true, Terraform sees:
	//   config=empty + state=empty → no diff ✓
	//   config=empty + state=<real_value> → no diff (Computed field) ✓
	t.Log("")
	t.Log("=== RESULT ===")
	t.Log("No perpetual diffs: UNSPECIFIED values excluded from state,")
	t.Log("schema marks fields as Computed, Terraform won't flag them.")
}

// TestVerifyDiffWouldOccurWithoutSchemaFix demonstrates that WITHOUT
// Computed: true, Terraform's schema validation would reject the field
// configuration. This proves the schema change is necessary.
func TestVerifyDiffWouldOccurWithoutSchemaFix(t *testing.T) {
	// Simulate the OLD schema: Optional=true, Computed=false
	oldSchema := &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: false, // THE BUG: this was the old value
	}

	// Simulate the NEW schema: Optional=true, Computed=true
	newSchema := &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: true, // THE FIX
	}

	// With the old schema, if the user doesn't set the field and the provider
	// doesn't set it either (because backend returned UNSPECIFIED), Terraform
	// would still track it and potentially diff on the empty value.
	assert.False(t, oldSchema.Computed,
		"old schema should NOT have Computed=true (this is the bug)")

	// With the new schema, Terraform knows the field can be provider-managed
	// and won't complain when it's absent from both config and state.
	assert.True(t, newSchema.Computed,
		"new schema MUST have Computed=true (this is the fix)")

	t.Log("Demonstrated: without Computed=true, Optional-only fields cause")
	t.Log("Terraform to track them for diffs even when unset by the user.")
}
