package qdrant

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qcBooking "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/booking/v2"
)

// Constant keys and descriptions for schema fields.
const (
	// Field keys.
	fieldPackages               = "packages"
	fieldID                     = "id"
	fieldName                   = "name"
	fieldCurrency               = "currency"
	fieldUnitIntPricePerHour    = "unit_int_price_per_hour"
	fieldResourceConfigurations = "resource_configurations"
	fieldAmount                 = "amount"
	fieldResourceType           = "resource_type"
	fieldResourceUnit           = "resource_unit"

	// Descriptions.
	descriptionPackages               = "List of packages"
	descriptionID                     = "The ID of the package"
	descriptionName                   = "The name of the package"
	descriptionCurrency               = "The currency of the package prices"
	descriptionUnitIntPricePerHour    = "The unit price per hour in integer format"
	descriptionResourceConfigurations = "The resource configurations of the package"
	descriptionAmount                 = "The amount of the resource"
	descriptionResourceType           = "The type of the resource"
	descriptionResourceUnit           = "The unit of the resource"
)

// packagesSchema defines the schema structure for a packages within the Terraform provider.
func packagesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		fieldPackages: {
			Description: descriptionPackages,
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: packageSchema(),
			},
		},
		clusterCloudProviderFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Cloud provider where the cluster resides"),
			Type:        schema.TypeString,
			Required:    true,
			Computed:    false,
		},
		clusterCloudRegionFieldName: {
			Description: fmt.Sprintf(clusterFieldTemplate, "Cloud region where the cluster resides"),
			Type:        schema.TypeString,
			Required:    true,
			Computed:    false,
		},
	}
}

// packageSchema defines the schema structure for a package within the Terraform provider.
func packageSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		fieldID: {
			Description: descriptionID,
			Type:        schema.TypeString,
			Computed:    true,
		},
		fieldName: {
			Description: descriptionName,
			Type:        schema.TypeString,
			Computed:    true,
		},
		fieldCurrency: {
			Description: descriptionCurrency,
			Type:        schema.TypeString,
			Computed:    true,
		},
		fieldUnitIntPricePerHour: {
			Description: descriptionUnitIntPricePerHour,
			Type:        schema.TypeInt,
			Computed:    true,
		},
		fieldResourceConfigurations: {
			Description: descriptionResourceConfigurations,
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: resourceConfigurationsSchema(true),
			},
		},
	}
}

// resourceConfigurationsSchema defines the schema structure for resource configurations.
func resourceConfigurationsSchema(asDataSource bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		fieldAmount: {
			Description: descriptionAmount,
			Type:        schema.TypeInt,
			Required:    !asDataSource,
			Computed:    asDataSource,
		},
		fieldResourceType: {
			Description: descriptionResourceType,
			Type:        schema.TypeString,
			Required:    !asDataSource,
			Computed:    asDataSource,
		},
		fieldResourceUnit: {
			Description: descriptionResourceUnit,
			Type:        schema.TypeString,
			Required:    !asDataSource,
			Computed:    asDataSource,
		},
	}
}

// flattenPackages flattens the package data into a format that Terraform can understand.
func flattenPackages(packages []*qcBooking.Package) []interface{} {
	var flattenedPackages []interface{}
	for _, p := range packages {
		flattenedPackages = append(flattenedPackages, map[string]interface{}{
			fieldID:                     p.GetId(),
			fieldName:                   p.GetName(),
			fieldCurrency:               p.GetCurrency(),
			fieldUnitIntPricePerHour:    int(p.GetUnitIntPricePerHour()),
			fieldResourceConfigurations: flattenResourceConfigurations(p.GetResourceConfigurations()),
		})
	}
	return flattenedPackages
}

// flattenResourceConfigurations flattens the resource configurations data into a format that Terraform can understand.
func flattenResourceConfigurations(rcs []*qcBooking.ResourceConfiguration) []interface{} {
	var flattenedResourceConfigurations []interface{}
	for _, rc := range rcs {
		flattenedResourceConfigurations = append(flattenedResourceConfigurations, map[string]interface{}{
			fieldAmount:       int(rc.GetAmount()),
			fieldResourceType: rc.GetResourceType(),
			fieldResourceUnit: rc.GetResourceUnit(),
		})
	}
	return flattenedResourceConfigurations
}
