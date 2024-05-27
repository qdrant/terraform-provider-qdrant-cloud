package qdrant

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	qc "terraform-provider-qdrant-cloud/v1/internal/client"
)

const (
	clusterFieldTemplate = "Cluster Resource %s field"

	clusterIdentifierFieldName = "id"
	clusterCreatedAtFieldName  = "created_at"
)

// resourceAccountsClusters constructs a Terraform resource for
// managing the creation, reading, and deletion of clusters associated with an account.
func resourceAccountsClusters() *schema.Resource {
	return &schema.Resource{
		Description:   "Account Cluster Resource",
		CreateContext: resourceClusterCreate,
		ReadContext:   resourceClusterRead,
		DeleteContext: resourceClusterDelete,
		UpdateContext: nil,
		Schema: map[string]*schema.Schema{
			clusterIdentifierFieldName: {
				Description: fmt.Sprintf(clusterFieldTemplate, "Identifier of the cluster"),
				Type:        schema.TypeString,
				Computed:    true,
			},
			clusterCreatedAtFieldName: {
				Description: fmt.Sprintf(clusterFieldTemplate, "Timestamp then the cluster is created"),
				Type:        schema.TypeString,
				Computed:    true,
			},
			"owner_id": {
				Description: fmt.Sprintf(clusterFieldTemplate, "Identifier of the owner"),
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"account_id": { // If set here, overrides account ID in provider
				Description: fmt.Sprintf(clusterFieldTemplate, "Identifier of the account"),
				Type:        schema.TypeString,
				Optional:    true,
			},
			"name": {
				Description: fmt.Sprintf(clusterFieldTemplate, "Name of the cluster"),
				Type:        schema.TypeString,
				Required:    true,
			},
			"cloud_provider": {
				Description: fmt.Sprintf(clusterFieldTemplate, "Cloud provider of the cluster"),
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"cloud_region": {
				Description: fmt.Sprintf(clusterFieldTemplate, "Cloud region of the cluster"),
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"cloud_region_az": {
				Description: fmt.Sprintf(clusterFieldTemplate, "Cloud region availability zone of the cluster"),
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"cloud_region_setup": {
				Description: fmt.Sprintf(clusterFieldTemplate, "Cloud region setup of the cluster"),
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"private_region_id": {
				Description: fmt.Sprintf(clusterFieldTemplate, "Identifier of the Private Region"),
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
					Type: schema.TypeString,
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
							ForceNew:    true,
						},
						"num_nodes": {
							Description: "TODO",
							Type:        schema.TypeInt,
							Required:    true,
							ForceNew:    true,
						},
						"node_configuration": {
							Description: "TODO",
							Type:        schema.TypeSet,
							Required:    true,
							ForceNew:    true,
							Elem: &schema.Resource{
								Description: "TODO",
								Schema: map[string]*schema.Schema{
									"package_id": {
										Description: "TODO",
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
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
					Description: "TODO",
					Type:        schema.TypeString,
				},
			},
			"total_extra_disk": {
				Description: "TODO",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

// resourceClusterRead reads the specific cluster's data from the API.
func resourceClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get an authenticated client
	apiClient, diagnostics := getClient(m)
	if diagnostics.HasError() {
		return diagnostics
	}

	accountID, err := uuid.Parse(d.Get("account_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	clusterID, err := uuid.Parse(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
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
	// Get an authenticated client
	apiClient, diagnostics := getClient(m)
	if diagnostics.HasError() {
		return diagnostics
	}
	accountID, err := uuid.Parse(d.Get("account_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	clusterID, err := uuid.Parse(d.Get("id").(string))
	if err != nil {
		return diag.FromErr(err)
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
	// Get an authenticated client
	apiClient, diagnostics := getClient(m)
	if diagnostics.HasError() {
		return diagnostics
	}
	accountID := d.Get("account_id").(string)

	name := d.Get("name")
	cloudProvider := d.Get("cloud_provider")
	cloudRegion := d.Get("cloud_region")

	cluster := qc.ClusterIn{
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

	accID, err := uuid.Parse(accountID)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := apiClient.CreateClusterWithResponse(ctx, accID, cluster)
	if err != nil {
		return diag.FromErr(err)
	}

	if resp.JSON422 != nil {
		return diag.FromErr(fmt.Errorf("error creating cluster: %v", resp.JSON422))
	}

	clusterOut := resp.JSON200
	if clusterOut == nil {
		return diag.FromErr(fmt.Errorf("error creating cluster: no cluter returned"))
	}

	// Set properties into Terraform state
	if err := d.Set("id", clusterOut.Id); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", clusterOut.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cloud_provider", clusterOut.CloudProvider); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cloud_region", clusterOut.CloudRegion); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("configuration", ccfg); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
