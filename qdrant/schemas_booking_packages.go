package qdrant

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
