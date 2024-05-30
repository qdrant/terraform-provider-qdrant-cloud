package qdrant

import (
	"testing"

	"github.com/stretchr/testify/assert"

	qc "terraform-provider-qdrant-cloud/v1/internal/client"
)

func TestFlattenPackages(t *testing.T) {
	packages := []qc.PackageOut{
		{
			Id:                   newString("packageID1"),
			Name:                 "packageName1",
			Status:               1,
			Currency:             "USD",
			UnitIntPricePerHour:  newInt(10),
			UnitIntPricePerDay:   newInt(200),
			UnitIntPricePerMonth: newInt(6000),
			UnitIntPricePerYear:  newInt(72000),
			RegionalMappingId:    newString("regionID1"),
			ResourceConfiguration: []qc.ResourceConfiguration{
				{
					ResourceOptionId: "optionID1",
					Amount:           1,
					ResourceOption: &qc.ResourceOptionOut{
						Id: "resourceOption1",
					},
				},
			},
		},
		{
			Id:                   newString("packageID2"),
			Name:                 "packageName2",
			Status:               2,
			Currency:             "EUR",
			UnitIntPricePerHour:  newInt(20),
			UnitIntPricePerDay:   newInt(400),
			UnitIntPricePerMonth: newInt(12000),
			UnitIntPricePerYear:  newInt(144000),
			RegionalMappingId:    newString("regionID2"),
			ResourceConfiguration: []qc.ResourceConfiguration{
				{
					ResourceOptionId: "optionID2",
					Amount:           2,
					ResourceOption: &qc.ResourceOptionOut{
						Id: "resourceOption2",
					},
				},
			},
		},
	}

	expected := []interface{}{
		map[string]interface{}{
			"id":                       "packageID1",
			"name":                     "packageName1",
			"status":                   1,
			"currency":                 "USD",
			"unit_int_price_per_hour":  10,
			"unit_int_price_per_day":   200,
			"unit_int_price_per_month": 6000,
			"unit_int_price_per_year":  72000,
			"regional_mapping_id":      "regionID1",
			"resource_configuration": []interface{}{
				map[string]interface{}{
					"resource_option_id": "optionID1",
					"amount":             1,
					"resource_option": []interface{}{
						map[string]interface{}{
							"id":                       "resourceOption1",
							"currency":                 "",
							"name":                     "",
							"resource_type":            "",
							"resource_unit":            "",
							"status":                   0,
							"unit_int_price_per_hour":  0,
							"unit_int_price_per_day":   0,
							"unit_int_price_per_month": 0,
							"unit_int_price_per_year":  0,
						},
					},
				},
			},
		},
		map[string]interface{}{
			"id":                       "packageID2",
			"name":                     "packageName2",
			"status":                   2,
			"currency":                 "EUR",
			"unit_int_price_per_hour":  20,
			"unit_int_price_per_day":   400,
			"unit_int_price_per_month": 12000,
			"unit_int_price_per_year":  144000,
			"regional_mapping_id":      "regionID2",
			"resource_configuration": []interface{}{
				map[string]interface{}{
					"resource_option_id": "optionID2",
					"amount":             2,
					"resource_option": []interface{}{
						map[string]interface{}{
							"id":                       "resourceOption2",
							"currency":                 "",
							"name":                     "",
							"resource_type":            "",
							"resource_unit":            "",
							"status":                   0,
							"unit_int_price_per_hour":  0,
							"unit_int_price_per_day":   0,
							"unit_int_price_per_month": 0,
							"unit_int_price_per_year":  0,
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expected, flattenPackages(packages))
}

func TestFlattenResourceConfiguration(t *testing.T) {
	rcs := []qc.ResourceConfiguration{
		{
			ResourceOptionId: "optionID1",
			Amount:           1,
			ResourceOption: &qc.ResourceOptionOut{
				Id:                   "resOptionID1",
				ResourceType:         "CPU",
				Status:               1,
				Name:                 newString("resourceName1"),
				ResourceUnit:         "unit1",
				Currency:             "USD",
				UnitIntPricePerHour:  newInt(10),
				UnitIntPricePerDay:   newInt(200),
				UnitIntPricePerMonth: newInt(6000),
				UnitIntPricePerYear:  newInt(72000),
			},
		},
		{
			ResourceOptionId: "optionID2",
			Amount:           2,
			ResourceOption: &qc.ResourceOptionOut{
				Id:                   "resOptionID2",
				ResourceType:         "Memory",
				Status:               2,
				Name:                 newString("resourceName2"),
				ResourceUnit:         "unit2",
				Currency:             "EUR",
				UnitIntPricePerHour:  newInt(20),
				UnitIntPricePerDay:   newInt(400),
				UnitIntPricePerMonth: newInt(12000),
				UnitIntPricePerYear:  newInt(144000),
			},
		},
	}

	expected := []interface{}{
		map[string]interface{}{
			"resource_option_id": "optionID1",
			"amount":             1,
			"resource_option": []interface{}{
				map[string]interface{}{
					"id":                       "resOptionID1",
					"resource_type":            "CPU",
					"status":                   1,
					"name":                     "resourceName1",
					"resource_unit":            "unit1",
					"currency":                 "USD",
					"unit_int_price_per_hour":  10,
					"unit_int_price_per_day":   200,
					"unit_int_price_per_month": 6000,
					"unit_int_price_per_year":  72000,
				},
			},
		},
		map[string]interface{}{
			"resource_option_id": "optionID2",
			"amount":             2,
			"resource_option": []interface{}{
				map[string]interface{}{
					"id":                       "resOptionID2",
					"resource_type":            "Memory",
					"status":                   2,
					"name":                     "resourceName2",
					"resource_unit":            "unit2",
					"currency":                 "EUR",
					"unit_int_price_per_hour":  20,
					"unit_int_price_per_day":   400,
					"unit_int_price_per_month": 12000,
					"unit_int_price_per_year":  144000,
				},
			},
		},
	}

	assert.Equal(t, expected, flattenResourceConfiguraton(rcs))
}

func TestFlattenResourceOption(t *testing.T) {
	ro := &qc.ResourceOptionOut{
		Id:                   "resOptionID",
		ResourceType:         "CPU",
		Status:               1,
		Name:                 newString("resourceName"),
		ResourceUnit:         "unit",
		Currency:             "USD",
		UnitIntPricePerHour:  newInt(10),
		UnitIntPricePerDay:   newInt(200),
		UnitIntPricePerMonth: newInt(6000),
		UnitIntPricePerYear:  newInt(72000),
	}

	expected := []interface{}{
		map[string]interface{}{
			"id":                       "resOptionID",
			"resource_type":            "CPU",
			"status":                   1,
			"name":                     "resourceName",
			"resource_unit":            "unit",
			"currency":                 "USD",
			"unit_int_price_per_hour":  10,
			"unit_int_price_per_day":   200,
			"unit_int_price_per_month": 6000,
			"unit_int_price_per_year":  72000,
		},
	}

	assert.Equal(t, expected, flattenResourceOption(ro))
}
