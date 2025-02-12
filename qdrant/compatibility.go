package qdrant

// ResourceType defines model for ResourceType.
type ResourceType string

// Defines values for ResourceType.
const (
	ResourceTypeComplimentaryDisk ResourceType = "complimentary_disk"
	ResourceTypeCpu               ResourceType = "cpu"
	ResourceTypeDisk              ResourceType = "disk"
	ResourceTypeRam               ResourceType = "ram"
	ResourceTypeSnapshot          ResourceType = "snapshot"
)

// ResourceUnit defines model for ResourceUnit.
type ResourceUnit string

// Defines values for ResourceUnit.
const (
	ResourceUnitGi ResourceUnit = "Gi"
	ResourceUnitM  ResourceUnit = "m"
)
