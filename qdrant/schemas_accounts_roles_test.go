package qdrant

import (
	"testing"

	"github.com/stretchr/testify/assert"

	qci "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/iam/v1"
)

func TestFlattenAccountRoles(t *testing.T) {
	roles := []*qci.Role{
		{
			Id:          "role-id-1",
			Name:        "Admin",
			Description: "Administrator role",
			RoleType:    qci.RoleType_ROLE_TYPE_SYSTEM,
			Permissions: []*qci.Permission{
				{Value: "read:clusters"},
				{Value: "write:clusters"},
			},
		},
		{
			Id:          "role-id-2",
			Name:        "Custom",
			Description: "Custom role",
			RoleType:    qci.RoleType_ROLE_TYPE_CUSTOM,
			Permissions: []*qci.Permission{
				{Value: "read:clusters"},
			},
		},
	}

	expected := []interface{}{
		map[string]interface{}{
			rolesIDFieldName:          "role-id-1",
			rolesNameFieldName:        "Admin",
			rolesDescriptionFieldName: "Administrator role",
			rolesRoleTypeFieldName:    "ROLE_TYPE_SYSTEM",
			rolesPermissionsFieldName: []interface{}{"read:clusters", "write:clusters"},
		},
		map[string]interface{}{
			rolesIDFieldName:          "role-id-2",
			rolesNameFieldName:        "Custom",
			rolesDescriptionFieldName: "Custom role",
			rolesRoleTypeFieldName:    "ROLE_TYPE_CUSTOM",
			rolesPermissionsFieldName: []interface{}{"read:clusters"},
		},
	}

	flattened := flattenAccountRoles(roles)
	assert.Equal(t, expected, flattened)
}

func TestFlattenAccountRolesEmpty(t *testing.T) {
	flattened := flattenAccountRoles(nil)
	assert.Equal(t, []interface{}{}, flattened)
	assert.NotNil(t, flattened)
}

func TestFlattenAccountRolesNoPermissions(t *testing.T) {
	roles := []*qci.Role{
		{
			Id:          "role-id-1",
			Name:        "Base",
			Description: "Base role",
			RoleType:    qci.RoleType_ROLE_TYPE_SYSTEM,
		},
	}

	expected := []interface{}{
		map[string]interface{}{
			rolesIDFieldName:          "role-id-1",
			rolesNameFieldName:        "Base",
			rolesDescriptionFieldName: "Base role",
			rolesRoleTypeFieldName:    "ROLE_TYPE_SYSTEM",
			rolesPermissionsFieldName: []interface{}{},
		},
	}

	flattened := flattenAccountRoles(roles)
	assert.Equal(t, expected, flattened)
}
