package qdrant

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	qci "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/iam/v1"
)

func TestSchema_RoleAccountIdOptionalComputed(t *testing.T) {
	s := accountsRoleSchema()
	f := s[roleAccountIdFieldName]
	require.NotNil(t, f, "account_id schema must exist")

	assert.True(t, f.Optional, "account_id should be Optional")
	assert.True(t, f.Computed, "account_id should be Computed")
	assert.Nil(t, f.Default, "account_id must NOT set a Default")
}

func TestSchema_RoleNameRequired(t *testing.T) {
	s := accountsRoleSchema()
	f := s[roleNameFieldName]
	require.NotNil(t, f, "name schema must exist")

	assert.Equal(t, schema.TypeString, f.Type)
	assert.True(t, f.Required)
}

func TestSchema_RoleDescriptionOptional(t *testing.T) {
	s := accountsRoleSchema()
	f := s[roleDescriptionFieldName]
	require.NotNil(t, f, "description schema must exist")

	assert.Equal(t, schema.TypeString, f.Type)
	assert.True(t, f.Optional)
}

func TestSchema_RoleTypeComputed(t *testing.T) {
	s := accountsRoleSchema()
	f := s[roleRoleTypeFieldName]
	require.NotNil(t, f, "role_type schema must exist")

	assert.Equal(t, schema.TypeString, f.Type)
	assert.True(t, f.Computed, "role_type should be Computed (provider sets ROLE_TYPE_CUSTOM)")
}

func TestSchema_RolePermissionsSetDefinition(t *testing.T) {
	s := accountsRoleSchema()
	f := s[rolePermissionsFieldName]
	require.NotNil(t, f, "permissions schema must exist")

	assert.Equal(t, schema.TypeSet, f.Type, "permissions should be a set")
	assert.True(t, f.Required, "permissions should be Required")
	assert.Equal(t, 1, f.MinItems, "permissions should require at least one item")

	// Nested schema checks (value required, category computed)
	ns, ok := f.Elem.(*schema.Resource)
	require.True(t, ok, "permissions Elem must be a *schema.Resource")
	ps := ns.Schema
	require.NotNil(t, ps)

	val := ps["value"]
	require.NotNil(t, val, "permissions.value must exist")
	assert.Equal(t, schema.TypeString, val.Type)
	assert.True(t, val.Required, "permissions.value should be Required")

	cat := ps["category"]
	require.NotNil(t, cat, "permissions.category must exist")
	assert.Equal(t, schema.TypeString, cat.Type)
	assert.True(t, cat.Computed, "permissions.category should be Computed (read-only)")
}

func TestExpandRole_Create_UsesProvidedAccountID(t *testing.T) {
	const acct = "00000000-1111-0000-0000-000000000001"

	// Provide a minimal valid config; include category so expandPermissions is happy if it's still read.
	d := schema.TestResourceDataRaw(t, accountsRoleSchema(), map[string]interface{}{
		roleNameFieldName:        "backup-operator",
		roleDescriptionFieldName: "Can create/restore backups",
		rolePermissionsFieldName: []interface{}{
			map[string]interface{}{"value": "read:backups", "category": "Cluster"},
		},
	})

	role, err := expandRole(d, acct, false)
	require.NoError(t, err)

	assert.Equal(t, acct, role.GetAccountId())
	assert.Equal(t, "backup-operator", role.GetName())
	assert.Equal(t, "Can create/restore backups", role.GetDescription())
	assert.Equal(t, qci.RoleType_ROLE_TYPE_CUSTOM, role.GetRoleType())

	perms := role.GetPermissions()
	require.Len(t, perms, 1)
	assert.Equal(t, "read:backups", perms[0].GetValue())
	// Category is server-populated; if provided in config we pass it through, otherwise it may be empty.
}

func TestExpandRole_Update_IncludesID(t *testing.T) {
	const acct = "00000000-2222-0000-0000-000000000002"

	d := schema.TestResourceDataRaw(t, accountsRoleSchema(), map[string]interface{}{
		roleNameFieldName:        "backup-operator",
		roleDescriptionFieldName: "desc",
		rolePermissionsFieldName: []interface{}{
			map[string]interface{}{"value": "read:backups", "category": "Cluster"},
		},
	})
	d.SetId("36badc1e-05ea-48ba-a5db-1df4d471fedb")

	role, err := expandRole(d, acct, true)
	require.NoError(t, err)

	assert.Equal(t, "36badc1e-05ea-48ba-a5db-1df4d471fedb", role.GetId())
	assert.Equal(t, acct, role.GetAccountId())
}

func TestFlattenRole_Minimal(t *testing.T) {
	created := timestamppb.New(time.Date(2025, 10, 8, 9, 28, 52, 0, time.UTC))
	modified := timestamppb.New(time.Date(2025, 10, 8, 9, 28, 52, 0, time.UTC))
	category := "Cluster"

	r := &qci.Role{
		Id:             "059a0aab-5010-43f1-a139-539bb25a5886",
		AccountId:      "222cda33-2c7a-4046-b2cc-0807170aed49",
		Name:           "backup-operator",
		Description:    "desc",
		RoleType:       qci.RoleType_ROLE_TYPE_CUSTOM,
		CreatedAt:      created,
		LastModifiedAt: modified,
		Permissions: []*qci.Permission{
			{Value: "read:backups", Category: &category},
			{Value: "write:backups", Category: &category},
		},
	}

	got := flattenRole(r)

	want := map[string]interface{}{
		roleIdFieldName:             r.GetId(),
		roleAccountIdFieldName:      r.GetAccountId(),
		roleNameFieldName:           r.GetName(),
		roleDescriptionFieldName:    r.GetDescription(),
		roleRoleTypeFieldName:       r.GetRoleType().String(),
		roleCreatedAtFieldName:      formatTime(r.GetCreatedAt()),
		roleLastModifiedAtFieldName: formatTime(r.GetLastModifiedAt()),
		rolePermissionsFieldName: []interface{}{
			map[string]interface{}{"value": "read:backups", "category": "Cluster"},
			map[string]interface{}{"value": "write:backups", "category": "Cluster"},
		},
	}

	assert.Equal(t, want, got)
}
