package qdrant

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AccountsClustersSchema defines the schema for a cluster resource.
// Returns a pointer to the schema.Resource object.
func AccountsClustersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"owner_id": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"account_id": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"cloud_provider": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"cloud_region": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"cloud_region_az": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"cloud_region_setup": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"private_region_id": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"current_configuration_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"encryption_key_id": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"marked_for_deletion_at": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"version": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"url": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"state": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"configuration": {
			Type:     schema.TypeSet,
			Required: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"num_nodes_max": {
						Type:     schema.TypeInt,
						Required: true,
					},
					"num_nodes": {
						Type:     schema.TypeInt,
						Required: true,
					},
					"node_configuration": {
						Type:     schema.TypeSet,
						Required: true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"package_id": {
									Type:     schema.TypeString,
									Required: true,
								},
							},
						},
					},
				},
			},
		},
		"resources": {
			Type:     schema.TypeMap,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"total_extra_disk": {
			Type:     schema.TypeInt,
			Computed: true,
		},
	}
}

type Cluster struct {
	AccountID              string                 `json:"account_id"`
	Name                   string                 `json:"name"`
	CloudProvider          string                 `json:"cloud_provider"`
	CloudRegion            string                 `json:"cloud_region"`
	Configuration          ClusterConfigurationIn `json:"configuration"`
	ID                     string                 `json:"id,omitempty"`
	CreatedAt              string                 `json:"created_at,omitempty"`
	OwnerID                string                 `json:"owner_id,omitempty"`
	CloudRegionAZ          string                 `json:"cloud_region_az,omitempty"`
	CloudRegionSetup       string                 `json:"cloud_region_setup,omitempty"`
	PrivateRegionID        string                 `json:"private_region_id,omitempty"`
	CurrentConfigurationID string                 `json:"current_configuration_id,omitempty"`
	EncryptionKeyID        string                 `json:"encryption_key_id,omitempty"`
	MarkedForDeletionAt    string                 `json:"marked_for_deletion_at,omitempty"`
	Version                string                 `json:"version,omitempty"`
	URL                    string                 `json:"url,omitempty"`
	State                  interface{}            `json:"state,omitempty"`
	Resources              interface{}            `json:"resources,omitempty"`
	TotalExtraDisk         int                    `json:"total_extra_disk,omitempty"`
	Schedule               *ScheduleIn            `json:"schedule,omitempty"`
	EncryptionConfig       *EncryptionConfigIn    `json:"encryption_config,omitempty"`
}

// ClusterConfigurationIn represents the input configuration for a cluster.
type ClusterConfigurationIn struct {
	NumNodes              int                     `json:"num_nodes"`
	NumNodesMax           int                     `json:"num_nodes_max"`
	NodeConfiguration     NodeConfiguration       `json:"node_configuration"`
	QdrantConfiguration   *map[string]interface{} `json:"qdrant_configuration,omitempty"`
	NodeSelector          *map[string]string      `json:"node_selector,omitempty"`
	Tolerations           *[]map[string]string    `json:"tolerations,omitempty"`
	ClusterAnnotations    *map[string]interface{} `json:"cluster_annotations,omitempty"`
	AllowedIPSourceRanges *[]string               `json:"allowed_ip_source_ranges,omitempty"`
}

// NodeConfiguration defines the configuration for a node within the cluster.
type NodeConfiguration struct {
	PackageID              string                   `json:"package_id"`
	Package                *PackageOut              `json:"package,omitempty"`
	ResourceConfigurations *[]ResourceConfiguration `json:"resource_configurations,omitempty"`
}

// PackageOut represents an output package information.
type PackageOut struct {
	ID                    *string                 `json:"id,omitempty"`
	ResourceConfiguration []ResourceConfiguration `json:"resource_configuration"`
	Name                  string                  `json:"name"`
	Status                BookingStatus           `json:"status"`
	Currency              Currency                `json:"currency"`
	UnitIntPricePerHour   *int                    `json:"unit_int_price_per_hour,omitempty"`
	UnitIntPricePerDay    *int                    `json:"unit_int_price_per_day,omitempty"`
	UnitIntPricePerMonth  *int                    `json:"unit_int_price_per_month,omitempty"`
	UnitIntPricePerYear   *int                    `json:"unit_int_price_per_year,omitempty"`
	RegionalMappingID     *string                 `json:"regional_mapping_id,omitempty"`
}

// ResourceConfiguration holds the resource configurations for a node.
type ResourceConfiguration struct {
	ResourceOptionID string             `json:"resource_option_id"`
	ResourceOption   *ResourceOptionOut `json:"resource_option,omitempty"`
	Amount           int                `json:"amount"`
}

// ResourceOptionOut represents the details of a resource option.
type ResourceOptionOut struct {
	ID                   string        `json:"id"`
	ResourceType         ResourceType  `json:"resource_type"`
	Status               BookingStatus `json:"status"`
	Name                 *string       `json:"name,omitempty"`
	ResourceUnit         string        `json:"resource_unit"`
	Currency             Currency      `json:"currency"`
	UnitIntPricePerHour  *int          `json:"unit_int_price_per_hour,omitempty"`
	UnitIntPricePerDay   *int          `json:"unit_int_price_per_day,omitempty"`
	UnitIntPricePerMonth *int          `json:"unit_int_price_per_month,omitempty"`
	UnitIntPricePerYear  *int          `json:"unit_int_price_per_year,omitempty"`
}

// ScheduleIn represents scheduling information for backups and other periodic tasks.
type ScheduleIn struct {
	CreatorUserID       *string        `json:"creator_user_id,omitempty"`
	AccountID           *string        `json:"account_id,omitempty"`
	Cron                string         `json:"cron"`
	Retention           int            `json:"retention"`
	PrivateRegionID     *string        `json:"private_region_id,omitempty"`
	MarkedForDeletionAt *time.Time     `json:"marked_for_deletion_at,omitempty"`
	Status              *ScheduleState `json:"status,omitempty"`
}

// EncryptionConfigIn defines the encryption settings for a cluster.
type EncryptionConfigIn struct {
	AWSEncryptionConfig *AWSEncryptionConfig `json:"aws_encryption_config,omitempty"`
}

// AWSEncryptionConfig contains AWS specific encryption configuration details.
type AWSEncryptionConfig struct {
	Managed         bool    `json:"managed,omitempty"`
	EncryptionKeyID *string `json:"encryption_key_id,omitempty"`
}

// BookingStatus represents the status of a booking.
type BookingStatus int

// Currency represents a currency type.
type Currency string

// ResourceType represents the type of resource.
type ResourceType string

// ScheduleState represents the state of a schedule.
type ScheduleState string
