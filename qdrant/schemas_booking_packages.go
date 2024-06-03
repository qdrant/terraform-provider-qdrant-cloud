package qdrant

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qc "terraform-provider-qdrant-cloud/v1/internal/client"
)

// Constant keys and descriptions for schema fields
const (
	// Field keys
	fieldPackages              = "packages"
	fieldID                    = "id"
	fieldName                  = "name"
	fieldStatus                = "status"
	fieldCurrency              = "currency"
	fieldUnitIntPricePerHour   = "unit_int_price_per_hour"
	fieldUnitIntPricePerDay    = "unit_int_price_per_day"
	fieldUnitIntPricePerMonth  = "unit_int_price_per_month"
	fieldUnitIntPricePerYear   = "unit_int_price_per_year"
	fieldRegionalMappingID     = "regional_mapping_id"
	fieldResourceConfiguration = "resource_configuration"
	fieldResourceOptionID      = "resource_option_id"
	fieldAmount                = "amount"
	fieldResourceOption        = "resource_option"
	fieldResourceType          = "resource_type"
	fieldResourceUnit          = "resource_unit"

	// Descriptions
	descriptionPackages              = "List of packages"
	descriptionID                    = "The ID of the package"
	descriptionName                  = "The name of the package"
	descriptionStatus                = "The status of the package"
	descriptionCurrency              = "The currency of the package prices"
	descriptionUnitIntPricePerHour   = "The unit price per hour in integer format"
	descriptionUnitIntPricePerDay    = "The unit price per day in integer format"
	descriptionUnitIntPricePerMonth  = "The unit price per month in integer format"
	descriptionUnitIntPricePerYear   = "The unit price per year in integer format"
	descriptionRegionalMappingID     = "The ID of the regional mapping"
	descriptionResourceConfiguration = "The resource configuration of the package"
	descriptionResourceOptionID      = "The ID of the resource option"
	descriptionAmount                = "The amount of the resource"
	descriptionResourceOption        = "The resource option details"
	descriptionResourceType          = "The type of the resource"
	descriptionResourceUnit          = "The unit of the resource"
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
		fieldStatus: {
			Description: descriptionStatus,
			Type:        schema.TypeInt,
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
		fieldUnitIntPricePerDay: {
			Description: descriptionUnitIntPricePerDay,
			Type:        schema.TypeInt,
			Computed:    true,
		},
		fieldUnitIntPricePerMonth: {
			Description: descriptionUnitIntPricePerMonth,
			Type:        schema.TypeInt,
			Computed:    true,
		},
		fieldUnitIntPricePerYear: {
			Description: descriptionUnitIntPricePerYear,
			Type:        schema.TypeInt,
			Computed:    true,
		},
		fieldRegionalMappingID: {
			Description: descriptionRegionalMappingID,
			Type:        schema.TypeString,
			Computed:    true,
		},
		fieldResourceConfiguration: {
			Description: descriptionResourceConfiguration,
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: resourceConfigurationSchema(),
			},
		},
	}
}

// resourceConfigurationSchema defines the schema structure for resource configurations.
func resourceConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		fieldResourceOptionID: {
			Description: descriptionResourceOptionID,
			Type:        schema.TypeString,
			Computed:    true,
		},
		fieldAmount: {
			Description: descriptionAmount,
			Type:        schema.TypeInt,
			Computed:    true,
		},
		fieldResourceOption: {
			Description: descriptionResourceOption,
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: resourceOptionSchema(),
			},
		},
	}
}

// resourceOptionSchema returns the schema for individual resource options.
func resourceOptionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		fieldID: {
			Description: descriptionID,
			Type:        schema.TypeString,
			Computed:    true,
		},
		fieldResourceType: {
			Description: descriptionResourceType,
			Type:        schema.TypeString,
			Computed:    true,
		},
		fieldStatus: {
			Description: descriptionStatus,
			Type:        schema.TypeInt,
			Computed:    true,
		},
		fieldName: {
			Description: descriptionName,
			Type:        schema.TypeString,
			Computed:    true,
		},
		fieldResourceUnit: {
			Description: descriptionResourceUnit,
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
		fieldUnitIntPricePerDay: {
			Description: descriptionUnitIntPricePerDay,
			Type:        schema.TypeInt,
			Computed:    true,
		},
		fieldUnitIntPricePerMonth: {
			Description: descriptionUnitIntPricePerMonth,
			Type:        schema.TypeInt,
			Computed:    true,
		},
		fieldUnitIntPricePerYear: {
			Description: descriptionUnitIntPricePerYear,
			Type:        schema.TypeInt,
			Computed:    true,
		},
	}
}

// ConvertToJSON takes a slice of interface{} and converts it to JSON for debugging purposes.
func ConvertToJSON(data []interface{}) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error converting to JSON")
		return "[]"
	}
	return string(jsonData)
}

// flattenPackages flattens the package data into a format that Terraform can understand.
func flattenPackages(packages []qc.PackageOut) []interface{} {
	var flattenedPackages []interface{}
	for _, p := range packages {
		flattenedPackages = append(flattenedPackages, map[string]interface{}{
			fieldID:                    derefString(p.Id),
			fieldName:                  p.Name,
			fieldStatus:                int(p.Status),
			fieldCurrency:              string(p.Currency),
			fieldUnitIntPricePerDay:    derefInt(p.UnitIntPricePerDay),
			fieldUnitIntPricePerHour:   derefInt(p.UnitIntPricePerHour),
			fieldUnitIntPricePerMonth:  derefInt(p.UnitIntPricePerMonth),
			fieldUnitIntPricePerYear:   derefInt(p.UnitIntPricePerYear),
			fieldRegionalMappingID:     derefString(p.RegionalMappingId),
			fieldResourceConfiguration: flattenResourceConfiguraton(p.ResourceConfiguration),
		})

	}
	return flattenedPackages
}

// flattenResourceConfiguraton flattens the resource configuration data into a format that Terraform can understand.
func flattenResourceConfiguraton(rcs []qc.ResourceConfiguration) []interface{} {
	var flattenedResourceConfigurations []interface{}
	for _, rc := range rcs {
		flattenedResourceConfigurations = append(flattenedResourceConfigurations, map[string]interface{}{
			fieldResourceOptionID: rc.ResourceOptionId,
			fieldAmount:           rc.Amount,
			fieldResourceOption:   flattenResourceOption(rc.ResourceOption),
		})

	}
	return flattenedResourceConfigurations
}

// flattenResourceOption flattens the resource option data into a format that Terraform can understand.
func flattenResourceOption(ro *qc.ResourceOptionOut) []interface{} {
	if ro == nil {
		return nil
	}
	flattenedResourceOption := []interface{}{
		map[string]interface{}{
			fieldID:                   ro.Id,
			fieldResourceType:         string(ro.ResourceType),
			fieldStatus:               int(ro.Status),
			fieldName:                 derefString(ro.Name),
			fieldResourceUnit:         ro.ResourceUnit,
			fieldCurrency:             string(ro.Currency),
			fieldUnitIntPricePerHour:  derefInt(ro.UnitIntPricePerHour),
			fieldUnitIntPricePerDay:   derefInt(ro.UnitIntPricePerDay),
			fieldUnitIntPricePerMonth: derefInt(ro.UnitIntPricePerMonth),
			fieldUnitIntPricePerYear:  derefInt(ro.UnitIntPricePerYear),
		},
	}
	return flattenedResourceOption
}
