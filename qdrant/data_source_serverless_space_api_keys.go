package qdrant

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	spaceauthv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/serverless/space/auth/v1"
)

// dataSourceServerlessSpaceAPIKeys lists the API keys of a serverless space (without the secrets).
func dataSourceServerlessSpaceAPIKeys() *schema.Resource {
	return &schema.Resource{
		Description: "Serverless Space API Keys Data Source",
		ReadContext: dataSourceServerlessSpaceAPIKeysRead,
		Schema:      serverlessSpaceAPIKeysDataSourceSchema(),
	}
}

func dataSourceServerlessSpaceAPIKeysRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error listing serverless space API keys"
	client, clientCtx, diags := getServiceClient(ctx, m, spaceauthv1.NewSpaceApiKeyServiceClient)
	if diags.HasError() {
		return diags
	}
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	spaceID := d.Get(slSpaceAPIKeySpaceIDFieldName).(string)
	var trailer metadata.MD
	resp, err := client.ListSpaceApiKeys(clientCtx, &spaceauthv1.ListSpaceApiKeysRequest{
		AccountId: accountUUID.String(),
		SpaceId:   spaceID,
	}, grpc.Trailer(&trailer))
	if err != nil {
		return grpcErrorDiag(errorPrefix, getRequestID(trailer), err)
	}
	if err := d.Set(slSpaceAPIKeysFieldName, flattenServerlessSpaceAPIKeys(resp.GetItems())); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	if err := d.Set(slSpaceAPIKeyAccountIDFieldName, accountUUID.String()); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	d.SetId(fmt.Sprintf("%s/%s", accountUUID.String(), spaceID))
	return nil
}
