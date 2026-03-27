package qdrant

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetRequestID(t *testing.T) {
	t.Run("no request id", func(t *testing.T) {
		md := metadata.MD{}
		assert.Empty(t, getRequestID(md))
	})

	t.Run("with request id", func(t *testing.T) {
		md := metadata.Pairs(requestIDTrailerField, "test-id")
		assert.Equal(t, " [test-id]", getRequestID(md))
	})

	t.Run("with multiple request ids", func(t *testing.T) {
		md := metadata.Pairs(requestIDTrailerField, "test-id-1", requestIDTrailerField, "test-id-2")
		assert.Equal(t, " [test-id-1|test-id-2]", getRequestID(md))
	})
}

func TestGetAccountUUID(t *testing.T) {
	providerConfig := &ProviderConfig{
		AccountID: "00000000-0000-0000-0000-000000000002",
	}

	t.Run("from resource data", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"account_id": {Type: schema.TypeString},
		}, map[string]interface{}{
			"account_id": "00000000-0000-0000-0000-000000000001",
		})

		id, err := getAccountUUID(d, providerConfig)
		require.NoError(t, err)
		assert.Equal(t, "00000000-0000-0000-0000-000000000001", id.String())
	})

	t.Run("from provider default", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"account_id": {Type: schema.TypeString},
		}, map[string]interface{}{})

		id, err := getAccountUUID(d, providerConfig)
		require.NoError(t, err)
		assert.Equal(t, "00000000-0000-0000-0000-000000000002", id.String())
	})

	t.Run("not found", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"account_id": {Type: schema.TypeString},
		}, map[string]interface{}{})

		_, err := getAccountUUID(d, &ProviderConfig{})
		require.Error(t, err)
		assert.Equal(t, "cannot find account ID", err.Error())
	})

	t.Run("invalid uuid", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
			"account_id": {Type: schema.TypeString},
		}, map[string]interface{}{
			"account_id": "not-a-uuid",
		})

		_, err := getAccountUUID(d, providerConfig)
		assert.Error(t, err)
	})
}

func TestGetDefaultAccountID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		config := &ProviderConfig{AccountID: "test-id"}
		assert.Equal(t, "test-id", getDefaultAccountID(config))
	})

	t.Run("invalid type", func(t *testing.T) {
		assert.Empty(t, getDefaultAccountID("not-a-config"))
	})
}

func TestTimeParsing(t *testing.T) {
	rfc3339Time := "2025-07-01T10:30:00Z"
	goTime, _ := time.Parse(time.RFC3339, rfc3339Time)
	protoTime := timestamppb.New(goTime)

	t.Run("parseTime success", func(t *testing.T) {
		parsed := parseTime(rfc3339Time)
		assert.NotNil(t, parsed)
		assert.Equal(t, protoTime.Seconds, parsed.Seconds)
	})

	t.Run("parseTime failure", func(t *testing.T) {
		assert.Nil(t, parseTime("invalid-time"))
	})

	t.Run("formatTime success", func(t *testing.T) {
		assert.Equal(t, rfc3339Time, formatTime(protoTime))
	})

	t.Run("formatTime nil", func(t *testing.T) {
		assert.Empty(t, formatTime(nil))
	})
}

func TestDurationParsing(t *testing.T) {
	durationStr := "72h0m0s"
	goDuration, _ := time.ParseDuration("72h")
	protoDuration := durationpb.New(goDuration)

	t.Run("parseDuration success", func(t *testing.T) {
		parsed := parseDuration("72h")
		assert.NotNil(t, parsed)
		assert.Equal(t, protoDuration.Seconds, parsed.Seconds)
	})

	t.Run("parseDuration failure", func(t *testing.T) {
		assert.Nil(t, parseDuration("invalid-duration"))
	})

	t.Run("formatDuration success", func(t *testing.T) {
		assert.Equal(t, durationStr, formatDuration(protoDuration))
	})

	t.Run("formatDuration nil", func(t *testing.T) {
		assert.Empty(t, formatDuration(nil))
	})
}

func TestGetAccountUUID_NilUUID(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"account_id": {Type: schema.TypeString},
	}, map[string]interface{}{})

	providerConfig := &ProviderConfig{
		AccountID: uuid.Nil.String(),
	}

	id, err := getAccountUUID(d, providerConfig)
	require.NoError(t, err)
	assert.Equal(t, uuid.Nil, id)
}

func TestGetAccountUUID_EmptyString(t *testing.T) {
	d := schema.TestResourceDataRaw(t, map[string]*schema.Schema{
		"account_id": {Type: schema.TypeString},
	}, map[string]interface{}{
		"account_id": "",
	})

	providerConfig := &ProviderConfig{
		AccountID: "00000000-0000-0000-0000-000000000002",
	}

	id, err := getAccountUUID(d, providerConfig)
	require.NoError(t, err)
	assert.Equal(t, "00000000-0000-0000-0000-000000000002", id.String())
}

func TestSuppressDurationDiff(t *testing.T) {
	t.Run("equal durations", func(t *testing.T) {
		assert.True(t, suppressDurationDiff("key", "72h", "72h0m0s", nil))
	})
	t.Run("unequal durations", func(t *testing.T) {
		assert.False(t, suppressDurationDiff("key", "72h", "73h", nil))
	})
	t.Run("invalid old duration", func(t *testing.T) {
		assert.False(t, suppressDurationDiff("key", "invalid", "72h", nil))
	})
	t.Run("invalid new duration", func(t *testing.T) {
		assert.False(t, suppressDurationDiff("key", "72h", "invalid", nil))
	})
}

func TestDiffStringSets_AddAndDel(t *testing.T) {
	desired := []string{"a", "b", "c"}
	current := []string{"b", "d"}

	toAdd, toDel := diffStringSets(desired, current)

	// Order is not guaranteed; compare as sets.
	assert.ElementsMatch(t, []string{"a", "c"}, toAdd)
	assert.ElementsMatch(t, []string{"d"}, toDel)
}

func TestDiffStringSets_Identical(t *testing.T) {
	desired := []string{"a", "b"}
	current := []string{"b", "a"} // different order, same set

	toAdd, toDel := diffStringSets(desired, current)
	assert.Empty(t, toAdd)
	assert.Empty(t, toDel)
}

func TestDiffStringSets_EmptyDesired(t *testing.T) {
	desired := []string{}
	current := []string{"x", "y"}

	toAdd, toDel := diffStringSets(desired, current)
	assert.Empty(t, toAdd)
	assert.ElementsMatch(t, []string{"x", "y"}, toDel)
}

func TestDiffStringSets_EmptyCurrent(t *testing.T) {
	desired := []string{"x", "y"}
	current := []string{}

	toAdd, toDel := diffStringSets(desired, current)
	assert.ElementsMatch(t, []string{"x", "y"}, toAdd)
	assert.Empty(t, toDel)
}

func TestSetToStringSlice_WithSchemaSet(t *testing.T) {
	// A schema.Set cannot contain nil values; HashString would panic.
	// Use duplicates to ensure we still get unique string handling from schema.Set.
	s := schema.NewSet(schema.HashString, []interface{}{"a", "b", "a"})

	out := setToStringSlice(s)

	// Order from Set.List() is not guaranteed; compare as a set.
	require.ElementsMatch(t, []string{"a", "b"}, out)
}

func TestSetToStringSlice_WithInterfaceSlice(t *testing.T) {
	in := []interface{}{"x", nil, "y"}

	got := setToStringSlice(in)
	require.Len(t, got, 2)
	assert.ElementsMatch(t, []string{"x", "y"}, got)
}

func TestSetToStringSlice_UnsupportedType(t *testing.T) {
	got := setToStringSlice(123) // not a *schema.Set or []interface{}
	assert.Nil(t, got)
}

func TestIntersectStrings_Basic(t *testing.T) {
	// Intersection should preserve the order of 'b'
	a := []string{"a", "b", "c"}
	b := []string{"c", "a", "d"}

	got := intersectStrings(a, b)
	assert.Equal(t, []string{"c", "a"}, got)
}

// TestGetInterfaceSliceFromSchemaValue tests the getInterfaceSliceFromSchemaValue function.
func TestGetInterfaceSliceFromSchemaValue(t *testing.T) {
	t.Run("nil input", func(t *testing.T) {
		result := getInterfaceSliceFromSchemaValue(nil)
		assert.Nil(t, result)
	})

	t.Run("schema.Set input", func(t *testing.T) {
		s := schema.NewSet(schema.HashString, []interface{}{"item1", "item2"})
		result := getInterfaceSliceFromSchemaValue(s)
		assert.ElementsMatch(t, []interface{}{"item1", "item2"}, result)
	})

	t.Run("[]interface{} input", func(t *testing.T) {
		l := []interface{}{"itemA", "itemB"}
		result := getInterfaceSliceFromSchemaValue(l)
		assert.Equal(t, l, result)
	})

	t.Run("unsupported input type", func(t *testing.T) {
		result := getInterfaceSliceFromSchemaValue("just a string")
		assert.Nil(t, result)
	})
}

// TestKeyValHashFunc tests the keyValHashFunc function.
func TestKeyValHashFunc(t *testing.T) {
	t.Run("basic key-value pair", func(t *testing.T) {
		kv1 := map[string]interface{}{"key": "name", "value": "test"}
		kv2 := map[string]interface{}{"value": "test", "key": "name"} // different order
		kv3 := map[string]interface{}{"key": "other", "value": "value"}

		hash1 := keyValHashFunc(kv1)
		hash2 := keyValHashFunc(kv2)
		hash3 := keyValHashFunc(kv3)

		assert.Equal(t, hash1, hash2, "identical key-value pairs should have same hash")
		assert.NotEqual(t, hash1, hash3, "different key-value pairs should have different hashes")
	})

	t.Run("empty values", func(t *testing.T) {
		kv1 := map[string]interface{}{"key": "empty", "value": ""}
		kv2 := map[string]interface{}{"key": "empty", "value": ""}
		hash1 := keyValHashFunc(kv1)
		hash2 := keyValHashFunc(kv2)
		assert.Equal(t, hash1, hash2)
	})
}

// TestTolerationHashFunc tests the tolerationHashFunc function.
func TestTolerationHashFunc(t *testing.T) {
	t.Run("toleration with all fields", func(t *testing.T) {
		tol1 := map[string]interface{}{
			tolerationKeyFieldName:      "key1",
			tolerationOperatorFieldName: "Equal",
			tolerationValueFieldName:    "value1",
			tolerationEffectFieldName:   "NoSchedule",
			tolerationSecondsFieldName:  300,
		}
		tol2 := map[string]interface{}{
			tolerationKeyFieldName:      "key1",
			tolerationOperatorFieldName: "Equal",
			tolerationValueFieldName:    "value1",
			tolerationEffectFieldName:   "NoSchedule",
			tolerationSecondsFieldName:  300,
		}
		tol3 := map[string]interface{}{
			tolerationKeyFieldName:      "key2",
			tolerationOperatorFieldName: "Exists",
			tolerationEffectFieldName:   "NoExecute",
		}

		hash1 := tolerationHashFunc(tol1)
		hash2 := tolerationHashFunc(tol2)
		hash3 := tolerationHashFunc(tol3)

		assert.Equal(t, hash1, hash2, "identical tolerations should have same hash")
		assert.NotEqual(t, hash1, hash3, "different tolerations should have different hashes")
	})

	t.Run("toleration with partial fields", func(t *testing.T) {
		tol1 := map[string]interface{}{tolerationKeyFieldName: "key1", tolerationEffectFieldName: "NoSchedule"}
		tol2 := map[string]interface{}{tolerationKeyFieldName: "key1", tolerationEffectFieldName: "NoSchedule"}
		hash1 := tolerationHashFunc(tol1)
		hash2 := tolerationHashFunc(tol2)
		assert.Equal(t, hash1, hash2)
	})
}

// TestTopologySpreadConstraintHashFunc tests the topologySpreadConstraintHashFunc function.
func TestTopologySpreadConstraintHashFunc(t *testing.T) {
	tsc1 := map[string]interface{}{
		topologySpreadConstraintMaxSkewFieldName:           1,
		topologySpreadConstraintTopologyKeyFieldName:       "zone",
		topologySpreadConstraintWhenUnsatisfiableFieldName: "DoNotSchedule",
	}
	tsc2 := map[string]interface{}{
		topologySpreadConstraintWhenUnsatisfiableFieldName: "DoNotSchedule",
		topologySpreadConstraintTopologyKeyFieldName:       "zone",
		topologySpreadConstraintMaxSkewFieldName:           1,
	}
	tsc3 := map[string]interface{}{
		topologySpreadConstraintMaxSkewFieldName:           2,
		topologySpreadConstraintTopologyKeyFieldName:       "node",
		topologySpreadConstraintWhenUnsatisfiableFieldName: "ScheduleAnyway",
	}

	hash1 := topologySpreadConstraintHashFunc(tsc1)
	hash2 := topologySpreadConstraintHashFunc(tsc2)
	hash3 := topologySpreadConstraintHashFunc(tsc3)

	assert.Equal(t, hash1, hash2, "identical constraints should have same hash")
	assert.NotEqual(t, hash1, hash3, "different constraints should have different hashes")
}

// TestPermissionHashFunc tests the permissionHashFunc function.
func TestPermissionHashFunc(t *testing.T) {
	perm1 := map[string]interface{}{"value": "read:clusters"}
	perm2 := map[string]interface{}{"value": "read:clusters"}
	perm3 := map[string]interface{}{"value": "write:clusters"}

	hash1 := permissionHashFunc(perm1)
	hash2 := permissionHashFunc(perm2)
	hash3 := permissionHashFunc(perm3)

	assert.Equal(t, hash1, hash2, "identical permissions should have same hash")
	assert.NotEqual(t, hash1, hash3, "different permissions should have different hashes")
}

func TestIntersectStrings_Overlap(t *testing.T) {
	a := []string{"a", "b", "c"}
	b := []string{"c", "d", "e"}

	got := intersectStrings(a, b)
	assert.Equal(t, []string{"c"}, got)
}
func TestIntersectStrings_NoOverlap(t *testing.T) {
	a := []string{"a", "b"}
	b := []string{"c", "d"}

	got := intersectStrings(a, b)
	assert.Empty(t, got)
}
