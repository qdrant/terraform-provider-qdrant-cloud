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

// dataSourceServerlessSpace retrieves a specific serverless space.
func dataSourceServerlessSpace() *schema.Resource {
	return &schema.Resource{
		Description: "Serverless Space Data Source",
		ReadContext: dataSourceServerlessSpaceRead,
		Schema:      serverlessSpaceSchema(true),
	}
}

func dataSourceServerlessSpaceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error reading serverless space"
	client, clientCtx, diags := getServiceClient(ctx, m, spacev1.NewSpaceServiceClient)
	if diags.HasError() {
		return diags
	}
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	var trailer metadata.MD
	resp, err := client.GetSpace(clientCtx, &spacev1.GetSpaceRequest{
		AccountId: accountUUID.String(),
		SpaceId:   d.Get(slSpaceIDFieldName).(string),
	}, grpc.Trailer(&trailer))
	if err != nil {
		return grpcErrorDiag(errorPrefix, getRequestID(trailer), err)
	}
	if err := setFlattened(d, flattenServerlessSpace(resp.GetSpace())); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	d.SetId(resp.GetSpace().GetId())
	return nil
}
