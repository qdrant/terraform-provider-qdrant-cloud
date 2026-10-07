package qdrant

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	slplatformv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/serverless/platform/v1"
)

const (
	slCloudRegionsFieldTemplate          = "Serverless Cloud Regions Schema %s field"
	slCloudRegionsAccountIDFieldName     = "account_id"
	slCloudRegionsRegionsFieldName       = "regions"
	slCloudRegionIDFieldName             = "id"
	slCloudRegionNameFieldName           = "name"
	slCloudRegionAvailableFieldName      = "available"
	slCloudRegionCountryISOCodeFieldName = "country_iso_code"
	slCloudRegionGeoSubRegionFieldName   = "geographical_sub_region"
)

// dataSourceServerlessCloudRegions lists the cloud regions where serverless spaces can be hosted.
func dataSourceServerlessCloudRegions() *schema.Resource {
	return &schema.Resource{
		Description: "Serverless Cloud Regions Data Source",
		ReadContext: dataSourceServerlessCloudRegionsRead,
		Schema: map[string]*schema.Schema{
			slCloudRegionsAccountIDFieldName: {
				Description: fmt.Sprintf(slCloudRegionsFieldTemplate, "Account ID"),
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			slCloudRegionsRegionsFieldName: {
				Description: fmt.Sprintf(slCloudRegionsFieldTemplate, "List of cloud regions"),
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						slCloudRegionIDFieldName: {
							Description: "Cloud region ID, to be used as `cloud_region_id` of a serverless space.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						slCloudRegionNameFieldName: {
							Description: "Human-readable name of the cloud region (e.g. Frankfurt).",
							Type:        schema.TypeString,
							Computed:    true,
						},
						slCloudRegionAvailableFieldName: {
							Description: "Whether the cloud region is available.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						slCloudRegionCountryISOCodeFieldName: {
							Description: "ISO 3166-1 alpha-2 country code of the region.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						slCloudRegionGeoSubRegionFieldName: {
							Description: "Geographic location grouping of the region (e.g. Europe).",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceServerlessCloudRegionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error listing serverless cloud regions"
	client, clientCtx, diags := getServiceClient(ctx, m, slplatformv1.NewPlatformServiceClient)
	if diags.HasError() {
		return diags
	}
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}

	var trailer metadata.MD
	resp, err := client.ListCloudRegions(clientCtx, &slplatformv1.ListCloudRegionsRequest{
		AccountId: accountUUID.String(),
	}, grpc.Trailer(&trailer))
	if err != nil {
		return grpcErrorDiag(errorPrefix, getRequestID(trailer), err)
	}

	if err := d.Set(slCloudRegionsRegionsFieldName, flattenServerlessCloudRegions(resp.GetItems())); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	if err := d.Set(slCloudRegionsAccountIDFieldName, accountUUID.String()); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	d.SetId(accountUUID.String())
	return nil
}

func flattenServerlessCloudRegions(regions []*slplatformv1.CloudRegion) []interface{} {
	result := make([]interface{}, 0, len(regions))
	for _, r := range regions {
		result = append(result, map[string]interface{}{
			slCloudRegionIDFieldName:             r.GetId(),
			slCloudRegionNameFieldName:           r.GetName(),
			slCloudRegionAvailableFieldName:      r.GetAvailable(),
			slCloudRegionCountryISOCodeFieldName: r.GetCountryIsoCode(),
			slCloudRegionGeoSubRegionFieldName:   r.GetGeographicalSubRegion(),
		})
	}
	return result
}
