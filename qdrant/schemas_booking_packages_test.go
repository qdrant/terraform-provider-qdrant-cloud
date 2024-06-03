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
			fieldID:                   "packageID1",
			fieldName:                 "packageName1",
			fieldStatus:               1,
			fieldCurrency:             "USD",
			fieldUnitIntPricePerHour:  10,
			fieldUnitIntPricePerDay:   200,
			fieldUnitIntPricePerMonth: 6000,
			fieldUnitIntPricePerYear:  72000,
			fieldRegionalMappingID:    "regionID1",
			fieldResourceConfiguration: []interface{}{
				map[string]interface{}{
					fieldResourceOptionID: "optionID1",
					fieldAmount:           1,
					fieldResourceOption: []interface{}{
						map[string]interface{}{
							fieldID:                   "resourceOption1",
							fieldCurrency:             "",
							fieldName:                 "",
							fieldResourceType:         "",
							fieldResourceUnit:         "",
							fieldStatus:               0,
							fieldUnitIntPricePerHour:  0,
							fieldUnitIntPricePerDay:   0,
							fieldUnitIntPricePerMonth: 0,
							fieldUnitIntPricePerYear:  0,
						},
					},
				},
			},
		},
		map[string]interface{}{
			fieldID:                   "packageID2",
			fieldName:                 "packageName2",
			fieldStatus:               2,
			fieldCurrency:             "EUR",
			fieldUnitIntPricePerHour:  20,
			fieldUnitIntPricePerDay:   400,
			fieldUnitIntPricePerMonth: 12000,
			fieldUnitIntPricePerYear:  144000,
			fieldRegionalMappingID:    "regionID2",
			fieldResourceConfiguration: []interface{}{
				map[string]interface{}{
					fieldResourceOptionID: "optionID2",
					fieldAmount:           2,
					fieldResourceOption: []interface{}{
						map[string]interface{}{
							fieldID:                   "resourceOption2",
							fieldCurrency:             "",
							fieldName:                 "",
							fieldResourceType:         "",
							fieldResourceUnit:         "",
							fieldStatus:               0,
							fieldUnitIntPricePerHour:  0,
							fieldUnitIntPricePerDay:   0,
							fieldUnitIntPricePerMonth: 0,
							fieldUnitIntPricePerYear:  0,
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
			fieldResourceOptionID: "optionID1",
			fieldAmount:           1,
			fieldResourceOption: []interface{}{
				map[string]interface{}{
					fieldID:                   "resOptionID1",
					fieldResourceType:         "CPU",
					fieldStatus:               1,
					fieldName:                 "resourceName1",
					fieldResourceUnit:         "unit1",
					fieldCurrency:             "USD",
					fieldUnitIntPricePerHour:  10,
					fieldUnitIntPricePerDay:   200,
					fieldUnitIntPricePerMonth: 6000,
					fieldUnitIntPricePerYear:  72000,
				},
			},
		},
		map[string]interface{}{
			fieldResourceOptionID: "optionID2",
			fieldAmount:           2,
			fieldResourceOption: []interface{}{
				map[string]interface{}{
					fieldID:                   "resOptionID2",
					fieldResourceType:         "Memory",
					fieldStatus:               2,
					fieldName:                 "resourceName2",
					fieldResourceUnit:         "unit2",
					fieldCurrency:             "EUR",
					fieldUnitIntPricePerHour:  20,
					fieldUnitIntPricePerDay:   400,
					fieldUnitIntPricePerMonth: 12000,
					fieldUnitIntPricePerYear:  144000,
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
			fieldID:                   "resOptionID",
			fieldResourceType:         "CPU",
			fieldStatus:               1,
			fieldName:                 "resourceName",
			fieldResourceUnit:         "unit",
			fieldCurrency:             "USD",
			fieldUnitIntPricePerHour:  10,
			fieldUnitIntPricePerDay:   200,
			fieldUnitIntPricePerMonth: 6000,
			fieldUnitIntPricePerYear:  72000,
		},
	}

	assert.Equal(t, expected, flattenResourceOption(ro))
}
