package qdrant

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	authv2 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/auth/v2"
)

const (
	authKeysV2FieldTemplate                  = "Auth Keys V2 Schema %s field"
	authKeysV2AccountIDFieldName             = "account_id"
	authKeysV2KeysFieldName                  = "keys"
	authKeysV2ClusterIDFieldName             = "cluster_id"
	authKeysV2IDFieldName                    = "id"
	authKeysV2NameFieldName                  = "name"
	authKeysV2CreatedAtFieldName             = "created_at"
	authKeysV2ExpiresAtFieldName             = "expires_at"
	authKeysV2CreatedByEmailFieldName        = "created_by_email"
	authKeysV2PostfixFieldName               = "postfix"
	authKeysV2TokenFieldName                 = "token"
	authKeysV2GlobalAccessRuleFieldName      = "global_access_rule"
	authKeysV2CollectionAccessRulesFieldName = "collection_access_rules"
	authKeysV2AccessTypeFieldName            = "access_type"
	authKeysV2CollectionNameFieldName        = "collection_name"
	authKeysV2PayloadFieldName               = "payload"
)

// accountsAuthKeysV2DataSourceSchema returns the schema for the auth keys data source.
func accountsAuthKeysV2DataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		authKeysV2AccountIDFieldName: {
			Description: fmt.Sprintf(authKeysV2FieldTemplate, "Account Identifier where all those Auth Keys belongs to"),
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		authKeysV2ClusterIDFieldName: {
			Description: fmt.Sprintf(authKeysV2FieldTemplate, "Cluster Identifier for which this Auth Key is attached"),
			Type:        schema.TypeString,
			Required:    true,
		},
		authKeysV2KeysFieldName: {
			Description: fmt.Sprintf(authKeysV2FieldTemplate, "List of Auth Keys"),
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Description: fmt.Sprintf(authKeysV2FieldTemplate, "Individual Auth Key"),
				Schema:      accountsAuthKeyV2ResourceSchema(true),
			},
		},
	}
}

// accountsAuthKeyV2ResourceSchema returns the schema for a single auth key resource.
func accountsAuthKeyV2ResourceSchema(isDataSource bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		authKeysV2IDFieldName: {
			Description: fmt.Sprintf(authKeysV2FieldTemplate, "Auth Key Identifier"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		authKeysV2AccountIDFieldName: {
			Description: fmt.Sprintf(authKeysV2FieldTemplate, "Account Identifier"),
			Type:        schema.TypeString,
			Optional:    !isDataSource,
			Computed:    true,
		},
		authKeysV2ClusterIDFieldName: {
			Description: fmt.Sprintf(authKeysV2FieldTemplate, "Cluster Identifier for which this Auth Key is attached"),
			Type:        schema.TypeString,
			Required:    !isDataSource,
			Computed:    isDataSource,
			ForceNew:    !isDataSource,
		},
		authKeysV2NameFieldName: {
			Description: fmt.Sprintf(authKeysV2FieldTemplate, "Auth Key Name"),
			Type:        schema.TypeString,
			Required:    !isDataSource,
			Computed:    isDataSource,
			ForceNew:    !isDataSource,
		},
		authKeysV2CreatedAtFieldName: {
			Description: fmt.Sprintf(authKeysV2FieldTemplate, "Timestamp when the Auth Key is created"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		authKeysV2ExpiresAtFieldName: {
			Description: fmt.Sprintf(authKeysV2FieldTemplate, "Timestamp when the Auth Key expires"),
			Type:        schema.TypeString,
			Optional:    !isDataSource,
			Computed:    true,
			ForceNew:    !isDataSource,
		},
		authKeysV2CreatedByEmailFieldName: {
			Description: fmt.Sprintf(authKeysV2FieldTemplate, "Email of the user who created the key"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		authKeysV2PostfixFieldName: {
			Description: fmt.Sprintf(authKeysV2FieldTemplate, "Postfix of the Auth Key"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		authKeysV2TokenFieldName: {
			Description: fmt.Sprintf(authKeysV2FieldTemplate, "Secret token for this Auth Key"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		authKeysV2GlobalAccessRuleFieldName: {
			Description:   "A rule granting global access to the entire database. Cannot be used together with `collection_access_rules`.",
			Type:          schema.TypeList,
			Optional:      !isDataSource,
			Computed:      isDataSource,
			ForceNew:      !isDataSource,
			MaxItems:      1,
			ConflictsWith: []string{authKeysV2CollectionAccessRulesFieldName},
			Elem: &schema.Resource{
				Schema: globalAccessRuleSchema(isDataSource),
			},
		},
		authKeysV2CollectionAccessRulesFieldName: {
			Description:   "A list of rules granting access to specific collections. Cannot be used together with `global_access_rule`.",
			Type:          schema.TypeList,
			Optional:      !isDataSource,
			Computed:      isDataSource,
			ForceNew:      !isDataSource,
			ConflictsWith: []string{authKeysV2GlobalAccessRuleFieldName},
			Elem: &schema.Resource{
				Schema: collectionAccessRuleSchema(isDataSource),
			},
		},
	}
}

// globalAccessRuleSchema defines the schema for a global access rule.
func globalAccessRuleSchema(isDataSource bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		authKeysV2AccessTypeFieldName: {
			Description: "Access type for global access. Can be `GLOBAL_ACCESS_RULE_ACCESS_TYPE_READ_ONLY` or `GLOBAL_ACCESS_RULE_ACCESS_TYPE_MANAGE`.",
			Type:        schema.TypeString,
			Required:    !isDataSource,
			Computed:    isDataSource,
		},
	}
}

// collectionAccessRuleSchema defines the schema for a collection-specific access rule.
func collectionAccessRuleSchema(isDataSource bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		authKeysV2CollectionNameFieldName: {
			Description: "Name of the collection.",
			Type:        schema.TypeString,
			Required:    !isDataSource,
			Computed:    isDataSource,
		},
		authKeysV2AccessTypeFieldName: {
			Description: "Access type for the collection. Can be `COLLECTION_ACCESS_RULE_ACCESS_TYPE_READ_ONLY` or `COLLECTION_ACCESS_RULE_ACCESS_TYPE_READ_WRITE`.",
			Type:        schema.TypeString,
			Required:    !isDataSource,
			Computed:    isDataSource,
		},
		authKeysV2PayloadFieldName: {
			Description: "Payload restrictions.",
			Type:        schema.TypeMap,
			Optional:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
	}
}

// flattenAuthKeysV2 flattens the API keys response into a slice of map[string]interface{}.
func flattenAuthKeysV2(keys []*authv2.DatabaseApiKey) []interface{} {
	var flattenedKeys []interface{}
	for _, key := range keys {
		flattenedKeys = append(flattenedKeys, flattenAuthKeyV2(key))
	}
	return flattenedKeys
}

// flattenAuthKeyV2 flattens the API key response into a map[string]interface{}.
func flattenAuthKeyV2(key *authv2.DatabaseApiKey) map[string]interface{} {
	data := map[string]interface{}{
		authKeysV2IDFieldName:             key.GetId(),
		authKeysV2AccountIDFieldName:      key.GetAccountId(),
		authKeysV2ClusterIDFieldName:      key.GetClusterId(),
		authKeysV2NameFieldName:           key.GetName(),
		authKeysV2CreatedAtFieldName:      formatTime(key.GetCreatedAt()),
		authKeysV2ExpiresAtFieldName:      formatTime(key.GetExpiresAt()),
		authKeysV2CreatedByEmailFieldName: key.GetCreatedByEmail(),
		authKeysV2PostfixFieldName:        key.GetPostfix(),
		authKeysV2TokenFieldName:          key.GetKey(),
	}

	globalRules, collectionRules := separateAccessRules(key.GetAccessRules())
	if len(globalRules) > 0 {
		data[authKeysV2GlobalAccessRuleFieldName] = flattenGlobalAccessRules(globalRules)
	}
	if len(collectionRules) > 0 {
		data[authKeysV2CollectionAccessRulesFieldName] = flattenCollectionAccessRules(collectionRules)
	}

	return data
}

// separateAccessRules splits a list of generic access rules into lists of global and collection-specific rules.
func separateAccessRules(rules []*authv2.AccessRule) ([]*authv2.GlobalAccessRule, []*authv2.CollectionAccessRule) {
	var globalRules []*authv2.GlobalAccessRule
	var collectionRules []*authv2.CollectionAccessRule
	for _, rule := range rules {
		if r := rule.GetGlobalAccess(); r != nil {
			globalRules = append(globalRules, r)
		}
		if r := rule.GetCollectionAccess(); r != nil {
			collectionRules = append(collectionRules, r)
		}
	}
	return globalRules, collectionRules
}

// flattenGlobalAccessRules flattens a list of global access rules into a format suitable for Terraform state.
func flattenGlobalAccessRules(rules []*authv2.GlobalAccessRule) []interface{} {
	if len(rules) == 0 {
		return nil
	}
	return []interface{}{
		map[string]interface{}{
			authKeysV2AccessTypeFieldName: rules[0].GetAccessType().String(),
		},
	}
}

// flattenCollectionAccessRules flattens a list of collection access rules into a format suitable for Terraform state.
func flattenCollectionAccessRules(rules []*authv2.CollectionAccessRule) []interface{} {
	flattened := make([]interface{}, len(rules))
	for i, rule := range rules {
		flattened[i] = map[string]interface{}{
			authKeysV2CollectionNameFieldName: rule.GetCollectionName(),
			authKeysV2AccessTypeFieldName:     rule.GetAccessType().String(),
			authKeysV2PayloadFieldName:        rule.GetPayload(),
		}
	}
	return flattened
}
