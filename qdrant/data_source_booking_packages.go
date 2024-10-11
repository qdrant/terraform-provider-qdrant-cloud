package qdrant

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qc "github.com/qdrant/terraform-provider-qdrant-cloud/v1/internal/client"
)

// dataSourceBookingPackages returns the schema for the data source qdrant_booking_packages.
func dataSourceBookingPackages() *schema.Resource {
	return &schema.Resource{
		Description: "Booking packages Data Source",
		ReadContext: dataBookingPackagesRead,
		Schema:      packagesSchema(),
	}
}

// dataBookingPackagesRead fetches and sets package data from the API into the Terraform state.
// d: The Terraform ResourceData object containing the state.
// m: The Terraform meta object containing the client configuration.
func dataBookingPackagesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	errorPrefix := "error listing packages"
	// Get an authenticated client
	apiClient, diagnostics := getClient(m)
	if diagnostics.HasError() {
		return diagnostics
	}

	params := qc.GetPackagesParams{
		Provider: qc.GetPackagesParamsProvider(d.Get("cloud_provider").(string)),
		Region:   qc.GetPackagesParamsRegion(d.Get("cloud_region").(string)),
	}

	// Get all packages
	resp, err := apiClient.GetPackagesWithResponse(ctx, &params)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	if resp.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("%s: %v", errorPrefix, getError(resp.JSON422)))
	}
	if resp.StatusCode() != 200 {
		return diag.FromErr(fmt.Errorf("%s", getErrorMessage(errorPrefix, resp.HTTPResponse)))
	}
	// Flatten packages
	packages := flattenPackages(*resp.JSON200)

	// Set the packages in the Terraform state.
	if err := d.Set("packages", packages); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	d.SetId(time.Now().Format(time.RFC3339))
	return nil
}
