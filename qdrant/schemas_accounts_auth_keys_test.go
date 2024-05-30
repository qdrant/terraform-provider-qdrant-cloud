package qdrant

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	qc "terraform-provider-qdrant-cloud/v1/internal/client"
)

func TestFlattenGetAuthKeys(t *testing.T) {
	createdAt := time.Now()
	keys := []qc.GetApiKeyOut{
		{
			Id:            "testID1",
			CreatedAt:     createdAt,
			UserId:        newString("testUserID1"),
			AccountId:     newString("testAccountID1"),
			ClusterIdList: &[]string{"cluster1", "cluster2"},
			Prefix:        "prefix1",
		},
		{
			Id:            "testID2",
			CreatedAt:     createdAt,
			UserId:        newString("testUserID2"),
			AccountId:     newString("testAccountID2"),
			ClusterIdList: &[]string{"cluster3"},
			Prefix:        "prefix2",
		},
	}

	expected := []interface{}{
		map[string]interface{}{
			authKeysKeysIDFieldName:         "testID1",
			authKeysKeysCreatedAtFieldName:  formatTime(createdAt),
			authKeysKeysUserIDFieldName:     "testUserID1",
			authKeysKeysAccountIDFieldName:  "testAccountID1",
			authKeysKeysClusterIDsFieldName: []string{"cluster1", "cluster2"},
			authKeysKeysPrefixFieldName:     "prefix1",
		},
		map[string]interface{}{
			authKeysKeysIDFieldName:         "testID2",
			authKeysKeysCreatedAtFieldName:  formatTime(createdAt),
			authKeysKeysUserIDFieldName:     "testUserID2",
			authKeysKeysAccountIDFieldName:  "testAccountID2",
			authKeysKeysClusterIDsFieldName: []string{"cluster3"},
			authKeysKeysPrefixFieldName:     "prefix2",
		},
	}

	flattened := flattenGetAuthKeys(keys)
	assert.Equal(t, expected, flattened)
}

func TestFlattenGetAuthKey(t *testing.T) {
	createdAt := time.Now()
	key := qc.GetApiKeyOut{
		Id:            "testID",
		CreatedAt:     createdAt,
		UserId:        newString("testUserID"),
		AccountId:     newString("testAccountID"),
		ClusterIdList: &[]string{"cluster1", "cluster2"},
		Prefix:        "prefix",
	}

	expected := map[string]interface{}{
		authKeysKeysIDFieldName:         "testID",
		authKeysKeysCreatedAtFieldName:  formatTime(createdAt),
		authKeysKeysUserIDFieldName:     "testUserID",
		authKeysKeysAccountIDFieldName:  "testAccountID",
		authKeysKeysClusterIDsFieldName: []string{"cluster1", "cluster2"},
		authKeysKeysPrefixFieldName:     "prefix",
	}

	assert.Equal(t, expected, flattenGetAuthKey(key))
}

func TestFlattenCreateAuthKey(t *testing.T) {
	createdAt := time.Now()
	key := qc.CreateApiKeyOut{
		Id:            "testID",
		CreatedAt:     createdAt,
		UserId:        newString("testUserID"),
		AccountId:     newString("testAccountID"),
		ClusterIdList: &[]string{"cluster1", "cluster2"},
		Prefix:        "prefix",
		Token:         "testToken",
	}

	expected := map[string]interface{}{
		authKeysKeysIDFieldName:         "testID",
		authKeysKeysCreatedAtFieldName:  formatTime(createdAt),
		authKeysKeysUserIDFieldName:     "testUserID",
		authKeysKeysAccountIDFieldName:  "testAccountID",
		authKeysKeysClusterIDsFieldName: []string{"cluster1", "cluster2"},
		authKeysKeysPrefixFieldName:     "prefix",
		authKeysKeysTokenFieldName:      "testToken",
	}

	assert.Equal(t, expected, flattenCreateAuthKey(key))
}
