package qdrant

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataBookingPackages returns the schema for the data source qdrant_booking_packages.
func dataBookingPackages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataBookingPackagesRead,
		Schema: map[string]*schema.Schema{
			"packages": {
				Type:     schema.TypeList,
				Computed: true,
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
	client := m.(ClientConfig)

	// Fetch the list of packages from the API.
	req, diags := newQdrantCloudRequest(client, "GET", "/booking/packages", nil)
	if diags.HasError() {
		return diags
	}

	// Execute the request and handle the response
	resp, diags := ExecuteRequest(client, req.WithContext(ctx))
	if diags.HasError() {
		return diags
	}

	// Decode the JSON response into the result structure.
	var packages []PackageOut
	if err := json.NewDecoder(resp.Body).Decode(&packages); err != nil {
		return diag.FromErr(fmt.Errorf("error decoding response: %s", err))
	}

	// Set the packages in the Terraform state.
	if err := d.Set("packages", mapPackages(packages)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(time.Now().Format(time.RFC3339))
	return nil
}
