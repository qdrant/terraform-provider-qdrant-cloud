package qdrant

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// accountsAuthKeysSchema returns the schema for the auth keys.
func accountsAuthKeysSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_id": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"keys": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"created_at": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"user_id": {
						Type:     schema.TypeString,
						Computed: true,
						Optional: true,
					},
					"prefix": {
						Type:     schema.TypeString,
						Computed: true,
					},
					"cluster_id_list": {
						Type:     schema.TypeList,
						Computed: true,
						Optional: true,
						Elem:     &schema.Schema{Type: schema.TypeString},
					},
					"account_id": {
						Type:     schema.TypeString,
						Computed: true,
						Optional: true,
					},
					"token": {
						Type:     schema.TypeString,
						Computed: true,
					},
				},
			},
		},
	}
}

// APIKey represents the structure of an API key as retrieved from the Qdrant Cloud API.
type APIKey struct {
	ID            string   `json:"id"`              // The unique identifier of the API key.
	CreatedAt     string   `json:"created_at"`      // The creation timestamp of the API key.
	UserID        *string  `json:"user_id"`         // The user identifier associated with the API key, if applicable.
	Prefix        string   `json:"prefix"`          // The prefix associated with the API key.
	ClusterIDList []string `json:"cluster_id_list"` // A list of cluster identifiers associated with the API key.
	AccountID     string   `json:"account_id"`      // The account identifier associated with the API key.
}

// ApiKeyCreate defines the structure for the output response for an API key
type ApiKeyCreate struct {
	APIKey        // Embed the APIKey struct
	Token  string `json:"token"` // The token associated with the API key
}

type ApiKeyCreateRequestBody struct {
	ClusterIDList []string `json:"cluster_id_list"`
}
