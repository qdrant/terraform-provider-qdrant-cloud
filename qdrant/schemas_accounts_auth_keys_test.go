package qdrant

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	qcAuth "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/auth/v2"
)

func TestFlattenAuthKeySchema(t *testing.T) {
	createdAt := timestamppb.New(time.Now())
	keys := []*qcAuth.ApiKey{
		{
			Id:         "00000000-0000-0000-0000-000000000001",
			CreatedAt:  createdAt,
			ClusterIds: []string{"00000000-0000-0000-0001-000000000001", "00000000-0000-0000-0002-000000000001"},
			Prefix:     "prefix1",
			Token:      "token1",
		},
		{
			Id:         "00000000-0000-0000-0000-000000000002",
			CreatedAt:  createdAt,
			ClusterIds: []string{"00000000-0000-0000-0003-000000000002"},
			Prefix:     "prefix2",
			Token:      "token2",
		},
	}

	expected := []interface{}{
		map[string]interface{}{
			authKeysKeysIDFieldName:         "00000000-0000-0000-0000-000000000001",
			authKeysKeysCreatedAtFieldName:  createdAt.String(),
			authKeysKeysClusterIDsFieldName: []string{"00000000-0000-0000-0001-000000000001", "00000000-0000-0000-0002-000000000001"},
			authKeysKeysPrefixFieldName:     "prefix1",
			authKeysKeysTokenFieldName:      "token1",
		},
		map[string]interface{}{
			authKeysKeysIDFieldName:         "00000000-0000-0000-0000-000000000002",
			authKeysKeysCreatedAtFieldName:  createdAt.String(),
			authKeysKeysClusterIDsFieldName: []string{"00000000-0000-0000-0003-000000000002"},
			authKeysKeysPrefixFieldName:     "prefix2",
			authKeysKeysTokenFieldName:      "token2",
		},
	}

	flattened := flattenAuthKeys(keys)
	assert.Equal(t, expected, flattened)
}

func TestFlattenCreateAuthKey(t *testing.T) {
	createdAt := timestamppb.New(time.Now())
	key := &qcAuth.ApiKey{
		Id:         "10000000-0000-0000-0000-000000000002",
		CreatedAt:  createdAt,
		ClusterIds: []string{"10000000-0000-0000-0003-000000000002"},
		Prefix:     "prefix3",
		Token:      "token3",
	}

	expected := map[string]interface{}{
		authKeysKeysIDFieldName:         "10000000-0000-0000-0000-000000000002",
		authKeysKeysCreatedAtFieldName:  createdAt.String(),
		authKeysKeysClusterIDsFieldName: []string{"10000000-0000-0000-0003-000000000002"},
		authKeysKeysPrefixFieldName:     "prefix3",
		authKeysKeysTokenFieldName:      "token3",
	}

	assert.Equal(t, expected, flattenAuthKey(key))
}
