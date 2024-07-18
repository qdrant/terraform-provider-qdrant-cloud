package qdrant

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qc "terraform-provider-qdrant-cloud/v1/internal/client"
)

// Constant keys and descriptions for schema fields
const (
	// Field keys
	fieldPackages               = "packages"
	fieldID                     = "id"
	fieldName                   = "name"
	fieldCurrency               = "currency"
	fieldUnitIntPricePerHour    = "unit_int_price_per_hour"
	fieldResourceConfigurations = "resource_configurations"
	fieldAmount                 = "amount"
	fieldResourceType           = "resource_type"
	fieldResourceUnit           = "resource_unit"

	// Descriptions
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
				Schema: resourceConfigurationsSchema(),
			},
		},
	}
}

// resourceConfigurationsSchema defines the schema structure for resource configurations.
func resourceConfigurationsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		fieldAmount: {
			Description: descriptionAmount,
			Type:        schema.TypeInt,
			Computed:    true,
		},
		fieldResourceType: {
			Description: descriptionResourceType,
			Type:        schema.TypeString,
			Computed:    true,
		},
		fieldResourceUnit: {
			Description: descriptionResourceUnit,
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

// flattenPackages flattens the package data into a format that Terraform can understand.
func flattenPackages(packages []qc.PackageSchema) []interface{} {
	var flattenedPackages []interface{}
	for _, p := range packages {
		flattenedPackages = append(flattenedPackages, map[string]interface{}{
			fieldID:                     p.Id.String(),
			fieldName:                   p.Name,
			fieldCurrency:               string(p.Currency),
			fieldUnitIntPricePerHour:    derefPointer(p.UnitIntPricePerHour),
			fieldResourceConfigurations: flattenResourceConfiguratons(p.ResourceConfigurations),
		})
	}
	return flattenedPackages
}

// flattenResourceConfiguratons flattens the resource configurations data into a format that Terraform can understand.
func flattenResourceConfiguratons(rcs []qc.ResourceConfigurationSchema) []interface{} {
	var flattenedResourceConfigurations []interface{}
	for _, rc := range rcs {
		flattenedResourceConfigurations = append(flattenedResourceConfigurations, map[string]interface{}{
			fieldAmount:       rc.Amount,
			fieldResourceType: string(rc.ResourceType),
			fieldResourceUnit: string(rc.ResourceUnit),
		})
	}
	return flattenedResourceConfigurations
}
