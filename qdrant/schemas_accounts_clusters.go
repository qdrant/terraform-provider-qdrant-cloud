package qdrant

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// AccountsClustersSchema defines the schema for a cluster resource.
// Returns a pointer to the schema.Resource object.
func AccountsClustersSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description: "TODO",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"created_at": {
			Description: "TODO",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"owner_id": {
			Description: "TODO",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"account_id": {
			Description: "TODO",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"name": {
			Description: "TODO",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"cloud_provider": {
			Description: "TODO",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"cloud_region": {
			Description: "TODO",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"cloud_region_az": {
			Description: "TODO",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"cloud_region_setup": {
			Description: "TODO",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"private_region_id": {
			Description: "TODO",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"current_configuration_id": {
			Description: "TODO",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"encryption_key_id": {
			Description: "TODO",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"marked_for_deletion_at": {
			Description: "TODO",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"version": {
			Description: "TODO",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"url": {
			Description: "TODO",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"state": {
			Description: "TODO",
			Type:        schema.TypeMap,
			Computed:    true,
			Elem: &schema.Schema{
				Description: "TODO",
				Type:        schema.TypeString,
			},
		},
		"configuration": {
			Description: "TODO",
			Type:        schema.TypeSet,
			Required:    true,
			ForceNew:    true,
			Elem: &schema.Resource{
				Description: "TODO",
				Schema: map[string]*schema.Schema{
					"num_nodes_max": {
						Description: "TODO",
						Type:        schema.TypeInt,
						Required:    true,
					},
					"num_nodes": {
						Description: "TODO",
						Type:        schema.TypeInt,
						Required:    true,
					},
					"node_configuration": {
						Description: "TODO",
						Type:        schema.TypeSet,
						Required:    true,
						Elem: &schema.Resource{
							Schema: map[string]*schema.Schema{
								"package_id": {
									Description: "TODO",
									Type:        schema.TypeString,
									Required:    true,
								},
							},
						},
					},
				},
			},
		},
		"resources": {
			Description: "TODO",
			Type:        schema.TypeMap,
			Computed:    true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"total_extra_disk": {
			Description: "TODO",
			Type:        schema.TypeInt,
			Computed:    true,
		},
	}
}
