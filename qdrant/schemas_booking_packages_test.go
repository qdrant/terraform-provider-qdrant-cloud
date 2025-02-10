package qdrant

import (
	"testing"

	"github.com/stretchr/testify/assert"

	qcBooking "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/booking/v2"
)

func TestFlattenPackages(t *testing.T) {
	packages := []*qcBooking.Package{
		{
			Id:   "00000000-0000-0000-0000-000000000001",
			Name: "packageName1",

			Currency:            "USD",
			UnitIntPricePerHour: 10,
			ResourceConfiguration: []*qcBooking.ResourceConfiguration{
				{
					Amount:       1,
					ResourceType: "type1",
					ResourceUnit: "unit1",
				},
			},
		},
		{
			Id:                  "00000000-0000-0000-0000-000000000002",
			Name:                "packageName2",
			Currency:            "EUR",
			UnitIntPricePerHour: 20,
			ResourceConfiguration: []*qcBooking.ResourceConfiguration{
				{
					Amount:       2,
					ResourceType: "type2",
					ResourceUnit: "unit2",
				},
			},
		},
	}

	expected := []interface{}{
		map[string]interface{}{
			fieldID:                  "00000000-0000-0000-0000-000000000001",
			fieldName:                "packageName1",
			fieldCurrency:            "USD",
			fieldUnitIntPricePerHour: int32(10),
			fieldResourceConfigurations: []interface{}{
				map[string]interface{}{
					fieldAmount:       int32(1),
					fieldResourceType: "type1",
					fieldResourceUnit: "unit1",
				},
			},
		},
		map[string]interface{}{
			fieldID:                  "00000000-0000-0000-0000-000000000002",
			fieldName:                "packageName2",
			fieldCurrency:            "EUR",
			fieldUnitIntPricePerHour: int32(20),
			fieldResourceConfigurations: []interface{}{
				map[string]interface{}{
					fieldAmount:       int32(2),
					fieldResourceType: "type2",
					fieldResourceUnit: "unit2",
				},
			},
		},
	}

	assert.Equal(t, expected, flattenPackages(packages))
}

func TestFlattenResourceConfiguration(t *testing.T) {
	rcs := []*qcBooking.ResourceConfiguration{
		{
			Amount:       1,
			ResourceType: "type1",
			ResourceUnit: "unit1",
		},
		{
			Amount:       2,
			ResourceType: "type2",
			ResourceUnit: "unit2",
		},
	}

	expected := []interface{}{
		map[string]interface{}{
			fieldAmount:       1,
			fieldResourceType: "type1",
			fieldResourceUnit: "unit1",
		},
		map[string]interface{}{
			fieldAmount:       2,
			fieldResourceType: "type2",
			fieldResourceUnit: "unit2",
		},
	}

	assert.Equal(t, expected, flattenResourceConfigurations(rcs))
}
