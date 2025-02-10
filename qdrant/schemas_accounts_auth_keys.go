package qdrant

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qcAuth "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/auth/v2"
)

const (
	authKeysFieldTemplate      = "Auth Keys Schema %s field"
	authKeysAccountIDFieldName = "account_id"
	authKeysKeysFieldName      = "keys"

	authKeysKeysFieldTemplate       = "Auth Keys Keys Schema %s field"
	authKeysKeysIDFieldName         = "id"
	authKeysKeysCreatedAtFieldName  = "created_at"
	authKeysKeysPrefixFieldName     = "prefix"
	authKeysKeysClusterIDsFieldName = "cluster_ids"
	authKeysKeysTokenFieldName      = "token"
)

// accountsAuthKeysSchema returns the schema for the auth keys.
func accountsAuthKeysSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		authKeysAccountIDFieldName: {
			Description: fmt.Sprintf(authKeysFieldTemplate, "Account Identifier where all those Auth Keys belongs to"),
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		authKeysKeysFieldName: {
			Description: fmt.Sprintf(authKeysFieldTemplate, "List of Auth Keys"),
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Description: fmt.Sprintf(authKeysFieldTemplate, "Individual Auth Keys"),
				Schema:      accountsAuthKeySchema(),
			},
		},
	}
}

// accountsAuthKeySchema returns the schema for the auth key.
func accountsAuthKeySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		authKeysKeysIDFieldName: {
			Description: fmt.Sprintf(authKeysKeysFieldTemplate, "Auth Key Identifier"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		authKeysKeysCreatedAtFieldName: {
			Description: fmt.Sprintf(authKeysKeysFieldTemplate, "Timestamp when the Auth Key is created"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		authKeysKeysPrefixFieldName: {
			Description: fmt.Sprintf(authKeysKeysFieldTemplate, "Prefix of the Auth Key (the first few bytes from the token)"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		authKeysKeysClusterIDsFieldName: {
			Description: fmt.Sprintf(authKeysKeysFieldTemplate, "Cluster Identifiers for which this Auth Key is attached"),
			Type:        schema.TypeList,
			Required:    true,
			ForceNew:    true,
			Elem: &schema.Schema{
				Description: fmt.Sprintf(authKeysKeysFieldTemplate, "Single Cluster Identifier for which this Auth Key is attached"),
				Type:        schema.TypeString,
			},
		},
		authKeysKeysTokenFieldName: {
			Description: fmt.Sprintf(authKeysKeysFieldTemplate, "Secret token for this Auth Key (handle with care!)"),
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

// flattenAuthKeys flattens the API keys response into a slice of map[string]interface{}.
func flattenAuthKeys(keys []*qcAuth.ApiKey) []interface{} {
	var flattenedKeys []interface{}
	for _, key := range keys {
		flattenedKeys = append(flattenedKeys, flattenAuthKey(key))
	}
	return flattenedKeys
}

// flattenAuthKey flattens the API key response into a slice of map[string]interface{}.
func flattenAuthKey(key *qcAuth.ApiKey) map[string]interface{} {
	result := map[string]interface{}{
		authKeysKeysIDFieldName:         key.GetId(),
		authKeysKeysCreatedAtFieldName:  formatTime(key.GetCreatedAt()),
		authKeysKeysClusterIDsFieldName: key.GetClusterIds(),
		authKeysKeysPrefixFieldName:     key.GetPrefix(),
		authKeysKeysTokenFieldName:      key.GetToken(),
	}
	return result
}
