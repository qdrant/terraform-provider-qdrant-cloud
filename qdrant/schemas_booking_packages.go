package qdrant

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qcBooking "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/booking/v1"
)

// Constant keys and descriptions for schema fields.
const (
	// Field keys.
	fieldPackages              = "packages"
	fieldID                    = "id"
	fieldName                  = "name"
	fieldType                  = "type"
	fieldCurrency              = "currency"
	fieldUnitIntPricePerHour   = "unit_int_price_per_hour"
	fieldStatus                = "status"
	fieldResourceConfiguration = "resource_configuration"
	fieldResourceRam           = "ram"
	fieldResourceCpu           = "cpu"
	fieldResourceDisk          = "disk"

	// Descriptions.
	descriptionPackages              = "List of packages"
	descriptionID                    = "The ID of the package"
	descriptionName                  = "The name of the package"
	descriptionType                  = "The type of the package"
	descriptionCurrency              = "The currency of the package prices"
	descriptionUnitIntPricePerHour   = "The unit price per hour in integer format"
	descriptionStatus                = "The status of the package"
	descriptionResourceConfiguration = "The resource configuration of the package"
	descriptionResourceRam           = "The amount of RAM (e.g., '1GiB')"
	descriptionResourceCpu           = "The amount of CPU (e.g., '1000m' (1 vCPU))"
	descriptionResourceDisk          = "The amount of disk (e.g., '100GiB')"
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
		fieldType: {
			Description: descriptionType,
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
		fieldResourceConfiguration: {
			Description: descriptionResourceConfiguration,
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: resourceConfigurationSchema(true),
			},
		},
		fieldStatus: {
			Description: descriptionStatus,
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

// resourceConfigurationsSchema defines the schema structure for resource configurations.
func resourceConfigurationSchema(asDataSource bool) map[string]*schema.Schema {
	return map[string]*schema.Schema{
		fieldResourceRam: {
			Description: descriptionResourceRam,
			Type:        schema.TypeString,
			Required:    !asDataSource,
			Computed:    asDataSource,
		},
		fieldResourceCpu: {
			Description: descriptionResourceCpu,
			Type:        schema.TypeString,
			Required:    !asDataSource,
			Computed:    asDataSource,
		},
		fieldResourceDisk: {
			Description: descriptionResourceDisk,
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
			fieldID:                    p.GetId(),
			fieldName:                  p.GetName(),
			fieldType:                  p.GetType(),
			fieldCurrency:              p.GetCurrency(),
			fieldUnitIntPricePerHour:   int(p.GetUnitIntPricePerHour()),
			fieldResourceConfiguration: flattenResourceConfiguration(p.GetResourceConfiguration()),
			fieldStatus:                p.GetStatus().String(),
		})
	}
	return flattenedPackages
}

// flattenResourceConfiguration flattens the resource configuration data into a format that Terraform can understand.
func flattenResourceConfiguration(rc *qcBooking.ResourceConfiguration) map[string]interface{} {
	return map[string]interface{}{
		fieldResourceRam:  rc.GetRam(),
		fieldResourceCpu:  rc.GetCpu(),
		fieldResourceDisk: rc.GetDisk(),
	}
}
