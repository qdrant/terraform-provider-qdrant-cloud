package qdrant

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	authv2 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/auth/v2"
)

func TestFlattenAuthKeyV2Schema(t *testing.T) {
	createdAt := timestamppb.New(time.Now())
	expiresAt := timestamppb.New(time.Now().Add(24 * time.Hour))

	keys := []*authv2.DatabaseApiKey{
		{
			Id:             "key-id-1",
			AccountId:      "account-id-1",
			ClusterId:      "cluster-id-1",
			Name:           "global-key",
			CreatedAt:      createdAt,
			ExpiresAt:      expiresAt,
			CreatedByEmail: "user1@example.com",
			Postfix:        "post1",
			Key:            "key1",
			AccessRules: []*authv2.AccessRule{
				{
					Scope: &authv2.AccessRule_GlobalAccess{
						GlobalAccess: &authv2.GlobalAccessRule{
							AccessType: authv2.GlobalAccessRuleAccessType_GLOBAL_ACCESS_RULE_ACCESS_TYPE_MANAGE,
						},
					},
				},
			},
		},
		{
			Id:             "key-id-2",
			AccountId:      "account-id-1",
			ClusterId:      "cluster-id-1",
			Name:           "collection-key",
			CreatedAt:      createdAt,
			ExpiresAt:      nil,
			CreatedByEmail: "user2@example.com",
			Postfix:        "post2",
			Key:            "key2",
			AccessRules: []*authv2.AccessRule{
				{
					Scope: &authv2.AccessRule_CollectionAccess{
						CollectionAccess: &authv2.CollectionAccessRule{
							CollectionName: "collection1",
							AccessType:     authv2.CollectionAccessRuleAccessType_COLLECTION_ACCESS_RULE_ACCESS_TYPE_READ_WRITE,
						},
					},
				},
				{
					Scope: &authv2.AccessRule_CollectionAccess{
						CollectionAccess: &authv2.CollectionAccessRule{
							CollectionName: "collection2",
							AccessType:     authv2.CollectionAccessRuleAccessType_COLLECTION_ACCESS_RULE_ACCESS_TYPE_READ_ONLY,
							Payload: map[string]string{
								"filter_key": "filter_value",
							},
						},
					},
				},
			},
		},
	}

	expected := []interface{}{
		map[string]interface{}{
			authKeysV2IDFieldName:             "key-id-1",
			authKeysV2AccountIDFieldName:      "account-id-1",
			authKeysV2ClusterIDFieldName:      "cluster-id-1",
			authKeysV2NameFieldName:           "global-key",
			authKeysV2CreatedAtFieldName:      formatTime(createdAt),
			authKeysV2ExpiresAtFieldName:      formatTime(expiresAt),
			authKeysV2CreatedByEmailFieldName: "user1@example.com",
			authKeysV2PostfixFieldName:        "post1",
			authKeysV2KeyFieldName:            "key1",
			authKeysV2GlobalAccessRuleFieldName: []interface{}{
				map[string]interface{}{
					authKeysV2AccessTypeFieldName: "GLOBAL_ACCESS_RULE_ACCESS_TYPE_MANAGE",
				},
			},
		},
		map[string]interface{}{
			authKeysV2IDFieldName:             "key-id-2",
			authKeysV2AccountIDFieldName:      "account-id-1",
			authKeysV2ClusterIDFieldName:      "cluster-id-1",
			authKeysV2NameFieldName:           "collection-key",
			authKeysV2CreatedAtFieldName:      formatTime(createdAt),
			authKeysV2ExpiresAtFieldName:      "",
			authKeysV2CreatedByEmailFieldName: "user2@example.com",
			authKeysV2PostfixFieldName:        "post2",
			authKeysV2KeyFieldName:            "key2",
			authKeysV2CollectionAccessRulesFieldName: []interface{}{
				map[string]interface{}{
					authKeysV2CollectionNameFieldName: "collection1",
					authKeysV2AccessTypeFieldName:     "COLLECTION_ACCESS_RULE_ACCESS_TYPE_READ_WRITE",
					authKeysV2PayloadFieldName:        map[string]string(nil),
				},
				map[string]interface{}{
					authKeysV2CollectionNameFieldName: "collection2",
					authKeysV2AccessTypeFieldName:     "COLLECTION_ACCESS_RULE_ACCESS_TYPE_READ_ONLY",
					authKeysV2PayloadFieldName:        map[string]string{"filter_key": "filter_value"},
				},
			},
		},
	}

	flattened := flattenAuthKeysV2(keys, true)
	assert.Equal(t, expected, flattened)
}
