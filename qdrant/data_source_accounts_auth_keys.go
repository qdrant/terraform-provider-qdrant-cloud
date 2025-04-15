package qdrant

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	qcAuth "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/auth/v1"
)

// dataSourceAccountsAuthKeys constructs a Terraform resource for managing the reading of API keys associated with an account.
// Returns a schema.Resource pointer configured with schema definitions and the read context function.
func dataSourceAccountsAuthKeys() *schema.Resource {
	return &schema.Resource{
		Description: "Account AuthKey Data Source",
		ReadContext: dataAccountsAuthKeysRead,
		Schema:      accountsAuthKeysSchema(),
	}
}

// dataAccountsAuthKeysRead performs a read operation to fetch all API keys associated with a specific account.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func dataAccountsAuthKeysRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error listing API Keys"
	// Get a client connection and context
	apiClientConn, clientCtx, diagnostics := getClientConnection(ctx, m)
	if diagnostics.HasError() {
		return diagnostics
	}
	// Get a client
	client := qcAuth.NewDatabaseApiKeyServiceClient(apiClientConn)
	// Get The account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// List the API Keys for the provided account
	var header metadata.MD
	resp, err := client.ListDatabaseApiKeys(clientCtx, &qcAuth.ListDatabaseApiKeysRequest{
		AccountId: accountUUID.String(),
	}, grpc.Header(&header))
	// enrich prefix with request ID
	errorPrefix += getRequestID(header)
	// Handle the response in case of error
	if err != nil {
		d := diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
		if d.HasError() {
			return d
		}
	}
	// Flatten cluster and store in Terraform state
	if err := d.Set(authKeysKeysFieldName, flattenAuthKeys(resp.GetItems())); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	d.SetId(time.Now().Format(time.RFC3339))
	return nil
}
