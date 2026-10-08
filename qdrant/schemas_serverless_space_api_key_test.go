package qdrant

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	spaceauthv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/serverless/space/auth/v1"
)

func TestServerlessAccessRuleSchemas(t *testing.T) {
	t.Run("global "+slSpaceAPIKeyAccessTypeFieldName, func(t *testing.T) {
		testGeneratedEnumValidation(t, slSpaceAPIKeyAccessTypeFieldName,
			serverlessGlobalAccessRuleSchema(false)[slSpaceAPIKeyAccessTypeFieldName],
			serverlessGlobalAccessRuleSchema(true)[slSpaceAPIKeyAccessTypeFieldName],
			spaceauthv1.GlobalAccessRuleAccessType_name)
	})
	t.Run("collection "+slSpaceAPIKeyAccessTypeFieldName, func(t *testing.T) {
		testGeneratedEnumValidation(t, slSpaceAPIKeyAccessTypeFieldName,
			serverlessCollectionAccessRuleSchema(false)[slSpaceAPIKeyAccessTypeFieldName],
			serverlessCollectionAccessRuleSchema(true)[slSpaceAPIKeyAccessTypeFieldName],
			spaceauthv1.CollectionAccessRuleAccessType_name)
	})
}

func TestServerlessSpaceAPIKeySchema(t *testing.T) {
	s := serverlessSpaceAPIKeySchema(false)
	t.Run("when used as resource, the secret key is sensitive", func(t *testing.T) {
		require.Contains(t, s, slSpaceAPIKeyKeyFieldName)
		assert.True(t, s[slSpaceAPIKeyKeyFieldName].Sensitive)
	})
	t.Run("when used as data source, the secret key is not exposed", func(t *testing.T) {
		assert.NotContains(t, serverlessSpaceAPIKeySchema(true), slSpaceAPIKeyKeyFieldName)
	})
	t.Run("when used as resource, all configurable fields force a new key", func(t *testing.T) {
		for k, f := range s {
			if f.Optional || f.Required {
				assert.True(t, f.ForceNew, k)
			}
		}
	})
}

func TestFlattenServerlessSpaceAPIKey(t *testing.T) {
	key := &spaceauthv1.SpaceApiKey{
		Id:        "key-id",
		AccountId: "account-id",
		SpaceId:   "space-id",
		Name:      "my key",
		Postfix:   "abcd",
		Key:       "secret",
		State:     &spaceauthv1.SpaceApiKeyState{Phase: spaceauthv1.SpaceApiKeyStatePhase_SPACE_API_KEY_STATE_PHASE_READY},
		AccessRules: []*spaceauthv1.AccessRule{
			{Scope: &spaceauthv1.AccessRule_CollectionAccess{CollectionAccess: &spaceauthv1.CollectionAccessRule{
				CollectionName: "c1",
				AccessType:     spaceauthv1.CollectionAccessRuleAccessType_COLLECTION_ACCESS_RULE_ACCESS_TYPE_READ_ONLY,
			}}},
		},
	}

	flat := flattenServerlessSpaceAPIKey(key)

	assert.Equal(t, "key-id", flat[slSpaceAPIKeyIDFieldName])
	assert.Equal(t, "space-id", flat[slSpaceAPIKeySpaceIDFieldName])
	assert.Equal(t, "my key", flat[slSpaceAPIKeyNameFieldName])
	assert.Equal(t, "abcd", flat[slSpaceAPIKeyPostfixFieldName])
	assert.Empty(t, flat[slSpaceAPIKeyExpiresAtFieldName])
	assert.NotContains(t, flat, slSpaceAPIKeyKeyFieldName, "the secret must only be set on create")
	assert.Equal(t, []interface{}{}, flat[slSpaceAPIKeyGlobalAccessRuleFieldName])
	assert.Equal(t, []interface{}{map[string]interface{}{
		slSpaceAPIKeyCollectionNameFieldName: "c1",
		slSpaceAPIKeyAccessTypeFieldName:     "COLLECTION_ACCESS_RULE_ACCESS_TYPE_READ_ONLY",
	}}, flat[slSpaceAPIKeyCollectionAccessRulesFieldName])

	d := schema.TestResourceDataRaw(t, serverlessSpaceAPIKeySchema(false), map[string]interface{}{})
	require.NoError(t, setFlattened(d, flat))
}

func TestExpandServerlessSpaceAPIKey(t *testing.T) {
	t.Run("when a global rule is configured, it builds a global access rule", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, serverlessSpaceAPIKeySchema(false), map[string]interface{}{
			slSpaceAPIKeySpaceIDFieldName:   "space-id",
			slSpaceAPIKeyNameFieldName:      "my key",
			slSpaceAPIKeyExpiresAtFieldName: "2030-01-01T00:00:00Z",
			slSpaceAPIKeyGlobalAccessRuleFieldName: []interface{}{map[string]interface{}{
				slSpaceAPIKeyAccessTypeFieldName: "GLOBAL_ACCESS_RULE_ACCESS_TYPE_READ_ONLY",
			}},
		})
		key := expandServerlessSpaceAPIKey(d, "account-id")
		assert.Equal(t, "account-id", key.GetAccountId())
		assert.Equal(t, "space-id", key.GetSpaceId())
		assert.Equal(t, "my key", key.GetName())
		assert.Equal(t, "2030-01-01T00:00:00Z", formatTime(key.GetExpiresAt()))
		require.Len(t, key.GetAccessRules(), 1)
		assert.Equal(t, spaceauthv1.GlobalAccessRuleAccessType_GLOBAL_ACCESS_RULE_ACCESS_TYPE_READ_ONLY,
			key.GetAccessRules()[0].GetGlobalAccess().GetAccessType())
	})
	t.Run("when collection rules are configured, it builds collection access rules", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, serverlessSpaceAPIKeySchema(false), map[string]interface{}{
			slSpaceAPIKeySpaceIDFieldName: "space-id",
			slSpaceAPIKeyNameFieldName:    "my key",
			slSpaceAPIKeyCollectionAccessRulesFieldName: []interface{}{
				map[string]interface{}{
					slSpaceAPIKeyCollectionNameFieldName: "c1",
					slSpaceAPIKeyAccessTypeFieldName:     "COLLECTION_ACCESS_RULE_ACCESS_TYPE_READ_WRITE",
				},
				map[string]interface{}{
					slSpaceAPIKeyCollectionNameFieldName: "c2",
					slSpaceAPIKeyAccessTypeFieldName:     "COLLECTION_ACCESS_RULE_ACCESS_TYPE_READ_ONLY",
				},
			},
		})
		key := expandServerlessSpaceAPIKey(d, "account-id")
		assert.Nil(t, key.ExpiresAt)
		require.Len(t, key.GetAccessRules(), 2)
		assert.Equal(t, "c1", key.GetAccessRules()[0].GetCollectionAccess().GetCollectionName())
		assert.Equal(t, spaceauthv1.CollectionAccessRuleAccessType_COLLECTION_ACCESS_RULE_ACCESS_TYPE_READ_WRITE,
			key.GetAccessRules()[0].GetCollectionAccess().GetAccessType())
		assert.Equal(t, "c2", key.GetAccessRules()[1].GetCollectionAccess().GetCollectionName())
	})
	t.Run("when no rules are configured, it sends no rules", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, serverlessSpaceAPIKeySchema(false), map[string]interface{}{
			slSpaceAPIKeySpaceIDFieldName: "space-id",
			slSpaceAPIKeyNameFieldName:    "my key",
		})
		assert.Empty(t, expandServerlessSpaceAPIKey(d, "account-id").GetAccessRules())
	})
}
