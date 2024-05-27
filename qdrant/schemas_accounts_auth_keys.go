package qdrant

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	authKeysFieldTemplate = "Auth Keys Schema %s field"

	authKeysIdentifierFieldName = "id"
	authKeysCreatedAtFieldName  = "created_at"
)

// accountsAuthKeysSchema returns the schema for the auth keys.
func accountsAuthKeysSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_id": {
			Description: fmt.Sprintf(authKeysFieldTemplate, "Account Identifier"),
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"keys": {
			Description: "TODO",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Description: "TODO",
				Schema: map[string]*schema.Schema{
					"id": {
						Description: "TODO",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"created_at": {
						Description: "TODO",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"user_id": {
						Description: "TODO",
						Type:        schema.TypeString,
						Computed:    true,
						Optional:    true,
					},
					"prefix": {
						Description: "TODO",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"cluster_id_list": {
						Description: "TODO",
						Type:        schema.TypeList,
						Computed:    true,
						Optional:    true,
						Elem: &schema.Schema{
							Description: "TODO",
							Type:        schema.TypeString,
						},
					},
					"account_id": {
						Description: "TODO",
						Type:        schema.TypeString,
						Computed:    true,
						Optional:    true,
					},
					"token": {
						Description: "TODO",
						Type:        schema.TypeString,
						Computed:    true,
					},
				},
			},
		},
	}
}
