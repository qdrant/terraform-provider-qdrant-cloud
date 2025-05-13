package qdrant

import (
	"testing"

	"github.com/stretchr/testify/assert"

	qcBooking "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/booking/v1"
)

func TestFlattenPackages(t *testing.T) {
	packages := []*qcBooking.Package{
		{
			Id:                  "00000000-0000-0000-0000-000000000001",
			Name:                "packageName1",
			Type:                "packageType1",
			Currency:            "USD",
			UnitIntPricePerHour: 10,
			ResourceConfiguration: &qcBooking.ResourceConfiguration{
				Ram:  "ram_1",
				Cpu:  "cpu_1",
				Disk: "disk_1",
			},
			Status: qcBooking.PackageStatus_PACKAGE_STATUS_ACTIVE,
		},
		{
			Id:                  "00000000-0000-0000-0000-000000000002",
			Name:                "packageName2",
			Type:                "packageType2",
			Currency:            "EUR",
			UnitIntPricePerHour: 20,
			ResourceConfiguration: &qcBooking.ResourceConfiguration{
				Ram:  "ram_2",
				Cpu:  "cpu_2",
				Disk: "disk_2",
			},
			Status: qcBooking.PackageStatus_PACKAGE_STATUS_DEACTIVATED,
		},
	}

	expected := []interface{}{
		map[string]interface{}{
			fieldID:                  "00000000-0000-0000-0000-000000000001",
			fieldName:                "packageName1",
			fieldType:                "packageType1",
			fieldCurrency:            "USD",
			fieldUnitIntPricePerHour: 10,
			fieldResourceConfiguration: map[string]interface{}{
				fieldResourceRam:  "ram_1",
				fieldResourceCpu:  "cpu_1",
				fieldResourceDisk: "disk_1",
			},
			fieldStatus: "PACKAGE_STATUS_ACTIVE",
		},
		map[string]interface{}{
			fieldID:                  "00000000-0000-0000-0000-000000000002",
			fieldName:                "packageName2",
			fieldType:                "packageType2",
			fieldCurrency:            "EUR",
			fieldUnitIntPricePerHour: 20,
			fieldResourceConfiguration: map[string]interface{}{
				fieldResourceRam:  "ram_2",
				fieldResourceCpu:  "cpu_2",
				fieldResourceDisk: "disk_2",
			},
			fieldStatus: "PACKAGE_STATUS_DEACTIVATED",
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

	expected := []interface{}{
		map[string]interface{}{
			fieldResourceRam:  "ram_3",
			fieldResourceCpu:  "cpu_3",
			fieldResourceDisk: "disk_3",
		},
	}

	assert.Equal(t, expected, flattenResourceConfiguration(rcs))
}
