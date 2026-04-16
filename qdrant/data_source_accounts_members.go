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
	if err := d.Set(membersMembersFieldName, flattenAccountMembers(resp.GetItems())); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	if err := d.Set(membersAccountIDFieldName, accountUUID.String()); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	d.SetId(time.Now().UTC().Format(time.RFC3339))
	return nil
}
