package qdrant

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qci "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/iam/v1"
)

const (
	roleFieldTemplate = "Role Schema %s field"

	// Writable fields.
	roleAccountIdFieldName   = "account_id"
	roleNameFieldName        = "name"
	roleDescriptionFieldName = "description"
	roleRoleTypeFieldName    = "role_type"
	rolePermissionsFieldName = "permissions"

	// Read-only fields.
	roleIdFieldName             = "id"
	roleCreatedAtFieldName      = "created_at"
	roleLastModifiedAtFieldName = "last_modified_at"
	roleSubTypeFieldName        = "sub_type"
)

// permissionsNestedSchema returns the nested Terraform schema for a single permission entry.
func permissionsNestedSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"value": {
			Description: "Permission value (e.g., \"read:backups\").",
			Type:        schema.TypeString,
			Required:    true,
		},
		"category": {
			Description: "Permission category (e.g., \"Cluster\", \"Account\").",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

// accountsRoleSchema defines the Terraform schema for an Role resource.
func accountsRoleSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		// Writable
		roleAccountIdFieldName: {
			Description: fmt.Sprintf(roleFieldTemplate, "Account ID"),
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
		},
		roleNameFieldName: {
			Description: fmt.Sprintf(roleFieldTemplate, "Role name (printable, length 4-64)"),
			Type:        schema.TypeString,
			Required:    true,
		},
		roleDescriptionFieldName: {
			Description: fmt.Sprintf(roleFieldTemplate, "Human-readable description (<=256 chars)"),
			Type:        schema.TypeString,
			Optional:    true,
		},
		roleRoleTypeFieldName: {
			Description: fmt.Sprintf(roleFieldTemplate, "Role type"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		rolePermissionsFieldName: {
			Description: "Permissions assigned to this role (unordered).",
			Type:        schema.TypeSet,
			Required:    true,
			MinItems:    1,
			Elem: &schema.Resource{
				Schema: permissionsNestedSchema(),
			},
		},

		// Read-only
		roleIdFieldName: {
			Description: fmt.Sprintf(roleFieldTemplate, "Role ID"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		roleCreatedAtFieldName: {
			Description: fmt.Sprintf(roleFieldTemplate, "Creation timestamp"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		roleLastModifiedAtFieldName: {
			Description: fmt.Sprintf(roleFieldTemplate, "Last modification timestamp"),
			Type:        schema.TypeString,
			Computed:    true,
		},
		roleSubTypeFieldName: {
			Description: fmt.Sprintf(roleFieldTemplate, "System role sub-type (if any)"),
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

// expandRole builds a qci.Role from Terraform resource data.
// The includeID flag controls whether to populate the ID field (required for updates).
func expandRole(d *schema.ResourceData, accountID string, includeID bool) (*qci.Role, error) {
	if accountID == "" {
		return nil, fmt.Errorf("account ID not specified")
	}

	role := &qci.Role{
		AccountId:   accountID,
		Name:        d.Get(roleNameFieldName).(string),
		Description: d.Get(roleDescriptionFieldName).(string),
		RoleType:    qci.RoleType_ROLE_TYPE_CUSTOM,
		Permissions: expandPermissions(d.Get(rolePermissionsFieldName)),
	}

	if includeID {
		role.Id = d.Id()
	}

	return role, nil
}

// expandPermissions converts the Terraform permissions attribute into []*qci.Permission.
func expandPermissions(v interface{}) []*qci.Permission {
	var items []interface{}

	switch t := v.(type) {
	case *schema.Set:
		items = t.List()
	case []interface{}:
		items = t
	default:
		return nil
	}

	if len(items) == 0 {
		return nil
	}

	out := make([]*qci.Permission, 0, len(items))
	for _, it := range items {
		if it == nil {
			continue
		}
		m := it.(map[string]interface{})
		out = append(out, &qci.Permission{
			Value:    m["value"].(string),
			Category: m["category"].(string),
		})
	}
	return out
}

// flattenRole maps a qci.Role into Terraform state fields.
func flattenRole(r *qci.Role) map[string]interface{} {
	if r == nil {
		return map[string]interface{}{}
	}

	out := map[string]interface{}{
		roleIdFieldName:          r.GetId(),
		roleAccountIdFieldName:   r.GetAccountId(),
		roleNameFieldName:        r.GetName(),
		roleDescriptionFieldName: r.GetDescription(),
		roleRoleTypeFieldName:    r.GetRoleType().String(),
		rolePermissionsFieldName: flattenPermissions(r.GetPermissions()),
	}

	if ts := r.GetCreatedAt(); ts != nil {
		out[roleCreatedAtFieldName] = formatTime(ts)
	}
	if ts := r.GetLastModifiedAt(); ts != nil {
		out[roleLastModifiedAtFieldName] = formatTime(ts)
	}
	if r.SubType != nil {
		out[roleSubTypeFieldName] = r.GetSubType().String()
	}

	return out
}

// flattenPermissions converts []*qci.Permission into the Terraform attribute representation.
func flattenPermissions(ps []*qci.Permission) []interface{} {
	if len(ps) == 0 {
		return []interface{}{}
	}
	out := make([]interface{}, 0, len(ps))
	for _, p := range ps {
		out = append(out, map[string]interface{}{
			"value":    p.GetValue(),
			"category": p.GetCategory(),
		})
	}
	return out
}
