package qdrant

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	qca "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/account/v1"
)

// dataSourceAccountsMembers constructs a Terraform data source for
// listing all members of a Qdrant Cloud account.
func dataSourceAccountsMembers() *schema.Resource {
	return &schema.Resource{
		Description: "Account Members Data Source. Lists all members in a Qdrant Cloud account.",
		ReadContext: dataSourceAccountsMembersRead,
		Schema:      accountsMembersDataSourceSchema(),
	}
}

// dataSourceAccountsMembersRead performs a read operation to fetch all members associated with a specific account.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func dataSourceAccountsMembersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error listing account members"
	client, clientCtx, diags := getServiceClient(ctx, m, qca.NewAccountServiceClient)
	if diags.HasError() {
		return diags
	}
	// Get the account ID as UUID.
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// List all members for the provided account.
	var trailer metadata.MD
	resp, err := client.ListAccountMembers(clientCtx, &qca.ListAccountMembersRequest{
		AccountId: accountUUID.String(),
	}, grpc.Trailer(&trailer))
	// Enrich prefix with request ID.
	errorPrefix += getRequestID(trailer)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Flatten members and store in Terraform state.
	if err := d.Set("members", flattenAccountMembers(resp.GetItems())); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	d.SetId(time.Now().UTC().Format(time.RFC3339))
	return nil
}

// flattenAccountMembers converts a list of AccountMember proto messages into a list of maps for Terraform state.
func flattenAccountMembers(members []*qca.AccountMember) []interface{} {
	var result []interface{}
	for _, member := range members {
		user := member.GetAccountMember()
		result = append(result, map[string]interface{}{
			"user_id":  user.GetId(),
			"email":    user.GetEmail(),
			"status":   user.GetStatus().String(),
			"is_owner": member.GetIsOwner(),
		})
	}
	return result
}

// accountsMembersDataSourceSchema defines the Terraform schema for the accounts_members data source.
func accountsMembersDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_id": {
			Description: "The account ID (UUID). Defaults to the provider-level account_id.",
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
		},
		"members": {
			Description: "List of account members.",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"user_id": {
						Description: "Unique identifier for the user (UUID).",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"email": {
						Description: "Email address of the user.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"status": {
						Description: "User status (USER_STATUS_ACTIVE, USER_STATUS_BLOCKED, USER_STATUS_DELETED).",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"is_owner": {
						Description: "Whether the user is the account owner.",
						Type:        schema.TypeBool,
						Computed:    true,
					},
				},
			},
		},
	}
}
