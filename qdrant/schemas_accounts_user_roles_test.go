package qdrant

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchema_UserRolesAccountIdOptionalComputed(t *testing.T) {
	s := accountsUserRolesSchema()
	f := s[userRolesAccountIdFieldName]
	require.NotNil(t, f, "account_id schema must exist")

	assert.Equal(t, schema.TypeString, f.Type)
	assert.True(t, f.Optional, "account_id should be Optional")
	assert.True(t, f.Computed, "account_id should be Computed")
	assert.Nil(t, f.Default, "account_id must NOT set a Default")
}

func TestSchema_UserRolesUserEmailRequired(t *testing.T) {
	s := accountsUserRolesSchema()
	f := s[userRolesUserEmailFieldName]
	require.NotNil(t, f, "user_email schema must exist")

	assert.Equal(t, schema.TypeString, f.Type)
	assert.True(t, f.Required, "user_email should be Required")
	assert.False(t, f.Computed, "user_email should not be Computed")
}

func TestSchema_UserRolesRoleIdsSetRequired(t *testing.T) {
	s := accountsUserRolesSchema()
	f := s[userRolesRoleIdsFieldName]
	require.NotNil(t, f, "role_ids schema must exist")

	assert.Equal(t, schema.TypeSet, f.Type, "role_ids should be a set")
	assert.True(t, f.Required, "role_ids should be Required")
	ns, ok := f.Elem.(*schema.Schema)
	require.True(t, ok, "role_ids Elem must be a *schema.Schema")
	assert.Equal(t, schema.TypeString, ns.Type, "role_ids elements must be strings")
}

func TestSchema_UserRolesKeepOnDestroyDefaultFalse(t *testing.T) {
	s := accountsUserRolesSchema()
	f := s[userRolesKeepOnDestroyName]
	require.NotNil(t, f, "keep_on_destroy schema must exist")

	assert.Equal(t, schema.TypeBool, f.Type)
	assert.True(t, f.Optional, "keep_on_destroy should be Optional")
	assert.EqualValues(t, false, f.Default, "keep_on_destroy default should be false")

	// Sanity check via ResourceData
	d := schema.TestResourceDataRaw(t, s, map[string]interface{}{
		userRolesUserEmailFieldName: "someone@example.com",
		userRolesRoleIdsFieldName:   []interface{}{"11111111-1111-1111-1111-111111111111"},
	})
	// Unset in config → should resolve to default false
	got := d.Get(userRolesKeepOnDestroyName)
	assert.EqualValues(t, false, got, "keep_on_destroy should read back as false by default")
}

func TestSchema_UserRolesUserIdComputed(t *testing.T) {
	s := accountsUserRolesSchema()
	f := s[userRolesUserIdFieldName]
	require.NotNil(t, f, "user_id schema must exist")

	assert.Equal(t, schema.TypeString, f.Type)
	assert.True(t, f.Computed, "user_id should be Computed")
	assert.False(t, f.Optional, "user_id should not be Optional")
	assert.False(t, f.Required, "user_id should not be Required")
}
