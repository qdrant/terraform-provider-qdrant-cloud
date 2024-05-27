package qdrant

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qc "terraform-provider-qdrant-cloud/v1/internal/client"
)

// dataBookingPackages returns the schema for the data source qdrant_booking_packages.
func dataBookingPackages() *schema.Resource {
	return &schema.Resource{
		Description: "Booking packages Data Source",
		ReadContext: dataBookingPackagesRead,
		Schema: map[string]*schema.Schema{
			"packages": {
				Description: "TODO",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: packageSchema(),
				},
			},
		},
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

	params := qc.GetPackagesParams{}
	response, err := apiClient.GetPackagesWithResponse(ctx, &params)
	if err != nil {
		d := diag.FromErr(fmt.Errorf("error listing packages: %v", err))
		if d.HasError() {
			return d
		}
	}

	if response.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("error listing packages: %v", response.JSON422))
	}

	packages := flattenPackages(*response.JSON200)

	// Set the packages in the Terraform state.
	if err := d.Set("packages", packages); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(time.Now().Format(time.RFC3339))
	return nil
}

// flattenPackages flattens the package data into a format that Terraform can understand.
func flattenPackages(packages []qc.PackageOut) []interface{} {
	var flattenedPackages []interface{}
	for _, p := range packages {
		flattenedPackages = append(flattenedPackages, map[string]interface{}{
			"id":   p.Id,
			"name": p.Name,
		})
	}
	return flattenedPackages
}
