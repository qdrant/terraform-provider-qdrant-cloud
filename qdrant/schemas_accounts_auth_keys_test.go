package qdrant

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	qc "github.com/qdrant/terraform-provider-qdrant-cloud/v1/internal/client"
)

func TestFlattenAuthKeySchema(t *testing.T) {
	createdAt := time.Now()
	keys := []qc.ApiKeySchema{
		{
			Id:         newPointer(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
			CreatedAt:  newPointer(createdAt),
			ClusterIds: []uuid.UUID{uuid.MustParse("00000000-0000-0000-0001-000000000001"), uuid.MustParse("00000000-0000-0000-0002-000000000001")},
			Prefix:     newPointer("prefix1"),
			Token:      newPointer("token1"),
		},
		{
			Id:         newPointer(uuid.MustParse("00000000-0000-0000-0000-000000000002")),
			CreatedAt:  newPointer(createdAt),
			ClusterIds: []uuid.UUID{uuid.MustParse("00000000-0000-0000-0003-000000000002")},
			Prefix:     newPointer("prefix2"),
			Token:      newPointer("token2"),
		},
	}

	expected := []interface{}{
		map[string]interface{}{
			authKeysKeysIDFieldName:         "00000000-0000-0000-0000-000000000001",
			authKeysKeysCreatedAtFieldName:  formatTime(createdAt),
			authKeysKeysClusterIDsFieldName: []string{"00000000-0000-0000-0001-000000000001", "00000000-0000-0000-0002-000000000001"},
			authKeysKeysPrefixFieldName:     "prefix1",
			authKeysKeysTokenFieldName:      "token1",
		},
		map[string]interface{}{
			authKeysKeysIDFieldName:         "00000000-0000-0000-0000-000000000002",
			authKeysKeysCreatedAtFieldName:  formatTime(createdAt),
			authKeysKeysClusterIDsFieldName: []string{"00000000-0000-0000-0003-000000000002"},
			authKeysKeysPrefixFieldName:     "prefix2",
			authKeysKeysTokenFieldName:      "token2",
		},
	}

	flattened := flattenAuthKeys(keys)
	assert.Equal(t, expected, flattened)
}

func TestFlattenCreateAuthKey(t *testing.T) {
	createdAt := time.Now()
	key := qc.ApiKeySchema{
		Id:         newPointer(uuid.MustParse("10000000-0000-0000-0000-000000000002")),
		CreatedAt:  newPointer(createdAt),
		ClusterIds: []uuid.UUID{uuid.MustParse("10000000-0000-0000-0003-000000000002")},
		Prefix:     newPointer("prefix3"),
		Token:      newPointer("token3"),
	}

	expected := map[string]interface{}{
		authKeysKeysIDFieldName:         "10000000-0000-0000-0000-000000000002",
		authKeysKeysCreatedAtFieldName:  formatTime(createdAt),
		authKeysKeysClusterIDsFieldName: []string{"10000000-0000-0000-0003-000000000002"},
		authKeysKeysPrefixFieldName:     "prefix3",
		authKeysKeysTokenFieldName:      "token3",
	}

	assert.Equal(t, expected, flattenAuthKey(key))
}
