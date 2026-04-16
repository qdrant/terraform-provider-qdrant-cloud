package qdrant

import (
	"testing"

	"github.com/stretchr/testify/assert"

	qca "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/account/v1"
	qci "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/iam/v1"
)

func TestFlattenAccountMembers(t *testing.T) {
	members := []*qca.AccountMember{
		{
			AccountMember: &qci.User{
				Id:     "user-id-1",
				Email:  "owner@example.com",
				Status: qci.UserStatus_USER_STATUS_ACTIVE,
			},
			IsOwner: true,
		},
		{
			AccountMember: &qci.User{
				Id:     "user-id-2",
				Email:  "member@example.com",
				Status: qci.UserStatus_USER_STATUS_ACTIVE,
			},
			IsOwner: false,
		},
	}

	expected := []interface{}{
		map[string]interface{}{
			membersUserIDFieldName:  "user-id-1",
			membersEmailFieldName:   "owner@example.com",
			membersStatusFieldName:  "USER_STATUS_ACTIVE",
			membersIsOwnerFieldName: true,
		},
		map[string]interface{}{
			membersUserIDFieldName:  "user-id-2",
			membersEmailFieldName:   "member@example.com",
			membersStatusFieldName:  "USER_STATUS_ACTIVE",
			membersIsOwnerFieldName: false,
		},
	}

	flattened := flattenAccountMembers(members)
	assert.Equal(t, expected, flattened)
}

func TestFlattenAccountMembersEmpty(t *testing.T) {
	flattened := flattenAccountMembers(nil)
	assert.Equal(t, []interface{}{}, flattened)
	assert.NotNil(t, flattened)
}
