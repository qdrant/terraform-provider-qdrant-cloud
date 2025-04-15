package qdrant

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	qcBooking "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/booking/v1"
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
	// Get a client connection and context
	apiClientConn, clientCtx, diagnostics := getClientConnection(ctx, m)
	if diagnostics.HasError() {
		return diagnostics
	}
	// Get a client
	client := qcBooking.NewBookingServiceClient(apiClientConn)
	// Get The account ID as UUID
	accountUUID, err := getAccountUUID(d, m)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Get all packages
	var header metadata.MD
	resp, err := client.ListPackages(clientCtx, &qcBooking.ListPackagesRequest{
		AccountId:     accountUUID.String(),
		CloudProvider: newPointer(d.Get("cloud_provider").(string)),
		CloudRegion:   d.Get("cloud_region").(string),
	}, grpc.Header(&header))
	// enrich prefix with request ID
	errorPrefix += getRequestID(header)
	if err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	// Flatten packages
	packages := flattenPackages(resp.GetItems())

	// Set the packages in the Terraform state.
	if err := d.Set("packages", packages); err != nil {
		return diag.FromErr(fmt.Errorf("%s: %w", errorPrefix, err))
	}
	d.SetId(time.Now().Format(time.RFC3339))
	return nil
}
