package qdrant

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	qc "terraform-provider-qdrant-cloud/v1/internal/client"
)

func TestFlattenPackages(t *testing.T) {
	packages := []qc.PackageSchema{
		{
			Id:   uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			Name: "packageName1",

			Currency:            "USD",
			UnitIntPricePerHour: newPointer(10),
			ResourceConfigurations: []qc.ResourceConfigurationSchema{
				{
					Amount:       1,
					ResourceType: "type1",
					ResourceUnit: "unit1",
				},
			},
		},
		{
			Id:                  uuid.MustParse("00000000-0000-0000-0000-000000000002"),
			Name:                "packageName2",
			Currency:            "EUR",
			UnitIntPricePerHour: newPointer(20),
			ResourceConfigurations: []qc.ResourceConfigurationSchema{
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
			fieldUnitIntPricePerHour: 10,
			fieldResourceConfigurations: []interface{}{
				map[string]interface{}{
					fieldAmount:       1,
					fieldResourceType: "type1",
					fieldResourceUnit: "unit1",
				},
			},
		},
		map[string]interface{}{
			fieldID:                  "00000000-0000-0000-0000-000000000002",
			fieldName:                "packageName2",
			fieldCurrency:            "EUR",
			fieldUnitIntPricePerHour: 20,
			fieldResourceConfigurations: []interface{}{
				map[string]interface{}{
					fieldAmount:       2,
					fieldResourceType: "type2",
					fieldResourceUnit: "unit2",
				},
			},
		},
	}

	assert.Equal(t, expected, flattenPackages(packages))
}

func TestFlattenResourceConfiguration(t *testing.T) {
	rcs := []qc.ResourceConfigurationSchema{
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
