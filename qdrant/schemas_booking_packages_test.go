package qdrant

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	qcBooking "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/booking/v1"
	commonv1 "github.com/qdrant/qdrant-cloud-public-api/gen/go/qdrant/cloud/common/v1"
)

func TestFlattenPackages(t *testing.T) {
	gpu3 := "gpu_3"
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
			Tier:   qcBooking.PackageTier_PACKAGE_TIER_STANDARD,
			AvailableAdditionalResources: &qcBooking.AvailableAdditionalResources{
				DiskPricePerHour: 5,
			},
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
			Tier:   qcBooking.PackageTier_PACKAGE_TIER_PREMIUM,
		},
		{
			Id:                  "00000000-0000-0000-0000-000000000003",
			Name:                "packageName3",
			Type:                "packageType3",
			Currency:            "JPY",
			UnitIntPricePerHour: 30,
			ResourceConfiguration: &qcBooking.ResourceConfiguration{
				Ram:  "ram_3",
				Cpu:  "cpu_3",
				Disk: "disk_3",
				Gpu:  &gpu3,
			},
			Status:  qcBooking.PackageStatus_PACKAGE_STATUS_ACTIVE,
			Tier:    qcBooking.PackageTier_PACKAGE_TIER_PREMIUM,
			MultiAz: true,
			AvailableStorageTierConfigurations: []*qcBooking.AvailableStoragePerformanceTierConfigurations{
				{StorageTierType: commonv1.StorageTierType_STORAGE_TIER_TYPE_PERFORMANCE, PricePerHour: 15},
				{StorageTierType: commonv1.StorageTierType_STORAGE_TIER_TYPE_COST_OPTIMISED, PricePerHour: 0},
			},
		},
	}

	expected := []interface{}{
		map[string]interface{}{
			fieldID:                  "00000000-0000-0000-0000-000000000001",
			fieldName:                "packageName1",
			fieldType:                "packageType1",
			fieldCurrency:            "USD",
			fieldUnitIntPricePerHour: 10,
			fieldResourceConfiguration: []interface{}{
				map[string]interface{}{
					fieldResourceRam:  "ram_1",
					fieldResourceCpu:  "cpu_1",
					fieldResourceDisk: "disk_1",
					fieldResourceGpu:  "",
				},
			},
			fieldStatus: "PACKAGE_STATUS_ACTIVE",
			fieldTier:   "PACKAGE_TIER_STANDARD",
			fieldAvailableAddResources: []interface{}{map[string]interface{}{
				"disk_price_per_hour": 5,
			},
			},
			fieldMultiAz: false,
		},
		map[string]interface{}{
			fieldID:                  "00000000-0000-0000-0000-000000000002",
			fieldName:                "packageName2",
			fieldType:                "packageType2",
			fieldCurrency:            "EUR",
			fieldUnitIntPricePerHour: 20,
			fieldResourceConfiguration: []interface{}{
				map[string]interface{}{
					fieldResourceRam:  "ram_2",
					fieldResourceCpu:  "cpu_2",
					fieldResourceDisk: "disk_2",
					fieldResourceGpu:  "",
				},
			},
			fieldStatus:  "PACKAGE_STATUS_DEACTIVATED",
			fieldTier:    "PACKAGE_TIER_PREMIUM",
			fieldMultiAz: false,
		},
		map[string]interface{}{
			fieldID:                  "00000000-0000-0000-0000-000000000003",
			fieldName:                "packageName3",
			fieldType:                "packageType3",
			fieldCurrency:            "JPY",
			fieldUnitIntPricePerHour: 30,
			fieldResourceConfiguration: []interface{}{
				map[string]interface{}{
					fieldResourceRam:  "ram_3",
					fieldResourceCpu:  "cpu_3",
					fieldResourceDisk: "disk_3",
					fieldResourceGpu:  "gpu_3",
				},
			},
			fieldStatus:  "PACKAGE_STATUS_ACTIVE",
			fieldTier:    "PACKAGE_TIER_PREMIUM",
			fieldMultiAz: true,
			fieldAvailableStorageTierConfigurations: []interface{}{
				map[string]interface{}{fieldStorageTierType: "STORAGE_TIER_TYPE_PERFORMANCE", fieldPricePerHour: 15},
				map[string]interface{}{fieldStorageTierType: "STORAGE_TIER_TYPE_COST_OPTIMISED", fieldPricePerHour: 0},
			},
		},
	}

	got := flattenPackages(packages)
	require.Len(t, got, len(expected))
	assert.Equal(t, expected, got)
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
			fieldResourceGpu:  "",
		},
	}

	assert.Equal(t, expected, flattenResourceConfiguration(rcs))
}

func TestFlattenResourceConfiguration_WithGpu(t *testing.T) {
	gpu := "gpu_4"
	rcs := &qcBooking.ResourceConfiguration{
		Ram:  "ram_4",
		Cpu:  "cpu_4",
		Disk: "disk_4",
		Gpu:  &gpu,
	}

	expected := []interface{}{
		map[string]interface{}{
			fieldResourceRam:  "ram_4",
			fieldResourceCpu:  "cpu_4",
			fieldResourceDisk: "disk_4",
			fieldResourceGpu:  "gpu_4",
		},
	}

	assert.Equal(t, expected, flattenResourceConfiguration(rcs))
}

func TestFlattenAvailableStorageTierConfigurations(t *testing.T) {
	configs := []*qcBooking.AvailableStoragePerformanceTierConfigurations{
		{
			StorageTierType: commonv1.StorageTierType_STORAGE_TIER_TYPE_PERFORMANCE,
			PricePerHour:    100,
		},
		{
			StorageTierType: commonv1.StorageTierType_STORAGE_TIER_TYPE_COST_OPTIMISED,
			PricePerHour:    50,
		},
	}

	expected := []interface{}{
		map[string]interface{}{fieldStorageTierType: "STORAGE_TIER_TYPE_PERFORMANCE", fieldPricePerHour: 100},
		map[string]interface{}{fieldStorageTierType: "STORAGE_TIER_TYPE_COST_OPTIMISED", fieldPricePerHour: 50},
	}

	assert.ElementsMatch(t, expected, flattenAvailableStorageTierConfigurations(configs))
}

func TestFlattenAvailableStorageTierConfigurations_Empty(t *testing.T) {
	assert.Nil(t, flattenAvailableStorageTierConfigurations(nil))
	assert.Nil(t, flattenAvailableStorageTierConfigurations([]*qcBooking.AvailableStoragePerformanceTierConfigurations{}))
}
