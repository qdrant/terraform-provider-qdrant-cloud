package qdrant

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	spacev1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/serverless/space/v1"
)

// dataSourceServerlessSpaces lists the serverless spaces of an account.
func dataSourceServerlessSpaces() *schema.Resource {
	return &schema.Resource{
		Description: "Serverless Spaces Data Source",
		ReadContext: dataSourceServerlessSpacesRead,
		Schema:      serverlessSpacesDataSourceSchema(),
	}
}

func dataSourceServerlessSpacesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error listing serverless spaces"
	client, clientCtx, diags := getServiceClient(ctx, m, spacev1.NewSpaceServiceClient)
	if diags.HasError() {
		return diags
	}
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	var cloudRegionID *string
	if v, ok := d.GetOk(slSpaceCloudRegionIDFieldName); ok {
		cloudRegionID = newPointer(v.(string))
	}
	var trailer metadata.MD
	spaces, err := listAllPages(func(pageToken *string) ([]*spacev1.Space, string, error) {
		resp, err := client.ListSpaces(clientCtx, &spacev1.ListSpacesRequest{
			AccountId:     accountUUID.String(),
			CloudRegionId: cloudRegionID,
			PageSize:      newPointer(serverlessListPageSize),
			PageToken:     pageToken,
		}, grpc.Trailer(&trailer))
		if err != nil {
			return nil, "", err
		}
		return resp.GetItems(), resp.GetNextPageToken(), nil
	})
	if err != nil {
		return grpcErrorDiag(errorPrefix, getRequestID(trailer), err)
	}
	if err := d.Set(slSpacesFieldName, flattenServerlessSpaces(spaces)); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	if err := d.Set(slSpaceAccountIDFieldName, accountUUID.String()); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	d.SetId(fmt.Sprintf("%s/%s", accountUUID.String(), d.Get(slSpaceCloudRegionIDFieldName).(string)))
	return nil
}
