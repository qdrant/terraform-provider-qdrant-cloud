package qdrant

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qc "terraform-provider-qdrant-cloud/v1/internal/client"
)

// packagesSchema defines the schema structure for a packages within the Terraform provider.
func packagesSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"packages": {
			Description: "TODO",
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
		"id": {
			Description: "TODO",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name":                     {Description: "TODO", Type: schema.TypeString, Computed: true},
		"status":                   {Description: "TODO", Type: schema.TypeInt, Computed: true},
		"currency":                 {Description: "TODO", Type: schema.TypeString, Computed: true},
		"unit_int_price_per_hour":  {Description: "TODO", Type: schema.TypeInt, Computed: true},
		"unit_int_price_per_day":   {Description: "TODO", Type: schema.TypeInt, Computed: true},
		"unit_int_price_per_month": {Description: "TODO", Type: schema.TypeInt, Computed: true},
		"unit_int_price_per_year":  {Description: "TODO", Type: schema.TypeInt, Computed: true},
		"regional_mapping_id":      {Description: "TODO", Type: schema.TypeString, Computed: true},
		"resource_configuration": {
			Description: "TODO",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Description: "TODO",
				Schema:      resourceConfigurationSchema(),
			},
		},
	}
}

// resourceConfigurationSchema defines the schema structure for resource configurations.
func resourceConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"resource_option_id": {Description: "TODO", Type: schema.TypeString, Computed: true},
		"amount":             {Description: "TODO", Type: schema.TypeInt, Computed: true},
		"resource_option": {
			Description: "TODO",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Description: "TODO",
				Schema:      resourceOptionSchema(),
			},
		},
	}
}

// resourceOptionSchema returns the schema for individual resource options.
func resourceOptionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id":                       {Description: "TODO", Type: schema.TypeString, Computed: true},
		"resource_type":            {Description: "TODO", Type: schema.TypeString, Computed: true},
		"status":                   {Description: "TODO", Type: schema.TypeInt, Computed: true},
		"name":                     {Description: "TODO", Type: schema.TypeString, Computed: true},
		"resource_unit":            {Description: "TODO", Type: schema.TypeString, Computed: true},
		"currency":                 {Description: "TODO", Type: schema.TypeString, Computed: true},
		"unit_int_price_per_hour":  {Description: "TODO", Type: schema.TypeInt, Computed: true},
		"unit_int_price_per_day":   {Description: "TODO", Type: schema.TypeInt, Computed: true},
		"unit_int_price_per_month": {Description: "TODO", Type: schema.TypeInt, Computed: true},
		"unit_int_price_per_year":  {Description: "TODO", Type: schema.TypeInt, Computed: true},
	}
}

// flattenPackages flattens the package data into a format that Terraform can understand.
func flattenPackages(packages []qc.PackageOut) []interface{} {
	var flattenedPackages []interface{}
	for _, p := range packages {
		flattenedPackages = append(flattenedPackages, map[string]interface{}{
			"id":                       derefString(p.Id),
			"name":                     p.Name,
			"status":                   int(p.Status),
			"currency":                 string(p.Currency),
			"unit_int_price_per_day":   derefInt(p.UnitIntPricePerDay),
			"unit_int_price_per_hour":  derefInt(p.UnitIntPricePerHour),
			"unit_int_price_per_month": derefInt(p.UnitIntPricePerMonth),
			"unit_int_price_per_year":  derefInt(p.UnitIntPricePerYear),
			"regional_mapping_id":      derefString(p.RegionalMappingId),
			"resource_configuration":   flattenResourceConfiguraton(p.ResourceConfiguration),
		})

	}
	return flattenedPackages
}

// flattenResourceConfiguraton flattens the resource configuration data into a format that Terraform can understand.
func flattenResourceConfiguraton(rcs []qc.ResourceConfiguration) []interface{} {
	var flattenedResourceConfigurations []interface{}
	for _, rc := range rcs {
		flattenedResourceConfigurations = append(flattenedResourceConfigurations, map[string]interface{}{
			"resource_option_id": rc.ResourceOptionId,
			"amount":             rc.Amount,
			"resource_option":    flattenResourceOption(rc.ResourceOption),
		})

	}
	return flattenedResourceConfigurations
}

// flattenResourceOption flattens the resource option data into a format that Terraform can understand.
func flattenResourceOption(ro *qc.ResourceOptionOut) []interface{} {
	if ro == nil {
		return []interface{}{}
	}
	flattenedResourceOption := []interface{}{
		map[string]interface{}{
			"id":                       ro.Id,
			"resource_type":            string(ro.ResourceType),
			"status":                   int(ro.Status),
			"name":                     derefString(ro.Name),
			"resource_unit":            ro.ResourceUnit,
			"currency":                 string(ro.Currency),
			"unit_int_price_per_hour":  derefInt(ro.UnitIntPricePerHour),
			"unit_int_price_per_day":   derefInt(ro.UnitIntPricePerDay),
			"unit_int_price_per_month": derefInt(ro.UnitIntPricePerMonth),
			"unit_int_price_per_year":  derefInt(ro.UnitIntPricePerYear),
		},
	}
	return flattenedResourceOption
}
