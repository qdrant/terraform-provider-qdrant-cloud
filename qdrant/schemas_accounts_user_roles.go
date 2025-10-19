package qdrant

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	userRolesFieldTemplate = "User Roles Schema %s field"

	// Writable fields.
	userRolesAccountIdFieldName = "account_id"
	userRolesUserEmailFieldName = "user_email"
	userRolesRoleIdsFieldName   = "role_ids"
	userRolesKeepOnDestroyName  = "keep_on_destroy"

	// Computed fields.
	userRolesUserIdFieldName = "user_id"
)

// accountsUserRolesSchema defines the Terraform schema for a "one resource per user" role assignment resource.
func accountsUserRolesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Writable
		userRolesAccountIdFieldName: {
			Description: fmt.Sprintf(userRolesFieldTemplate, "Account ID"),
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
		},
		userRolesUserEmailFieldName: {
			Description: fmt.Sprintf(userRolesFieldTemplate, "User email (will be resolved to user_id via AccountService.ListAccountMembers)"),
			Type:        schema.TypeString,
			Required:    true,
		},
		userRolesRoleIdsFieldName: {
			Description: "Authoritative set of role IDs assigned to this user (unordered).",
			Type:        schema.TypeSet,
			Required:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		userRolesKeepOnDestroyName: {
			Description: "If true, the provider will not revoke roles on destroy.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},

		// Computed
		userRolesUserIdFieldName: {
			Description: fmt.Sprintf(userRolesFieldTemplate, "Resolved User ID"),
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}
