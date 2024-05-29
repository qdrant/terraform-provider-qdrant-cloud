package qdrant

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qc "terraform-provider-qdrant-cloud/v1/internal/client"
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
	// Get an authenticated client
	apiClient, diagnostics := getClient(m)
	if diagnostics.HasError() {
		return diagnostics
	}
	// Get all packages
	response, err := apiClient.GetPackagesWithResponse(ctx, &qc.GetPackagesParams{})
	if err != nil {
		return diag.FromErr(fmt.Errorf("error listing packages: %v", err))
	}
	if response.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("error listing packages: %v", getError(response.JSON422)))
	}
	if response.StatusCode() != 200 {
		return diag.FromErr(fmt.Errorf("error listing packages: [%d] - %s", response.StatusCode(), response.Status()))
	}
	// Flatten packages
	packages := flattenPackages(*response.JSON200)
	// Set the packages in the Terraform state.
	if err := d.Set("packages", packages); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(time.Now().Format(time.RFC3339))
	return nil
}
