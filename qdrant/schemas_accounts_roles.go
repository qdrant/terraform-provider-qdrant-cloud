package qdrant

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qci "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/iam/v1"
)

const (
	rolesFieldTemplate = "Account Roles Schema %s field"

	rolesAccountIDFieldName   = "account_id"
	rolesRolesFieldName       = "roles"
	rolesIDFieldName          = "id"
	rolesNameFieldName        = "name"
	rolesDescriptionFieldName = "description"
	rolesRoleTypeFieldName    = "role_type"
	rolesPermissionsFieldName = "permissions"
)

// accountsRolesDataSourceSchema defines the Terraform schema for the accounts_roles data source.
func accountsRolesDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		rolesAccountIDFieldName: {
			Description: fmt.Sprintf(rolesFieldTemplate, "Account ID"),
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
		},
		rolesRolesFieldName: {
			Description: "List of account roles (system and custom).",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					rolesIDFieldName: {
						Description: fmt.Sprintf(rolesFieldTemplate, "Role ID (UUID)"),
						Type:        schema.TypeString,
						Computed:    true,
					},
					rolesNameFieldName: {
						Description: fmt.Sprintf(rolesFieldTemplate, "Role name"),
						Type:        schema.TypeString,
						Computed:    true,
					},
					rolesDescriptionFieldName: {
						Description: fmt.Sprintf(rolesFieldTemplate, "Human-readable description"),
						Type:        schema.TypeString,
						Computed:    true,
					},
					rolesRoleTypeFieldName: {
						Description: fmt.Sprintf(rolesFieldTemplate, "Role type (ROLE_TYPE_SYSTEM or ROLE_TYPE_CUSTOM)"),
						Type:        schema.TypeString,
						Computed:    true,
					},
					rolesPermissionsFieldName: {
						Description: fmt.Sprintf(rolesFieldTemplate, "List of permission values"),
						Type:        schema.TypeList,
						Computed:    true,
						Elem: &schema.Schema{
							Type: schema.TypeString,
						},
					},
				},
			},
		},
	}
}

// flattenAccountRoles converts a list of Role proto messages into a list of maps for Terraform state.
func flattenAccountRoles(roles []*qci.Role) []interface{} {
	result := make([]interface{}, 0, len(roles))
	for _, role := range roles {
		permissions := make([]interface{}, 0, len(role.GetPermissions()))
		for _, perm := range role.GetPermissions() {
			permissions = append(permissions, perm.GetValue())
		}
		result = append(result, map[string]interface{}{
			rolesIDFieldName:          role.GetId(),
			rolesNameFieldName:        role.GetName(),
			rolesDescriptionFieldName: role.GetDescription(),
			rolesRoleTypeFieldName:    role.GetRoleType().String(),
			rolesPermissionsFieldName: permissions,
		})
	}
	return result
}
