package qdrant

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	qci "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/iam/v1"
)

// dataSourceAccountsRoles constructs a Terraform data source for
// listing all roles (system and custom) of a Qdrant Cloud account.
func dataSourceAccountsRoles() *schema.Resource {
	return &schema.Resource{
		Description: "Account Roles Data Source. Lists all roles (system and custom) in a Qdrant Cloud account.",
		ReadContext: dataSourceAccountsRolesRead,
		Schema:      accountsRolesDataSourceSchema(),
	}
}

// dataSourceAccountsRolesRead performs a read operation to fetch all roles associated with a specific account.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func dataSourceAccountsRolesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error listing account roles"
	client, clientCtx, diags := getServiceClient(ctx, m, qci.NewIAMServiceClient)
	if diags.HasError() {
		return diags
	}
	// Get the account ID as UUID.
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// List all roles for the provided account.
	var trailer metadata.MD
	resp, err := client.ListRoles(clientCtx, &qci.ListRolesRequest{
		AccountId: accountUUID.String(),
	}, grpc.Trailer(&trailer))
	// Enrich prefix with request ID.
	errorPrefix += getRequestID(trailer)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Flatten roles and store in Terraform state.
	if err := d.Set("roles", flattenAccountRoles(resp.GetItems())); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	d.SetId(time.Now().UTC().Format(time.RFC3339))
	return nil
}

// flattenAccountRoles converts a list of Role proto messages into a list of maps for Terraform state.
func flattenAccountRoles(roles []*qci.Role) []interface{} {
	var result []interface{}
	for _, role := range roles {
		permissions := make([]interface{}, 0, len(role.GetPermissions()))
		for _, perm := range role.GetPermissions() {
			permissions = append(permissions, perm.GetValue())
		}
		result = append(result, map[string]interface{}{
			"id":          role.GetId(),
			"name":        role.GetName(),
			"description": role.GetDescription(),
			"role_type":   role.GetRoleType().String(),
			"permissions": permissions,
		})
	}
	return result
}

// accountsRolesDataSourceSchema defines the Terraform schema for the accounts_roles data source.
func accountsRolesDataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"account_id": {
			Description: "The account ID (UUID). Defaults to the provider-level account_id.",
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
		},
		"roles": {
			Description: "List of account roles (system and custom).",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Description: "Unique identifier for the role (UUID).",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"name": {
						Description: "Name of the role.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"description": {
						Description: "Human-readable description of the role.",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"role_type": {
						Description: "Role type (ROLE_TYPE_SYSTEM or ROLE_TYPE_CUSTOM).",
						Type:        schema.TypeString,
						Computed:    true,
					},
					"permissions": {
						Description: "List of permission values (e.g., read:clusters, write:roles).",
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
