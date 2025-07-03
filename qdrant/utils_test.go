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
