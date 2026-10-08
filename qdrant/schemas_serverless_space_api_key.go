package qdrant

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	spaceauthv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/serverless/space/auth/v1"
)

const (
	slSpaceAPIKeyFieldTemplate = "Serverless Space API Key Schema %s field"

	slSpaceAPIKeyIDFieldName                    = "id"
	slSpaceAPIKeyAccountIDFieldName             = "account_id"
	slSpaceAPIKeySpaceIDFieldName               = "space_id"
	slSpaceAPIKeyNameFieldName                  = "name"
	slSpaceAPIKeyExpiresAtFieldName             = "expires_at"
	slSpaceAPIKeyCreatedAtFieldName             = "created_at"
	slSpaceAPIKeyPostfixFieldName               = "postfix"
	slSpaceAPIKeyKeyFieldName                   = "key"
	slSpaceAPIKeyStateFieldName                 = "state"
	slSpaceAPIKeyGlobalAccessRuleFieldName      = "global_access_rule"
	slSpaceAPIKeyCollectionAccessRulesFieldName = "collection_access_rules"
	slSpaceAPIKeyAccessTypeFieldName            = "access_type"
	slSpaceAPIKeyCollectionNameFieldName        = "collection_name"
	slSpaceAPIKeysFieldName                     = "keys"
)

var (
	slSpaceAPIKeyNameRegexp = regexp.MustCompile(`^[\w\s-]+$`)
)

// serverlessSpaceAPIKeySchema returns the schema of a serverless space API key.
// All configurable fields are ForceNew, as API keys cannot be updated.
func serverlessSpaceAPIKeySchema(asDataSource bool) map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		slSpaceAPIKeyIDFieldName: {
			Description: fmt.Sprintf(slSpaceAPIKeyFieldTemplate, "API Key ID"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		slSpaceAPIKeyAccountIDFieldName: {
			Description: fmt.Sprintf(slSpaceAPIKeyFieldTemplate, "Account ID (defaults to the provider account ID)"),
			Type:        schema.TypeString,
			Optional:    !asDataSource,
			Computed:    true,
			ForceNew:    !asDataSource,
		},
		slSpaceAPIKeySpaceIDFieldName: {
			Description: fmt.Sprintf(slSpaceAPIKeyFieldTemplate, "ID of the serverless space this API key grants access to"),
			Type:        schema.TypeString,
			Required:    !asDataSource,
			Computed:    asDataSource,
			ForceNew:    !asDataSource,
		},
		slSpaceAPIKeyNameFieldName: {
			Description: fmt.Sprintf(slSpaceAPIKeyFieldTemplate, "Name of the API key (4-128 characters: letters, numbers, spaces, `_` and `-`)"),
			Type:        schema.TypeString,
			Required:    !asDataSource,
			Computed:    asDataSource,
			ForceNew:    !asDataSource,
		},
		slSpaceAPIKeyExpiresAtFieldName: {
			Description: fmt.Sprintf(slSpaceAPIKeyFieldTemplate, "Expiration time (RFC3339) of the API key. If not set, the key does not expire"),
			Type:        schema.TypeString,
			Optional:    !asDataSource,
			Computed:    true,
			ForceNew:    !asDataSource,
		},
		slSpaceAPIKeyCreatedAtFieldName: {
			Description: fmt.Sprintf(slSpaceAPIKeyFieldTemplate, "Creation time"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		slSpaceAPIKeyPostfixFieldName: {
			Description: fmt.Sprintf(slSpaceAPIKeyFieldTemplate, "Last characters of the key, to identify it"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		slSpaceAPIKeyStateFieldName: {
			Description: fmt.Sprintf(slSpaceAPIKeyFieldTemplate, "Current state of the API key"),
			Type:        schema.TypeList,
			Computed:    true,
			Elem:        &schema.Resource{Schema: serverlessStateSchema()},
		},
		slSpaceAPIKeyGlobalAccessRuleFieldName: {
			Description: "A rule granting access to the entire space. Cannot be used with `collection_access_rules`. If no rules are provided, global `MANAGE` access is granted.",
			Type:        schema.TypeList,
			Optional:    !asDataSource,
			Computed:    true,
			ForceNew:    !asDataSource,
			MaxItems:    configurableMaxItems(1, asDataSource),
			Elem:        &schema.Resource{Schema: serverlessGlobalAccessRuleSchema(asDataSource)},
		},
		slSpaceAPIKeyCollectionAccessRulesFieldName: {
			Description: "A list of rules granting access to specific collections (max 20). Cannot be used with `global_access_rule`.",
			Type:        schema.TypeList,
			Optional:    !asDataSource,
			Computed:    true,
			ForceNew:    !asDataSource,
			MaxItems:    configurableMaxItems(20, asDataSource),
			Elem:        &schema.Resource{Schema: serverlessCollectionAccessRuleSchema(asDataSource)},
		},
	}
	if !asDataSource {
		s[slSpaceAPIKeyKeyFieldName] = &schema.Schema{
			Description: fmt.Sprintf(slSpaceAPIKeyFieldTemplate, "Secret key, only available after creation (empty after import)"),
			Type:        schema.TypeString,
			Computed:    true,
			Sensitive:   true,
		}
		s[slSpaceAPIKeyNameFieldName].ValidateDiagFunc = validation.ToDiagFunc(validation.All(
			validation.StringLenBetween(4, 128),
			validation.StringMatch(slSpaceAPIKeyNameRegexp, "must only contain letters, numbers, spaces, '_' and '-'"),
		))
		s[slSpaceAPIKeyExpiresAtFieldName].ValidateDiagFunc = validation.ToDiagFunc(validation.IsRFC3339Time)
		s[slSpaceAPIKeyGlobalAccessRuleFieldName].ConflictsWith = []string{slSpaceAPIKeyCollectionAccessRulesFieldName}
		s[slSpaceAPIKeyCollectionAccessRulesFieldName].ConflictsWith = []string{slSpaceAPIKeyGlobalAccessRuleFieldName}
	}
	return s
}

// serverlessGlobalAccessRuleSchema defines the schema for a serverless global access rule.
func serverlessGlobalAccessRuleSchema(asDataSource bool) map[string]*schema.Schema {
	validAccessTypes := protoEnumNames(spaceauthv1.GlobalAccessRuleAccessType_name)
	accessType := &schema.Schema{
		Description: fmt.Sprintf("Access type for the entire space. Must be one of: %s.", strings.Join(validAccessTypes, ", ")),
		Type:        schema.TypeString,
		Required:    !asDataSource,
		Computed:    asDataSource,
	}
	if !asDataSource {
		accessType.ValidateDiagFunc = validation.ToDiagFunc(validation.StringInSlice(validAccessTypes, false))
	}
	return map[string]*schema.Schema{
		slSpaceAPIKeyAccessTypeFieldName: accessType,
	}
}

// serverlessCollectionAccessRuleSchema defines the schema for a serverless collection access rule.
func serverlessCollectionAccessRuleSchema(asDataSource bool) map[string]*schema.Schema {
	validAccessTypes := protoEnumNames(spaceauthv1.CollectionAccessRuleAccessType_name)
	accessType := &schema.Schema{
		Description: fmt.Sprintf("Access type for the collection. Must be one of: %s.", strings.Join(validAccessTypes, ", ")),
		Type:        schema.TypeString,
		Required:    !asDataSource,
		Computed:    asDataSource,
	}
	collectionName := &schema.Schema{
		Description: "Name of the collection.",
		Type:        schema.TypeString,
		Required:    !asDataSource,
		Computed:    asDataSource,
	}
	if !asDataSource {
		accessType.ValidateDiagFunc = validation.ToDiagFunc(validation.StringInSlice(validAccessTypes, false))
		collectionName.ValidateDiagFunc = serverlessCollectionNameValidation()
	}
	return map[string]*schema.Schema{
		slSpaceAPIKeyCollectionNameFieldName: collectionName,
		slSpaceAPIKeyAccessTypeFieldName:     accessType,
	}
}

func serverlessSpaceAPIKeysDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		slSpaceAPIKeyAccountIDFieldName: {
			Description: fmt.Sprintf(slSpaceAPIKeyFieldTemplate, "Account ID (defaults to the provider account ID)"),
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
		},
		slSpaceAPIKeySpaceIDFieldName: {
			Description: fmt.Sprintf(slSpaceAPIKeyFieldTemplate, "ID of the serverless space to list the API keys for"),
			Type:        schema.TypeString,
			Required:    true,
		},
		slSpaceAPIKeysFieldName: {
			Description: "List of API keys of the serverless space",
			Type:        schema.TypeList,
			Computed:    true,
			Elem:        &schema.Resource{Schema: serverlessSpaceAPIKeySchema(true)},
		},
	}
}

// flattenServerlessSpaceAPIKey flattens the API key into a map, the secret key is never included.
func flattenServerlessSpaceAPIKey(key *spaceauthv1.SpaceApiKey) map[string]interface{} {
	globalRules, collectionRules := []interface{}{}, []interface{}{}
	for _, rule := range key.GetAccessRules() {
		if r := rule.GetGlobalAccess(); r != nil {
			globalRules = append(globalRules, map[string]interface{}{
				slSpaceAPIKeyAccessTypeFieldName: flattenProtoEnum(r.GetAccessType()),
			})
		}
		if r := rule.GetCollectionAccess(); r != nil {
			collectionRules = append(collectionRules, map[string]interface{}{
				slSpaceAPIKeyCollectionNameFieldName: r.GetCollectionName(),
				slSpaceAPIKeyAccessTypeFieldName:     flattenProtoEnum(r.GetAccessType()),
			})
		}
	}
	return map[string]interface{}{
		slSpaceAPIKeyIDFieldName:                    key.GetId(),
		slSpaceAPIKeyAccountIDFieldName:             key.GetAccountId(),
		slSpaceAPIKeySpaceIDFieldName:               key.GetSpaceId(),
		slSpaceAPIKeyNameFieldName:                  key.GetName(),
		slSpaceAPIKeyExpiresAtFieldName:             formatTime(key.GetExpiresAt()),
		slSpaceAPIKeyCreatedAtFieldName:             formatTime(key.GetCreatedAt()),
		slSpaceAPIKeyPostfixFieldName:               key.GetPostfix(),
		slSpaceAPIKeyStateFieldName:                 flattenServerlessState(flattenProtoEnum(key.GetState().GetPhase()), key.GetState().GetReason()),
		slSpaceAPIKeyGlobalAccessRuleFieldName:      globalRules,
		slSpaceAPIKeyCollectionAccessRulesFieldName: collectionRules,
	}
}

func flattenServerlessSpaceAPIKeys(keys []*spaceauthv1.SpaceApiKey) []interface{} {
	result := make([]interface{}, 0, len(keys))
	for _, key := range keys {
		result = append(result, flattenServerlessSpaceAPIKey(key))
	}
	return result
}

// expandServerlessSpaceAPIKey builds the API key to create from the resource data.
func expandServerlessSpaceAPIKey(d *schema.ResourceData, accountID string) *spaceauthv1.SpaceApiKey {
	key := &spaceauthv1.SpaceApiKey{
		AccountId: accountID,
		SpaceId:   d.Get(slSpaceAPIKeySpaceIDFieldName).(string),
		Name:      d.Get(slSpaceAPIKeyNameFieldName).(string),
	}
	if v, ok := d.GetOk(slSpaceAPIKeyExpiresAtFieldName); ok {
		key.ExpiresAt = parseTime(v.(string))
	}
	key.AccessRules = expandServerlessAccessRules(
		d.Get(slSpaceAPIKeyGlobalAccessRuleFieldName).([]interface{}),
		d.Get(slSpaceAPIKeyCollectionAccessRulesFieldName).([]interface{}),
	)
	return key
}

func expandServerlessAccessRules(global, collections []interface{}) []*spaceauthv1.AccessRule {
	var rules []*spaceauthv1.AccessRule
	for _, v := range global {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		rules = append(rules, &spaceauthv1.AccessRule{Scope: &spaceauthv1.AccessRule_GlobalAccess{
			GlobalAccess: &spaceauthv1.GlobalAccessRule{
				AccessType: spaceauthv1.GlobalAccessRuleAccessType(
					spaceauthv1.GlobalAccessRuleAccessType_value[item[slSpaceAPIKeyAccessTypeFieldName].(string)]),
			},
		}})
	}
	for _, v := range collections {
		item, ok := v.(map[string]interface{})
		if !ok {
			continue
		}
		rules = append(rules, &spaceauthv1.AccessRule{Scope: &spaceauthv1.AccessRule_CollectionAccess{
			CollectionAccess: &spaceauthv1.CollectionAccessRule{
				CollectionName: item[slSpaceAPIKeyCollectionNameFieldName].(string),
				AccessType: spaceauthv1.CollectionAccessRuleAccessType(
					spaceauthv1.CollectionAccessRuleAccessType_value[item[slSpaceAPIKeyAccessTypeFieldName].(string)]),
			},
		}})
	}
	return rules
}
