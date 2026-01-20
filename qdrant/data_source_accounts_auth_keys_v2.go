package qdrant

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	authv2 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/cluster/auth/v2"
)

// dataSourceAccountsAuthKeysV2 constructs a Terraform data source for reading v2 API keys.
// Returns a schema.Resource pointer configured with schema definitions and the read context function.
func dataSourceAccountsAuthKeysV2() *schema.Resource {
	return &schema.Resource{
		Description: "Account Database API Keys Data Source (v2)",
		ReadContext: dataAccountsAuthKeysV2Read,
		Schema:      accountsAuthKeysV2DataSourceSchema(),
	}
}

// dataAccountsAuthKeysV2Read performs a read operation to fetch all v2 API keys for a specific cluster.
// ctx: Context to carry deadlines, cancellation signals, and other request-scoped values across API calls.
// d: Resource data which is used to manage the state of the resource.
// m: The interface where the configured client is passed.
// Returns diagnostic information encapsulating any runtime issues encountered during the API call.
func dataAccountsAuthKeysV2Read(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error listing API Keys (v2)"
	client, clientCtx, diags := getServiceClient(ctx, m, authv2.NewDatabaseApiKeyServiceClient)
	if diags.HasError() {
		return diags
	}
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	clusterID := d.Get(authKeysV2ClusterIDFieldName).(string)

	var trailer metadata.MD
	resp, err := client.ListDatabaseApiKeys(clientCtx, &authv2.ListDatabaseApiKeysRequest{
		AccountId: accountUUID.String(),
		ClusterId: &clusterID,
	}, grpc.Trailer(&trailer))
	errorPrefix += getRequestID(trailer)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	if err := d.Set(authKeysV2KeysFieldName, flattenAuthKeysV2(resp.GetItems(), false)); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	if err := d.Set(authKeysV2AccountIDFieldName, accountUUID.String()); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	if err := d.Set(authKeysV2ClusterIDFieldName, clusterID); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	d.SetId(fmt.Sprintf("%s/%s", accountUUID.String(), clusterID))
	return nil
}
