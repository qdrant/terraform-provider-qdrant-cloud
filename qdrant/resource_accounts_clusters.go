package qdrant

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	_ "github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
	client := m.(ClientConfig)
	accountId := d.Get("account_id").(string)
	clusterId := d.Get("cluster_id").(string)

	// Construct the request using the helper function.
	req, diags := newQdrantCloudRequest(client, "GET", fmt.Sprintf("/accounts/%s/clusters/%s", accountId, clusterId), nil)
	if diags.HasError() {
		return diags
	}

	// Execute the request and handle the response.
	resp, diags := ExecuteRequest(client, req.WithContext(ctx))
	if diags.HasError() {
		return diags
	}

	var raw json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return diag.FromErr(fmt.Errorf("failed to decode response body: %s", err))
	}

	fmt.Println(string(raw))

	// Decode the response body into a map to match the schema structure.

	if string(raw) == "[]" {
		if err := d.Set("clusters", Cluster{}); err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	var cluster Cluster

	if err := json.Unmarshal(raw, &cluster); err != nil {
		return diag.FromErr(fmt.Errorf("failed to decode response body: %s", err))
	}

	if err := d.Set("clusters", cluster); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(ClientConfig)
	accountId := d.Get("account_id").(string)
	clusterId := d.Get("cluster_id").(string)

	// Utilize helper to construct the request according to OpenAPI
	req, diags := newQdrantCloudRequest(client, "DELETE", fmt.Sprintf("/accounts/%s/clusters/%s", accountId, clusterId), nil)
	if diags.HasError() {
		return diags
	}

	// Execute the request and handle the response
	_, diags = ExecuteRequest(client, req.WithContext(ctx))
	if diags.HasError() {
		return diags
	}

	d.SetId("")
	return nil
}

func expandNodeConfigurationIn(v *schema.Set) *NodeConfiguration {
	if len(v.List()) == 0 {
		return nil
	}

	m := v.List()[0].(map[string]interface{})
	config := NodeConfiguration{}

	if v, ok := m["package_id"]; ok {
		config.PackageID = v.(string)
	}

	return &config
}

func expandClusterConfigurationIn(v *schema.Set) *ClusterConfigurationIn {
	if len(v.List()) == 0 {
		return nil
	}

	m := v.List()[0].(map[string]interface{})
	config := ClusterConfigurationIn{}

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

func resourceClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(ClientConfig)

	name := d.Get("name")
	cloudProvider := d.Get("cloud_provider")
	cloudRegion := d.Get("cloud_region")
	accountId := d.Get("account_id")

	cluster := &Cluster{
		Name:          name.(string),
		CloudProvider: cloudProvider.(string),
		CloudRegion:   cloudRegion.(string),
		AccountID:     accountId.(string),
	}

	var ccfg *schema.Set
	if v, ok := d.GetOk("configuration"); ok {
		ccfg = v.(*schema.Set)
		configuration := expandClusterConfigurationIn(v.(*schema.Set))
		cluster.Configuration = *configuration
	}

	// Using newQdrantCloudRequest to handle HTTP request creation and error checking.
	req, diags := newQdrantCloudRequest(client, "POST", fmt.Sprintf("/accounts/%s/clusters", cluster.AccountID), cluster)
	if diags.HasError() {
		return diags
	}

	resp, diags := ExecuteRequest(client, req.WithContext(ctx))
	if diags.HasError() {
		return diags
	}

	var raw json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return diag.FromErr(fmt.Errorf("failed to decode response body: %s", err))
	}

	// Decode the response body into a map to match the schema structure.
	if strings.Contains(string(raw), "\"detail\"") {
		return diag.FromErr(fmt.Errorf("failed to create cluster: %s", raw))
	}

	if string(raw) == "[]" {
		if err := d.Set("clusters", Cluster{}); err != nil {
			return diag.FromErr(err)
		}
		return nil
	}

	// Decode JSON response into cluster struct
	var clusterOut Cluster
	if err := json.Unmarshal(raw, &clusterOut); err != nil {
		return diag.FromErr(fmt.Errorf("failed to decode response body: %s", err))
	}

	// Set properties into Terraform state
	if err := d.Set("id", clusterOut.ID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", cluster.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cloud_provider", cluster.CloudProvider); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cloud_region", cluster.CloudRegion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("configuration", ccfg); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
