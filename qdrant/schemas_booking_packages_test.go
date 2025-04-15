package qdrant

import (
	"testing"

	"github.com/stretchr/testify/assert"

	qcBooking "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/booking/v1"
)

func TestFlattenPackages(t *testing.T) {
	packages := []*qcBooking.Package{
		{
			Id:   "00000000-0000-0000-0000-000000000001",
			Name: "packageName1",

			Currency:            "USD",
			UnitIntPricePerHour: 10,
			ResourceConfiguration: &qcBooking.ResourceConfiguration{
				Ram:  "ram_1",
				Cpu:  "cpu_1",
				Disk: "disk_1",
			},
		},
		{
			Id:                  "00000000-0000-0000-0000-000000000002",
			Name:                "packageName2",
			Currency:            "EUR",
			UnitIntPricePerHour: 20,
			ResourceConfiguration: &qcBooking.ResourceConfiguration{
				Ram:  "ram_2",
				Cpu:  "cpu_2",
				Disk: "disk_2",
			},
		},
	}

	expected := []interface{}{
		map[string]interface{}{
			fieldID:                  "00000000-0000-0000-0000-000000000001",
			fieldName:                "packageName1",
			fieldCurrency:            "USD",
			fieldUnitIntPricePerHour: 10,
			fieldResourceConfiguration: map[string]interface{}{
				fieldResourceRam:  "ram_1",
				fieldResourceCpu:  "cpu_1",
				fieldResourceDisk: "disk_1",
			},
		},
		map[string]interface{}{
			fieldID:                  "00000000-0000-0000-0000-000000000002",
			fieldName:                "packageName2",
			fieldCurrency:            "EUR",
			fieldUnitIntPricePerHour: 20,
			fieldResourceConfiguration: map[string]interface{}{
				fieldResourceRam:  "ram_2",
				fieldResourceCpu:  "cpu_2",
				fieldResourceDisk: "disk_2",
			},
		},
	}

	assert.Equal(t, expected, flattenPackages(packages))
}

func TestFlattenResourceConfiguration(t *testing.T) {
	rcs := &qcBooking.ResourceConfiguration{
		Ram:  "ram_3",
		Cpu:  "cpu_3",
		Disk: "disk_3",
	}

	expected :=
		map[string]interface{}{
			fieldResourceRam:  "ram_3",
			fieldResourceCpu:  "cpu_3",
			fieldResourceDisk: "disk_3",
		}

	assert.Equal(t, expected, flattenResourceConfiguration(rcs))
}
