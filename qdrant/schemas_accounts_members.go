package qdrant

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qca "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/account/v1"
)

const (
	membersFieldTemplate = "Account Members Schema %s field"

	membersAccountIDFieldName = "account_id"
	membersMembersFieldName   = "members"
	membersUserIDFieldName    = "user_id"
	membersEmailFieldName     = "email"
	membersStatusFieldName    = "status"
	membersIsOwnerFieldName   = "is_owner"
)

// accountsMembersDataSourceSchema defines the Terraform schema for the accounts_members data source.
func accountsMembersDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		membersAccountIDFieldName: {
			Description: "The account ID (UUID). Defaults to the provider-level account_id.",
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
		},
		membersMembersFieldName: {
			Description: "List of account members.",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					membersUserIDFieldName: {
						Description: fmt.Sprintf(membersFieldTemplate, "User ID (UUID)"),
						Type:        schema.TypeString,
						Computed:    true,
					},
					membersEmailFieldName: {
						Description: fmt.Sprintf(membersFieldTemplate, "Email address"),
						Type:        schema.TypeString,
						Computed:    true,
					},
					membersStatusFieldName: {
						Description: fmt.Sprintf(membersFieldTemplate, "User status"),
						Type:        schema.TypeString,
						Computed:    true,
					},
					membersIsOwnerFieldName: {
						Description: fmt.Sprintf(membersFieldTemplate, "Whether the user is the account owner"),
						Type:        schema.TypeBool,
						Computed:    true,
					},
				},
			},
		},
	}
}

// flattenAccountMembers converts a list of AccountMember proto messages into a list of maps for Terraform state.
func flattenAccountMembers(members []*qca.AccountMember) []interface{} {
	result := make([]interface{}, 0, len(members))
	for _, member := range members {
		user := member.GetAccountMember()
		result = append(result, map[string]interface{}{
			membersUserIDFieldName:  user.GetId(),
			membersEmailFieldName:   user.GetEmail(),
			membersStatusFieldName:  user.GetStatus().String(),
			membersIsOwnerFieldName: member.GetIsOwner(),
		})
	}
	return result
}
