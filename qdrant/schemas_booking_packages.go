package qdrant

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// packageSchema defines the schema structure for a package within the Terraform provider.
func packageSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id":                       {Type: schema.TypeString, Computed: true},
		"name":                     {Type: schema.TypeString, Computed: true},
		"status":                   {Type: schema.TypeInt, Computed: true},
		"currency":                 {Type: schema.TypeString, Computed: true},
		"unit_int_price_per_hour":  {Type: schema.TypeInt, Computed: true},
		"unit_int_price_per_day":   {Type: schema.TypeInt, Computed: true},
		"unit_int_price_per_month": {Type: schema.TypeInt, Computed: true},
		"unit_int_price_per_year":  {Type: schema.TypeInt, Computed: true},
		"regional_mapping_id":      {Type: schema.TypeString, Computed: true},
		"resource_configuration": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: resourceConfigurationSchema(),
			},
		},
	}
}

// resourceOptionSchema returns the schema for individual resource options.
func resourceOptionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id":                       {Type: schema.TypeString, Computed: true},
		"resource_type":            {Type: schema.TypeString, Computed: true},
		"status":                   {Type: schema.TypeInt, Computed: true},
		"name":                     {Type: schema.TypeString, Computed: true},
		"resource_unit":            {Type: schema.TypeString, Computed: true},
		"currency":                 {Type: schema.TypeString, Computed: true},
		"unit_int_price_per_hour":  {Type: schema.TypeInt, Computed: true},
		"unit_int_price_per_day":   {Type: schema.TypeInt, Computed: true},
		"unit_int_price_per_month": {Type: schema.TypeInt, Computed: true},
		"unit_int_price_per_year":  {Type: schema.TypeInt, Computed: true},
	}
}

// mapPackages maps API package data to Terraform-readable format.
// packages: A list of packages from the API.
func mapPackages(packages []PackageOut) []interface{} {
	var result []interface{}
	mappedPackages := make([]map[string]interface{}, len(packages))

	for i, pkg := range packages {
		mappedPackages[i] = map[string]interface{}{
			"id":                       pkg.ID,
			"name":                     pkg.Name,
			"status":                   pkg.Status,
			"currency":                 pkg.Currency,
			"unit_int_price_per_hour":  pkg.UnitIntPricePerHour,
			"unit_int_price_per_day":   pkg.UnitIntPricePerDay,
			"unit_int_price_per_month": pkg.UnitIntPricePerMonth,
			"unit_int_price_per_year":  pkg.UnitIntPricePerYear,
			"regional_mapping_id":      pkg.RegionalMappingID,
			"resource_configuration":   mapResourceConfigurations(pkg.ResourceConfiguration),
		}
	}
	return result
}

// mapResourceConfigurations maps API resource configurations to Terraform-readable format.
// rcs: A list of resource configurations.
func mapResourceConfigurations(rcs []ResourceConfiguration) []interface{} {
	var out []interface{}
	for _, rc := range rcs {
		out = append(out, map[string]interface{}{
			"resource_option_id": rc.ResourceOptionID,
			"resource_option": map[string]interface{}{
				"id":                       rc.ResourceOption.ID,
				"resource_type":            rc.ResourceOption.ResourceType,
				"status":                   rc.ResourceOption.Status,
				"name":                     rc.ResourceOption.Name,
				"resource_unit":            rc.ResourceOption.ResourceUnit,
				"currency":                 rc.ResourceOption.Currency,
				"unit_int_price_per_hour":  rc.ResourceOption.UnitIntPricePerHour,
				"unit_int_price_per_day":   rc.ResourceOption.UnitIntPricePerDay,
				"unit_int_price_per_month": rc.ResourceOption.UnitIntPricePerMonth,
				"unit_int_price_per_year":  rc.ResourceOption.UnitIntPricePerYear,
			},
		})
	}
	return out
}

// resourceConfigurationSchema defines the schema structure for resource configurations.
func resourceConfigurationSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"resource_option_id": {Type: schema.TypeString, Computed: true},
		"resource_option": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: resourceOptionSchema(),
			},
		},
	}
}
