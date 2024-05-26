package qdrant

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	qc "terraform-provider-qdrant-cloud/v1/internal/client"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceAccountsClusters constructs a Terraform resource for
// managing the creation, reading, and deletion of clusters associated with an account.
func resourceAccountsClusters() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterCreate,
		ReadContext:   resourceClusterRead,
		DeleteContext: resourceClusterDelete,
		UpdateContext: nil,
		Schema: map[string]*schema.Schema{
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
							ForceNew: true,
						},
						"num_nodes": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
						},
						"node_configuration": {
							Type:     schema.TypeSet,
							Required: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"package_id": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
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
		},
	}
}

// resourceClusterRead reads the specific cluster's data from the API.
func resourceClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountID, err := uuid.Parse(d.Get("account_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	clusterID, err := uuid.Parse(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	apiClient, err, diagnostics, done := GetClient(m)
	if done {
		return diagnostics
	}

	resp, err := apiClient.GetClusterWithResponse(ctx, accountID, clusterID)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("error reading cluster: %v", resp.JSON422))
	}

	if err := d.Set("clusters", resp.JSON200); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(time.Now().Format(time.RFC3339))
	return nil
}

// resourceClusterDelete performs a delete operation to remove a cluster associated with an account.
func resourceClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountID, err := uuid.Parse(d.Get("account_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	clusterID, err := uuid.Parse(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	apiClient, err, diagnostics, done := GetClient(m)
	if done {
		return diagnostics
	}

	deleteBackups := true
	params := qc.DeleteClusterParams{
		DeleteBackups: &deleteBackups,
	}
	resp, err := apiClient.DeleteClusterWithResponse(ctx, accountID, clusterID, &params)
	if err != nil {
		return diag.FromErr(err)
	}
	if resp.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("error deleting cluster: %v", resp.JSON422))
	}

	d.SetId("")
	return nil
}

func expandClusterConfigurationIn(v *schema.Set) *qc.ClusterConfigurationIn {
	if len(v.List()) == 0 {
		return nil
	}

	m := v.List()[0].(map[string]interface{})
	config := qc.ClusterConfigurationIn{}

	if v, ok := m["num_nodes_max"]; ok {
		config.NumNodesMax = v.(int)
	}
	if v, ok := m["num_nodes"]; ok {
		config.NumNodes = v.(int)
	}

	if v, ok := m["node_configuration"]; ok {
		nodeConfig := expandNodeConfigurationIn(v.(*schema.Set))
		if nodeConfig != nil {
			config.NodeConfiguration = *nodeConfig
		}
	}

	return &config
}

func expandNodeConfigurationIn(v *schema.Set) *qc.NodeConfiguration {
	if len(v.List()) == 0 {
		return nil
	}

	m := v.List()[0].(map[string]interface{})
	config := qc.NodeConfiguration{}

	if v, ok := m["package_id"]; ok {
		config.PackageId = v.(string)
	}

	return &config
}

func resourceClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	accountID := d.Get("account_id").(string)

	name := d.Get("name")
	cloudProvider := d.Get("cloud_provider")
	cloudRegion := d.Get("cloud_region")

	cluster := &qc.ClusterIn{
		Name:          name.(string),
		CloudProvider: qc.ClusterInCloudProvider(cloudProvider.(string)),
		CloudRegion:   qc.ClusterInCloudRegion(cloudRegion.(string)),
		AccountId:     &accountID,
	}

	var ccfg *schema.Set
	if v, ok := d.GetOk("configuration"); ok {
		ccfg = v.(*schema.Set)
		configuration := expandClusterConfigurationIn(v.(*schema.Set))
		cluster.Configuration = *configuration
	}

	requestBody, err := json.Marshal(cluster)
	if err != nil {
		return diag.FromErr(err)
	}

	body := bytes.NewReader(requestBody)

	apiClient, err, diagnostics, done := GetClient(m)
	if done {
		return diagnostics
	}

	accID, err := uuid.Parse(accountID)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := apiClient.CreateClusterWithBodyWithResponse(ctx, accID, "application/json", body)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("error creating cluster: %v", resp.JSON422))
	}

	// Set properties into Terraform state
	if err := d.Set("id", resp.JSON200.Id); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", resp.JSON200.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cloud_provider", resp.JSON200.CloudProvider); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cloud_region", resp.JSON200.CloudRegion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("configuration", ccfg); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
